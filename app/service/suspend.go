// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/internal/ent/subscribesuspend"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "time"
)

type suspendService struct {
    ctx      context.Context
    modifier *model.Modifier
    orm      *ent.SubscribeSuspendClient
}

func NewSuspend() *suspendService {
    return &suspendService{
        ctx: context.Background(),
        orm: ent.Database.SubscribeSuspend,
    }
}

func NewSuspendWithModifier(m *model.Modifier) *suspendService {
    s := NewSuspend()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *suspendService) Current(subscribeID uint64) *ent.SubscribeSuspend {
    sus, _ := s.orm.Query().Where(subscribesuspend.SubscribeID(subscribeID), subscribesuspend.EndAtIsNil()).First(s.ctx)
    return sus
}

// Suspend 暂停扣费
// 当骑士卡暂停时无法办理包含寄存在内的任何业务
// 当寄存中时, 可以办理暂停计费, 此时寄存天数需要减去实际发生的暂停天数
func (s *suspendService) Suspend(req *model.SuspendReq) {
    sub := NewSubscribeWithModifier(s.modifier).QueryX(req.ID)
    if sub.SuspendAt != nil {
        snag.Panic("已经处于暂停中")
    }

    if sub.EnterpriseID != nil {
        snag.Panic("团签用户无法操作")
    }

    pause, _ := ent.Database.SubscribePause.QueryNotDeleted().Where(subscribepause.SubscribeID(sub.ID), subscribepause.EndAtIsNil()).First(s.ctx)

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        now := time.Now()
        // 判断是否寄存中, 如果是寄存中的话时间计算为当日寄存生效时间
        if pause != nil {
            // 判断当日0点是否是寄存状态
            beginning := carbon.Now().StartOfDay().Carbon2Time()
            if beginning.After(pause.StartAt) {
                now = beginning
            }
        }
        _, err = tx.SubscribeSuspend.Create().SetStartAt(time.Now()).SetRemark(req.Remark).SetStartAt(now).SetSubscribeID(sub.ID).SetCityID(sub.CityID).SetRiderID(sub.RiderID).SetPause(pause).Save(s.ctx)
        if err != nil {
            snag.Panic("暂停失败")
        }

        return tx.Subscribe.UpdateOne(sub).SetSuspendAt(now).Exec(s.ctx)
    })

    // 记录日志
    go logging.NewOperateLog().
        SetRef(NewRider().Query(sub.ID)).
        SetOperate(model.OperateSuspend).
        SetModifier(s.modifier).
        SetDiff(model.SubscribeStatusText(sub.Status), "暂停扣费").
        Send()
}

// UnSuspend 取消暂停扣费
func (s *suspendService) UnSuspend(req *model.SuspendReq) {
    sub := NewSubscribeWithModifier(s.modifier).QueryX(req.ID)
    if sub.EnterpriseID != nil {
        snag.Panic("团签用户无法操作")
    }

    sus := s.Current(req.ID)
    if sus == nil {
        snag.Panic("未处于暂停中")
    }

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        now := time.Now()
        days, _, _ := sus.GetAdditionalDays()
        err = tx.SubscribeSuspend.UpdateOne(sus).SetDays(days).SetEndAt(now).SetEndModifier(s.modifier).SetEndReason(req.Remark).Exec(s.ctx)
        if err != nil {
            snag.Panic("继续计费操作失败")
        }
        return tx.Subscribe.UpdateOne(sub).ClearSuspendAt().AddSuspendDays(days).Exec(s.ctx)
    })

    // 记录日志
    go logging.NewOperateLog().
        SetRef(NewRider().Query(sub.ID)).
        SetOperate(model.OperateResume).
        SetModifier(s.modifier).
        SetDiff("暂停扣费", model.SubscribeStatusText(sub.Status)).
        Send()
}
