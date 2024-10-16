package service

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/model"
	pm "github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/goods"
	"github.com/auroraride/aurservd/internal/ent/purchaseorder"
	"github.com/auroraride/aurservd/pkg/tools"
)

type orderService struct {
	orm *ent.PurchaseOrderClient
}

func NewOrder() *orderService {
	return &orderService{
		orm: ent.Database.PurchaseOrder,
	}
}

// 查询商品
func (s *orderService) queryGoods(ctx context.Context, goodsID uint64) (*ent.Goods, error) {
	g, _ := ent.Database.Goods.QueryNotDeleted().Where(goods.ID(goodsID)).First(ctx)
	if g == nil {
		return nil, errors.New("商品不存在")
	}
	return g, nil
}

// IsAuthed 是否已认证
func (s *orderService) IsAuthed(u *ent.Rider) bool {
	return u.Edges.Person != nil && !model.PersonAuthStatus(u.Edges.Person.Status).RequireAuth()
}

// Create 创建订单
func (s *orderService) Create(ctx context.Context, r *ent.Rider, req *pm.OrderCreateReq) error {
	// 判定是否实名认证
	if s.IsAuthed(r) {
		return errors.New("未实名认证, 请先实名认证")
	}
	// 查询商品
	g, _ := s.queryGoods(ctx, req.GoodsID)
	if g == nil {
		return errors.New("商品不存在")
	}

	if req.PlanIndex == nil {
		return errors.New("分期方案有误")
	}

	if len(g.PaymentPlans) == 0 {
		return errors.New("分期计划不存在")
	}

	// 创建订单
	err := s.orm.Create().
		SetGoodsID(g.ID).
		SetRiderID(r.ID).
		SetInstallmentPlan(g.PaymentPlans.Plan(*req.PlanIndex)).
		SetSn(tools.NewUnique().NewSN()).
		SetStatus(purchaseorder.StatusPending).
		SetInstallmentTotal(len(g.PaymentPlans.Plan(*req.PlanIndex))).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// QueryOrderById 通过id查询订单
func (s *orderService) QueryOrderById(ctx context.Context, id uint64) (*ent.PurchaseOrder, error) {
	o, _ := s.orm.QueryNotDeleted().Where(purchaseorder.ID(id)).WithGoods().WithRider().First(ctx)
	if o == nil {
		return nil, errors.New("订单不存在")
	}
	return o, nil
}
