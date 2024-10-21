package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	pm "github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/goods"
	"github.com/auroraride/aurservd/internal/ent/purchaseorder"
	"github.com/auroraride/aurservd/internal/ent/purchasepayment"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/auroraride/aurservd/pkg/utils"
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
func (s *orderService) Create(ctx context.Context, r *ent.Rider, req *pm.PurchaseOrderCreateReq) error {
	// 判定是否实名认证
	if !s.IsAuthed(r) {
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

// List 订单列表
func (s *orderService) List(req *pm.PurchaseOrderListReq) (res *model.PaginationRes) {
	q := s.orm.QueryNotDeleted().
		WithPayments(
			func(query *ent.PurchasePaymentQuery) {
				query.Order(ent.Asc(purchasepayment.FieldIndex))
			},
		).
		WithRider().
		WithStore().
		WithGoods().
		Order(ent.Desc(purchaseorder.FieldCreatedAt))
	s.listFilter(q, req.PurchaseOrderListFilter)
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.PurchaseOrder) (result pm.PurchaseOrderListRes) {
		return s.detail(item)
	})
	return
}

// listFilter 订单列表筛选
func (s *orderService) listFilter(q *ent.PurchaseOrderQuery, req pm.PurchaseOrderListFilter) {
	if req.Keyword != nil {
		q.Where(
			purchaseorder.Or(
				purchaseorder.ActiveNameContains(*req.Keyword),
				purchaseorder.ActivePhoneContains(*req.Keyword),
				purchaseorder.HasRiderWith(rider.NameContains(*req.Keyword)),
				purchaseorder.HasRiderWith(rider.PhoneContains(*req.Keyword)),
			),
		)
	}
	if req.ID != nil {
		q.Where(purchaseorder.ID(*req.ID))
	}
	if req.Sn != nil {
		q.Where(purchaseorder.SnContains(*req.Sn))
	}
	if req.Status != nil {
		q.Where(purchaseorder.StatusEQ(purchaseorder.Status(req.Status.Value())))
	}
	if req.BillStatus != nil {
		switch *req.BillStatus {
		case pm.BillStatusNormal:
			// 正常：订单未激活，订单已激活未超过付款日
			q.Where(
				purchaseorder.Or(
					purchaseorder.StartDateIsNil(),
					purchaseorder.And(
						purchaseorder.StartDateNotNil(),
						purchaseorder.HasPaymentsWith(
							purchasepayment.StatusEQ(purchasepayment.Status(pm.PaymentStatusObligation)),
							purchasepayment.BillingDateGTE(time.Now().AddDate(0, 0, -1)),
						),
					),
				),
			)

		case pm.BillStatusOverdue:
			// 逾期：订单已激活且已超过付款日
			q.Where(
				purchaseorder.StartDateNotNil(),
				purchaseorder.HasPaymentsWith(
					purchasepayment.StatusEQ(purchasepayment.Status(pm.PaymentStatusObligation)),
					purchasepayment.BillingDateLT(time.Now().AddDate(0, 0, -1)),
				),
			)
		}

	}
	if req.StoreID != nil {
		q.Where(purchaseorder.StoreID(*req.StoreID))
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			purchaseorder.CreatedAtGTE(start),
			purchaseorder.CreatedAtLTE(end),
		)
	}
	if req.RiderID != nil {
		q.Where(purchaseorder.RiderID(*req.RiderID))
	}
}

// detail 详情数据
func (s *orderService) detail(item *ent.PurchaseOrder) (res pm.PurchaseOrderListRes) {
	res = pm.PurchaseOrderListRes{
		ID:               item.ID,
		Status:           pm.OrderStatus(item.Status),
		InstallmentTotal: item.InstallmentTotal,
		InstallmentStage: item.InstallmentStage,
		Sn:               item.Sn,
		Color:            item.Color,
		CreatedAt:        item.CreatedAt.Format(carbon.DateTimeLayout),
		Remark:           item.Remark,
		BillStatus:       pm.BillStatusNormal,
		Signed:           item.Signed,
	}

	// 订单信息关键字段完善
	if item.StartDate != nil {
		res.StartDate = silk.String(item.StartDate.Format(carbon.DateTimeLayout))
	}
	if item.ContractURL != "" {
		res.ContractUrl = silk.String(item.ContractURL)
	}
	if item.DocID != "" {
		encryptDocID, _ := utils.EncryptAES([]byte(ar.Config.Contract.EncryptKey), item.DocID)
		if encryptDocID != "" {
			res.DocID = silk.String(encryptDocID)
		}
	}
	if item.ActiveName != "" && item.ActivePhone != "" {
		res.ActiveName = silk.String(item.ActiveName)
		res.ActivePhone = silk.String(item.ActivePhone)
	}

	// 商品信息
	if item.Edges.Goods != nil {
		g := item.Edges.Goods
		// 解析付款方案数据
		payPlans := make([][]float64, 0)
		for _, p := range g.PaymentPlans {
			payPlan := make([]float64, 0)
			for _, o := range p {
				payPlan = append(payPlan, o.Amount)
			}
			payPlans = append(payPlans, payPlan)
		}
		res.Goods = &definition.Goods{
			ID:           g.ID,
			Sn:           g.Sn,
			Name:         g.Name,
			Type:         definition.GoodsType(g.Type),
			Lables:       g.Lables,
			Price:        g.Price,
			Weight:       g.Weight,
			HeadPic:      g.HeadPic,
			Photos:       g.Photos,
			Intro:        g.Intro,
			CreatedAt:    item.CreatedAt.Format(carbon.DateTimeLayout),
			Status:       definition.GoodsStatus(g.Status),
			Remark:       item.Remark,
			PaymentPlans: payPlans,
		}
		// 未激活时订单金额默认为商品价格
		res.Amount = res.Goods.Price
		// 当前订单默认分期方案索引
		res.PlanIndex = silk.Int(g.PaymentPlans.PlanIndex(item.InstallmentPlan))
	}

	// 门店信息
	if item.Edges.Store != nil {
		res.StoreID = item.Edges.Store.ID
		res.StoreName = item.Edges.Store.Name
	}

	// 当订单状态为分期中（已激活）时
	if item.Status == purchaseorder.StatusStaging {
		// 支付金额信息
		for _, p := range item.Edges.Payments {
			res.Amount += p.Amount
			if p.Status.String() == pm.PaymentStatusPaid.Value() {
				res.PaidAmount += p.Amount
			}
			// 订单已激活且分期账单有未付款逾期数据
			if p.Status == purchasepayment.StatusObligation &&
				p.BillingDate.Before(carbon.Now().StartOfDay().StdTime()) {
				res.BillStatus = pm.BillStatusOverdue
			}
		}
	}

	// 骑手信息
	if item.Edges.Rider != nil {
		res.RiderName = item.Edges.Rider.Name
		res.RiderPhone = item.Edges.Rider.Phone
	}

	// 违约说明（暂时先固定文本返回前段）
	res.Formula = pm.PurchaseOrderFormula

	// 解析分期方案数据
	insPlan := make([]float64, 0)
	for _, o := range item.InstallmentPlan {
		insPlan = append(insPlan, o.Amount)
	}
	res.InstallmentPlan = insPlan

	return
}

// Detail 订单详情
func (s *orderService) Detail(id uint64) (res pm.PurchaseOrderDetail) {
	item, _ := s.orm.QueryNotDeleted().
		Where(
			purchaseorder.ID(id),
		).
		WithPayments(
			func(query *ent.PurchasePaymentQuery) {
				query.Order(ent.Asc(purchasepayment.FieldIndex))
			},
		).
		WithRider().
		WithStore().
		WithGoods().
		WithFollows().
		First(context.Background())

	// 订单详情数据
	res.PurchaseOrderListRes = s.detail(item)

	// 分期订单数据
	payments := make([]*pm.PaymentDetail, 0)
	for _, p := range item.Edges.Payments {
		status := pm.PaymentStatus(p.Status)

		// 未支付逾期判定，当期逾期未支付
		if p.Status == purchasepayment.StatusObligation && p.BillingDate.Before(carbon.Now().StartOfDay().StdTime()) {
			status = pm.PaymentStatusOverdue
		}

		payment := &pm.PaymentDetail{
			ID:         p.ID,
			Total:      p.Total,
			Amount:     p.Amount,
			Forfeit:    p.Forfeit,
			Payway:     pm.Payway(p.Payway),
			OutTradeNo: p.OutTradeNo,
			Status:     status,
		}

		// 账单日期
		if !p.BillingDate.IsZero() {
			payment.BillingDate = p.BillingDate.Format(carbon.DateLayout)
			// 逾期天数
			if time.Now().After(p.BillingDate) {
				payment.OverdueDays = int(time.Since(p.BillingDate).Hours() / 24)
			}
		}

		// 支付时间
		if p.PaymentDate != nil {
			payment.PaymentDate = p.PaymentDate.Format(carbon.DateTimeLayout)
		}
		payments = append(payments, payment)
	}
	res.Payments = payments

	// 跟进数据
	follows := make([]*pm.PurchaseOrderFollow, 0)
	for _, f := range item.Edges.Follows {
		follows = append(follows, &pm.PurchaseOrderFollow{
			ID:        f.ID,
			Content:   f.Content,
			Pics:      f.Pics,
			Modifier:  f.LastModifier,
			CreatedAt: f.CreatedAt.Format(carbon.DateTimeLayout),
		})
	}
	res.Follows = follows

	return
}

// Follow 跟进订单
func (s *orderService) Follow(ctx context.Context, req *pm.PurchaseOrderFollowReq, md *model.Modifier) (err error) {
	order, _ := NewOrder().QueryOrderById(ctx, req.ID)
	if order == nil {
		return errors.New("订单不存在")
	}

	return ent.Database.PurchaseFollow.Create().
		SetOrderID(order.ID).
		SetContent(req.Content).
		SetPics(req.Pics).
		SetCreator(md).
		SetLastModifier(md).
		Exec(ctx)
}

// Active 激活订单
func (s *orderService) Active(ctx context.Context, req *pm.PurchaseOrderActiveReq, md *model.Modifier) (err error) {
	order, _ := NewOrder().QueryOrderById(ctx, req.ID)
	if order == nil {
		return errors.New("订单不存在")
	}
	if order.Status.String() == pm.OrderStatusCancelled.Value() {
		return errors.New("订单已取消")
	}
	// 判断订单是否已激活
	if order.StartDate != nil {
		return errors.New("订单已激活")
	}
	// 订单更新
	ou := s.orm.Update().
		Where(purchaseorder.ID(req.ID)).
		SetStoreID(req.StoreID).
		SetSn(req.Sn).
		SetColor(req.Color).
		SetRemark(req.Remark).
		SetActiveName(md.Name).
		SetActivePhone(md.Phone).
		SetLastModifier(md)

	// 查询订单默认商品
	g := order.Edges.Goods
	if g == nil {
		return errors.New("订单商品不存在")
	}

	// 激活的商品是否为默认下单商品
	var goodsCheckNew bool
	var goodPlan model.GoodsPaymentPlan

	if g.ID != req.GoodsID {
		var newGoods *ent.Goods
		newGoods, _ = ent.Database.Goods.QueryNotDeleted().
			Where(goods.ID(req.GoodsID)).
			First(ctx)
		if newGoods == nil {
			return errors.New("激活商品不存在")
		}
		goodPlan = newGoods.PaymentPlans.Plan(*req.PlanIndex)
		goodsCheckNew = true

	} else {
		goodPlan = g.PaymentPlans.Plan(*req.PlanIndex)
	}
	if len(goodPlan) == 0 {
		return errors.New("分期计划不存在")
	}

	// 对比分期计划是否相等，不相等更新分期计划数据（更换了商品直接更新新的分期方案）
	if goodsCheckNew || !order.InstallmentPlan.Equal(goodPlan) {
		ou.SetInstallmentPlan(goodPlan)
		ou.SetInstallmentTotal(len(goodPlan))
	}
	err = ou.Exec(ctx)
	if err != nil {
		return err
	}

	// 创建分期计划订单数据
	return NewPayment().Create(ctx, &pm.PaymentPlanCreateReq{OrderID: req.ID}, md)
}

// Cancel 订单取消
func (s *orderService) Cancel(ctx context.Context, id uint64, md *model.Modifier) (err error) {
	order, _ := NewOrder().QueryOrderById(ctx, id)
	if order == nil {
		return errors.New("订单不存在")
	}
	if order.Status.String() == pm.OrderStatusCancelled.Value() {
		return errors.New("订单已取消")
	}
	// 更新订单
	err = s.orm.Update().Where(purchaseorder.ID(order.ID)).
		SetStatus(purchaseorder.Status(pm.OrderStatusCancelled)).
		SetLastModifier(md).
		Exec(ctx)
	if err != nil {
		return
	}
	// 更新分期订单
	ent.Database.PurchasePayment.Update().
		Where(purchasepayment.OrderID(order.ID)).
		SetStatus(purchasepayment.Status(pm.PaymentStatusCanceled)).SetLastModifier(md)

	return

}

// Export 订单导出
func (s *orderService) Export(req *pm.PurchaseOrderExportReq, md *model.Modifier) model.ExportRes {
	q := s.orm.QueryNotDeleted().
		WithPayments(
			func(query *ent.PurchasePaymentQuery) {
				query.Order(ent.Asc(purchasepayment.FieldIndex))
			},
		).
		WithRider().
		WithStore().
		WithGoods()
	s.listFilter(q, req.PurchaseOrderListFilter)
	return service.NewExportWithModifier(md).Start("购车订单", req.PurchaseOrderListFilter, nil, req.Remark, func(path string) {
		items, _ := q.All(context.Background())
		var rows tools.ExcelItems
		title := []any{
			"订单编号",
			"订单状态",
			"商品名称",
			"商品编号",
			"订单金额",
			"已支付",
			"分期期数",
			"还款状态",
			"骑手",
			"骑手电话",
			"提车门店",
			"车架号",
			"激活人",
			"订单时间",
			"激活时间",
			"备注",
		}
		rows = append(rows, title)
		for _, item := range items {
			detail := s.detail(item)
			var activeInfo string
			if detail.ActiveName != nil && detail.ActivePhone != nil {
				activeInfo = *detail.ActiveName + "-" + *detail.ActivePhone
			}
			row := []any{
				detail.ID,
				detail.Status.String(),
				detail.Goods.Name,
				detail.Goods.Sn,
				detail.Amount,
				detail.PaidAmount,
				detail.InstallmentTotal,
				detail.BillStatus.String(),
				detail.RiderName,
				detail.RiderPhone,
				detail.StoreName,
				detail.Sn,
				activeInfo,
				detail.CreatedAt,
				detail.StartDate,
				detail.Remark,
			}
			rows = append(rows, row)
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}
