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
    "github.com/auroraride/aurservd/internal/ent/commission"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/orderrefund"
    "github.com/auroraride/aurservd/internal/payment"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
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
    if NewRiderOrder().FindNotActived(s.rider.ID) != nil {
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
    no := tools.NewUnique().NewSonyflakeID()
    result.OutTradeNo = no
    // 生成订单字段
    price := plan.Price
    // DEBUG 模式支付一分钱
    mode := ar.Config.App.Mode
    if mode == "debug" {
        price = 0.01
        if deposit > 0 {
            deposit = 0.01
        }
    }
    total, _ := decimal.NewFromFloat(price).Add(decimal.NewFromFloat(deposit)).Float64()
    prepay := &model.PaymentCache{
        CacheType: model.PaymentCacheTypePlan,
        Plan: &model.PaymentPlan{
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
        },
    }
    // 订单缓存
    err := cache.Set(s.ctx, no, prepay, 20*time.Minute).Err()
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

// Refund 退款
func (s *orderService) Refund(riderID uint64, req *model.OrderRefundReq) {
    if (req.Deposit == nil && req.OrderID == nil) || (req.Deposit != nil && req.OrderID != nil) {
        snag.Panic("退款参数错误")
    }
    q := s.orm.QueryNotDeleted().
        Where(
            order.RiderID(riderID),
            order.Status(model.OrderStatusPaid),
        )
    if req.OrderID != nil {
        q.Where(order.ID(*req.OrderID))
    }
    if req.Deposit != nil && *req.Deposit {
        // 查找是否满足退押金条件
        if exist, _ := s.orm.QueryNotDeleted().Where(
            order.RiderID(riderID),
            order.EndAtIsNil(),
            order.StatusIn(model.OrderStatusPaid, model.OrderStatusRefundPending),
            order.TypeIn(model.OrderRiderPlan...),
        ).Exist(s.ctx); exist {
            snag.Panic("当前无法退款")
        }
        q.Where(order.Type(model.OrderTypeDeposit))
    }

    o, _ := q.Only(s.ctx)

    if o == nil {
        snag.Panic("未找到有效订单")
    }
    no := tools.NewUnique().NewSonyflakeID()
    refund := &model.PaymentRefund{
        OrderID:      o.ID,
        TradeNo:      o.TradeNo,
        Total:        o.Total,
        RefundAmount: o.Amount,
        OutRefundNo:  no,
        Reason:       "用户申请",
    }
    // 缓存
    cache.Set(s.ctx, no, &model.PaymentCache{CacheType: model.PaymentCacheTypeRefund, Refund: refund}, 0)

    switch o.Payway {
    case model.OrderPaywayAlipay:
        payment.NewAlipay().Refund(refund)
        break
    case model.OrderPaywayWechat:
        payment.NewWechat().Refund(refund)
        break
    }

    if refund.Request {
        o.Update().SetStatus(model.OrderStatusRefundPending).SaveX(s.ctx)
        // 保存订单
        ar.Ent.OrderRefund.Create().
            SetOrder(o).
            SetAmount(refund.RefundAmount).
            SetOutRefundNo(no).
            SetReason(refund.Reason).
            SetStatus(model.OrderRefundStatusPending).
            SaveX(s.ctx)

        if refund.Success {
            s.RefundSuccess(refund)
        }
    }
}

// DoPayment 处理支付
func (s *orderService) DoPayment(pc *model.PaymentCache) {
    if pc == nil {
        return
    }
    switch pc.CacheType {
    case model.PaymentCacheTypePlan:
        s.OrderPaid(pc.Plan)
        break
    case model.PaymentCacheTypeRefund:
        s.RefundSuccess(pc.Refund)
    }
}

// OrderPaid 订单成功支付
// TODO 业绩提成逻辑
// TODO 2续签 3重签 4更改电池 5救援 6滞纳金 逻辑
func (s *orderService) OrderPaid(trade *model.PaymentPlan) {
    // 查询订单是否已存在
    if exists, err := ar.Ent.Order.Query().Where(order.OutTradeNo(trade.OutTradeNo)).Exist(s.ctx); err == nil && exists {
        return
    }
    o, err := ar.Ent.Order.Create().
        SetPayway(trade.Payway).
        SetPlanID(trade.Plan.ID).
        SetRiderID(trade.RiderID).
        SetAmount(decimal.NewFromFloat(trade.Amount).Sub(decimal.NewFromFloat(trade.Deposit)).InexactFloat64()).
        SetTotal(trade.Amount).
        SetOutTradeNo(trade.OutTradeNo).
        SetTradeNo(trade.TradeNo).
        SetStatus(model.OrderStatusPaid).
        SetNillablePlanDetail(trade.Plan).
        SetType(trade.OrderType).
        SetCityID(trade.CityID).
        SetDays(trade.Plan.Days).
        Save(s.ctx)
    if err != nil {
        log.Errorf("[ORDER PAID ERROR]: %s, Trade: %#v", err.Error(), trade)
        return
    }
    // 如果有押金, 创建押金订单
    if trade.Deposit > 0 {
        _, err = ar.Ent.Order.Create().
            SetStatus(model.OrderStatusPaid).
            SetPayway(trade.Payway).
            SetType(model.OrderTypeDeposit).
            SetOutTradeNo(trade.OutTradeNo).
            SetTradeNo(trade.TradeNo).
            SetAmount(trade.Deposit).
            SetTotal(trade.Amount).
            SetCityID(trade.CityID).
            SetRiderID(trade.RiderID).
            SetParentID(o.ID).
            Save(s.ctx)
        if err != nil {
            log.Errorf("[DEPOSIT PAID ERROR]: %s, Trade: %#v", err.Error(), trade)
            return
        }
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
    cache.Del(context.Background(), trade.OutTradeNo)
}

// RefundSuccess 成功退款
func (s *orderService) RefundSuccess(req *model.PaymentRefund) {
    ctx := context.Background()
    tx, _ := ar.Ent.Tx(ctx)
    tx.Order.UpdateOneID(req.OrderID).
        SetStatus(model.OrderStatusRefundSuccess).
        SetEndAt(req.Time).
        SaveX(ctx)

    tx.OrderRefund.Update().
        Where(orderrefund.OutRefundNo(req.OutRefundNo)).
        SetStatus(model.OrderRefundStatusSuccess).
        SetRefundAt(req.Time).
        SaveX(ctx)

    // 删除提成订单
    tx.Commission.Delete().Where(commission.OrderID(req.OrderID)).ExecX(ctx)

    _ = tx.Commit()

    // 删除缓存
    cache.Del(ctx, req.OutRefundNo)
}
