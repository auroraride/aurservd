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
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    log "github.com/sirupsen/logrus"
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

func (s *riderMgrService) QuerySubscribeWithRider(subscribeID uint64) *ent.Subscribe {
    item, _ := ar.Ent.Subscribe.QueryNotDeleted().Where(subscribe.ID(subscribeID)).WithRider().Only(s.ctx)
    if item == nil {
        snag.Panic("未找到对应订阅")
    }
    return item
}

// PauseSubscribe 暂停计费
func (s *riderMgrService) PauseSubscribe(subscribeID uint64) {
    sub := s.QuerySubscribeWithRider(subscribeID)
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
        SetRiderID(sub.RiderID).
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
        SetRef(sub.Edges.Rider).
        SetModifier(s.modifier).
        SetOperate(logging.OperateSubscribePause).
        SetDiff("计费中", "暂停计费").
        Send()
}

// ContinueSubscribe 继续计费
func (s *riderMgrService) ContinueSubscribe(subscribeID uint64) {
    sub := s.QuerySubscribeWithRider(subscribeID)

    sp, _ := ar.Ent.SubscribePause.QueryNotDeleted().
        Where(subscribepause.SubscribeID(sub.ID), subscribepause.EndAtIsNil()).
        Order(ent.Desc(subscribepause.FieldCreatedAt)).
        First(s.ctx)

    if sp == nil || sub.Status != model.SubscribeStatusPaused {
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
    _, err = tx.Subscribe.UpdateOne(sub).
        SetStatus(model.SubscribeStatusUsing).
        AddPauseDays(days).
        ClearPausedAt().
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _ = tx.Commit()

    // 记录日志
    go logging.NewOperateLog().
        SetRef(sub.Edges.Rider).
        SetModifier(s.modifier).
        SetOperate(logging.OperateSubscribePause).
        SetDiff("暂停计费", "计费中").
        Send()
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
        SetOperate(logging.OperateSubscribePause).
        SetDiff(before, "强制退租").
        Send()
}
