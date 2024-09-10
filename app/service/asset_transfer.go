package service

import (
	"context"
	"errors"
	"slices"
	"sort"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetTransferService struct {
	orm *ent.AssetTransferClient
}

func NewAssetTransfer(params ...any) *assetTransferService {
	return &assetTransferService{
		orm: ent.Database.AssetTransfer,
	}
}

// Transfer 调拨
func (s *assetTransferService) Transfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (sn string, failed []string, err error) {
	// 调拨限制
	err = s.transferLimit(ctx, req)
	if err != nil {
		return "", nil, err
	}

	var assetIDs []uint64
	var newTime = time.Now()
	bulk := make([]*ent.AssetTransferDetailsCreate, 0, len(assetIDs))
	// 创建调拨记录
	q := s.orm.Create()
	// 已经入库资产调拨
	if req.FromLocationType != nil && req.FromLocationID != nil {
		assetIDs, failed = s.stockTransfer(ctx, req, modifier)
		if len(failed) > 0 {
			return "", failed, nil
		}
		q.SetFromLocationType(req.FromLocationType.Value()).
			SetFromLocationID(*req.FromLocationID).
			SetStatus(model.AssetTransferStatusDelivering.Value()).
			SetOutTimeAt(newTime).
			SetOutOperateID(req.OperatorID).
			SetOutOperateType(req.OperatorType.Value()).
			SetOutNum(uint(len(assetIDs))).
			SetType(req.AssetTransferType.Value())
		for _, id := range assetIDs {
			d := ent.Database.AssetTransferDetails.Create().
				SetAssetID(id).
				SetCreator(modifier).
				SetLastModifier(modifier).
				SetRemark(req.OperatorType.String() + "调拨")
			bulk = append(bulk, d)
		}

		// 修改资产状态 只有已入库的资产才会修改状态
		_, err = ent.Database.Asset.Update().Where(asset.IDIn(assetIDs...)).
			SetStatus(model.AssetStatusDelivering.Value()).
			SetLastModifier(modifier).
			Save(ctx)
		if err != nil {
			return "", nil, err
		}
	}
	// 初始调拨
	if req.FromLocationType == nil {
		assetIDs, failed = s.initialTransfer(ctx, req, modifier)
		if len(failed) > 0 {
			return "", failed, nil
		}
		q.SetInNum(uint(len(assetIDs))).
			SetStatus(model.AssetTransferStatusStock.Value()).
			SetType(model.AssetTransferTypeInitial.Value()).
			SetRemark("后台初始调拨")
		for _, id := range assetIDs {
			d := ent.Database.AssetTransferDetails.Create().
				SetAssetID(id).
				SetCreator(modifier).
				SetLastModifier(modifier)
			if req.FromLocationType == nil {
				d.SetInTimeAt(newTime).
					SetInOperateID(req.OperatorID).
					SetInOperateType(req.OperatorType.Value()).
					SetIsIn(true).
					SetRemark("后台初始调拨")
			}
			bulk = append(bulk, d)
		}
	}

	details, _ := ent.Database.AssetTransferDetails.CreateBulk(bulk...).Save(ctx)
	if len(details) == 0 {
		return "", nil, errors.New("调拨失败")
	}

	save, err := q.SetToLocationType(req.ToLocationType.Value()).
		SetToLocationID(req.ToLocationID).
		SetSn(tools.NewUnique().NewSN28()).
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetReason(req.Reason).
		AddTransferDetails(details...).
		Save(ctx)
	if err != nil {
		return "", nil, err
	}

	// 自动入库
	if req.AutoIn && req.FromLocationType != nil {
		assetTransferReceive := make([]model.AssetTransferReceiveDetail, 0)
		for _, v := range req.Details {
			assetTransferReceive = append(assetTransferReceive, model.AssetTransferReceiveDetail{
				AssetType:  v.AssetType,
				SN:         v.SN,
				Num:        v.Num,
				MaterialID: v.MaterialID,
			})
		}
		err = s.TransferReceive(ctx, &model.AssetTransferReceiveBatchReq{
			OperateType: req.OperatorType,
			AssetTransferReceive: []model.AssetTransferReceiveReq{
				{
					ID:     save.ID,
					Remark: silk.String("自动入库"),
					Detail: assetTransferReceive,
				},
			},
		}, modifier)
		if err != nil {
			return "", nil, err
		}
	}
	return save.Sn, failed, nil
}

// 有资产编号的物资调拨
func (s *assetTransferService) transferAssetWithSN(ctx context.Context, locType model.AssetLocationsType, locID uint64, req model.AssetTransferCreateDetail, modifier *model.Modifier) (assetIDs []uint64, err error) {
	q := ent.Database.Asset.QueryNotDeleted().Where(
		asset.LocationsType(locType.Value()),
		asset.LocationsID(locID),
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
	)
	if req.SN == nil || *req.SN == "" {
		return nil, errors.New("资产编号不能为空")
	}
	// 查询物资是否存在
	item, _ := q.Where(asset.Sn(*req.SN)).First(ctx)
	if item == nil {
		return nil, errors.New(*req.SN + "物资不存在或不在库存中")
	}
	// 在盘点中的物资不能调拨
	if item.CheckAt != nil {
		return nil, errors.New(*req.SN + "物资正在盘点中")
	}

	assetIDs = append(assetIDs, item.ID)
	return
}

// 无资产编号的物资调拨
func (s *assetTransferService) transferAssetWithoutSN(ctx context.Context, locType model.AssetLocationsType, locID uint64, req model.AssetTransferCreateDetail, modifier *model.Modifier) (assetIds []uint64, err error) {
	if req.Num == nil || *req.Num == 0 {
		return nil, errors.New("调拨数量不能为空")
	}
	if req.AssetType == model.AssetTypeNonSmartBattery {
		// 非智能电池调拨
		if req.ModelID == nil || *req.ModelID == 0 {
			return nil, errors.New(req.AssetType.String() + "型号ID不能为空")
		}
		item, _ := ent.Database.BatteryModel.Query().Where(batterymodel.ID(*req.ModelID)).First(ctx)
		if item == nil {
			return nil, errors.New(req.AssetType.String() + "型号不存在")
		}
	} else {
		if req.MaterialID == nil || *req.MaterialID == 0 {
			return nil, errors.New(req.AssetType.String() + "分类ID不能为空")
		}
		// 判定其它物资类型是否存在
		item, _ := ent.Database.Material.QueryNotDeleted().Where(
			material.ID(*req.MaterialID),
			material.Type(req.AssetType.Value()),
		).First(ctx)
		if item == nil {
			return nil, errors.New(req.AssetType.String() + "分类不存在")
		}
	}
	q := ent.Database.Asset.QueryNotDeleted().Where(
		asset.LocationsType(locType.Value()),
		asset.LocationsID(locID),
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
	)
	if req.ModelID != nil {
		q = q.Where(asset.ModelID(*req.ModelID))
	}
	if req.MaterialID != nil {
		q = q.Where(asset.MaterialID(*req.MaterialID))
	}
	// 查询其它物资是否充足
	all, _ := q.Limit(int(*req.Num)).All(ctx)
	// 查询出的物资数量小于调拨数量 则调拨失败
	if len(all) < int(*req.Num) {
		return nil, errors.New(req.AssetType.String() + "物资不足或不在库存中")
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
	case model.AssetLocationsTypeCabinet:
		exists, _ = ent.Database.Cabinet.Query().Where(cabinet.ID(locationID)).Exist(ctx)
	case model.AssetLocationsTypeRider:
		exists, _ = ent.Database.Rider.Query().Where(rider.ID(locationID)).Exist(ctx)
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
		// 仓库限制（仓库、门店、站点、运维、电柜）
		if *req.FromLocationType == model.AssetLocationsTypeWarehouse {
			switch req.ToLocationType {
			case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeStore, model.AssetLocationsTypeStation, model.AssetLocationsTypeOperation, model.AssetLocationsTypeCabinet:
				if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
					return err
				}
			default:
				return errors.New("调拨目标地点不合法")
			}
		}
		// 门店（仓库、门店、运维、电柜） 运维（仓库、门店、运维、电柜）
		if *req.FromLocationType == model.AssetLocationsTypeStore || *req.FromLocationType == model.AssetLocationsTypeOperation {
			switch req.ToLocationType {
			case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeStore, model.AssetLocationsTypeOperation, model.AssetLocationsTypeCabinet:
				if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
					return err
				}
			default:
				return errors.New("调拨目标地点不合法")
			}
		}
		// 站点（仓库、相同代理商其他站点、电柜）
		if *req.FromLocationType == model.AssetLocationsTypeStation {
			switch req.ToLocationType {
			case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeCabinet:
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
		// 电柜（仓库、门店、站点、运维）
		if *req.FromLocationType == model.AssetLocationsTypeCabinet {
			switch req.ToLocationType {
			case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeStore, model.AssetLocationsTypeStation, model.AssetLocationsTypeOperation:
				if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
					return err
				}
			default:
				return errors.New("调拨目标地点不合法")
			}
		}

	}
	if req.FromLocationType == nil {
		// 初始调拨只能调仓库  todo 暂时去掉吧
		// if req.ToLocationType != model.AssetLocationsTypeWarehouse  || req.ToLocationType!={
		// 	return errors.New("调拨目标地点不存在")
		// }
		if err = s.checkTargetLocationExists(ctx, req.ToLocationType, req.ToLocationID); err != nil {
			return err
		}
	}
	if req.ToLocationType == model.AssetLocationsTypeOperation {
		for _, v := range req.Details {
			if v.AssetType != model.AssetTypeSmartBattery && v.AssetType != model.AssetTypeNonSmartBattery {
				return errors.New("运维只支持电池类型物资调拨")
			}
		}
	}
	return nil
}

// 已入库调拨
func (s *assetTransferService) stockTransfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (assetIDs []uint64, failed []string) {
	var err error
	for _, v := range req.Details {
		var iDs []uint64
		switch v.AssetType {
		case model.AssetTypeEbike, model.AssetTypeSmartBattery:
			iDs, err = s.transferAssetWithSN(ctx, *req.FromLocationType, *req.FromLocationID, v, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		case model.AssetTypeNonSmartBattery, model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeOtherAccessory:
			iDs, err = s.transferAssetWithoutSN(ctx, *req.FromLocationType, *req.FromLocationID, v, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		default:
		}
		assetIDs = append(assetIDs, iDs...)
	}
	// 去重
	var realIds []uint64
	idMap := make(map[uint64]bool)
	for _, v := range assetIDs {
		if !idMap[v] {
			realIds = append(realIds, v)
			idMap[v] = true
		}
	}

	return realIds, failed
}

// 初始调拨
func (s *assetTransferService) initialTransfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (assetIDs []uint64, failed []string) {
	var err error
	// 创建物资
	for _, v := range req.Details {
		var iDs []uint64
		switch v.AssetType {
		case model.AssetTypeNonSmartBattery, model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeOtherAccessory:
			iDs, err = s.initialTransferWithoutSN(ctx, v, req.ToLocationID, req.ToLocationType, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		case model.AssetTypeEbike, model.AssetTypeSmartBattery:
			iDs, err = s.initialTransferWithSN(ctx, v, req.ToLocationID, req.ToLocationType, modifier)
			if err != nil {
				failed = append(failed, err.Error())
				continue
			}
		default:
			failed = append(failed, v.AssetType.String()+"物资类型不合法,已跳过")
		}
		assetIDs = append(assetIDs, iDs...)
	}
	return assetIDs, failed
}

// initialTransferWithoutSN 无编号资产初始化调拨
func (s *assetTransferService) initialTransferWithoutSN(ctx context.Context, req model.AssetTransferCreateDetail, toLocationID uint64, toLocationType model.AssetLocationsType, modifier *model.Modifier) (assetIDs []uint64, err error) {
	if req.Num == nil || *req.Num == 0 {
		return nil, errors.New("调拨数量不能为空")
	}
	var name string
	if req.AssetType == model.AssetTypeNonSmartBattery {
		// 非智能电池调拨
		if req.ModelID == nil || *req.ModelID == 0 {
			return nil, errors.New(req.AssetType.String() + "型号ID不能为空")
		}
		item, _ := ent.Database.BatteryModel.Query().Where(batterymodel.ID(*req.ModelID)).First(ctx)
		if item == nil {
			return nil, errors.New(req.AssetType.String() + "型号不存在")
		}
		name = item.Model
	} else {
		if req.MaterialID == nil || *req.MaterialID == 0 {
			return nil, errors.New(req.AssetType.String() + "分类ID不能为空")
		}
		// 判定其它物资类型是否存在
		item, _ := ent.Database.Material.QueryNotDeleted().Where(
			material.ID(*req.MaterialID),
			material.Type(req.AssetType.Value()),
		).First(ctx)
		if item == nil {
			return nil, errors.New(req.AssetType.String() + "分类不存在")
		}
		name = item.Name
	}

	// 创建物资
	bulk := make([]*ent.AssetCreate, 0, int(*req.Num))
	for i := 0; i < int(*req.Num); i++ {
		bulk = append(bulk, ent.Database.Asset.Create().
			SetType(req.AssetType.Value()).
			SetNillableMaterialID(req.MaterialID).
			SetNillableModelID(req.ModelID).
			SetStatus(model.AssetStatusStock.Value()).
			SetEnable(true).
			SetCreator(modifier).
			SetLastModifier(modifier).
			SetLocationsType(toLocationType.Value()).
			SetLocationsID(toLocationID).
			SetNillableModelID(req.ModelID).
			SetName(name))
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

// 有编号资产初始化调拨
func (s *assetTransferService) initialTransferWithSN(ctx context.Context, req model.AssetTransferCreateDetail, toLocationID uint64, toLocationType model.AssetLocationsType, modifier *model.Modifier) (assetIDs []uint64, err error) {
	if req.SN == nil || *req.SN == "" {
		return nil, errors.New("资产编号不能为空")
	}
	// 查询物资是否存在
	item, _ := ent.Database.Asset.QueryNotDeleted().Where(
		asset.Status(model.AssetStatusDelivering.Value()),
		asset.Sn(*req.SN),
	).First(ctx)
	if item == nil {
		return nil, errors.New(*req.SN + "物资不存在或不在待入库中")
	}
	// 修改状态
	_, err = ent.Database.Asset.Update().Where(asset.ID(item.ID)).
		SetLocationsType(toLocationType.Value()).
		SetLocationsID(toLocationID).
		SetLastModifier(modifier).
		SetStatus(model.AssetStatusStock.Value()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	assetIDs = append(assetIDs, item.ID)
	return
}

// TransferList 调拨列表
func (s *assetTransferService) TransferList(ctx context.Context, req *model.AssetTransferListReq) (res *model.PaginationRes, err error) {
	q := ent.Database.AssetTransfer.QueryNotDeleted().WithTransferDetails().
		WithOutOperateAgent().WithOutOperateManager().WithOutOperateStore().WithOutOperateMaintainer().WithOutOperateRider().WithOutOperateCabinet().
		WithFromLocationWarehouse().WithFromLocationStore().WithFromLocationStation().WithFromLocationOperator().WithFromLocationRider().WithFromLocationCabinet().
		WithToLocationWarehouse().WithToLocationStore().WithToLocationStation().WithToLocationOperator().WithToLocationRider().WithToLocationCabinet()
	s.filter(ctx, q, &req.AssetTransferFilter)
	q.Order(ent.Desc(assettransfer.FieldCreatedAt))
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetTransfer) (res *model.AssetTransferListRes) {
		return s.TransferInfo(nil, item)
	}), nil

}

// TransferInfo 补充完善调拨信息
func (s *assetTransferService) TransferInfo(atu *model.AssetTransferUserId, item *ent.AssetTransfer) (res *model.AssetTransferListRes) {
	res = &model.AssetTransferListRes{
		ID:                item.ID,
		SN:                item.Sn,
		Reason:            item.Reason,
		Remark:            item.Remark,
		Status:            model.AssetTransferStatus(item.Status).String(),
		AssetTransferType: model.AssetTransferType(item.Type),
		OutNum:            item.OutNum,
		InNum:             item.InNum,
	}

	if item.OutTimeAt != nil {
		res.OutTimeAt = item.OutTimeAt.Format("2006-01-02 15:04:05")
	}
	if item.FromLocationType != nil && item.FromLocationID != nil {
		res.FromLocationType = *item.FromLocationType
		res.FromLocationID = *item.FromLocationID
		switch model.AssetLocationsType(*item.FromLocationType) {
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.FromLocationWarehouse != nil {
				res.FromLocationName = "[仓库]" + item.Edges.FromLocationWarehouse.Name
			}
		case model.AssetLocationsTypeStore:
			if item.Edges.FromLocationStore != nil {
				res.FromLocationName = "[门店]" + item.Edges.FromLocationStore.Name
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.FromLocationStation != nil {
				res.FromLocationName = "[站点]" + item.Edges.FromLocationStation.Name
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.FromLocationOperator != nil {
				res.FromLocationName = "[运维]" + item.Edges.FromLocationOperator.Name
			}
		case model.AssetLocationsTypeRider:
			if item.Edges.FromLocationRider != nil {
				res.FromLocationName = "[骑手]" + item.Edges.FromLocationRider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if item.Edges.FromLocationCabinet != nil {
				res.FromLocationName = "[电柜]" + item.Edges.FromLocationCabinet.Name
			}
		default:
		}
	}

	if item.ToLocationType != 0 && item.ToLocationID != 0 {
		res.ToLocationType = item.ToLocationType
		res.ToLocationID = item.ToLocationID
		switch model.AssetLocationsType(item.ToLocationType) {
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.ToLocationWarehouse != nil {
				res.ToLocationName = "[仓库]" + item.Edges.ToLocationWarehouse.Name
			}
		case model.AssetLocationsTypeStore:
			if item.Edges.ToLocationStore != nil {
				res.ToLocationName = "[门店]" + item.Edges.ToLocationStore.Name
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.ToLocationStation != nil {
				res.ToLocationName = "[站点]" + item.Edges.ToLocationStation.Name
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.ToLocationOperator != nil {
				res.ToLocationName = "[运维]" + item.Edges.ToLocationOperator.Name
			}
		case model.AssetLocationsTypeRider:
			if item.Edges.ToLocationRider != nil {
				res.ToLocationName = "[骑手]" + item.Edges.ToLocationRider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if item.Edges.ToLocationCabinet != nil {
				res.ToLocationName = "[电柜]" + item.Edges.ToLocationCabinet.Name
			}
		default:
		}
	}

	// 出库操作人
	if item.OutOperateType != nil && item.OutOperateID != nil {
		switch model.OperatorType(*item.OutOperateType) {
		case model.OperatorTypeAssetManager:
			if item.Edges.OutOperateManager != nil {
				res.OutOperateName = "[仓管]" + item.Edges.OutOperateManager.Name
			}
		case model.OperatorTypeEmployee:
			if item.Edges.OutOperateStore != nil {
				res.OutOperateName = "[门店]" + item.Edges.OutOperateStore.Name
			}
		case model.OperatorTypeAgent:
			if item.Edges.OutOperateAgent != nil {
				res.OutOperateName = "[代理]" + item.Edges.OutOperateAgent.Name
			}
		case model.OperatorTypeMaintainer:
			if item.Edges.OutOperateMaintainer != nil {
				res.OutOperateName = "[运维]" + item.Edges.OutOperateMaintainer.Name
			}
		case model.OperatorTypeRider:
			if item.Edges.OutOperateRider != nil {
				res.OutOperateName = "[骑手]" + item.Edges.OutOperateRider.Name
			}
		case model.OperatorTypeCabinet:
			if item.Edges.OutOperateCabinet != nil {
				res.OutOperateName = "[电柜]" + item.Edges.OutOperateCabinet.Name
			}
		default:
		}
	}

	res.InOut = s.checkInOutType(atu, item)
	return res
}

// 筛选
func (s *assetTransferService) checkInOutType(atu *model.AssetTransferUserId, item *ent.AssetTransfer) string {
	if atu == nil || item == nil {
		return ""
	}
	var in, out bool

	if item.FromLocationType != nil && item.FromLocationID != nil {
		switch model.AssetLocationsType(*item.FromLocationType) {
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.FromLocationWarehouse != nil {
				// if data, _ := ent.Database.Warehouse.QueryNotDeleted().
				// 	Where(
				// 		warehouse.ID(*item.FromLocationID),
				// 		warehouse.HasBelongAssetManagersWith(assetmanager.ID(atu.AssetManagerID)),
				// 	).First(context.Background()); data != nil {
				// 	out = true
				// }

				// 根据上班地点确认inout
				if data, _ := ent.Database.AssetManager.QueryNotDeleted().Where(
					assetmanager.ID(atu.AssetManagerID),
					assetmanager.WarehouseID(*item.FromLocationID),
				).First(context.Background()); data != nil {
					out = true
				}
			}
		case model.AssetLocationsTypeStore:
			if item.Edges.FromLocationStore != nil {
				// if data, _ := ent.Database.Store.QueryNotDeleted().
				// 	Where(
				// 		store.ID(*item.FromLocationID),
				// 		store.HasEmployeesWith(employee.ID(atu.EmployeeID)),
				// 	).First(context.Background()); data != nil {
				// 	out = true
				// }

				// 根据上班地点确认inout
				if data, _ := ent.Database.Employee.QueryNotDeleted().Where(
					employee.ID(atu.EmployeeID),
					employee.DutyStoreID(*item.FromLocationID),
				).First(context.Background()); data != nil {
					out = true
				}
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.FromLocationStation != nil {
				// 查询代理人员配置的代理站点
				ag, _ := ent.Database.Agent.QueryNotDeleted().
					WithEnterprise(func(query *ent.EnterpriseQuery) {
						query.WithStations()
					}).
					Where(
						agent.ID(atu.AgentID),
					).First(context.Background())
				if ag != nil && ag.Edges.Enterprise != nil {
					for _, v := range ag.Edges.Enterprise.Edges.Stations {
						if *item.FromLocationID == v.ID {
							out = true
							break
						}
					}
				}
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.FromLocationOperator != nil && atu.MaintainerID == *item.FromLocationID {
				out = true
			}
		default:
		}

	}

	if item.ToLocationType != 0 && item.ToLocationID != 0 {
		switch model.AssetLocationsType(item.ToLocationType) {
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.ToLocationWarehouse != nil {
				// if data, _ := ent.Database.Warehouse.QueryNotDeleted().
				// 	Where(
				// 		warehouse.ID(item.ToLocationID),
				// 		warehouse.HasBelongAssetManagersWith(assetmanager.ID(atu.AssetManagerID)),
				// 	).First(context.Background()); data != nil {
				// 	in = true
				// }

				// 根据上班地点确认inout
				if data, _ := ent.Database.AssetManager.QueryNotDeleted().Where(
					assetmanager.ID(atu.AssetManagerID),
					assetmanager.WarehouseID(item.ToLocationID),
				).First(context.Background()); data != nil {
					in = true
				}
			}
		case model.AssetLocationsTypeStore:
			if item.Edges.ToLocationStore != nil {
				// if data, _ := ent.Database.Store.QueryNotDeleted().
				// 	Where(
				// 		store.ID(item.ToLocationID),
				// 		store.HasEmployeesWith(employee.ID(atu.EmployeeID)),
				// 	).First(context.Background()); data != nil {
				// 	in = true
				// }

				// 根据上班地点确认inout
				if data, _ := ent.Database.Employee.QueryNotDeleted().Where(
					employee.ID(atu.EmployeeID),
					employee.DutyStoreID(item.ToLocationID),
				).First(context.Background()); data != nil {
					in = true
				}
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.ToLocationStation != nil {
				// 查询代理人员配置的代理站点
				ag, _ := ent.Database.Agent.QueryNotDeleted().
					WithEnterprise(func(query *ent.EnterpriseQuery) {
						query.WithStations()
					}).
					Where(
						agent.ID(atu.AgentID),
					).First(context.Background())
				if ag != nil && ag.Edges.Enterprise != nil {
					for _, v := range ag.Edges.Enterprise.Edges.Stations {
						if item.ToLocationID == v.ID {
							in = true
							break
						}
					}
				}
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.ToLocationOperator != nil && atu.MaintainerID == item.ToLocationID {
				in = true
			}
		default:
		}

	}

	switch {
	case in && out:
		return model.AssetTransferBoundTypeALl
	case !in && out:
		return model.AssetTransferBoundTypeOut
	case in && !out:
		return model.AssetTransferBoundTypeIn
	default:
		return ""

	}
}

// 筛选
func (s *assetTransferService) filter(ctx context.Context, q *ent.AssetTransferQuery, req *model.AssetTransferFilter) {
	// 查询调拨
	if req.FromLocationsType != nil && req.FromLocationsID != nil {
		q.Where(
			assettransfer.FromLocationType((*req.FromLocationsType).Value()),
			assettransfer.FromLocationID(*req.FromLocationsID),
		)
	}
	if req.ToLocationsType != nil && req.ToLocationsID != nil {
		q.Where(
			assettransfer.ToLocationType((*req.ToLocationsType).Value()),
			assettransfer.ToLocationID(*req.ToLocationsID),
		)
	}
	if req.Status != nil {
		q.Where(assettransfer.Status((*req.Status).Value()))
	}

	if req.OutStart != nil {
		q.Where(assettransfer.OutTimeAtGTE(tools.NewTime().ParseDateStringX(*req.OutStart)))
	}
	if req.OutEnd != nil {
		q.Where(assettransfer.OutTimeAtLTE(tools.NewTime().ParseNextDateStringX(*req.OutEnd)))
	}

	if req.Keyword != nil {
		q.Where(
			assettransfer.Or(
				assettransfer.SnContains(*req.Keyword),
				assettransfer.ReasonContains(*req.Keyword),
			),
		)
	}
	if req.AssetManagerID != 0 {
		wq := warehouse.HasBelongAssetManagersWith(assetmanager.ID(req.AssetManagerID))
		// 判断是否首页跳转查询
		if req.MainPage == nil {
			// 检查是否选择出入方
			switch {
			case req.FromLocationsID == nil && req.ToLocationsID == nil:
				q.Where(
					assettransfer.Or(
						assettransfer.HasFromLocationWarehouseWith(wq),
						assettransfer.HasToLocationWarehouseWith(wq),
					),
				)
			case req.FromLocationsID != nil && req.ToLocationsID == nil:
				q.Where(
					assettransfer.HasToLocationWarehouseWith(wq),
				)
			case req.FromLocationsID == nil && req.ToLocationsID != nil:
				q.Where(
					assettransfer.HasFromLocationWarehouseWith(wq),
				)
			default:

			}

		} else {
			// 检查是否选择出入方
			switch {
			case req.FromLocationsID == nil && req.ToLocationsID == nil:
				switch *req.MainPage {
				case model.AssetTransferMainPageReceive:
					q.Where(assettransfer.HasToLocationWarehouseWith(wq))
				case model.AssetTransferMainPageDeliver:
					q.Where(assettransfer.HasFromLocationWarehouseWith(wq))
				}

			case req.FromLocationsID != nil && req.ToLocationsID == nil:
				switch *req.MainPage {
				case model.AssetTransferMainPageReceive:
					q.Where(
						assettransfer.HasFromLocationWarehouseWith(warehouse.ID(*req.FromLocationsID)),
						assettransfer.HasToLocationWarehouseWith(wq),
					)
				case model.AssetTransferMainPageDeliver:
					q.Where(
						assettransfer.HasFromLocationWarehouseWith(warehouse.ID(*req.FromLocationsID)),
					)
				}

			case req.FromLocationsID == nil && req.ToLocationsID != nil:
				switch *req.MainPage {
				case model.AssetTransferMainPageReceive:
					q.Where(
						assettransfer.HasToLocationWarehouseWith(warehouse.ID(*req.ToLocationsID)),
					)
				case model.AssetTransferMainPageDeliver:
					q.Where(
						assettransfer.HasFromLocationWarehouseWith(wq),
						assettransfer.HasToLocationWarehouseWith(warehouse.ID(*req.ToLocationsID)),
					)
				}
			default:

			}
		}

	}
	if req.EmployeeID != 0 {
		wq := store.HasEmployeesWith(employee.ID(req.EmployeeID))

		// 判断是否首页跳转查询
		if req.MainPage == nil {
			// 检查是否选择出入方
			switch {
			case req.FromLocationsID == nil && req.ToLocationsID == nil:
				q.Where(
					assettransfer.Or(
						assettransfer.HasFromLocationStoreWith(wq),
						assettransfer.HasToLocationStoreWith(wq),
					),
				)
			default:
			}

		} else {
			// 检查是否选择出入方
			switch {
			case req.FromLocationsID == nil && req.ToLocationsID == nil:
				switch *req.MainPage {
				case model.AssetTransferMainPageReceive:
					q.Where(assettransfer.HasToLocationStoreWith(wq))
				case model.AssetTransferMainPageDeliver:
					q.Where(assettransfer.HasFromLocationStoreWith(wq))
				}

			case req.FromLocationsID != nil && req.ToLocationsID == nil:
				switch *req.MainPage {
				case model.AssetTransferMainPageReceive:
					q.Where(
						assettransfer.HasFromLocationStoreWith(store.ID(*req.FromLocationsID)),
						assettransfer.HasToLocationStoreWith(wq),
					)
				case model.AssetTransferMainPageDeliver:
					q.Where(
						assettransfer.HasFromLocationStoreWith(store.ID(*req.FromLocationsID)),
					)
				}

			case req.FromLocationsID == nil && req.ToLocationsID != nil:
				switch *req.MainPage {
				case model.AssetTransferMainPageReceive:
					q.Where(
						assettransfer.HasToLocationStoreWith(store.ID(*req.ToLocationsID)),
					)
				case model.AssetTransferMainPageDeliver:
					q.Where(
						assettransfer.HasFromLocationStoreWith(wq),
						assettransfer.HasToLocationStoreWith(store.ID(*req.ToLocationsID)),
					)
				}
			default:

			}
		}
	}
	if req.AgentID != 0 {
		// 查询代理人员配置的代理站点
		ids := make([]uint64, 0)
		ag, _ := ent.Database.Agent.QueryNotDeleted().
			WithEnterprise(func(query *ent.EnterpriseQuery) {
				query.WithStations()
			}).
			Where(
				agent.ID(req.AgentID),
			).First(context.Background())
		if ag != nil && ag.Edges.Enterprise != nil {
			for _, v := range ag.Edges.Enterprise.Edges.Stations {
				ids = append(ids, v.ID)
			}
		}
		q.Where(
			assettransfer.Or(
				assettransfer.HasFromLocationStationWith(enterprisestation.IDIn(ids...)),
				assettransfer.HasToLocationStationWith(enterprisestation.IDIn(ids...)),
			),
		)

	}
	if req.MaintainerID != 0 {
		q.Where(
			assettransfer.Or(
				assettransfer.HasFromLocationOperatorWith(maintainer.ID(req.MaintainerID)),
				assettransfer.HasToLocationOperatorWith(maintainer.ID(req.MaintainerID)),
			),
		)
	}

}

// TransferDetail 调拨详情
func (s *assetTransferService) TransferDetail(ctx context.Context, req *model.AssetTransferDetailReq) (res []*model.AssetTransferDetail, err error) {
	assetSnMap := make(map[string]*ent.Asset)
	ebikeSnAstMap := make(map[string]*model.TransferAssetDetail)
	sBSnAstMap := make(map[string]*model.TransferAssetDetail)
	nSbModelAstMap := make(map[string]*model.TransferAssetDetail)
	cabAccNameAstMap := make(map[string]*model.TransferAssetDetail)
	ebikeAccNameAstMap := make(map[string]*model.TransferAssetDetail)
	otherAccNameAstMap := make(map[string]*model.TransferAssetDetail)

	atds, err := ent.Database.AssetTransferDetails.QueryNotDeleted().
		WithInOperateAgent().WithInOperateManager().WithInOperateStore().WithInOperateMaintainer().WithInOperateCabinet().WithInOperateRider().
		Where(
			assettransferdetails.TransferID(req.ID),
			assettransferdetails.HasAssetWith(
				asset.DeletedAtIsNil(),
			),
		).WithAsset(func(query *ent.AssetQuery) {
		query.WithBrand().WithModel().WithMaterial().WithValues()
	}).All(ctx)
	if err != nil {
		return nil, err
	}

	// 分类统计调拨详情资产数据
	for _, atd := range atds {
		var inOpName, inTimeAt string
		inOpName = s.GetOperaterInfo(atd)
		if atd.InTimeAt != nil {
			inTimeAt = atd.InTimeAt.Format("2006-01-02 15:04:05")
		}

		ast := atd.Edges.Asset
		if ast != nil {
			switch ast.Type {
			case model.AssetTypeEbike.Value():
				// 电车类型资产需要以单辆车为单位进行统计
				NewAssetTransferDetail().TransferDetailCount(ebikeSnAstMap, ast.Sn, atd.IsIn, *ast.BrandID, inOpName, inTimeAt)
				assetSnMap[ast.Sn] = ast
			case model.AssetTypeSmartBattery.Value():
				// 智能电池需要以单独序列号为单位进行统计
				NewAssetTransferDetail().TransferDetailCount(sBSnAstMap, ast.Sn, atd.IsIn, *ast.ModelID, inOpName, inTimeAt)
				assetSnMap[ast.Sn] = ast
			case model.AssetTypeNonSmartBattery.Value():
				if ast.Edges.Model != nil {
					modelName := ast.Edges.Model.Model
					NewAssetTransferDetail().TransferDetailCount(nSbModelAstMap, modelName, atd.IsIn, *ast.ModelID, inOpName, inTimeAt)
				}
			case model.AssetTypeCabinetAccessory.Value():
				if ast.Edges.Material != nil {
					materialName := ast.Edges.Material.Name
					NewAssetTransferDetail().TransferDetailCount(cabAccNameAstMap, materialName, atd.IsIn, *ast.MaterialID, inOpName, inTimeAt)
				}
			case model.AssetTypeEbikeAccessory.Value():
				if ast.Edges.Material != nil {
					materialName := ast.Edges.Material.Name
					NewAssetTransferDetail().TransferDetailCount(ebikeAccNameAstMap, materialName, atd.IsIn, *ast.MaterialID, inOpName, inTimeAt)
				}
			case model.AssetTypeOtherAccessory.Value():
				if ast.Edges.Material != nil {
					materialName := ast.Edges.Material.Name
					NewAssetTransferDetail().TransferDetailCount(otherAccNameAstMap, materialName, atd.IsIn, *ast.MaterialID, inOpName, inTimeAt)
				}
			}

		}
	}

	// 分组组装调拨详情数据
	ebikeResList := make([]*model.AssetTransferDetail, 0)
	for k, v := range ebikeSnAstMap {
		item := assetSnMap[k]
		if item != nil && item.Edges.Brand != nil {
			// 查询属性值
			attributeValue, _ := item.QueryValues().WithAttribute().All(context.Background())
			assetAttributeMap := make(map[uint64]model.AssetAttribute)
			for _, v2 := range attributeValue {
				var attributeName, attributeKey string
				if v2.Edges.Attribute != nil {
					attributeName = v2.Edges.Attribute.Name
					attributeKey = v2.Edges.Attribute.Key
				}
				assetAttributeMap[v2.AttributeID] = model.AssetAttribute{
					AttributeID:      v2.AttributeID,
					AttributeValue:   v2.Value,
					AttributeName:    attributeName,
					AttributeKey:     attributeKey,
					AttributeValueID: v2.ID,
				}
			}
			ebikeResList = append(ebikeResList, &model.AssetTransferDetail{
				AssetType:     model.AssetTypeEbike,
				SN:            k,
				Name:          assetSnMap[k].Edges.Brand.Name,
				OutNum:        v.Outbound,
				InNum:         v.Inbound,
				InTimeAt:      v.InTimeAt,
				InOperateName: v.InOperateName,
				Attribute:     assetAttributeMap,
			})
		}

	}

	sBResList := make([]*model.AssetTransferDetail, 0)
	for k, v := range sBSnAstMap {
		if assetSnMap[k] != nil && assetSnMap[k].Edges.Model != nil {
			sBResList = append(sBResList, &model.AssetTransferDetail{
				AssetType:     model.AssetTypeSmartBattery,
				SN:            k,
				Name:          assetSnMap[k].Edges.Model.Model,
				OutNum:        v.Outbound,
				InNum:         v.Inbound,
				InTimeAt:      v.InTimeAt,
				InOperateName: v.InOperateName,
			})
		}
	}

	nSbResList := make([]*model.AssetTransferDetail, 0)
	for _, v := range nSbModelAstMap {
		nSbResList = append(nSbResList, &model.AssetTransferDetail{
			AssetType:     model.AssetTypeNonSmartBattery,
			Name:          v.Name,
			OutNum:        v.Outbound,
			InNum:         v.Inbound,
			ModelID:       v.ID,
			InTimeAt:      v.InTimeAt,
			InOperateName: v.InOperateName,
		})
	}

	cabAccResList := make([]*model.AssetTransferDetail, 0)
	for _, v := range cabAccNameAstMap {
		cabAccResList = append(cabAccResList, &model.AssetTransferDetail{
			AssetType:     model.AssetTypeCabinetAccessory,
			Name:          v.Name,
			OutNum:        v.Outbound,
			InNum:         v.Inbound,
			MaterialID:    v.ID,
			InTimeAt:      v.InTimeAt,
			InOperateName: v.InOperateName,
		})
	}

	ebikeAccResList := make([]*model.AssetTransferDetail, 0)
	for _, v := range ebikeAccNameAstMap {
		ebikeAccResList = append(ebikeAccResList, &model.AssetTransferDetail{
			AssetType:     model.AssetTypeEbikeAccessory,
			Name:          v.Name,
			OutNum:        v.Outbound,
			InNum:         v.Inbound,
			MaterialID:    v.ID,
			InTimeAt:      v.InTimeAt,
			InOperateName: v.InOperateName,
		})
	}

	otherAccResList := make([]*model.AssetTransferDetail, 0)
	for _, v := range otherAccNameAstMap {
		otherAccResList = append(otherAccResList, &model.AssetTransferDetail{
			AssetType:     model.AssetTypeOtherAccessory,
			Name:          v.Name,
			OutNum:        v.Outbound,
			InNum:         v.Inbound,
			MaterialID:    v.ID,
			InTimeAt:      v.InTimeAt,
			InOperateName: v.InOperateName,
		})
	}

	// 排序分组结果集
	sort.Slice(ebikeResList, func(i, j int) bool {
		return strings.Compare(ebikeResList[i].Name, ebikeResList[j].Name) < 0
	})
	sort.Slice(sBResList, func(i, j int) bool {
		return strings.Compare(sBResList[i].Name, sBResList[j].Name) < 0
	})
	sort.Slice(nSbResList, func(i, j int) bool {
		return strings.Compare(nSbResList[i].Name, nSbResList[j].Name) < 0
	})
	sort.Slice(cabAccResList, func(i, j int) bool {
		return strings.Compare(cabAccResList[i].Name, cabAccResList[j].Name) < 0
	})
	sort.Slice(ebikeAccResList, func(i, j int) bool {
		return strings.Compare(ebikeAccResList[i].Name, ebikeAccResList[j].Name) < 0
	})
	sort.Slice(otherAccResList, func(i, j int) bool {
		return strings.Compare(otherAccResList[i].Name, otherAccResList[j].Name) < 0
	})

	// 整合详情数据
	res = slices.Concat(ebikeResList, sBResList, nSbResList, cabAccResList, ebikeAccResList, otherAccResList)

	return res, nil
}

// GetOperaterInfo 获取操作人信息
func (s *assetTransferService) GetOperaterInfo(item *ent.AssetTransferDetails) string {
	switch model.OperatorType(item.InOperateType) {
	case model.OperatorTypeAssetManager:
		if item.Edges.InOperateManager != nil {
			return "[仓管]" + item.Edges.InOperateManager.Name
		}
	case model.OperatorTypeEmployee:
		if item.Edges.InOperateStore != nil {
			return "[门店]" + item.Edges.InOperateStore.Name
		}
	case model.OperatorTypeAgent:
		if item.Edges.InOperateAgent != nil {
			return "[代理]" + item.Edges.InOperateAgent.Name
		}
	case model.OperatorTypeMaintainer:
		if item.Edges.InOperateMaintainer != nil {
			return "[运维]" + item.Edges.InOperateMaintainer.Name
		}
	case model.OperatorTypeCabinet:
		if item.Edges.InOperateCabinet != nil {
			return "[电柜]" + item.Edges.InOperateCabinet.Name
		}
	case model.OperatorTypeRider:
		if item.Edges.InOperateRider != nil {
			return "[骑手]" + item.Edges.InOperateRider.Name
		}
	default:
	}

	return ""
}

// TransferCancel 取消调拨
func (s *assetTransferService) TransferCancel(ctx context.Context, req *model.AssetTransferDetailReq, modifier *model.Modifier) (err error) {
	item, err := ent.Database.AssetTransfer.QueryNotDeleted().WithTransferDetails(
		func(query *ent.AssetTransferDetailsQuery) {
			query.WithAsset()
		},
	).Where(assettransfer.ID(req.ID)).First(ctx)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("调拨单不存在")
	}
	if item.Status == model.AssetTransferStatusStock.Value() {
		return errors.New("已入库的调拨单不能取消")
	}
	if item.Status == model.AssetTransferStatusCancel.Value() {
		return errors.New("调拨单已取消")
	}

	// 初始调拨删除资产
	// if item.Type == model.AssetTransferTypeInitial.Value() {
	// 	if item.Edges.TransferDetails != nil {
	// 		for _, v := range item.Edges.TransferDetails {
	// 			_ = NewAsset().Delete(ctx, v.AssetID)
	// 		}
	// 	}
	// }

	// 修改调拨单状态
	_, err = item.Update().
		SetStatus(model.AssetTransferStatusCancel.Value()).
		SetLastModifier(modifier).
		Save(ctx)
	if err != nil {
		return err
	}
	// 修改资产状态
	ids := make([]uint64, 0, len(item.Edges.TransferDetails))
	for _, v := range item.Edges.TransferDetails {
		ids = append(ids, v.AssetID)
	}
	_, err = ent.Database.Asset.Update().Where(asset.IDIn(ids...)).SetStatus(model.AssetStatusStock.Value()).SetLastModifier(modifier).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// TransferReceive 接收资产
func (s *assetTransferService) TransferReceive(ctx context.Context, req *model.AssetTransferReceiveBatchReq, modifier *model.Modifier) (err error) {
	timeNow := time.Now()
	for _, v := range req.AssetTransferReceive {
		// 查询调拨单
		item, _ := ent.Database.AssetTransfer.QueryNotDeleted().
			Where(
				assettransfer.ID(v.ID),
				assettransfer.StatusNEQ(model.AssetTransferStatusCancel.Value()),
			).First(ctx)
		if item == nil {
			return errors.New("调拨单不存在,或已取消")
		}

		iDs := make([]uint64, 0)
		assetIDs := make([]uint64, 0)
		for _, vl := range v.Detail {
			switch vl.AssetType {
			case model.AssetTypeSmartBattery, model.AssetTypeEbike:
				// 此类资产可以分批次接收
				receiveiDs, receiveAssetIDs, err := s.receiveAssetWithSN(ctx, vl, v.ID)
				if err != nil {
					return err
				}
				iDs = append(iDs, receiveiDs...)
				assetIDs = append(assetIDs, receiveAssetIDs...)
			case model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeNonSmartBattery, model.AssetTypeOtherAccessory:
				// 此类资产只能一次性接收
				if item.Status == model.AssetTransferStatusStock.Value() {
					return errors.New("已入库的调拨单不能接收")
				}
				receiveiDs, receiveAssetIDs, err := s.receiveAssetWithoutSN(ctx, vl, v.ID)
				if err != nil {
					return err
				}
				iDs = append(iDs, receiveiDs...)
				assetIDs = append(assetIDs, receiveAssetIDs...)
			}
		}
		// 修改调拨详情状态
		var remark string
		if v.Remark != nil {
			remark = *v.Remark
		}
		err = ent.Database.AssetTransferDetails.Update().Where(assettransferdetails.IDIn(iDs...)).
			SetIsIn(true).
			SetInTimeAt(timeNow).
			SetInOperateID(modifier.ID).
			SetInOperateType(req.OperateType.Value()).
			SetLastModifier(modifier).
			SetRemark(remark).
			Exec(ctx)
		if err != nil {
			return err
		}
		// 修改调拨单状态
		err = ent.Database.AssetTransfer.Update().
			Where(assettransfer.ID(v.ID)).
			SetStatus(model.AssetTransferStatusStock.Value()).
			SetLastModifier(modifier).
			SetInNum(uint(len(iDs))).Exec(ctx)
		if err != nil {
			return err
		}
		// 修改资产状态
		err = ent.Database.Asset.Update().
			Where(asset.IDIn(assetIDs...)).
			SetLocationsType(item.ToLocationType).
			SetLocationsID(item.ToLocationID).
			SetStatus(model.AssetStatusStock.Value()).
			SetLastModifier(modifier).Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// 有编号资产接收
func (s *assetTransferService) receiveAssetWithSN(ctx context.Context, req model.AssetTransferReceiveDetail, assetTransferID uint64) (iDs []uint64, assetIDs []uint64, err error) {
	if req.SN == nil || *req.SN == "" {
		return nil, nil, errors.New("资产编号不能为空")
	}
	item, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().WithAsset().
		Where(
			assettransferdetails.TransferID(assetTransferID),
			assettransferdetails.IsIn(false),
			assettransferdetails.HasAssetWith(
				asset.Status(model.AssetStatusDelivering.Value()),
				asset.Sn(*req.SN),
			),
			assettransferdetails.HasTransferWith(
				assettransfer.StatusNEQ(model.AssetTransferStatusCancel.Value()),
			)).First(ctx)
	if item == nil {
		return nil, nil, errors.New(*req.SN + "物资不存在或不处于配送中")
	}
	iDs = append(iDs, item.ID)
	assetIDs = append(assetIDs, item.Edges.Asset.ID)
	return
}

// receiveAssetWithoutSN 无编号资产接收
func (s *assetTransferService) receiveAssetWithoutSN(ctx context.Context, req model.AssetTransferReceiveDetail, assetTransferID uint64) (iDs []uint64, assetIDs []uint64, err error) {
	if req.Num == nil {
		return nil, nil, errors.New("接收数量不能为空")
	}

	q := ent.Database.AssetTransferDetails.QueryNotDeleted().WithAsset().
		Where(
			assettransferdetails.IsIn(false),
			assettransferdetails.TransferID(assetTransferID),
			assettransferdetails.HasTransferWith(
				assettransfer.StatusNEQ(model.AssetTransferStatusCancel.Value()),
			),
		).Limit(int(*req.Num))

	switch {
	case req.AssetType == model.AssetTypeNonSmartBattery:
		if req.ModelID == nil {
			return nil, nil, errors.New("电池型号ID不能为空")
		}
		q.Where(
			assettransferdetails.HasAssetWith(
				asset.Status(model.AssetStatusDelivering.Value()),
				asset.ModelID(*req.ModelID),
			),
		)
	case req.AssetType == model.AssetTypeCabinetAccessory || req.AssetType == model.AssetTypeEbikeAccessory || req.AssetType == model.AssetTypeOtherAccessory:
		if req.MaterialID == nil {
			return nil, nil, errors.New("物资分类ID不能为空")
		}
		q.Where(
			assettransferdetails.HasAssetWith(
				asset.Status(model.AssetStatusDelivering.Value()),
				asset.MaterialID(*req.MaterialID),
			),
		)
	default:
		return nil, nil, errors.New("当前物资类型不支持")
	}

	list, _ := q.All(ctx)
	if len(list) == 0 {
		return nil, nil, errors.New("物资不存在或已入库")
	}
	for _, v := range list {
		iDs = append(iDs, v.ID)
		if v.Edges.Asset != nil {
			assetIDs = append(assetIDs, v.Edges.Asset.ID)
		}
	}
	return iDs, assetIDs, nil
}

// GetTransferBySN 根据sn查询未入库的调拨单
func (s *assetTransferService) GetTransferBySN(assetSignInfo definition.AssetSignInfo, ctx context.Context, req *model.GetTransferBySNReq) (res *model.AssetTransferListRes, err error) {
	item, _ := ent.Database.AssetTransferDetails.Query().Where(
		assettransferdetails.HasAssetWith(
			asset.Sn(req.SN),
			asset.Status(model.AssetStatusDelivering.Value()),
		),
		assettransferdetails.IsIn(false),
		assettransferdetails.HasTransferWith(
			assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
		),
	).WithTransfer(func(query *ent.AssetTransferQuery) {
		query.
			WithFromLocationOperator().WithFromLocationStation().WithFromLocationStore().WithFromLocationWarehouse().
			WithToLocationOperator().WithToLocationStation().WithToLocationStore().WithToLocationWarehouse()
	}).WithAsset().First(ctx)
	if item == nil {
		return nil, errors.New("物资不存在或未调拨")
	}

	if item.Edges.Transfer != nil {
		// 检查库管上班权限
		switch {
		case assetSignInfo.AssetManager != nil:
			if am, _ := ent.Database.AssetManager.QueryNotDeleted().Where(
				assetmanager.ID(assetSignInfo.AssetManager.ID),
				assetmanager.WarehouseID(item.Edges.Transfer.ToLocationID),
			).First(ctx); am == nil {
				return nil, errors.New("该入库单不属于当前上班位置")
			}
		case assetSignInfo.Employee != nil:
			if am, _ := ent.Database.Employee.QueryNotDeleted().Where(
				employee.ID(assetSignInfo.Employee.ID),
				employee.DutyStoreID(item.Edges.Transfer.ToLocationID),
			).First(ctx); am == nil {
				return nil, errors.New("该入库单不属于当前上班位置")
			}
		}
	}

	res = &model.AssetTransferListRes{
		ID:                item.Edges.Transfer.ID,
		SN:                item.Edges.Transfer.Sn,
		Reason:            item.Edges.Transfer.Reason,
		Status:            model.AssetTransferStatus(item.Edges.Transfer.Status).String(),
		AssetTransferType: model.AssetTransferType(item.Edges.Transfer.Type),
		OutNum:            item.Edges.Transfer.OutNum,
		InNum:             item.Edges.Transfer.InNum,
	}

	var fromLocationName, toLocationName string
	if item.Edges.Transfer.FromLocationType != nil && item.Edges.Transfer.FromLocationID != nil {
		res.FromLocationType = *item.Edges.Transfer.FromLocationType
		res.FromLocationID = *item.Edges.Transfer.FromLocationID
		switch model.AssetLocationsType(*item.Edges.Transfer.FromLocationType) {
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.Transfer.Edges.FromLocationWarehouse != nil {
				fromLocationName = "[仓库]" + item.Edges.Transfer.Edges.FromLocationWarehouse.Name
			}
		case model.AssetLocationsTypeStore:
			if item.Edges.Transfer.Edges.FromLocationStore != nil {
				fromLocationName = "[门店]" + item.Edges.Transfer.Edges.FromLocationStore.Name
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.Transfer.Edges.FromLocationStation != nil {
				fromLocationName = "[站点]" + item.Edges.Transfer.Edges.FromLocationStation.Name
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.Transfer.Edges.FromLocationOperator != nil {
				fromLocationName = "[运维]" + item.Edges.Transfer.Edges.FromLocationOperator.Name
			}
		default:
		}
	}
	if item.Edges.Transfer.ToLocationType != 0 && item.Edges.Transfer.ToLocationID != 0 {
		res.ToLocationID = item.Edges.Transfer.ToLocationID
		res.ToLocationType = item.Edges.Transfer.ToLocationType
		switch model.AssetLocationsType(item.Edges.Transfer.ToLocationType) {
		case model.AssetLocationsTypeWarehouse:
			if item.Edges.Transfer.Edges.ToLocationWarehouse != nil {
				toLocationName = "[仓库]" + item.Edges.Transfer.Edges.ToLocationWarehouse.Name
			}
		case model.AssetLocationsTypeStore:
			if item.Edges.Transfer.Edges.ToLocationStore != nil {
				toLocationName = "[门店]" + item.Edges.Transfer.Edges.ToLocationStore.Name
			}
		case model.AssetLocationsTypeStation:
			if item.Edges.Transfer.Edges.ToLocationStation != nil {
				toLocationName = "[站点]" + item.Edges.Transfer.Edges.ToLocationStation.Name
			}
		case model.AssetLocationsTypeOperation:
			if item.Edges.Transfer.Edges.ToLocationOperator != nil {
				toLocationName = "[运维]" + item.Edges.Transfer.Edges.ToLocationOperator.Name
			}
		default:
		}
	}
	if item.Edges.Transfer.OutTimeAt != nil {
		res.OutTimeAt = item.Edges.Transfer.OutTimeAt.Format("2006-01-02 15:04:05")
	}
	res.ToLocationName = toLocationName
	res.FromLocationName = fromLocationName

	if item.Edges.Asset != nil {
		res.AssetDetail = &model.AssetTransferDetail{
			AssetType: model.AssetType(item.Edges.Asset.Type),
			SN:        item.Edges.Asset.Sn,
			Name:      item.Edges.Asset.Name,
			Attribute: make(map[uint64]model.AssetAttribute),
		}
		// 查询属性值
		attributeValue, _ := item.Edges.Asset.QueryValues().WithAttribute().All(context.Background())
		for _, v := range attributeValue {
			var attributeName, attributeKey string
			if v.Edges.Attribute != nil {
				attributeName = v.Edges.Attribute.Name
				attributeKey = v.Edges.Attribute.Key
			}
			res.AssetDetail.Attribute[v.AttributeID] = model.AssetAttribute{
				AttributeID:      v.AttributeID,
				AttributeValue:   v.Value,
				AttributeName:    attributeName,
				AttributeKey:     attributeKey,
				AttributeValueID: v.ID,
			}
		}
	}

	return res, nil
}

// Flow 电池流转明细
func (s *assetTransferService) Flow(ctx context.Context, req *model.AssetTransferFlowReq) []*model.AssetTransferFlow {
	q := ent.Database.AssetTransferDetails.QueryNotDeleted().
		Where(
			assettransferdetails.IsIn(true),
			assettransferdetails.HasTransferWith(
				assettransfer.StatusNEQ(model.AssetTransferStatusCancel.Value()),
			),
			assettransferdetails.HasAssetWith(
				asset.Sn(req.SN),
			),
		).WithTransfer(func(query *ent.AssetTransferQuery) {
		query.
			WithFromLocationOperator().WithFromLocationStation().WithFromLocationStore().WithFromLocationWarehouse().WithFromLocationCabinet().WithFromLocationRider().
			WithToLocationOperator().WithToLocationStation().WithToLocationStore().WithToLocationWarehouse().WithToLocationCabinet().WithToLocationRider()
	}).WithInOperateAgent().WithInOperateManager().WithInOperateStore().WithInOperateMaintainer().WithInOperateCabinet().WithInOperateRider()
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			assettransferdetails.InTimeAtGTE(start),
			assettransferdetails.InTimeAtLTE(end),
			assettransferdetails.Or(
				assettransferdetails.HasTransferWith(
					assettransfer.CreatedAtGTE(start),
					assettransfer.CreatedAtLTE(end),
				),
			),
		)
	}
	if req.AssetType != nil {
		q.Where(assettransferdetails.HasAssetWith(asset.Type(req.AssetType.Value())))
	}
	all, _ := q.Order(ent.Desc(assettransferdetails.FieldCreatedAt)).All(ctx)
	res := make([]*model.AssetTransferFlow, 0)
	for _, item := range all {
		var fromoperateName, toOperateName, fromLocationName, toLocationName string
		var transferType model.AssetTransferType
		var fromLocationType, toLocationType model.AssetLocationsType
		// 入库操作人
		switch model.OperatorType(item.InOperateType) {
		case model.OperatorTypeAssetManager:
			if item.Edges.InOperateManager != nil {
				if r, _ := item.Edges.InOperateManager.QueryRole().First(ctx); r != nil {
					toOperateName = "[" + r.Name + "]" + item.Edges.InOperateManager.Name
				}
			}
		case model.OperatorTypeEmployee:
			if item.Edges.InOperateStore != nil {
				toOperateName = "[门店]" + item.Edges.InOperateStore.Name
			}
		case model.OperatorTypeAgent:
			if item.Edges.InOperateAgent != nil {
				toOperateName = "[代理]" + item.Edges.InOperateAgent.Name
			}
		case model.OperatorTypeMaintainer:
			if item.Edges.InOperateMaintainer != nil {
				toOperateName = "[运维]" + item.Edges.InOperateMaintainer.Name
			}
		case model.OperatorTypeCabinet:
			if item.Edges.InOperateCabinet != nil {
				toOperateName = "[电柜]" + item.Edges.InOperateCabinet.Name
			}
		case model.OperatorTypeRider:
			if item.Edges.InOperateRider != nil {
				toOperateName = "[骑手]" + item.Edges.InOperateRider.Name
			}
		default:
		}

		if item.Edges.Transfer != nil {
			if item.Edges.Transfer.FromLocationType != nil && item.Edges.Transfer.FromLocationID != nil {
				fromLocationType = model.AssetLocationsType(*item.Edges.Transfer.FromLocationType)
				switch model.AssetLocationsType(*item.Edges.Transfer.FromLocationType) {
				case model.AssetLocationsTypeWarehouse:
					if item.Edges.Transfer.Edges.FromLocationWarehouse != nil {
						fromLocationName = "[仓库]" + item.Edges.Transfer.Edges.FromLocationWarehouse.Name
					}
				case model.AssetLocationsTypeStore:
					if item.Edges.Transfer.Edges.FromLocationStore != nil {
						fromLocationName = "[门店]" + item.Edges.Transfer.Edges.FromLocationStore.Name
					}
				case model.AssetLocationsTypeStation:
					if item.Edges.Transfer.Edges.FromLocationStation != nil {
						fromLocationName = "[站点]" + item.Edges.Transfer.Edges.FromLocationStation.Name
					}
				case model.AssetLocationsTypeOperation:
					if item.Edges.Transfer.Edges.FromLocationOperator != nil {
						fromLocationName = "[运维]" + item.Edges.Transfer.Edges.FromLocationOperator.Name
					}
				case model.AssetLocationsTypeCabinet:
					if item.Edges.Transfer.Edges.FromLocationCabinet != nil {
						fromLocationName = "[电柜]" + item.Edges.Transfer.Edges.FromLocationCabinet.Name
					}
				case model.AssetLocationsTypeRider:
					if item.Edges.Transfer.Edges.FromLocationRider != nil {
						fromLocationName = "[骑手]" + item.Edges.Transfer.Edges.FromLocationRider.Name
					}
				default:
				}
			}
			// 入库位置
			if item.Edges.Transfer.ToLocationType != 0 && item.Edges.Transfer.ToLocationID != 0 {
				toLocationType = model.AssetLocationsType(item.Edges.Transfer.ToLocationType)
				switch model.AssetLocationsType(item.Edges.Transfer.ToLocationType) {
				case model.AssetLocationsTypeWarehouse:
					if item.Edges.Transfer.Edges.ToLocationWarehouse != nil {
						toLocationName = "[仓库]" + item.Edges.Transfer.Edges.ToLocationWarehouse.Name
					}
				case model.AssetLocationsTypeStore:
					if item.Edges.Transfer.Edges.ToLocationStore != nil {
						toLocationName = "[门店]" + item.Edges.Transfer.Edges.ToLocationStore.Name
					}
				case model.AssetLocationsTypeStation:
					if item.Edges.Transfer.Edges.ToLocationStation != nil {
						toLocationName = "[站点]" + item.Edges.Transfer.Edges.ToLocationStation.Name
					}
				case model.AssetLocationsTypeOperation:
					if item.Edges.Transfer.Edges.ToLocationOperator != nil {
						toLocationName = "[运维]" + item.Edges.Transfer.Edges.ToLocationOperator.Name
					}
				case model.AssetLocationsTypeCabinet:
					if item.Edges.Transfer.Edges.ToLocationCabinet != nil {
						toLocationName = "[电柜]" + item.Edges.Transfer.Edges.ToLocationCabinet.Name
					}
				case model.AssetLocationsTypeRider:
					if item.Edges.Transfer.Edges.ToLocationRider != nil {
						toLocationName = "[骑手]" + item.Edges.Transfer.Edges.ToLocationRider.Name
					}
				default:
				}
			}

			// 出库操作人
			if item.Edges.Transfer.OutOperateType != nil && item.Edges.Transfer.OutOperateID != nil {
				switch model.OperatorType(*item.Edges.Transfer.OutOperateType) {
				case model.OperatorTypeAssetManager:
					if item.Edges.Transfer.Edges.OutOperateManager != nil {
						if r, _ := item.Edges.Transfer.Edges.OutOperateManager.QueryRole().First(ctx); r != nil {
							fromoperateName = "[" + r.Name + "]" + item.Edges.Transfer.Edges.OutOperateManager.Name
						}
					}
				case model.OperatorTypeEmployee:
					if item.Edges.Transfer.Edges.OutOperateStore != nil {
						fromoperateName = "[门店]" + item.Edges.Transfer.Edges.OutOperateStore.Name
					}
				case model.OperatorTypeAgent:
					if item.Edges.Transfer.Edges.OutOperateAgent != nil {
						fromoperateName = "[代理]" + item.Edges.Transfer.Edges.OutOperateAgent.Name
					}
				case model.OperatorTypeMaintainer:
					if item.Edges.Transfer.Edges.OutOperateMaintainer != nil {
						fromoperateName = "[运维]" + item.Edges.Transfer.Edges.OutOperateMaintainer.Name
					}
				case model.OperatorTypeCabinet:
					if item.Edges.Transfer.Edges.OutOperateCabinet != nil {
						fromoperateName = "[电柜]" + item.Edges.Transfer.Edges.OutOperateCabinet.Name
					}
				case model.OperatorTypeRider:
					if item.Edges.Transfer.Edges.OutOperateRider != nil {
						fromoperateName = "[骑手]" + item.Edges.Transfer.Edges.OutOperateRider.Name
					}
				default:
				}
			}
			// 调拨类型
			transferType = model.AssetTransferType(item.Edges.Transfer.Type)
		}
		var outTimeAt, inTimeAt string
		if item.Edges.Transfer.OutTimeAt != nil {
			outTimeAt = item.Edges.Transfer.OutTimeAt.Format("2006-01-02 15:04:05")
		}
		if item.InTimeAt != nil {
			inTimeAt = item.InTimeAt.Format("2006-01-02 15:04:05")
		}

		var out *model.AssetTransferFlowDetail

		if fromLocationName != "" {
			out = &model.AssetTransferFlowDetail{
				OperatorName:     fromoperateName,
				LocationsName:    fromLocationName,
				TimeAt:           outTimeAt,
				TransferTypeName: "[出库]" + transferType.String(),
				TransferType:     transferType.Value(),
				LocationsType:    fromLocationType.Value(),
			}
		}

		in := model.AssetTransferFlowDetail{
			OperatorName:     toOperateName,
			LocationsName:    toLocationName,
			TimeAt:           inTimeAt,
			TransferTypeName: "[入库]" + transferType.String(),
			TransferType:     transferType.Value(),
			LocationsType:    toLocationType.Value(),
		}

		res = append(res, &model.AssetTransferFlow{
			Out: out,
			In:  &in,
		})
	}
	return res
}

// TransferDetailsList 出入库明细列表
func (s *assetTransferService) TransferDetailsList(ctx context.Context, req *model.AssetTransferDetailListReq) (res *model.PaginationRes, err error) {
	q := ent.Database.AssetTransfer.QueryNotDeleted().
		Where(
			assettransfer.StatusNEQ(model.AssetTransferStatusCancel.Value()),
		).WithTransferDetails(func(query *ent.AssetTransferDetailsQuery) {
		query.WithInOperateAgent().WithInOperateManager().WithInOperateStore().WithInOperateMaintainer().WithInOperateCabinet().WithInOperateRider().
			WithAsset(func(query *ent.AssetQuery) {
				query.WithMaterial().WithCity().WithModel().WithBrand()
			})
	}).WithFromLocationOperator().WithFromLocationStation().WithFromLocationStore().WithFromLocationWarehouse().WithFromLocationCabinet().WithFromLocationRider().
		WithToLocationOperator().WithToLocationStation().WithToLocationStore().WithToLocationWarehouse().WithToLocationCabinet().WithToLocationRider().
		WithOutOperateAgent().WithOutOperateManager().WithOutOperateStore().WithOutOperateMaintainer().WithOutOperateCabinet().WithOutOperateRider().
		Order(ent.Desc(asset.FieldCreatedAt))

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			assettransfer.CreatedAtGTE(start),
			assettransfer.CreatedAtLTE(end),
			assettransfer.Or(
				assettransfer.HasTransferDetailsWith(
					assettransferdetails.InTimeAtGTE(start),
					assettransferdetails.InTimeAtLTE(end),
				),
			),
		)
	}
	if req.AssetType != nil {
		q.Where(assettransfer.HasTransferDetailsWith(assettransferdetails.HasAssetWith(asset.Type(req.AssetType.Value()))))
	}
	if req.CityID != nil {
		q.Where(assettransfer.HasTransferDetailsWith(assettransferdetails.HasAssetWith(asset.CityID(*req.CityID))))
	}
	if req.SN != nil {
		q.Where(assettransfer.HasTransferDetailsWith(assettransferdetails.HasAssetWith(asset.Sn(*req.SN))))
	}
	if req.LocationsID != nil && req.LocationsType != nil {
		q.Where(
			func(selector *sql.Selector) {
				switch *req.LocationsType {
				case model.AssetLocationsTypeRider, model.AssetLocationsTypeCabinet:
					if req.LocationsKeyword != nil && *req.LocationsType == model.AssetLocationsTypeCabinet {
						q.Where(
							assettransfer.Or(
								assettransfer.HasFromLocationCabinetWith(cabinet.SnContains(*req.LocationsKeyword)),
								assettransfer.HasToLocationCabinetWith(cabinet.SnContains(*req.LocationsKeyword)),
							),
						)
					}
					if req.LocationsKeyword != nil && *req.LocationsType == model.AssetLocationsTypeRider {
						q.Where(
							assettransfer.Or(
								assettransfer.HasFromLocationRiderWith(rider.NameContains(*req.LocationsKeyword)),
								assettransfer.HasToLocationRiderWith(rider.NameContains(*req.LocationsKeyword)),
							),
						)
					}
				case model.AssetLocationsTypeWarehouse, model.AssetLocationsTypeStore, model.AssetLocationsTypeStation, model.AssetLocationsTypeOperation:
					selector.Where(
						sql.Or(
							sql.And(
								sql.EQ(assettransfer.FieldFromLocationID, *req.LocationsID),
								sql.EQ(assettransfer.FieldFromLocationType, req.LocationsType.Value()),
							),
							sql.And(
								sql.EQ(assettransfer.FieldToLocationID, *req.LocationsID),
								sql.EQ(assettransfer.FieldToLocationType, req.LocationsType.Value()),
							),
						),
					)
				}
			},
		)
	}
	if req.AssetTransferType != nil {
		q.Where(assettransfer.Type(req.AssetTransferType.Value()))
	}
	if req.CabinetSN != nil {
		q.Where(
			assettransfer.Or(
				assettransfer.HasTransferDetailsWith(assettransferdetails.HasInOperateCabinetWith(cabinet.Sn(*req.CabinetSN))),
				assettransfer.HasOutOperateCabinetWith(cabinet.Sn(*req.CabinetSN)),
			),
		)
	}
	if req.AssetManagerID != 0 {
		// 查询库管人员配置的仓库数据
		wIds := make([]uint64, 0)
		am, _ := ent.Database.AssetManager.QueryNotDeleted().WithBelongWarehouses().
			Where(
				assetmanager.ID(req.AssetManagerID),
				assetmanager.HasBelongWarehousesWith(warehouse.DeletedAtIsNil()),
			).First(context.Background())
		if am != nil {
			for _, wh := range am.Edges.BelongWarehouses {
				wIds = append(wIds, wh.ID)
			}
		}
		q.Where(
			assettransfer.Or(
				assettransfer.HasFromLocationWarehouseWith(warehouse.IDIn(wIds...)),
				assettransfer.HasToLocationWarehouseWith(warehouse.IDIn(wIds...)),
			),
		)
	}
	if req.EmployeeID != 0 {
		// 查询门店人员配置的门店数据
		sIds := make([]uint64, 0)
		ep, _ := ent.Database.Employee.QueryNotDeleted().WithStores().
			Where(
				employee.ID(req.EmployeeID),
				employee.HasStoresWith(store.DeletedAtIsNil()),
			).First(context.Background())
		if ep != nil {
			for _, st := range ep.Edges.Stores {
				sIds = append(sIds, st.ID)
			}
		}
		q.Where(
			assettransfer.Or(
				assettransfer.HasFromLocationStoreWith(store.IDIn(sIds...)),
				assettransfer.HasToLocationStoreWith(store.IDIn(sIds...)),
			),
		)
	}
	if req.AgentID != 0 {
		// 查询代理人员配置的代理站点
		ids := make([]uint64, 0)
		ag, _ := ent.Database.Agent.QueryNotDeleted().
			WithEnterprise(func(query *ent.EnterpriseQuery) {
				query.WithStations()
			}).
			Where(
				agent.ID(req.AgentID),
			).First(context.Background())
		if ag != nil && ag.Edges.Enterprise != nil {
			for _, v := range ag.Edges.Enterprise.Edges.Stations {
				ids = append(ids, v.ID)
			}
		}
		q.Where(
			assettransfer.Or(
				assettransfer.HasFromLocationStationWith(enterprisestation.IDIn(ids...)),
				assettransfer.HasToLocationStationWith(enterprisestation.IDIn(ids...)),
			),
		)

	}
	if req.MaintainerID != 0 {
		q.Where(
			assettransfer.Or(
				assettransfer.HasFromLocationOperatorWith(maintainer.ID(req.MaintainerID)),
				assettransfer.HasToLocationOperatorWith(maintainer.ID(req.MaintainerID)),
			),
		)
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetTransfer) (res *model.AssetTransferDetailListRes) {
		var fromOperateName, toOperateName, fromLocationName, toLocationName, outTimeAt, inTimeAt, cityName, assetName, inRemark string
		if item.OutTimeAt != nil {
			outTimeAt = item.OutTimeAt.Format("2006-01-02 15:04:05")
		}
		// 入库操作人
		if item.Edges.TransferDetails != nil {
			// 其它物资类型有很多记录 所以只取第一条
			details := item.Edges.TransferDetails[0]
			if details != nil {
				if details.InTimeAt != nil {
					inTimeAt = details.InTimeAt.Format("2006-01-02 15:04:05")
				}
				switch model.OperatorType(details.InOperateType) {
				case model.OperatorTypeAssetManager:
					if details.Edges.InOperateManager != nil {
						if r, _ := details.Edges.InOperateManager.QueryRole().First(ctx); r != nil {
							toOperateName = "[" + r.Name + "]" + details.Edges.InOperateManager.Name
						}
					}
				case model.OperatorTypeEmployee:
					if details.Edges.InOperateStore != nil {
						toOperateName = "[门店]" + details.Edges.InOperateStore.Name
					}
				case model.OperatorTypeAgent:
					if details.Edges.InOperateAgent != nil {
						toOperateName = "[代理]" + details.Edges.InOperateAgent.Name
					}
				case model.OperatorTypeMaintainer:
					if details.Edges.InOperateMaintainer != nil {
						toOperateName = "[运维]" + details.Edges.InOperateMaintainer.Name
					}
				case model.OperatorTypeCabinet:
					if details.Edges.InOperateCabinet != nil {
						toOperateName = "[电柜]" + details.Edges.InOperateCabinet.Name
					}
				case model.OperatorTypeRider:
					if details.Edges.InOperateRider != nil {
						toOperateName = "[骑手]" + details.Edges.InOperateRider.Name
					}
				default:
				}
				if details.Edges.Asset != nil {
					if details.Edges.Asset.Edges.City != nil {
						cityName = details.Edges.Asset.Edges.City.Name
					}
					switch model.AssetType(details.Edges.Asset.Type) {
					case model.AssetTypeSmartBattery:
						if details.Edges.Asset.Edges.Model != nil {
							assetName = "[" + details.Edges.Asset.Edges.Model.Model + "]" + details.Edges.Asset.Sn
						}
					case model.AssetTypeEbike:
						if details.Edges.Asset.Edges.Brand != nil {
							assetName = "[" + details.Edges.Asset.Edges.Brand.Name + "]" + details.Edges.Asset.Sn
						}
					case model.AssetTypeNonSmartBattery, model.AssetTypeCabinetAccessory, model.AssetTypeEbikeAccessory, model.AssetTypeOtherAccessory:
						if details.Edges.Asset.Edges.Material != nil {
							assetName = details.Edges.Asset.Edges.Material.Name
						}
					}
				}
				inRemark = details.Remark
			}
		}
		if item.FromLocationType != nil && item.FromLocationID != nil {
			switch model.AssetLocationsType(*item.FromLocationType) {
			case model.AssetLocationsTypeWarehouse:
				if item.Edges.FromLocationWarehouse != nil {
					fromLocationName = "[仓库]" + item.Edges.FromLocationWarehouse.Name
				}
			case model.AssetLocationsTypeStore:
				if item.Edges.FromLocationStore != nil {
					fromLocationName = "[门店]" + item.Edges.FromLocationStore.Name
				}
			case model.AssetLocationsTypeStation:
				if item.Edges.FromLocationStation != nil {
					fromLocationName = "[站点]" + item.Edges.FromLocationStation.Name
				}
			case model.AssetLocationsTypeOperation:
				if item.Edges.FromLocationOperator != nil {
					fromLocationName = "[运维]" + item.Edges.FromLocationOperator.Name
				}
			case model.AssetLocationsTypeCabinet:
				if item.Edges.FromLocationCabinet != nil {
					fromLocationName = "[电柜]" + item.Edges.FromLocationCabinet.Name
				}
			case model.AssetLocationsTypeRider:
				if item.Edges.FromLocationRider != nil {
					fromLocationName = "[骑手]" + item.Edges.FromLocationRider.Name
				}
			default:
			}
		}
		// 入库位置
		if item.ToLocationType != 0 && item.ToLocationID != 0 {
			switch model.AssetLocationsType(item.ToLocationType) {
			case model.AssetLocationsTypeWarehouse:
				if item.Edges.ToLocationWarehouse != nil {
					toLocationName = "[仓库]" + item.Edges.ToLocationWarehouse.Name
				}
			case model.AssetLocationsTypeStore:
				if item.Edges.ToLocationStore != nil {
					toLocationName = "[门店]" + item.Edges.ToLocationStore.Name
				}
			case model.AssetLocationsTypeStation:
				if item.Edges.ToLocationStation != nil {
					toLocationName = "[站点]" + item.Edges.ToLocationStation.Name
				}
			case model.AssetLocationsTypeOperation:
				if item.Edges.ToLocationOperator != nil {
					toLocationName = "[运维]" + item.Edges.ToLocationOperator.Name
				}
			case model.AssetLocationsTypeCabinet:
				if item.Edges.ToLocationCabinet != nil {
					toLocationName = "[电柜]" + item.Edges.ToLocationCabinet.Name
				}
			case model.AssetLocationsTypeRider:
				if item.Edges.ToLocationRider != nil {
					toLocationName = "[骑手]" + item.Edges.ToLocationRider.Name
				}
			default:
			}
		}
		// 出库操作人
		if item.OutOperateType != nil && item.OutOperateID != nil {
			switch model.OperatorType(*item.OutOperateType) {
			case model.OperatorTypeAssetManager:
				if item.Edges.OutOperateManager != nil {
					if r, _ := item.Edges.OutOperateManager.QueryRole().First(ctx); r != nil {
						fromOperateName = "[" + r.Name + "]" + item.Edges.OutOperateManager.Name
					}
				}
			case model.OperatorTypeEmployee:
				if item.Edges.OutOperateStore != nil {
					fromOperateName = "[门店]" + item.Edges.OutOperateStore.Name
				}
			case model.OperatorTypeAgent:
				if item.Edges.OutOperateAgent != nil {
					fromOperateName = "[代理]" + item.Edges.OutOperateAgent.Name
				}
			case model.OperatorTypeMaintainer:
				if item.Edges.OutOperateMaintainer != nil {
					fromOperateName = "[运维]" + item.Edges.OutOperateMaintainer.Name
				}
			case model.OperatorTypeCabinet:
				if item.Edges.OutOperateCabinet != nil {
					fromOperateName = "[电柜]" + item.Edges.OutOperateCabinet.Name
				}
			case model.OperatorTypeRider:
				if item.Edges.OutOperateRider != nil {
					fromOperateName = "[骑手]" + item.Edges.OutOperateRider.Name
				}
			default:
			}
		}

		in := model.AssetTransferDetailList{
			OperatorName:  toOperateName,
			LocationsName: toLocationName,
			TimeAt:        inTimeAt,
			Remark:        inRemark,
			Num:           item.InNum,
		}
		res = &model.AssetTransferDetailListRes{
			CityName:         cityName,
			AssetName:        assetName,
			In:               in,
			TransferTypeName: model.AssetTransferType(item.Type).String(),
			TransferType:     model.AssetTransferType(item.Type),
		}
		if fromLocationName != "" && fromOperateName != "" {
			res.Out = &model.AssetTransferDetailList{
				LocationsName: fromLocationName,
				OperatorName:  fromOperateName,
				Remark:        item.Remark,
				TimeAt:        outTimeAt,
				Num:           item.OutNum,
			}
		}
		return res
	}), nil
}

// Modify 编辑调拨
func (s *assetTransferService) Modify(ctx context.Context, req *model.AssetTransferModifyReq, modifier *model.Modifier) (err error) {
	item, _ := s.orm.QueryNotDeleted().
		Where(
			assettransfer.ID(req.ID),
			assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
		).
		WithFromLocationStation(func(query *ent.EnterpriseStationQuery) {
			query.WithAgents()
		}).WithFromLocationWarehouse(
		func(query *ent.WarehouseQuery) {
			query.WithBelongAssetManagers()
		}).WithFromLocationStore(func(query *ent.StoreQuery) {
		query.WithEmployees()
	}).
		First(ctx)
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
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
