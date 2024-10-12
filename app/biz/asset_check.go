// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-21, by aurb

package biz

import (
	"context"

	"github.com/pkg/errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/pkg/silk"
)

type assetCheckBiz struct {
	orm *ent.AssetCheckClient
	ctx context.Context
}

func NewAssetCheck() *assetCheckBiz {
	return &assetCheckBiz{
		orm: ent.Database.AssetCheck,
		ctx: context.Background(),
	}
}

// GetAssetBySN 通过sn查询资产
func (b *assetCheckBiz) GetAssetBySN(assetSignInfo definition.AssetSignInfo, req *definition.AssetCheckByAssetSnReq) (res *model.AssetCheckByAssetSnRes, err error) {
	newReq := model.AssetCheckByAssetSnReq{
		SN: req.SN,
	}

	if assetSignInfo.AssetManager != nil {
		wType := model.OperatorTypeAssetManager
		newReq.OperatorType = wType
		newReq.OperatorID = assetSignInfo.AssetManager.ID
		// 上班确定仓库位置
		am, _ := ent.Database.AssetManager.QueryNotDeleted().WithDutyWarehouse().
			Where(
				assetmanager.ID(assetSignInfo.AssetManager.ID),
			).First(b.ctx)
		if am == nil {
			return nil, errors.New("无效仓管员")
		}

		if am.Edges.DutyWarehouse == nil {
			return nil, errors.New("未找到仓管员上班信息")
		}
		newReq.LocationsType = model.AssetLocationsTypeWarehouse
		newReq.LocationsID = am.Edges.DutyWarehouse.ID
	}
	if assetSignInfo.Employee != nil {
		sType := model.OperatorTypeEmployee
		newReq.OperatorType = sType
		newReq.OperatorID = assetSignInfo.Employee.ID
		// 上班确定门店位置
		ep, _ := ent.Database.Employee.QueryNotDeleted().WithDutyStore().
			Where(
				employee.ID(assetSignInfo.Employee.ID),
			).First(b.ctx)
		if ep == nil {
			return nil, errors.New("无效店员")
		}

		if ep.Edges.DutyStore == nil {
			return nil, errors.New("未找到店员上班信息")
		}

		newReq.LocationsType = model.AssetLocationsTypeStore
		newReq.LocationsID = ep.Edges.DutyStore.ID
	}
	if assetSignInfo.Agent != nil {
		if req.StationID == nil {
			return nil, errors.New("未选择具体站点")
		}
		sType := model.OperatorTypeAgent
		newReq.OperatorType = sType
		newReq.OperatorID = assetSignInfo.Agent.ID
		newReq.LocationsType = model.AssetLocationsTypeStation
		newReq.LocationsID = *req.StationID
	}

	if assetSignInfo.Maintainer != nil {
		sType := model.OperatorTypeAgent
		newReq.OperatorType = sType
		newReq.OperatorID = assetSignInfo.Maintainer.ID
		newReq.LocationsType = model.AssetLocationsTypeOperation
		newReq.LocationsID = assetSignInfo.Maintainer.ID
	}

	return service.NewAssetCheck().GetAssetBySN(b.ctx, &newReq)
}

// Create 创建资产盘点
func (b *assetCheckBiz) Create(assetSignInfo definition.AssetSignInfo, req *definition.AssetCheckCreateReq) (res *definition.AssetCheckCreateRes, err error) {
	var md model.Modifier

	newReq := model.AssetCheckCreateReq{
		AssetCheckCreateDetail: req.AssetCheckCreateDetail,
		StartAt:                req.StartAt,
		EndAt:                  req.EndAt,
	}

	if assetSignInfo.AssetManager != nil {
		wType := model.AssetLocationsTypeWarehouse
		newReq.LocationsType = wType
		am, _ := ent.Database.AssetManager.QueryNotDeleted().WithDutyWarehouse().
			Where(
				assetmanager.ID(assetSignInfo.AssetManager.ID),
			).First(b.ctx)
		if am == nil {
			return nil, errors.New("无效仓管员")
		}

		if am.Edges.DutyWarehouse == nil {
			return nil, errors.New("未找到仓管员上班信息")
		}
		newReq.LocationsID = am.Edges.DutyWarehouse.ID
		newReq.OperatorID = assetSignInfo.AssetManager.ID
		newReq.OperatorType = model.OperatorTypeAssetManager

		md = model.Modifier{
			ID:    assetSignInfo.AssetManager.ID,
			Name:  assetSignInfo.AssetManager.Name,
			Phone: assetSignInfo.AssetManager.Phone,
		}
	}
	if assetSignInfo.Employee != nil {
		sType := model.AssetLocationsTypeStore
		newReq.LocationsType = sType
		ep, _ := ent.Database.Employee.QueryNotDeleted().WithDutyStore().
			Where(
				employee.ID(assetSignInfo.Employee.ID),
			).First(b.ctx)
		if ep == nil {
			return nil, errors.New("无效店员")
		}

		if ep.Edges.DutyStore == nil {
			return nil, errors.New("未找到店员上班信息")
		}

		newReq.LocationsID = ep.Edges.DutyStore.ID
		newReq.OperatorID = assetSignInfo.Employee.ID
		newReq.OperatorType = model.OperatorTypeEmployee
		md = model.Modifier{
			ID:    assetSignInfo.Employee.ID,
			Name:  assetSignInfo.Employee.Name,
			Phone: assetSignInfo.Employee.Phone,
		}
	}
	if assetSignInfo.Agent != nil {
		if req.StationID == nil {
			return nil, errors.New("站点ID数据有误")
		}

		sType := model.AssetLocationsTypeStation
		newReq.LocationsType = sType
		newReq.LocationsID = *req.StationID
		newReq.OperatorID = assetSignInfo.Agent.ID
		newReq.OperatorType = model.OperatorTypeAgent
		md = model.Modifier{
			ID:    assetSignInfo.Agent.ID,
			Name:  assetSignInfo.Agent.Name,
			Phone: assetSignInfo.Agent.Phone,
		}
	}

	cId, err := service.NewAssetCheck().CreateAssetCheck(b.ctx, &newReq, &md)
	if err != nil {
		return nil, err
	}

	return &definition.AssetCheckCreateRes{ID: cId}, nil
}

// List 盘点记录
func (b *assetCheckBiz) List(assetSignInfo definition.AssetSignInfo, req *definition.AssetCheckListReq) (res *model.PaginationRes, err error) {
	newReq := model.AssetCheckListReq{
		PaginationReq: req.PaginationReq,
		AssetCheckListFilter: model.AssetCheckListFilter{
			Keyword:     req.Keyword,
			StartAt:     req.StartAt,
			EndAt:       req.EndAt,
			CheckResult: req.CheckResult,
		},
	}

	if assetSignInfo.AssetManager != nil {
		wType := model.AssetLocationsTypeWarehouse
		newReq.LocationsType = &wType

		am, _ := ent.Database.AssetManager.QueryNotDeleted().WithDutyWarehouse().
			Where(
				assetmanager.ID(assetSignInfo.AssetManager.ID),
			).First(b.ctx)
		if am == nil {
			return nil, errors.New("无效仓管员")
		}

		if am.Edges.DutyWarehouse == nil {
			return nil, errors.New("未找到仓管员上班信息")
		}

		newReq.LocationsID = silk.UInt64(am.Edges.DutyWarehouse.ID)
	}
	if assetSignInfo.Employee != nil {
		sType := model.AssetLocationsTypeStore
		newReq.LocationsType = &sType

		ep, _ := ent.Database.Employee.QueryNotDeleted().WithDutyStore().
			Where(
				employee.ID(assetSignInfo.Employee.ID),
			).First(b.ctx)
		if ep == nil {
			return nil, errors.New("无效店员")
		}

		if ep.Edges.DutyStore == nil {
			return nil, errors.New("未找到店员上班信息")
		}

		newReq.LocationsID = silk.UInt64(ep.Edges.DutyStore.ID)
	}
	if assetSignInfo.Agent != nil {
		sType := model.AssetLocationsTypeStation
		newReq.LocationsType = &sType

		list, _ := ent.Database.EnterpriseStation.QueryNotDeleted().
			Where(
				enterprisestation.HasAgentsWith(agent.ID(assetSignInfo.Agent.ID)),
			).All(b.ctx)
		ids := make([]uint64, 0)
		for _, v := range list {
			ids = append(ids, v.ID)
		}

		newReq.LocationsIds = ids
	}

	return service.NewAssetCheck().List(b.ctx, &newReq)
}
