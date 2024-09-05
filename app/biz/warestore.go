// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package biz

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/LucaTheHacker/go-haversine"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
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
		am, err = ent.Database.AssetManager.QueryNotDeleted().WithRole().
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

// CheckDuty 库管端检查上班范围
func (b *warestoreBiz) CheckDuty(assetSignInfo definition.AssetSignInfo, req *definition.WarestoreDutyReq) (res *definition.WarestoreCheckDutyRes, err error) {
	switch {
	case assetSignInfo.AssetManager != nil:
		// 检查是否可上班
		var wh *ent.Warehouse
		wh, err = b.checkWarehouseDuty(assetSignInfo.AssetManager.ID, req.Sn, req.Lat, req.Lng)
		if err != nil {
			return
		}
		return &definition.WarestoreCheckDutyRes{Name: wh.Name}, nil

	case assetSignInfo.Employee != nil:
		// 检查是否可上班
		var st *ent.Store
		st, err = b.checkStoreDuty(assetSignInfo.Employee.ID, req.Sn, req.Lat, req.Lng)
		if err != nil {
			return
		}
		return &definition.WarestoreCheckDutyRes{Name: st.Name}, nil
	}
	return
}

// Duty 库管端上班
func (b *warestoreBiz) Duty(assetSignInfo definition.AssetSignInfo, req *definition.WarestoreDutyReq) (err error) {
	switch {
	case assetSignInfo.AssetManager != nil:
		// 检查是否可上班
		var wh *ent.Warehouse
		wh, err = b.checkWarehouseDuty(assetSignInfo.AssetManager.ID, req.Sn, req.Lat, req.Lng)
		if err != nil {
			return
		}
		// 上班更新
		err = ent.Database.Warehouse.Update().ClearAssetManagerID().Where(warehouse.AssetManagerID(assetSignInfo.AssetManager.ID)).Exec(b.ctx)
		if err != nil {
			return err
		}

		return wh.Update().SetAssetManagerID(assetSignInfo.AssetManager.ID).Exec(b.ctx)

	case assetSignInfo.Employee != nil:
		// 检查是否可上班
		var st *ent.Store
		st, err = b.checkStoreDuty(assetSignInfo.Employee.ID, req.Sn, req.Lat, req.Lng)
		if err != nil {
			return
		}
		// 上班更新
		err = ent.Database.Store.Update().ClearEmployeeID().Where(store.EmployeeID(assetSignInfo.Employee.ID)).Exec(b.ctx)
		if err != nil {
			return err
		}
		return st.Update().SetEmployeeID(assetSignInfo.Employee.ID).Exec(b.ctx)
	}
	return
}

// checkWarehouseDuty 检查仓库上班信息
func (b *warestoreBiz) checkWarehouseDuty(amId uint64, sn string, lat, lng float64) (wh *ent.Warehouse, err error) {
	wh = NewWarehouse().QuerySn(sn)
	// 判断距离
	if wh == nil || wh.Lat == 0 || wh.Lng == 0 {
		return wh, errors.New("未找到门店地理信息")
	}
	if wh.AssetManagerID != nil {
		return wh, errors.New("当前已有员工上班")
	}

	// 检查当前账号与目标是否绑定关系
	if ewh, _ := ent.Database.Warehouse.QueryNotDeleted().
		Where(
			warehouse.ID(wh.ID),
			warehouse.HasAssetManagersWith(assetmanager.ID(amId)),
		).First(b.ctx); ewh == nil {
		return wh, errors.New("当前用户未拥有该仓库权限")
	}

	distance := haversine.Distance(haversine.NewCoordinates(lat, lng), haversine.NewCoordinates(wh.Lat, wh.Lng))
	meters := distance.Kilometers() * 1000
	if meters > 1000 {
		return wh, errors.New("距离过远")
	}
	return wh, nil
}

// checkStoreDuty 检查门店上班信息
func (b *warestoreBiz) checkStoreDuty(epId uint64, sn string, lat, lng float64) (st *ent.Store, err error) {
	st = service.NewStore().QuerySn(sn)
	if st == nil {
		return st, errors.New("当前门店未找到")
	}
	bc := service.NewBranch().Query(st.BranchID)

	if st.EmployeeID != nil {
		return st, errors.New("当前已有员工上班")
	}
	// 判断距离
	if bc == nil || bc.Lat == 0 || bc.Lng == 0 {
		return st, errors.New("未找到门店地理信息")
	}
	// 检查当前账号与目标是否绑定关系
	if est, _ := ent.Database.Store.QueryNotDeleted().
		Where(
			store.ID(st.ID),
			store.HasEmployeesWith(employee.ID(epId)),
		).First(b.ctx); est == nil {
		return st, errors.New("当前用户未拥有该门店权限")
	}

	distance := haversine.Distance(haversine.NewCoordinates(lat, lng), haversine.NewCoordinates(bc.Lat, bc.Lng))
	meters := distance.Kilometers() * 1000
	if meters > 1000 {
		return st, errors.New("距离过远")
	}
	return st, nil
}

// Profile 仓管资料
func (b *warestoreBiz) Profile(am *ent.AssetManager, ep *ent.Employee, platType definition.PlatType) definition.WarestorePeopleProfile {
	switch {
	case platType == definition.PlatTypeWarehouse && am != nil:
		res := definition.WarestorePeopleProfile{
			ID:       am.ID,
			Phone:    am.Phone,
			Name:     am.Name,
			PlatType: platType,
			RoleName: "仓库管理员",
		}

		v, _ := ent.Database.Warehouse.QueryNotDeleted().
			Where(
				warehouse.HasAssetManagerWith(assetmanager.ID(am.ID)),
			).First(b.ctx)
		if v != nil {
			res.Duty = true
			res.DutyLocationID = v.ID
			res.DutyLocation = "[仓库]" + v.Name
		}

		return res

	case platType == definition.PlatTypeStore && ep != nil:
		res := definition.WarestorePeopleProfile{
			ID:       ep.ID,
			Phone:    ep.Phone,
			Name:     ep.Name,
			PlatType: platType,
			RoleName: "门店店员",
		}
		v, _ := ent.Database.Store.QueryNotDeleted().
			Where(
				store.HasEmployeeWith(employee.ID(ep.ID)),
			).First(b.ctx)
		if v != nil {
			res.Duty = true
			res.DutyLocationID = v.ID
			res.DutyLocation = "[门店]" + v.Name
		}

		return res
	default:
		return definition.WarestorePeopleProfile{}
	}
}

// AssetCount 仓管资产统计
func (b *warestoreBiz) AssetCount(assetSignInfo definition.AssetSignInfo) (res definition.AssetCountRes) {
	switch {
	case assetSignInfo.AssetManager != nil:
		// 确认为仓库管理员
		return b.assetCountForWarehouse(assetSignInfo.AssetManager.ID)
	case assetSignInfo.Employee != nil:
		// 确认为门店管理员
		return b.assetCountForStore(assetSignInfo.Employee.ID)
	case assetSignInfo.Maintainer != nil:
		// 确认为运维人员
		return b.assetCountForMaintainer(assetSignInfo.Maintainer.ID)
	}
	return
}

// assetCountForWarehouse 仓管资产统计
func (b *warestoreBiz) assetCountForWarehouse(id uint64) (res definition.AssetCountRes) {
	// 待接收
	rList, _ := ent.Database.AssetTransfer.QueryNotDeleted().
		WithTransferDetails(func(query *ent.AssetTransferDetailsQuery) { query.WithAsset() }).
		Where(
			assettransfer.ToLocationType(model.AssetLocationsTypeWarehouse.Value()),
			assettransfer.HasToLocationWarehouseWith(warehouse.HasAssetManagersWith(assetmanager.ID(id))),
			assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
		).All(b.ctx)

	for _, v := range rList {
		res.ReceivingCount += 1
		for _, td := range v.Edges.TransferDetails {
			if td.Edges.Asset != nil {
				switch td.Edges.Asset.Type {
				case model.AssetTypeEbikeAccessory.Value():
					res.EbikeAsset.DeliverCount += 1
				case model.AssetTypeSmartBattery.Value():
					res.SmartBatteryAsset.DeliverCount += 1
				case model.AssetTypeNonSmartBattery.Value():
					res.NonSmartBatteryAsset.DeliverCount += 1
				}
			}
		}
	}

	// 配送中
	dList, _ := ent.Database.AssetTransfer.QueryNotDeleted().Where(
		assettransfer.FromLocationType(model.AssetLocationsTypeWarehouse.Value()),
		assettransfer.HasFromLocationWarehouseWith(warehouse.HasAssetManagersWith(assetmanager.ID(id))),
		assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
	).All(b.ctx)

	res.DeliveringCount = len(dList)

	// 异常告警 暂时不做

	// 从资产数据中查询当前账号所属仓库/门店各个类别统计数据

	aList, _ := ent.Database.Asset.QueryNotDeleted().Where(
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
		asset.LocationsType(model.AssetLocationsTypeWarehouse.Value()),
		asset.HasWarehouseWith(warehouse.HasAssetManagersWith(assetmanager.ID(id))),
	).All(b.ctx)
	for _, v := range aList {
		switch v.Type {
		case model.AssetTypeEbikeAccessory.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.EbikeAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.EbikeAsset.FaultCount += 1
			}
		case model.AssetTypeSmartBattery.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.SmartBatteryAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.SmartBatteryAsset.FaultCount += 1
			}
		case model.AssetTypeNonSmartBattery.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.NonSmartBatteryAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.NonSmartBatteryAsset.FaultCount += 1
			}
		}
	}

	// 各个类别合计
	res.EbikeAsset.TotalCount = res.EbikeAsset.StockCount + res.EbikeAsset.DeliverCount + res.EbikeAsset.FaultCount
	res.SmartBatteryAsset.TotalCount = res.SmartBatteryAsset.StockCount + res.SmartBatteryAsset.DeliverCount + res.SmartBatteryAsset.FaultCount
	res.NonSmartBatteryAsset.TotalCount = res.NonSmartBatteryAsset.StockCount + res.NonSmartBatteryAsset.DeliverCount + res.NonSmartBatteryAsset.FaultCount

	return
}

// assetCountForWarehouse 仓管资产统计
func (b *warestoreBiz) assetCountForStore(id uint64) (res definition.AssetCountRes) {
	// 待接收
	rList, _ := ent.Database.AssetTransfer.QueryNotDeleted().
		WithTransferDetails(func(query *ent.AssetTransferDetailsQuery) { query.WithAsset() }).
		Where(
			assettransfer.ToLocationType(model.AssetLocationsTypeStore.Value()),
			assettransfer.HasToLocationStoreWith(store.HasEmployeesWith(employee.ID(id))),
			assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
		).All(b.ctx)

	for _, v := range rList {
		res.ReceivingCount += 1
		for _, td := range v.Edges.TransferDetails {
			if td.Edges.Asset != nil {
				switch td.Edges.Asset.Type {
				case model.AssetTypeEbikeAccessory.Value():
					res.EbikeAsset.DeliverCount += 1
				case model.AssetTypeSmartBattery.Value():
					res.SmartBatteryAsset.DeliverCount += 1
				case model.AssetTypeNonSmartBattery.Value():
					res.NonSmartBatteryAsset.DeliverCount += 1
				}
			}
		}
	}

	// 配送中
	dList, _ := ent.Database.AssetTransfer.QueryNotDeleted().Where(
		assettransfer.FromLocationType(model.AssetLocationsTypeStore.Value()),
		assettransfer.HasFromLocationStoreWith(store.HasEmployeesWith(employee.ID(id))),
		assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
	).All(b.ctx)

	res.DeliveringCount = len(dList)

	// 异常告警 暂时不做

	// 从资产数据中查询当前账号所属仓库/门店各个类别统计数据

	aList, _ := ent.Database.Asset.QueryNotDeleted().Where(
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
		asset.LocationsType(model.AssetLocationsTypeStore.Value()),
		asset.HasStoreWith(store.HasEmployeesWith(employee.ID(id))),
	).All(b.ctx)
	for _, v := range aList {
		switch v.Type {
		case model.AssetTypeEbikeAccessory.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.EbikeAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.EbikeAsset.FaultCount += 1
			}
		case model.AssetTypeSmartBattery.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.SmartBatteryAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.SmartBatteryAsset.FaultCount += 1
			}
		case model.AssetTypeNonSmartBattery.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.NonSmartBatteryAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.NonSmartBatteryAsset.FaultCount += 1
			}
		}
	}

	// 各个类别合计
	res.EbikeAsset.TotalCount = res.EbikeAsset.StockCount + res.EbikeAsset.DeliverCount + res.EbikeAsset.FaultCount
	res.SmartBatteryAsset.TotalCount = res.SmartBatteryAsset.StockCount + res.SmartBatteryAsset.DeliverCount + res.SmartBatteryAsset.FaultCount
	res.NonSmartBatteryAsset.TotalCount = res.NonSmartBatteryAsset.StockCount + res.NonSmartBatteryAsset.DeliverCount + res.NonSmartBatteryAsset.FaultCount

	return
}

// assetCountForWarehouse 仓管资产统计
func (b *warestoreBiz) assetCountForMaintainer(id uint64) (res definition.AssetCountRes) {
	// 待接收
	rList, _ := ent.Database.AssetTransfer.QueryNotDeleted().
		WithTransferDetails(func(query *ent.AssetTransferDetailsQuery) { query.WithAsset() }).
		Where(
			assettransfer.ToLocationType(model.AssetLocationsTypeOperation.Value()),
			assettransfer.HasToLocationOperatorWith(maintainer.ID(id)),
			assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
		).All(b.ctx)

	for _, v := range rList {
		res.ReceivingCount += 1
		for _, td := range v.Edges.TransferDetails {
			if td.Edges.Asset != nil {
				switch td.Edges.Asset.Type {
				case model.AssetTypeEbikeAccessory.Value():
					res.EbikeAsset.DeliverCount += 1
				case model.AssetTypeSmartBattery.Value():
					res.SmartBatteryAsset.DeliverCount += 1
				case model.AssetTypeNonSmartBattery.Value():
					res.NonSmartBatteryAsset.DeliverCount += 1
				}
			}
		}
	}

	// 配送中
	dList, _ := ent.Database.AssetTransfer.QueryNotDeleted().Where(
		assettransfer.FromLocationType(model.AssetLocationsTypeOperation.Value()),
		assettransfer.HasFromLocationOperatorWith(maintainer.ID(id)),
		assettransfer.Status(model.AssetTransferStatusDelivering.Value()),
	).All(b.ctx)

	res.DeliveringCount = len(dList)

	// 异常告警 暂时不做

	// 从资产数据中查询当前账号所属仓库/门店各个类别统计数据

	aList, _ := ent.Database.Asset.QueryNotDeleted().Where(
		asset.StatusIn(model.AssetStatusStock.Value(), model.AssetStatusFault.Value()),
		asset.LocationsType(model.AssetLocationsTypeOperation.Value()),
		asset.LocationsID(id),
	).All(b.ctx)
	for _, v := range aList {
		switch v.Type {
		case model.AssetTypeEbikeAccessory.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.EbikeAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.EbikeAsset.FaultCount += 1
			}
		case model.AssetTypeSmartBattery.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.SmartBatteryAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.SmartBatteryAsset.FaultCount += 1
			}
		case model.AssetTypeNonSmartBattery.Value():
			switch v.Status {
			case model.AssetStatusStock.Value():
				res.NonSmartBatteryAsset.StockCount += 1
			case model.AssetStatusFault.Value():
				res.NonSmartBatteryAsset.FaultCount += 1
			}
		}
	}

	// 各个类别合计
	res.EbikeAsset.TotalCount = res.EbikeAsset.StockCount + res.EbikeAsset.DeliverCount + res.EbikeAsset.FaultCount
	res.SmartBatteryAsset.TotalCount = res.SmartBatteryAsset.StockCount + res.SmartBatteryAsset.DeliverCount + res.SmartBatteryAsset.FaultCount
	res.NonSmartBatteryAsset.TotalCount = res.NonSmartBatteryAsset.StockCount + res.NonSmartBatteryAsset.DeliverCount + res.NonSmartBatteryAsset.FaultCount

	return
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
func (b *warestoreBiz) Assets(assetSignInfo definition.AssetSignInfo, req *definition.WarestoreAssetsReq) *model.PaginationRes {
	switch {
	case assetSignInfo.AssetManager != nil:
		// 确认为仓库管理员
		return b.assetsForWarehouse(assetSignInfo.AssetManager.ID, req)
	case assetSignInfo.Employee != nil:
		// 确认为门店管理员
		return b.assetsForStore(assetSignInfo.Employee.ID, req)
	case assetSignInfo.Agent != nil:
		// 确认为代理员
		return b.assetsForAgent(assetSignInfo.Agent.ID, req)
	case assetSignInfo.Maintainer != nil:
		// 确认为运维
		return b.assetsForMaintainer(assetSignInfo.Maintainer.ID, req)
	}

	return nil
}

// assetsForWarehouse 仓库物资数据
func (b *warestoreBiz) assetsForWarehouse(amId uint64, req *definition.WarestoreAssetsReq) *model.PaginationRes {
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

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Warehouse) *definition.WarestoreAssetRes {
		// 查询仓库资产详情
		res := &definition.WarestoreAssetRes{
			ID:     item.ID,
			Name:   item.Name,
			Detail: b.assetTotal(item.ID, model.AssetLocationsTypeWarehouse),
		}
		if item.Edges.City != nil {
			res.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		return res
	})

}

// assetsForStore 门店物资数据
func (b *warestoreBiz) assetsForStore(epId uint64, req *definition.WarestoreAssetsReq) *model.PaginationRes {
	// 门店数据
	q := ent.Database.Store.QueryNotDeleted().WithCity()
	if req.StoreID != nil {
		q.Where(store.ID(*req.StoreID))
	}

	ep, _ := ent.Database.Employee.QueryNotDeleted().WithStores().WithGroup().
		Where(
			employee.ID(epId),
			employee.HasStoresWith(store.DeletedAtIsNil()),
		).First(b.ctx)
	if ep != nil {
		// 判断是配置的门店集合还是门店数据
		if ep.Edges.Group != nil {
			q.Where(store.GroupID(ep.Edges.Group.ID))
		} else {
			sIds := make([]uint64, 0)
			for _, st := range ep.Edges.Stores {
				sIds = append(sIds, st.ID)
			}
			q.Where(store.IDIn(sIds...))
		}

	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Store) *definition.WarestoreAssetRes {
		// 查询仓库资产详情
		res := &definition.WarestoreAssetRes{
			ID:     item.ID,
			Name:   item.Name,
			Detail: b.assetTotal(item.ID, model.AssetLocationsTypeStore),
		}
		if item.Edges.City != nil {
			res.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		return res
	})

}

// assetsForAgent 代理物资数据
func (b *warestoreBiz) assetsForAgent(agId uint64, req *definition.WarestoreAssetsReq) *model.PaginationRes {
	// 门店数据
	q := ent.Database.EnterpriseStation.QueryNotDeleted().WithCity()
	if req.StationID != nil {
		q.Where(enterprisestation.ID(*req.StoreID))
	}

	ag, _ := ent.Database.Agent.QueryNotDeleted().WithStations().
		Where(
			agent.ID(agId),
			agent.HasStationsWith(enterprisestation.DeletedAtIsNil()),
		).First(b.ctx)
	if ag != nil {
		// 判断是配置的站点数据
		sIds := make([]uint64, 0)
		for _, st := range ag.Edges.Stations {
			sIds = append(sIds, st.ID)
		}
		q.Where(enterprisestation.IDIn(sIds...))
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.EnterpriseStation) *definition.WarestoreAssetRes {
		// 查询仓库资产详情
		res := &definition.WarestoreAssetRes{
			ID:     item.ID,
			Name:   item.Name,
			Detail: b.assetTotal(item.ID, model.AssetLocationsTypeStation),
		}
		if item.Edges.City != nil {
			res.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		return res
	})

}

// assetsForAgent 代理物资数据
func (b *warestoreBiz) assetsForMaintainer(mtId uint64, req *definition.WarestoreAssetsReq) *model.PaginationRes {
	// 数据
	q := ent.Database.Maintainer.Query().Where(maintainer.ID(mtId))

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Maintainer) *definition.WarestoreAssetRes {
		// 查询仓库资产详情
		res := &definition.WarestoreAssetRes{
			ID:     item.ID,
			Name:   item.Name,
			Detail: b.assetTotal(item.ID, model.AssetLocationsTypeOperation),
		}
		return res
	})

}

// assetTotal 物资数据统计
func (b *warestoreBiz) assetTotal(id uint64, aType model.AssetLocationsType) (res definition.WarestoreAssetDetail) {
	// 查询仓库所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.LocationsType(aType.Value()),
			asset.LocationsIDIn(id),
			asset.Status(model.AssetStatusStock.Value()),
		).WithModel().WithBrand().WithMaterial()

	list, _ := q.All(b.ctx)

	res = b.commonAssetTotalDetail(list)

	return
}

// wAssetTotal 仓库物资数据统计
func (b *warestoreBiz) commonAssetTotalDetail(list []*ent.Asset) (res definition.WarestoreAssetDetail) {
	// 查询统计资产数据
	res = definition.WarestoreAssetDetail{
		Ebikes:             make([]*definition.WarestoreMaterial, 0),
		SmartBatteries:     make([]*definition.WarestoreMaterial, 0),
		NonSmartBatteries:  make([]*definition.WarestoreMaterial, 0),
		CabinetAccessories: make([]*definition.WarestoreMaterial, 0),
		EbikeAccessories:   make([]*definition.WarestoreMaterial, 0),
		OtherAssets:        make([]*definition.WarestoreMaterial, 0),
	}

	ebikeBrandMap := make(map[string]*definition.WarestoreMaterial)
	sBatModelMap := make(map[string]*definition.WarestoreMaterial)
	nSBatModelMap := make(map[string]*definition.WarestoreMaterial)
	cabAccNameMap := make(map[string]*definition.WarestoreMaterial)
	ebikeAccNameMap := make(map[string]*definition.WarestoreMaterial)
	otherNameMap := make(map[string]*definition.WarestoreMaterial)

	for _, v := range list {
		switch {
		case v.Type == model.AssetTypeEbike.Value() && v.Edges.Brand != nil:
			res.EbikeTotal += 1
			if ebikeBrandMap[v.Edges.Brand.Name] != nil {
				ebikeBrandMap[v.Edges.Brand.Name].Num += 1
			} else {
				ebikeBrandMap[v.Edges.Brand.Name] = &definition.WarestoreMaterial{
					ID:   v.Edges.Brand.ID,
					Name: v.Edges.Brand.Name,
					Num:  1,
				}
			}

		case v.Type == model.AssetTypeSmartBattery.Value() && v.Edges.Model != nil:
			res.SmartBatteryTotal += 1
			if sBatModelMap[v.Edges.Model.Model] != nil {
				sBatModelMap[v.Edges.Model.Model].Num += 1
			} else {
				sBatModelMap[v.Edges.Model.Model] = &definition.WarestoreMaterial{
					ID:   v.Edges.Model.ID,
					Name: v.Edges.Model.Model,
					Num:  1,
				}
			}
		case v.Type == model.AssetTypeNonSmartBattery.Value() && v.Edges.Model != nil:
			res.NonSmartBatteryTotal += 1
			if nSBatModelMap[v.Edges.Model.Model] != nil {
				nSBatModelMap[v.Edges.Model.Model].Num += 1
			} else {
				nSBatModelMap[v.Edges.Model.Model] = &definition.WarestoreMaterial{
					ID:   v.Edges.Model.ID,
					Name: v.Edges.Model.Model,
					Num:  1,
				}
			}
		case v.Type == model.AssetTypeCabinetAccessory.Value() && v.Edges.Material != nil:
			res.CabinetAccessoryTotal += 1
			if cabAccNameMap[v.Edges.Material.Name] != nil {
				cabAccNameMap[v.Edges.Material.Name].Num += 1
			} else {
				cabAccNameMap[v.Edges.Material.Name] = &definition.WarestoreMaterial{
					ID:   v.Edges.Material.ID,
					Name: v.Edges.Material.Name,
					Num:  1,
				}
			}
		case v.Type == model.AssetTypeEbikeAccessory.Value() && v.Edges.Material != nil:
			res.EbikeAccessoryTotal += 1
			if ebikeAccNameMap[v.Edges.Material.Name] != nil {
				ebikeAccNameMap[v.Edges.Material.Name].Num += 1
			} else {
				ebikeAccNameMap[v.Edges.Material.Name] = &definition.WarestoreMaterial{
					ID:   v.Edges.Material.ID,
					Name: v.Edges.Material.Name,
					Num:  1,
				}
			}
		case v.Type == model.AssetTypeOtherAccessory.Value() && v.Edges.Material != nil:
			res.OtherAssetTotal += 1
			if otherNameMap[v.Edges.Material.Name] != nil {
				otherNameMap[v.Edges.Material.Name].Num += 1
			} else {
				otherNameMap[v.Edges.Material.Name] = &definition.WarestoreMaterial{
					ID:   v.Edges.Material.ID,
					Name: v.Edges.Material.Name,
					Num:  1,
				}
			}
		}
	}

	// 组装数据
	for _, v := range ebikeBrandMap {
		res.Ebikes = append(res.Ebikes, v)
	}
	for _, v := range sBatModelMap {
		res.SmartBatteries = append(res.SmartBatteries, v)
	}
	for _, v := range nSBatModelMap {
		res.NonSmartBatteries = append(res.NonSmartBatteries, v)
	}
	for _, v := range cabAccNameMap {
		res.CabinetAccessories = append(res.CabinetAccessories, v)
	}
	for _, v := range ebikeAccNameMap {
		res.EbikeAccessories = append(res.EbikeAccessories, v)
	}
	for _, v := range otherNameMap {
		res.OtherAssets = append(res.OtherAssets, v)
	}

	// 排序
	sort.Slice(res.Ebikes, func(i, j int) bool {
		return strings.Compare(res.Ebikes[i].Name, res.Ebikes[j].Name) < 0
	})
	sort.Slice(res.SmartBatteries, func(i, j int) bool {
		return strings.Compare(res.SmartBatteries[i].Name, res.SmartBatteries[j].Name) < 0
	})
	sort.Slice(res.NonSmartBatteries, func(i, j int) bool {
		return strings.Compare(res.NonSmartBatteries[i].Name, res.NonSmartBatteries[j].Name) < 0
	})
	sort.Slice(res.CabinetAccessories, func(i, j int) bool {
		return strings.Compare(res.CabinetAccessories[i].Name, res.CabinetAccessories[j].Name) < 0
	})
	sort.Slice(res.EbikeAccessories, func(i, j int) bool {
		return strings.Compare(res.EbikeAccessories[i].Name, res.EbikeAccessories[j].Name) < 0
	})
	sort.Slice(res.OtherAssets, func(i, j int) bool {
		return strings.Compare(res.OtherAssets[i].Name, res.OtherAssets[j].Name) < 0
	})
	return
}

// AssetsCommon 物资数据
func (b *warestoreBiz) AssetsCommon(assetSignInfo definition.AssetSignInfo, req *definition.WarestoreAssetsCommonReq) *model.PaginationRes {
	q := ent.Database.Asset.QueryNotDeleted().WithCabinet().WithCity().WithStation().WithModel().WithOperator().WithValues().WithStore().WithWarehouse().WithBrand().WithValues()

	b.assetsCommonFilter(assetSignInfo, q, req)

	q.Order(ent.Desc(asset.FieldCreatedAt))
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Asset) *model.AssetListRes {
		return service.NewAsset().DetailForList(item)
	})
}

// assetsCommonFilter 物资数据
func (b *warestoreBiz) assetsCommonFilter(assetSignInfo definition.AssetSignInfo, q *ent.AssetQuery, req *definition.WarestoreAssetsCommonReq) {
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
	case assetSignInfo.AssetManager != nil:
		// 仓库管理查询

		// 查询库管人员配置的仓库数据
		wIds := make([]uint64, 0)
		am, _ := ent.Database.AssetManager.QueryNotDeleted().WithWarehouses().
			Where(
				assetmanager.ID(assetSignInfo.AssetManager.ID),
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
	case assetSignInfo.Employee != nil:
		// 门店管理查询

		// 查询门店人员配置的门店数据
		sIds := make([]uint64, 0)
		ep, _ := ent.Database.Employee.QueryNotDeleted().WithStores().
			Where(
				employee.ID(assetSignInfo.Employee.ID),
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
	case assetSignInfo.Agent != nil:
		// 查询代理人员配置的代理站点
		ids := make([]uint64, 0)
		ag, _ := ent.Database.Agent.QueryNotDeleted().WithStations().
			Where(
				agent.ID(assetSignInfo.Agent.ID),
				agent.HasStationsWith(enterprisestation.DeletedAtIsNil()),
			).First(context.Background())
		if ag != nil {
			for _, v := range ag.Edges.Stations {
				ids = append(ids, v.ID)
			}
		}
		q.Where(
			asset.LocationsType(model.AssetLocationsTypeStation.Value()),
			asset.LocationsIDIn(ids...),
		)
	case assetSignInfo.Maintainer != nil:
		q.Where(
			asset.LocationsType(model.AssetLocationsTypeOperation.Value()),
			asset.LocationsID(assetSignInfo.Maintainer.ID),
		)
	default:
	}
}

// Modify 修改资产
func (b *warestoreBiz) Modify(assetSignInfo definition.AssetSignInfo, req *model.AssetModifyReq) error {
	var md model.Modifier
	if assetSignInfo.AssetManager != nil {
		md = model.Modifier{
			ID:    assetSignInfo.AssetManager.ID,
			Name:  assetSignInfo.AssetManager.Name,
			Phone: assetSignInfo.AssetManager.Phone,
		}
	}
	if assetSignInfo.Employee != nil {
		md = model.Modifier{
			ID:    assetSignInfo.Employee.ID,
			Name:  assetSignInfo.Employee.Name,
			Phone: assetSignInfo.Employee.Phone,
		}
	}
	if assetSignInfo.Agent != nil {
		md = model.Modifier{
			ID:    assetSignInfo.Agent.ID,
			Name:  assetSignInfo.Agent.Name,
			Phone: assetSignInfo.Agent.Phone,
		}
	}
	if assetSignInfo.Maintainer != nil {
		md = model.Modifier{
			ID:    assetSignInfo.Maintainer.ID,
			Name:  assetSignInfo.Maintainer.Name,
			Phone: assetSignInfo.Maintainer.Phone,
		}
	}
	return service.NewAsset().Modify(b.ctx, req, &md)
}
