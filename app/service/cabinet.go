// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
)

type cabinetService struct {
    ctx context.Context
    orm *ent.CabinetClient
}

func NewCabinet() *cabinetService {
    return &cabinetService{
        ctx: context.Background(),
        orm: ar.Ent.Cabinet,
    }
}

// CreateCabinet 创建电柜
func (s *cabinetService) CreateCabinet(modifier *model.Modifier, req *model.CabinetCreateReq) (res *model.CabinetItem) {
    q := s.orm.Create().
        SetName(req.Name).
        SetSerial(req.Serial).
        SetSn(shortuuid.New()).
        SetStatus(uint(req.Status)).
        SetDoors(req.Doors).
        SetLastModifier(modifier).
        SetCreator(modifier).
        SetNillableRemark(req.Remark).
        SetBrand(req.Brand.String())
    if req.BranchID != nil {
        q.SetBranchID(*req.BranchID)
    }

    // 查询设置电池型号
    bms := make([]model.BatteryModel, len(req.Models))
    for i, id := range req.Models {
        bm := ar.Ent.BatteryModel.Query().Where(batterymodel.ID(id)).OnlyX(s.ctx)
        if bm == nil {
            snag.Panic("未找到电池型号")
        }
        q.AddBms(bm)
        bms[i] = model.BatteryModel{
            ID:       id,
            Voltage:  bm.Voltage,
            Capacity: bm.Capacity,
        }
    }
    q.SetModels(bms)

    item := q.SaveX(s.ctx)
    res = new(model.CabinetItem)
    _ = copier.Copy(res, item)
    res.Models = bms
    return
}

// Query 查询电柜
func (s *cabinetService) Query(req *model.CabinetQueryReq) (res *model.PaginationRes) {
    res = new(model.PaginationRes)
    q := s.orm.Query().WithBranch(
        func(bq *ent.BranchQuery) {
            bq.WithCity()
        },
    )
    if req.Serial != nil {
        q.Where(cabinet.SerialContainsFold(*req.Serial))
    }
    if req.Name != nil {
        q.Where(cabinet.NameContainsFold(*req.Name))
    }
    if req.CityId != nil {
        q.Where(cabinet.HasBranchWith(branch.CityID(*req.CityId)))
    }
    if req.Brand != nil {
        q.Where(cabinet.Brand(*req.Brand))
    }
    if req.Status != nil {
        q.Where(cabinet.Status(uint(*req.Status)))
    }

    res.Pagination = q.PaginationResult(req.PaginationReq)

    items := q.Pagination(req.PaginationReq).AllX(s.ctx)
    out := make([]model.CabinetItem, len(items))
    for i, item := range items {
        city := item.Edges.Branch.Edges.City
        _ = copier.Copy(&out[i], item)
        out[i].City = model.City{
            ID:   city.ID,
            Name: city.Name,
        }
    }
    res.Items = out
    return
}
