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
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "strings"
    "time"
)

type riderMgrService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
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
    s.employeeInfo = &model.Employee{
        ID:    s.employee.ID,
        Name:  s.employee.Name,
        Phone: s.employee.Phone,
    }
    return s
}

func (s *riderMgrService) QuerySubscribeWithRider(subscribeID uint64) *ent.Subscribe {
    item, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(subscribeID)).WithRider().Only(s.ctx)
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

    var bls *businessLogService

    sub := s.QuerySubscribeWithRider(subscribeID)

    if sub == nil || sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效订阅")
    }
    if sub.EnterpriseID != nil {
        snag.Panic("团签用户无法办理")
    }

    tx, _ := ent.Database.Tx(s.ctx)

    spc := tx.SubscribePause.Create().
        SetStartAt(time.Now()).
        SetRiderID(sub.RiderID).
        SetSubscribeID(sub.ID)

    stockReq := &model.StockWithRiderReq{
        RiderID:   sub.RiderID,
        Model:     sub.Model,
        StockType: model.StockTypeRiderPause,
    }

    lg := logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetOperate(model.OperateSubscribePause).
        SetDiff("计费中", "暂停计费")

    if s.modifier != nil {
        stockReq.ManagerID = s.modifier.ID
        lg.SetModifier(s.modifier)
        bls = NewBusinessLogWithModifier(s.modifier, sub)
    }

    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID)
        spc.SetEmployee(s.employee)
        stockReq.EmployeeID = s.employee.ID
        stockReq.StoreID = s.employee.Edges.Store.ID
        lg.SetEmployee(s.employeeInfo)
        bls = NewBusinessLogWithEmployee(s.employee, sub)
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

    // 记录业务日志
    bls.SaveAsync(business.TypePause)
}

// ContinueSubscribe 继续计费
// TODO 后台操作出库怎么处理
func (s *riderMgrService) ContinueSubscribe(subscribeID uint64) {
    var bls *businessLogService

    sub := s.QuerySubscribeWithRider(subscribeID)

    sp, _ := ent.Database.SubscribePause.QueryNotDeleted().
        Where(subscribepause.SubscribeID(sub.ID), subscribepause.EndAtIsNil()).
        Order(ent.Desc(subscribepause.FieldCreatedAt)).
        First(s.ctx)

    if sp == nil || sub.Status != model.SubscribeStatusPaused {
        snag.Panic("无暂停中的骑士卡")
    }

    tx, _ := ent.Database.Tx(s.ctx)

    now := time.Now()

    // 已寄存天数
    days := NewSubscribe().PausedDays(sp.StartAt, now)

    spu := tx.SubscribePause.UpdateOne(sp).SetDays(days).SetEndAt(now)

    stockReq := &model.StockWithRiderReq{
        RiderID:   sub.RiderID,
        Model:     sub.Model,
        StockType: model.StockTypeRiderContinue,
    }

    lg := logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetOperate(model.OperateSubscribeContinue).
        SetDiff(fmt.Sprintf("暂停计费 (%s - %s 共%d天)", sp.StartAt.Format(carbon.DateTimeLayout), now.Format(carbon.DateTimeLayout), days), "计费中")

    if s.modifier != nil {
        stockReq.ManagerID = s.modifier.ID
        lg.SetModifier(s.modifier)
        bls = NewBusinessLogWithModifier(s.modifier, sub)
    }

    if s.employee != nil {
        NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID)

        spu.SetContinueEmployee(s.employee)
        stockReq.EmployeeID = s.employee.ID
        stockReq.StoreID = s.employee.Edges.Store.ID
        lg.SetEmployee(s.employeeInfo)

        bls = NewBusinessLogWithEmployee(s.employee, sub)
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

    // 记录业务日志
    bls.SaveAsync(business.TypeContinue)
}

// UnSubscribe 退租
// 会抹去欠费情况
// TODO 管理端强制退组库存如何操作
func (s *riderMgrService) UnSubscribe(subscribeID uint64) {
    var bls *businessLogService

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
        stockReq.ManagerID = s.modifier.ID

        bls = NewBusinessLogWithModifier(s.modifier, sub)
    }
    if s.employee != nil {
        reason = "店员操作退租"
        lgr.SetOperate(model.OperateUnsubscribe).SetEmployee(s.employeeInfo)
        stockReq.EmployeeID = s.employeeInfo.ID
        stockReq.StoreID = s.employee.Edges.Store.ID

        bls = NewBusinessLogWithEmployee(s.employee, sub)
    }

    tx, _ := ent.Database.Tx(s.ctx)

    _, err := tx.Subscribe.
        UpdateOneID(sub.ID).
        SetEndAt(time.Now()).
        SetRemaining(0).
        SetStatus(model.SubscribeStatusUnSubscribed).
        SetUnsubscribeReason(reason).
        Save(s.ctx)

    snag.PanicIfErrorX(err, tx.Rollback)

    // 电池入库
    snag.PanicIfErrorX(NewStockWithEmployee(s.employee).BatteryWithRider(
        tx.Stock.Create(),
        stockReq,
    ), tx.Rollback)

    _ = tx.Commit()

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
    bls.SaveAsync(business.TypeUnsubscribe)

    // 更新企业账单
    if sub.EnterpriseID != nil {
        go NewEnterprise().UpdateStatementByID(*sub.EnterpriseID)
    }
}

// Deposit 手动调整押金
func (s *riderMgrService) Deposit(req *model.RiderMgrDepositReq) {
    r := NewRider().Query(req.ID)
    o, _ := ent.Database.Order.QueryNotDeleted().
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
    tx, _ := ent.Database.Tx(s.ctx)
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

    tx, _ := ent.Database.Tx(s.ctx)

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

func (s *riderMgrService) QueryPhone(phone string) model.RiderEmployeeSearchRes {
    r, _ := ent.Database.Rider.QueryNotDeleted().WithPerson().Where(rider.Phone(phone)).WithEnterprise().First(s.ctx)
    if r == nil {
        snag.Panic("未找到骑手")
    }

    subd, _ := NewSubscribe().RecentDetail(r.ID)

    res := model.RiderEmployeeSearchRes{
        ID:              r.ID,
        Phone:           r.Phone,
        Overview:        NewExchange().Overview(r.ID),
        Status:          NewRider().Status(r),
        SubscribeStatus: subd.Status,
    }

    p := r.Edges.Person
    if p != nil {
        res.Name = p.Name
        res.AuthStatus = model.PersonAuthStatus(p.Status)
    }

    e := r.Edges.Enterprise
    if e != nil {
        res.Enterprise = &model.EnterpriseBasic{
            ID:   e.ID,
            Name: e.Name,
        }
    }

    res.Plan = subd.Plan

    return res
}
