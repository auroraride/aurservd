// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/pkg/snag"
)

type enterpriseStationService struct {
	ctx          context.Context
	modifier     *model.Modifier
	rider        *ent.Rider
	employee     *ent.Employee
	orm          *ent.EnterpriseStationClient
	employeeInfo *model.Employee
}

func NewEnterpriseStation() *enterpriseStationService {
	return &enterpriseStationService{
		ctx: context.Background(),
		orm: ent.Database.EnterpriseStation,
	}
}

func NewEnterpriseStationWithRider(r *ent.Rider) *enterpriseStationService {
	s := NewEnterpriseStation()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, r)
	s.rider = r
	return s
}

func NewEnterpriseStationWithModifier(m *model.Modifier) *enterpriseStationService {
	s := NewEnterpriseStation()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func NewEnterpriseStationWithEmployee(e *ent.Employee) *enterpriseStationService {
	s := NewEnterpriseStation()
	if e != nil {
		s.employee = e
		s.employeeInfo = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
		s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
	}
	return s
}

func (s *enterpriseStationService) QueryX(id uint64) *ent.EnterpriseStation {
	item, _ := s.orm.QueryNotDeleted().Where(enterprisestation.ID(id)).First(s.ctx)
	if item == nil {
		snag.Panic("未找到站点")
	}
	return item
}

// Create 创建站点
func (s *enterpriseStationService) Create(req *model.EnterpriseStationCreateReq) uint64 {
	es := s.orm.Create().SetEnterpriseID(req.EnterpriseID).SetName(req.Name).SetCityID(req.CityID).SaveX(s.ctx)
	return es.ID
}

// Modify 修改站点
func (s *enterpriseStationService) Modify(req *model.EnterpriseStationModifyReq) {
	s.QueryX(req.ID).Update().SetName(req.Name).SetCityID(req.CityID).SaveX(s.ctx)
}

func (s *enterpriseStationService) List(req *model.EnterpriseStationListReq) (res []model.EnterpriseStationListRes) {
	res = make([]model.EnterpriseStationListRes, 0)
	items, _ := s.orm.QueryNotDeleted().Where(enterprisestation.EnterpriseID(req.EnterpriseID)).WithCity().All(s.ctx)
	for _, item := range items {
		res = append(res, model.EnterpriseStationListRes{
			EnterpriseStation: model.EnterpriseStation{
				ID:   item.ID,
				Name: item.Name,
			},
			City: model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			},
		})
	}
	return
}
