package service

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	pm "github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/purchaseorder"
	"github.com/auroraride/aurservd/internal/ent/purchasepayment"
	"github.com/auroraride/aurservd/internal/payment/alipay"
	"github.com/auroraride/aurservd/internal/payment/wechat"
	"github.com/auroraride/aurservd/pkg/tools"
)

type paymentService struct {
	orm *ent.PurchasePaymentClient
}

func NewPayment() *paymentService {
	return &paymentService{
		orm: ent.Database.PurchasePayment,
	}
}

// QueryByIndex 查询某个订单的某个支付计划
func (s *paymentService) QueryByIndex(ctx context.Context, orderID uint64, index int) (*ent.PurchasePayment, error) {
	return s.orm.QueryNotDeleted().Where(purchasepayment.OrderID(orderID), purchasepayment.Index(index)).First(ctx)
}

// Create 创建支付计划(支付计划以订单的分期计划为准)
func (s *paymentService) Create(ctx context.Context, req *pm.PaymentPlanCreateReq) error {
	// 如果已经有支付计划，则不再创建
	b, _ := s.orm.QueryNotDeleted().Where(purchasepayment.OrderID(req.OrderID)).Exist(ctx)
	if b {
		return nil
	}
	o, _ := NewOrder().QueryOrderById(ctx, req.OrderID)
	if o == nil {
		return errors.New("订单不存在")
	}
	if o.Status != purchaseorder.StatusPending {
		return errors.New("订单状态不正确")
	}
	if len(o.InstallmentPlan) == 0 {
		return errors.New("付款分期计划不存在")
	}
	if o.Edges.Rider == nil {
		return errors.New("订单用户不存在")
	}
	paymentBulk := make([]*ent.PurchasePaymentCreate, 0)
	dates := o.InstallmentPlan.BillingDates(time.Now())
	for k, plan := range o.InstallmentPlan {
		paymentBulk = append(paymentBulk, s.orm.Create().
			SetOutTradeNo(tools.NewUnique().NewSN()).
			SetIndex(k).
			SetStatus(purchasepayment.StatusObligation).
			SetTotal(plan.Amount).
			SetAmount(plan.Amount).
			SetBillingDate(dates[k]),
		)
	}
	_, err := s.orm.CreateBulk(paymentBulk...).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Pay 支付
func (s *paymentService) Pay(ctx context.Context, req *pm.PaymentReq) (*model.PurchasePayRes, error) {
	order, _ := NewOrder().QueryOrderById(ctx, req.OrderID)
	if order == nil {
		return nil, errors.New("订单不存在")
	}
	if !(order.Status == purchaseorder.StatusPending || order.Status == purchaseorder.StatusStaging) {
		return nil, errors.New("订单状态不正确")
	}
	if order.Edges.Goods == nil {
		return nil, errors.New("订单商品不存在")
	}
	stage, _ := s.QueryByIndex(ctx, req.OrderID, *req.PlanIndex)
	if stage == nil {
		return nil, errors.New("支付计划不存在")
	}
	payreq := &model.PurchasePayReq{
		OutTradeNo: stage.OutTradeNo,
		Subject:    order.Edges.Goods.Name,
		Amount:     tools.NewDecimal().Sum(stage.Amount, stage.Forfeit),
	}

	var err error
	var prepay string
	switch req.Payway {
	case pm.Alipay:
		prepay, err = alipay.NewApp().AppPayPurchase(payreq)
		if err != nil {
			zap.L().Error("支付宝支付失败", zap.Error(err))
			return nil, err
		}
	case pm.Wechat:
		prepay, err = wechat.NewApp().AppPayPurchase(payreq)
		if err != nil {
			zap.L().Error("微信支付失败", zap.Error(err))
			return nil, err
		}
	default:
		return nil, errors.New("支付方式不正确")
	}
	return &model.PurchasePayRes{
		Prepay:     prepay,
		OutTradeNo: stage.OutTradeNo,
	}, err
}
