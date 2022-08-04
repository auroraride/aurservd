// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
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
    "time"
)

type businessRiderService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    cabinet      *ent.Cabinet
}

func NewBusinessRider() *businessRiderService {
    return &businessRiderService{
        ctx: context.Background(),
    }
}

func NewBusinessRiderWithRider(r *ent.Rider, serial string) *businessRiderService {
    s := NewBusinessRider()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    s.cabinet = NewCabinet().QueryOneSerialX(serial)
    return s
}

func NewBusinessRiderWithModifier(m *model.Modifier) *businessRiderService {
    s := NewBusinessRider()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewBusinessRiderWithEmployee(e *ent.Employee) *businessRiderService {
    s := NewBusinessRider()
    s.ctx = context.WithValue(s.ctx, "employee", e)
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
    if p.Status != model.PersonAuthenticated.Raw() {
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

// Active 激活订阅
func (s *businessRiderService) Active(info *model.SubscribeActiveInfo, sub *ent.Subscribe) {
    var newsub *ent.Subscribe

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        var storeID, employeeID, cabinetID *uint64

        if s.employee != nil {
            employeeID = tools.NewPointerInterface(s.employee.ID)
            storeID = tools.NewPointerInterface(s.employee.Edges.Store.ID)
        }

        if s.cabinet != nil {
            cabinetID = tools.NewPointerInterface(s.cabinet.ID)
        }

        // 激活
        newsub, err = tx.Subscribe.UpdateOneID(info.ID).
            SetStatus(model.SubscribeStatusUsing).
            SetStartAt(time.Now()).
            SetNillableEmployeeID(employeeID).
            SetNillableStoreID(storeID).
            SetNillableCabinetID(cabinetID).
            Save(s.ctx)

        if err != nil {
            return
        }

        // 提成
        if info.CommissionID != nil && employeeID != nil {
            _, err = tx.Commission.UpdateOneID(*info.CommissionID).SetEmployeeID(*employeeID).Save(s.ctx)
            if err != nil {
                return
            }
        }

        // 调出库存
        return NewStockWithEmployee(s.employee).BatteryWithRider(
            tx.Stock.Create(),
            &model.StockWithRiderReq{
                RiderID:   info.Rider.ID,
                Model:     info.Model,
                StockType: model.StockTypeRiderObtain,

                StoreID:    storeID,
                EmployeeID: employeeID,
            },
        )
    })

    if info.EnterpriseID != nil {
        // 更新团签账单
        go NewEnterprise().UpdateStatement(sub.Edges.Enterprise)
    } else {
        // 更新个签订阅
        go NewSubscribe().UpdateStatus(newsub)
    }

    // 保存业务日志
    NewBusinessLog(sub).SetModifier(s.modifier).SetEmployee(s.employee).SaveAsync(business.TypeActive)
}

// UnSubscribe 退租
// 会抹去欠费情况
func (s *businessRiderService) UnSubscribe(subscribeID uint64) {
    sub := s.QuerySubscribeWithRider(subscribeID)
    if sub == nil || sub.EndAt != nil {
        snag.Panic("未找到订阅")
    }

    if sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无法退订, 骑士卡当前状态错误")
    }

    lgr := logging.NewOperateLog()
    stockReq := &model.StockWithRiderReq{
        RiderID:   sub.RiderID,
        Model:     sub.Model,
        StockType: model.StockTypeRiderUnSubscribe,
    }

    var reason string
    if s.modifier != nil {
        reason = "管理员操作强制退租"
        lgr.SetOperate(model.OperateHalt).SetModifier(s.modifier)
        stockReq.ManagerID = tools.NewPointerInterface(s.modifier.ID)
    }
    if s.employee != nil {
        reason = "店员操作退租"
        lgr.SetOperate(model.OperateUnsubscribe).SetEmployee(s.employeeInfo)
        stockReq.EmployeeID = tools.NewPointerInterface(s.employeeInfo.ID)
        stockReq.StoreID = tools.NewPointerInterface(s.employee.Edges.Store.ID)
    }

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        _, err = tx.Subscribe.
            UpdateOneID(sub.ID).
            SetEndAt(time.Now()).
            SetRemaining(0).
            SetStatus(model.SubscribeStatusUnSubscribed).
            SetUnsubscribeReason(reason).
            Save(s.ctx)
        snag.PanicIfError(err)

        // 标记需要签约
        _, err = tx.Rider.UpdateOneID(sub.RiderID).SetContractual(false).Save(s.ctx)
        snag.PanicIfError(err)

        // 电池入库
        return NewStockWithEmployee(s.employee).BatteryWithRider(
            tx.Stock.Create(),
            stockReq,
        )
    })

    // 查询并标记用户合同为失效
    _, _ = ent.Database.Contract.Update().Where(contract.RiderID(sub.RiderID)).SetEffective(false).Save(s.ctx)

    before := fmt.Sprintf(
        "%s剩余天数: %d",
        model.SubscribeStatusText(sub.Status),
        sub.Remaining,
    )

    // 记录日志
    go lgr.SetDiff(before, reason).Send()

    // 记录业务日志
    NewBusinessLog(sub).SetModifier(s.modifier).SetEmployee(s.employee).SaveAsync(business.TypeUnsubscribe)

    // 更新企业账单
    if sub.EnterpriseID != nil {
        go NewEnterprise().UpdateStatementByID(*sub.EnterpriseID)
    }
}

// PauseSubscribe 暂停计费
func (s *businessRiderService) PauseSubscribe(subscribeID uint64) {
    if s.employee == nil && s.modifier == nil && s.cabinet == nil {
        snag.Panic("操作权限校验失败")
    }

    sub := s.QuerySubscribeWithRider(subscribeID)

    if sub == nil || sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效订阅")
    }
    if sub.EnterpriseID != nil {
        snag.Panic("团签用户无法办理")
    }

    lg := logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetOperate(model.OperateSubscribePause).
        SetDiff("计费中", "暂停计费")

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        spc := tx.SubscribePause.Create().
            SetStartAt(time.Now()).
            SetRiderID(sub.RiderID).
            SetSubscribeID(sub.ID)

        stockReq := &model.StockWithRiderReq{
            RiderID:   sub.RiderID,
            Model:     sub.Model,
            StockType: model.StockTypeRiderPause,
        }

        if s.modifier != nil {
            stockReq.ManagerID = tools.NewPointerInterface(s.modifier.ID)
            lg.SetModifier(s.modifier)
        }

        if s.employee != nil {
            NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID)
            spc.SetEmployee(s.employee)
            stockReq.EmployeeID = tools.NewPointerInterface(s.employee.ID)
            stockReq.StoreID = tools.NewPointerInterface(s.employee.Edges.Store.ID)
            lg.SetEmployee(s.employeeInfo)
        }

        _, err = spc.Save(s.ctx)
        snag.PanicIfError(err)

        _, err = tx.Subscribe.UpdateOne(sub).
            SetPausedAt(time.Now()).
            SetStatus(model.SubscribeStatusPaused).
            Save(s.ctx)
        snag.PanicIfError(err)

        // 电池入库
        return NewStockWithEmployee(s.employee).BatteryWithRider(
            tx.Stock.Create(),
            stockReq,
        )
    })

    // 记录日志
    go lg.Send()

    // 记录业务日志
    NewBusinessLog(sub).
        SetModifier(s.modifier).
        SetEmployee(s.employee).
        SaveAsync(business.TypePause)
}

// ContinueSubscribe 继续计费
func (s *businessRiderService) ContinueSubscribe(subscribeID uint64) {
    sub := s.QuerySubscribeWithRider(subscribeID)

    sp, _ := ent.Database.SubscribePause.QueryNotDeleted().
        Where(subscribepause.SubscribeID(sub.ID), subscribepause.EndAtIsNil()).
        Order(ent.Desc(subscribepause.FieldCreatedAt)).
        First(s.ctx)

    if sp == nil || sub.Status != model.SubscribeStatusPaused {
        snag.Panic("无暂停中的骑士卡")
    }

    // 当前时间
    now := time.Now()

    // 已寄存天数
    days := NewSubscribe().PausedDays(sp.StartAt, now)

    lg := logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetOperate(model.OperateSubscribeContinue).
        SetDiff(fmt.Sprintf("暂停计费 (%s - %s 共%d天)", sp.StartAt.Format(carbon.DateTimeLayout), now.Format(carbon.DateTimeLayout), days), "计费中")

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        spu := tx.SubscribePause.UpdateOne(sp).SetDays(days).SetEndAt(now)

        stockReq := &model.StockWithRiderReq{
            RiderID:   sub.RiderID,
            Model:     sub.Model,
            StockType: model.StockTypeRiderContinue,
        }

        if s.modifier != nil {
            stockReq.ManagerID = tools.NewPointerInterface(s.modifier.ID)
            lg.SetModifier(s.modifier)
        }

        if s.employee != nil {
            NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID)

            spu.SetContinueEmployee(s.employee)
            stockReq.EmployeeID = tools.NewPointerInterface(s.employee.ID)
            stockReq.StoreID = tools.NewPointerInterface(s.employee.Edges.Store.ID)
            lg.SetEmployee(s.employeeInfo)
        }

        _, err = spu.Save(s.ctx)
        snag.PanicIfError(err)

        // 更新订阅
        _, err = tx.Subscribe.UpdateOne(sub).
            SetStatus(model.SubscribeStatusUsing).
            AddPauseDays(days).
            ClearPausedAt().
            Save(s.ctx)
        snag.PanicIfError(err)

        // 电池出库
        return NewStockWithEmployee(s.employee).BatteryWithRider(
            tx.Stock.Create(),
            stockReq,
        )
    })

    // 记录日志
    go lg.Send()

    // 记录业务日志
    NewBusinessLog(sub).SetModifier(s.modifier).SetEmployee(s.employee).SaveAsync(business.TypeContinue)
}
