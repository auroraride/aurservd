// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
)

type riderPermissionService struct {
    ctx   context.Context
    rider *ent.Rider
}

func NewRiderPermission() *riderPermissionService {
    return &riderPermissionService{
        ctx: context.Background(),
    }
}

func NewRiderPermissionWithRider(r *ent.Rider) *riderPermissionService {
    s := NewRiderPermission()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func (s *riderPermissionService) Business() (err error) {
    err = NewRider().Business(s.rider)
    if err != nil {
        return
    }

    // 检查是否企业用户, 企业状态
    if s.rider.EnterpriseID != nil {
        e := s.rider.Edges.Enterprise
        if e == nil {
            e = NewEnterprise().QueryX(*s.rider.EnterpriseID)
        }

        err = NewEnterprise().Business(e)
    }

    return
}

func (s *riderPermissionService) BusinessX() {
    err := s.Business()
    if err != nil {
        snag.Panic(err)
    }
}