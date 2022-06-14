// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    log "github.com/sirupsen/logrus"
    "strings"
    "time"
)

type riderMgrService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
}

func NewRiderMgr() *riderMgrService {
    return &riderMgrService{
        ctx: context.Background(),
    }
}

func NewRiderMgrWithRider(r *ent.Rider) *riderMgrService {
    s := NewRiderMgr()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewRiderMgrWithModifier(m *model.Modifier) *riderMgrService {
    s := NewRiderMgr()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewRiderMgrWithEmployee(e *ent.Employee) *riderMgrService {
    s := NewRiderMgr()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *riderMgrService) QuerySubscribeWithRider(subscribeID uint64) *ent.Subscribe {
    item, _ := ar.Ent.Subscribe.QueryNotDeleted().Where(subscribe.ID(subscribeID)).WithRider().Only(s.ctx)
    if item == nil {
        snag.Panic("未找到对应订阅")
    }
    return item
}

// PauseSubscribe 暂停计费
// TODO 后台操作入库怎么处理
func (s *riderMgrService) PauseSubscribe(subscribeID uint64) {
    if s.employee == nil && s.modifier == nil {
        snag.Panic("操作权限校验失败")
    }

    sub := s.QuerySubscribeWithRider(subscribeID)

    if sub == nil || sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效订阅")
    }
    if sub.EnterpriseID != nil {
        snag.Panic("团签用户无法办理")
    }

    tx, _ := ar.Ent.Tx(s.ctx)

    spc := tx.SubscribePause.Create().
        SetStartAt(time.Now()).
        SetRiderID(sub.RiderID).
        SetSubscribeID(sub.ID)

    stockReq := &model.StockWithRiderReq{
        RiderID:   sub.RiderID,
        Voltage:   sub.Voltage,
        StockType: model.StockTypeRiderPause,
    }

    lg := logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetOperate(model.OperateSubscribePause).
        SetDiff("计费中", "暂停计费")

    if s.modifier != nil {
        stockReq.ManagerID = s.modifier.ID
        lg.SetModifier(s.modifier)
    }

    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID)
        spc.SetEmployee(s.employee)
        stockReq.EmployeeID = s.employee.ID
        stockReq.StoreID = s.employee.Edges.Store.ID
        lg.SetEmployee(&model.Employee{
            ID:    s.employee.ID,
            Name:  s.employee.Name,
            Phone: s.employee.Phone,
        })
    }

    _, err := spc.Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _, err = tx.Subscribe.UpdateOne(sub).
        SetPausedAt(time.Now()).
        SetStatus(model.SubscribeStatusPaused).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 电池入库
    snag.PanicIfErrorX(NewStockWithEmployee(s.employee).BatteryWithRider(
        tx.Stock.Create(),
        stockReq,
    ), tx.Rollback)

    _ = tx.Commit()

    // 记录日志
    go lg.Send()
}

// ContinueSubscribe 继续计费
// TODO 后台操作出库怎么处理
func (s *riderMgrService) ContinueSubscribe(subscribeID uint64) {
    sub := s.QuerySubscribeWithRider(subscribeID)

    sp, _ := ar.Ent.SubscribePause.QueryNotDeleted().
        Where(subscribepause.SubscribeID(sub.ID), subscribepause.EndAtIsNil()).
        Order(ent.Desc(subscribepause.FieldCreatedAt)).
        First(s.ctx)

    if sp == nil || sub.Status != model.SubscribeStatusPaused {
        snag.Panic("无暂停中的骑士卡")
    }

    tx, _ := ar.Ent.Tx(s.ctx)

    now := time.Now()

    days := tools.NewTime().DiffDaysOfStart(now, sp.StartAt)

    spu := tx.SubscribePause.UpdateOne(sp).SetDays(days).SetEndAt(now)

    stockReq := &model.StockWithRiderReq{
        RiderID:   sub.RiderID,
        Voltage:   sub.Voltage,
        StockType: model.StockTypeRiderContinue,
    }

    lg := logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetOperate(model.OperateSubscribeContinue).
        SetDiff(fmt.Sprintf("暂停计费 (共%d天)", tools.NewTime().DiffDaysOfNextDayToNow(*sub.PausedAt)), "计费中")

    if s.modifier != nil {
        stockReq.ManagerID = s.modifier.ID
        lg.SetModifier(s.modifier)
    }

    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID)

        spu.SetContinueEmployee(s.employee)
        stockReq.EmployeeID = s.employee.ID
        stockReq.StoreID = s.employee.Edges.Store.ID
        lg.SetEmployee(&model.Employee{
            ID:    s.employee.ID,
            Name:  s.employee.Name,
            Phone: s.employee.Phone,
        })
    }

    _, err := spu.Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 更新订阅
    _, err = tx.Subscribe.UpdateOne(sub).
        SetStatus(model.SubscribeStatusUsing).
        AddPauseDays(days).
        ClearPausedAt().Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 电池出库
    snag.PanicIfErrorX(NewStockWithEmployee(s.employee).BatteryWithRider(
        tx.Stock.Create(),
        stockReq,
    ), tx.Rollback)

    _ = tx.Commit()

    // 记录日志
    go lg.Send()
}

// HaltSubscribe 强制退租
// 会抹去欠费情况
// TODO 是否适用于团签骑手
func (s *riderMgrService) HaltSubscribe(subscribeID uint64) {
    sub := s.QuerySubscribeWithRider(subscribeID)
    if sub == nil || sub.EndAt != nil {
        snag.Panic("未找到订阅")
    }

    _, err := sub.Update().SetEndAt(time.Now()).SetRemaining(0).SetStatus(model.SubscribeStatusCanceled).SetRemark("管理员操作强制退租").Save(s.ctx)
    if err != nil {
        log.Error(err)
        snag.Panic("操作失败")
    }

    before := fmt.Sprintf(
        "%s剩余天数: %d",
        model.SubscribeStatusText(sub.Status),
        sub.Remaining,
    )

    // 记录日志
    go logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetModifier(s.modifier).
        SetOperate(model.OperateSubscribePause).
        SetDiff(before, "强制退租").
        Send()
}

// Deposit 手动调整押金
func (s *riderMgrService) Deposit(req *model.RiderMgrDepositReq) {
    r := NewRider().Query(req.ID)
    o, _ := ar.Ent.Order.QueryNotDeleted().
        Where(
            order.RiderID(req.ID),
            order.Status(model.OrderStatusPaid),
            order.Type(model.OrderTypeDeposit),
            order.DeletedAtIsNil(),
        ).
        First(s.ctx)
    var before float64
    // 判断押金是否骑手自行缴纳
    if o != nil && o.Creator == nil {
        snag.Panic("用户有实际缴纳的押金订单, 无法继续修改")
    }

    var err error
    tx, _ := ar.Ent.Tx(s.ctx)
    if o != nil {
        before = o.Amount
        _, err = tx.Order.SoftDeleteOne(o).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)
    }

    if req.Amount > 0 {
        _, err = tx.Order.Create().
            SetRiderID(req.ID).
            SetType(model.OrderTypeDeposit).
            SetStatus(model.OrderStatusPaid).
            SetRemark("管理员修改").
            SetAmount(req.Amount).
            SetTotal(req.Amount).
            SetPayway(model.OrderPaywayManual).
            SetOutTradeNo(tools.NewUnique().NewSN28()).
            SetTradeNo(tools.NewUnique().NewSN28()).
            Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)
    }

    _ = tx.Commit()

    // 记录日志
    go logging.NewOperateLog().
        SetRef(r).
        SetModifier(s.modifier).
        SetOperate(model.OperateDeposit).
        SetDiff(fmt.Sprintf("%.2f元", before), fmt.Sprintf("%.2f元", req.Amount)).
        Send()
}

// Modify 修改骑手资料
func (s *riderMgrService) Modify(req *model.RiderMgrModifyReq) {
    if req.Contact == nil && req.Phone == nil && req.AuthStatus == nil {
        snag.Panic("参数错误")
    }

    tx, _ := ar.Ent.Tx(s.ctx)

    var err error
    r := NewRiderWithModifier(s.modifier).Query(req.ID)
    p := r.Edges.Person

    if p == nil && req.AuthStatus != nil {
        snag.Panic("用户还未提交个人信息")
    }

    pu := tx.Person.UpdateOne(p)
    ru := tx.Rider.UpdateOne(r)

    var before, after []string

    if req.Phone != nil {
        ru.SetPhone(*req.Phone)
        before = append(before, fmt.Sprintf("电话: %s", r.Phone))
        after = append(after, fmt.Sprintf("电话: %s", *req.Phone))
    }

    if req.Contact != nil {
        ru.SetContact(req.Contact)
        if r.Contact == nil {
            before = append(before, "联系人: 无")
        } else {
            before = append(before, fmt.Sprintf("联系人: %s, %s, %s", r.Contact.Relation, r.Contact.Phone, r.Contact.Name))
        }
        after = append(after, fmt.Sprintf("联系人: %s, %s, %s", req.Contact.Relation, req.Contact.Phone, req.Contact.Name))
    }

    if req.AuthStatus != nil {
        pu.SetStatus(req.AuthStatus.Raw())
        before = append(before, fmt.Sprintf("认证状态: %s", model.PersonAuthStatus(p.Status).String()))
        after = append(after, fmt.Sprintf("认证状态: %s", req.AuthStatus.String()))
    }

    _, err = pu.Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _, err = ru.Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 记录日志
    go logging.NewOperateLog().
        SetRef(r).
        SetModifier(s.modifier).
        SetOperate(model.OperateProfile).
        SetDiff(strings.Join(before, "\n"), strings.Join(after, "\n")).
        Send()
}
