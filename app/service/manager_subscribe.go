// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/pkg/silk"
)

type managerSubscribeService struct {
    *BaseService
}

func NewManagerSubscribe(params ...any) *managerSubscribeService {
    return &managerSubscribeService{
        BaseService: newService(params...),
    }
}

func (s *managerSubscribeService) Active(req *model.ManagerSubscribeActive) {
    var bikeID *uint64
    if req.EbikeKeyword != nil {
        bike := NewAllocate().UnallocatedEbikeInfo(*req.EbikeKeyword)
        bikeID = silk.UInt64(bike.ID)
    }

    NewAllocate(s.modifier).Create(&model.AllocateCreateReq{
        EbikeID:     bikeID,
        SubscribeID: req.ID,
        StoreID:     req.StoreID,
    })
}
