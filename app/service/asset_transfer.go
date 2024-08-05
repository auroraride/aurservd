package service

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/manager"
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
		assetIDs, failed = s.stockTransfer(ctx, req, modifier)
		if len(failed) > 0 {
			return failed, nil
		}
		q.SetFromLocationType(req.FromLocationType.Value()).
			SetFromLocationID(*req.FromLocationID).
			SetStatus(model.AssetTransferStatusDelivering.Value()).
			SetOutTimeAt(time.Now()).
			SetOutOperateID(modifier.ID).
			SetOutOperateType(model.AssetOperateRoleManager.Value()).
			SetOutNum(uint(len(assetIDs)))
	}
	// 初始调拨
	if req.FromLocationType == nil {
		assetIDs, failed = s.initialTransfer(ctx, req, modifier)
		if len(failed) > 0 {
			return failed, nil
		}
		q.SetInNum(uint(len(assetIDs))).
			SetInTimeAt(time.Now()).
			SetStatus(model.AssetTransferStatusStock.Value()).
			SetInOperateID(modifier.ID).
			SetInOperateType(model.AssetOperateRoleManager.Value()).
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
func (s *assetTransferService) stockTransfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (assetIDs []uint64, failed []string) {
	// 查询物资是否充足
	q := ent.Database.Asset.QueryNotDeleted().Where(
		asset.LocationsType((*req.FromLocationType).Value()),
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
	)

	var err error
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
	return assetIDs, failed
}

// 初始调拨
func (s *assetTransferService) initialTransfer(ctx context.Context, req *model.AssetTransferCreateReq, modifier *model.Modifier) (assetIDs []uint64, failed []string) {
	var err error
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
	return assetIDs, failed
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

// TransferList 调拨列表
func (s *assetTransferService) TransferList(ctx context.Context, req *model.AssetTransferListReq) (res *model.PaginationRes, err error) {
	q := ent.Database.AssetTransfer.QueryNotDeleted().WithDetails()
	s.filter(ctx, q, &req.AssetTransferFilter)

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetTransfer) (res *model.AssetTransferListRes) {
		res = &model.AssetTransferListRes{
			ID:        item.ID,
			SN:        item.Sn,
			Reason:    item.Reason,
			Remark:    item.Remark,
			Status:    model.AssetTransferStatus(item.Status).String(),
			OutNum:    item.OutNum,
			InNum:     item.InNum,
			OutTimeAt: item.OutTimeAt.Format("2006-01-02 15:04:05"),
			InTimeAt:  item.InTimeAt.Format("2006-01-02 15:04:05"),
		}

		if item.FromLocationType != nil && item.FromLocationID != nil {
			switch model.AssetLocationsType(*item.FromLocationType) {
			case model.AssetLocationsTypeWarehouse:
				if item.Edges.LocationWarehouse != nil {
					res.FromLocationName = "[仓库]" + item.Edges.LocationWarehouse.Name
				}
			case model.AssetLocationsTypeStore:
				if item.Edges.LocationStore != nil {
					res.FromLocationName = "[门店]" + item.Edges.LocationStore.Name
				}
			case model.AssetLocationsTypeStation:
				if item.Edges.LocationStation != nil {
					res.FromLocationName = "[站点]" + item.Edges.LocationStation.Name
				}
			case model.AssetLocationsTypeOperation:
				if item.Edges.LocationOperator != nil {
					res.FromLocationName = "[运维]" + item.Edges.LocationOperator.Name
				}
			default:
			}
		}

		if item.ToLocationType != 0 && item.ToLocationID != 0 {
			switch model.AssetLocationsType(item.ToLocationType) {
			case model.AssetLocationsTypeWarehouse:
				if item.Edges.LocationWarehouse != nil {
					res.ToLocationName = "[仓库]" + item.Edges.LocationWarehouse.Name
				}
			case model.AssetLocationsTypeStore:
				if item.Edges.LocationStore != nil {
					res.ToLocationName = "[门店]" + item.Edges.LocationStore.Name
				}
			case model.AssetLocationsTypeStation:
				if item.Edges.LocationStation != nil {
					res.ToLocationName = "[站点]" + item.Edges.LocationStation.Name
				}
			case model.AssetLocationsTypeOperation:
				if item.Edges.LocationOperator != nil {
					res.ToLocationName = "[运维]" + item.Edges.LocationOperator.Name
				}
			default:
			}
		}

		// 出库操作人
		if item.OutOperateType != nil && item.OutOperateID != nil {
			switch model.AssetOperateRoleType(*item.OutOperateType) {
			case model.AssetOperateRoleManager:
				if item.Edges.OutOperateManager != nil {
					// 查询角色
					var roleName string
					if role, _ := item.Edges.OutOperateManager.QueryRole().First(ctx); role != nil {
						roleName = role.Name
					}
					res.OutOperateName = "[" + roleName + "]" + item.Edges.OutOperateManager.Name
				}
			case model.AssetOperateRoleStore:
				if item.Edges.OutOperateStore != nil {
					res.OutOperateName = "[门店]" + item.Edges.OutOperateStore.Name
				}
			case model.AssetOperateRoleAgent:
				if item.Edges.OutOperateAgent != nil {
					res.OutOperateName = "[代理]" + item.Edges.OutOperateAgent.Name
				}
			case model.AssetOperateRoleOperation:
				if item.Edges.OutOperateMaintainer != nil {
					res.OutOperateName = "[运维]" + item.Edges.OutOperateMaintainer.Name
				}
			default:
			}
		}

		// 入库操作人
		if item.InOperateType != 0 && item.InOperateID != 0 {
			switch model.AssetOperateRoleType(item.InOperateType) {
			case model.AssetOperateRoleManager:
				if item.Edges.InOperateManager != nil {
					res.InOperateName = item.Edges.InOperateManager.Name
				}
			case model.AssetOperateRoleStore:
				if item.Edges.InOperateStore != nil {
					res.InOperateName = item.Edges.InOperateStore.Name
				}
			case model.AssetOperateRoleAgent:
				if item.Edges.InOperateAgent != nil {
					res.InOperateName = item.Edges.InOperateAgent.Name
				}
			case model.AssetOperateRoleOperation:
				if item.Edges.InOperateMaintainer != nil {
					res.InOperateName = item.Edges.InOperateMaintainer.Name
				}
			default:
			}
		}
		return res
	}), nil

}

// 筛选
func (s *assetTransferService) filter(ctx context.Context, q *ent.AssetTransferQuery, req *model.AssetTransferFilter) {
	if req.FromLocationType != nil {
		q.Where(assettransfer.FromLocationType((*req.FromLocationType).Value()))
	}
	if req.FromLocationID != nil {
		q.Where(assettransfer.FromLocationID(*req.FromLocationID))
	}
	if req.ToLocationType != nil {
		q.Where(assettransfer.ToLocationType((*req.ToLocationType).Value()))
	}
	if req.ToLocationID != nil {
		q.Where(assettransfer.ToLocationID(*req.ToLocationID))
	}
	if req.Status != nil {
		q.Where(assettransfer.Status((*req.Status).Value()))
	}
	if req.OutStart != nil && req.OutEnd != nil {
		start := tools.NewTime().ParseDateStringX(*req.OutStart)
		end := tools.NewTime().ParseNextDateStringX(*req.OutEnd)
		q.Where(assettransfer.OutTimeAtGTE(start), assettransfer.OutTimeAtLTE(end))
	}
	if req.InStart != nil && req.InEnd != nil {
		start := tools.NewTime().ParseDateStringX(*req.InStart)
		end := tools.NewTime().ParseNextDateStringX(*req.InEnd)
		q.Where(assettransfer.InTimeAtGTE(start), assettransfer.InTimeAtLTE(end))
	}
	if req.Keyword != nil {
		q.Where(
			assettransfer.Or(
				assettransfer.SnContains(*req.Keyword),
				assettransfer.ReasonContains(*req.Keyword),
				// 资产后台
				assettransfer.HasInOperateManagerWith(manager.NameContains(*req.Keyword)),
				assettransfer.HasOutOperateManagerWith(manager.NameContains(*req.Keyword)),
				// 门店
				assettransfer.HasInOperateStoreWith(store.NameContains(*req.Keyword)),
				assettransfer.HasOutOperateStoreWith(store.NameContains(*req.Keyword)),
				// 代理
				assettransfer.HasInOperateAgentWith(agent.NameContains(*req.Keyword)),
				assettransfer.HasOutOperateAgentWith(agent.NameContains(*req.Keyword)),
				// 运维
				assettransfer.HasInOperateMaintainerWith(maintainer.NameContains(*req.Keyword)),
				assettransfer.HasOutOperateMaintainerWith(maintainer.NameContains(*req.Keyword)),
			),
		)
	}
}

// TransferDetail 调拨详情
func (s *assetTransferService) TransferDetail(ctx context.Context, req *model.AssetTransferDetailReq) (res *model.AssetTransferDetail, err error) {
	// todo 详情
	// var result struct {
	// 	OutNum    int    `json:"out_num"`    // 出库数量
	// 	InNum     int    `json:"in_num"`     // 入库数量
	// 	Name      string `json:"name"`       // 资产名称
	// 	SN        string `json:"sn"`         // 资产编号
	// 	AssetType int    `json:"asset_type"` // 资产类型
	// }
	//
	// q := ent.Database.AssetTransferDetails.
	// 	QueryNotDeleted().
	// 	Where(assettransferdetails.TransferID(req.ID)).
	// 	Modify(func(sel *sql.Selector) {
	// 		a := sql.Table(asset.Table)
	// 		m := sql.Table(material.Table)
	// 		sel.LeftJoin(a).On(sel.C(assettransferdetails.FieldAssetID), a.C(asset.FieldID))
	// 		sel.LeftJoin(m).On(a.C(asset.FieldMaterialID), m.C(material.FieldID))
	// 		sel.Select(a.C(asset.FieldSn), a.C(asset.FieldType), a.C(asset.FieldName), m.C(material.FieldName))
	// 	}).GroupBy(asset.FieldSn, asset.FieldMaterialID, asset.FieldType)
	//
	// q.Aggregate(ent.Sum(stock.FieldNum)).
	// 	Scan(ctx, &result)
	return res, nil
}
