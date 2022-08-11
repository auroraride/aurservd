// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type riderCabinetService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.CabinetClient
}

func NewRiderCabinet(r *ent.Rider) *riderCabinetService {
    s := &riderCabinetService{
        ctx:   context.Background(),
        orm:   ent.Database.Cabinet,
        rider: r,
    }
    s.ctx = context.WithValue(s.ctx, "rider", r)
    return s
}
