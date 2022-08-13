// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-13
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type reserveService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewReserve() *reserveService {
    return &reserveService{
        ctx: context.Background(),
    }
}

func NewReserveWithRider(r *ent.Rider) *reserveService {
    s := NewReserve()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewReserveWithModifier(m *model.Modifier) *reserveService {
    s := NewReserve()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *reserveService) CabinetNum() {
}
