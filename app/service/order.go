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
    "github.com/auroraride/aurservd/internal/payment"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
    log "github.com/sirupsen/logrus"
    "time"
)

type orderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewOrder() *orderService {
    return &orderService{
        ctx: context.Background(),
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
    // 查询套餐是否存在
    plan := NewPlan().QueryEffective(req.PlanID)
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
    no := shortuuid.New()
    result.OutTradeNo = no
    // 生成订单字段
    price := plan.Price
    // DEBUG 模式支付一分钱
    mode := ar.Config.App.Mode
    if mode == "debug" {
        price = 0.01
    }
    prepay := &model.OrderCache{
        OrderType:  req.OrderType,
        OutTradeNo: no,
        RiderID:    s.rider.ID,
        Plan:       cp,
        Name:       "购买" + plan.Name,
        Amount:     price,
        Payway:     req.Payway,
        Expire:     time.Now().Add(10 * time.Minute),
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
        SetAmount(trade.Amount).
        SetOutTradeNo(trade.OutTradeNo).
        SetTradeNo(trade.TradeNo).
        SetStatus(model.OrderStatusPaid).
        SetNillablePlanDetail(trade.Plan).
        SetType(trade.OrderType).
        Save(s.ctx)
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
