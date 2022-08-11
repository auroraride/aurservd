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
    "github.com/auroraride/aurservd/internal/ent/commission"
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
        past = tools.NewPointer().Int(int(carbon.Time2Carbon(*sub.EndAt).AddDay().DiffInDays(carbon.Now())))
    }
    // 判定退订时间是否超出设置天数
    if model.NewRecentSubscribePastDays(*past).Commission() {
        state = model.OrderTypeNewly
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
    otype := req.OrderType
    switch otype {
    case model.OrderTypeNewly:
    case model.OrderTypeAgain:
        // 新签/重签判定
        if req.Model == "" || req.CityID == 0 {
            snag.Panic("请求参数错误")
        }
        otype, past = s.PreconditionNewly(sub)
        break
    case model.OrderTypeRenewal:
        // 续签判定
        s.PreconditionRenewal(sub)
        if sub.Remaining < 0 && int(op.Days)+sub.Remaining < 0 {
            snag.Panic("无法继续, 逾期天数大于套餐天数")
        }
        subID = tools.NewPointer().UInt64(sub.ID)
        orderID = tools.NewPointer().UInt64(sub.InitialOrderID)
        req.Model = sub.Model
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
    // 生成订单字段
    price := op.Price
    // TODO DEBUG 模式支付一分钱
    mode := ar.Config.App.Mode
    if mode == "debug" || mode == "next" {
        price = 0.01
        if deposit > 0 {
            deposit = 0.01
        }
    }
    // TODO DEBUG 记得删除

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
            PlanID:      op.ID,
            Deposit:     deposit,
            PastDays:    past,
            Commission:  op.Commission,
            Model:       req.Model,
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
    if sub == nil || sub.Remaining > 0 {
        snag.Panic("未找到逾期骑士卡信息")
    }

    fee, _, o := NewSubscribe().OverdueFee(riderID, sub)
    // TODO DEBUG 模式支付一分钱
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

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        var o *ent.Order
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
            SetNillableSubscribeID(trade.SubscribeID).
            SetNillablePastDays(trade.PastDays)
        o, err = oc.Save(s.ctx)
        if err != nil {
            log.Errorf("[ORDER PAID %s ERROR]: %s", trade.OutTradeNo, err.Error())
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
                SetModel(trade.Model).
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
            _, err = sc.Save(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
                return
            }
        }

        // 续签
        if trade.OrderType == model.OrderTypeRenewal {
            _, err = tx.Subscribe.UpdateOneID(*trade.SubscribeID).
                AddRenewalDays(int(trade.Days)).
                AddRemaining(int(trade.Days)).
                SetStatus(model.SubscribeStatusUsing).
                Save(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID %s SUBSCRIBE(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
                return
            }
        }

        // 当新签和重签的时候有提成
        if trade.OrderType == model.OrderTypeNewly {
            // 创建提成
            _, err = tx.Commission.Create().SetOrderID(o.ID).SetAmount(trade.Commission).SetStatus(model.CommissionStatusPending).Save(s.ctx)
            if err != nil {
                log.Errorf("[ORDER PAID %s COMMISSION(%d) ERROR]: %s", trade.OutTradeNo, o.ID, err.Error())
                return
            }
        }
        return
    })

    // 删除缓存
    cache.Del(context.Background(), trade.OutTradeNo)
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

func (s *orderService) listFilter(req model.OrderListFilter) (*ent.OrderQuery, map[string]interface{}) {
    tt := tools.NewTime()
    q := s.orm.QueryNotDeleted().
        Order(ent.Desc(order.FieldCreatedAt)).
        WithCity().
        WithPlan().
        WithCity().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithEmployee().WithStore()
        }).
        WithRefund()
    if req.Start != nil {
        q.Where(order.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }
    if req.End != nil {
        q.Where(order.CreatedAtLT(tt.ParseNextDateStringX(*req.End)))
    }
    if req.Type != nil {
        q.Where(order.Type(*req.Type))
    }
    if req.RiderID != nil {
        q.Where(order.RiderID(*req.RiderID))
    }
    if req.Keyword != nil {
        q.Where(order.HasRiderWith(
            rider.Or(
                rider.HasPersonWith(person.NameContainsFold(*req.Keyword)),
                rider.PhoneContainsFold(*req.Keyword),
            ),
        ))
    }
    if req.CityID != nil {
        q.Where(order.CityID(*req.CityID))
    }
    if req.EmployeeID != nil {
        q.Where(order.HasSubscribeWith(subscribe.EmployeeID(*req.EmployeeID)))
    }
    if req.StoreName != nil {
        q.Where(order.HasSubscribeWith(subscribe.HasStoreWith(store.NameContainsFold(*req.StoreName))))
    }
    if req.Model != nil {
        q.Where(order.HasSubscribeWith(subscribe.Model(*req.Model)))
    }
    if req.Days != nil {
        q.Where(order.InitialDaysGTE(*req.Days))
    }
    if req.Payway != nil {
        q.Where(order.Payway(*req.Payway))
    }
    if req.Refund != nil && *req.Refund > 0 {
        switch *req.Refund {
        case 1:
            q.Where(order.StatusNotIn(model.OrderStatusRefundPending, model.OrderStatusRefundRefused, model.OrderStatusRefundSuccess))
            break
        case 2:
            q.Where(order.StatusIn(model.OrderStatusRefundPending, model.OrderStatusRefundRefused, model.OrderStatusRefundSuccess)).WithRefund()
            break
        }
    }
    return q, nil
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
        for _, item := range items {
            detail := NewRiderOrder().Detail(item)
            var rows [][]any

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

            tools.NewExcel(path).AddValues(rows).Done()
        }
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
