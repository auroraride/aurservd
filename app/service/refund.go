// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"time"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/orderrefund"
	"github.com/auroraride/aurservd/internal/payment/alipay"
	"github.com/auroraride/aurservd/internal/payment/wechat"
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
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, r)
	s.rider = r
	return s
}

func NewRefundWithModifier(m *model.Modifier) *refundService {
	s := NewRefund()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
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
		s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
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
	if o == nil || o.TradeNo == "" && o.AuthNo == "" {
		snag.Panic("原始订单查询失败")
	}

	// 退款处理
	s.DoRefund(o, or, req.Status, req.Remark)

	// 退押金
	depositOrder := ent.Database.Order.
		QueryNotDeleted().
		Where(
			order.Or(
				order.ParentID(o.ID),
				order.ID(o.ID),
			),
			order.Type(model.OrderTypeDeposit),
			order.Status(model.OrderStatusRefundPending),
		).
		WithRefund().
		FirstX(s.ctx)
	if depositOrder != nil && depositOrder.Edges.Refund != nil {
		// 判定押金退款状态
		if depositOrder.Edges.Refund.Status == model.RefundStatusPending {
			s.DoRefund(depositOrder, depositOrder.Edges.Refund, req.Status, req.Remark)
		}
	}
}

// DoRefund  退款处理
func (s *refundService) DoRefund(o *ent.Order, or *ent.OrderRefund, status uint8, remark string) {
	var os uint8

	after := "同意退款"
	if o.Type == model.OrderTypeDeposit {
		after = "订阅已退订，系统自动退押"
	}

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

	if status == 1 {

		// 订单缓存key (退款单号)
		var no string
		no = or.OutRefundNo

		// 预支付订单号
		if o.TradePayAt == nil && o.OutOrderNo != "" {
			no = o.OutOrderNo
		}

		err := cache.Set(s.ctx, no, prepay, 20*time.Minute).Err()
		if err != nil {
			snag.Panic("退款处理失败")
		}

		// 处理退款
		switch o.Payway {
		case model.OrderPaywayAlipay:
			alipay.NewApp().Refund(prepay.Refund)
		case model.OrderPaywayWechat:
			wechat.NewApp().Refund(prepay.Refund)
		case model.OrderPaywayAlipayAuthFreeze:
			// 芝麻免押参数不一样
			var isDeposit bool
			if o.Type == model.OrderTypeDeposit {
				isDeposit = true
			}
			// 如果只是预支付直接解冻
			err = alipay.NewApp().FandAuthUnfreeze(prepay.Refund, &definition.FandAuthUnfreezeReq{
				AuthNo:       o.AuthNo,
				Amount:       or.Amount,
				OutRequestNo: or.OutRefundNo,
				Remark:       or.Reason,
				IsDeposit:    isDeposit,
			})
			if err != nil {
				snag.Panic("退款处理失败")
			}
		case model.OrderPaywayAlipayMiniProgram:
			alipay.NewMiniProgram().Refund(prepay.Refund)
		default:
			snag.Panic("退款处理失败")
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
			_, err = tx.OrderRefund.UpdateOne(or).SetStatus(status).SetRemark(remark).Save(s.ctx)
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
