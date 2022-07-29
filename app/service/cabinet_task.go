// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type cabinetTaskService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewCabinetTask() *cabinetTaskService {
    return &cabinetTaskService{
        ctx: context.Background(),
    }
}

func NewCabinetTaskWithRider(r *ent.Rider) *cabinetTaskService {
    s := NewCabinetTask()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewCabinetTaskWithModifier(m *model.Modifier) *cabinetTaskService {
    s := NewCabinetTask()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *cabinetTaskService) CreateExchange() {

}
