package service

import (
	"context"
	"fmt"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetcheck"
	"github.com/auroraride/aurservd/internal/ent/assetcheckdetails"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetCheckService struct {
	orm *ent.AssetCheckClient
}

func NewAssetCheck() *assetCheckService {
	return &assetCheckService{
		orm: ent.Database.AssetCheck,
	}
}

// CreateAssetCheck 创建资产盘点
func (s *assetCheckService) CreateAssetCheck(ctx context.Context, req *model.AssetCheckCreateReq, modifier *model.Modifier) error {
	// 应盘资产
	var assetIDs []uint64
	// 查询应盘电车资产
	allEbike, err := s.getCheckAsset(ctx, model.GetCheckAssetReq{
		OpratorType:   req.OpratorType,
		OpratorID:     req.OpratorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
		AssetType:     model.AssetTypeEbike,
	})
	if err != nil {
		return err
	}
	bulk := make([]*ent.AssetCheckDetailsCreate, 0)
	for _, v := range allEbike {
		bulk = append(bulk, ent.Database.AssetCheckDetails.Create().
			SetAsset(v).
			SetCreator(modifier).
			SetLastModifier(modifier).
			SetStatus(model.AssetCheckDetailsStatusUntreated.Value()).
			SetLocationsID(v.LocationsID).
			SetLocationsType(v.LocationsType))
		assetIDs = append(assetIDs, v.ID)
	}
	// 查询应盘电池资产
	allBattery, err := s.getCheckAsset(ctx, model.GetCheckAssetReq{
		OpratorType:   req.OpratorType,
		OpratorID:     req.OpratorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
		AssetType:     model.AssetTypeSmartBattery,
	})
	if err != nil {
		return err
	}
	for _, v := range allBattery {
		bulk = append(bulk, ent.Database.AssetCheckDetails.Create().
			SetAsset(v).
			SetCreator(modifier).
			SetLastModifier(modifier).
			SetStatus(model.AssetCheckDetailsStatusUntreated.Value()).
			SetLocationsID(v.LocationsID).
			SetLocationsType(v.LocationsType))
		assetIDs = append(assetIDs, v.ID)
	}
	var realEbikeNum, realBatteryNum uint
	var realAssetIDs []uint64
	for _, v := range req.AssetCheckCreateDetail {
		// 限制只能盘点自己的资产
		err = s.checkAssetCheckOwner(ctx, model.CheckAssetCheckOwnerReq{
			AssetID:       v.AssetID,
			AssetType:     v.AssetType,
			OpratorType:   req.OpratorType,
			OpratorID:     req.OpratorID,
			LocationsType: req.LocationsType,
			LocationsID:   req.LocationsID,
		})
		if err != nil {
			return err
		}
		// 获取实际盘点数量
		if v.AssetType == model.AssetTypeEbike {
			realEbikeNum++
		} else {
			realBatteryNum++
		}
		realAssetIDs = append(realAssetIDs, v.AssetID)
	}
	b, err := ent.Database.AssetCheckDetails.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return err
	}

	start := tools.NewTime().ParseDateStringX(req.StartAt)
	end := tools.NewTime().ParseNextDateStringX(req.EndAt)

	c, err := s.orm.Create().
		SetCreator(modifier).
		SetLastModifier(modifier).
		SetOperateID(req.OpratorID).
		SetOperateType(req.OpratorType.Value()).
		SetLocationsID(req.LocationsID).
		SetLocationsType(req.LocationsType.Value()).
		SetEbikeNum(uint(len(allEbike))).
		SetBatteryNum(uint(len(allBattery))).
		SetBatteryNumReal(realBatteryNum).
		SetEbikeNumReal(realEbikeNum).
		SetStartAt(start).
		SetEndAt(end).
		AddCheckDetails(b...).
		Save(ctx)
	if err != nil {
		return err
	}
	// 多出来的资产ID 和 未盘点的资产ID
	missingIDs, extraIDs := s.findMissingAndExtraIDs(assetIDs, realAssetIDs)
	_, err = ent.Database.AssetCheckDetails.Update().
		Where(
			assetcheckdetails.AssetIDIn(missingIDs...),
			assetcheckdetails.CheckID(c.ID),
		).
		SetResult(model.AssetCheckResultLoss.Value()).
		SetRealLocationsID(req.LocationsID).
		SetRealLocationsType(req.LocationsType.Value()).
		Save(ctx)
	if err != nil {
		return err
	}
	_, err = ent.Database.AssetCheckDetails.Update().
		Where(
			assetcheckdetails.AssetIDIn(extraIDs...),
			assetcheckdetails.CheckID(c.ID),
		).
		SetResult(model.AssetCheckResultSurplus.Value()).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *assetCheckService) findMissingAndExtraIDs(assetIDs []uint64, realAssetIDs []uint64) ([]uint64, []uint64) {
	// 用 map 记录实际盘点资产的 ID
	realAssetMap := make(map[uint64]bool)
	for _, assetID := range realAssetIDs {
		realAssetMap[assetID] = true
	}

	// 分别记录未被盘点和多出来的资产ID
	missingIDs := make([]uint64, 0)
	extraIDs := make([]uint64, 0)

	// 找出未被盘点的资产ID
	for _, assetID := range assetIDs {
		if _, ok := realAssetMap[assetID]; !ok {
			missingIDs = append(missingIDs, assetID)
		} else {
			delete(realAssetMap, assetID) // 从 map 中删除已经盘点的资产ID
		}
	}

	// 剩下的就是多出来的资产ID
	for assetID := range realAssetMap {
		extraIDs = append(extraIDs, assetID)
	}

	return missingIDs, extraIDs
}

// 查询应盘点资产
func (s *assetCheckService) getCheckAsset(ctx context.Context, req model.GetCheckAssetReq) (res []*ent.Asset, err error) {
	t, _, err := s.getAssetByOperateRole(ctx, model.GetAssetByOperateRole{
		OpratorType:   req.OpratorType,
		OpratorID:     req.OpratorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
	})
	if err != nil {
		return nil, err
	}
	if !t {
		return nil, fmt.Errorf("当前用户无操作权限")
	}
	items, _ := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.LocationsType(req.LocationsType.Value()),
			asset.Type(req.AssetType.Value()),
			asset.LocationsID(req.LocationsID),
		).All(ctx)
	if items == nil {
		return nil, fmt.Errorf("未找到对应的资产")
	}
	return items, nil
}

// 根据操作角色找对应能操作的资产
func (s *assetCheckService) getAssetByOperateRole(ctx context.Context, req model.GetAssetByOperateRole) (b bool, ids []uint64, err error) {
	ids = make([]uint64, 0)
	b = false

	switch req.OpratorType {
	case model.AssetOperateRoleTypeAgent:
		item, _ := ent.Database.Agent.QueryNotDeleted().Where(agent.ID(req.OpratorID)).WithStations().First(ctx)
		if item == nil {
			return b, nil, fmt.Errorf("未找到对应的代理商")
		}
		if len(item.Edges.Stations) == 0 {
			return b, nil, fmt.Errorf("代理商未绑定站点")
		}
		for _, v := range item.Edges.Stations {
			if v.ID == req.LocationsID {
				b = true
			}
			ids = append(ids, v.ID)
		}
	case model.AssetOperateRoleTypeStore:
		// 查询能操作哪些门店的资产
		// todo 门店集合没写！还有上班限制
		b = true
	case model.AssetOperateRoleTypeManager:
		// 查询能操作哪些仓库的资产
		// todo 仓库集合没写！还有上班限制
		b = true
	default:
		return b, nil, fmt.Errorf("未知的操作人类型")
	}
	return b, ids, nil
}

// 验证资产盘点是否属于当前用户
func (s *assetCheckService) checkAssetCheckOwner(ctx context.Context, req model.CheckAssetCheckOwnerReq) error {
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.ID(req.AssetID),
			asset.Status(model.AssetStatusStock.Value()),
			asset.Type(req.AssetType.Value()),
			asset.LocationsType(req.LocationsType.Value()),
			asset.LocationsID(req.LocationsID),
		)
	item, _ := q.First(ctx)
	if item == nil {
		return fmt.Errorf("资产不存在或不属于当前用户")
	}
	return nil
}

// GetAssetBySN 通过sn查询资产
func (s *assetCheckService) GetAssetBySN(ctx context.Context, req *model.AssetCheckByAssetSnReq) (res *model.AssetCheckByAssetSnRes, err error) {
	item, _ := ent.Database.Asset.QueryNotDeleted().Where(asset.Sn(req.SN)).WithModel().WithBrand().First(ctx)
	if item == nil {
		return nil, fmt.Errorf("未找到对应的资产")
	}
	// 验证资产是否属于当前用户
	err = s.checkAssetCheckOwner(ctx, model.CheckAssetCheckOwnerReq{
		AssetID:       item.ID,
		AssetType:     model.AssetType(item.Type),
		OpratorType:   req.OpratorType,
		OpratorID:     req.OpratorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
	})
	if err != nil {
		return nil, err
	}
	res = &model.AssetCheckByAssetSnRes{
		AssetID:   item.ID,
		AssetSN:   item.Sn,
		AssetType: model.AssetType(item.Type),
	}
	if item.Edges.Model != nil {
		res.Model = item.Edges.Model.Model
	}
	if item.Edges.Brand != nil {
		res.BrandName = item.Edges.Brand.Name
	}
	return res, nil
}

// MarkStartOrEndCheck 标记开始盘点或结束盘点
func (s *assetCheckService) MarkStartOrEndCheck(ctx context.Context, req *model.MarkStartOrEndCheckReq) (err error) {
	t, _, err := s.getAssetByOperateRole(ctx, model.GetAssetByOperateRole{
		OpratorType:   req.OpratorType,
		OpratorID:     req.OpratorID,
		LocationsType: req.LocationsType,
		LocationsID:   req.LocationsID,
	})
	if err != nil {
		return err
	}
	if !t {
		return fmt.Errorf("当前用户无操作权限")
	}
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.LocationsType(req.LocationsType.Value()),
			asset.LocationsID(req.LocationsID),

			asset.Status(model.AssetStatusStock.Value()),
		)
	if req.Enable {
		// 标记开始盘点
		q.Where(asset.CheckAtIsNil())
	} else {
		// 标记结束盘点
		q.Where(asset.CheckAtNotNil())
	}
	all, _ := q.All(ctx)
	if all == nil {
		return fmt.Errorf("未找到对应的资产")
	}
	assetIDs := make([]uint64, 0)
	for _, v := range all {
		assetIDs = append(assetIDs, v.ID)
	}
	update := ent.Database.Asset.Update().Where(asset.IDIn(assetIDs...))
	if req.Enable {
		update.SetCheckAt(time.Now())
	} else {
		update.SetNillableCheckAt(nil)
	}
	_, err = update.Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// List 盘点列表
func (s *assetCheckService) List(ctx context.Context, req *model.AssetCheckListReq) (*model.PaginationRes, error) {
	q := ent.Database.AssetCheck.QueryNotDeleted().WithStore().WithStation().WithWarehouse().
		WithOperateStore().WithOperateManager().WithOperateAgent().WithCheckDetails(func(query *ent.AssetCheckDetailsQuery) {
		query.WithAsset()
	})
	if req.LocationsID != nil && req.LocationsType != nil {
		q.Where(
			assetcheck.LocationsID(*req.LocationsID),
			assetcheck.LocationsType(req.LocationsType.Value()),
		)
	}

	if req.LocationsID != nil && req.LocationsType != nil {
		q.Where(
			assetcheck.LocationsID(*req.LocationsID),
			assetcheck.LocationsType(req.LocationsType.Value()),
		)
	}

	if req.StartAt != nil && req.EndAt != nil {
		start := tools.NewTime().ParseDateStringX(*req.StartAt)
		end := tools.NewTime().ParseNextDateStringX(*req.EndAt)
		q.Where(assetcheck.EndAtGTE(start), assetcheck.EndAtLTE(end))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetCheck) *model.AssetCheckListRes {
		var start, end string
		if item.StartAt != nil {
			start = item.StartAt.Format("2006-01-02 15:04:05")
		}
		if item.EndAt != nil {
			end = item.EndAt.Format("2006-01-02 15:04:05")
		}
		opratorName := ""
		switch model.AssetOperateRoleType(item.OperateType) {
		case model.AssetOperateRoleTypeAgent:
			opratorName = item.Edges.OperateAgent.Name
		case model.AssetOperateRoleTypeStore:
			opratorName = item.Edges.OperateStore.Name
		case model.AssetOperateRoleTypeManager:
			opratorName = item.Edges.OperateManager.Name
		default:
			opratorName = ""
		}
		locationsName := ""
		switch model.AssetLocationsType(item.LocationsType) {
		case model.AssetLocationsTypeStore:
			locationsName = "[平台门店]-" + item.Edges.Store.Name
		case model.AssetLocationsTypeStation:
			locationsName = "[代理站点]-" + item.Edges.Station.Name
		case model.AssetLocationsTypeWarehouse:
			locationsName = "[平台仓库]-" + item.Edges.Warehouse.Name
		default:
			locationsName = ""
		}

		var checkResult bool
		if item.BatteryNum == item.BatteryNumReal && item.EbikeNum == item.EbikeNumReal {
			checkResult = true
		} else {
			checkResult = false
		}

		res := &model.AssetCheckListRes{
			StartAt:        start,
			EndAt:          end,
			OpratorID:      item.OperateID,
			OpratorName:    opratorName,
			BatteryNum:     item.BatteryNum,
			BatteryNumReal: item.BatteryNumReal,
			EbikeNum:       item.EbikeNum,
			EbikeNumReal:   item.EbikeNumReal,
			LocationsID:    item.LocationsID,
			LocationsType:  item.LocationsType,
			LocationsName:  locationsName,
			CheckResult:    checkResult,
		}
		return res
	}), nil
}

// ListAbnormal 查询盘点异常资产
func (s *assetCheckService) ListAbnormal(ctx context.Context, req *model.AssetCheckListAbnormalReq) (res []*model.AssetCheckAbnormal, err error) {
	// 查询异常资产
	abnormalAll, _ := ent.Database.AssetCheckDetails.QueryNotDeleted().
		Where(
			assetcheckdetails.ResultNEQ(model.AssetCheckResultNormal.Value()),
			assetcheckdetails.CheckID(req.AssetCheckID),
		).WithAsset(func(query *ent.AssetQuery) {
		query.WithModel().WithBrand().WithMaterial()
	}).WithWarehouse().WithStore().WithStation().WithRider().WithCabinet().WithOperator().
		WithRealWarehouse().WithRealStation().WithRealStore().WithRealCabinet().WithRealRider().WithRealOperator().
		All(context.Background())

	for _, v := range abnormalAll {
		var name, sn, modelName, brandName, locationsName, realLocationsName, opratorName string
		if v.Edges.Asset != nil {
			name = v.Edges.Asset.Name
			sn = v.Edges.Asset.Sn
			if v.Edges.Asset.Edges.Model != nil {
				modelName = v.Edges.Asset.Edges.Model.Model
			}
			if v.Edges.Asset.Edges.Brand != nil {
				brandName = v.Edges.Asset.Edges.Brand.Name
			}
		}

		switch model.AssetLocationsType(v.LocationsType) {
		case model.AssetLocationsTypeStore:
			if v.Edges.Store != nil {
				locationsName = "[平台门店]-" + v.Edges.Store.Name
			}
		case model.AssetLocationsTypeStation:
			if v.Edges.Station != nil {
				locationsName = "[代理站点]-" + v.Edges.Station.Name
			}
		case model.AssetLocationsTypeWarehouse:
			if v.Edges.Warehouse != nil {
				locationsName = "[平台仓库]-" + v.Edges.Warehouse.Name
			}
		case model.AssetLocationsTypeRider:
			if v.Edges.Rider != nil {
				locationsName = "[骑手]-" + v.Edges.Rider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if v.Edges.Cabinet != nil {
				locationsName = "[电柜]-" + v.Edges.Cabinet.Name
			}
		case model.AssetLocationsTypeOperation:
			if v.Edges.Operator != nil {
				locationsName = "[运维]-" + v.Edges.Operator.Name
			}
		default:
			locationsName = ""
		}

		switch model.AssetLocationsType(v.RealLocationsType) {
		case model.AssetLocationsTypeStore:
			if v.Edges.RealStore != nil {
				realLocationsName = "[平台门店]-" + v.Edges.RealStore.Name
			}
		case model.AssetLocationsTypeStation:
			if v.Edges.RealStation != nil {
				realLocationsName = "[代理站点]-" + v.Edges.RealStation.Name
			}
		case model.AssetLocationsTypeWarehouse:
			if v.Edges.RealWarehouse != nil {
				realLocationsName = "[平台仓库]-" + v.Edges.RealWarehouse.Name
			}
		case model.AssetLocationsTypeRider:
			if v.Edges.RealRider != nil {
				realLocationsName = "[骑手]-" + v.Edges.RealRider.Name
			}
		case model.AssetLocationsTypeCabinet:
			if v.Edges.RealCabinet != nil {
				realLocationsName = "[电柜]-" + v.Edges.RealCabinet.Name
			}
		case model.AssetLocationsTypeOperation:
			if v.Edges.RealOperator != nil {
				realLocationsName = "[运维]-" + v.Edges.RealOperator.Name
			}
		default:
			realLocationsName = ""
		}

		if v.Edges.Maintainer != nil {
			opratorName = v.Edges.Maintainer.Name
		}

		res = append(res, &model.AssetCheckAbnormal{
			AssetID:           v.AssetID,
			Name:              name,
			Model:             modelName,
			Brand:             brandName,
			Result:            model.AssetCheckResult(v.Result).String(),
			SN:                sn,
			LocationsName:     locationsName,
			RealLocationsName: realLocationsName,
			Status:            model.AssetCheckDetailsStatus(v.Status).String(),
			OpratorName:       opratorName,
			OpratorAt:         v.OperateAt.Format("2006-01-02 15:04:05"),
		})
	}
	return res, nil
}
