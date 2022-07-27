// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type cabinetBusinessService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewCabinetBusiness() *cabinetBusinessService {
    return &cabinetBusinessService{
        ctx: context.Background(),
    }
}

func NewCabinetBusinessWithRider(r *ent.Rider) *cabinetBusinessService {
    s := NewCabinetBusiness()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewCabinetBusinessWithModifier(m *model.Modifier) *cabinetBusinessService {
    s := NewCabinetBusiness()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}
