// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/pointlog"
    "github.com/auroraride/aurservd/internal/ent/rider"
)

type pointService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    orm          *ent.PointLogClient
}

func NewPoint() *pointService {
    return &pointService{
        ctx: context.Background(),
        orm: ent.Database.PointLog,
    }
}

func NewPointWithRider(r *ent.Rider) *pointService {
    s := NewPoint()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewPointWithModifier(m *model.Modifier) *pointService {
    s := NewPoint()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewPointWithEmployee(e *ent.Employee) *pointService {
    s := NewPoint()
    if e != nil {
        s.employee = e
        s.employeeInfo = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
        s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    }
    return s
}

// Modify 修改积分
func (s *pointService) Modify(req *model.PointModifyReq) error {
    r := NewRider().Query(req.RiderID)
    after := r.Points + req.Points
    if after < 0 {
        return errors.New("积分余额不能小于0")
    }
    return ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
        err = tx.Rider.UpdateOne(r).SetPoints(after).Exec(s.ctx)
        if err != nil {
            return
        }
        return tx.PointLog.Create().SetRiderID(req.RiderID).SetPoints(req.Points).SetReason(req.Reason).SetType(req.Type.Value()).SetAfter(after).Exec(s.ctx)
    })
}

// LogList 积分变动日志
func (s *pointService) LogList(req *model.PointLogListReq) *model.PaginationRes {
    q := s.orm.Query()
    if req.RiderID == 0 {
        if req.Keyword != "" {
            q.Where(
                pointlog.HasRiderWith(rider.Or(
                    rider.NameContainsFold(req.Keyword),
                    rider.PhoneContainsFold(req.Keyword),
                )),
            )
        }
    } else {
        q.Where(pointlog.RiderID(req.RiderID))
    }
    if req.Type != 0 {
        q.Where(pointlog.Type(req.Type.Value()))
    }
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.PointLog) model.PointLogListRes {
        var mp, mm string
        if item.Attach != nil {
            if item.Attach.Plan != nil {
                mp = fmt.Sprintf("%s-%d天", item.Attach.Plan.Name, item.Attach.Plan.Days)
            }
        }
        return model.PointLogListRes{
            ID:       item.ID,
            Type:     model.PointLogType(item.Type).String(),
            Plan:     mp,
            Points:   item.Points,
            Reason:   item.Reason,
            After:    item.After,
            Modifier: mm,
        }
    })
}

// Batch 批量发放积分
func (s *pointService) Batch(req *model.PointBatchReq) []string {
    riders, _, notfound := NewRider().QueryPhones(req.Phones)

    for i, p := range notfound {
        notfound[i] = fmt.Sprintf("未找到: %s", p)
    }

    for _, r := range riders {
        err := s.Modify(&model.PointModifyReq{
            RiderID: r.ID,
            Points:  req.Points,
            Reason:  req.Reason,
            Type:    req.Type,
        })
        if err != nil {
            notfound = append(notfound, fmt.Sprintf("%v: %s", err, r.Phone))
        }
    }

    return notfound
}
