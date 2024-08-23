// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-20, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
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
func (b *assetTransferBiz) TransferList(am *ent.AssetManager, ep *ent.Employee, req *definition.TransferListReq) (res *model.PaginationRes, err error) {
	newReq := model.AssetTransferListReq{
		PaginationReq:       req.PaginationReq,
		AssetTransferFilter: req.AssetTransferFilter,
	}

	if am != nil {
		newReq.AssetManagerID = am.ID
	}

	if ep != nil {
		newReq.EmployeeID = ep.ID
	}

	return service.NewAssetTransfer().TransferList(context.Background(), &newReq)
}

// TransferDetail 调拨记录详情
func (b *assetTransferBiz) TransferDetail(ctx context.Context, req *model.AssetTransferDetailReq) (res *definition.TransferDetailRes, err error) {
	var t *ent.AssetTransfer
	t, err = ent.Database.AssetTransfer.QueryNotDeleted().Where(assettransfer.ID(req.ID)).First(ctx)
	if err != nil {
		return nil, err
	}

	details, err := service.NewAssetTransfer().TransferDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	atr := service.NewAssetTransfer().TransferInfo(t)
	res = &definition.TransferDetailRes{
		AssetTransferListRes: *atr,
		Detail:               details,
	}
	return
}

// TransferReceive 接收资产
func (b *assetTransferBiz) TransferReceive(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetTransferReceiveBatchReq) (err error) {
	var md model.Modifier

	newReq := model.AssetTransferReceiveBatchReq{
		AssetTransferReceive: req.AssetTransferReceive,
	}

	if am != nil {
		newReq.OperateType = model.AssetOperateRoleTypeManager
		md = model.Modifier{
			ID:    am.ID,
			Name:  am.Name,
			Phone: am.Phone,
		}
	}
	if ep != nil {
		newReq.OperateType = model.AssetOperateRoleTypeStore
		md = model.Modifier{
			ID:    ep.ID,
			Name:  ep.Name,
			Phone: ep.Phone,
		}
	}

	return service.NewAssetTransfer().TransferReceive(b.ctx, &newReq, &md)
}

// Transfer 创建资产调拨
func (b *assetTransferBiz) Transfer(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetTransferCreateReq) (failed []string, err error) {

	var md model.Modifier

	aType := model.AssetTransferTypeTransfer
	newReq := model.AssetTransferCreateReq{
		ToLocationType:    req.ToLocationType,
		ToLocationID:      req.ToLocationID,
		Details:           req.Details,
		Reason:            req.Reason,
		AssetTransferType: &aType,
	}

	if am != nil {
		wType := model.AssetLocationsTypeWarehouse
		newReq.FromLocationType = &wType
		if req.FromLocationID == nil {
			req.FromLocationID = silk.UInt64(am.ID)
		} else {
			newReq.FromLocationID = req.FromLocationID
		}

		md = model.Modifier{
			ID:    am.ID,
			Name:  am.Name,
			Phone: am.Phone,
		}
	}
	if ep != nil {
		sType := model.AssetLocationsTypeStore
		newReq.FromLocationType = &sType
		if req.FromLocationID == nil {
			req.FromLocationID = silk.UInt64(ep.ID)
		} else {
			newReq.FromLocationID = req.FromLocationID
		}
		md = model.Modifier{
			ID:    ep.ID,
			Name:  ep.Name,
			Phone: ep.Phone,
		}
	}

	return service.NewAssetTransfer().Transfer(b.ctx, &newReq, &md)
}

// Flow 资产流转明细
func (b *assetTransferBiz) Flow(req *model.AssetTransferFlowReq) []*model.AssetTransferFlow {
	return service.NewAssetTransfer().Flow(b.ctx, req)
}

// GetTransferBySn 扫码入库根据Sn获取调拨信息
func (b *assetTransferBiz) GetTransferBySn(req *model.GetTransferBySNReq) (res *model.AssetTransferListRes, err error) {
	return service.NewAssetTransfer().GetTransferBySN(b.ctx, req)
}

// TransferDetailsList 出入库明细列表
func (b *assetTransferBiz) TransferDetailsList(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetTransferDetailListReq) (res *model.PaginationRes, err error) {
	newReq := model.AssetTransferDetailListReq{
		PaginationReq:     req.PaginationReq,
		AssetTransferType: req.AssetTransferType,
		Start:             req.Start,
		End:               req.End,
		AssetType:         req.AssetType,
		CabinetSN:         req.Keyword,
		SN:                req.Keyword,
	}
	if am != nil {
		newReq.AssetManagerID = am.ID
	}
	if ep != nil {
		newReq.EmployeeID = ep.ID
	}

	return service.NewAssetTransfer().TransferDetailsList(b.ctx, &newReq)
}
