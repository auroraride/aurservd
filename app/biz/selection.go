// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-21, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/material"
)

type selectionBiz struct {
	ctx context.Context
}

func NewSelection() *selectionBiz {
	return &selectionBiz{
		ctx: context.Background(),
	}
}

// MaintainerList 维护人员筛选数据
func (b *selectionBiz) MaintainerList() (res []*definition.MaintainerDetail) {
	ms, _ := ent.Database.Maintainer.Query().Where(
		maintainer.Enable(true),
	).All(b.ctx)
	for _, m := range ms {
		res = append(res, &definition.MaintainerDetail{
			ID:   m.ID,
			Name: m.Name,
		})
	}
	return
}

// StationList 站点筛选数据
func (b *selectionBiz) StationList() (res []model.CascaderOptionLevel2) {
	stList, _ := ent.Database.EnterpriseStation.QueryNotDeleted().WithCity().All(b.ctx)
	cityIds := make([]uint64, 0)
	cityIdMap := make(map[uint64]*ent.City)
	cityIdListMap := make(map[uint64][]model.SelectOption)
	for _, st := range stList {
		if st.Edges.City != nil {
			cityID := st.Edges.City.ID
			cityIds = append(cityIds, cityID)
			cityIdMap[cityID] = st.Edges.City
			cityIdListMap[cityID] = append(cityIdListMap[cityID], model.SelectOption{
				Label: st.Name,
				Value: st.ID,
			})
		}
	}

	for _, cityId := range cityIds {
		if cityIdMap[cityId] != nil && len(cityIdListMap[cityId]) != 0 {
			res = append(res, model.CascaderOptionLevel2{
				SelectOption: model.SelectOption{
					Value: cityIdMap[cityId].ID,
					Label: cityIdMap[cityId].Name,
				},
				Children: cityIdListMap[cityId],
			})
		}
	}

	return
}

// // StationList 站点筛选数据
// func (b *selectionBiz) StationList() (res []model.EnterpriseStationListRes) {
// 	res = make([]model.EnterpriseStationListRes, 0)
// 	items, _ := s.orm.QueryNotDeleted().Where(enterprisestation.EnterpriseID(req.EnterpriseID)).WithCity().All(s.ctx)
// 	for _, item := range items {
// 		res = append(res, model.EnterpriseStationListRes{
// 			EnterpriseStation: model.EnterpriseStation{
// 				ID:   item.ID,
// 				Name: item.Name,
// 			},
// 			City: model.City{
// 				ID:   item.Edges.City.ID,
// 				Name: item.Edges.City.Name,
// 			},
// 		})
// 	}
// 	return
// }

// MaterialSelect 资产分类筛选数据
func (b *selectionBiz) MaterialSelect(req *model.SelectMaterialReq) (res []model.SelectOption) {
	q := ent.Database.Material.QueryNotDeleted()
	if req.Type != nil {
		q.Where(material.Type(req.Type.Value()))
	}
	items, _ := q.All(b.ctx)
	for _, item := range items {
		res = append(res, model.SelectOption{
			Label: item.Name,
			Value: item.ID,
		})
	}
	return
}
