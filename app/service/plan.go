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
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/pkg/snag"
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

// Query 查找骑士卡
func (s *planService) Query(id uint64) *ent.Plan {
    item, err := s.orm.QueryNotDeleted().Where(plan.ID(id)).Only(s.ctx)
    if err != nil || item == nil {
        snag.Panic("未找到有效的骑士卡")
    }
    return item
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

// UpdateEnable 修改骑士卡状态
func (s *planService) UpdateEnable(m *model.Modifier, req *model.PlanEnableModifyReq) {
    item := s.Query(req.ID)
    s.orm.UpdateOne(item).SetEnable(*req.Enable).SetLastModifier(m).SaveX(s.ctx)
}

// Delete 软删除骑士卡
func (s *planService) Delete(m *model.Modifier, req *model.IDParamReq) {
    item := s.Query(req.ID)
    s.orm.UpdateOne(item).SetDeletedAt(time.Now()).SetLastModifier(m).SaveX(s.ctx)
}
