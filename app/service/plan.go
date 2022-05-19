// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/golang-module/carbon/v2"
    "time"
)

type planService struct {
    ctx context.Context
    orm *ent.PlanClient
}

func NewPlan() *planService {
    return &planService{
        ctx: context.Background(),
        orm: ar.Ent.Plan,
    }
}

// CreatePlan 创建骑士卡
func (s *planService) CreatePlan(m *model.Modifier, req *model.PlanCreateReq) {
    start, _ := time.Parse(carbon.DateLayout, req.Start)
    end, _ := time.Parse(carbon.DateLayout, req.End)
    s.orm.Create().
        AddCityIDs(req.Cities...).
        AddPmIDs(req.Models...).
        SetStart(start).
        SetEnd(end).
        SetDays(req.Days).
        SetEnable(req.Enable).
        SetName(req.Name).
        SetPrice(req.Price).
        SetCommission(req.Commission).
        SetLastModifier(m).
        SaveX(s.ctx)
}

func (s *planService) UpdateEnable() {
}
