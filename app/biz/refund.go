package biz

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/orderrefund"
	"github.com/auroraride/aurservd/pkg/tools"
)

type refundBiz struct {
	orm *ent.OrderRefundClient
	ctx context.Context
}

func NewRefundBiz() *refundBiz {
	return &refundBiz{
		ctx: context.Background(),
		orm: ent.Database.OrderRefund,
	}
}

// Refund 申请退款
func (s *refundBiz) Refund(rId uint64, req *definition.RefundReq) (err error) {
	sub := service.NewSubscribe().Recent(rId)
	// 查询订单和押金
	o, _ := ent.Database.Order.QueryNotDeleted().Where(
		order.Or(
			order.And(
				order.ParentID(sub.InitialOrderID),
				order.Type(model.OrderTypeDeposit),
				order.Status(model.OrderStatusPaid),
			),
			order.And(
				order.IDEQ(sub.InitialOrderID),
				order.Status(model.OrderStatusPaid),
			),
		),
	).All(s.ctx)
	if len(o) == 0 {
		return errors.New("未找到有效订单")
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		if req.SubscribeID != nil {
			if sub.Edges.InitialOrder == nil {
				return errors.New("未找到有效订单")
			}
			if sub.Edges.InitialOrder.Status != model.OrderStatusPaid {
				return errors.New("订单状态异常")
			}

			if sub.ID != *req.SubscribeID {
				return errors.New("未找到有效骑士卡")
			}
			if sub.Status != model.SubscribeStatusInactive {
				return errors.New("骑士卡已激活, 无法退款")
			}

			for _, v := range o {
				or, _ := ent.Database.OrderRefund.
					QueryNotDeleted().
					Where(
						orderrefund.OrderID(v.ID),
						orderrefund.StatusIn(model.RefundStatusRefused, model.RefundStatusFail),
					).First(s.ctx)
				if or != nil {
					// 修改退款订单状态
					_, err = tx.OrderRefund.UpdateOneID(or.ID).SetStatus(model.RefundStatusPending).SetUpdatedAt(time.Now()).Save(s.ctx)
				} else {
					no := tools.NewUnique().NewSN28()
					orc := tx.OrderRefund.Create().SetOutRefundNo(no).SetStatus(model.RefundStatusPending)
					_, err = orc.SetOrderID(v.ID).
						SetAmount(v.Amount).
						SetReason("骑士卡未激活, 用户申请").
						Save(s.ctx)
					if err != nil {
						return errors.New("退款申请失败")
					}
				}
				// 更新订单
				_, err = tx.Order.UpdateOneID(v.ID).SetStatus(model.OrderStatusRefundPending).Save(s.ctx)
				if err != nil {
					return errors.New("退款申请失败")
				}
			}
		}
		return
	})
	return
}
