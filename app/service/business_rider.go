// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type businessRiderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewBusinessRider() *businessRiderService {
    return &businessRiderService{
        ctx: context.Background(),
    }
}

func (s *businessRiderService) Active() {
}
