// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type riderPlanService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewRiderPlan() *riderPlanService {
    return &riderPlanService{
        ctx: context.Background(),
    }
}

func NewRiderPlanWithRider(rider *ent.Rider) *orderService {
    s := NewOrder()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewRiderPlanWithModifier(m *model.Modifier) *riderPlanService {
    s := NewRiderPlan()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}
