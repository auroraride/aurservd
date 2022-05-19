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
    "github.com/auroraride/aurservd/internal/ent/cabinetfault"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "time"
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

// Query 查找故障
func (s *cabinetFaultService) Query(id uint64) *ent.CabinetFault {
    cf, err := s.orm.Query().Where(cabinetfault.ID(id)).Only(s.ctx)
    if err != nil || cf == nil {
        snag.Panic("未找到故障")
    }
    return cf
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
        // SetBrand(ca.Brand).
        // SetCity(model.City{
        //     ID:   city.ID,
        //     Name: city.Name,
        // }).
        SetCabinetID(ca.ID).
        SetBranchID(branch.ID).
        SetRiderID(rider.ID).
        SetCityID(city.ID).
        // SetCabinetName(ca.Name).
        // SetSerial(ca.Serial).
        // SetModels(ca.Models).
        SetDescription(req.Description).
        SetAttachments(attachments).
        SetFault(req.Fault).
        SaveX(s.ctx)
    return true
}

// List 分页列举故障列表
func (s *cabinetFaultService) List(req *model.CabinetFaultListReq) (res *model.PaginationRes) {
    cq := ar.Ent.Cabinet.Query()
    q := s.orm.Query().
        WithBranch().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithCity()
    if req.CityID != nil {
        q.Where(cabinetfault.CityID(*req.CityID))
    }
    if req.CabinetName != nil {
        cq.Where(cabinet.NameContainsFold(*req.CabinetName))
    }
    if req.Serial != nil {
        cq.Where(cabinet.SerialContainsFold(*req.Serial))
    }
    if req.Status != nil {
        q.Where(cabinetfault.Status(*req.Status))
    }
    if req.Start != nil {
        start, err := time.Parse(carbon.DateLayout, *req.Start)
        if err != nil {
            snag.Panic("日期格式错误")
        }
        q.Where(cabinetfault.CreatedAtGTE(start))
    }
    if req.End != nil {
        end, err := time.Parse(carbon.DateLayout, *req.End)
        if err != nil {
            snag.Panic("日期格式错误")
        }
        end.AddDate(0, 0, 1)
        q.Where(cabinetfault.CreatedAtLT(end))
    }
    q.WithCabinet(func(query *ent.CabinetQuery) {
        query = cq
    })
    res = &model.PaginationRes{Pagination: q.PaginationResult(req.PaginationReq)}
    items := q.Pagination(req.PaginationReq).AllX(s.ctx)
    out := make([]model.CabinetFaultItem, len(items))
    for i, item := range items {
        _ = copier.Copy(&out[i], item)
        _ = copier.Copy(&out[i].City, item.Edges.City)
        _ = copier.Copy(&out[i].Cabinet, item.Edges.Cabinet)
        out[i].Rider = NewRider().GetRiderSampleInfo(item.Edges.Rider)
    }
    res.Items = out
    return
}

// Deal 处理故障
func (s *cabinetFaultService) Deal(m *model.Modifier, req *model.CabinetFaultDealReq) {
    s.orm.UpdateOne(s.Query(*req.ID)).
        SetRemark(*req.Remark).
        SetStatus(*req.Status).
        SetLastModifier(m).
        SaveX(s.ctx)
}
