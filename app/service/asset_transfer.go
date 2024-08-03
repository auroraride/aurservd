package service

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
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

// Transfer 调拨
func (s *assetTransferService) Transfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (failed []string, err error) {
	// 调拨限制
	err = s.transferLimit(ctx, req)
	if err != nil {
		return nil, err
	}

	var assetIDs []uint64
	// 创建调拨记录
	q := s.orm.Create()
	// 已经入库资产调拨
	if req.FromLocationType != nil && req.FromLocationID != nil {
		assetIDs, failed, err = s.stockTransfer(ctx, req, modifier)
		if err != nil || len(failed) > 0 {
			return failed, err
		}
		q.SetFromLocationType(req.FromLocationType.Value()).
			SetFromLocationID(*req.FromLocationID).
			SetStatus(model.AssetTransferStatusDelivering.Value()).
			SetOutTimeAt(time.Now()).
			SetOutUserID(modifier.ID).
			SetOutRoleType(model.AssetOperateRoleAdmin.Value()).
			SetOutNum(uint(len(assetIDs)))
	}
	// 初始调拨
	if req.FromLocationType == nil {
		assetIDs, failed, err = s.initialTransfer(ctx, req, modifier)
		if err != nil || len(failed) > 0 {
			return failed, err
		}
		q.SetInNum(uint(len(assetIDs))).
			SetInTimeAt(time.Now()).
			SetStatus(model.AssetTransferStatusStock.Value()).
			SetInUserID(modifier.ID).
			SetInRoleType(model.AssetOperateRoleAdmin.Value()).
			SetRemark("初始化调拨")
	}

	bulk := make([]*ent.AssetTransferDetailsCreate, 0, len(assetIDs))
	for _, id := range assetIDs {
		bulk = append(bulk, ent.Database.AssetTransferDetails.Create().SetAssetID(id).SetCreator(modifier).SetLastModifier(modifier))
	}
	details, _ := ent.Database.AssetTransferDetails.CreateBulk(bulk...).Save(ctx)
	if len(details) == 0 {
		return nil, errors.New("调拨失败")
	}

	err = q.SetToLocationType(req.ToLocationType.Value()).
		SetToLocationID(req.ToLocationID).
		SetSn(tools.NewUnique().NewSN28()).
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetReason(req.Reason).
		AddDetails(details...).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return
}

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
		return nil, errors.New(req.AssetType.String() + "分类ID不能为空")
	}
	// 判定其它物资类型是否存在
	item, _ := ent.Database.Material.QueryNotDeleted().Where(material.ID(*req.MaterialID), material.Type(req.AssetType.Value())).First(ctx)
	if item == nil {
		return nil, errors.New(req.AssetType.String() + "分类不存在")
	}
	// 查询其它物资是否充足
	all, _ := q.Where(asset.MaterialID(*req.MaterialID)).Limit(int(*req.Num)).All(ctx)
	// 查询出的物资数量小于调拨数量 则调拨失败
	if len(all) < int(*req.Num) {
		return nil, errors.New(req.AssetType.String() + "物资不足")
	}
	assetIds = make([]uint64, 0, len(all))
	for _, v := range all {
		assetIds = append(assetIds, v.ID)
	}
	return assetIds, nil
}

// 检查目标位置是否存在的函数
func (s *assetTransferService) checkTargetLocationExists(ctx context.Context, locationType model.AssetLocationsType, locationID uint64) error {
	var exists bool
	switch locationType {
	case model.AssetLocationsTypeWarehouse:
		exists, _ = ent.Database.Warehouse.QueryNotDeleted().Where(warehouse.ID(locationID)).Exist(ctx)
	case model.AssetLocationsTypeStore:
		exists, _ = ent.Database.Store.QueryNotDeleted().Where(store.ID(locationID)).Exist(ctx)
	case model.AssetLocationsTypeStation:
		exists, _ = ent.Database.EnterpriseStation.QueryNotDeleted().Where(enterprisestation.ID(locationID)).Exist(ctx)
	case model.AssetLocationsTypeOperation:
		exists, _ = ent.Database.Maintainer.Query().Where(maintainer.ID(locationID)).Exist(ctx)
	default:
		return errors.New("调拨目标地点不存在")
	}
	if !exists {
		return errors.New("调拨目标不存在")
	}
	return nil
}

// 调拨限制
func (s *assetTransferService) transferLimit(ctx context.Context, req *model.AssetTransferCreateReq) (err error) {
	if req.FromLocationType != nil {
		if *req.FromLocationID == req.ToLocationID {
			return errors.New("无法调拨到相同位置")
		}
		// 仓库限制（仓库、门店、站点、运维）
		if *req.FromLocationType == model.AssetLocationsTypeWarehouse {
			switch req.ToLocationType {
			case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeStore, model.AssetLocationsTypeStation, model.AssetLocationsTypeOperation:
				if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
					return err
				}
			default:
				return errors.New("调拨目标地点不合法")
			}
		}
		// 门店（仓库、门店、运维） 运维（仓库、门店、运维）
		if *req.FromLocationType == model.AssetLocationsTypeStore || *req.FromLocationType == model.AssetLocationsTypeOperation {
			switch req.ToLocationType {
			case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeStore, model.AssetLocationsTypeOperation:
				if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
					return err
				}
			default:
				return errors.New("调拨目标地点不合法")
			}
		}
		// 站点（仓库、相同代理商其他站点）
		if *req.FromLocationType == model.AssetLocationsTypeStation {
			switch req.ToLocationType {
			case model.AssetLocationsTypeWarehouse:
				if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
					return err
				}
			case model.AssetLocationsTypeStation:
				item, _ := ent.Database.EnterpriseStation.QueryNotDeleted().WithEnterprise(func(query *ent.EnterpriseQuery) {
					query.WithStations()
				}).Where(enterprisestation.ID(req.ToLocationID)).First(ctx)
				if item == nil {
					return errors.New("站点不存在")
				}
				// 查询出该代理所有站点
				if item.Edges.Enterprise != nil && item.Edges.Enterprise.Edges.Stations != nil {
					var in bool
					for _, v := range item.Edges.Enterprise.Edges.Stations {
						if v.ID == req.ToLocationID {
							in = true
							break
						}
					}
					if !in {
						return errors.New("只能调拨到相同代理商的站点")
					}
				}
			default:
				return errors.New("调拨目标地点不合法")
			}
		}
	}
	if req.FromLocationType == nil {
		// 初始调拨只能调仓库
		if req.ToLocationType != model.AssetLocationsTypeWarehouse {
			return errors.New("调拨目标地点不存在")
		}
		if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
			return err
		}
	}
	return nil
}

// 已入库调拨
func (s *assetTransferService) stockTransfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (assetIDs []uint64, failed []string, err error) {
	// 查询物资是否充足
	q := ent.Database.Asset.QueryNotDeleted().Where(
		asset.LocationsType((*req.FromLocationType).Value()),
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
	)

	for _, v := range req.Details {
		var iDs []uint64
		switch v.AssetType {
		case model.AssetTypeEbike, model.AssetTypeSmartBattery:
			iDs, err = s.transferAssetWithSN(ctx, q, v, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		case model.AssetTypeNonSmartBattery, model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeOtherAccessory:
			iDs, err = s.transferAssetWithoutSN(ctx, q, v, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		default:
		}
		assetIDs = append(assetIDs, iDs...)
	}
	return assetIDs, nil, nil
}

// 初始调拨
func (s *assetTransferService) initialTransfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (assetIDs []uint64, failed []string, err error) {
	// 创建物资
	for _, v := range req.Details {
		var iDs []uint64
		switch v.AssetType {
		case model.AssetTypeNonSmartBattery, model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeOtherAccessory:
			iDs, err = s.initialTransferWithoutSN(ctx, v, req.ToLocationID, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		default:
			failed = append(failed, v.AssetType.String()+"物资类型不合法,已跳过")
		}
		assetIDs = append(assetIDs, iDs...)
	}
	return assetIDs, failed, nil
}

// initialTransferWithoutSN 无编号资产初始化调拨
func (s *assetTransferService) initialTransferWithoutSN(ctx context.Context, req model.AssetTransferCreateDetail, toLocationID uint64, modifier *model.Modifier) (assetIDs []uint64, err error) {
	if req.Num == nil || *req.Num == 0 {
		return nil, errors.New("调拨数量不能为空")
	}
	if req.MaterialID == nil || *req.MaterialID == 0 {
		return nil, errors.New(req.AssetType.String() + "分类ID不能为空")
	}
	// 判定其它物资类型是否存在
	item, _ := ent.Database.Material.QueryNotDeleted().Where(material.ID(*req.MaterialID), material.Type(req.AssetType.Value())).First(ctx)
	if item == nil {
		return nil, errors.New(req.AssetType.String() + "分类不存在")
	}
	// 创建物资
	bulk := make([]*ent.AssetCreate, 0, int(*req.Num))
	for i := 0; i < int(*req.Num); i++ {
		bulk = append(bulk, ent.Database.Asset.Create().
			SetType(req.AssetType.Value()).
			SetMaterialID(*req.MaterialID).
			SetStatus(model.AssetStatusStock.Value()).
			SetEnable(true).
			SetCreator(modifier).
			SetLastModifier(modifier).
			SetLocationsType(model.AssetLocationsTypeWarehouse.Value()).
			SetLocationsID(toLocationID).
			SetName(item.Name))
	}
	assets, _ := ent.Database.Asset.CreateBulk(bulk...).Save(ctx)
	if len(assets) == 0 {
		return nil, errors.New("创建资产失败")
	}
	for _, v := range assets {
		assetIDs = append(assetIDs, v.ID)
	}
	return assetIDs, nil
}
