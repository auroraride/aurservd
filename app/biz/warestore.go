// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package biz

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
	"github.com/auroraride/aurservd/pkg/utils"
)

type warestoreBiz struct {
	ctx                    context.Context
	warestoreTokenCacheKey string
}

func NewWarestore() *warestoreBiz {
	return &warestoreBiz{
		ctx:                    context.Background(),
		warestoreTokenCacheKey: ar.Config.Environment.UpperString() + ":" + "WARESTORE:TOKEN",
	}
}

// signin 仓库、门店登录
func (b *warestoreBiz) signin(am *ent.AssetManager, ep *ent.Employee, platType definition.PlatType) (res *definition.WarestorePeopleSigninRes, err error) {
	var token string
	tokenKey := b.warestoreTokenCacheKey
	switch {
	case platType == definition.PlatTypeWarehouse && am != nil:
		idstr := definition.SignTokenWarehouse + "-" + strconv.FormatUint(am.ID, 10)
		// 查询并删除旧token key
		exists := ar.Redis.HGet(b.ctx, tokenKey, idstr).Val()
		if exists != "" {
			ar.Redis.HDel(b.ctx, tokenKey, exists)
		}

		// 生成token
		token = utils.NewEcdsaToken()

		// 存储登录token和ID进行对应
		ar.Redis.HSet(b.ctx, tokenKey, token, idstr)
		ar.Redis.HSet(b.ctx, tokenKey, idstr, token)
	case platType == definition.PlatTypeStore && ep != nil:
		idstr := definition.SignTokenStore + "-" + strconv.FormatUint(ep.ID, 10)
		// 查询并删除旧token key
		exists := ar.Redis.HGet(b.ctx, tokenKey, idstr).Val()
		if exists != "" {
			ar.Redis.HDel(b.ctx, tokenKey, exists)
		}

		// 生成token
		token = utils.NewEcdsaToken()

		// 存储登录token和ID进行对应
		ar.Redis.HSet(b.ctx, tokenKey, token, idstr)
		ar.Redis.HSet(b.ctx, tokenKey, idstr, token)
	default:
		return nil, errors.New("登录平台失败")
	}

	return &definition.WarestorePeopleSigninRes{
		Profile: b.Profile(am, ep, platType),
		Token:   token,
	}, nil
}

// Signin 登录
func (b *warestoreBiz) Signin(req *definition.WarestorePeopleSigninReq) (res *definition.WarestorePeopleSigninRes, err error) {
	am := new(ent.AssetManager)
	ep := new(ent.Employee)
	switch req.PlatType {
	case definition.PlatTypeWarehouse:
		am, err = ent.Database.AssetManager.QueryNotDeleted().
			Where(
				assetmanager.Phone(req.Phone),
				assetmanager.MiniEnable(true),
			).First(b.ctx)
		if am == nil || err != nil {
			return nil, errors.New("账号不存在")
		}

		// 比对密码
		if !utils.PasswordCompare(req.Password, am.Password) {
			return nil, errors.New(ar.UserAuthenticationFailed)
		}

	case definition.PlatTypeStore:
		ep, err = ent.Database.Employee.QueryNotDeleted().Where(employee.Phone(req.Phone)).First(b.ctx)
		if ep == nil || err != nil {
			return nil, errors.New("账号不存在")
		}

		// 比对密码
		if !utils.PasswordCompare(req.Password, ep.Password) {
			return nil, errors.New(ar.UserAuthenticationFailed)
		}
	}

	return b.signin(am, ep, req.PlatType)
}

// Profile 仓管资料
func (b *warestoreBiz) Profile(am *ent.AssetManager, ep *ent.Employee, platType definition.PlatType) definition.WarestorePeopleProfile {
	switch {
	case platType == definition.PlatTypeWarehouse && am != nil:
		// todo 上班信息
		return definition.WarestorePeopleProfile{
			ID:    am.ID,
			Phone: am.Phone,
			Name:  am.Name,
		}

	case platType == definition.PlatTypeStore && ep != nil:
		// todo 上班信息
		return definition.WarestorePeopleProfile{
			ID:    ep.ID,
			Phone: ep.Phone,
			Name:  ep.Name,
		}
	default:
		return definition.WarestorePeopleProfile{}
	}
}

// AssetCount 仓管资产统计
func (b *warestoreBiz) AssetCount(am *ent.AssetManager, ep *ent.Employee) definition.AssetCountRes {
	return definition.AssetCountRes{}
}

// TokenVerify Token校验
func (b *warestoreBiz) TokenVerify(token string) (am *ent.AssetManager, ep *ent.Employee) {
	// 获取token对应ID
	tokenVal := ar.Redis.HGet(b.ctx, b.warestoreTokenCacheKey, token).Val()
	vals := strings.Split(tokenVal, "-")
	// 解析的数据不为两组数据则直接返回
	if len(vals) != 2 {
		return
	}
	platType := vals[0]
	wsId, _ := strconv.Atoi(vals[1])
	// 判断仓管类型取出人员信息
	switch platType {
	case definition.SignTokenWarehouse:
		// 反向校验token是否正确
		if ar.Redis.HGet(b.ctx, b.warestoreTokenCacheKey, definition.SignTokenWarehouse+"-"+strconv.FormatUint(uint64(wsId), 10)).Val() != token {
			return
		}
		// 获取库管人员
		am, _ = ent.Database.AssetManager.QueryNotDeleted().Where(assetmanager.ID(uint64(wsId)), assetmanager.MiniEnable(true)).First(b.ctx)
	case definition.SignTokenStore:
		// 反向校验token是否正确
		if ar.Redis.HGet(b.ctx, b.warestoreTokenCacheKey, definition.SignTokenStore+"-"+strconv.FormatUint(uint64(wsId), 10)).Val() != token {
			return
		}
		// 获取门店人员
		ep, _ = ent.Database.Employee.QueryNotDeleted().Where(employee.ID(uint64(wsId))).First(b.ctx)

	}

	return
}

// Assets 物资数据
func (b *warestoreBiz) Assets(am *ent.AssetManager, ep *ent.Employee, req *definition.WarestoreAssetsReq) (res []*definition.WarestoreAssetRes, err error) {
	switch {
	case am != nil && ep == nil:
		// 确认为仓库管理员
		return b.assetsForWarehouse(am.ID, req)
	case am == nil && ep != nil:
		// 确认为门店管理员
		return b.assetsForStore(ep.ID, req)
	default:
		return nil, errors.New(ar.UserNotFound)
	}
}

// assetsForWarehouse 仓库物资数据
func (b *warestoreBiz) assetsForWarehouse(amId uint64, req *definition.WarestoreAssetsReq) (res []*definition.WarestoreAssetRes, err error) {
	// 查询仓库数据
	q := ent.Database.Warehouse.QueryNotDeleted().WithCity()

	if req.WarehouseID != nil {
		q.Where(warehouse.ID(*req.WarehouseID))
	}

	// 查询仓管人员负责的仓库信息
	am, _ := ent.Database.AssetManager.QueryNotDeleted().WithWarehouses().
		Where(
			assetmanager.ID(amId),
			assetmanager.MiniEnable(true),
			assetmanager.HasWarehousesWith(warehouse.DeletedAtIsNil()),
		).First(b.ctx)
	if am != nil && len(am.Edges.Warehouses) != 0 {
		wIds := make([]uint64, 0)
		for _, wh := range am.Edges.Warehouses {
			wIds = append(wIds, wh.ID)
		}
		q.Where(warehouse.IDIn(wIds...))
	}

	whs, err := q.All(b.ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range whs {
		// 查询仓库资产详情
		detail := NewWarehouse().AssetsDetail(item.ID)

		wa := &definition.WarestoreAssetRes{
			ID:     item.ID,
			Name:   item.Name,
			Detail: *detail,
		}
		if item.Edges.City != nil {
			wa.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		res = append(res, wa)
	}

	return
}

// assetsForStore 门店物资数据
func (b *warestoreBiz) assetsForStore(epId uint64, req *definition.WarestoreAssetsReq) (res []*definition.WarestoreAssetRes, err error) {
	// 门店数据
	q := ent.Database.Store.QueryNotDeleted().WithCity()
	if req.StoreID != nil {
		q.Where(store.ID(*req.StoreID))
	}
	// 查询门店人员负责的门店信息 todo 门店集合筛选
	ep, _ := ent.Database.Employee.QueryNotDeleted().WithStores().
		Where(
			employee.ID(epId),
			employee.HasStoresWith(store.DeletedAtIsNil()),
		).First(b.ctx)
	if ep != nil && len(ep.Edges.Stores) != 0 {
		sIds := make([]uint64, 0)
		for _, st := range ep.Edges.Stores {
			sIds = append(sIds, st.ID)
		}
		q.Where(store.IDIn(sIds...))
	}

	sts, err := q.All(b.ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range sts {
		// 查询仓库资产详情
		detail := NewStoreAsset().AssetDetail(item.ID)

		wa := &definition.WarestoreAssetRes{
			ID:     item.ID,
			Name:   item.Name,
			Detail: *detail,
		}
		if item.Edges.City != nil {
			wa.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		res = append(res, wa)
	}

	return
}

// AssetsCommon 物资数据
func (b *warestoreBiz) AssetsCommon(am *ent.AssetManager, ep *ent.Employee, req *definition.WarestoreAssetsCommonReq) *model.PaginationRes {
	q := ent.Database.Asset.QueryNotDeleted().WithCabinet().WithCity().WithStation().WithModel().WithOperator().WithValues().WithStore().WithWarehouse().WithBrand().WithValues()

	b.assetsCommonFilter(am, ep, q, req)

	q.Order(ent.Desc(asset.FieldCreatedAt))
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Asset) *model.AssetListRes {
		return service.NewAsset().DetailForList(item)
	})
}

// assetsCommonFilter 物资数据
func (b *warestoreBiz) assetsCommonFilter(am *ent.AssetManager, ep *ent.Employee, q *ent.AssetQuery, req *definition.WarestoreAssetsCommonReq) {
	switch req.Type {
	case definition.CommonAssetTypeEbike:
		q.Where(asset.Type(model.AssetTypeEbike.Value()))
	case definition.CommonAssetTypeBattery:
		q.Where(asset.TypeIn(model.AssetTypeSmartBattery.Value(), model.AssetTypeNonSmartBattery.Value()))
	}

	if req.WarehouseID != nil {
		q.Where(
			asset.LocationsType(model.AssetLocationsTypeWarehouse.Value()),
			asset.LocationsID(*req.WarehouseID),
		)
	}

	if req.StoreID != nil {
		q.Where(asset.Status(req.Status.Value()))
	}

	if req.Status != nil {
		q.Where(asset.Status(req.Status.Value()))
	}

	if req.ModelID != nil {
		q.Where(asset.ModelID(*req.ModelID))
	}

	if req.BrandID != nil {
		q.Where(asset.BrandID(*req.BrandID))
	}

	if req.BatteryKeyword != nil {
		q.Where(asset.SnContains(*req.BatteryKeyword))
	}

	if req.EbikeKeyword != nil {
		q.Where(asset.SnContains(*req.EbikeKeyword))
	}

	switch {
	case am != nil && ep == nil:
		// 仓库管理查询

		// 查询库管人员配置的仓库数据
		wIds := make([]uint64, 0)
		am, _ = ent.Database.AssetManager.QueryNotDeleted().WithWarehouses().
			Where(
				assetmanager.ID(am.ID),
				assetmanager.HasWarehousesWith(warehouse.DeletedAtIsNil()),
			).First(context.Background())
		if am != nil {
			for _, wh := range am.Edges.Warehouses {
				wIds = append(wIds, wh.ID)
			}
		}
		q.Where(
			asset.LocationsType(model.AssetLocationsTypeWarehouse.Value()),
			asset.LocationsIDIn(wIds...),
		)

	case am == nil && ep != nil:
		// 门店管理查询

		// 查询门店人员配置的门店数据
		sIds := make([]uint64, 0)
		ep, _ = ent.Database.Employee.QueryNotDeleted().WithStores().
			Where(
				employee.ID(ep.ID),
				employee.HasStoresWith(store.DeletedAtIsNil()),
			).First(context.Background())
		if ep != nil {
			for _, st := range ep.Edges.Stores {
				sIds = append(sIds, st.ID)
			}
		}
		q.Where(
			asset.LocationsType(model.AssetLocationsTypeStore.Value()),
			asset.LocationsIDIn(sIds...),
		)

	default:

	}
}
