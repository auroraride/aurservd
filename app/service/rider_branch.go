// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
)

type riderBranchService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewRiderBranch(r *ent.Rider) *riderBranchService {
    s := &riderBranchService{
        ctx:   context.Background(),
        rider: r,
    }
    s.ctx = context.WithValue(s.ctx, "rider", r)
    return s
}

func (s *riderBranchService) cabinetDetail(req *model.RiderBranchDetailReq) any {
    items, _ := ent.Database.Cabinet.QueryNotDeleted().
        Where(
            cabinet.BranchID(req.BranchID),
            cabinet.HasModelsWith(batterymodel.Model(req.Model)),
        ).
        All(s.ctx)
    return items
}

func (s *riderBranchService) storeDetail(req *model.RiderBranchDetailReq) any {
    return nil
}

func (s *riderBranchService) Detail(req *model.RiderBranchDetailReq) any {
    if req.Type == 1 {
        return s.storeDetail(req)
    } else {
        return s.cabinetDetail(req)
    }
}
