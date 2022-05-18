// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/pkg/snag"
)

type cabinetFaultService struct {
    ctx context.Context
    orm *ent.CabinetFaultClient
}

func NewCabinetFault() *cabinetFaultService {
    return &cabinetFaultService{
        ctx: context.Background(),
        orm: ar.Ent.CabinetFault,
    }
}

// Report 骑手故障上报
func (s *cabinetFaultService) Report(rider *ent.Rider, req *model.CabinetFaultReportReq) bool {
    // 获取电柜信息
    ca, _ := ar.Ent.Cabinet.Query().
        Where(cabinet.ID(req.CabinetID)).
        WithBranch(func(bq *ent.BranchQuery) {
            bq.WithCity()
        }).
        Only(s.ctx)
    if ca == nil {
        snag.Panic("未找到电柜")
    }
    if ca.Edges.Branch == nil || ca.Edges.Branch.Edges.City == nil {
        snag.Panic("电柜未投放, 无法上报")
    }
    attachments := make([]string, 0)
    if len(req.Attachments) > 0 {
        attachments = req.Attachments
    }
    branch := ca.Edges.Branch
    city := branch.Edges.City
    s.orm.Create().
        SetBrand(ca.Brand).
        SetCity(model.City{
            ID:   city.ID,
            Name: city.Name,
        }).
        SetCabinetID(ca.ID).
        SetBranchID(branch.ID).
        SetRiderID(rider.ID).
        SetCabinetName(ca.Name).
        SetSerial(ca.Serial).
        SetModels(ca.Models).
        SetDescription(req.Description).
        SetAttachments(attachments).
        SetFault(req.Fault).
        SaveX(s.ctx)
    return true
}
