// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type permissionService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewPermission() *permissionService {
    return &permissionService{
        ctx: context.Background(),
    }
}

func NewPermissionWithRider(r *ent.Rider) *permissionService {
    s := NewPermission()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewPermissionWithModifier(m *model.Modifier) *permissionService {
    s := NewPermission()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *permissionService) Create() {
}
