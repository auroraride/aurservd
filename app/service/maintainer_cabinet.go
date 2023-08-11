// Copyright (C) liasica. 2023-present.
//
// Created at 2023-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
)

type maintainerCabinetService struct {
	*BaseService
	orm *ent.CabinetClient
}

func NewMaintainerCabinet(params ...any) *maintainerCabinetService {
	return &maintainerCabinetService{
		BaseService: newService(params...),
		orm:         ent.Database.Cabinet,
	}
}

// List 运维归属电柜列表
func (s *maintainerCabinetService) List(cityIDs []uint64, pg *model.PaginationReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Where(cabinet.CityIDIn(cityIDs...))
	return model.ParsePaginationResponse(q, *pg, func(cab *ent.Cabinet) model.CabinetListByDistanceRes {
		return model.CabinetListByDistanceRes{
			CabinetBasicInfo: model.CabinetBasicInfo{
				ID:     cab.ID,
				Brand:  cab.Brand,
				Serial: cab.Serial,
				Name:   cab.Name,
			},
			Status: cab.Status,
			Health: cab.Health,
		}
	}, func(cabinets ent.Cabinets) {
		NewCabinet().SyncCabinets(cabinets)
	})
}

// Detail 获取电柜详情
func (s *maintainerCabinetService) Detail(serial string) *model.CabinetDetailRes {
	return NewCabinet().DetailFromSerial(serial)
}
