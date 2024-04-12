// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-13
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/business"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/reserve"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type reserveService struct {
	ctx      context.Context
	modifier *model.Modifier
	rider    *ent.Rider
	orm      *ent.ReserveClient
	max      time.Duration
}

func NewReserve() *reserveService {
	max := time.Duration(cache.Int(model.SettingReserveDurationKey))
	return &reserveService{
		ctx: context.Background(),
		orm: ent.Database.Reserve,
		max: max,
	}
}

func NewReserveWithRider(r *ent.Rider) *reserveService {
	s := NewReserve()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, r)
	s.rider = r
	return s
}

func NewReserveWithModifier(m *model.Modifier) *reserveService {
	s := NewReserve()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func (s *reserveService) Query(id uint64) (*ent.Reserve, error) {
	return s.orm.QueryNotDeleted().Where(reserve.ID(id)).First(s.ctx)
}

func (s *reserveService) QueryX(id uint64) *ent.Reserve {
	rev, _ := s.Query(id)
	if rev == nil {
		snag.Panic("未找到预约")
	}
	return rev
}

// RiderUnfinished 查询用户未完成的预约
func (s *reserveService) RiderUnfinished(riderID uint64) *ent.Reserve {
	rev, _ := s.orm.QueryNotDeleted().
		Where(
			reserve.StatusIn(model.ReserveStatusPending.Value(), model.ReserveStatusProcessing.Value()),
			reserve.CreatedAtGTE(time.Now().Add(-s.max*time.Minute)),
			reserve.RiderID(riderID),
		).
		First(s.ctx)
	return rev
}

func (s *reserveService) CabinetUnfinished(cabinetID uint64) []*ent.Reserve {
	items, _ := s.orm.QueryNotDeleted().
		WithRider().
		Where(
			reserve.StatusIn(model.ReserveStatusPending.Value(), model.ReserveStatusProcessing.Value()),
			reserve.CreatedAtGTE(time.Now().Add(-s.max*time.Minute)),
			reserve.CabinetID(cabinetID),
		).
		All(s.ctx)
	if items == nil {
		items = make([]*ent.Reserve, 0)
	}
	return items
}

func (s *reserveService) RiderUnfinishedDetail(riderID uint64) *model.ReserveUnfinishedRes {
	rev := s.RiderUnfinished(riderID)
	if rev == nil {
		return nil
	}
	return s.Detail(rev)
}

func (s *reserveService) Detail(rev *ent.Reserve) *model.ReserveUnfinishedRes {
	return &model.ReserveUnfinishedRes{
		ID:        rev.ID,
		CabinetID: rev.CabinetID,
		Business:  rev.Type,
		Time:      rev.CreatedAt.Format(carbon.DateTimeLayout),
		Status:    model.ReserveStatus(rev.Status),
		Fid:       NewBranch().EncodeCabinetID(rev.CabinetID),
	}
}

// Timeout 标记任务超时, 条件为: 一个小时之前未开始
func (s *reserveService) Timeout() {
	_, _ = s.orm.Update().
		Where(
			reserve.Status(model.ReserveStatusPending.Value()),
			reserve.CreatedAtLT(time.Now().Add(-s.max*time.Minute)),
		).
		SetStatus(model.ReserveStatusTimeout.Value()).
		Save(s.ctx)
	// log.Infof("标记任务超时数量: %d", x)
}

// CabinetCounts 获取电柜当前预约数量
// 条件: 执行中的预约或一个小时以内的未执行预约
func (s *reserveService) CabinetCounts(ids []uint64, typ business.Type) (data map[uint64]int) {
	var results []struct {
		CabinetID uint64 `json:"cabinet_id"`
		Count     int    `json:"count"`
	}
	_ = s.orm.QueryNotDeleted().
		Where(
			reserve.CabinetIDIn(ids...),
			reserve.Or(
				reserve.And(
					reserve.Status(model.ReserveStatusPending.Value()),
					reserve.CreatedAtGTE(time.Now().Add(-s.max*time.Minute)),
					reserve.Type(typ.String()),
				),
				reserve.Status(model.ReserveStatusProcessing.Value()),
			),
		).
		GroupBy(reserve.FieldCabinetID).
		Aggregate(ent.Count()).
		Scan(s.ctx, &results)
	data = make(map[uint64]int)
	for _, result := range results {
		data[result.CabinetID] = result.Count
	}
	return
}

// Create 创建预约
func (s *reserveService) Create(req *model.ReserveCreateReq) *model.ReserveUnfinishedRes {
	typ := business.Type(req.Business)

	// 检查订阅状态
	sub := NewSubscribeWithRider(s.rider).RecentX(s.rider.ID)

	// 检查骑手权限以及是否可办理业务
	NewRiderPermissionWithRider(s.rider).BusinessX().SubscribeX(model.RiderPermissionTypeBusiness, sub)

	// 骑士卡状态
	if !NewRiderBusiness(s.rider).Executable(sub, typ) {
		snag.Panic("骑士卡状态错误")
	}

	// 判断电柜是否可预约
	cab := NewCabinet().QueryOne(req.CabinetID)
	// 同步电柜并返回电柜详情
	NewCabinet().Sync(cab)
	m := s.CabinetCounts([]uint64{cab.ID}, typ)
	if !cab.ReserveAble(typ, m[cab.ID]) {
		snag.Panic("电柜无法预约")
	}

	// 判断骑手是否可预约
	if s.RiderUnfinished(s.rider.ID) != nil {
		snag.Panic("当前已有其他预约")
	}

	// 创建预约
	rev, err := s.orm.Create().
		SetRiderID(s.rider.ID).
		SetCabinetID(cab.ID).
		SetCityID(*cab.CityID).
		SetType(req.Business).
		Save(s.ctx)
	if err != nil {
		snag.Panic("预约失败")
	}

	return s.Detail(rev)
}

// Cancel 取消预约
func (s *reserveService) Cancel(req *model.IDParamReq) {
	_, _ = s.QueryX(req.ID).Update().SetStatus(model.ReserveStatusCancel.Value()).Save(s.ctx)
}

func (s *reserveService) listFilter(req model.ReserveListFilter) (q *ent.ReserveQuery, info ar.Map) {
	q = s.orm.QueryNotDeleted().WithCity().WithCabinet().WithRider().Order(ent.Desc(reserve.FieldCreatedAt))
	info = make(ar.Map)
	if req.CityID != 0 {
		q.Where(reserve.CityID(req.CityID))
		info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
	}
	if req.RiderID != 0 {
		q.Where(reserve.RiderID(req.RiderID))
		info["骑手"] = ent.NewExportInfo(req.RiderID, rider.Table)
	}
	if req.Serial != "" {
		q.Where(reserve.HasCabinetWith(cabinet.Serial(req.Serial)))
		info["电柜编码"] = req.Serial
	}
	if req.Start != "" {
		q.Where(reserve.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
		info["开始日期"] = req.Start
	}
	if req.End != "" {
		q.Where(reserve.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
		info["结束日期"] = req.End
	}
	if req.CabinetID != 0 {
		q.Where(reserve.CabinetID(req.CabinetID))
		info["电柜"] = ent.NewExportInfo(req.CabinetID, cabinet.Table)
	}
	if req.Business != "" {
		q.Where(reserve.Type(req.Business))
		info["业务"] = model.BusinessTypeText(req.Business)
	}
	return
}

func (s *reserveService) List(req *model.ReserveListReq) *model.PaginationRes {
	q, _ := s.listFilter(req.ReserveListFilter)
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Reserve) model.ReserveListRes {
		return s.listDetail(item)
	})
}

func (s *reserveService) listDetail(item *ent.Reserve) (res model.ReserveListRes) {
	res = model.ReserveListRes{
		City:      item.Edges.City.Name,
		Name:      item.Edges.Rider.Name,
		Phone:     item.Edges.Rider.Phone,
		Business:  model.BusinessTypeText(item.Type),
		Status:    model.ReserveStatus(item.Status).String(),
		CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
	}
	cab := item.Edges.Cabinet
	if cab != nil {
		res.CabinetName = cab.Name
		res.Serial = cab.Serial
	}
	return
}
