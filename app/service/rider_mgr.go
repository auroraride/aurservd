// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "time"
)

type riderMgrService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
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

func NewRiderMgrWithEmployee(e *model.Employee) *riderMgrService {
    s := NewRiderMgr()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// PauseSubscribe 暂停计费
func (s *riderMgrService) PauseSubscribe(riderID uint64) {
    r := NewRider().Query(riderID)
    sub := NewSubscribeWithModifier(s.modifier).Recent(riderID)
    if sub == nil || sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效订阅")
    }
    if sub.EnterpriseID != nil {
        snag.Panic("团签用户无法办理")
    }
    tx, _ := ar.Ent.Tx(s.ctx)
    _, err := tx.SubscribePause.Create().
        SetStartAt(time.Now()).
        SetRemark("管理员操作").
        SetRiderID(riderID).
        SetSubscribeID(sub.ID).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _, err = tx.Subscribe.UpdateOne(sub).
        SetPausedAt(time.Now()).
        SetStatus(model.SubscribeStatusPaused).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _ = tx.Commit()

    // 记录日志
    go logging.NewOperateLog().
        SetRef(r).
        SetModifier(s.modifier).
        SetOperate(logging.OperateSubscribePause).
        SetDiff("计费中", "暂停计费").
        Send()
}

// ContinueSubscribe 继续计费
func (s *riderMgrService) ContinueSubscribe(riderID uint64) {
    r := NewRider().Query(riderID)

    sp, _ := ar.Ent.SubscribePause.QueryNotDeleted().
        Where(subscribepause.RiderID(riderID), subscribepause.EndAtIsNil()).
        Order(ent.Desc(subscribepause.FieldCreatedAt)).
        WithSubscribe().
        First(s.ctx)

    if sp == nil || sp.Edges.Subscribe == nil || sp.Edges.Subscribe.Status != model.SubscribeStatusPaused {
        snag.Panic("无暂停中的订阅")
    }

    tx, err := ar.Ent.Tx(s.ctx)
    if err != nil {
        snag.Panic(err)
    }

    now := time.Now()

    days := tools.NewTime().DiffDaysOfStart(now, sp.StartAt)
    _, err = tx.SubscribePause.UpdateOne(sp).SetDays(days).SetEndAt(now).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 更新订阅
    _, err = tx.Subscribe.UpdateOne(sp.Edges.Subscribe).
        SetStatus(model.SubscribeStatusUsing).
        AddPauseDays(days).
        ClearPausedAt().
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _ = tx.Commit()

    // 记录日志
    go logging.NewOperateLog().
        SetRef(r).
        SetModifier(s.modifier).
        SetOperate(logging.OperateSubscribePause).
        SetDiff("暂停计费", "计费中").
        Send()

}
