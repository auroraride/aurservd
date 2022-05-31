// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "encoding/json"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
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

// PreconditionNewly 新签订单前置条件判定
// 返回值为最终判定的订单类型(有可能是重签)
func (s *orderService) PreconditionNewly(sub *model.Subscribe) (state uint) {
    if sub == nil {
        return model.OrderTypeNewly
    }
    switch sub.Status {
    case model.SubscribeStatusInactive:
        snag.Panic("当前有未激活的骑士卡")
        break
    case model.SubscribeStatusCanceled:
        return model.OrderTypeNewly
    case model.SubscribeStatusUnSubscribed:
        return model.OrderTypeAgain
    default:
        snag.Panic("不满足新签条件")
        break
    }
    return model.OrderTypeNewly
}

// PreconditionRenewal 是否满足续签条件
func (s *orderService) PreconditionRenewal(sub *model.Subscribe) {
    if sub == nil ||
        sub.Status == model.SubscribeStatusInactive ||
        sub.Status == model.SubscribeStatusUnSubscribed ||
        sub.Status == model.SubscribeStatusCanceled {
        snag.Panic("当前无使用中的骑士卡")
    }
    if sub.Remaining < 0 {
        snag.Panic("请先缴纳逾期费用")
    }
}

// Create 创建订单
// TODO 需要做更改电池逻辑, 退款 + 推送订单
func (s *orderService) Create(req *model.OrderCreateReq) (result model.OrderCreateRes) {
    // 查询骑手是否签约过
    if !NewContract().Effective(s.rider) {
        snag.Panic("请先签约")
    }
    sub := NewSubscribe().Recent(s.rider.ID)
    // 判定类型条件
    var subID, orderID *uint64
    otype := req.OrderType
    switch otype {
    case model.OrderTypeNewly:
        // 新签判定
        otype = s.PreconditionNewly(sub)
        break
    case model.OrderTypeRenewal:
        // 续签判定
        s.PreconditionRenewal(sub)
        subID = tools.NewPointer().UInt64(sub.ID)
        orderID = tools.NewPointer().UInt64(sub.Order.ID)
        break
    default:
        snag.Panic("未知的支付请求")
        break
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

    // 距离上次订阅过去的时间
    var pastDays int
    if otype == model.OrderTypeAgain {
        pastDays = tools.NewTime().SubDaysToNowString(sub.EndAt)
    }

    // 判定用户是否需要缴纳押金
    deposit := NewRider().Deposit(s.rider)
    no := tools.NewUnique().NewSN28()
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
        Subscribe: &model.PaymentSubscribe{
            CityID:      req.CityID,
            OrderType:   otype,
            OutTradeNo:  no,
            RiderID:     s.rider.ID,
            Name:        "购买" + plan.Name,
            Amount:      total,
            Payway:      req.Payway,
            Expire:      time.Now().Add(10 * time.Minute),
            PlanID:      cp.ID,
            Deposit:     deposit,
            PastDays:    pastDays,
            Commission:  plan.Commission,
            Voltage:     plan.Edges.Pms[0].Voltage,
            Days:        plan.Days,
            OrderID:     orderID,
            SubscribeID: subID,
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
// TODO 重构退款申请逻辑
func (s *orderService) Refund(riderID uint64, req *model.OrderRefundReq) {
    // if (req.Deposit == nil && req.OrderID == nil) || (req.Deposit != nil && req.OrderID != nil) {
    //     snag.Panic("退款参数错误")
    // }
    // q := s.orm.QueryNotDeleted().
    //     Where(
    //         order.RiderID(riderID),
    //         order.Status(model.OrderStatusPaid),
    //     )
    // if req.OrderID != nil {
    //     q.Where(order.ID(*req.OrderID))
    // }
    // if req.Deposit != nil && *req.Deposit {
    //     // 是否满足退押金条件
    //     if exist, _ := s.orm.QueryNotDeleted().Where(
    //         order.RiderID(riderID),
    //         order.EndAtIsNil(),
    //         order.StatusIn(model.OrderStatusPaid, model.OrderStatusRefundPending),
    //         order.TypeIn(model.OrderSubscribeTypes...),
    //     ).Exist(s.ctx); exist {
    //         snag.Panic("当前无法退押金")
    //     }
    //     q.Where(order.Type(model.OrderTypeDeposit))
    // }
    //
    // o, _ := q.Only(s.ctx)
    //
    // if o == nil {
    //     snag.Panic("未找到有效订单")
    // }
    // no := tools.NewUnique().NewSN28()
    // refund := &model.PaymentRefund{
    //     OrderID:      o.ID,
    //     TradeNo:      o.TradeNo,
    //     Total:        o.Total,
    //     RefundAmount: o.Amount,
    //     OutRefundNo:  no,
    //     Reason:       "用户申请",
    // }
    // // 缓存
    // cache.Set(s.ctx, no, &model.PaymentCache{CacheType: model.PaymentCacheTypeRefund, Refund: refund}, 0)
    //
    // switch o.Payway {
    // case model.OrderPaywayAlipay:
    //     payment.NewAlipay().Refund(refund)
    //     break
    // case model.OrderPaywayWechat:
    //     payment.NewWechat().Refund(refund)
    //     break
    // }
    //
    // if refund.Request {
    //     o.Update().SetStatus(model.OrderStatusRefundPending).SaveX(s.ctx)
    //     // 保存订单
    //     ar.Ent.OrderRefund.Create().
    //         SetOrder(o).
    //         SetAmount(refund.RefundAmount).
    //         SetOutRefundNo(no).
    //         SetReason(refund.Reason).
    //         SetStatus(model.OrderRefundStatusPending).
    //         SaveX(s.ctx)
    //
    //     if refund.Success {
    //         s.RefundSuccess(refund)
    //     }
    // }
}

// DoPayment 处理支付
func (s *orderService) DoPayment(pc *model.PaymentCache) {
    if pc == nil {
        return
    }
    switch pc.CacheType {
    case model.PaymentCacheTypePlan:
        s.OrderPaid(pc.Subscribe)
        break
    case model.PaymentCacheTypeRefund:
        s.RefundSuccess(pc.Refund)
    }
}

// OrderPaid 订单成功支付
// TODO 重构
// TODO 业绩提成逻辑
// TODO 2续签 3重签 4更改电池 5救援 6滞纳金 逻辑
func (s *orderService) OrderPaid(trade *model.PaymentSubscribe) {
    // 查询订单是否已存在
    if exists, err := ar.Ent.Order.Query().Where(order.OutTradeNo(trade.OutTradeNo)).Exist(s.ctx); err == nil && exists {
        return
    }

    j, _ := json.MarshalIndent(trade, "", "  ")
    log.Infof("[ORDER PAID %s] %s", trade.OutTradeNo, j)

    ctx := context.Background()
    tx, _ := ar.Ent.Tx(ctx)

    // 创建订单
    oc := tx.Order.Create().
        SetPayway(trade.Payway).
        SetPlanID(trade.PlanID).
        SetRiderID(trade.RiderID).
        SetAmount(decimal.NewFromFloat(trade.Amount).Sub(decimal.NewFromFloat(trade.Deposit)).InexactFloat64()).
        SetTotal(trade.Amount).
        SetOutTradeNo(trade.OutTradeNo).
        SetTradeNo(trade.TradeNo).
        SetStatus(model.OrderStatusPaid).
        SetType(trade.OrderType).
        SetCityID(trade.CityID).
        SetNillableParentID(trade.OrderID).
        SetNillableSubscribeID(trade.SubscribeID)
    o, err := oc.Save(s.ctx)
    if err != nil {
        log.Errorf("[ORDER PAID %s ERROR]: %s", trade.OutTradeNo, err.Error())
        _ = tx.Rollback()
        return
    }

    // 如果有押金, 创建押金订单
    var do *ent.Order
    if trade.Deposit > 0 {
        do, err = tx.Order.Create().
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
            log.Errorf("[ORDER PAID %s DEPOSIT ERROR]: %s", trade.OutTradeNo, err.Error())
            _ = tx.Rollback()
            return
        }
    }

    // 创建或更新subscribe
    // 新签或重签
    if trade.OrderType == model.OrderTypeNewly || trade.OrderType == model.OrderTypeAgain {
        // 创建subscribe
        // var sub *ent.Subscribe
        sc := tx.Subscribe.Create().
            SetType(model.OrderTypeNewly).
            SetRiderID(trade.RiderID).
            SetVoltage(trade.Voltage).
            SetDays(trade.Days).
            SetPlanID(trade.PlanID).
            SetCityID(trade.CityID).
            SetInitialOrderID(o.ID).
            AddOrders(o)
        if do != nil {
            sc.AddOrders(do)
        }
        _, err = sc.Save(ctx)
        if err != nil {
            log.Errorf("[ORDER PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
            _ = tx.Rollback()
            return
        }
    }

    // 续签
    if trade.OrderType == model.OrderTypeRenewal {
        _, err = tx.Subscribe.UpdateOneID(*trade.SubscribeID).AddDays(int(trade.Days)).Save(ctx)
        if err != nil {
            log.Errorf("[ORDER PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
            _ = tx.Rollback()
            return
        }
    }

    // 当新签和重签的时候有提成
    if trade.OrderType == model.OrderTypeNewly || (trade.OrderType == model.OrderTypeAgain && model.NewRecentSubscribePastDays(trade.PastDays).Commission()) {
        // 创建提成
        _, err = tx.Commission.Create().SetOrderID(o.ID).SetAmount(trade.Commission).SetStatus(model.CommissionStatusPending).Save(s.ctx)
        if err != nil {
            log.Errorf("[ORDER PAID %s COMMISSION(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
            _ = tx.Rollback()
            return
        }
    }

    // 删除缓存
    cache.Del(context.Background(), trade.OutTradeNo)

    // 提交事务
    _ = tx.Commit()
}

// RefundSuccess 成功退款
func (s *orderService) RefundSuccess(req *model.PaymentRefund) {
    // ctx := context.Background()
    // tx, _ := ar.Ent.Tx(ctx)
    // tx.Order.UpdateOneID(req.OrderID).
    //     SetStatus(model.OrderStatusRefundSuccess).
    //     SetEndAt(req.Time).
    //     SaveX(ctx)
    //
    // tx.OrderRefund.Update().
    //     Where(orderrefund.OutRefundNo(req.OutRefundNo)).
    //     SetStatus(model.OrderRefundStatusSuccess).
    //     SetRefundAt(req.Time).
    //     SaveX(ctx)
    //
    // // 删除提成订单
    // tx.Commission.Delete().Where(commission.OrderID(req.OrderID)).ExecX(ctx)
    //
    // _ = tx.Commit()
    //
    // // 删除缓存
    // cache.Del(ctx, req.OutRefundNo)
}
