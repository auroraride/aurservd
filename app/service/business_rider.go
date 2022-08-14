// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type businessRiderService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    cabinet      *ent.Cabinet
    cabinetInfo  *model.CabinetBasicInfo
    store        *ent.Store
    subscribe    *ent.Subscribe
    reserve      *ent.Reserve

    task func() *ec.BinInfo // 电柜任务

    storeID, employeeID, cabinetID, subscribeID *uint64
}

func NewBusinessRider(r *ent.Rider) *businessRiderService {
    return &businessRiderService{
        ctx:   context.Background(),
        rider: r,
    }
}

func NewBusinessRiderWithModifier(m *model.Modifier) *businessRiderService {
    s := NewBusinessRider(nil)
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *businessRiderService) SetCabinet(cab *ent.Cabinet) *businessRiderService {
    if cab != nil {
        s.cabinet = cab
        s.cabinetInfo = &model.CabinetBasicInfo{
            ID:     s.cabinet.ID,
            Brand:  model.CabinetBrand(s.cabinet.Brand),
            Serial: s.cabinet.Serial,
            Name:   s.cabinet.Name,
        }
    }
    return s
}

func (s *businessRiderService) SetCabinetID(id *uint64) *businessRiderService {
    if id != nil {
        s.SetCabinet(NewCabinet().QueryOne(*id))
    }
    return s
}

func (s *businessRiderService) SetStoreID(id *uint64) *businessRiderService {
    if id != nil {
        s.store = NewStore().Query(*id)
    }
    return s
}

func (s *businessRiderService) SetTask(task func() *ec.BinInfo) *businessRiderService {
    if task != nil {
        s.task = task
    }
    return s
}

func NewBusinessRiderWithEmployee(e *ent.Employee) *businessRiderService {
    s := NewBusinessRider(nil)
    s.ctx = context.WithValue(s.ctx, "employee", e)
    store := e.Edges.Store
    if store == nil {
        snag.Panic("未找到所属门店")
    }
    s.store = store
    s.employee = e
    s.employeeInfo = &model.Employee{
        ID:    s.employee.ID,
        Name:  s.employee.Name,
        Phone: s.employee.Phone,
    }
    return s
}

// QuerySubscribeWithRider 查询订阅信息
func (s *businessRiderService) QuerySubscribeWithRider(subscribeID uint64) *ent.Subscribe {
    item, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(subscribeID)).WithRider().Only(s.ctx)
    if item == nil {
        snag.Panic("未找到对应订阅")
    }
    return item
}

// Inactive 获取待激活订阅信息
func (s *businessRiderService) Inactive(id uint64) (*model.SubscribeActiveInfo, *ent.Subscribe) {
    // 查询订单状态
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().
        Where(
            subscribe.ID(id),
            subscribe.RefundAtIsNil(),
            subscribe.StartAtIsNil(),
            subscribe.Or(
                subscribe.Type(0),
                subscribe.TypeIn(model.OrderTypeNewly, model.OrderTypeAgain),
            ),
            subscribe.Status(model.SubscribeStatusInactive),
        ).
        WithInitialOrder(func(oq *ent.OrderQuery) {
            oq.WithPlan().WithCommission()
        }).
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithEnterprise().
        WithCity().
        Only(s.ctx)

    if sub == nil {
        snag.Panic("未找到待激活骑士卡")
    }

    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID)
    }

    r := sub.Edges.Rider
    if r == nil {
        snag.Panic("骑手信息获取失败")
    }

    NewRiderPermissionWithRider(r).BusinessX()

    p := sub.Edges.Rider.Edges.Person
    if p.Status != model.PersonAuthenticated.Value() {
        snag.Panic("骑手还未认证")
    }

    res := &model.SubscribeActiveInfo{
        ID:           sub.ID,
        EnterpriseID: sub.EnterpriseID,
        Model:        sub.Model,
        CommissionID: nil,
        Rider: model.RiderBasic{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  p.Name,
        },
        City: model.City{
            ID:   sub.Edges.City.ID,
            Name: sub.Edges.City.Name,
        },
    }

    if sub.EnterpriseID == nil {
        o := sub.Edges.InitialOrder
        if o == nil || o.Status != model.OrderStatusPaid {
            snag.Panic("订单状态异常") // TODO 退款被拒绝的如何操作
        }
        res.Order = &model.SubscribeOrderInfo{
            ID:      o.ID,
            Status:  o.Status,
            PayAt:   o.CreatedAt.Format(carbon.DateTimeLayout),
            Payway:  o.Payway,
            Amount:  o.Amount,
            Deposit: o.Total - o.Amount,
            Total:   o.Total,
        }

        c := sub.Edges.InitialOrder.Edges.Commission
        if c != nil {
            res.CommissionID = &c.ID
        }
    } else {
        res.Enterprise = &model.EnterpriseBasic{
            ID:   sub.Edges.Enterprise.ID,
            Name: sub.Edges.Enterprise.Name,
        }
    }
    return res, sub
}

// preprocess 预处理数据
func (s *businessRiderService) preprocess(typ business.Type, sub *ent.Subscribe) {
    s.subscribe = sub
    s.subscribeID = tools.NewPointerInterface(sub.ID)

    r := sub.Edges.Rider
    if r == nil {
        r, _ = sub.QueryRider().First(s.ctx)
    }

    if r == nil {
        snag.Panic("骑手查询失败")
    }

    s.rider = r

    if s.store != nil {
        s.storeID = tools.NewPointerInterface(s.store.ID)
    }

    if s.cabinet != nil {
        s.cabinetID = tools.NewPointerInterface(s.cabinet.ID)
    }

    if s.store == nil && s.cabinet == nil {
        snag.Panic("条件不满足")
    }

    if s.employee != nil {
        s.employeeID = tools.NewPointerInterface(s.employee.ID)
    }

    if s.employee == nil && s.modifier == nil && s.cabinet == nil {
        snag.Panic("操作权限校验失败")
    }

    // 校验权限
    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(s.subscribe.CityID)
    }

    // 预约检查
    rev := NewReserveWithRider(r).RiderUnfinished(r.ID)
    if rev != nil {
        if s.cabinet == nil || s.cabinet.ID != rev.CabinetID || typ.String() != rev.Type {
            _, _ = rev.Update().SetStatus(model.ReserveStatusInvalid.Value()).Save(s.ctx)
        } else {
            // 预约处理中
            s.reserve, _ = rev.Update().SetStatus(model.ReserveStatusProcessing.Value()).Save(s.ctx)
        }
    }
}

// doTask 处理电柜任务
func (s *businessRiderService) doTask() (bin *ec.BinInfo, err error) {
    defer func() {
        if v := recover(); v != nil {
            err = fmt.Errorf("%v", v)
        }
    }()

    bin = s.task()
    return
}

// do 处理业务
func (s *businessRiderService) do(bt business.Type, cb func(tx *ent.Tx)) {
    sts := map[business.Type]uint8{
        business.TypeActive:      model.StockTypeRiderObtain,
        business.TypeUnsubscribe: model.StockTypeRiderUnSubscribe,
        business.TypePause:       model.StockTypeRiderPause,
        business.TypeContinue:    model.StockTypeRiderContinue,
    }

    ops := map[business.Type]model.Operate{
        business.TypeActive:      model.OperateActive,
        business.TypeUnsubscribe: model.OperateUnsubscribe,
        business.TypePause:       model.OperateSubscribePause,
        business.TypeContinue:    model.OperateSubscribeContinue,
    }

    bfs := map[business.Type]string{
        business.TypeActive:      "未激活",
        business.TypeUnsubscribe: "生效中",
        business.TypePause:       "计费中",
        business.TypeContinue:    "寄存中",
    }

    afs := map[business.Type]string{
        business.TypeActive:      "已激活",
        business.TypeUnsubscribe: "已退租",
        business.TypePause:       "已寄存",
        business.TypeContinue:    "计费中",
    }

    var bin *ec.BinInfo
    var err error

    // 放入电池优先执行
    if s.task != nil && (bt == business.TypePause || bt == business.TypeUnsubscribe) {
        bin, err = s.doTask()
        if err != nil {
            snag.Panic(err)
        }
    }

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        cb(tx)

        return NewStockWithModifier(s.modifier).BatteryWithRider(
            tx.Stock.Create(),
            &model.StockBusinessReq{
                RiderID:   s.subscribe.RiderID,
                Model:     s.subscribe.Model,
                CityID:    s.subscribe.CityID,
                StockType: sts[bt],

                StoreID:     s.storeID,
                EmployeeID:  s.employeeID,
                CabinetID:   s.cabinetID,
                SubscribeID: s.subscribeID,
            },
        )
    })

    // 取出电池优后执行
    if s.task != nil && (bt == business.TypeActive || bt == business.TypeContinue) {
        bin, err = s.doTask()
        if err != nil {
            log.Error(err)
        }
    }

    // 保存业务日志
    var b *ent.Business
    b, err = NewBusinessLog(s.subscribe).
        SetModifier(s.modifier).
        SetEmployee(s.employee).
        SetCabinet(s.cabinet).
        SetBinInfo(bin).
        Save(bt)
    var bussinessID *uint64
    revStatus := model.ReserveStatusFail
    if b != nil {
        revStatus = model.ReserveStatusSuccess
        bussinessID = tools.NewPointerInterface(b.ID)
    }

    // 更新预约
    if s.reserve != nil {
        _, _ = s.reserve.Update().
            SetStatus(revStatus.Value()).
            SetNillableBusinessID(bussinessID).
            Save(s.ctx)
    }

    // 记录日志
    go logging.NewOperateLog().
        SetRef(s.rider).
        SetOperate(ops[bt]).
        SetEmployee(s.employeeInfo).
        SetModifier(s.modifier).
        SetCabinet(s.cabinetInfo).
        SetDiff(bfs[bt], afs[bt]).
        Send()

    if err != nil {
        snag.Panic(err)
    }
}

// Active 激活订阅
func (s *businessRiderService) Active(info *model.SubscribeActiveInfo, sub *ent.Subscribe) {
    s.preprocess(business.TypeActive, sub)

    s.do(business.TypeActive, func(tx *ent.Tx) {
        // 激活
        var err error
        s.subscribe, err = tx.Subscribe.UpdateOneID(info.ID).
            SetStatus(model.SubscribeStatusUsing).
            SetStartAt(time.Now()).
            SetNillableEmployeeID(s.employeeID).
            SetNillableStoreID(s.storeID).
            SetNillableCabinetID(s.cabinetID).
            Save(s.ctx)
        snag.PanicIfError(err)

        // 提成
        if info.CommissionID != nil && s.employeeID != nil {
            _, err = tx.Commission.UpdateOneID(*info.CommissionID).SetEmployeeID(*s.employeeID).Save(s.ctx)
            snag.PanicIfError(err)
        }
    })

    if info.EnterpriseID != nil {
        // 更新团签账单
        go NewEnterprise().UpdateStatement(sub.Edges.Enterprise)
    }
}

// UnSubscribe 退租
// 会抹去欠费情况
func (s *businessRiderService) UnSubscribe(subscribeID uint64) {
    s.preprocess(business.TypeUnsubscribe, s.QuerySubscribeWithRider(subscribeID))

    sub := s.subscribe

    err := NewSubscribe().UpdateStatus(sub)
    if err != nil {
        snag.Panic(err)
    }

    // 判定退租是否满足条件
    if s.modifier == nil {
        if sub.Remaining < 0 {
            snag.Panic("欠费中, 无法继续办理")
        }
    } else {
        if sub.Remaining < 0 || sub.Status == model.SubscribeStatusOverdue {
            sub.Status = model.SubscribeStatusUsing
        }
    }

    if sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无法退租, 骑士卡当前非使用中")
    }

    s.do(business.TypeUnsubscribe, func(tx *ent.Tx) {
        var reason string
        if s.cabinet != nil {
            reason = "用户电柜退租"
        }
        if s.modifier != nil {
            reason = "管理员操作强制退租"
        }
        if s.employee != nil {
            reason = "店员操作退租"
        }

        _, err = tx.Subscribe.
            UpdateOneID(sub.ID).
            SetEndAt(time.Now()).
            SetStatus(model.SubscribeStatusUnSubscribed).
            SetUnsubscribeReason(reason).
            Save(s.ctx)
        snag.PanicIfError(err)

        // 标记需要签约
        _, err = tx.Rider.UpdateOneID(sub.RiderID).SetContractual(false).Save(s.ctx)
        snag.PanicIfError(err)

        // 查询并标记用户合同为失效
        _, _ = tx.Contract.Update().Where(contract.RiderID(sub.RiderID)).SetEffective(false).Save(s.ctx)
    })

    // 更新企业账单
    if sub.EnterpriseID != nil {
        go NewEnterprise().UpdateStatementByID(*sub.EnterpriseID)
    }
}

// Pause 寄存
func (s *businessRiderService) Pause(subscribeID uint64) {
    s.preprocess(business.TypePause, s.QuerySubscribeWithRider(subscribeID))

    if s.subscribe == nil || s.subscribe.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效订阅")
    }
    if s.subscribe.EnterpriseID != nil {
        snag.Panic("团签用户无法办理")
    }

    s.do(business.TypePause, func(tx *ent.Tx) {
        _, err := tx.SubscribePause.Create().
            SetStartAt(time.Now()).
            SetRiderID(s.subscribe.RiderID).
            SetSubscribeID(s.subscribe.ID).
            SetCityID(s.subscribe.CityID).
            SetNillableStoreID(s.storeID).
            SetNillableCabinetID(s.cabinetID).
            SetNillableEmployeeID(s.employeeID).Save(s.ctx)
        snag.PanicIfError(err)

        _, err = tx.Subscribe.UpdateOne(s.subscribe).
            SetPausedAt(time.Now()).
            SetStatus(model.SubscribeStatusPaused).
            Save(s.ctx)
        snag.PanicIfError(err)
    })
}

// Continue 继续计费
func (s *businessRiderService) Continue(subscribeID uint64) {
    s.preprocess(business.TypeContinue, s.QuerySubscribeWithRider(subscribeID))

    sp, _ := ent.Database.SubscribePause.QueryNotDeleted().
        Where(subscribepause.SubscribeID(s.subscribe.ID), subscribepause.EndAtIsNil()).
        Order(ent.Desc(subscribepause.FieldCreatedAt)).
        First(s.ctx)

    if sp == nil || s.subscribe.Status != model.SubscribeStatusPaused {
        snag.Panic("骑士卡状态错误")
    }

    // 当前时间
    now := time.Now()

    // 已寄存天数
    days, overdue := NewSubscribe().PausedDays(sp.StartAt, now)

    s.do(business.TypeContinue, func(tx *ent.Tx) {
        _, err := tx.SubscribePause.
            UpdateOne(sp).
            SetDays(days).
            SetEndAt(now).
            SetNillableEndEmployeeID(s.employeeID).
            SetEndModifier(s.modifier).
            SetOverdueDays(overdue).
            SetNillableEndStoreID(s.storeID).
            SetNillableEndCabinetID(s.cabinetID).
            Save(s.ctx)
        snag.PanicIfError(err)

        // 更新订阅
        _, err = tx.Subscribe.UpdateOne(s.subscribe).
            SetStatus(model.SubscribeStatusUsing).
            AddPauseDays(days).
            ClearPausedAt().
            Save(s.ctx)
        snag.PanicIfError(err)
    })
}
