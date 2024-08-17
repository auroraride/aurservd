package service

import (
	"context"
	"fmt"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
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
	bulk := make([]*ent.AssetCheckDetailsCreate, 0)
	var realEbikeNum, realBatteryNum uint
	for _, v := range req.AssetCheckCreateDetail {
		err := s.checkAssetCheckOwner(ctx, model.CheckAssetCheckOwnerReq{})
		if err != nil {
			return err
		}
		if v.AssetType == model.AssetTypeEbike {
			realEbikeNum++
		} else {
			realBatteryNum++
		}

		bulk = append(bulk, ent.Database.AssetCheckDetails.Create().SetAssetID(v.AssetID).SetCreator(modifier).SetLastModifier(modifier))
	}
	b, err := ent.Database.AssetCheckDetails.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return err
	}
	// 查询应盘电车资产
	allEbike, err := s.getCheckAsset(ctx, model.GetCheckAssetReq{
		OpratorType:   req.OpratorType,
		OpratorID:     req.OpratorID,
		LocationsType: model.AssetLocationsTypeStation,
		LocationsID:   req.LocationsID,
		AssetType:     model.AssetTypeEbike,
	})
	if err != nil {
		return err
	}
	// 查询应盘电池资产
	allBattery, err := s.getCheckAsset(ctx, model.GetCheckAssetReq{
		OpratorType:   req.OpratorType,
		OpratorID:     req.OpratorID,
		LocationsType: model.AssetLocationsTypeStation,
		LocationsID:   req.LocationsID,
		AssetType:     model.AssetTypeSmartBattery,
	})
	if err != nil {
		return err
	}

	err = s.orm.Create().
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
		AddCheckDetails(b...).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 查询应盘点资产
func (s *assetCheckService) getCheckAsset(ctx context.Context, req model.GetCheckAssetReq) (res []*ent.Asset, err error) {
	t, err := s.getAssetByOperateRole(ctx, model.GetAssetByOperateRole{
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
func (s *assetCheckService) getAssetByOperateRole(ctx context.Context, req model.GetAssetByOperateRole) (res bool, err error) {
	switch req.OpratorType {
	case model.AssetOperateRoleTypeAgent:
		item, _ := ent.Database.Agent.QueryNotDeleted().Where(agent.ID(req.OpratorID)).WithStations().First(ctx)
		if item == nil {
			return false, fmt.Errorf("未找到对应的代理商")
		}
		if len(item.Edges.Stations) == 0 {
			return false, fmt.Errorf("代理商未绑定站点")
		}
		for _, v := range item.Edges.Stations {
			if v.ID == req.LocationsID {
				return true, nil
			}
		}
	case model.AssetOperateRoleTypeStore:
		// 查询能操作哪些门店的资产
		// todo 门店集合没写！还有上班限制
		return true, err
	case model.AssetOperateRoleTypeManager:
		// 查询能操作哪些仓库的资产
		// todo 仓库集合没写！还有上班限制
		return true, err
	default:
		return false, fmt.Errorf("未知的操作角色类型")
	}
	return false, nil
}

// 验证资产盘点是否属于当前用户
func (s *assetCheckService) checkAssetCheckOwner(ctx context.Context, req model.CheckAssetCheckOwnerReq) error {
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.ID(req.AssetID),
			asset.Status(model.AssetStatusStock.Value()),
			asset.Type(req.AssetType.Value()),
			asset.Type(req.LocationsType.Value()),
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

// GetAssetCheck 盘点资产查询
func (s *assetCheckService) GetAssetCheck(ctx context.Context, req *model.AssetCheckGetReq) (res *model.AssetCheckGetRes, err error) {
	// assetID, err := ent.Database.Asset.QueryNotDeleted().Where(asset.Sn(req.SN)).First(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// res = &model.AssetCheckGetRes{
	// 	StartAt:        "",
	// 	EndAt:          "",
	// 	OpratorID:      0,
	// 	OpratorName:    "",
	// 	BatteryNum:     0,
	// 	BatteryNumReal: 0,
	// 	EbikeNum:       0,
	// 	EbikeNumReal:   0,
	// 	LocationsID:    0,
	// 	LocationsType:  0,
	// 	CheckResult:    0,
	// 	Abnormal:       nil,
	// }

	return nil, nil
}
