// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/internal/ent/subscribesuspend"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
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
        var pauseID *uint64
        // 判断是否寄存中, 如果是寄存中的话时间计算为当日寄存生效时间
        if pause != nil {
            // 判断当日0点是否是寄存状态
            beginning := carbon.Now().StartOfDay().Carbon2Time()
            if beginning.After(pause.StartAt) {
                now = beginning
            }
            pauseID = tools.Pointer(pause.ID)
        }
        _, err = tx.SubscribeSuspend.Create().SetStartAt(time.Now()).SetRemark(req.Remark).SetStartAt(now).SetSubscribeID(sub.ID).SetCityID(sub.CityID).SetRiderID(sub.RiderID).SetNillablePauseID(pauseID).Save(s.ctx)
        if err != nil {
            snag.Panic("暂停失败")
        }

        return tx.Subscribe.UpdateOne(sub).SetSuspendAt(now).Exec(s.ctx)
    })

    // 记录日志
    go logging.NewOperateLog().
        SetRef(NewRider().Query(sub.RiderID)).
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
        days, _, _, _ := sus.GetAdditionalDays()
        err = tx.SubscribeSuspend.UpdateOne(sus).SetDays(days).SetEndAt(now).SetEndModifier(s.modifier).SetEndReason(req.Remark).Exec(s.ctx)
        if err != nil {
            snag.Panic("继续计费操作失败")
        }
        return tx.Subscribe.UpdateOne(sub).ClearSuspendAt().AddSuspendDays(days).Exec(s.ctx)
    })

    // 记录日志
    go logging.NewOperateLog().
        SetRef(NewRider().Query(sub.RiderID)).
        SetOperate(model.OperateResume).
        SetModifier(s.modifier).
        SetDiff("暂停扣费", model.SubscribeStatusText(sub.Status)).
        Send()
}

func (s *suspendService) listFilter(req model.SuspendListFilter) (q *ent.SubscribeSuspendQuery, info ar.Map) {
    info = make(ar.Map)
    q = s.orm.Query().WithRider(func(query *ent.RiderQuery) {
        query.WithPerson()
    }).WithCity().WithSubscribe(func(query *ent.SubscribeQuery) {
        query.WithPlan()
    }).Order(ent.Desc(subscribesuspend.FieldStartAt))
    if req.CityID != 0 {
        q.Where(subscribesuspend.CityID(req.CityID))
        info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
    }
    if req.RiderID != 0 {
        q.Where(subscribesuspend.RiderID(req.RiderID))
        info["骑手"] = ent.NewExportInfo(req.RiderID, rider.Table)
    }
    if req.Start != "" {
        q.Where(subscribesuspend.StartAt(tools.NewTime().ParseDateStringX(req.Start)))
        info["开始日期"] = req.Start
    }
    if req.End != "" {
        q.Where(
            subscribesuspend.Or(
                subscribesuspend.EndAtIsNil(),
                subscribesuspend.EndAtLT(tools.NewTime().ParseNextDateStringX(req.End)),
            ),
        )
        info["结束日期"] = req.End
    }
    return
}

func (s *suspendService) List(req *model.SuspendListReq) *model.PaginationRes {
    q, _ := s.listFilter(req.SuspendListFilter)
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.SubscribeSuspend) (res model.SuspendListRes) {
        status := "暂停中"
        if !item.EndAt.IsZero() {
            status = "已结束"
        }
        sub := item.Edges.Subscribe
        res = model.SuspendListRes{
            City:            item.Edges.City.Name,
            Name:            item.Edges.Rider.Edges.Person.Name,
            Phone:           item.Edges.Rider.Phone,
            Plan:            item.Edges.Subscribe.Edges.Plan.Name,
            Status:          status,
            SubscribeDays:   sub.Remaining,
            SubscribeStatus: model.SubscribeStatusText(sub.Status),
            Days:            item.Days,
            Start:           item.StartAt.Format(carbon.DateTimeLayout),
        }
        if item.Creator != nil {
            res.StartBy = item.Creator.Name
        }
        if !item.EndAt.IsZero() {
            res.End = item.EndAt.Format(carbon.DateTimeLayout)
            if item.EndModifier != nil {
                res.EndBy = item.EndModifier.Name
            }
        }
        return
    })
}
