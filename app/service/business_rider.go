// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/auroraride/adapter/async"
	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/commission"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/orderrefund"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribepause"
	"github.com/auroraride/aurservd/internal/payment/alipay"
	"github.com/auroraride/aurservd/internal/payment/wechat"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type businessRiderService struct {
	ctx          context.Context
	modifier     *model.Modifier
	rider        *ent.Rider
	employee     *ent.Employee
	employeeInfo *model.Employee
	cabinet      *ent.Cabinet
	cabinetInfo  *model.CabinetBasicInfo
	store        *ent.Store
	ebikeStore   *ent.Store
	batStore     *ent.Store
	subscribe    *ent.Subscribe
	reserve      *ent.Reserve

	cabTask func() (*model.BinInfo, *model.Battery, error) // 电柜任务

	storeID, employeeID, cabinetID, subscribeID, agentID, ebikeStoreID, batStoreID, stationID *uint64

	// 电车信息
	ebikeInfo *model.EbikeBusinessInfo

	// 智能电池信息
	battery *ent.Asset

	*BaseService
}

func NewBusinessRider(params ...any) *businessRiderService {
	return &businessRiderService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
}

func (s *businessRiderService) SetModifier(m *model.Modifier) *businessRiderService {
	if m != nil {
		s.modifier = m
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	}
	return s
}

// SetCabinet 设置电柜
func (s *businessRiderService) SetCabinet(cab *ent.Cabinet) *businessRiderService {
	if cab != nil {
		s.cabinet = cab
		s.cabinetInfo = &model.CabinetBasicInfo{
			ID:     s.cabinet.ID,
			Brand:  s.cabinet.Brand,
			Serial: s.cabinet.Serial,
			Name:   s.cabinet.Name,
		}
	}
	return s
}

// SetCabinetID 设置电柜
func (s *businessRiderService) SetCabinetID(id *uint64) *businessRiderService {
	if id != nil {
		s.SetCabinet(NewCabinet().QueryOne(*id))
	}
	return s
}

// SetStoreID 设置门店ID
func (s *businessRiderService) SetStoreID(id *uint64) *businessRiderService {
	if id != nil {
		s.store = NewStore().Query(*id)
	}
	return s
}

// SetEbikeStoreID 设置电车门店ID（退租使用）
func (s *businessRiderService) SetEbikeStoreID(id *uint64) *businessRiderService {
	if id != nil {
		s.ebikeStore = NewStore().Query(*id)
	}
	return s
}

// SetBatStoreID 设置电池门店ID（退租使用）
func (s *businessRiderService) SetBatStoreID(id *uint64) *businessRiderService {
	if id != nil {
		s.batStore = NewStore().Query(*id)
	}
	return s
}

// SetAgentID 设置代理商ID
func (s *businessRiderService) SetAgentID(id *uint64) *businessRiderService {
	if id != nil {
		s.agentID = id
	}
	return s
}

func (s *businessRiderService) SetEmployeeID(id *uint64) *businessRiderService {
	if id != nil {
		s.employee, _ = NewEmployee().Query(*id)
		if s.employee != nil {
			s.employeeID = id
			s.employeeInfo = &model.Employee{
				ID:    s.employee.ID,
				Name:  s.employee.Name,
				Phone: s.employee.Phone,
			}
		}
	}
	return s
}

// SetBatteryID 设置电池
func (s *businessRiderService) SetBatteryID(id *uint64) *businessRiderService {
	if id == nil {
		return s
	}
	// 查找电池
	bat, err := NewAsset().QueryID(*id)
	if err != nil {
		snag.Panic("电池查询失败")
	}
	s.battery = bat
	return s
}

// SetEbikeID 设置电车
func (s *businessRiderService) SetEbikeID(id *uint64) *businessRiderService {
	if id == nil {
		return s
	}
	bike, _ := ent.Database.Asset.Query().Where(asset.ID(*id), asset.Type(model.AssetTypeEbike.Value())).WithBrand().First(s.ctx)
	if bike == nil || bike.Edges.Brand == nil {
		snag.Panic("电车信息查询失败")
		return nil
	}
	brand := bike.Edges.Brand
	s.ebikeInfo = &model.EbikeBusinessInfo{
		ID:           bike.ID,
		BrandID:      brand.ID,
		BrandName:    brand.Name,
		Sn:           bike.Sn,
		LocationType: model.AssetLocationsType(bike.LocationsType),
		LocationID:   &bike.LocationsID,
	}

	return s
}

func (s *businessRiderService) SetEbike(info *model.EbikeBusinessInfo) *businessRiderService {
	s.ebikeInfo = info
	return s
}

func (s *businessRiderService) SetCabinetTask(task func() (*model.BinInfo, *model.Battery, error)) *businessRiderService {
	if task != nil {
		s.cabTask = task
	}
	return s
}

// NewBusinessRiderWithParams 设置参数并初始化
func NewBusinessRiderWithParams(params ...any) *businessRiderService {
	s := &businessRiderService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
	for _, param := range params {
		switch v := param.(type) {
		case *model.Modifier:
			s.modifier = v
			s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, v)
		case *ent.Rider:
			s.rider = v
		case *ent.Store:
			s.store = v
		case *ent.Asset:
			s.battery = v
		case *ent.Subscribe:
			s.subscribe = v
		case *ent.Cabinet:
			s.cabinet = v
		case *ent.Employee:
			s.employee = v
			s.employeeInfo = &model.Employee{
				ID:    v.ID,
				Name:  v.Name,
				Phone: v.Phone,
			}
		case *model.EbikeBusinessInfo:
			s.ebikeInfo = v
		}
	}
	return s
}

func NewBusinessRiderWithEmployee(e *ent.Employee) *businessRiderService {
	s := NewBusinessRider(nil)
	if e == nil {
		snag.Panic("店员错误")
	}
	store := e.Edges.Store
	if store == nil {
		snag.Panic("未找到所属门店")
	}
	s.store = store
	s.employee = e
	s.employeeInfo = &model.Employee{
		ID:    e.ID,
		Name:  e.Name,
		Phone: e.Phone,
	}
	s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
	return s
}

// QuerySubscribeWithRider 查询订阅信息
func (s *businessRiderService) QuerySubscribeWithRider(subscribeID uint64) *ent.Subscribe {
	item, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(subscribeID)).WithEnterprise().WithRider().First(s.ctx)
	if item == nil {
		snag.Panic("未找到对应订阅")
	}
	return item
}

// Inactive 获取待激活订阅信息
func (s *businessRiderService) Inactive(id uint64) (*model.SubscribeActiveInfo, *ent.Subscribe) {
	// 查询订单状态
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(
			subscribe.ID(id),
			subscribe.RefundAtIsNil(),
			subscribe.StartAtIsNil(),
			subscribe.Or(
				subscribe.Type(0),
				subscribe.TypeIn(model.OrderTypeNewly, model.OrderTypeAgain),
			),
			subscribe.Status(model.SubscribeStatusInactive),
		).
		WithInitialOrder(func(oq *ent.OrderQuery) {
			oq.WithPlan().WithCommission()
		}).
		WithRider().
		WithPlan().
		WithEnterprise().
		WithCity().
		WithBrand().
		First(s.ctx)

	if sub == nil {
		snag.Panic("未找到待激活骑士卡")
	}

	if s.employee != nil {
		NewBusinessWithEmployee(s.employee).CheckCity(sub.CityID, s.store)
	}

	r := sub.Edges.Rider
	if r == nil {
		snag.Panic("骑手信息获取失败")
	}

	res := &model.SubscribeActiveInfo{
		ID:           sub.ID,
		EnterpriseID: sub.EnterpriseID,
		Model:        sub.Model,
		CommissionID: nil,
		Rider: model.Rider{
			ID:    r.ID,
			Phone: r.Phone,
			Name:  r.Name,
		},
		City: model.City{
			ID:   sub.Edges.City.ID,
			Name: sub.Edges.City.Name,
		},
	}

	if sub.EnterpriseID == nil {
		o := sub.Edges.InitialOrder
		if o == nil || o.Status != model.OrderStatusPaid {
			snag.Panic("订单状态异常") // TODO 退款被拒绝的如何操作
		}
		res.Order = &model.SubscribeOrderInfo{
			ID:      o.ID,
			Status:  o.Status,
			PayAt:   o.CreatedAt.Format(carbon.DateTimeLayout),
			Payway:  o.Payway,
			Amount:  o.Amount,
			Deposit: o.Total - o.Amount,
			Total:   o.Total,
		}

		c := sub.Edges.InitialOrder.Edges.Commission
		if c != nil {
			res.CommissionID = &c.ID
		}
	} else {
		en := sub.Edges.Enterprise
		res.Enterprise = &model.Enterprise{
			ID:    en.ID,
			Name:  en.Name,
			Agent: en.Agent,
		}
	}

	if sub.BrandID != nil {
		brand := sub.Edges.Brand
		if brand == nil {
			snag.Panic("电车型号查询失败")
		}
		res.EbikeBrand = &model.EbikeBrand{
			ID:    brand.ID,
			Name:  brand.Name,
			Cover: brand.Cover,
		}
	}

	return res, sub
}

// Preprocess 预处理数据
func (s *businessRiderService) Preprocess(bt model.BusinessType, sub *ent.Subscribe) {
	s.subscribe = sub

	if sub.EnterpriseID != nil {
		en := sub.Edges.Enterprise
		if en == nil {
			snag.Panic("未找到团签信息")
		}
		// 判定是否寄存或取消寄存业务
		if bt == model.BusinessTypePause || bt == model.BusinessTypeContinue {
			snag.Panic("团签用户无法办理")
		}
		// 判定代理是否可使用门店
		if en.Agent && !en.UseStore && s.employee != nil {
			snag.Panic("代理无法在门店办理业务")
		}

	}

	s.subscribeID = silk.Pointer(sub.ID)

	r := sub.Edges.Rider
	if r == nil {
		r, _ = sub.QueryRider().First(s.ctx)
	}

	if r == nil {
		snag.Panic("骑手查询失败")
	}

	// 骑士卡状态
	if !NewRiderBusiness(r).Executable(sub, bt) {
		snag.Panic("骑士卡状态错误")
	}

	// 检查用户是否可以办理业务
	NewRiderPermissionWithRider(r, s.modifier).BusinessX().SubscribeX(model.RiderPermissionTypeBusiness, sub)

	s.rider = r

	if s.store != nil {
		s.storeID = silk.Pointer(s.store.ID)
	}

	if s.cabinet != nil {
		s.cabinetID = silk.Pointer(s.cabinet.ID)
	}

	// 以租代购管理员强制退租专用
	if s.ebikeStore != nil {
		s.ebikeStoreID = silk.Pointer(s.ebikeStore.ID)
	}
	if s.batStore != nil {
		s.batStoreID = silk.Pointer(s.batStore.ID)
	}

	// 判定是否满足业务调键
	// 代理站点的骑手无需门店或电柜
	// 反之则必须门店或电柜
	if sub.StationID == nil && s.store == nil && s.cabinet == nil && s.ebikeStore == nil && s.batStore == nil {
		snag.Panic("条件不满足")
	}

	if s.employee != nil {
		s.employeeID = silk.Pointer(s.employee.ID)
	}

	if s.employee == nil && s.store == nil && s.modifier == nil && s.cabinet == nil && sub.StationID == nil && s.ebikeStore == nil && s.batStore == nil {
		snag.Panic("操作权限校验失败")
	}

	// 校验权限
	if s.employee != nil {
		NewBusinessWithEmployee(s.employee).CheckCity(s.subscribe.CityID, s.store)
		// 以租代购专用检验电车退租门店、电池退租门店
		NewBusinessWithEmployee(s.employee).CheckCity(s.subscribe.CityID, s.ebikeStore)
		NewBusinessWithEmployee(s.employee).CheckCity(s.subscribe.CityID, s.batStore)
	}

	// 车电订阅检查
	if sub.BrandID != nil {
		// 车电订阅无法办理寄存相关业务
		// 车电订阅无法使用电柜（以租代购退租时，车电业务存在电车还门店、电池还电柜/门店情况）
		if bt != model.BusinessTypeActive && bt != model.BusinessTypeUnsubscribe {
			snag.Panic("车电订阅无法办理此业务")
		}
	}

	// 如果是车电订阅, 查询并设置电车信息
	if sub.EbikeID != nil && s.ebikeInfo == nil {
		s.SetEbikeID(sub.EbikeID)
	}

	// 预约检查
	rev := NewReserveWithRider(r).RiderUnfinished(r.ID)
	if rev != nil {
		if s.cabinet == nil || s.cabinet.ID != rev.CabinetID || bt.String() != rev.Type {
			_, _ = rev.Update().SetStatus(model.ReserveStatusInvalid.Value()).Save(s.ctx)
		} else {
			// 预约处理中
			s.reserve, _ = rev.Update().SetStatus(model.ReserveStatusProcessing.Value()).Save(s.ctx)
		}
	}

}

// doTask 处理电柜任务
func (s *businessRiderService) doTask() (bin *model.BinInfo, bat *model.Battery, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

	bin, bat, err = s.cabTask()
	return
}

// do 处理业务
func (s *businessRiderService) do(doReq model.BusinessRiderServiceDoReq, cb func(tx *ent.Tx)) {
	async.WithTask(func() {
		sts := map[model.BusinessType]model.AssetTransferType{
			model.BusinessTypeActive:      model.AssetTransferTypeActive,
			model.BusinessTypeUnsubscribe: model.AssetTransferTypeUnSubscribe,
			model.BusinessTypePause:       model.AssetTransferTypePause,
			model.BusinessTypeContinue:    model.AssetTransferTypeContinue,
		}

		ops := map[model.BusinessType]model.Operate{
			model.BusinessTypeActive:      model.OperateActive,
			model.BusinessTypeUnsubscribe: model.OperateUnsubscribe,
			model.BusinessTypePause:       model.OperateSubscribePause,
			model.BusinessTypeContinue:    model.OperateSubscribeContinue,
		}

		bfs := map[model.BusinessType]string{
			model.BusinessTypeActive:      "未激活",
			model.BusinessTypeUnsubscribe: "生效中",
			model.BusinessTypePause:       "计费中",
			model.BusinessTypeContinue:    "寄存中",
		}

		afs := map[model.BusinessType]string{
			model.BusinessTypeActive:      "已激活",
			model.BusinessTypeUnsubscribe: "已退租",
			model.BusinessTypePause:       "已寄存",
			model.BusinessTypeContinue:    "计费中",
		}

		var bin *model.BinInfo
		var err error

		// 放入电池优先执行
		var bat *model.Battery
		if s.battery != nil {
			var m string
			if s.battery.Edges.Model != nil {
				m = s.battery.Edges.Model.Model
			}
			bat = &model.Battery{
				ID:    s.battery.ID,
				SN:    s.battery.Sn,
				Model: m,
			}
		}

		// 放入电池优先执行电柜任务
		if s.cabTask != nil && (doReq.Type == model.BusinessTypePause || doReq.Type == model.BusinessTypeUnsubscribe) {
			bin, bat, err = s.doTask()
			if err != nil {
				snag.Panic(err)
			}
		}

		// 激活业务查找提成
		var co *ent.Commission
		if doReq.Type == model.BusinessTypeActive {
			co, _ = ent.Database.Commission.QueryNotDeleted().Where(commission.SubscribeID(s.subscribe.ID)).First(s.ctx)
		}

		// 库存管理
		// var batSk *ent.Stock
		ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
			cb(tx)

			// 需要进行业务出入库
			if s.cabinetID != nil || s.storeID != nil || s.subscribe.StationID != nil || s.ebikeStoreID != nil || s.batStoreID != nil {
				err = NewAsset(s.modifier, s.operator).RiderBusiness(
					&model.StockBusinessReq{
						RiderID:           s.subscribe.RiderID,
						Model:             s.subscribe.Model,
						CityID:            s.subscribe.CityID,
						AssetTransferType: sts[doReq.Type],
						BusinessType:      doReq.Type,

						StoreID:     s.storeID,
						EmployeeID:  s.employeeID,
						CabinetID:   s.cabinetID,
						SubscribeID: s.subscribeID,

						StationID:    s.stationID,
						EnterpriseID: s.subscribe.EnterpriseID,
						AgentID:      s.agentID,

						Ebike:   s.ebikeInfo,
						Battery: bat,

						Rto:          doReq.Rto,
						EbikeStoreID: s.ebikeStoreID,
						BatStoreID:   s.batStoreID,
					},
				)

				if err != nil {
					zap.L().Error("骑手业务出入库失败: "+doReq.Type.String(), zap.Error(err))
				}
			}

			return
		})

		// 取出电池滞后执行电柜任务
		if s.cabTask != nil && (doReq.Type == model.BusinessTypeActive || doReq.Type == model.BusinessTypeContinue) {
			bin, bat, err = s.doTask()
			if err != nil {
				zap.L().Error("骑手业务取出电池后任务执行失败: "+doReq.Type.String(), zap.Error(err))
				return
			}
			b, _ := NewAsset().QueryID(bat.ID)
			if b == nil {
				zap.L().Error(fmt.Sprintf("业务：电池查询失败  sn:%s, id:%d", bat.SN, bat.ID))
				return
			}
			// 查询调拨单
			t, _ := NewAssetTransfer().QueryTransferByAssetID(s.ctx, bat.ID)
			if t == nil {
				zap.L().Error("调拨单查询失败")
				return
			}
			detail := make([]model.AssetTransferReceiveDetail, 0)
			if b.Type == model.AssetTypeSmartBattery.Value() {
				detail = append(detail, model.AssetTransferReceiveDetail{
					AssetType: model.AssetType(b.Type),
					SN:        silk.String(b.Sn),
				})
			}
			if b.Type == model.AssetTypeNonSmartBattery.Value() {
				detail = append(detail, model.AssetTransferReceiveDetail{
					AssetType: model.AssetType(b.Type),
					Num:       silk.UInt(1),
					ModelID:   b.ModelID,
				})
			}

			err = NewAssetTransfer(s.operator).TransferReceive(context.Background(), &model.AssetTransferReceiveBatchReq{
				OperateType: s.operator.Type,
				AssetTransferReceive: []model.AssetTransferReceiveReq{
					{
						ID:     t.ID,
						Detail: detail,
						Remark: silk.String("骑手业务" + doReq.Type.String() + "电池接收"),
					},
				},
			}, &model.Modifier{
				ID:    s.operator.ID,
				Name:  s.operator.Name,
				Phone: s.operator.Phone,
			})
			if err != nil {
				return
			}
		}

		// 保存业务日志
		var b *ent.Business
		var bq *businessLogService
		bq = NewBusinessLog(s.subscribe).
			SetModifier(s.modifier).
			SetEmployee(s.employee).
			SetCabinet(s.cabinet).
			SetStore(s.store).
			SetBinInfo(bin).
			// SetAssetTransfer(batSk).
			SetBattery(bat).
			SetAgentId(s.agentID)
		// 电池归还门店必须保持与电车归还门店一致
		if s.batStoreID != nil && s.ebikeStoreID != nil && *s.batStoreID != *s.ebikeStoreID {
			snag.Panic("电池退租门店必须与车辆选择门店一致")
		}
		// 电池单独归还门店(单电退租，满足以租代购退租)
		if s.batStoreID != nil {
			bq.SetStore(s.batStore)
		}
		// 电车归还门店
		if s.ebikeStore != nil {
			bq.SetStore(s.ebikeStore)
		}
		// 满足以租代购
		if doReq.Rto && s.ebikeInfo != nil {
			bq.SetRtoEbikeID(s.ebikeInfo.ID)
		}

		if doReq.Remark != nil {
			bq.SetRemark(*doReq.Remark)
		}

		b, _ = bq.Save(doReq.Type)

		var bussinessID *uint64
		revStatus := model.ReserveStatusFail
		if b != nil {
			revStatus = model.ReserveStatusSuccess
			bussinessID = silk.Pointer(b.ID)
		}

		// 更新预约
		if s.reserve != nil {
			_, _ = s.reserve.Update().
				SetStatus(revStatus.Value()).
				SetNillableBusinessID(bussinessID).
				Save(s.ctx)
		}

		// 更新提成
		if doReq.Type == model.BusinessTypeActive && co != nil && b != nil && s.employeeID != nil {
			_, _ = co.Update().SetBusiness(b).SetEmployeeID(*s.employeeID).Save(s.ctx)
		}

		// 记录日志
		olog := logging.NewOperateLog().
			SetRef(s.rider).
			SetOperate(ops[doReq.Type]).
			SetEmployee(s.employeeInfo).
			SetModifier(s.modifier).
			SetCabinet(s.cabinetInfo).
			SetDiff(bfs[doReq.Type], afs[doReq.Type])
		if doReq.Rto && s.ebikeInfo != nil {
			olog.SetRemark("满足以租代购条件, 电车归属骑手 (" + s.ebikeInfo.Sn + ")")
		}
		go olog.Send()

		if err != nil {
			snag.Panic(err)
		}
	})
}

// Active 激活订阅
func (s *businessRiderService) Active(sub *ent.Subscribe, allo *ent.Allocate) {
	// 设置代理id
	s.agentID = allo.AgentID

	s.Preprocess(model.BusinessTypeActive, sub)

	if NewSubscribe().NeedContract(sub) {
		snag.Panic("还未签约, 请签约")
	}

	s.do(model.BusinessRiderServiceDoReq{Type: model.BusinessTypeActive}, func(tx *ent.Tx) {
		var err error

		// 更新分配
		err = tx.Allocate.UpdateOne(allo).SetStatus(model.AllocateStatusSigned.Value()).Exec(s.ctx)
		if err != nil {
			return
		}

		var (
			aend *time.Time
		)

		endAt := silk.Pointer(tools.NewTime().WillEnd(time.Now(), sub.InitialDays))

		// 如果是代理商, 计算骑士卡代理商结束时间
		if sub.EnterpriseID != nil {
			if sub.Edges.Enterprise == nil {
				sub.Edges.Enterprise = sub.QueryEnterprise().FirstX(s.ctx)
			}
			if sub.Edges.Enterprise == nil {
				snag.Panic("未找到团签信息")
			}
			if sub.Edges.Enterprise.Agent {
				aend = endAt
			}
		}

		updater := tx.Subscribe.UpdateOneID(sub.ID).
			SetStatus(model.SubscribeStatusUsing).
			SetStartAt(time.Now()).
			UpdateTarget(s.cabinetID, s.storeID, s.employeeID).
			SetNillableAgentEndAt(aend).
			SetNeedContract(false)

		// 设置订阅电车
		if s.ebikeInfo != nil {
			updater.SetEbikeID(s.ebikeInfo.ID).SetBrandID(s.ebikeInfo.BrandID)
		}

		// 激活
		s.subscribe, err = updater.Save(s.ctx)
		snag.PanicIfError(err)

		// 更新电车
		if s.ebikeInfo != nil {
			// // 更新电车所属
			eb, _ := NewAsset().QuerySn(s.ebikeInfo.Sn)
			if eb == nil {
				snag.Panic("未找到电车信息")
			}
			err = eb.Update().SetSubscribeID(s.subscribe.ID).Exec(s.ctx)
			if err != nil {
				return
			}
		}
		// 后台操作设置电池编码
		if s.battery != nil && s.cabinet == nil {
			err = s.battery.Update().SetSubscribeID(s.subscribe.ID).SetNillableOrdinal(nil).Exec(s.ctx)
			if err != nil {
				snag.Panic(err)
			}
		}

		if sub.EnterpriseID == nil && sub.InitialOrderID != 0 {
			// 若有退款, 则标记更新状态为失败
			of, _ := tx.OrderRefund.QueryNotDeleted().Where(orderrefund.OrderID(sub.InitialOrderID)).First(s.ctx)
			if of != nil {
				err = tx.OrderRefund.UpdateOne(of).SetReason("激活订阅, 自动拒绝退款").SetStatus(model.OrderStatusRefundRefused).Exec(s.ctx)
				if err != nil {
					return
				}
				_ = tx.Order.UpdateOneID(of.OrderID).SetStatus(model.OrderStatusPaid).Exec(s.ctx)
			}
			// 更新订阅到期时间
			err = tx.Order.UpdateOneID(sub.InitialOrderID).SetSubscribeEndAt(*endAt).Exec(s.ctx)
			if err != nil {
				zap.L().Error("更新订阅到期时间失败", zap.Error(err))
				return
			}
		}
		// 更新出租位置 查询当前资产是否在骑手
		if s.ebikeInfo != nil {
			as, _ := ent.Database.Asset.QueryNotDeleted().
				Where(
					asset.Sn(s.ebikeInfo.Sn),
					asset.LocationsType(model.AssetLocationsTypeRider.Value()),
					asset.LocationsID(sub.RiderID),
				).First(s.ctx)
			if as != nil {
				_ = as.Update().SetRentLocationsType(s.ebikeInfo.LocationType.Value()).SetRentLocationsID(*s.ebikeInfo.LocationID).Exec(s.ctx)
			}
		}
	})

	if sub.EnterpriseID != nil {
		// 更新团签账单
		go NewEnterprise().UpdateStatement(sub.Edges.Enterprise)
	}

	// 返佣计算
	// 团签用户不返佣 个签用户返佣 新签和重签
	if sub.EnterpriseID == nil && sub.Type == model.OrderTypeNewly || sub.Type == model.OrderTypeAgain {
		go NewPromotionCommissionService().RiderActivateCommission(sub)
	}
}

// UnSubscribe 退租
// 会抹去欠费情况
func (s *businessRiderService) UnSubscribe(req *model.BusinessSubscribeReq, fns ...func(sub *ent.Subscribe)) {
	sub := s.QuerySubscribeWithRider(req.ID)
	// 处理单电无电池退租
	if !sub.QueryBattery().ExistX(s.ctx) && sub.EbikeID == nil {
		s.NoBatteryUnsubscribe(req, sub)
		return
	}

	// 预处理业务信息
	s.Preprocess(model.BusinessTypeUnsubscribe, sub)

	if len(fns) > 0 {
		fns[0](s.subscribe)
	}

	// 代理商操作退租
	if req.AgentID != nil {
		if req.StationID == nil {
			snag.Panic("退租必须指定站点")
		}
		s.agentID = req.AgentID
		s.stationID = req.StationID
	}

	// 查找电池
	s.battery, _ = ent.Database.Asset.Query().Where(
		asset.SubscribeID(sub.ID),
		asset.TypeIn(model.AssetTypeSmartBattery.Value(), model.AssetTypeNonSmartBattery.Value()),
	).First(s.ctx)

	err := NewSubscribe().UpdateStatus(sub, false)
	if err != nil {
		snag.Panic(err)
	}

	// 判定退租是否满足条件
	if s.modifier == nil {
		if sub.Remaining < 0 {
			snag.Panic("欠费中, 无法继续办理")
		}
	} else {
		if sub.Remaining < 0 || sub.Status == model.SubscribeStatusOverdue {
			sub.Status = model.SubscribeStatusUsing
		}
	}

	// 校验代理商站点
	if s.rider.StationID != nil && s.cabinet != nil && *s.rider.StationID != *s.cabinet.StationID {
		snag.Panic("请到指定站点退租")
	}

	// 判定是否以租代购
	doReq := model.BusinessRiderServiceDoReq{
		Rto:    req.Rto != nil && *req.Rto,
		Type:   model.BusinessTypeUnsubscribe,
		Remark: req.RtoRemark,
	}

	s.do(doReq, func(tx *ent.Tx) {
		var reason string
		if s.cabinet != nil {
			reason = "用户电柜退租"
		}
		if s.modifier != nil {
			reason = "管理员操作强制退租"
		}
		if s.employee != nil {
			reason = "店员操作退租"
		}

		// 代理商小程序退租
		if s.subscribe.EnterpriseID != nil && s.modifier == nil {
			reason = "代理商操作退租"
		}

		_, err = tx.Subscribe.
			UpdateOneID(sub.ID).
			SetEndAt(time.Now()).
			SetStatus(model.SubscribeStatusUnSubscribed).
			SetUnsubscribeReason(reason).
			Save(s.ctx)
		snag.PanicIfError(err)

		// 查询并标记用户合同为失效
		_, err = tx.Contract.Update().Where(contract.RiderID(sub.RiderID)).SetEffective(false).Save(s.ctx)
		snag.PanicIfError(err)

		// 更新电车
		if sub.EbikeID != nil {
			eb, _ := NewAsset().QueryID(*sub.EbikeID)
			if eb == nil {
				snag.Panic("未找到电车信息")
			}
			if doReq.Rto {
				// 当前属于以租代购
				err = eb.Update().SetRtoRiderID(sub.RiderID).Exec(s.ctx)
				snag.PanicIfError(err)
			}
			err = eb.Update().ClearSubscribeID().Exec(s.ctx)
			if err != nil {
				snag.Panic(err)
			}
		}

		// 删除电池
		if s.battery != nil {
			// 这里调拨单已经完成位置调整 只需要清除订阅信息即可
			err = s.battery.Update().ClearSubscribeID().Exec(s.ctx)
			if err != nil {
				snag.Panic(err)
			}
			// snag.PanicIfError(NewBattery(s.operator, s.modifier).Unallocate(s.battery, toLocationType, toLocationID, model.AssetTransferTypeUnSubscribe))
		}
	})

	// 更新企业账单
	if sub.EnterpriseID != nil {
		go NewEnterprise().UpdateStatementByID(*sub.EnterpriseID)
	}

	// 处理退款
	err = s.ForceUnsubscribe(req, sub.ID, sub.RiderID)
	if err != nil {
		zap.L().Error("退款失败", zap.Error(err))
		snag.Panic(err)
	}

	// 清除电车出租位置
	if s.ebikeInfo != nil {
		_ = ent.Database.Asset.UpdateOneID(s.ebikeInfo.ID).ClearRentLocationsID().ClearRentLocationsType().Exec(s.ctx)
	}
}

// Pause 寄存
func (s *businessRiderService) Pause(subscribeID uint64) {
	s.Preprocess(model.BusinessTypePause, s.QuerySubscribeWithRider(subscribeID))

	if s.subscribe.Remaining < 1 {
		snag.Panic("当前剩余时间不足, 无法寄存")
	}

	s.do(model.BusinessRiderServiceDoReq{Type: model.BusinessTypePause}, func(tx *ent.Tx) {
		_, err := tx.SubscribePause.Create().
			SetStartAt(time.Now()).
			SetRiderID(s.subscribe.RiderID).
			SetSubscribeID(s.subscribe.ID).
			SetCityID(s.subscribe.CityID).
			SetNillableStoreID(s.storeID).
			SetNillableCabinetID(s.cabinetID).
			SetNillableEmployeeID(s.employeeID).Save(s.ctx)
		snag.PanicIfError(err)

		_, err = tx.Subscribe.UpdateOne(s.subscribe).
			SetPausedAt(time.Now()).
			SetStatus(model.SubscribeStatusPaused).
			Save(s.ctx)
		snag.PanicIfError(err)
	})
}

// Continue 取消寄存
func (s *businessRiderService) Continue(subscribeID uint64) {
	s.Preprocess(model.BusinessTypeContinue, s.QuerySubscribeWithRider(subscribeID))

	// 更新订阅信息
	err := NewSubscribe().UpdateStatus(s.subscribe, false)
	if err != nil {
		zap.L().Error("骑士卡更新失败", zap.Error(err))
		snag.Panic("骑士卡更新失败")
	}

	pr, _ := s.subscribe.GetAdditionalItems()

	sp, _ := ent.Database.SubscribePause.QueryNotDeleted().
		Where(subscribepause.SubscribeID(s.subscribe.ID), subscribepause.EndAtIsNil()).
		Order(ent.Desc(subscribepause.FieldCreatedAt)).
		First(s.ctx)

	if sp == nil || pr.Current == nil || pr.Current.ID != sp.ID {
		snag.Panic("未找到寄存信息")
	}

	// 当前时间
	now := time.Now()

	s.do(model.BusinessRiderServiceDoReq{Type: model.BusinessTypeContinue}, func(tx *ent.Tx) {
		_, err = tx.SubscribePause.
			UpdateOne(sp).
			SetDays(pr.CurrentDays).
			SetEndAt(now).
			SetNillableEndEmployeeID(s.employeeID).
			SetEndModifier(s.modifier).
			SetOverdueDays(pr.CurrentOverdueDays).
			SetNillableEndStoreID(s.storeID).
			SetNillableEndCabinetID(s.cabinetID).
			SetSuspendDays(pr.CurrentDuplicateDays).
			Save(s.ctx)
		snag.PanicIfError(err)

		// 更新订阅
		_, err = tx.Subscribe.UpdateOne(s.subscribe).
			SetStatus(model.SubscribeStatusUsing).
			SetPauseDays(pr.Days).
			ClearPausedAt().
			Save(s.ctx)
		snag.PanicIfError(err)
	})
}

// ForceUnsubscribe 处理强制退租预支付和退款
func (s *businessRiderService) ForceUnsubscribe(req *model.BusinessSubscribeReq, subscribeID uint64, riderID uint64) error {
	err := ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
		// 处理预支付订单
		o, _ := ent.Database.Order.QueryNotDeleted().Where(order.SubscribeID(subscribeID)).First(s.ctx)
		// 预授权支付的订单转支付
		if o != nil && o.Payway == model.OrderPaywayAlipayAuthFreeze && (o.Type == model.OrderTypeNewly || o.Type == model.OrderTypeAgain) {
			// 转支付
			err = NewOrder().TradePay(o)
			if err != nil {
				return err
			}
		}

		sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.ID(subscribeID)).WithInitialOrder().First(s.ctx)
		if sub == nil || sub.Edges.InitialOrder == nil {
			zap.L().Info("强制退租, 未找到订阅信息", zap.Uint64("subscribe_id", subscribeID))
			return
		}

		var depositOrder *ent.Order
		if sub.Edges.InitialOrder != nil {
			depositOrder, _ = ent.Database.Order.Query().Where(
				order.Or(
					order.ParentID(sub.Edges.InitialOrder.ID),
					order.ID(sub.Edges.InitialOrder.ID),
				),
				order.Status(model.OrderStatusPaid),
				order.Type(model.OrderTypeDeposit)).First(s.ctx)
			if depositOrder == nil {
				zap.L().Info("强制退租, 未找到押金订单", zap.Uint64("orderId", sub.Edges.InitialOrder.ID))
				return
			}
		}

		if depositOrder != nil && req.RefundDeposit != nil {
			var remainAmount float64
			var reason string
			var refundStatus, orderStatus uint8

			if *req.RefundDeposit {
				reason = "订阅已退订，系统自动退押"
				refundStatus = model.RefundStatusPending
				orderStatus = model.OrderStatusRefundPending
			} else {
				reason = "拒绝退押"
				refundStatus = model.RefundStatusRefused
				orderStatus = model.OrderStatusRefundRefused
			}

			// 如果没传押金金额 默认全退
			if req.DepositAmount != nil {
				if *req.DepositAmount > depositOrder.Amount {
					return fmt.Errorf("退款金额超出押金金额")
				}
				reason = "人工退押"
				// 退部分押金
				if *req.DepositAmount != depositOrder.Amount {
					// 剩余金额
					remainAmount = tools.NewDecimal().Sub(depositOrder.Amount, *req.DepositAmount)
					// 退款金额
					depositOrder.Amount = *req.DepositAmount
					reason = "人工退押，部分退押"
				}
			}

			var or *ent.OrderRefund

			outRefundNo := tools.NewUnique().NewSN28()
			or, err = ent.Database.OrderRefund.Create().
				SetOrderID(depositOrder.ID).
				SetOutRefundNo(outRefundNo).
				SetAmount(depositOrder.Amount).
				SetReason(reason).
				SetOrderID(depositOrder.ID).
				SetStatus(refundStatus).
				SetRemainAmount(remainAmount).
				SetNillableRemark(req.Remark).
				Save(s.ctx)
			if err != nil {
				return err
			}

			// 更新订单状态
			_, err = depositOrder.Update().SetStatus(orderStatus).Save(s.ctx)
			if err != nil {
				return err
			}

			prepay := &model.PaymentCache{
				CacheType: model.PaymentCacheTypeRefund,
				Refund: &model.PaymentRefund{
					OrderID:      depositOrder.ID,
					TradeNo:      depositOrder.TradeNo,
					Total:        depositOrder.Total,
					RefundAmount: or.Amount,
					Reason:       or.Reason,
					OutRefundNo:  or.OutRefundNo,
				},
			}

			var no string
			no = or.OutRefundNo

			// 预支付订单号
			if depositOrder.TradePayAt == nil && depositOrder.OutOrderNo != "" {
				no = depositOrder.OutOrderNo
			}

			err = cache.Set(s.ctx, no, prepay, 20*time.Minute).Err()
			if err != nil {
				return err
			}

			// 当为支付押金时退款
			if (depositOrder.Payway == model.OrderPaywayAlipay || depositOrder.Payway == model.OrderPaywayWechat) && *req.RefundDeposit {
				switch depositOrder.Payway {
				case model.OrderPaywayAlipay:
					alipay.NewApp().Refund(prepay.Refund)
				case model.OrderPaywayWechat:
					wechat.NewApp().Refund(prepay.Refund)
				default:
					return fmt.Errorf("不支持的支付方式")
				}

				if prepay.Refund.Success { // 退款支付宝同步返回
					NewOrder().RefundSuccess(prepay.Refund)
				}
			}

			// 预授权 退款未转支付的订单
			if depositOrder.Payway == model.OrderPaywayAlipayAuthFreeze && depositOrder.TradePayAt == nil {
				if *req.RefundDeposit {
					// 退款 解冻押金
					err = NewOrder().FandAuthUnfreeze(prepay.Refund, depositOrder)
					// 如果退部分押金 还有部分要退的金额转支付
					if remainAmount > 0 {
						depositOrder.Amount = remainAmount
						err = NewOrder().TradePay(depositOrder)
					}
				} else {
					// 不退款 押金转支付
					err = NewOrder().TradePay(depositOrder)
				}

				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// NoBatteryUnsubscribe 单电无电池退租 只能由业务后台操作
func (s *businessRiderService) NoBatteryUnsubscribe(req *model.BusinessSubscribeReq, sub *ent.Subscribe) {
	if s.modifier == nil {
		snag.Panic("无权限操作")
	}

	_, err := ent.Database.Subscribe.
		UpdateOneID(sub.ID).
		SetEndAt(time.Now()).
		SetStatus(model.SubscribeStatusUnSubscribed).
		SetUnsubscribeReason("无电池后台强制退租").
		Save(s.ctx)
	snag.PanicIfError(err)

	// 查询并标记用户合同为失效
	_, err = ent.Database.Contract.Update().Where(contract.RiderID(sub.RiderID)).SetEffective(false).Save(s.ctx)
	snag.PanicIfError(err)

	// 更新企业账单
	if sub.EnterpriseID != nil {
		go NewEnterprise().UpdateStatementByID(*sub.EnterpriseID)
	}

	// 处理退款
	err = s.ForceUnsubscribe(req, sub.ID, sub.RiderID)
	if err != nil {
		zap.L().Error("退款失败", zap.Error(err))
		snag.Panic(err)
	}
}
