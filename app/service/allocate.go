// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"time"

	"github.com/auroraride/adapter"
	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/socket"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/allocate"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/person"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type allocateService struct {
	*BaseService
	orm *ent.AllocateClient
}

func NewAllocate(params ...any) *allocateService {
	return &allocateService{
		BaseService: newService(params...),
		orm:         ent.Database.Allocate,
	}
}

func (s *allocateService) QueryID(id uint64) (*ent.Allocate, error) {
	return s.orm.Query().Where(allocate.ID(id)).WithBrand().WithRider().First(s.ctx)
}

func (s *allocateService) QueryIDX(id uint64) *ent.Allocate {
	al, _ := s.QueryID(id)
	if al == nil {
		snag.Panic("未找到信息")
	}
	return al
}

// QueryEffectiveSubscribeID 查询生效中的分配信息
func (s *allocateService) QueryEffectiveSubscribeID(subscribeID uint64) (*ent.Allocate, error) {
	return s.orm.Query().
		Where(
			allocate.SubscribeID(subscribeID),
			allocate.TimeGTE(carbon.Now().SubSeconds(model.AllocateExpiration).StdTime()),
		).
		First(s.ctx)
}

func (s *allocateService) QueryEffectiveSubscribeIDX(subscribeID uint64) *ent.Allocate {
	al, _ := s.QueryEffectiveSubscribeID(subscribeID)
	if al == nil {
		snag.Panic("未找到有效分配信息")
	}
	return al
}

// Create 订单激活分配
func (s *allocateService) Create(params *model.AllocateCreateParams) model.AllocateCreateRes {
	// 禁止骑手请求
	if s.rider != nil {
		snag.Panic("请求方式错误")
	}

	// 查找订阅
	_, sub := NewBusinessRider(nil).Inactive(*params.SubscribeID)

	if sub == nil {
		snag.Panic("未找到订阅信息")
	}

	// 判定非智能约束
	if !sub.Intelligent && params.BatteryID != nil {
		snag.Panic("非智能订阅无法绑定电池")
	}

	// 判定条件
	// 必须有 门店 / 电柜 / 站点其一
	if params.StoreID == nil && sub.StationID == nil {
		snag.Panic("必须由门店或站点激活")
	}

	if params.StoreID != nil && sub.StationID != nil {
		snag.Panic("门店和站点不能同时存在")
	}

	// 车电套餐(V2) 只能由骑手选择的当前门店激活
	// 后台激活不受限制
	if s.modifier == nil {
		if sub.Edges.Plan != nil && sub.Edges.Plan.Type == model.PlanTypeEbikeWithBattery.Value() &&
			sub.StoreID != nil && params.StoreID != nil && *sub.StoreID != *params.StoreID {
			snag.Panic("请前往App选择的门店激活或者更改门店")
		}
	}

	// 是否需要分配电车
	if sub.BrandID != nil && !params.EbikeParam.Exists() {
		snag.Panic("需要分配电车")
	}

	// 获取骑手
	r := sub.Edges.Rider
	if r == nil {
		snag.Panic("骑手查询失败")
		return model.AllocateCreateRes{}
	}

	// 判定骑手是否实名认证
	if exists, _ := r.QueryPerson().Where(person.Status(model.PersonAuthenticated.Value())).Exist(s.ctx); !exists {
		snag.Panic("骑手未实名认证")
	}

	// todo v2版本合同 单电激活流程更新了（先分配在签约 所以后台激活走不通）
	// 查询是否已签约
	if exists, _ := ent.Database.Contract.QueryNotDeleted().Where(
		contract.SubscribeID(sub.ID),
		contract.Status(model.ContractStatusSuccess.Value()),
		contract.Effective(true),
	).Exist(s.ctx); exists {
		snag.Panic("已签约, 无法重新分配")
	}

	// 是否被分配过
	if exists, _ := s.orm.Query().
		Where(
			allocate.SubscribeID(*params.SubscribeID),
			allocate.TimeGT(carbon.Now().SubSeconds(model.AllocateExpiration).StdTime()),
			allocate.StatusIn(model.AllocateStatusPending.Value(), model.AllocateStatusSigned.Value()),
		).
		Exist(s.ctx); exists {
		snag.Panic("已被分配过")
	}

	e := sub.Edges.Enterprise
	if e != nil {
		// 禁止门店端激活 --by 曹博文 2023-10-25
		if e.Agent && s.employee != nil {
			snag.Panic("无法使用门店激活")
		}

		if e.Payment == model.EnterprisePaymentPrepay && e.Balance <= 0 {
			snag.Panic("余额不足,请充值")
		}
	}

	var (
		cityID        uint64
		entStore      *ent.Store
		entStation    *ent.EnterpriseStation
		locationsType model.AssetLocationsType
		locationsID   uint64
	)

	// 门店
	if params.StoreID != nil {
		entStore = NewStore().Query(*params.StoreID)
		cityID = entStore.CityID
		locationsType = model.AssetLocationsTypeStore
		locationsID = *params.StoreID
		checkAsset, _ := NewAsset().CheckAsset(locationsType, *params.StoreID, sub.Model)
		// 判定门店非智能电池库存
		if params.BatteryID == nil && checkAsset == nil {
			snag.Panic("电池库存不足")
		}
	}

	// 站点
	if sub.StationID != nil {
		entStation = NewEnterpriseStation().QueryX(*sub.StationID)
		cityID = *entStation.CityID
		locationsType = model.AssetLocationsTypeStation
		locationsID = *sub.StationID
		checkAsset, _ := NewAsset().CheckAsset(model.AssetLocationsTypeStation, *sub.StationID, sub.Model)
		if params.BatteryID == nil && checkAsset == nil {
			snag.Panic("站点电池库存不足")
		}
	}

	// 判定城市
	if sub.CityID != cityID {
		snag.Panic("无法跨城市操作")
	}

	// 获取并判定电池
	var bat *ent.Asset
	var err error

	// 如果为非智能套餐 找一个非智能电池电池
	if params.BatteryID == nil && !sub.Intelligent {
		md := ent.Database.BatteryModel.Query().Where(batterymodel.Model(sub.Model)).FirstX(s.ctx)
		if md == nil {
			snag.Panic("未找到电池型号")
			return model.AllocateCreateRes{}
		}
		bat, _ = NewAsset().QueryNonSmartBattery(&model.QueryAssetBatteryReq{
			LocationsType: &locationsType,
			LocationsID:   &locationsID,
			ModelID:       md.ID,
		})
	} else {
		// 智能电池
		bat, _ = NewAsset().QueryAssetByLocation(model.QueryAssetReq{
			LocationsType: locationsType,
			LocationsID:   locationsID,
			ID:            params.BatteryID,
		})
	}
	if bat == nil {
		snag.Panic("未找到电池信息")
		return model.AllocateCreateRes{}
	}
	// 盘点电池位置
	switch bat.LocationsType {
	case model.AssetLocationsTypeStation.Value():
		if sub.StationID == nil {
			snag.Panic("未找到站点信息")
		}
		if bat.LocationsID != *sub.StationID {
			snag.Panic("电池站点不符")
		}
	case model.AssetLocationsTypeCabinet.Value():
		snag.Panic("电池在电柜中无法使用")
	case model.AssetLocationsTypeRider.Value():
		snag.Panic("电池已被骑手绑定")
	}
	if bat.Edges.Model != nil && bat.Edges.Model.Model != sub.Model {
		snag.Panic("电池型号不符")
	}
	// 默认单电类型
	typ := allocate.TypeBattery

	// 查找电车
	var bikeID, brandID *uint64
	var bikeInfo *model.EbikeBusinessInfo
	if params.EbikeParam.Exists() {
		// 车电类型
		typ = allocate.TypeEbike

		// 车电必须有门店或站点
		if params.StoreID == nil && sub.StationID == nil {
			snag.Panic("车电套餐调键判定不足")
		}

		// 查找电车
		bike := NewEbike().UnallocatedX(&model.EbikeUnallocatedParams{
			ID:        params.EbikeParam.ID,
			StoreID:   params.StoreID,
			StationID: sub.StationID,
			Keyword:   params.EbikeParam.Keyword,
		})

		if bike == nil {
			snag.Panic("未找到车辆信息")
		}

		if sub.StationID != nil && bike.LocationID != nil && bike.LocationType == model.AssetLocationsTypeStore && *sub.StationID != *bike.LocationID {
			snag.Panic("车辆站点归属不一致")
		}

		// 判定电池是不是在同一个门店 或者站点
		if (bike.LocationType == model.AssetLocationsTypeStore || bike.LocationType == model.AssetLocationsTypeStation) && bike.LocationID != nil && bat.LocationsID != *bike.LocationID {
			snag.Panic("电池与车辆不在同一位置")
		}

		// 比对型号
		if bike.Brand.ID != *sub.BrandID {
			snag.Panic("待分配车辆型号错误")
		}

		bikeID = silk.UInt64(bike.ID)
		brandID = silk.UInt64(bike.Brand.ID)
		bikeInfo = &model.EbikeBusinessInfo{
			ID:           bike.ID,
			BrandID:      bike.Brand.ID,
			BrandName:    bike.Brand.Name,
			Sn:           bike.SN,
			LocationType: bike.LocationType,
			LocationID:   bike.LocationID,
		}
	}

	// 存储分配信息
	status := model.AllocateStatusPending.Value()
	if !sub.NeedContract {
		status = model.AllocateStatusSigned.Value()
	}

	// 强制删除原有分配信息
	s.SubscribeDeleteIfExists(sub.ID)

	// 分配信息
	var allo *ent.Allocate

	// 激活失败作废分配信息
	defer func() {
		if v := recover(); v != nil {
			_ = allo.Update().SetStatus(model.AllocateStatusVoid.Value()).Exec(s.ctx)
			panic(v)
		}
	}()

	// 保存分配信息
	allo, err = s.orm.Create().
		SetType(typ).
		SetNillableEmployeeID(params.EmployeeID).
		SetNillableStoreID(params.StoreID).
		SetNillableStationID(sub.StationID).
		SetNillableBatteryID(params.BatteryID).
		SetNillableEbikeID(bikeID).
		SetNillableBrandID(brandID).
		SetSubscribe(sub).
		SetRider(r).
		SetStatus(status).
		SetTime(time.Now()).
		SetModel(sub.Model).
		SetNillableAgentID(params.AgentID).
		Save(s.ctx)

	if err != nil {
		zap.L().Error("分配失败", zap.Error(err))
		snag.Panic("分配失败")
	}

	// 指定电话号码不需要签约
	debugPhones := ar.Config.App.Debug.Phone
	if debugPhones[r.Phone] && ar.Config.Environment != adapter.Production {
		sub.NeedContract = false
	}

	switch sub.NeedContract {
	case true:
		// 需要签约, 推送签约消息
		socket.SendMessage(NewRiderSocket(), r.ID, &model.RiderSocketMessage{ContractSign: &model.ContractSignReq{
			SubscribeID: sub.ID,
		}})
	default:
		// 无须签约, 直接激活
		NewBusinessRiderWithParams(
			s.modifier,
			entStore,
			entStation,
			bikeInfo,
			bat,
			s.operator,
		).
			SetEmployeeID(allo.EmployeeID).
			SetAgentID(params.AgentID).
			Active(sub, allo)
	}

	return model.AllocateCreateRes{
		ID:           allo.ID,
		NeedContract: sub.NeedContract,
	}
}

// func (s *allocateService) detail(al *ent.Allocate) model.AllocateDetail {
// 	r := al.Edges.Rider
// 	res := model.AllocateDetail{
// 		ID:     al.ID,
// 		Type:   al.Type.String(),
// 		Status: model.AllocateStatus(al.Status),
// 		Time:   al.Time.Format(carbon.DateTimeLayout),
// 		Model:  al.Model,
// 		Rider: model.Rider{
// 			ID:    r.ID,
// 			Phone: r.Phone,
// 			Name:  r.Name,
// 		},
// 		Ebike: NewEbike().Detail(al.Edges.Ebike, al.Edges.Brand),
// 	}
//
// 	if time.Since(al.Time).Seconds() > model.AllocateExpiration && res.Status == model.AllocateStatusPending {
// 		res.Status = model.AllocateStatusVoid
// 	}
//
// 	return res
// }

// Info 分配信息
// func (s *allocateService) Info(req *model.IDParamReq) model.AllocateDetail {
// 	al := s.QueryIDX(req.ID)
// 	return s.detail(al)
// }

// EmployeeList 电车分配店员列表
// func (s *allocateService) EmployeeList(req *model.AllocateEmployeeListReq) *model.PaginationRes {
// 	status := req.Status
// 	if req.Status != 2 {
// 		status = model.AllocateStatusPending
// 	}
// 	q := s.orm.Query().
// 		WithRider().
// 		WithEbike().
// 		WithBrand().
// 		Where(allocate.EmployeeID(s.employee.ID), allocate.Status(status.Value())).
// 		Order(ent.Desc(allocate.FieldTime))
// 	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Allocate) model.AllocateDetail {
// 		return s.detail(item)
// 	})
// }

// LoopStatus 长连接查询是否已分配
func (s *allocateService) LoopStatus(riderID, subscribeID uint64) (res model.AllocateRiderRes) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	start := time.Now()
	for ; true; <-ticker.C {

		allo, _ := s.orm.Query().Where(
			allocate.RiderID(riderID),
			allocate.SubscribeID(subscribeID),
			allocate.TimeGT(carbon.Now().SubSeconds(model.AllocateExpiration).StdTime()),
		).First(s.ctx)

		// 如果有分配信息 并且 状态为待签约 并且 非电柜扫码
		if allo != nil && allo.Status == model.AllocateStatusPending.Value() && allo.CabinetID == nil {
			res.Allocated = true
			return
		}

		if time.Since(start).Seconds() > 50 {
			return
		}
	}

	return
}

// SubscribeDeleteIfExists 根据subscribeID强制删除分配信息
func (s *allocateService) SubscribeDeleteIfExists(subscribeID uint64) {
	_, err := s.orm.Delete().Where(allocate.SubscribeID(subscribeID)).Exec(s.ctx)
	if err != nil {
		zap.L().Error("分配信息强制删除失败", zap.Error(err))
	}
}
