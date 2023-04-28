// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/workwx"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/cabinetfault"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
)

type cabinetFaultService struct {
	ctx      context.Context
	orm      *ent.CabinetFaultClient
	modifier *model.Modifier
}

func NewCabinetFault() *cabinetFaultService {
	return &cabinetFaultService{
		ctx: context.Background(),
		orm: ent.Database.CabinetFault,
	}
}

func NewCabinetFaultWithModifier(m *model.Modifier) *cabinetFaultService {
	s := NewCabinetFault()
	s.ctx = context.WithValue(s.ctx, "modifier", m)
	s.modifier = m
	return s
}

// Query 查找故障
func (s *cabinetFaultService) Query(id uint64) *ent.CabinetFault {
	cf, err := s.orm.QueryNotDeleted().Where(cabinetfault.ID(id)).First(s.ctx)
	if err != nil || cf == nil {
		snag.Panic("未找到故障")
	}
	return cf
}

// Report 骑手故障上报
func (s *cabinetFaultService) Report(rider *ent.Rider, req *model.CabinetFaultReportReq) bool {
	// 获取电柜信息
	cab, _ := ent.Database.Cabinet.QueryNotDeleted().
		Where(
			cabinet.ID(req.CabinetID),
			cabinet.Status(model.CabinetStatusNormal.Value()),
		).
		WithCity().
		WithBranch().
		First(s.ctx)
	if cab == nil {
		snag.Panic("未找到运营中的电柜")
	}
	attachments := make([]string, 0)
	if len(req.Attachments) > 0 {
		attachments = req.Attachments
	}
	city := cab.Edges.City
	s.orm.Create().
		// SetBrand(ca.Brand).
		// SetCity(model.City{
		//     ID:   city.ID,
		//     Name: city.Name,
		// }).
		SetCabinetID(cab.ID).
		SetBranchID(*cab.BranchID).
		SetRiderID(rider.ID).
		SetCityID(city.ID).
		// SetCabinetName(ca.Name).
		// SetSerial(ca.Serial).
		// SetModels(ca.Models).
		SetDescription(req.Description).
		SetAttachments(attachments).
		SetFault(req.Fault).
		SaveX(s.ctx)
	go workwx.New().SendCabinetFault(model.CabinetFaultNotice{
		City:        city.Name,
		Branch:      cab.Edges.Branch.Name,
		Name:        cab.Name,
		Serial:      cab.Serial,
		Phone:       rider.Phone,
		Fault:       req.Fault,
		Description: req.Description,
	})
	return true
}

// List 分页列举故障列表
func (s *cabinetFaultService) List(req *model.CabinetFaultListReq) (res *model.PaginationRes) {
	q := s.orm.QueryNotDeleted().
		WithBranch().
		WithRider().
		WithCity().
		WithCabinet().
		Order(ent.Asc(cabinetfault.FieldStatus), ent.Desc(cabinetfault.FieldCreatedAt))
	if req.CityID != nil {
		q.Where(cabinetfault.CityID(*req.CityID))
	}
	if req.CabinetName != nil {
		q.Where(cabinetfault.HasCabinetWith(cabinet.NameContainsFold(*req.CabinetName)))
	}
	if req.Serial != nil {
		q.Where(cabinetfault.HasCabinetWith(cabinet.SerialContainsFold(*req.Serial)))
	}
	if req.Status != nil {
		q.Where(cabinetfault.Status(*req.Status))
	}
	if req.Start != nil {
		start := carbon.ParseByLayout(*req.Start, carbon.DateLayout)
		if start.Error != nil {
			snag.Panic("日期格式错误")
		}
		q.Where(cabinetfault.CreatedAtGTE(start.Carbon2Time()))
	}
	if req.End != nil {
		end := carbon.ParseByLayout(*req.End, carbon.DateLayout)
		if end.Error != nil {
			snag.Panic("日期格式错误")
		}
		end.AddDay()
		q.Where(cabinetfault.CreatedAtLT(end.Carbon2Time()))
	}
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
func (s *cabinetFaultService) Deal(req *model.CabinetFaultDealReq) {
	s.orm.UpdateOne(s.Query(*req.ID)).
		SetRemark(*req.Remark).
		SetStatus(*req.Status).
		SaveX(s.ctx)
}
