// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/workwx"
	"github.com/auroraride/aurservd/internal/amap"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/baidu"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assistance"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/payment/alipay"
	"github.com/auroraride/aurservd/internal/payment/wechat"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assistanceService struct {
	ctx          context.Context
	modifier     *model.Modifier
	rider        *ent.Rider
	employee     *ent.Employee
	orm          *ent.AssistanceClient
	employeeInfo *model.Employee
}

func NewAssistance() *assistanceService {
	return &assistanceService{
		ctx: context.Background(),
		orm: ent.Database.Assistance,
	}
}

func NewAssistanceWithRider(r *ent.Rider) *assistanceService {
	s := NewAssistance()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, r)
	s.rider = r
	return s
}

func NewAssistanceWithModifier(m *model.Modifier) *assistanceService {
	s := NewAssistance()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func NewAssistanceWithEmployee(e *ent.Employee) *assistanceService {
	s := NewAssistance()
	if e != nil {
		s.employee = e
		s.employeeInfo = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
		s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
	}
	return s
}

func (s *assistanceService) Breakdown() any {
	return NewSetting().GetSetting(model.SettingRescueReasonKey)
}

// Unpaid 是否有未支付的救援订单
func (s *assistanceService) Unpaid(riderID uint64) *ent.Assistance {
	ass, _ := s.orm.QueryNotDeleted().
		Where(assistance.Status(model.AssistanceStatusUnpaid), assistance.RiderID(riderID)).
		First(s.ctx)
	return ass
}

// List 列举
func (s *assistanceService) List(req *model.AssistanceListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().
		WithRider().
		WithCity().
		WithStore().
		WithEmployee().
		Order(func(sel *sql.Selector) {
			sel.OrderExpr(sql.Raw(fmt.Sprintf(
				"POSITION(%s::text in '%d,%d,%d,%d,%d,%d')",
				sel.C(store.FieldStatus),
				model.AssistanceStatusPending,
				model.AssistanceStatusAllocated,
				model.AssistanceStatusUnpaid,
				model.AssistanceStatusFailed,
				model.AssistanceStatusRefused,
				model.AssistanceStatusSuccess,
			)))
		}).
		Unique(false).
		Order(ent.Desc(assistance.FieldCreatedAt))
	tt := tools.NewTime()
	if req.Start != "" {
		q.Where(assistance.CreatedAtGTE(tt.ParseDateStringX(req.Start)))
	}
	if req.End != "" {
		q.Where(assistance.CreatedAtLT(tt.ParseNextDateStringX(req.End)))
	}
	if req.CityID != 0 {
		q.Where(assistance.CityID(req.CityID))
	}
	if req.Keyword != "" {
		q.Where(
			assistance.HasRiderWith(rider.Or(
				rider.NameContainsFold(req.Keyword),
				rider.PhoneContainsFold(req.Keyword),
			)),
		)
	}
	if req.Status != nil {
		q.Where(assistance.Status(*req.Status))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Assistance) model.AssistanceListRes {
		return s.BasicInfo(item)
	})
}

// BasicInfo 基础信息
func (s *assistanceService) BasicInfo(item *ent.Assistance) model.AssistanceListRes {
	r := item.Edges.Rider
	c := item.Edges.City
	res := model.AssistanceListRes{
		ID:       item.ID,
		Status:   item.Status,
		Cost:     item.Cost,
		Distance: item.Distance,
		Time:     item.CreatedAt.Format(carbon.DateTimeLayout),
		Rider: model.Rider{
			ID:    r.ID,
			Phone: r.Phone,
			Name:  r.Name,
		},
		City: model.City{
			ID:   c.ID,
			Name: c.Name,
		},
	}

	e := item.Edges.Employee
	st := item.Edges.Store
	if e != nil {
		res.Employee = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
	}

	if st != nil {
		res.Store = &model.Store{
			ID:   st.ID,
			Name: st.Name,
		}
	}

	return res
}

func (s *assistanceService) Query(id uint64) (*ent.Assistance, error) {
	return s.orm.QueryNotDeleted().
		Where(assistance.ID(id)).
		First(s.ctx)
}

func (s *assistanceService) QueryX(id uint64) *ent.Assistance {
	item, _ := s.Query(id)
	if item == nil {
		snag.Panic("未找到救援信息")
	}
	return item
}

func (s *assistanceService) QueryDetail(id uint64) (*ent.Assistance, error) {
	return s.orm.QueryNotDeleted().
		WithRider().
		WithCity().
		WithStore().
		WithEmployee().
		Where(assistance.ID(id)).
		First(s.ctx)
}

func (s *assistanceService) QueryDetailX(id uint64) *ent.Assistance {
	item, _ := s.QueryDetail(id)
	if item == nil {
		snag.Panic("未找到救援信息")
	}
	return item
}

// Detail 救援详情
func (s *assistanceService) Detail(id uint64) model.AssistanceDetail {
	item := s.QueryDetailX(id)
	res := model.AssistanceDetail{
		AssistanceListRes: s.BasicInfo(item),
		Lng:               item.Lng,
		Lat:               item.Lat,
		Address:           item.Address,
		Breakdown:         item.Breakdown,
		BreakdownDesc:     item.BreakdownDesc,
		BreakdownPhotos:   item.BreakdownPhotos,
		Reason:            item.Reason,
		DetectPhoto:       item.DetectPhoto,
		JointPhoto:        item.JointPhoto,
		RefusedDesc:       item.RefusedDesc,
		FreeReason:        item.FreeReason,
		FailReason:        item.FailReason,
	}
	if item.PayAt != nil {
		res.PayAt = silk.Pointer(item.PayAt.Format(carbon.DateTimeLayout))
	}

	return res
}

// Current 获取当前进行中的救援
func (s *assistanceService) Current(riderID uint64) *ent.Assistance {
	item, _ := s.orm.QueryNotDeleted().
		Where(
			assistance.RiderID(riderID),
			assistance.StatusIn(model.AssistanceStatusProcessing...),
		).
		First(s.ctx)
	return item
}

// CurrentMessage 获取骑手进行中的救援信息
func (s *assistanceService) CurrentMessage(riderID uint64) *model.AssistanceSocketMessage {
	ass := s.Current(riderID)
	if ass == nil {
		return nil
	}
	message, err := NewAssistanceSocket().Detail(ass)
	if err != nil {
		return nil
	}
	return message
}

// Create 发起救援订单
// 救援订单未支付的禁止办理所有业务
// 救援订单支付状态可以直接在后台修改为不需要支付
func (s *assistanceService) Create(req *model.AssistanceCreateReq) model.AssistanceCreateRes {
	sub := NewSubscribe().Recent(s.rider.ID)

	// 检查是否可发起救援
	NewRiderPermissionWithRider(s.rider).BusinessX().SubscribeX(model.RiderPermissionTypeAssistance, sub)

	// 检查是否已有救援订单
	if exists := s.Current(s.rider.ID); exists != nil {
		snag.Panic("当前有进行中的救援订单")
	}

	rg, err := amap.New().ReGeo(req.Lng, req.Lat)
	if err != nil {
		snag.Panic("救援发起失败: " + err.Error())
	}
	adcode := rg.AddressComponent.Adcode
	cityId, _ := strconv.ParseUint(adcode[:len(adcode)-2]+"00", 10, 64)

	as, _ := s.orm.Create().
		SetStatus(model.AssistanceStatusPending).
		SetLng(req.Lng).
		SetLat(req.Lat).
		SetAddress(req.Address).
		SetBreakdown(req.Breakdown).
		SetBreakdownPhotos(req.BreakdownPhotos).
		SetBreakdownDesc(req.BreakdownDesc).
		SetRiderID(s.rider.ID).
		SetSubscribeID(sub.ID).
		SetCityID(cityId).
		Save(s.ctx)

	if as == nil {
		snag.Panic("救援发起失败")
	}

	go workwx.New().SendAssistance(model.AssistanceNotice{
		Phone:         s.rider.Phone,
		Reason:        req.Breakdown,
		Address:       rg.FormattedAddress,
		AddressDetail: req.Address,
	})

	return model.AssistanceCreateRes{ID: as.ID}
}

// Nearby 救援订单附近门店
func (s *assistanceService) Nearby(req *model.IDQueryReq) (items []model.AssistanceNearbyRes) {
	ass, _ := s.orm.QueryNotDeleted().Where(assistance.ID(req.ID)).First(s.ctx)
	if ass == nil {
		snag.Panic("未找到救援订单")
	}

	var temps []struct {
		ID       uint64  `json:"id"`
		Name     string  `json:"name"`
		Lng      float64 `json:"lng"`
		Lat      float64 `json:"lat"`
		Distance float64 `json:"distance"`
		EID      uint64  `json:"eid"`
		EName    string  `json:"e_ame"`
		EPhone   string  `json:"e_hone"`
	}

	err := ent.Database.Store.QueryNotDeleted().
		Where(store.EmployeeIDNotNil(), store.Status(model.StoreStatusOpen.Value())).
		Modify(
			func(sel *sql.Selector) {
				sel.Select(sel.C(store.FieldID), sel.C(store.FieldLng), sel.C(store.FieldLat), sel.C(store.FieldName)).
					Where(sql.P(func(b *sql.Builder) {
						b.WriteString(fmt.Sprintf(`ST_DWithin(ST_GeogFromText('SRID=4326;POINT('||%s||' '||%s||')'), ST_GeogFromText('POINT(%f %f)'), %f)`, sel.C(store.FieldLng), sel.C(store.FieldLat), ass.Lng, ass.Lat, model.DefaultMaxDistance))
					})).
					AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeogFromText('SRID=4326;POINT('||%s||' '||%s||')'), ST_GeogFromText('POINT(%f %f)'))`, sel.C(store.FieldLng), sel.C(store.FieldLat), ass.Lng, ass.Lat)), "distance").
					OrderBy(sql.Asc("distance"))

				// 查找employee
				emt := sql.Table(employee.Table)
				sel.Join(emt).
					On(sel.C(store.FieldEmployeeID), emt.C(employee.FieldID)).
					AppendSelect(
						sql.As(emt.C(employee.FieldID), "eid"),
						sql.As(emt.C(employee.FieldName), "e_ame"),
						sql.As(emt.C(employee.FieldPhone), "e_hone"),
					)

				// // 查找city
				// ct := sql.Table(city.Table)
				// sel.Join(ct).
				//     On(sel.C(store.FieldCityID), ct.C(city.FieldID)).
				//     AppendSelect(
				//         sql.As(ct.C(city.FieldID), "cid"),
				//         sql.As(ct.C(city.FieldName), "c_name"),
				//     )
			},
		).
		Scan(s.ctx, &temps)

	if err != nil {
		zap.L().Error("门店查询失败", zap.Error(err))
	}

	// 查找执行中的救援员工ID
	ids, _ := s.orm.QueryNotDeleted().
		Where(assistance.Status(model.AssistanceStatusAllocated), assistance.EmployeeIDNotNil()).
		Select(assistance.FieldEmployeeID).
		Ints(s.ctx)
	idsm := make(map[uint64]bool)
	for _, id := range ids {
		idsm[uint64(id)] = true
	}

	items = make([]model.AssistanceNearbyRes, 0)

	for _, temp := range temps {
		if idsm[temp.EID] {
			continue
		}
		items = append(items, model.AssistanceNearbyRes{
			ID:       temp.ID,
			Name:     temp.Name,
			Lng:      temp.Lng,
			Lat:      temp.Lat,
			Distance: temp.Distance,
			Employee: model.Employee{
				ID:    temp.EID,
				Name:  temp.EName,
				Phone: temp.EPhone,
			},
		})
	}

	return
}

// Allocate 分配救援任务
func (s *assistanceService) Allocate(req *model.AssistanceAllocateReq) {
	st, _ := ent.Database.Store.QueryNotDeleted().
		Where(
			store.ID(req.StoreID),
			store.EmployeeIDNotNil(),
			store.Status(model.StoreStatusOpen.Value()),
			store.HasEmployeeWith(
				employee.Not(employee.HasAssistancesWith(
					assistance.Status(model.AssistanceStatusAllocated),
				)),
			),
		).
		First(s.ctx)
	if st == nil {
		snag.Panic("未找到营业中的门店")
	}

	item, _ := s.orm.QueryNotDeleted().Where(assistance.ID(req.ID), assistance.StatusIn(model.AssistanceStatusPending, model.AssistanceStatusAllocated)).First(s.ctx)
	if item == nil {
		snag.Panic("未找到有效救援")
	}
	before := item.Status

	var err error

	duration, distance, polylines := baidu.NewMap().RidingPlanX(fmt.Sprintf("%f,%f", st.Lat, st.Lng), fmt.Sprintf("%f,%f", item.Lat, item.Lng))

	item, err = item.Update().
		SetDistance(float64(distance)).
		SetNaviDuration(duration).
		SetNaviPolylines(polylines).
		SetStoreID(st.ID).
		SetEmployeeID(*st.EmployeeID).
		SetStatus(model.AssistanceStatusAllocated).
		SetAllocateAt(time.Now()).
		SetWait(int(time.Since(item.CreatedAt).Seconds())).
		Save(s.ctx)

	if err != nil {
		snag.Panic("分配失败")
	}

	// 救援处理接单响应
	// 发送信息给骑手
	go NewAssistanceSocket().SendRider(item.RiderID, item)

	// 发送消息给门店
	go NewAssistanceSocket().SenderEmployee(*st.EmployeeID, item)

	// 记录日志
	go logging.NewOperateLog().
		SetRef(item).
		SetModifier(s.modifier).
		SetOperate(model.OperateAssistanceAllocate).
		SetDiff(model.AssistanceStatus(before), model.AssistanceStatus(item.Status)).
		Send()
}

// Free 救援免费
func (s *assistanceService) Free(req *model.AssistanceFreeReq) {
	item, _ := s.orm.QueryNotDeleted().Where(assistance.Status(model.AssistanceStatusUnpaid)).First(s.ctx)
	if item == nil {
		snag.Panic("未找到待支付订单")
	}

	before := fmt.Sprintf("%s (%.2f元)", model.AssistanceStatus(item.Status), item.Cost)
	var err error

	item, err = item.Update().
		SetFreeReason(req.Reason).
		SetCost(0).
		SetStatus(model.AssistanceStatusSuccess).
		Save(s.ctx)
	if err != nil {
		snag.Panic("处理失败")
	}

	// 救援处理免费响应
	// 发送信息给骑手
	go NewAssistanceSocket().SendRider(item.RiderID, item)

	// 记录日志
	go logging.NewOperateLog().
		SetRef(item).
		SetModifier(s.modifier).
		SetOperate(model.OperateAssistanceFree).
		SetDiff(
			before,
			fmt.Sprintf("%s (%s)", model.AssistanceStatus(model.AssistanceStatusSuccess), req.Reason),
		).
		Send()
}

// Refuse 拒绝救援
func (s *assistanceService) Refuse(req *model.AssistanceRefuseReq) {
	item := s.QueryX(req.ID)
	if item.Status == model.AssistanceStatusSuccess || item.Status == model.AssistanceStatusUnpaid {
		snag.Panic("救援状态错误")
	}

	before := model.AssistanceStatus(item.Status)

	var err error
	item, err = item.Update().
		ClearEmployeeID().
		ClearStoreID().
		ClearCost().
		SetStatus(model.AssistanceStatusRefused).
		SetRefusedDesc(req.Reason).
		Save(s.ctx)

	if err != nil {
		snag.Panic("操作失败")
	}

	// 救援处理拒绝响应
	// 发送信息给骑手
	go NewAssistanceSocket().SendRider(item.RiderID, item)

	// 记录日志
	go logging.NewOperateLog().
		SetRef(item).
		SetModifier(s.modifier).
		SetOperate(model.OperateAssistanceRefuse).
		SetDiff(
			before,
			fmt.Sprintf("%s (%s)", model.AssistanceStatus(model.AssistanceStatusRefused), req.Reason),
		).
		Send()
}

// Cancel 取消救援
func (s *assistanceService) Cancel(req *model.AssistanceCancelReq) {
	// 查找未分配的救援
	item, _ := s.orm.QueryNotDeleted().Where(assistance.ID(req.ID), assistance.RiderID(s.rider.ID)).First(s.ctx)
	if item == nil {
		snag.Panic("未找到救援信息")
	}
	if item.Status != model.AssistanceStatusPending {
		snag.Panic("救援状态错误")
	}
	_, err := s.orm.SoftDeleteOne(item).
		SetCancelReason(req.Reason).
		SetCancelReasonDesc(req.Desc).
		ClearCost().
		Save(s.ctx)
	if err != nil {
		snag.Panic("取消失败:" + err.Error())
	}
}

// EmployeeDetail 门店端显示救援详情
func (s *assistanceService) EmployeeDetail(id uint64) (res model.AssistanceEmployeeDetailRes) {
	ass, _ := s.orm.QueryNotDeleted().
		Where(assistance.ID(id), assistance.EmployeeID(s.employee.ID)).
		WithRider().
		WithStore().
		WithSubscribe().
		First(s.ctx)
	if ass == nil {
		snag.Panic("未找到救援信息")
	}
	_ = copier.Copy(&res, ass)

	r := ass.Edges.Rider
	if r == nil {
		return
	}
	res.Rider = model.Rider{
		ID:    r.ID,
		Phone: r.Phone,
		Name:  r.Name,
	}

	// 救援原因
	res.Configure.Breakdown = s.Breakdown().([]interface{})

	st := ass.Edges.Store
	if st != nil {
		res.Store = model.StoreLngLat{
			Store: model.Store{
				ID:   st.ID,
				Name: st.Name,
			},
			Lng: st.Lng,
			Lat: st.Lat,
		}
	}

	sub := ass.Edges.Subscribe
	if sub != nil {
		res.Model = sub.Model
	}

	res.Time = ass.CreatedAt.Format(carbon.DateTimeLayout)
	if ass.Distance > 0 {
		res.Distance = fmt.Sprintf("%.2fkm", ass.Distance/1000.0)
	}

	if ass.ProcessAt != nil {
		res.Minutes = math.Round(ass.ProcessAt.Sub(ass.CreatedAt).Minutes())
	} else if ass.Status == model.AssistanceStatusAllocated {
		res.Minutes = math.Round(time.Since(ass.CreatedAt).Minutes())
	}

	for _, polyline := range ass.NaviPolylines {
		res.Polylines = append(res.Polylines, strings.Split(strings.TrimSpace(polyline), ";")...)
	}

	return
}

// Process 处理救援
func (s *assistanceService) Process(req *model.AssistanceProcessReq) (res model.AssistanceProcessRes) {
	ass, _ := s.orm.QueryNotDeleted().
		Where(assistance.ID(req.ID), assistance.EmployeeID(s.employee.ID), assistance.Status(model.AssistanceStatusAllocated)).
		WithRider().
		First(s.ctx)
	if ass == nil {
		snag.Panic("未找到有效救援信息")
	}

	up := ass.Update().SetProcessAt(time.Now())

	if req.Success {
		if req.Reason == "" || req.JointPhoto == "" || req.DetectPhoto == "" {
			snag.Panic("参数错误")
		}

		up.SetReason(req.Reason).
			SetJointPhoto(req.JointPhoto).
			SetDetectPhoto(req.DetectPhoto)

		status := model.AssistanceStatusUnpaid
		var cost, price float64

		if req.Pay {
			price = cache.Float64(model.SettingRescueFeeKey)
			cost = tools.NewDecimal().Mul(price, ass.Distance/1000.0)
		}

		if cost == 0 {
			status = model.AssistanceStatusSuccess
		}

		up.SetPrice(price).SetCost(cost).SetStatus(status)

		res.Cost = cost
	} else {
		if req.FailReason == "" {
			snag.Panic("参数错误")
		}
		up.SetStatus(model.AssistanceStatusFailed).SetFailReason(req.FailReason)
	}

	_, err := up.Save(s.ctx)
	if err != nil {
		snag.Panic("处理失败")
	}

	return
}

func (s *assistanceService) SimpleInfo(ass *ent.Assistance) model.AssistanceSimpleListRes {
	res := model.AssistanceSimpleListRes{
		ID:       ass.ID,
		Status:   ass.Status,
		Rider:    model.Rider{},
		Cost:     ass.Cost,
		Time:     ass.CreatedAt.Format(carbon.DateTimeLayout),
		Reason:   ass.Reason,
		Distance: fmt.Sprintf("%.2fkm", ass.Distance/1000.0),
	}

	sub := ass.Edges.Subscribe
	if sub != nil {
		res.Model = sub.Model
	}

	r := ass.Edges.Rider
	if r != nil {
		res.Rider = model.Rider{
			ID:    r.ID,
			Phone: r.Phone,
			Name:  r.Name,
		}
	}

	st := ass.Edges.Store
	if st != nil {
		res.Store = &model.Store{
			ID:   st.ID,
			Name: st.Name,
		}
	}

	return res
}

// Pay 支付救援
func (s *assistanceService) Pay(req *model.AssistancePayReq) model.AssistancePayRes {
	ass, _ := s.orm.QueryNotDeleted().
		Where(
			assistance.EmployeeID(s.employee.ID),
			assistance.Status(model.AssistanceStatusUnpaid),
			assistance.CostGT(0),
			assistance.ID(req.ID),
		).
		WithRider().
		WithSubscribe().
		WithStore().
		First(s.ctx)

	if ass == nil {
		snag.Panic("未找到待支付救援详情")
	}

	no := tools.NewUnique().NewSN()
	cost := ass.Cost

	// Development模式支付一分钱
	if ar.Config.Environment.IsDevelopment() {
		cost = 0.01
	}

	pc := &model.PaymentCache{
		CacheType: model.PaymentCacheTypeAssistance,
		Assistance: &model.PaymentAssistance{
			Subject:    fmt.Sprintf("%.2f公里救援", ass.Distance/1000),
			Payway:     req.Payway,
			Cost:       cost,
			OutTradeNo: no,
			ID:         ass.ID,
		},
	}

	var qr string
	var err error

	// 生成预支付订单
	switch req.Payway {
	case model.OrderPaywayAlipay:
		qr, err = alipay.NewApp().Native(pc)
	case model.OrderPaywayWechat:
		qr, err = wechat.NewApp().Native(pc)
	default:
		snag.Panic("unhandled default case")
	}

	if err != nil {
		snag.Panic(err)
	}

	if qr == "" {
		snag.Panic("支付二维码生成失败")
	}

	err = cache.Set(s.ctx, no, pc, 20*time.Minute).Err()
	if err != nil {
		snag.Panic("支付二维码生成失败")
	}

	res := model.AssistancePayRes{
		AssistanceSimpleListRes: s.SimpleInfo(ass),

		QR:         qr,
		OutTradeNo: no,
	}

	return res
}

// Paid 支付回调
func (s *assistanceService) Paid(trade *model.PaymentAssistance) {
	ass, _ := s.orm.QueryNotDeleted().Where(assistance.ID(trade.ID)).First(s.ctx)
	if ass == nil {
		return
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) error {
		o, err := tx.Order.Create().
			SetPayway(trade.Payway).
			SetAmount(trade.Cost).
			SetOutTradeNo(trade.OutTradeNo).
			SetTradeNo(trade.TradeNo).
			SetTotal(trade.Cost).
			SetCityID(ass.CityID).
			SetSubscribeID(ass.SubscribeID).
			SetType(model.OrderTypeAssistance).
			SetRiderID(ass.RiderID).
			Save(s.ctx)
		if err != nil {
			zap.L().Error("救援订单支付失败: "+trade.OutTradeNo, zap.Error(err))
			return err
		}

		_, err = tx.Assistance.UpdateOne(ass).SetStatus(model.AssistanceStatusSuccess).SetPayAt(time.Now()).SetOrderID(o.ID).Save(s.ctx)
		if err != nil {
			zap.L().Error("救援订单更新失败: "+trade.OutTradeNo+" -> "+strconv.FormatUint(ass.ID, 10), zap.Error(err))
			return err
		}
		return nil
	})
}

// SimpleList 简单列表
func (s *assistanceService) SimpleList(req model.PaginationReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Order(ent.Desc(assistance.FieldCreatedAt)).WithSubscribe().WithRider().WithStore()

	if s.rider != nil {
		q.Where(assistance.RiderID(s.rider.ID))
	}

	if s.employee != nil {
		q.Where(assistance.EmployeeID(s.employee.ID))
	}

	return model.ParsePaginationResponse(q, req, func(item *ent.Assistance) model.AssistanceSimpleListRes {
		return s.SimpleInfo(item)
	})
}

func (s *assistanceService) EmployeeOverview(employeeID uint64) model.AssistanceEmployeeOverview {
	items, _ := s.orm.QueryNotDeleted().
		Where(assistance.EmployeeID(employeeID)).
		All(s.ctx)
	res := model.AssistanceEmployeeOverview{Times: len(items)}
	for _, item := range items {
		if item.Status == model.AssistanceStatusUnpaid || item.Status == model.AssistanceStatusSuccess {
			res.Success += 1
			res.Distance += item.Distance / 1000.0
		}
	}
	res.Distance, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", res.Distance), 64)
	return res
}
