// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

// import (
//     "context"
//     "github.com/auroraride/aurservd/app/model"
//     "github.com/auroraride/aurservd/internal/ar"
//     "github.com/auroraride/aurservd/internal/ent"
//     "github.com/auroraride/aurservd/internal/ent/order"
//     "github.com/auroraride/aurservd/pkg/tools"
//     "github.com/golang-module/carbon/v2"
//     "time"
// )
//
// type riderOrderService struct {
//     ctx      context.Context
//     modifier *model.Modifier
//     rider    *ent.Rider
//     orm      *ent.OrderClient
// }
//
// func NewRiderOrder() *riderOrderService {
//     return &riderOrderService{
//         ctx: context.Background(),
//         orm: ar.Ent.Order,
//     }
// }
//
// func NewRiderOrderWithRider(rider *ent.Rider) *orderService {
//     s := NewOrder()
//     s.ctx = context.WithValue(s.ctx, "rider", rider)
//     s.rider = rider
//     return s
// }
//
// func NewRiderOrderWithModifier(m *model.Modifier) *riderOrderService {
//     s := NewRiderOrder()
//     s.ctx = context.WithValue(s.ctx, "modifier", m)
//     s.modifier = m
//     return s
// }
//
// func (s *riderOrderService) CalculateDays() {
//
// }
//
// // FindNotActived 查找用户未激活的新订单
// func (s *riderOrderService) FindNotActived(riderID uint64) *ent.Order {
//     o, _ := s.orm.QueryNotDeleted().Where(
//         order.RiderID(riderID),
//         order.Status(model.OrderStatusPaid),
//         order.TypeIn(model.OrderTypeNewPlan, model.OrderTypeRenewal),
//         order.StartAtIsNil(),
//         order.ParentIDIsNil(),
//     ).Only(s.ctx)
//     return o
// }
//
// // FindUsingOrder 查找使用中的骑士卡
// func (s *riderOrderService) FindUsingOrder(riderID uint64) *ent.Order {
//     o, _ := s.orm.QueryNotDeleted().
//         WithAlters().
//         WithPauses().
//         Where(
//             order.RiderID(riderID),
//             order.TypeIn(model.OrderTypeNewPlan, model.OrderTypeRenewal),
//             order.ParentIDIsNil(),
//             order.StartAtNotNil(),
//             order.EndAtIsNil(),
//         ).
//         Only(s.ctx)
//     return o
// }
//
// // Detail 计算骑手当前订单详情
// func (s *riderOrderService) Detail(o *ent.Order) *model.RiderRecentOrder {
//     now := time.Now()
//     tt := tools.NewTime()
//     ro := &model.RiderRecentOrder{
//         ID: o.ID,
//         // Days:   int(o.PlanDetail.Days),
//         PayAt:  o.CreatedAt.Format(carbon.DateTimeLayout),
//         Amount: o.Amount,
//         Total:  o.Amount,
//         Payway: o.Payway,
//         City: model.City{
//             ID:   o.Edges.City.ID,
//             Name: o.Edges.City.Name,
//         },
//         Plan: model.OrderPlan{
//             ID:   o.PlanID,
//             Name: o.PlanDetail.Name,
//         },
//     }
//
//     if len(o.Edges.Children) > 0 {
//         ro.Deposit = o.Edges.Children[0].Amount
//     }
//     ro.Total += ro.Deposit
//
//     for i, pm := range o.Edges.Plan.Edges.Pms {
//         if i == 0 {
//             ro.Voltage = pm.Voltage
//         }
//         ro.Models = append(ro.Models, model.BatteryModel{
//             ID:       pm.ID,
//             Voltage:  pm.Voltage,
//             Capacity: pm.Capacity,
//         })
//     }
//
//     // 计算改动天数
//     alter := 0
//     for _, oa := range o.Edges.Alters {
//         alter += oa.Days
//     }
//
//     // 计算暂停天数
//     paused := 0
//     for _, op := range o.Edges.Pauses {
//         if op.EndAt.IsZero() {
//             paused += tt.SubDays(now, op.StartAt)
//         } else {
//             paused += op.Days
//         }
//     }
//
//     // 距离激活已过去多少天
//     past := 0
//     if o.StartAt != nil {
//         ro.StartAt = o.StartAt.Format(carbon.DateLayout)
//         past = tt.SubDays(now, *o.StartAt)
//     }
//
//     remaining := int(o.Days) + alter + paused - past
//
//     switch true {
//     case o.Status == model.OrderStatusRefundSuccess:
//         // 已退款
//         ro.Status = model.RiderOrderStatusRefunded
//         ro.EndAt = o.RefundAt.Format(carbon.DateLayout)
//         remaining = 0
//         break
//     case o.StartAt == nil:
//         // 未激活
//         ro.Status = model.RiderOrderStatusPending
//         break
//     case o.PausedAt == nil && o.StartAt != nil && o.EndAt == nil:
//         if remaining >= 0 {
//             // 计费中
//             ro.Status = model.RiderOrderStatusNormal
//         } else {
//             ro.Status = model.RiderOrderStatusOverdue
//         }
//         break
//     case o.PausedAt != nil:
//         // 暂停中
//         ro.Status = model.RiderOrderStatusPaused
//         break
//     case o.EndAt != nil:
//         // 已归还
//         ro.Status = model.RiderOrderStatusRemanded
//         ro.EndAt = o.EndAt.Format(carbon.DateLayout)
//         remaining = 0
//         break
//     }
//
//     if remaining >= 0 && o.EndAt == nil {
//         ro.EndAt = time.Now().AddDate(0, 0, remaining).Format(carbon.DateLayout)
//     }
//     // ro.Remaining = remaining
//     // ro.PausedDays = paused
//     // ro.AlterDays = alter
//
//     return ro
// }
//
// // Recent 获取最近的订单
// func (s *riderOrderService) Recent(riderID uint64) *model.RiderRecentOrder {
//     o, _ := ar.Ent.Order.
//         QueryNotDeleted().
//         WithAlters().
//         WithCity().
//         WithPauses().
//         WithPlan(
//             func(pq *ent.PlanQuery) {
//                 pq.WithPms()
//             },
//         ).
//         Where(
//             order.RiderID(riderID),
//             order.StatusIn(model.OrderStatusPaid, model.OrderStatusRefundPending),
//             order.TypeIn(model.OrderRiderPlan...),
//             order.ParentIDIsNil(),
//         ).
//         WithChildren(func(cq *ent.OrderQuery) {
//             cq.Where(order.Type(model.OrderTypeDeposit))
//         }).
//         Only(s.ctx)
//     if o == nil {
//         return nil
//     }
//
//     return s.Detail(o)
// }
//
// // List 分页获取骑手订单
// func (s *riderOrderService) List(riderID uint64, req *model.PaginationReq) *model.PaginationRes {
//     if req.PageSize > 10 {
//         req.PageSize = 10
//     }
//     q := s.orm.QueryNotDeleted().
//         WithAlters().
//         WithCity().
//         WithPauses().
//         WithPlan(
//             func(pq *ent.PlanQuery) {
//                 pq.WithPms()
//             },
//         ).
//         Order(ent.Desc(order.FieldCreatedAt)).
//         Where(
//             order.RiderID(riderID),
//             order.TypeIn(model.OrderRiderPlan...),
//         )
//     return model.ParsePaginationResponse[model.RiderRecentOrder](
//         q.PaginationResult(*req),
//         func() []model.RiderRecentOrder {
//             items := q.Pagination(*req).AllX(s.ctx)
//             out := make([]model.RiderRecentOrder, len(items))
//             for i, item := range items {
//                 out[i] = *s.Detail(item)
//             }
//             return out
//         },
//     )
// }
