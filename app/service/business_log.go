// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-16
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/business"
	"github.com/auroraride/aurservd/pkg/snag"
)

type businessLogService struct {
	ctx          context.Context
	modifier     *model.Modifier
	employee     *ent.Employee
	orm          *ent.BusinessClient
	creator      *ent.BusinessCreate
	employeeInfo *model.Employee
}

func NewBusinessLog(sub *ent.Subscribe) *businessLogService {
	s := &businessLogService{
		ctx:     context.Background(),
		orm:     ent.Database.Business,
		creator: ent.Database.Business.Create().SetNillableEnterpriseID(sub.EnterpriseID).SetNillableStationID(sub.StationID),
	}
	s.setSubscribe(sub)
	return s
}

// SetModifier 设置管理员
func (s *businessLogService) SetModifier(m *model.Modifier) *businessLogService {
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
		s.creator.SetLastModifier(m)
	}
	return s
}

// SetEmployee 设置店员和门店
func (s *businessLogService) SetEmployee(e *ent.Employee) *businessLogService {
	if e != nil {
		s.employeeInfo = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
		s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
		s.employee = e
		s.creator.SetEmployee(e)
		if e.Edges.Store != nil {
			s.creator.SetStoreID(e.Edges.Store.ID)
		}
	}
	return s
}

func (s *businessLogService) SetStore(sto *ent.Store) *businessLogService {
	if sto != nil {
		s.creator.SetStore(sto)
	}
	return s
}

// SetCabinet 设置电柜
func (s *businessLogService) SetCabinet(cab *ent.Cabinet) *businessLogService {
	if cab != nil {
		s.creator.SetCabinet(cab)
	}
	return s
}

func (s *businessLogService) SetBinInfo(bin *model.BinInfo) *businessLogService {
	if bin != nil {
		s.creator.SetBinInfo(bin)
	}
	return s
}

func (s *businessLogService) SetStock(sk *ent.Stock) *businessLogService {
	if sk != nil {
		s.creator.SetStockSn(sk.Sn)
	}
	return s
}

func (s *businessLogService) setSubscribe(sub *ent.Subscribe) {
	s.creator.SetRiderID(sub.RiderID).
		SetSubscribeID(sub.ID).
		SetCityID(sub.CityID).
		SetNillableEnterpriseID(sub.EnterpriseID).
		SetNillableStationID(sub.StationID).
		SetNillablePlanID(sub.PlanID)
}

func (s *businessLogService) SetBattery(bat *model.Battery) *businessLogService {
	if bat != nil {
		s.creator.SetBatteryID(bat.ID)
	}
	return s
}

func (s *businessLogService) Save(typ business.Type) (*ent.Business, error) {
	return s.creator.SetType(typ).Save(s.ctx)
}

func (s *businessLogService) SaveX(typ business.Type) *ent.Business {
	biz, err := s.Save(typ)
	if err != nil {
		snag.Panic("日志保存失败")
	}
	return biz
}

// SetAgentId 设置代理人
func (s *businessLogService) SetAgentId(agentId *uint64) *businessLogService {
	s.creator.SetNillableAgentID(agentId)
	return s
}
