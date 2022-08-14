// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
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

// Business 判定骑手业务状态是否满足条件
func (s *riderPermissionService) Business() (err error) {
    err = NewRider().Permission(s.rider)
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

func (s *riderPermissionService) BusinessX() *riderPermissionService {
    err := s.Business()
    if err != nil {
        snag.Panic(err)
    }
    return s
}

// SubscribeX 检查骑士卡权限
// TODO 暂停扣费中 -> 骑士卡暂停中
func (s *riderPermissionService) SubscribeX(typ model.RiderPermissionType, sub *ent.Subscribe) {
    if sub == nil {
        snag.Panic("未找到有效骑士卡")
    }

    switch typ {
    case model.RiderPermissionTypeAssistance:
        if sub.Status != model.SubscribeStatusUsing {
            snag.Panic("无法发起救援")
        }
        return
    case model.RiderPermissionTypeBusiness:
        return
    case model.RiderPermissionTypeExchange:
        if sub.Status != model.SubscribeStatusUsing {
            snag.Panic("骑士卡状态异常")
        }
        return
    }
}
