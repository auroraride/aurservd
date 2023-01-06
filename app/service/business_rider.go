// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/adapter/async"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/internal/ent/commission"
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/internal/ent/orderrefund"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/silk"
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

    task func() (*model.BinInfo, *model.Battery, error) // 电柜任务

    storeID, employeeID, cabinetID, subscribeID *uint64

    // 电车信息
    ebikeInfo *model.EbikeBusinessInfo
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

// SetCabinet 设置电柜
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

// SetCabinetID 设置电柜
func (s *businessRiderService) SetCabinetID(id *uint64) *businessRiderService {
    if id != nil {
        s.SetCabinet(NewCabinet().QueryOne(*id))
    }
    return s
}

// SetStoreID 设置门店
func (s *businessRiderService) SetStoreID(id *uint64) *businessRiderService {
    if id != nil {
        s.store = NewStore().Query(*id)
    }
    return s
}

func (s *businessRiderService) SetEmployeeID(id *uint64) *businessRiderService {
    if id != nil {
        s.employee, _ = NewEmployee().Query(*id)
        if s.employee != nil {
            s.employeeID = id
            s.employeeInfo = &model.Employee{
                ID:    s.employee.ID,
                Name:  s.employee.Name,
                Phone: s.employee.Phone,
            }
        }
    }
    return s
}

// SetEbikeID 设置电车
func (s *businessRiderService) SetEbikeID(id *uint64) *businessRiderService {
    if id == nil {
        return s
    }
    bike, _ := ent.Database.Ebike.Query().Where(ebike.ID(*id)).WithBrand().First(s.ctx)
    if bike == nil || bike.Edges.Brand == nil {
        snag.Panic("电车信息查询失败")
    }
    brand := bike.Edges.Brand
    s.ebikeInfo = &model.EbikeBusinessInfo{
        ID:        bike.ID,
        BrandID:   brand.ID,
        BrandName: brand.Name,
    }

    return s
}

func (s *businessRiderService) SetEbike(info *model.EbikeBusinessInfo) *businessRiderService {
    s.ebikeInfo = info
    return s
}

func (s *businessRiderService) SetTask(task func() (*model.BinInfo, *model.Battery, error)) *businessRiderService {
    if task != nil {
        s.task = task
    }
    return s
}

func NewBusinessRiderWithEmployee(e *ent.Employee) *businessRiderService {
    s := NewBusinessRider(nil)
    if e == nil {
        snag.Panic("店员错误")
    }
    store := e.Edges.Store
    if store == nil {
        snag.Panic("未找到所属门店")
    }
    s.store = store
    s.employee = e
    s.employeeInfo = &model.Employee{
        ID:    e.ID,
        Name:  e.Name,
        Phone: e.Phone,
    }
    s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    return s
}

// QuerySubscribeWithRider 查询订阅信息
func (s *businessRiderService) QuerySubscribeWithRider(subscribeID uint64) *ent.Subscribe {
    item, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(subscribeID)).WithEnterprise().WithRider().Only(s.ctx)
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
        WithRider().
        WithEnterprise().
        WithCity().
        WithBrand().
        First(s.ctx)

    if sub == nil {
        snag.Panic("未找到待激活骑士卡")
    }

    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID, s.store)
    }

    r := sub.Edges.Rider
    if r == nil {
        snag.Panic("骑手信息获取失败")
    }

    res := &model.SubscribeActiveInfo{
        ID:           sub.ID,
        EnterpriseID: sub.EnterpriseID,
        Model:        sub.Model,
        CommissionID: nil,
        Rider: model.Rider{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  r.Name,
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
        en := sub.Edges.Enterprise
        res.Enterprise = &model.Enterprise{
            ID:    en.ID,
            Name:  en.Name,
            Agent: en.Agent,
        }
    }

    if sub.BrandID != nil {
        brand := sub.Edges.Brand
        if brand == nil {
            snag.Panic("电车型号查询失败")
        }
        res.EbikeBrand = &model.EbikeBrand{
            ID:    brand.ID,
            Name:  brand.Name,
            Cover: brand.Cover,
        }
    }

    return res, sub
}

// preprocess 预处理数据
func (s *businessRiderService) preprocess(bt business.Type, sub *ent.Subscribe) {
    s.subscribe = sub

    if sub.EnterpriseID != nil {
        en := sub.Edges.Enterprise
        if en == nil {
            snag.Panic("未找到团签信息")
        }
        // 判定是否寄存或取消寄存业务
        if bt == business.TypePause || bt == business.TypeContinue {
            snag.Panic("团签用户无法办理")
        }
        // 判定代理是否可使用门店
        if en.Agent && !en.UseStore && s.employee != nil {
            snag.Panic("代理无法在门店办理业务")
        }
    }

    s.subscribeID = silk.Pointer(sub.ID)

    r := sub.Edges.Rider
    if r == nil {
        r, _ = sub.QueryRider().First(s.ctx)
    }

    if r == nil {
        snag.Panic("骑手查询失败")
    }

    // 骑士卡状态
    if !NewRiderBusiness(r).Executable(sub, bt) {
        snag.Panic("骑士卡状态错误")
    }

    // 检查用户是否可以办理业务
    NewRiderPermissionWithRider(r, s.modifier).BusinessX().SubscribeX(model.RiderPermissionTypeBusiness, sub)

    s.rider = r

    if s.store != nil {
        s.storeID = silk.Pointer(s.store.ID)
    }

    if s.cabinet != nil {
        s.cabinetID = silk.Pointer(s.cabinet.ID)
    }

    if s.store == nil && s.cabinet == nil {
        snag.Panic("条件不满足")
    }

    if s.employee != nil {
        s.employeeID = silk.Pointer(s.employee.ID)
    }

    if s.employee == nil && s.store == nil && s.modifier == nil && s.cabinet == nil {
        snag.Panic("操作权限校验失败")
    }

    // 校验权限
    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(s.subscribe.CityID, s.store)
    }

    // 车电订阅检查
    if sub.BrandID != nil {
        // 车电订阅无法办理寄存相关业务
        // 车电订阅无法使用电柜
        if (bt != business.TypeActive && bt != business.TypeUnsubscribe) || s.cabinetID != nil {
            snag.Panic("车电订阅无法办理此业务")
        }
    }

    // 如果是车电订阅, 查询并设置电车信息
    if sub.EbikeID != nil && s.ebikeInfo == nil {
        s.SetEbikeID(sub.EbikeID)
    }

    // 预约检查
    rev := NewReserveWithRider(r).RiderUnfinished(r.ID)
    if rev != nil {
        if s.cabinet == nil || s.cabinet.ID != rev.CabinetID || bt.String() != rev.Type {
            _, _ = rev.Update().SetStatus(model.ReserveStatusInvalid.Value()).Save(s.ctx)
        } else {
            // 预约处理中
            s.reserve, _ = rev.Update().SetStatus(model.ReserveStatusProcessing.Value()).Save(s.ctx)
        }
    }
}

// doTask 处理电柜任务
func (s *businessRiderService) doTask() (bin *model.BinInfo, bat *model.Battery, err error) {
    defer func() {
        if v := recover(); v != nil {
            err = fmt.Errorf("%v", v)
        }
    }()

    bin, bat, err = s.task()
    return
}

// do 处理业务
func (s *businessRiderService) do(bt business.Type, cb func(tx *ent.Tx)) {
    async.WithTask(func() {
        sts := map[business.Type]uint8{
            business.TypeActive:      model.StockTypeRiderActive,
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

        var bin *model.BinInfo
        var err error

        // 放入电池优先执行
        var bat *model.Battery
        if s.task != nil && (bt == business.TypePause || bt == business.TypeUnsubscribe) {
            bin, bat, err = s.doTask()
            if err != nil {
                snag.Panic(err)
            }
        }

        // 激活业务查找提成
        var co *ent.Commission
        if bt == business.TypeActive {
            co, _ = ent.Database.Commission.QueryNotDeleted().Where(commission.SubscribeID(s.subscribe.ID)).First(s.ctx)
        }

        // 库存管理
        // TODO 智能电池
        var sk *ent.Stock
        ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
            cb(tx)

            sk, err = NewStockWithModifier(s.modifier).RiderBusiness(
                tx,
                &model.StockBusinessReq{
                    RiderID:   s.subscribe.RiderID,
                    Model:     s.subscribe.Model,
                    CityID:    s.subscribe.CityID,
                    StockType: sts[bt],

                    StoreID:     s.storeID,
                    EmployeeID:  s.employeeID,
                    CabinetID:   s.cabinetID,
                    SubscribeID: s.subscribeID,

                    Ebike:   s.ebikeInfo,
                    Battery: bat,
                },
            )

            if err != nil {
                log.Errorf("骑手业务出入库失败: %v", err)
            }

            return err
        })

        // 取出电池滞后执行
        if s.task != nil && (bt == business.TypeActive || bt == business.TypeContinue) {
            bin, bat, err = s.doTask()
            if err != nil {
                log.Error(err)
            }
            if bat != nil && s.cabinet.Intelligent {
                _ = sk.Update().SetBatteryID(bat.ID).Exec(s.ctx)
            }
        }

        // 保存业务日志
        var b *ent.Business
        b, err = NewBusinessLog(s.subscribe).
            SetModifier(s.modifier).
            SetEmployee(s.employee).
            SetCabinet(s.cabinet).
            SetStore(s.store).
            SetBinInfo(bin).
            SetStock(sk).
            SetBattery(bat).
            Save(bt)
        var bussinessID *uint64
        revStatus := model.ReserveStatusFail
        if b != nil {
            revStatus = model.ReserveStatusSuccess
            bussinessID = silk.Pointer(b.ID)
        }

        // 更新预约
        if s.reserve != nil {
            _, _ = s.reserve.Update().
                SetStatus(revStatus.Value()).
                SetNillableBusinessID(bussinessID).
                Save(s.ctx)
        }

        // 更新提成
        if bt == business.TypeActive && co != nil && b != nil && s.employeeID != nil {
            _, _ = co.Update().SetBusiness(b).SetEmployeeID(*s.employeeID).Save(s.ctx)
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
    })
}

// Active 激活订阅
func (s *businessRiderService) Active(sub *ent.Subscribe, allo *ent.Allocate) {
    s.preprocess(business.TypeActive, sub)

    if NewSubscribe().NeedContract(sub) {
        snag.Panic("还未签约, 请签约")
    }

    s.do(business.TypeActive, func(tx *ent.Tx) {
        var err error

        // 更新分配
        err = tx.Allocate.UpdateOne(allo).SetStatus(model.AllocateStatusSigned.Value()).Exec(s.ctx)
        if err != nil {
            return
        }

        var (
            aend *time.Time
        )
        // 如果是代理商, 计算骑士卡代理商结束时间
        if sub.EnterpriseID != nil {
            if sub.Edges.Enterprise == nil {
                sub.Edges.Enterprise = sub.QueryEnterprise().FirstX(s.ctx)
            }
            if sub.Edges.Enterprise == nil {
                snag.Panic("未找到团签信息")
            }
            if sub.Edges.Enterprise.Agent {
                aend = silk.Pointer(tools.NewTime().WillEnd(time.Now(), sub.InitialDays))
            }
        }

        updater := tx.Subscribe.UpdateOneID(sub.ID).
            SetStatus(model.SubscribeStatusUsing).
            SetStartAt(time.Now()).
            UpdateTarget(s.cabinetID, s.storeID, s.employeeID).
            SetNillableAgentEndAt(aend).
            SetNeedContract(false)

        // 设置订阅电车
        if s.ebikeInfo != nil {
            updater.SetEbikeID(s.ebikeInfo.ID).SetBrandID(s.ebikeInfo.BrandID)
        }

        // 激活
        s.subscribe, err = updater.Save(s.ctx)
        snag.PanicIfError(err)

        // 更新电车
        if s.ebikeInfo != nil {
            // 更新电车所属
            err = tx.Ebike.UpdateOneID(s.ebikeInfo.ID).SetRiderID(sub.RiderID).SetStatus(model.EbikeStatusUsing).Exec(s.ctx)
        }

        // 更新退款失效
        if sub.EnterpriseID == nil && sub.InitialOrderID != 0 {
            of, _ := tx.OrderRefund.QueryNotDeleted().Where(orderrefund.OrderID(sub.InitialOrderID)).First(s.ctx)
            if of != nil {
                err = tx.OrderRefund.UpdateOne(of).SetReason("激活订阅, 自动拒绝退款").SetStatus(model.OrderStatusRefundRefused).Exec(s.ctx)
            }
        }
    })

    if sub.EnterpriseID != nil {
        // 更新团签账单
        go NewEnterprise().UpdateStatement(sub.Edges.Enterprise)
    }
}

// UnSubscribe 退租
// 会抹去欠费情况
func (s *businessRiderService) UnSubscribe(subscribeID uint64, fns ...func(sub *ent.Subscribe)) {

    s.preprocess(business.TypeUnsubscribe, s.QuerySubscribeWithRider(subscribeID))

    if len(fns) > 0 {
        fns[0](s.subscribe)
    }

    sub := s.subscribe

    err := NewSubscribe().UpdateStatus(sub, false)
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
        _, err = tx.Rider.UpdateOneID(sub.RiderID).Save(s.ctx)
        snag.PanicIfError(err)

        // 查询并标记用户合同为失效
        _, err = tx.Contract.Update().Where(contract.RiderID(sub.RiderID)).SetEffective(false).Save(s.ctx)
        snag.PanicIfError(err)

        // 更新电车
        if sub.EbikeID != nil {
            // 删除电车所属
            err = tx.Ebike.UpdateOneID(*sub.EbikeID).ClearRiderID().SetStatus(model.EbikeStatusInStock).SetNillableStoreID(s.storeID).Exec(s.ctx)
        }
    })

    // 更新企业账单
    if sub.EnterpriseID != nil {
        go NewEnterprise().UpdateStatementByID(*sub.EnterpriseID)
    }
}

// Pause 寄存
func (s *businessRiderService) Pause(subscribeID uint64) {
    s.preprocess(business.TypePause, s.QuerySubscribeWithRider(subscribeID))

    if s.subscribe.Remaining < 1 {
        snag.Panic("当前剩余时间不足, 无法寄存")
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

// Continue 取消寄存
func (s *businessRiderService) Continue(subscribeID uint64) {
    s.preprocess(business.TypeContinue, s.QuerySubscribeWithRider(subscribeID))

    // 更新订阅信息
    err := NewSubscribe().UpdateStatus(s.subscribe, false)
    if err != nil {
        log.Error(err)
        snag.Panic("骑士卡更新失败")
    }

    pr, _ := s.subscribe.GetAdditionalItems()

    sp, _ := ent.Database.SubscribePause.QueryNotDeleted().
        Where(subscribepause.SubscribeID(s.subscribe.ID), subscribepause.EndAtIsNil()).
        Order(ent.Desc(subscribepause.FieldCreatedAt)).
        First(s.ctx)

    if sp == nil || pr.Current == nil || pr.Current.ID != sp.ID {
        snag.Panic("未找到寄存信息")
    }

    // 当前时间
    now := time.Now()

    s.do(business.TypeContinue, func(tx *ent.Tx) {
        _, err = tx.SubscribePause.
            UpdateOne(sp).
            SetDays(pr.CurrentDays).
            SetEndAt(now).
            SetNillableEndEmployeeID(s.employeeID).
            SetEndModifier(s.modifier).
            SetOverdueDays(pr.CurrentOverdueDays).
            SetNillableEndStoreID(s.storeID).
            SetNillableEndCabinetID(s.cabinetID).
            SetSuspendDays(pr.CurrentDuplicateDays).
            Save(s.ctx)
        snag.PanicIfError(err)

        // 更新订阅
        _, err = tx.Subscribe.UpdateOne(s.subscribe).
            SetStatus(model.SubscribeStatusUsing).
            SetPauseDays(pr.Days).
            ClearPausedAt().
            Save(s.ctx)
        snag.PanicIfError(err)
    })
}
