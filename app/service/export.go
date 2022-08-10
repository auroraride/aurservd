// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type exportService struct{
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewExport() *exportService {
    return &exportService{
        ctx: context.Background(),
    }
}

func NewExportWithRider(r *ent.Rider) *exportService {
    s := NewExport()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewExportWithModifier(m *model.Modifier) *exportService {
    s := NewExport()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}
