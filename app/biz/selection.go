// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-21, by aurb

package biz

import (
	"context"

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
func (b *selectionBiz) MaintainerList() (res []model.SelectOption) {
	res = make([]model.SelectOption, 0)
	ms, _ := ent.Database.Maintainer.Query().Where(
		maintainer.Enable(true),
	).All(b.ctx)
	for _, m := range ms {
		res = append(res, model.SelectOption{
			Value: m.ID,
			Label: m.Name,
		})
	}
	return
}

// StationList 团签站点筛选数据
func (b *selectionBiz) StationList() (res []model.CascaderOptionLevel2) {
	res = make([]model.CascaderOptionLevel2, 0)
	stList, _ := ent.Database.EnterpriseStation.QueryNotDeleted().WithEnterprise().All(b.ctx)
	eIds := make([]uint64, 0)
	eIdMap := make(map[uint64]*ent.Enterprise)
	eIdListMap := make(map[uint64][]model.SelectOption)
	for _, st := range stList {
		if st.Edges.Enterprise != nil {
			eId := st.Edges.Enterprise.ID
			if eIdMap[eId] == nil {
				eIds = append(eIds, eId)
				eIdMap[eId] = st.Edges.Enterprise
			}
			eIdListMap[eId] = append(eIdListMap[eId], model.SelectOption{
				Label: st.Name,
				Value: st.ID,
			})
		}
	}

	for _, eId := range eIds {
		if eIdMap[eId] != nil && len(eIdListMap[eId]) != 0 {
			res = append(res, model.CascaderOptionLevel2{
				SelectOption: model.SelectOption{
					Value: eIdMap[eId].ID,
					Label: eIdMap[eId].Name,
				},
				Children: eIdListMap[eId],
			})
		}
	}

	return
}

// CityStationList 城市站点筛选数据
func (b *selectionBiz) CityStationList() (res []model.CascaderOptionLevel2) {
	stList, _ := ent.Database.EnterpriseStation.QueryNotDeleted().WithCity().All(b.ctx)
	cIds := make([]uint64, 0)
	cIdMap := make(map[uint64]*ent.City)
	cIdListMap := make(map[uint64][]model.SelectOption)
	for _, st := range stList {
		if st.Edges.City != nil {
			cId := st.Edges.City.ID
			if cIdMap[cId] == nil {
				cIds = append(cIds, cId)
				cIdMap[cId] = st.Edges.City
			}
			cIdListMap[cId] = append(cIdListMap[cId], model.SelectOption{
				Label: st.Name,
				Value: st.ID,
			})
		}
	}

	for _, cId := range cIds {
		if cIdMap[cId] != nil && len(cIdListMap[cId]) != 0 {
			res = append(res, model.CascaderOptionLevel2{
				SelectOption: model.SelectOption{
					Value: cIdMap[cId].ID,
					Label: cIdMap[cId].Name,
				},
				Children: cIdListMap[cId],
			})
		}
	}

	return
}

// MaterialSelect 资产分类筛选数据
func (b *selectionBiz) MaterialSelect(req *model.SelectMaterialReq) (res []model.SelectOption) {
	res = make([]model.SelectOption, 0)
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
