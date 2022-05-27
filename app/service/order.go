// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/payment"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
    "github.com/shopspring/decimal"
    log "github.com/sirupsen/logrus"
    "time"
)

type orderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.OrderClient
}

func NewOrder() *orderService {
    return &orderService{
        ctx: context.Background(),
        orm: ar.Ent.Order,
    }
}

func NewOrderWithRider(rider *ent.Rider) *orderService {
    s := NewOrder()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewOrderWithModifier(m *model.Modifier) *orderService {
    s := NewOrder()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// Create 创建订单
func (s *orderService) Create(req *model.OrderCreateReq) (result model.OrderCreateRes) {
    // 查询骑手是否签约过
    if !NewContract().Effective(s.rider) {
        snag.Panic("请先签约")
    }
    // 查询是否存在未使用的骑士卡
    if s.RiderNotActivedExists(s.rider.ID) {
        snag.Panic("当前有未使用的骑士卡")
    }
    // 查询套餐是否存在
    plan := NewPlan().QueryEffectiveWithID(req.PlanID)
    cp := new(model.PlanItem)
    _ = copier.CopyWithOption(cp, plan, copier.Option{Converters: []copier.TypeConverter{
        {
            SrcType: time.Time{},
            DstType: copier.String,
            Fn: func(src interface{}) (interface{}, error) {
                t, ok := src.(time.Time)
                if !ok {
                    return "", nil
                }
                return t.Format(carbon.DateLayout), nil
            },
        },
    }})
    // 判定用户是否需要缴纳押金
    deposit := NewRider().Deposit(s.rider)
    no := shortuuid.New()
    result.OutTradeNo = no
    // 生成订单字段
    price := plan.Price
    // DEBUG 模式支付一分钱
    mode := ar.Config.App.Mode
    if mode == "debug" {
        deposit = 0.01
        price = 0.01
    }
    total, _ := decimal.NewFromFloat(price).Add(decimal.NewFromFloat(deposit)).Float64()
    prepay := &model.OrderCache{
        OrderType:  req.OrderType,
        OutTradeNo: no,
        RiderID:    s.rider.ID,
        Plan:       cp,
        Name:       "购买" + plan.Name,
        Amount:     total,
        Deposit:    deposit,
        Payway:     req.Payway,
        Expire:     time.Now().Add(10 * time.Minute),
        CityID:     req.CityID,
    }
    // 订单缓存
    err := cache.Set(s.ctx, "ORDER_"+no, prepay, 20*time.Minute).Err()
    if err != nil {
        log.Error(err)
        snag.Panic("订单创建失败")
    }
    var str string
    switch req.Payway {
    case model.OrderPaywayAlipay:
        // 使用支付宝支付
        str, err = payment.NewAlipay().AppPay(prepay)
        if err != nil {
            log.Error(err)
            snag.Panic("支付宝支付请求失败")
        }
        result.Prepay = str
        break
    case model.OrderPaywayWechat:
        // 使用微信支付
        str, err = payment.NewWechat().AppPay(prepay)
        if err != nil {
            log.Error(err)
            snag.Panic("微信支付请求失败")
        }
        result.Prepay = str
        break
    default:
        snag.Panic("支付方式错误")
        break
    }
    return
}

// OrderPaid 订单成功支付
// TODO 业绩提成逻辑
// TODO 2续签 3重签 4更改电池 5救援 6滞纳金 逻辑
func (s *orderService) OrderPaid(trade *model.OrderCache) {
    o, err := ar.Ent.Order.Create().
        SetPayway(trade.Payway).
        SetPlanID(trade.Plan.ID).
        SetRiderID(trade.RiderID).
        SetAmount(decimal.NewFromFloat(trade.Amount).Sub(decimal.NewFromFloat(trade.Deposit)).InexactFloat64()).
        SetOutTradeNo(trade.OutTradeNo).
        SetTradeNo(trade.TradeNo).
        SetStatus(model.OrderStatusPaid).
        SetNillablePlanDetail(trade.Plan).
        SetType(trade.OrderType).
        SetCityID(trade.CityID).
        Save(s.ctx)
    // 如果有押金, 创建押金订单
    if trade.Deposit > 0 {
        _, _ = ar.Ent.Order.Create().
            SetPayway(trade.Payway).
            SetRiderID(trade.RiderID).
            SetAmount(trade.Deposit).
            SetOutTradeNo(model.SettingDeposit + "-" + trade.OutTradeNo).
            SetTradeNo(model.SettingDeposit + "-" + trade.TradeNo).
            SetStatus(model.OrderStatusPaid).
            SetType(model.OrderTypeDeposit).
            Save(s.ctx)
    }
    if err != nil {
        log.Errorf("[ORDER PAID ERROR]: %s, Trade: %#v", err.Error(), trade)
        return
    }
    // 当新签和重签的时候有提成
    if trade.OrderType == model.OrderTypeNewPlan || trade.OrderType == model.OrderTypeRenewal {
        // 创建提成
        _, err = ar.Ent.Commission.Create().SetOrder(o).SetAmount(trade.Plan.Commission).SetStatus(model.CommissionStatusPending).Save(s.ctx)
        if err != nil {
            log.Errorf("订单提成创建失败: %d: %s", o.ID, err.Error())
        }
    }

    // 删除缓存
    cache.Del(context.Background(), "ORDER_"+trade.OutTradeNo)
}

// RiderNotActivedExists 查找用户未激活的新订单
// 已支付 && 开始时间为空 && 新签 && 重签
func (s *orderService) RiderNotActivedExists(riderID uint64) bool {
    o, _ := s.orm.QueryNotDeleted().Where(
        order.RiderID(riderID),
        order.Status(model.OrderStatusPaid),
        order.TypeIn(model.OrderTypeNewPlan, model.OrderTypeRenewal),
        order.StartIsNil(),
    ).Exist(s.ctx)
    return o
}

// RiderNotActived 获取未激活骑士卡详情
func (s *orderService) RiderNotActived(riderID uint64) *model.OrderNotActived {
    o, err := s.orm.QueryNotDeleted().
        Where(
            order.RiderID(riderID),
            order.Status(model.OrderStatusPaid),
            order.TypeIn(model.OrderTypeNewPlan, model.OrderTypeRenewal),
            order.StartIsNil(),
        ).
        WithCity().
        WithPlan(func(pq *ent.PlanQuery) {
            pq.WithPms()
        }).
        WithChildren(func(cq *ent.OrderQuery) {
            cq.Where(order.Type(model.OrderTypeDeposit))
        }).
        Only(s.ctx)
    if err != nil {
        return nil
    }
    item := &model.OrderNotActived{
        ID:     o.ID,
        Amount: o.Amount,
        Total:  o.Amount,
        Payway: o.Payway,
        Plan:   o.PlanDetail,
        City: model.City{
            ID:   o.Edges.City.ID,
            Name: o.Edges.City.Name,
        },
        Time: o.CreatedAt.Format(carbon.DateLayout),
    }
    for _, pm := range o.Edges.Plan.Edges.Pms {
        item.Models = append(item.Models, model.BatteryModel{
            ID:       pm.ID,
            Voltage:  pm.Voltage,
            Capacity: pm.Capacity,
        })
    }
    if len(o.Edges.Children) > 0 {
        item.Deposit = o.Edges.Children[0].Amount
    }
    item.Total += item.Deposit
    return item
}
