// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/commission"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/orderrefund"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/payment"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type orderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.OrderClient
    employee *ent.Employee
}

func NewOrder() *orderService {
    return &orderService{
        ctx: context.Background(),
        orm: ent.Database.Order,
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
func (s *orderService) PreconditionNewly(sub *ent.Subscribe) (state uint, past *int) {
    if sub == nil {
        return model.OrderTypeNewly, nil
    }
    switch sub.Status {
    case model.SubscribeStatusInactive:
        snag.Panic("当前有未激活的骑士卡")
        break
    case model.SubscribeStatusUnSubscribed:
        state = model.OrderTypeAgain
        break
    default:
        snag.Panic("不满足新签条件")
        break
    }
    // 距离上次订阅过去的时间(从退订的第二天0点开始计算,不满一天算0天)
    if sub.EndAt != nil {
        past = silk.Int(int(carbon.Time2Carbon(*sub.EndAt).AddDay().DiffInDays(carbon.Now())))
        // 判定退订时间是否超出设置天数
        if model.NewRecentSubscribePastDays(*past).Commission() {
            state = model.OrderTypeNewly
        }
    }
    return
}

// PreconditionRenewal 是否满足续签条件
func (s *orderService) PreconditionRenewal(sub *ent.Subscribe) {
    if sub == nil ||
        sub.Status == model.SubscribeStatusInactive ||
        sub.Status == model.SubscribeStatusUnSubscribed ||
        sub.Status == model.SubscribeStatusCanceled {
        snag.Panic("当前无使用中的骑士卡")
    }
}

// Create 创建订单
// 仅新签和续签骑士卡可使用优惠和积分
// 滞纳金订单无法使用优惠券和积分
// 押金不参与任何优惠手段
func (s *orderService) Create(req *model.OrderCreateReq) (result *model.OrderCreateRes) {
    if req.OrderType == model.OrderTypeFee {
        return s.CreateFee(s.rider.ID, req.Payway)
    }

    result = new(model.OrderCreateRes)

    // 查询套餐是否存在
    p := NewPlan().QueryEffectiveWithID(req.PlanID)

    if p.BrandID != nil && p.Edges.Brand == nil {
        snag.Panic("骑士卡错误")
    }

    // 查询是否企业骑手
    if s.rider.EnterpriseID != nil {
        snag.Panic("团签用户无法购买")
    }

    // 查询是否有退款中的押金
    if exists, _ := ent.Database.Order.QueryNotDeleted().Where(
        order.RiderID(s.rider.ID),
        order.Status(model.OrderStatusRefundPending),
    ).Exist(s.ctx); exists {
        snag.Panic("当前有退款中的订单")
    }

    var past *int
    sub := NewSubscribe().Recent(s.rider.ID, *s.rider.PersonID)
    // 判定类型条件
    var subID, orderID *uint64
    t := req.OrderType
    switch t {
    case model.OrderTypeNewly, model.OrderTypeAgain:
        // 新签/重签判定
        if req.CityID == 0 {
            snag.Panic("请求参数错误")
        }
        t, past = s.PreconditionNewly(sub)
        break
    case model.OrderTypeRenewal:
        // 续签判定
        s.PreconditionRenewal(sub)
        if sub.Remaining < 0 && int(p.Days)+sub.Remaining < 0 {
            snag.Panic("无法继续, 逾期天数大于套餐天数")
        }
        subID = silk.UInt64(sub.ID)
        orderID = silk.UInt64(sub.InitialOrderID)
        req.CityID = sub.CityID
        break
    default:
        snag.Panic("未知的支付请求")
        break
    }

    // 判定用户是否需要缴纳押金
    deposit := NewRider().Deposit(s.rider.ID)
    no := tools.NewUnique().NewSN28()
    result.OutTradeNo = no

    // 计算需要支付金额
    // 1. 计算新签优惠
    // 2. 计算优惠券金额
    // 3. 计算积分抵扣
    price := p.Price

    // 计算新签优惠
    var ramount float64
    if t == model.OrderTypeNewly && p.DiscountNewly > 0 {
        ramount = p.DiscountNewly
        price = tools.NewDecimal().Sub(price, ramount)
    }

    // 获取优惠券
    var camount float64
    now := time.Now()
    if len(req.Coupons) > 0 {
        coupons := NewCoupon().QueryIDs(req.Coupons)
        if len(req.Coupons) != len(coupons) {
            snag.Panic("优惠券选择错误")
        }
        var isExclusive bool
        cm := make(map[uint64]uint64)
        for _, c := range coupons {
            // 校验有效期
            if c.ExpiresAt.Before(now) {
                snag.Panic("优惠券已失效")
            }

            // 是否互斥
            if c.Rule == model.CouponRuleExclusive.Value() {
                isExclusive = true
            }

            // 是否叠加
            if _, ok := cm[c.TemplateID]; ok {
                snag.Panic("优惠券无法叠加")
            }
            cm[c.TemplateID] = c.ID

            // 累加优惠券金额
            camount += c.Amount
        }
        if isExclusive && len(req.Coupons) > 0 {
            snag.Panic("所选优惠券互斥")
        }
        price = tools.NewDecimal().Sub(price, camount)
        // TODO price = 0时 直接支付成功
        // 暂时处理成支付一分钱
        if price <= 0 {
            price = 0.01
        }
    }

    // 积分抵扣
    var points int64
    if req.Point && price > 0.01 {
        pointServ := NewPoint()
        realPoints := pointServ.Real(s.rider)
        if realPoints > 0 {
            cents := int64(price / model.PointRatio)
            if realPoints < cents {
                // 若积分小于所需积分, 则全部扣除
                points = realPoints
            } else {
                // 若剩余积分大于所需金额, 则扣除剩余金额积分数量
                points = cents
            }
            price = tools.NewDecimal().Sub(price, float64(points)*model.PointRatio)
            _, err := pointServ.PreConsume(s.rider, points)
            if err != nil {
                snag.Panic("订单创建失败: %v", err)
            }
        }
        if price <= 0 {
            price = 0.01
        }
    }

    if price < 0 {
        snag.Panic("支付金额错误")
    }

    // DEBUG 模式支付一分钱
    mode := ar.Config.App.Mode
    if mode == "debug" || mode == "next" {
        price = 0.01
        if deposit > 0 {
            deposit = 0.01
        }
    }

    // 总计支付金额
    total := tools.NewDecimal().Sum(price, deposit)

    // 订单字段
    prepay := &model.PaymentCache{
        CacheType: model.PaymentCacheTypePlan,
        Subscribe: &model.PaymentSubscribe{
            CityID:        req.CityID,
            OrderType:     t,
            OutTradeNo:    no,
            RiderID:       s.rider.ID,
            Name:          "购买" + p.Name,
            Amount:        total,
            Payway:        req.Payway,
            Deposit:       deposit,
            PastDays:      past,
            Commission:    p.Commission,
            Model:         p.Model,
            Days:          p.Days,
            OrderID:       orderID,
            SubscribeID:   subID,
            Points:        points,
            PointRatio:    model.PointRatio,
            CouponAmount:  camount,
            Coupons:       req.Coupons,
            DiscountNewly: p.DiscountNewly,
            EbikeBrandID:  p.BrandID,
            Plan: &model.Plan{
                ID:   p.ID,
                Name: p.Name,
                Days: p.Days,
            },
        },
    }

    s.Prepay(req.Payway, no, prepay, result)

    return
}

// CreateFee 创建滞纳金订单
func (s *orderService) CreateFee(riderID uint64, payway uint8) *model.OrderCreateRes {
    result := new(model.OrderCreateRes)

    sub, _ := NewSubscribe().QueryEffective(riderID)
    if sub == nil || sub.Remaining > 0 {
        snag.Panic("未找到逾期骑士卡信息")
    }

    fee, _, o := NewSubscribe().OverdueFee(riderID, sub)

    // DEBUG 模式支付一分钱
    mode := ar.Config.App.Mode
    if mode == "debug" || mode == "next" {
        fee = 0.01
    }

    no := tools.NewUnique().NewSN28()
    prepay := &model.PaymentCache{
        CacheType: model.PaymentCacheTypeOverdueFee,
        OverDueFee: &model.PaymentOverdueFee{
            Subject:     fmt.Sprintf("逾期%d天费用", 0-sub.Remaining),
            OutTradeNo:  no,
            OrderType:   model.OrderTypeFee,
            Days:        0 - sub.Remaining,
            Amount:      fee,
            RiderID:     riderID,
            PlanID:      *sub.PlanID,
            SubscribeID: sub.ID,
            CityID:      sub.CityID,
            Payway:      payway,
        },
    }

    if o != nil {
        prepay.OverDueFee.OrderID = o.ID
    }

    s.Prepay(payway, no, prepay, result)
    result.OutTradeNo = no

    return result
}

// Prepay 预支付订单
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
    case model.PaymentCacheTypeAssistance:
        NewAssistance().Paid(pc.Assistance)
        break
    case model.PaymentCacheTypeRefund:
        s.RefundSuccess(pc.Refund)
        break
    case model.PaymentCacheTypeOverdueFee:
        s.FeePaid(pc.OverDueFee)
        break
    }
}

// FeePaid 滞纳金支付
func (s *orderService) FeePaid(trade *model.PaymentOverdueFee) {

    // 查询欠费订单是否已存在
    if exists, err := ent.Database.Order.Query().Where(order.OutTradeNo(trade.OutTradeNo)).Exist(s.ctx); err == nil && exists {
        return
    }

    j, _ := json.MarshalIndent(trade, "", "  ")
    log.Infof("[FEE PAID %s] %s", trade.OutTradeNo, j)

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        _, err = tx.Order.Create().
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
            Save(s.ctx)
        if err != nil {
            log.Errorf("[FEE PAID %s ERROR]: %s", trade.OutTradeNo, err.Error())
            return
        }

        // 更新订阅
        _, err = tx.Subscribe.
            UpdateOneID(trade.SubscribeID).
            SetStatus(model.SubscribeStatusUsing).
            SetRemaining(0).
            AddOverdueDays(trade.Days).
            Save(s.ctx)
        if err != nil {
            log.Errorf("[FEE PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, trade.SubscribeID, err.Error())
        }
        return
    })

    // 删除缓存
    cache.Del(context.Background(), trade.OutTradeNo)
}

// OrderPaid 订单成功支付
func (s *orderService) OrderPaid(trade *model.PaymentSubscribe) {
    // 查询订单是否已存在
    if exists, err := ent.Database.Order.Query().Where(order.OutTradeNo(trade.OutTradeNo)).Exist(s.ctx); err == nil && exists {
        return
    }

    j, _ := json.MarshalIndent(trade, "", "  ")
    log.Infof("[ORDER PAID %s] %s", trade.OutTradeNo, j)

    var sub *ent.Subscribe
    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        var o *ent.Order
        // 创建订单
        oc := tx.Order.Create().
            SetPayway(trade.Payway).
            SetPlanID(trade.Plan.ID).
            SetRiderID(trade.RiderID).
            SetAmount(tools.NewDecimal().Sub(trade.Amount, trade.Deposit)).
            SetTotal(trade.Amount).
            SetOutTradeNo(trade.OutTradeNo).
            SetTradeNo(trade.TradeNo).
            SetStatus(model.OrderStatusPaid).
            SetType(trade.OrderType).
            SetCityID(trade.CityID).
            SetInitialDays(int(trade.Days)).
            SetNillableParentID(trade.OrderID).
            SetNillableSubscribeID(trade.SubscribeID).
            SetNillablePastDays(trade.PastDays).
            SetPoints(trade.Points).
            SetPointRatio(trade.PointRatio).
            SetCouponAmount(trade.CouponAmount).
            SetDiscountNewly(trade.DiscountNewly).
            SetNillableBrandID(trade.EbikeBrandID)
        if len(trade.Coupons) > 0 {
            oc.AddCouponIDs(trade.Coupons...)
            // 更新优惠券使用状态
            for _, couponID := range trade.Coupons {
                err = tx.Coupon.UpdateOneID(couponID).SetPlanID(trade.Plan.ID).SetUsedAt(time.Now()).Exec(s.ctx)
                if err != nil {
                    log.Errorf("[ORDER PAID %s COUPON id = %d ERROR]: %s", trade.OutTradeNo, couponID, err.Error())
                }
            }
        }
        o, err = oc.Save(s.ctx)
        if err != nil {
            log.Errorf("[ORDER PAID %s ERROR]: %s", trade.OutTradeNo, err.Error())
            return
        }

        // 更新积分
        if trade.Points > 0 {
            var r *ent.Rider
            r, err = tx.Rider.UpdateOneID(trade.RiderID).AddPoints(-trade.Points).Save(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID POINT UPDATE %s ERROR]: %s", trade.OutTradeNo, err.Error())
                return
            }

            err = tx.PointLog.Create().
                SetPoints(trade.Points).
                SetRiderID(trade.RiderID).
                SetReason("订阅骑士卡").
                SetType(model.PointLogTypeConsume.Value()).
                SetAfter(r.Points).
                SetAttach(&model.PointLogAttach{Plan: trade.Plan}).
                SetOrder(o).
                Exec(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID POINT LOG %s ERROR]: %s", trade.OutTradeNo, err.Error())
                return
            }

            NewPoint().RemovePreConsume(r, trade.Points)
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
                return
            }
        }

        // 创建或更新subscribe
        // 新签或重签
        if trade.OrderType == model.OrderTypeNewly || trade.OrderType == model.OrderTypeAgain {
            // 创建subscribe
            sq := tx.Subscribe.Create().
                SetType(trade.OrderType).
                SetRiderID(trade.RiderID).
                SetModel(trade.Model).
                SetRemaining(int(trade.Days)).
                SetInitialDays(int(trade.Days)).
                SetStatus(model.SubscribeStatusInactive).
                SetPlanID(trade.Plan.ID).
                SetCityID(trade.CityID).
                SetInitialOrderID(o.ID).
                AddOrders(o).
                SetNillableBrandID(trade.EbikeBrandID).
                SetNeedContract(true)
            if do != nil {
                sq.AddOrders(do)
            }
            sub, err = sq.Save(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
                return
            }
        }

        // 续签
        if trade.OrderType == model.OrderTypeRenewal {
            status := model.SubscribeStatusUsing
            // 查询状态
            sub, _ = NewSubscribe().Query(*trade.SubscribeID)
            if sub.Status == model.SubscribeStatusPaused {
                status = model.SubscribeStatusPaused
            }

            sub, err = tx.Subscribe.UpdateOneID(*trade.SubscribeID).
                AddRenewalDays(int(trade.Days)).
                AddRemaining(int(trade.Days)).
                SetStatus(status).
                Save(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
                return
            }
        }

        // 当新签和重签的时候有提成
        if trade.OrderType == model.OrderTypeNewly && trade.Commission > 0 && sub != nil {
            // 创建提成
            _, err = tx.Commission.Create().SetOrderID(o.ID).SetPlanID(trade.Plan.ID).SetAmount(trade.Commission).SetStatus(model.CommissionStatusPending).SetRiderID(sub.RiderID).SetSubscribe(sub).Save(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID %s COMMISSION(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
                return
            }
        }
        return
    })

    // 删除缓存
    cache.Del(context.Background(), trade.OutTradeNo)

    if sub != nil && trade.OrderType == model.OrderTypeRenewal {
        go func() {
            _ = NewSubscribe().UpdateStatus(sub, false)
        }()
    }
}

// RefundSuccess 成功退款
func (s *orderService) RefundSuccess(req *model.PaymentRefund) {
    ctx := context.Background()

    // 删除缓存
    cache.Del(ctx, req.OutRefundNo)

    log.Infof("%s(OrderID:%d) [退款]退款完成, 实际退款时间: %s", req.OutRefundNo, req.OrderID, req.Time)

    // 更新订单
    _, err := ent.Database.Order.UpdateOneID(req.OrderID).
        SetStatus(model.OrderStatusRefundSuccess).
        SetRefundAt(req.Time).
        Save(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]原订单更新完成", req.OutRefundNo, req.OrderID)

    // 更新退款订单
    _, err = ent.Database.OrderRefund.Update().
        Where(orderrefund.OutRefundNo(req.OutRefundNo)).
        SetStatus(model.RefundStatusSuccess).
        SetRefundAt(req.Time).
        Save(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]退款订单更新完成", req.OutRefundNo, req.OrderID)

    // 更新骑士卡
    _, err = ent.Database.Subscribe.Update().Where(subscribe.InitialOrderID(req.OrderID)).SetRefundAt(req.Time).SetStatus(model.SubscribeStatusCanceled).Save(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]骑士卡更新完成", req.OutRefundNo, req.OrderID)

    // 删除提成订单
    err = ent.Database.Commission.SoftDelete().Where(commission.OrderID(req.OrderID)).SetRemark("用户已退款").Exec(ctx)
    if err != nil {
        log.Error(err)
    }
    log.Infof("%s(OrderID:%d) [退款]提成订单更新完成", req.OutRefundNo, req.OrderID)
}

func (s *orderService) listFilter(req model.OrderListFilter) (*ent.OrderQuery, ar.Map) {
    info := make(ar.Map)
    tt := tools.NewTime()
    q := s.orm.QueryNotDeleted().
        Order(ent.Desc(order.FieldCreatedAt)).
        WithCity().
        WithPlan().
        WithCity().
        WithRider().
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithEmployee().WithStore()
        }).
        WithRefund().
        WithCoupons().
        WithEbike().
        WithBrand()
    if req.Start != nil {
        info["开始日期"] = req.Start
        q.Where(order.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }
    if req.End != nil {
        info["结束日期"] = req.Start
        q.Where(order.CreatedAtLT(tt.ParseNextDateStringX(*req.End)))
    }
    if req.Type != nil && *req.Type != 0 {
        info["类别"] = model.OrderTypes[*req.Type]
        q.Where(order.Type(*req.Type))
    }
    if req.RiderID != nil {
        info["骑手"] = ent.NewExportInfo(*req.RiderID, rider.Table)
        q.Where(order.RiderID(*req.RiderID))
    }
    if req.Keyword != nil {
        info["关键词"] = *req.Keyword
        q.Where(order.HasRiderWith(
            rider.Or(
                rider.NameContainsFold(*req.Keyword),
                rider.PhoneContainsFold(*req.Keyword),
            ),
        ))
    }
    if req.CityID != nil {
        info["城市"] = ent.NewExportInfo(*req.CityID, city.Table)
        q.Where(order.CityID(*req.CityID))
    }
    if req.EmployeeID != nil {
        info["店员"] = ent.NewExportInfo(*req.EmployeeID, employee.Table)
        q.Where(order.HasSubscribeWith(subscribe.EmployeeID(*req.EmployeeID)))
    }
    if req.StoreName != nil {
        info["门店关键词"] = *req.StoreName
        q.Where(order.HasSubscribeWith(subscribe.HasStoreWith(store.NameContainsFold(*req.StoreName))))
    }
    if req.Model != nil {
        info["型号"] = *req.Model
        q.Where(order.HasSubscribeWith(subscribe.Model(*req.Model)))
    }
    if req.Days != nil {
        info["最小天数"] = *req.Days
        q.Where(order.InitialDaysGTE(*req.Days))
    }
    if req.Payway != nil {
        info["支付方式"] = model.OrderPayways[*req.Payway]
        q.Where(order.Payway(*req.Payway))
    }
    if req.Refund != nil && *req.Refund > 0 {
        k := "退款申请"
        v := " - "
        switch *req.Refund {
        case 1:
            q.Where(order.StatusNotIn(model.OrderStatusRefundPending, model.OrderStatusRefundRefused, model.OrderStatusRefundSuccess))
            v = "未申请"
            break
        case 2:
            q.Where(order.StatusIn(model.OrderStatusRefundPending, model.OrderStatusRefundRefused, model.OrderStatusRefundSuccess)).WithRefund()
            v = "已申请"
            break
        }
        info[k] = v
    }
    return q, info
}

// List 获取订单列表
func (s *orderService) List(req *model.OrderListReq) *model.PaginationRes {
    q, _ := s.listFilter(req.OrderListFilter)
    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        NewRiderOrder().Detail,
    )
}

func (s *orderService) Export(req *model.OrderListExport) model.ExportRes {
    q, info := s.listFilter(req.OrderListFilter)
    return NewExportWithModifier(s.modifier).Start("订单", req.OrderListFilter, info, req.Remark, func(path string) {
        items, _ := q.All(s.ctx)
        var rows tools.ExcelItems
        title := []any{
            "类型",
            "骑手",
            "电话",
            "骑士卡",
            "天数",
            "型号",
            "城市",
            "门店",
            "店员",
            "支付状态",
            "订单编号",
            "支付编号",
            "支付方式",
            "支付金额",
            "支付时间",
        }
        rows = append(rows, title)

        for _, item := range items {
            detail := NewRiderOrder().Detail(item)

            st := ""
            if detail.Store != nil {
                st = detail.Store.Name
            }

            em := ""
            if detail.Employee != nil {
                em = detail.Employee.Name
            }

            var pl string
            var pd uint
            if detail.Plan != nil {
                pl = detail.Plan.Name
                pd = detail.Plan.Days
            }

            row := []any{
                model.OrderTypes[detail.Type],
                detail.Rider.Name,
                detail.Rider.Phone,
                pl,
                pd,
                detail.Model,
                detail.City.Name,
                st,
                em,
                model.OrderStatuses[detail.Status],
                detail.OutTradeNo,
                detail.TradeNo,
                model.OrderPayways[detail.Payway],
                fmt.Sprintf("%.2f", detail.Amount),
                detail.PayAt,
            }

            rows = append(rows, row)
        }

        tools.NewExcel(path).AddValues(rows).Done()
    })
}

// QueryStatus 查询订单状态
func (s *orderService) QueryStatus(req *model.OrderStatusReq) (res model.OrderStatusRes) {
    now := time.Now()
    res = model.OrderStatusRes{
        OutTradeNo: req.OutTradeNo,
        Paid:       false,
    }
    for {
        o, _ := ent.Database.Order.QueryNotDeleted().Where(order.OutTradeNo(req.OutTradeNo)).First(s.ctx)

        if o != nil && o.Status == model.OrderStatusPaid {
            res.Paid = true
            return
        }

        if time.Now().Sub(now).Seconds() > 30 {
            return
        }

        time.Sleep(1 * time.Second)
    }
}
