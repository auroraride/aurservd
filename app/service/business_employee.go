// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "strconv"
    "strings"
)

type businessEmployeeService struct {
    *BaseService
}

func NewBusinessEmployee(params ...any) *businessEmployeeService {
    return &businessEmployeeService{
        BaseService: newService(params...),
    }
}

// Inactive 获取骑手待激活订阅详情
func (s *businessEmployeeService) Inactive(qr string) (*model.SubscribeActiveInfo, *ent.Subscribe) {
    if strings.HasPrefix(qr, "SUBSCRIBE:") {
        qr = strings.ReplaceAll(qr, "SUBSCRIBE:", "")
    }
    id, _ := strconv.ParseUint(strings.TrimSpace(qr), 10, 64)
    return NewBusinessRiderWithEmployee(s.entEmployee).Inactive(id)
}

// Active 激活订阅
func (s *businessEmployeeService) Active(req *model.AllocateCreateReq) (res model.IDPostReq) {
    if s.entStore == nil {
        snag.Panic("当前未上班")
    }

    req.EmployeeID = silk.UInt64(s.entEmployee.ID)
    req.StoreID = silk.UInt64(s.entStore.ID)

    return NewAllocate().Create(req)
}
