// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-20, by aurb

package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
	"github.com/auroraride/aurservd/pkg/silk"
)

type assetTransferBiz struct {
	orm *ent.AssetTransferClient
	ctx context.Context
}

func NewAssetTransfer() *assetTransferBiz {
	return &assetTransferBiz{
		orm: ent.Database.AssetTransfer,
		ctx: context.Background(),
	}
}

// TransferList 调拨记录列表
func (b *assetTransferBiz) TransferList(assetSignInfo definition.AssetSignInfo, req *definition.TransferListReq) (res *model.PaginationRes, err error) {
	newReq := model.AssetTransferListReq{
		PaginationReq:       req.PaginationReq,
		AssetTransferFilter: req.AssetTransferFilter,
	}

	if assetSignInfo.AssetManager != nil {
		newReq.AssetManagerID = assetSignInfo.AssetManager.ID
	}

	if assetSignInfo.Employee != nil {
		newReq.EmployeeID = assetSignInfo.Employee.ID
	}

	if assetSignInfo.Agent != nil {
		newReq.AgentID = assetSignInfo.Agent.ID
	}

	if assetSignInfo.Maintainer != nil {
		newReq.MaintainerID = assetSignInfo.Maintainer.ID
	}

	return service.NewAssetTransfer().TransferList(context.Background(), &newReq)
}

// TransferDetail 调拨记录详情
func (b *assetTransferBiz) TransferDetail(assetSignInfo definition.AssetSignInfo, req *model.AssetTransferDetailReq) (res *definition.TransferDetailRes, err error) {
	var t *ent.AssetTransfer
	t, err = ent.Database.AssetTransfer.QueryNotDeleted().WithTransferDetails().
		WithOutOperateAgent().WithOutOperateManager().WithOutOperateEmployee().WithOutOperateMaintainer().WithOutOperateAssetManager().
		WithFromLocationWarehouse().WithFromLocationStore().WithFromLocationStation().WithFromLocationOperator().
		WithToLocationWarehouse().WithToLocationStore().WithToLocationStation().WithToLocationOperator().
		Where(assettransfer.ID(req.ID)).First(b.ctx)
	if err != nil {
		return nil, err
	}

	details, err := service.NewAssetTransfer().TransferDetail(b.ctx, req)
	if err != nil {
		return nil, err
	}

	var atu model.AssetTransferUserId
	switch {
	case assetSignInfo.AssetManager != nil:
		atu.AssetManagerID = assetSignInfo.AssetManager.ID
	case assetSignInfo.Employee != nil:
		atu.EmployeeID = assetSignInfo.Employee.ID
	case assetSignInfo.Agent != nil:
		atu.AgentID = assetSignInfo.Agent.ID
	case assetSignInfo.Maintainer != nil:
		atu.MaintainerID = assetSignInfo.Maintainer.ID
	}
	atr := service.NewAssetTransfer().TransferInfo(&atu, t)
	res = &definition.TransferDetailRes{
		AssetTransferListRes: *atr,
		Detail:               details,
	}
	return
}

// TransferReceive 接收资产
func (b *assetTransferBiz) TransferReceive(assetSignInfo definition.AssetSignInfo, req *definition.AssetTransferReceiveBatchReq) (err error) {
	var md model.Modifier

	newReq := model.AssetTransferReceiveBatchReq{
		AssetTransferReceive: req.AssetTransferReceive,
	}

	if assetSignInfo.AssetManager != nil {
		newReq.OperateType = model.OperatorTypeAssetManager
		md = model.Modifier{
			ID:    assetSignInfo.AssetManager.ID,
			Name:  assetSignInfo.AssetManager.Name,
			Phone: assetSignInfo.AssetManager.Phone,
		}
	}
	if assetSignInfo.Employee != nil {
		newReq.OperateType = model.OperatorTypeEmployee
		md = model.Modifier{
			ID:    assetSignInfo.Employee.ID,
			Name:  assetSignInfo.Employee.Name,
			Phone: assetSignInfo.Employee.Phone,
		}
	}
	if assetSignInfo.Agent != nil {
		newReq.OperateType = model.OperatorTypeAgent
		md = model.Modifier{
			ID:    assetSignInfo.Agent.ID,
			Name:  assetSignInfo.Agent.Name,
			Phone: assetSignInfo.Agent.Phone,
		}
	}
	if assetSignInfo.Maintainer != nil {
		newReq.OperateType = model.OperatorTypeMaintainer
		md = model.Modifier{
			ID:    assetSignInfo.Maintainer.ID,
			Name:  assetSignInfo.Maintainer.Name,
			Phone: assetSignInfo.Maintainer.Phone,
		}
	}

	return service.NewAssetTransfer().TransferReceive(b.ctx, &newReq, &md)
}

// Transfer 创建资产调拨
func (b *assetTransferBiz) Transfer(assetSignInfo definition.AssetSignInfo, req *definition.AssetTransferCreateReq) (err error) {

	var md model.Modifier

	newReq := model.AssetTransferCreateReq{
		ToLocationType:    req.ToLocationType,
		ToLocationID:      req.ToLocationID,
		Details:           req.Details,
		Reason:            req.Reason,
		AssetTransferType: model.AssetTransferTypeTransfer,
	}

	if assetSignInfo.AssetManager != nil {
		wType := model.AssetLocationsTypeWarehouse
		newReq.FromLocationType = &wType
		if req.FromLocationID == nil {
			// 当前上班ID
			loc, _ := ent.Database.AssetManager.QueryNotDeleted().Where(assetmanager.ID(assetSignInfo.AssetManager.ID)).First(b.ctx)
			if loc == nil {
				return errors.New("当前未上班")
			}
			req.FromLocationID = loc.WarehouseID
			newReq.FromLocationID = req.FromLocationID
		} else {
			newReq.FromLocationID = req.FromLocationID
		}

		md = model.Modifier{
			ID:    assetSignInfo.AssetManager.ID,
			Name:  assetSignInfo.AssetManager.Name,
			Phone: assetSignInfo.AssetManager.Phone,
		}

		newReq.OperatorType = model.OperatorTypeAssetManager
		newReq.OperatorID = assetSignInfo.AssetManager.ID
	}
	if assetSignInfo.Employee != nil {
		sType := model.AssetLocationsTypeStore
		newReq.FromLocationType = &sType
		if req.FromLocationID == nil {
			// 当前上班ID
			loc, _ := ent.Database.Employee.QueryNotDeleted().Where(employee.ID(assetSignInfo.Employee.ID)).First(b.ctx)
			if loc == nil {
				return errors.New("当前未上班")
			}
			req.FromLocationID = loc.DutyStoreID
		} else {
			newReq.FromLocationID = req.FromLocationID
		}
		md = model.Modifier{
			ID:    assetSignInfo.Employee.ID,
			Name:  assetSignInfo.Employee.Name,
			Phone: assetSignInfo.Employee.Phone,
		}
		newReq.OperatorType = model.OperatorTypeEmployee
		newReq.OperatorID = assetSignInfo.Employee.ID
	}
	if assetSignInfo.Agent != nil {
		sType := model.AssetLocationsTypeStation
		newReq.FromLocationType = &sType
		if req.FromLocationID == nil {
			return errors.New("未选择出库位置")
		} else {
			newReq.FromLocationID = req.FromLocationID
		}
		md = model.Modifier{
			ID:    assetSignInfo.Agent.ID,
			Name:  assetSignInfo.Agent.Name,
			Phone: assetSignInfo.Agent.Phone,
		}
		newReq.OperatorType = model.OperatorTypeAgent
		newReq.OperatorID = assetSignInfo.Agent.ID
	}
	if assetSignInfo.Maintainer != nil {
		for _, v := range req.Details {
			if v.AssetType != model.AssetTypeSmartBattery && v.AssetType != model.AssetTypeNonSmartBattery {
				return errors.New("运维端当前只支持电池调拨")
			}
		}

		sType := model.AssetLocationsTypeOperation
		newReq.FromLocationType = &sType
		if req.FromLocationID == nil {
			req.FromLocationID = silk.UInt64(assetSignInfo.Maintainer.ID)
			newReq.FromLocationID = req.FromLocationID
		} else {
			newReq.FromLocationID = req.FromLocationID
		}
		md = model.Modifier{
			ID:    assetSignInfo.Maintainer.ID,
			Name:  assetSignInfo.Maintainer.Name,
			Phone: assetSignInfo.Maintainer.Phone,
		}
		newReq.OperatorType = model.OperatorTypeMaintainer
		newReq.OperatorID = assetSignInfo.Maintainer.ID
	}

	// 扫码出库限制为同类型
	if err = b.transferLimit(req); err != nil {
		return err
	}

	_, failed, err := service.NewAssetTransfer().Transfer(b.ctx, &newReq, &md)
	if err != nil {
		return err
	}
	if len(failed) > 0 {
		return errors.New(failed[0])
	}

	return
}

// TransferForUse 运维确认领用创建调拨单
func (b *assetTransferBiz) TransferForUse(assetSignInfo definition.AssetSignInfo, req *definition.AssetTransferCreateReq) (err error) {
	if assetSignInfo.Maintainer == nil {
		return errors.New("运维人员信息有误")
	}
	var md model.Modifier

	// 接收资产map
	var receiveMap = make(map[uint64][]*model.AssetTransferReceiveDetail)
	// 未调拨资产切片
	var notTransferSlice []string
	for _, v := range req.Details {
		if v.AssetType != model.AssetTypeSmartBattery && v.AssetType != model.AssetTypeNonSmartBattery {
			return errors.New("运维端当前只支持电池调拨")
		}
		if v.SN != nil {
			// 查询是否已经调拨
			at, _ := service.NewAssetTransfer().QueryTransferBySN(b.ctx, *v.SN)
			if at != nil {
				receiveMap[at.ID] = append(receiveMap[at.ID], &model.AssetTransferReceiveDetail{
					SN:         v.SN,
					AssetType:  v.AssetType,
					MaterialID: v.MaterialID,
					ModelID:    v.ModelID,
					Num:        v.Num,
				})
			} else {
				notTransferSlice = append(notTransferSlice, *v.SN)
			}
		}
	}

	md = model.Modifier{
		ID:    assetSignInfo.Maintainer.ID,
		Name:  assetSignInfo.Maintainer.Name,
		Phone: assetSignInfo.Maintainer.Phone,
	}

	// 扫码出库限制为同类型
	if err = b.transferLimit(req); err != nil {
		return err
	}

	if len(notTransferSlice) > 0 {
		newReq := model.AssetTransferCreateReq{
			FromLocationType:  req.FromLocationType,
			FromLocationID:    req.FromLocationID,
			ToLocationType:    req.ToLocationType,
			ToLocationID:      req.ToLocationID,
			Details:           req.Details,
			Reason:            req.Reason,
			AssetTransferType: model.AssetTransferTypeTransfer,
			OperatorType:      model.OperatorTypeMaintainer,
			OperatorID:        assetSignInfo.Maintainer.ID,
			AutoIn:            true,
		}
		_, failed, err := service.NewAssetTransfer().Transfer(b.ctx, &newReq, &md)
		if err != nil {
			return err
		}
		if len(failed) > 0 {
			return errors.New(failed[0])
		}
	}
	if len(receiveMap) > 0 {
		assetTransferReceiveReq := make([]model.AssetTransferReceiveReq, 0)
		for k, v := range receiveMap {
			assetTransferReceiveDetail := make([]model.AssetTransferReceiveDetail, 0)
			for _, vl := range v {
				d := vl
				assetTransferReceiveDetail = append(assetTransferReceiveDetail, *d)
			}
			assetTransferReceiveReq = append(assetTransferReceiveReq, model.AssetTransferReceiveReq{
				ID:     k,
				Detail: assetTransferReceiveDetail,
			})
		}
		err = NewAssetTransfer().TransferReceive(assetSignInfo, &definition.AssetTransferReceiveBatchReq{
			AssetTransferReceive: assetTransferReceiveReq,
		})
		if err != nil {
			return err
		}
	}

	return
}

// Transfer 创建资产调拨
func (b *assetTransferBiz) transferLimit(req *definition.AssetTransferCreateReq) (err error) {
	var assetType model.AssetType
	for k, v := range req.Details {
		if k == 0 {
			assetType = v.AssetType
		}
		if v.AssetType != assetType {
			return errors.New("单次调拨只支持同种物资类型")
		}
	}
	return
}

// Flow 资产流转明细
func (b *assetTransferBiz) Flow(req *model.AssetTransferFlowReq) *model.PaginationRes {
	return service.NewAssetTransfer().Flow(b.ctx, req)
}

// GetTransferBySn 扫码入库根据Sn获取调拨信息
func (b *assetTransferBiz) GetTransferBySn(assetSignInfo definition.AssetSignInfo, req *model.GetTransferBySNReq) (res *model.AssetTransferListRes, err error) {
	return service.NewAssetTransfer().GetTransferBySN(assetSignInfo, b.ctx, req)
}

// TransferDetailsList 出入库明细列表
func (b *assetTransferBiz) TransferDetailsList(assetSignInfo definition.AssetSignInfo, req *definition.AssetTransferDetailListReq) (res *model.PaginationRes, err error) {
	newReq := model.AssetTransferDetailListReq{
		PaginationReq:     req.PaginationReq,
		AssetTransferType: req.AssetTransferType,
		Start:             req.Start,
		End:               req.End,
		AssetType:         req.AssetType,
		CabinetSN:         req.Keyword,
		SN:                req.Keyword,
	}
	if assetSignInfo.AssetManager != nil {
		newReq.AssetManagerID = assetSignInfo.AssetManager.ID
	}
	if assetSignInfo.Employee != nil {
		newReq.EmployeeID = assetSignInfo.Employee.ID
	}
	if assetSignInfo.Agent != nil {
		newReq.AgentID = assetSignInfo.Agent.ID
	}
	if assetSignInfo.Maintainer != nil {
		newReq.MaintainerID = assetSignInfo.Maintainer.ID
	}

	return service.NewAssetTransfer().TransferDetailsList(b.ctx, &newReq)
}

// checkPermission 检查操作权限
func (b *assetTransferBiz) checkPermission(assetSignInfo definition.AssetSignInfo, atID uint64) (md model.Modifier, err error) {

	if assetSignInfo.AssetManager != nil {
		// 检验是否有修改权限
		if v, _ := b.orm.QueryNotDeleted().Where(
			assettransfer.ID(atID),
			assettransfer.HasFromLocationWarehouseWith(warehouse.HasBelongAssetManagersWith(assetmanager.ID(assetSignInfo.AssetManager.ID))),
		).First(b.ctx); v == nil {
			return md, errors.New("当前调拨单无权限")
		}
		md = model.Modifier{
			ID:    assetSignInfo.AssetManager.ID,
			Name:  assetSignInfo.AssetManager.Name,
			Phone: assetSignInfo.AssetManager.Phone,
		}
	}
	if assetSignInfo.Employee != nil {
		// 检验是否有修改权限
		if v, _ := b.orm.QueryNotDeleted().Where(
			assettransfer.ID(atID),
			assettransfer.HasFromLocationStoreWith(store.HasEmployeesWith(employee.ID(assetSignInfo.Employee.ID))),
		).First(b.ctx); v == nil {
			return md, errors.New("当前调拨单无权限")
		}
		md = model.Modifier{
			ID:    assetSignInfo.Employee.ID,
			Name:  assetSignInfo.Employee.Name,
			Phone: assetSignInfo.Employee.Phone,
		}
	}
	if assetSignInfo.Agent != nil {
		// 检验是否有修改权限
		// 查询代理人员配置的代理站点
		ids := make([]uint64, 0)
		ag, _ := ent.Database.Agent.QueryNotDeleted().
			WithEnterprise(func(query *ent.EnterpriseQuery) {
				query.WithStations()
			}).
			Where(
				agent.ID(assetSignInfo.Agent.ID),
			).First(context.Background())
		if ag != nil && ag.Edges.Enterprise != nil {
			for _, v := range ag.Edges.Enterprise.Edges.Stations {
				ids = append(ids, v.ID)
			}
		}
		if v, _ := b.orm.QueryNotDeleted().Where(
			assettransfer.ID(atID),
			assettransfer.HasFromLocationStationWith(enterprisestation.IDIn(ids...)),
		).First(b.ctx); v == nil {
			return md, errors.New("当前调拨单无权限")
		}
		md = model.Modifier{
			ID:    assetSignInfo.Agent.ID,
			Name:  assetSignInfo.Agent.Name,
			Phone: assetSignInfo.Agent.Phone,
		}
	}
	if assetSignInfo.Maintainer != nil {
		// 检验是否有修改权限
		if v, _ := b.orm.QueryNotDeleted().Where(
			assettransfer.ID(atID),
			assettransfer.HasFromLocationOperatorWith(maintainer.ID(assetSignInfo.Maintainer.ID)),
		).First(b.ctx); v == nil {
			return md, errors.New("当前调拨单无权限")
		}
		md = model.Modifier{
			ID:    assetSignInfo.Maintainer.ID,
			Name:  assetSignInfo.Maintainer.Name,
			Phone: assetSignInfo.Maintainer.Phone,
		}
	}

	return
}

// Modify 编辑调拨
func (b *assetTransferBiz) Modify(assetSignInfo definition.AssetSignInfo, req *definition.AssetTransferModifyReq) (err error) {
	newReq := model.AssetTransferModifyReq{
		ID:             req.ID,
		ToLocationType: req.ToLocationType,
		ToLocationID:   req.ToLocationID,
		Reason:         req.Reason,
		Remark:         req.Remark,
	}

	var md model.Modifier
	md, err = b.checkPermission(assetSignInfo, req.ID)
	if err != nil {
		return
	}

	return b.modify(&newReq, &md)
}

// modify 根据角色编辑调拨单
func (b *assetTransferBiz) modify(req *model.AssetTransferModifyReq, modifier *model.Modifier) (err error) {
	item, _ := b.orm.QueryNotDeleted().
		Where(
			assettransfer.ID(req.ID),
			assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
		).
		First(b.ctx)
	if item == nil {
		return errors.New("调拨单不存在或已入库")
	}

	// 修改调拨单
	_, err = item.Update().
		SetReason(req.Reason).
		SetNillableRemark(req.Remark).
		SetToLocationID(req.ToLocationID).
		SetToLocationType(req.ToLocationType.Value()).
		SetLastModifier(modifier).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return nil
}

// TransferCancel 取消调拨
func (b *assetTransferBiz) TransferCancel(assetSignInfo definition.AssetSignInfo, req *model.AssetTransferDetailReq) (err error) {
	newReq := model.AssetTransferDetailReq{
		ID: req.ID,
	}
	var md model.Modifier
	md, err = b.checkPermission(assetSignInfo, req.ID)
	if err != nil {
		return
	}

	return service.NewAssetTransfer().TransferCancel(b.ctx, &newReq, &md)
}
