package service

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetTransferService struct {
	orm *ent.AssetTransferClient
}

func NewAssetTransfer() *assetTransferService {
	return &assetTransferService{
		orm: ent.Database.AssetTransfer,
	}
}

// TransferAsset 调拨
func (s *assetTransferService) TransferAsset(ctx context.Context, req model.AssetTransferCreateReq, modifier *model.Modifier) (failed []string, err error) {
	// 查询物资是否充足
	q := ent.Database.Asset.QueryNotDeleted().Where(
		asset.LocationsType(req.FromLocationType.Value()),
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
	)
	var assetIDs []uint64
	for _, v := range req.Details {
		var iDs []uint64
		switch v.AssetType {
		case model.AssetTypeEbike, model.AssetTypeSmartBattery:
			iDs, err = s.transferAssetWithSN(ctx, q, v, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		case model.AssetTypeCabinetAccessory, model.AssetTypeOtherAccessory, model.AssetTypeNonSmartBattery:
			iDs, err = s.transferAssetWithoutSN(ctx, q, v, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		default:
		}
		assetIDs = append(assetIDs, iDs...)
	}
	bulk := make([]*ent.AssetTransferDetailsCreate, 0, len(assetIDs))
	for _, id := range assetIDs {
		bulk = append(bulk, ent.Database.AssetTransferDetails.Create().SetAssetID(id).SetCreator(modifier).SetLastModifier(modifier))
	}
	details, _ := ent.Database.AssetTransferDetails.CreateBulk(bulk...).Save(ctx)
	if len(details) == 0 {
		return failed, errors.New("调拨失败")
	}
	locationType := (*req.FromLocationType).Value()

	// 创建调拨记录
	err = s.orm.Create().
		SetNillableFromLocationType(&locationType).
		SetNillableFromLocationID(req.FromLocationID).
		SetToLocationType(req.ToLocationType.Value()).
		SetToLocationID(req.ToLocationID).
		SetStatus(model.AssetTransferStatusDelivering.Value()).
		SetSn(tools.NewUnique().NewSN28()).
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetInNum(uint(len(assetIDs))).
		SetInTimeAt(time.Now()).
		SetInUserID(modifier.ID).
		SetInRoleType(model.AssetOperateRoleAdmin.Value()).
		AddDetails(details...).
		Exec(ctx)
	return
}

// 判定调拨类型 todo 有问题
// func (s *assetTransferService) transferAssetType(ctx context.Context, req model.AssetTransferCreateReq) (res *model.AssetTransferType, err error) {
//
// 	// 初始入库，指平台新增的物资，第一次入库；初始入库，没有调出目标
// 	var transferType model.AssetTransferType
// 	if req.FromLocationType == nil && req.ToLocationType == model.AssetLocationsTypeWarehouse {
// 		transferType = model.AssetTransferTypeInitial
// 		return &transferType, nil
// 	}
//
// 	// 只有初始入库才有调拨前位置类型为空
// 	if req.FromLocationType != nil {
// 		return nil, errors.New("调拨前位置类型不能为空")
// 	}
//
// 	// 门店调拨，指门店间的出入库
// 	if *req.FromLocationType == model.AssetLocationsTypeStore && req.ToLocationType == model.AssetLocationsTypeStore {
// 		transferType = model.AssetTransferTypeStore
// 		return &transferType, nil
// 	}
//
// 	// 平台调拨，指仓库与其他库存点的出入库
// 	if *req.FromLocationType == model.AssetLocationsTypeWarehouse && req.ToLocationType != model.AssetLocationsTypeWarehouse {
// 		transferType = model.AssetTransferTypePlatform
// 		return &transferType, nil
// 	}
//
// 	// 运维领取，指运维从仓库/门店领取物资；运维退还，指运维将物资还回仓库或门店
// 	return
// }

// 有资产编号的物资调拨
func (s *assetTransferService) transferAssetWithSN(ctx context.Context, q *ent.AssetQuery, req model.AssetTransferCreateDetail, modifier *model.Modifier) (assetIDs []uint64, err error) {
	if req.SN == nil || *req.SN == "" {
		return nil, errors.New("资产编号不能为空")
	}
	// 查询物资是否存在
	item, _ := q.Where(asset.Sn(*req.SN)).First(ctx)
	if item == nil {
		return nil, errors.New("物资不存在")
	}

	assetIDs = append(assetIDs, item.ID)

	return
}

// 无资产编号的物资调拨
func (s *assetTransferService) transferAssetWithoutSN(ctx context.Context, q *ent.AssetQuery, req model.AssetTransferCreateDetail, modifier *model.Modifier) (assetIds []uint64, err error) {
	if req.Num == nil || *req.Num == 0 {
		return nil, errors.New("调拨数量不能为空")
	}
	if req.MaterialID == nil || *req.MaterialID == 0 {
		return nil, errors.New("其它物资分类ID不能为空")
	}
	// 查询其它物资是否充足
	all, _ := q.Where(asset.MaterialID(*req.MaterialID)).Limit(int(*req.Num)).All(ctx)
	// 查询出的物资数量小于调拨数量 则调拨失败
	if len(all) < int(*req.Num) {
		return nil, errors.New("其它物资不足")
	}
	assetIds = make([]uint64, 0, len(all))
	for _, v := range all {
		assetIds = append(assetIds, v.ID)
	}
	return assetIds, nil
}
