// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/assistance"
    "github.com/auroraride/aurservd/internal/ent/branch"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
)

type assistanceService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.AssistanceClient
}

func NewAssistance() *assistanceService {
    return &assistanceService{
        ctx: context.Background(),
        orm: ar.Ent.Assistance,
    }
}

func NewAssistanceWithRider(r *ent.Rider) *assistanceService {
    s := NewAssistance()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewAssistanceWithModifier(m *model.Modifier) *assistanceService {
    s := NewAssistance()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewAssistanceWithEmployee(e *ent.Employee) *assistanceService {
    s := NewAssistance()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *assistanceService) Breakdown() any {
    return NewSetting().GetSetting(model.SettingRescueReason)
}

// Unpaid 是否有未支付的救援订单
func (s *assistanceService) Unpaid(riderID uint64) *ent.Assistance {
    ass, _ := s.orm.QueryNotDeleted().
        Where(assistance.Status(model.AssistanceStatusSuccess), assistance.CostGT(0), assistance.PayAtIsNil(), assistance.RiderID(riderID)).
        First(s.ctx)
    return ass
}

// Create 发起救援订单
// TODO 救援订单未支付的禁止办理所有业务
// TODO 救援订单支付状态可以直接在后台修改为不需要支付
func (s *assistanceService) Create(req *model.AssistanceCreateReq) model.AssistanceCreateRes {
    sub := NewSubscribe().Recent(s.rider.ID)
    if sub.Status != model.SubscribeStatusUsing {
        snag.Panic("无法发起救援")
    }

    // 检查是否可发起救援
    NewRiderPermissionWithRider(s.rider).BusinessX()

    // 检查是否已有救援订单
    if exists, _ := s.orm.QueryNotDeleted().Where(assistance.RiderID(s.rider.ID)).Exist(s.ctx); exists {
        snag.Panic("当前有进行中的救援订单")
    }

    as, _ := s.orm.Create().
        SetStatus(model.AssistanceStatusPending).
        SetLng(req.Lng).
        SetLat(req.Lat).
        // SetDistance(haversine.Distance(haversine.NewCoordinates(req.Lat, req.Lng), haversine.NewCoordinates(stb.Lat, stb.Lng)).Miles()).
        SetAddress(req.Address).
        SetBreakdown(req.Breakdown).
        SetBreakdownPhotos(req.BreakdownPhotos).
        SetBreakdownDesc(req.BreakdownDesc).
        SetOutTradeNo(tools.NewUnique().NewSN28()).
        SetRiderID(s.rider.ID).
        SetSubscribeID(sub.ID).
        Save(s.ctx)

    if as == nil {
        snag.Panic("救援发起失败")
    }

    return model.AssistanceCreateRes{OutTradeNo: as.OutTradeNo}
}

func (s *assistanceService) Nearby(req *model.IDQueryReq) any {
    ass, _ := s.orm.QueryNotDeleted().Where(assistance.ID(req.ID)).First(s.ctx)
    if ass == nil {
        snag.Panic("未找到救援订单")
    }

    var btemps []struct {
        ID       uint64          `json:"id"`
        Address  string          `json:"address"`
        Lat      float64         `json:"lat"`
        Lng      float64         `json:"lng"`
        Distance float64         `json:"distance"`
        Edges    ent.BranchEdges `json:"edges"`
    }

    _ = ar.Ent.Branch.QueryNotDeleted().
        WithStores().
        Where(branch.HasStoresWith(store.HasEmployee(), store.Status(model.StoreStatusOpen))).
        Modify(func(sel *sql.Selector) {
            bt := sql.Table(branch.Table)
            sel.Select(bt.C(branch.FieldID), bt.C(branch.FieldAddress), bt.C(branch.FieldLat), bt.C(branch.FieldLng)).
                Where(sql.P(func(b *sql.Builder) {
                    b.WriteString(fmt.Sprintf(`ST_DWithin(%s, ST_GeogFromText('POINT(%f %f)'), %f)`, branch.FieldGeom, ass.Lng, ass.Lat, 50000.0))
                })).
                AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(%s, ST_GeogFromText('POINT(%f %f)'))`, branch.FieldGeom, ass.Lng, ass.Lat)), "distance").
                GroupBy(bt.C(branch.FieldID)).
                OrderBy(sql.Asc("distance"))
        }).
        Scan(s.ctx, &btemps)

    return nil
}
