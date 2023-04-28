// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"time"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/orderrefund"
	"github.com/auroraride/aurservd/internal/payment"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type refundService struct {
	ctx          context.Context
	modifier     *model.Modifier
	rider        *ent.Rider
	employee     *ent.Employee
	employeeInfo *model.Employee
}

func NewRefund() *refundService {
	return &refundService{
		ctx: context.Background(),
	}
}

func NewRefundWithRider(r *ent.Rider) *refundService {
	s := NewRefund()
	s.ctx = context.WithValue(s.ctx, "rider", r)
	s.rider = r
	return s
}

func NewRefundWithModifier(m *model.Modifier) *refundService {
	s := NewRefund()
	s.ctx = context.WithValue(s.ctx, "modifier", m)
	s.modifier = m
	return s
}

func NewRefundWithEmployee(e *ent.Employee) *refundService {
	s := NewRefund()
	if e != nil {
		s.employee = e
		s.employeeInfo = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
		s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
	}
	return s
}

// Refund 退款申请
func (s *refundService) Refund(riderID uint64, req *model.RefundReq) (res model.RefundRes) {
	if (req.Deposit == nil && req.SubscribeID == nil) || (req.Deposit != nil && req.SubscribeID != nil) {
		snag.Panic("退款参数错误")
	}

	sub := NewSubscribe().Recent(riderID)
	no := tools.NewUnique().NewSN28()

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		var id uint64
		orc := tx.OrderRefund.Create().SetOutRefundNo(no).SetStatus(model.RefundStatusPending)
		if req.SubscribeID != nil {
			if sub.Edges.InitialOrder == nil {
				snag.Panic("未找到有效订单")
			}
			if sub.Edges.InitialOrder.Status != model.OrderStatusPaid {
				snag.Panic("订单状态错误")
			}
			id = sub.Edges.InitialOrder.ID
			if sub.ID != *req.SubscribeID {
				snag.Panic("未找到有效骑士卡")
			}
			if sub.Status != model.SubscribeStatusInactive {
				snag.Panic("骑士卡已激活, 无法退款")
			}
			orc.SetOrderID(sub.InitialOrderID).
				SetAmount(sub.Edges.InitialOrder.Amount).
				SetReason("骑士卡未激活, 用户申请")
		}

		if req.Deposit != nil {
			// 押金退款逻辑
			if sub.Status != model.SubscribeStatusUnSubscribed && sub.Status != model.SubscribeStatusCanceled {
				snag.Panic("骑士卡正在使用中, 无法退押金")
			}

			// 获取骑手押金订单
			o := NewRider().DepositOrder(s.rider.ID)

			if o == nil {
				snag.Panic("未找到押金订单")
			}

			id = o.ID
			if o.Status != model.OrderStatusPaid {
				snag.Panic("订单状态错误")
			}

			orc.SetOrderID(o.ID).SetAmount(o.Amount).SetReason("无使用中的骑士卡, 用户申请退还押金")
		}

		_, err = orc.Save(s.ctx)
		if err != nil {
			snag.Panic("退款申请失败")
		}

		// 更新订单
		_, err = tx.Order.UpdateOneID(id).SetStatus(model.OrderStatusRefundPending).Save(s.ctx)
		return
	})

	return model.RefundRes{OutRefundNo: no}
}

// RefundAudit 退款审核
func (s *refundService) RefundAudit(req *model.RefundAuditReq) {
	// 查询退款订单
	or, _ := ent.Database.OrderRefund.QueryNotDeleted().Where(orderrefund.OutRefundNo(req.OutRefundNo)).WithOrder().First(s.ctx)

	if or == nil {
		snag.Panic("未找到退款订单")
	}

	if or.Status != model.RefundStatusPending {
		snag.Panic("退款申请已处理过")
	}

	if req.Status == model.RefundStatusRefused && req.Remark == "" {
		snag.Panic("拒绝理由不能为空")
	}

	// 原始订单
	o := or.Edges.Order
	if o == nil || o.TradeNo == "" {
		snag.Panic("原始订单查询失败")
	}

	status := req.Status
	var os uint8
	after := "同意退款"

	prepay := &model.PaymentCache{
		CacheType: model.PaymentCacheTypeRefund,
		Refund: &model.PaymentRefund{
			OrderID:      o.ID,
			TradeNo:      o.TradeNo,
			Total:        o.Total,
			RefundAmount: or.Amount,
			Reason:       or.Reason,
			OutRefundNo:  or.OutRefundNo,
		},
	}

	if req.Status == 1 {

		// 订单缓存 (原始订单号key)
		err := cache.Set(s.ctx, o.OutTradeNo, prepay, 20*time.Minute).Err()
		if err != nil {
			snag.Panic("退款处理失败")
		}

		// 处理退款
		switch o.Payway {
		case model.OrderPaywayAlipay:
			payment.NewAlipay().Refund(prepay.Refund)
			break
		case model.OrderPaywayWechat:
			payment.NewWechat().Refund(prepay.Refund)
			break
		}

		// 原路退款请求是否成功
		if !prepay.Refund.Request {
			status = model.RefundStatusFail
		} else {
			os = model.OrderStatusRefundSuccess
		}
	} else {
		os = model.OrderStatusPaid
		after = "拒绝退款"
	}

	// 判定退款是否成功 (支付宝同步返回有可能比异步慢)
	if prepay.Refund.Success {
		NewOrder().RefundSuccess(prepay.Refund)
	} else {
		ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
			_, err = tx.OrderRefund.UpdateOne(or).SetStatus(status).SetRemark(req.Remark).Save(s.ctx)
			if err != nil {
				return
			}

			_, err = tx.Order.UpdateOne(o).SetStatus(os).Save(s.ctx)
			return
		})
	}

	// 记录日志
	go logging.NewOperateLog().
		SetRef(o).
		SetModifier(s.modifier).
		SetOperate(model.OperateRefund).
		SetDiff("退款处理中", after).
		Send()
}
