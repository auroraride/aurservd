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
    "github.com/auroraride/aurservd/internal/ent/commission"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/orderrefund"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/payment"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
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

// RencentSubscribeOrder 获取骑手最近的骑士卡订单
func (s *orderService) RencentSubscribeOrder(riderID uint64) (*ent.Order, error) {
    return s.orm.QueryNotDeleted().
        Where(
            order.RiderID(riderID),
            order.TypeIn(model.OrderTypeNewly, model.OrderTypeRenewal, model.OrderTypeAgain, model.OrderTypeTransform),
            order.Status(model.OrderStatusPaid),
        ).
        Order(ent.Desc(order.FieldCreatedAt)).
        WithPlan().
        First(s.ctx)
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
}

// Create 创建订单
// TODO 需要做更改电池逻辑, 退款 + 推送订单
func (s *orderService) Create(req *model.OrderCreateReq) (result *model.OrderCreateRes) {
    if req.OrderType == model.OrderTypeFee {
        return s.CreateFee(s.rider.ID, req.Payway)
    }

    result = new(model.OrderCreateRes)

    // 查询套餐是否存在
    op := NewPlan().QueryEffectiveWithID(req.PlanID)

    // 查询是否企业骑手
    if s.rider.EnterpriseID != nil {
        snag.Panic("团签用户无法购买")
    }

    // 查询骑手是否签约过
    if !NewContract().Effective(s.rider) {
        snag.Panic("请先签约")
    }

    // 查询是否有退款中的押金
    if exists, _ := ar.Ent.Order.QueryNotDeleted().Where(
        order.RiderID(s.rider.ID),
        // order.Type(model.OrderTypeDeposit),
        order.Status(model.OrderStatusRefundPending),
    ).Exist(s.ctx); exists {
        snag.Panic("当前有退款中的订单")
    }

    sub := NewSubscribe().RecentDetail(s.rider.ID)
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
        if sub.Remaining < 0 && int(op.Days)+sub.Remaining < 0 {
            snag.Panic("无法继续, 逾期天数大于套餐天数")
        }
        subID = tools.NewPointer().UInt64(sub.ID)
        orderID = tools.NewPointer().UInt64(sub.Order.ID)
        break
    default:
        snag.Panic("未知的支付请求")
        break
    }

    // 距离上次订阅过去的时间(从退订的第二天0点开始计算,不满一天算0天)
    var pastDays int
    if otype == model.OrderTypeAgain {
        pastDays = int(carbon.Parse(sub.EndAt).AddDay().DiffInDays(carbon.Now()))
    }

    // 判定用户是否需要缴纳押金
    deposit := NewRider().Deposit(s.rider.ID)
    no := tools.NewUnique().NewSN28()
    result.OutTradeNo = no
    // 生成订单字段
    price := op.Price
    // DEBUG 模式支付一分钱
    mode := ar.Config.App.Mode
    if mode == "debug" || mode == "next" {
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
            Name:        "购买" + op.Name,
            Amount:      total,
            Payway:      req.Payway,
            Expire:      time.Now().Add(10 * time.Minute),
            PlanID:      op.ID,
            Deposit:     deposit,
            PastDays:    pastDays,
            Commission:  op.Commission,
            Voltage:     op.Edges.Pms[0].Voltage,
            Days:        op.Days,
            OrderID:     orderID,
            SubscribeID: subID,
        },
    }

    s.Prepay(req.Payway, no, prepay, result)

    return
}

func (s *orderService) CreateFee(riderID uint64, payway uint8) *model.OrderCreateRes {
    result := new(model.OrderCreateRes)

    sub, _ := NewSubscribe().QueryEffective(riderID)
    if sub == nil || sub.Remaining < 0 {
        snag.Panic("为找到逾期骑士卡信息")
    }

    fee, _, o := NewSubscribe().OverdueFee(riderID, sub.Remaining)
    no := tools.NewUnique().NewSN28()
    prepay := &model.PaymentCache{
        CacheType: model.PaymentCacheTypeOverdueFee,
        OverDueFee: &model.PaymentOverdueFee{
            OutTradeNo:  no,
            OrderType:   model.OrderTypeFee,
            Days:        0 - sub.Remaining,
            Amount:      fee,
            RiderID:     riderID,
            PlanID:      *sub.PlanID,
            OrderID:     o.ID,
            SubscribeID: sub.ID,
            CityID:      sub.CityID,
            Payway:      payway,
        },
    }

    s.Prepay(payway, no, prepay, result)

    return result
}

func (s *orderService) Prepay(payway uint8, no string, prepay *model.PaymentCache, result *model.OrderCreateRes) {
    // 订单缓存
    err := cache.Set(s.ctx, no, prepay, 20*time.Minute).Err()
    if err != nil {
        log.Error(err)
        snag.Panic("订单创建失败")
    }
    var str string
    switch payway {
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
        break
    case model.PaymentCacheTypeOverdueFee:
        s.FeePaid(pc.OverDueFee)
        break
    }
}

func (s *orderService) FeePaid(trade *model.PaymentOverdueFee) {

    // 查询欠费订单是否已存在
    if exists, err := ar.Ent.Order.Query().Where(order.OutTradeNo(trade.OutTradeNo)).Exist(s.ctx); err == nil && exists {
        return
    }

    j, _ := json.MarshalIndent(trade, "", "  ")
    log.Infof("[FEE PAID %s] %s", trade.OutTradeNo, j)

    ctx := context.Background()
    tx, _ := ar.Ent.Tx(ctx)

    _, err := tx.Order.Create().
        SetPayway(trade.Payway).
        SetPlanID(trade.PlanID).
        SetParentID(trade.OrderID).
        SetAmount(trade.Amount).
        SetOutTradeNo(trade.OutTradeNo).
        SetTradeNo(trade.TradeNo).
        SetTotal(trade.Amount).
        SetCityID(trade.CityID).
        SetSubscribeID(trade.SubscribeID).
        SetType(trade.OrderType).
        SetInitialDays(trade.Days).
        SetRiderID(trade.RiderID).
        Save(ctx)
    if err != nil {
        log.Errorf("[FEE PAID %s ERROR]: %s", trade.OutTradeNo, err.Error())
        _ = tx.Rollback()
        return
    }

    // 更新订阅
    _, err = tx.Subscribe.
        UpdateOneID(trade.SubscribeID).
        SetStatus(model.SubscribeStatusUsing).
        SetRemaining(0).
        SetOverdueDays(trade.Days).
        Save(ctx)
    if err != nil {
        log.Errorf("[FEE PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, trade.SubscribeID, err.Error())
        _ = tx.Rollback()
        return
    }

    _ = tx.Commit()

    // 删除缓存
    cache.Del(context.Background(), trade.OutTradeNo)
}

// OrderPaid 订单成功支付
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
        SetInitialDays(int(trade.Days)).
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
        sc := tx.Subscribe.Create().
            SetType(trade.OrderType).
            SetRiderID(trade.RiderID).
            SetVoltage(trade.Voltage).
            SetRemaining(int(trade.Days)).
            SetInitialDays(int(trade.Days)).
            SetStatus(model.SubscribeStatusInactive).
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
        _, err = tx.Subscribe.UpdateOneID(*trade.SubscribeID).
            AddRenewalDays(int(trade.Days)).
            AddRemaining(int(trade.Days)).
            Save(ctx)
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
    ctx := context.Background()

    // 删除缓存
    cache.Del(ctx, req.OutRefundNo)

    log.Infof("%s(OrderID:%d) [退款]退款完成, 实际退款时间: %s", req.OutRefundNo, req.OrderID, req.Time)

    // 更新订单
    _, err := ar.Ent.Order.UpdateOneID(req.OrderID).
        SetStatus(model.OrderStatusRefundSuccess).
        SetRefundAt(req.Time).
        Save(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]原订单更新完成", req.OutRefundNo, req.OrderID)

    // 更新退款订单
    _, err = ar.Ent.OrderRefund.Update().
        Where(orderrefund.OutRefundNo(req.OutRefundNo)).
        SetStatus(model.RefundStatusSuccess).
        SetRefundAt(req.Time).
        Save(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]退款订单更新完成", req.OutRefundNo, req.OrderID)

    // 更新骑士卡
    _, err = ar.Ent.Subscribe.Update().Where(subscribe.InitialOrderID(req.OrderID)).SetRefundAt(req.Time).SetStatus(model.SubscribeStatusCanceled).Save(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]骑士卡更新完成", req.OutRefundNo, req.OrderID)

    // 删除提成订单
    err = ar.Ent.Commission.SoftDelete().Where(commission.OrderID(req.OrderID)).SetRemark("用户已退款").Exec(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]提成订单更新完成", req.OutRefundNo, req.OrderID)
}

// List 获取订单列表
func (s *orderService) List(req *model.OrderListReq) *model.PaginationRes {
    tt := tools.NewTime()
    q := s.orm.QueryNotDeleted().
        Order(ent.Desc(order.FieldCreatedAt)).
        WithCity().
        WithPlan(func(pq *ent.PlanQuery) {
            pq.WithPms()
        }).
        WithCity().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithEmployee().WithStore()
        }).
        WithRefund()
    if req.RiderID != nil {
        q.Where(order.RiderID(*req.RiderID))
    }
    if req.Start != nil {
        q.Where(order.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }
    if req.End != nil {
        q.Where(order.CreatedAtLTE(tt.ParseDateStringX(*req.End)))
    }
    if req.Type != nil {
        q.Where(order.Type(*req.Type))
    }
    if req.RiderName != nil {
        q.Where(order.HasRiderWith(rider.HasPersonWith(person.NameContainsFold(*req.RiderName))))
    }
    if req.RiderPhone != nil {
        q.Where(order.HasRiderWith(rider.PhoneContainsFold(*req.RiderPhone)))
    }
    if req.CityID != nil {
        q.Where(order.CityID(*req.CityID))
    }
    // TODO 救援订单
    if req.EmployeeName != nil {
        q.Where(order.HasSubscribeWith(subscribe.HasEmployeeWith(employee.NameContainsFold(*req.EmployeeName))))
    }
    // TODO 救援订单
    if req.StoreName != nil {
        q.Where(order.HasSubscribeWith(subscribe.HasStoreWith(store.NameContainsFold(*req.StoreName))))
    }
    if req.Voltage != nil {
        q.Where(order.HasSubscribeWith(subscribe.Voltage(*req.Voltage)))
    }
    if req.Days != nil {
        q.Where(order.InitialDaysGTE(*req.Days))
    }
    if req.Refund > 0 {
        switch req.Refund {
        case 1:
            q.Where(order.StatusNotIn(model.OrderStatusRefundPending, model.OrderStatusRefundRefused, model.OrderStatusRefundSuccess))
            break
        case 2:
            q.Where(order.StatusIn(model.OrderStatusRefundPending, model.OrderStatusRefundRefused, model.OrderStatusRefundSuccess)).WithRefund()
            break
        }
    }
    return model.ParsePaginationResponse[model.RiderOrder, ent.Order](
        q,
        req.PaginationReq,
        NewRiderOrder().Detail,
    )
}
