// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-24
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"errors"
	"strconv"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/rpc/pb"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/silk"
)

type batteryService struct {
	*BaseService
	orm *ent.AssetClient
}

func NewBattery(params ...any) *batteryService {
	return &batteryService{
		BaseService: newService(params...),
		orm:         ent.Database.Asset,
	}
}

// QuerySn 查询电池
func (s *batteryService) QuerySn(sn string) (bat *ent.Asset, err error) {
	_, err = adapter.ParseBatterySN(sn)
	if err != nil {
		zap.L().Error("查询电池失败，电池编码错误: "+sn, zap.Error(err))
		return
	}

	bat, _ = NewAsset().QuerySn(sn)
	if bat == nil {
		zap.L().Error("查询电池失败，未找到电池: " + sn)
		return nil, errors.New("未找到电池: " + sn)
	}
	return
}

// LoadOrCreate 加载电池, 若电池不存在则先创建电池, 若电池存在, 则不更新电池直接返回
func (s *batteryService) LoadOrCreate(sn string, params ...any) (bat *ent.Asset, err error) {
	bat, _ = NewAsset().QuerySn(sn)
	if bat != nil {
		return
	}

	var (
		cabID   *uint64
		ordinal *int
	)

	for _, param := range params {
		switch v := param.(type) {
		case *model.BatteryInCabinet:
			cabID = silk.UInt64(v.CabinetID)
			ordinal = silk.Int(v.Ordinal)
		}
	}
	cab, _ := ent.Database.Cabinet.QueryNotDeleted().Where(cabinet.ID(*cabID)).First(s.ctx)
	if cab == nil {
		zap.L().Error("未找到电柜: " + strconv.Itoa(int(*cabID)))
		return nil, errors.New("未找到电柜: " + strconv.Itoa(int(*cabID)))
	}

	// s.orm.Create().SetModel(ab.Model).SetSn(sn).SetBrand(ab.Brand).SetNillableCabinetID(cabID).SetNillableOrdinal(ordinal).Save(s.ctx)
	// 创建资产
	err = NewAsset().Create(s.ctx, &model.AssetCreateReq{
		AssetType:     model.AssetTypeSmartBattery,
		CityID:        cab.CityID,
		SN:            silk.String(sn),
		LocationsType: model.AssetLocationsTypeCabinet,
		LocationsID:   cab.ID,
		Enable:        silk.Bool(true),
	}, &model.Modifier{
		ID:    cab.ID,
		Name:  cab.Name,
		Phone: cab.Serial,
	})
	if err != nil {
		return nil, err
	}
	bat, err = NewAsset().QuerySn(sn)
	if err != nil {
		return nil, err
	}
	err = bat.Update().SetNillableOrdinal(ordinal).Exec(s.ctx)
	if err != nil {
		return nil, err
	}
	return bat, nil
}

// Create 创建电池
// func (s *batteryService) Create(req *model.BatteryCreateReq) {
// 	enable := true
// 	if req.Enable != nil {
// 		enable = *req.Enable
// 	}
//
// 	// 解析电池编号
// 	ab, err := adapter.ParseBatterySN(req.SN)
// 	if err != nil || ab.Model == "" {
// 		snag.Panic("电池编号解析失败, 请擦亮你的双眼")
// 	}
// 	_, err = s.orm.Create().
// 		SetSn(req.SN).
// 		SetBrand(ab.Brand).
// 		SetModel(ab.Model).
// 		SetEnable(enable).
// 		SetCityID(req.CityID).
// 		Save(s.ctx)
// 	if err != nil {
// 		snag.Panic("电池创建失败: " + err.Error())
// 	}
// }

// BatchCreate 批量创建电池
// 0-城市:city 1-型号:model 2-编号:sn
// func (s *batteryService) BatchCreate(c echo.Context) []string {
// 	rows, sns, failed, err := s.BaseService.GetXlsxRows(c, 2, 2, 2)
// 	if err != nil {
// 		snag.Panic(err)
// 	}
//
// 	// 查重
// 	items, _ := s.orm.Query().Where(battery.SnIn(sns...)).All(s.ctx)
// 	m := make(map[string]bool)
// 	for _, item := range items {
// 		m[item.Sn] = true
// 	}
//
// 	// 查询城市
// 	cs := make(map[string]struct{})
// 	for _, row := range rows {
// 		cs[row[0]] = struct{}{}
// 	}
// 	var (
// 		cids []string
// 		cm   = make(map[string]uint64)
// 	)
// 	for k := range cs {
// 		cids = append(cids, k)
// 	}
// 	cities, _ := ent.Database.City.Query().Where(city.NameIn(cids...)).All(s.ctx)
// 	for _, ci := range cities {
// 		cm[ci.Name] = ci.ID
// 	}
//
// 	for _, row := range rows {
// 		sn := row[1]
// 		if m[sn] {
// 			failed = append(failed, fmt.Sprintf("编号%s已存在", sn))
// 			continue
// 		}
//
// 		// 解析电池编号
// 		ab, err := adapter.ParseBatterySN(sn)
// 		if err != nil || ab.Model == "" {
// 			failed = append(failed, fmt.Sprintf("电池编号%s解析失败, 请擦亮你的双眼", sn))
// 			continue
// 		}
//
// 		creator := s.orm.Create()
//
// 		// 城市
// 		if cid, ok := cm[row[0]]; ok {
// 			creator.SetCityID(cid)
// 		} else {
// 			failed = append(failed, fmt.Sprintf("城市%s查询失败", row[0]))
// 			continue
// 		}
//
// 		_, err = creator.SetModel(ab.Model).SetBrand(ab.Brand).SetSn(sn).Save(s.ctx)
// 		if err != nil {
// 			failed = append(failed, fmt.Sprintf("%s保存失败: %v", sn, err))
// 		}
//
// 	}
//
// 	return failed
// }

//	func (s *batteryService) Modify(req *model.BatteryModifyReq) {
//		// 查找电池
//		b := s.QueryIDX(req.ID)
//		u := b.Update()
//		if req.Enable != nil {
//			u.SetEnable(*req.Enable)
//		}
//		if req.CityID != nil {
//			u.SetCityID(*req.CityID)
//		}
//		_, _ = u.Save(s.ctx)
//	}
//
// func (s *batteryService) listFilter(req model.BatteryFilter) (q *ent.BatteryQuery, info ar.Map) {
// 	q = s.orm.Query().WithRider().WithCity().WithCabinet().WithStation().WithEnterprise().Order(ent.Desc(battery.FieldCreatedAt))
// 	info = make(ar.Map)
//
// 	var (
// 		status     = 1
// 		statusText = map[int]string{
// 			0: "全部",
// 			1: "启用",
// 			2: "禁用",
// 		}
// 	)
// 	if req.Status != nil {
// 		status = *req.Status
// 	}
// 	info["状态"] = statusText[status]
// 	switch status {
// 	case 1:
// 		q.Where(asset.Enable(true))
// 	case 2:
// 		q.Where(battery.Enable(false))
// 	default:
// 		info["状态"] = "-"
// 	}
//
// 	if req.CityID != 0 {
// 		info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
// 		q.Where(battery.CityID(req.CityID))
// 	}
//
// 	if req.SN != "" {
// 		info["编号"] = req.SN
// 		q.Where(battery.SnContainsFold(req.SN))
// 	}
//
// 	if req.Model != "" {
// 		m := strings.ToUpper(req.Model)
// 		info["型号"] = m
// 		q.Where(battery.Model(m))
// 	}
// 	if req.EnterpriseID != nil {
// 		info["团签"] = ent.NewExportInfo(*req.EnterpriseID, enterprise.Table)
// 		q.Where(battery.EnterpriseID(*req.EnterpriseID))
// 	}
// 	if req.StationID != nil {
// 		info["站点"] = ent.NewExportInfo(*req.StationID, enterprisestation.Table)
// 		q.Where(battery.StationID(*req.StationID))
// 	}
// 	if req.CabinetName != nil {
// 		info["电柜名称"] = *req.CabinetName
// 		q.Where(battery.HasCabinetWith(cabinet.Name(*req.CabinetName)))
// 	}
// 	if req.Keyword != nil {
// 		info["关键字"] = *req.Keyword
// 		q.Where(battery.HasRiderWith(rider.NameContainsFold(*req.Keyword)))
// 	}
//
// 	if req.OwnerType != nil {
// 		var name string
// 		switch *req.OwnerType {
// 		case 1: // 平台
// 			q.Where(battery.EnterpriseIDIsNil())
// 			name = "平台"
// 		case 2: // 代理商
// 			q.Where(battery.EnterpriseIDNotNil())
// 			name = "代理商"
// 		}
// 		info["归属"] = name
// 	}
//
// 	if req.CabinetID != nil {
// 		info["电柜"] = ent.NewExportInfo(*req.CabinetID, cabinet.Table)
// 		q.Where(battery.CabinetID(*req.CabinetID))
// 	}
//
// 	if req.RiderID != nil {
// 		info["骑手"] = ent.NewExportInfo(*req.RiderID, rider.Table)
// 		q.Where(battery.RiderID(*req.RiderID))
// 	}
//
// 	switch req.Goal {
// 	case model.BatteryStation:
// 		info["查询目标"] = "站点"
// 		q.Where(
// 			battery.StationIDNotNil(),
// 			battery.CabinetIDIsNil(),
// 			battery.RiderIDIsNil(),
// 		)
// 	case model.BatteryCabinet:
// 		info["查询目标"] = "电柜"
// 		q.Where(
// 			battery.CabinetIDNotNil(),
// 			battery.RiderIDIsNil(),
// 		)
// 	case model.BatteryRider:
// 		info["查询目标"] = "骑手"
// 		q.Where(
// 			battery.RiderIDNotNil(),
// 			battery.CabinetIDIsNil(),
// 		)
// 	}
//
// 	return
// }
//
// func (s *batteryService) List(req *model.BatteryListReq) (res *model.PaginationRes) {
// 	q, _ := s.listFilter(req.BatteryFilter)
// 	snmap := make(map[adapter.BatteryBrand][]string)
// 	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Battery) (res *model.BatteryListRes) {
// 		snmap[item.Brand] = append(snmap[item.Brand], item.Sn)
// 		res = &model.BatteryListRes{
// 			ID:     item.ID,
// 			Brand:  item.Brand,
// 			Model:  item.Model,
// 			Enable: item.Enable,
// 			SN:     item.Sn,
// 		}
//
// 		c := item.Edges.City
// 		if c != nil {
// 			res.City = &model.City{
// 				ID:   c.ID,
// 				Name: c.Name,
// 			}
// 		}
//
// 		r := item.Edges.Rider
// 		if r != nil {
// 			res.Rider = &model.Rider{
// 				ID:    r.ID,
// 				Phone: r.Phone,
// 				Name:  r.Name,
// 			}
// 		}
//
// 		cab := item.Edges.Cabinet
// 		if cab != nil {
// 			res.Cabinet = &model.CabinetBasicInfo{
// 				ID:     cab.ID,
// 				Brand:  cab.Brand,
// 				Serial: cab.Serial,
// 				Name:   cab.Name,
// 			}
// 		}
// 		if item.Edges.Station != nil {
// 			res.StationName = item.Edges.Station.Name
// 		}
// 		if item.Edges.Enterprise != nil {
// 			res.EnterpriseName = item.Edges.Enterprise.Name
// 		}
// 		return
// 	})
//
// 	// 请求bms rpc
// 	result := make(map[string]*pb.BatteryItem)
// 	for br, list := range snmap {
// 		r := rpc.BmsBatch(br, &pb.BatteryBatchRequest{Sn: list})
// 		if r == nil {
// 			continue
// 		}
//
// 		for _, rb := range r.Items {
// 			result[rb.Sn] = rb
// 		}
// 	}
//
// 	for _, data := range res.Items.([]*model.BatteryListRes) {
// 		if rb, ok := result[data.SN]; ok {
// 			if len(rb.Heartbeats) > 0 {
// 				data.BmsBattery = model.NewBmsBattery(rb.Heartbeats[0])
// 			}
// 		}
// 	}
//
// 	return
// }

// RiderBusiness 骑手业务操作电池
func (s *batteryService) RiderBusiness(putin bool, sn string, r *model.Rider, cab *ent.Cabinet, ordinal int) {
	var before, after string
	var op model.Operate

	target := "电池: " + sn + ", 电柜: " + cab.Serial + ", " + strconv.Itoa(ordinal) + " 号仓"

	if putin {
		op = model.OperateRiderPutin
		after = target
	} else {
		op = model.OperateRiderPutout
		before = target
	}

	go logging.NewOperateLog().
		SetOperate(op).
		SetCabinet(&model.CabinetBasicInfo{
			ID:     cab.ID,
			Brand:  cab.Brand,
			Serial: cab.Serial,
			Name:   cab.Name,
		}).
		SetRefManually(rider.Table, r.ID).
		SetDiff(before, after).
		Send()
}

// RiderDetail 获取电池详情
func (s *batteryService) RiderDetail(riderID uint64) (res model.BatteryDetail) {
	bat, _ := NewAsset().QueryRiderID(riderID)
	if bat != nil {
		// 获取电量信息
		var soc float64 = 0
		sr := rpc.BmsSample(adapter.BatteryBrand(bat.BrandName), &pb.BatterySnRequest{Sn: bat.Sn})
		if sr != nil {
			soc = float64(sr.Soc)
		}

		var modelName string
		if bat.Edges.Model != nil {
			modelName = bat.Edges.Model.Model
		}
		res = model.BatteryDetail{
			ID:    bat.ID,
			Model: modelName,
			SN:    bat.Sn,
			Soc:   soc,
		}
	}
	return
}

// BindRequest 电池绑定至骑手
func (s *batteryService) BindRequest(req *model.BatteryBind) error {
	// 查找订阅
	sub := NewSubscribe().QueryEffectiveX(req.RiderID, ent.SubscribeQueryWithRider)

	// 查找电池
	bat, _ := NewAsset().QueryID(req.BatteryID)
	if bat == nil {
		return errors.New("未找到电池")
	}

	// 查看是否冲突
	if bat.LocationsType == model.AssetLocationsTypeRider.Value() && bat.LocationsID != 0 || bat.SubscribeID != nil && *bat.SubscribeID != sub.ID {
		return errors.New("当前电池有绑定中的骑手, 无法重复绑定")
	}

	if bat.LocationsType == model.AssetLocationsTypeCabinet.Value() && bat.LocationsID != 0 {
		return errors.New("电柜中的电池无法手动绑定骑手")
	}
	err := s.Bind(bat, sub, sub.Edges.Rider)
	if err != nil {
		return err
	}
	return nil
}

// Bind 绑定电池
func (s *batteryService) Bind(bat *ent.Asset, sub *ent.Subscribe, rd *ent.Rider) error {
	err := s.Allocate(bat, sub, model.AssetTransferTypeTransfer)
	if err != nil {
		return errors.New("绑定失败" + err.Error())
	}
	go logging.NewOperateLog().
		SetOperate(model.OperateBindBattery).
		SetRef(rd).
		SetDiff("", "新电池: "+bat.Sn).
		SetModifier(s.modifier).
		Send()
	return nil
}

// Unbind 解绑电池
func (s *batteryService) Unbind(req *model.BatteryUnbindRequest) error {
	// 查找订阅
	sub := NewSubscribe().QueryEffectiveX(req.RiderID, ent.SubscribeQueryWithRider)
	if sub == nil {
		return errors.New("未找到订阅")
	}
	r, _ := sub.QueryRider().WithBattery(func(query *ent.AssetQuery) {
		query.Where(asset.TypeIn(model.AssetTypeSmartBattery.Value(), model.AssetTypeNonSmartBattery.Value()))
	}).First(s.ctx)
	if r == nil {
		return errors.New("未找到骑手")
	}
	bat := r.Edges.Battery
	if bat == nil {
		return errors.New("未找到绑定的电池")
	}
	err := s.Unallocate(bat, req.ToLocationType, req.ToLocationID, model.AssetTransferTypeTransfer)
	if err != nil {
		return errors.New("解绑失败" + err.Error())
	}
	go logging.NewOperateLog().
		SetOperate(model.OperateUnbindBattery).
		SetDiff("旧电池: "+bat.Sn, "无电池").
		SetModifier(s.modifier).
		SetRef(sub.Edges.Rider).
		Send()

	return nil
}

// Allocate 将电池分配给骑手
// @param createFlowIfError 即使遇到错误也创建电池流转信息
// func (s *batteryService) Allocate(tx *ent.Tx, bat *ent.Battery, sub *ent.Subscribe, createFlowIfError bool) (err error) {
// 	// 删除原有骑手电池信息
// 	err = tx.Battery.Update().Where(battery.SubscribeID(sub.ID)).ClearRiderID().ClearSubscribeID().Exec(s.ctx)
// 	if err != nil {
// 		return
// 	}
//
// 	// 分配电池给骑手
// 	updater := tx.Battery.UpdateOne(bat).SetRiderID(sub.RiderID).SetSubscribeID(sub.ID)
//
//
//
// 	return
// }

// Allocate 电池分配
func (s *batteryService) Allocate(bat *ent.Asset, sub *ent.Subscribe, transferType model.AssetTransferType) (err error) {
	err = ent.Database.Asset.Update().Where(asset.SubscribeID(sub.ID)).ClearSubscribeID().Exec(s.ctx)
	if err != nil {
		return err
	}
	var details []model.AssetTransferCreateDetail
	// 非智能电池
	if bat.Type == model.AssetTypeNonSmartBattery.Value() {
		details = append(details, model.AssetTransferCreateDetail{
			AssetType: model.AssetTypeNonSmartBattery,
			Num:       silk.UInt(1),
			ModelID:   bat.ModelID,
		})
	}
	// 智能电池
	if bat.Type == model.AssetTypeSmartBattery.Value() {
		details = append(details, model.AssetTransferCreateDetail{
			AssetType: model.AssetTypeSmartBattery,
			SN:        silk.String(bat.Sn),
		})
	}

	// 出库方信息
	fromLocationType := model.AssetLocationsType(bat.LocationsType)
	fromLocationID := bat.LocationsID
	// 入库方信息
	toLocationType := model.AssetLocationsTypeRider
	toLocationID := sub.RiderID

	// 判定是否自动入库
	var autoIn bool
	if transferType == model.AssetTransferTypeExchange || transferType == model.AssetTransferTypeTransfer {
		autoIn = true
	}

	// 解绑电池并产生调拨流转
	_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
		FromLocationType:  &fromLocationType,
		FromLocationID:    &fromLocationID,
		ToLocationType:    toLocationType,
		ToLocationID:      toLocationID,
		Details:           details,
		Reason:            "分配电池调拨",
		AssetTransferType: transferType,
		OperatorID:        s.operator.ID,
		OperatorType:      s.operator.Type,
		AutoIn:            autoIn,
		SkipLimit:         true,
	}, &model.Modifier{
		ID:    s.operator.ID,
		Name:  s.operator.Name,
		Phone: s.operator.Phone,
	})
	if err != nil {
		return err
	}
	if failed != nil {
		return errors.New("分配调拨失败" + failed[0])
	}
	err = bat.Update().SetSubscribeID(sub.ID).Exec(s.ctx)
	if err != nil {
		return errors.New("分配失败" + err.Error())
	}
	return nil
}

// Unallocate 清除骑手电池信息
func (s *batteryService) Unallocate(bat *ent.Asset, toLocationType model.AssetLocationsType, toLocationID uint64, transferType model.AssetTransferType) (err error) {
	// 出库方信息
	fromLocationType := model.AssetLocationsType(bat.LocationsType)
	fromLocationID := bat.LocationsID
	var details []model.AssetTransferCreateDetail
	// 非智能电池
	if bat.Type == model.AssetTypeNonSmartBattery.Value() {
		details = append(details, model.AssetTransferCreateDetail{
			AssetType: model.AssetTypeNonSmartBattery,
			Num:       silk.UInt(1),
			ModelID:   bat.ModelID,
		})
	}
	// 智能电池
	if bat.Type == model.AssetTypeSmartBattery.Value() {
		details = append(details, model.AssetTransferCreateDetail{
			AssetType: model.AssetTypeSmartBattery,
			SN:        silk.String(bat.Sn),
		})
	}
	// 解绑电池并产生调拨流转
	_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
		FromLocationType:  &fromLocationType,
		FromLocationID:    &fromLocationID,
		ToLocationType:    toLocationType,
		ToLocationID:      toLocationID,
		Details:           details,
		Reason:            "解绑电池调拨",
		AssetTransferType: transferType,
		OperatorID:        s.operator.ID,
		OperatorType:      s.operator.Type,
		AutoIn:            true, // 自动入库
		SkipLimit:         true,
	}, &model.Modifier{
		ID:    s.operator.ID,
		Name:  s.operator.Name,
		Phone: s.operator.Phone,
	})
	if err != nil {
		return err
	}
	if failed != nil {
		return errors.New("分配调拨失败" + failed[0])
	}
	// 这里调拨单已经完成位置调整 只需要清除订阅信息即可
	err = bat.Update().ClearSubscribeID().Exec(s.ctx)
	if err != nil {
		return err
	}
	return
}
