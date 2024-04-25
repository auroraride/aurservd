package biz

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/orderrefund"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionreferrals"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/payment"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type orderBiz struct {
	orm *ent.OrderClient
	ctx context.Context
}

func NewOrderBiz() *orderBiz {
	return &orderBiz{
		orm: ent.Database.Order,
		ctx: context.Background(),
	}
}

// DepositCredit 押金免支付订单
func (s *orderBiz) DepositCredit(r *ent.Rider, req *definition.OrderDepositCreditReq) (result *definition.OrderDepositCreditRes, err error) {
	result = &definition.OrderDepositCreditRes{}

	no := tools.NewUnique().NewSN28()
	// 查询套餐是否存在
	p := service.NewPlan().QueryEffectiveWithID(req.PlanID)

	// 判定套餐开启押金
	if !p.Deposit {
		return nil, errors.New("套餐不支持免押金")
	}

	if p.Deposit {
		if req.Payway == model.OrderPaywayAlipayAuthFreeze && !p.DepositAlipayAuthFreeze {
			return nil, errors.New("套餐不支持支付宝免押金")
		}
		if req.Payway == model.OrderPaywayWechatDeposit && !p.DepositWechatPayscore {
			return nil, errors.New("套餐不支持微信免押金")
		}
		if (req.Payway == model.OrderPaywayWechat || req.Payway == model.OrderPaywayAlipay) && !p.DepositPay {
			return nil, errors.New("套餐不支持支付宝或微信免押金")
		}
	}

	depositFree := &model.DepositCredit{
		RiderID: r.ID,
		Amount:  p.DepositAmount,
		Plan:    p.BasicInfo(),
		PlanID:  req.PlanID,
		Payway:  req.Payway,
	}
	if req.Payway == model.OrderPaywayAlipay || req.Payway == model.OrderPaywayWechat {
		depositFree.OutTradeNo = no
		result.OutTradeNo = no
	} else {
		depositFree.OutOrderNo = no
		result.OutOrderNo = no
	}

	// 订单字段
	prepay := &model.PaymentCache{
		CacheType:     model.PaymentCacheTypeDeposit,
		DepositCredit: depositFree,
	}

	// 订单缓存
	err = cache.Set(s.ctx, no, prepay, 20*time.Minute).Err()
	if err != nil {
		zap.L().Error("获取免押订单信息失败", zap.Error(err))
		return nil, err
	}

	// 发起请求
	var str string
	switch req.Payway {
	// case model.OrderPaywayAlipay:
	// 使用支付宝支付
	// str, err = payment.NewAlipay().AppPay(prepay)
	// if err != nil {
	// 	return nil, err
	// }
	// case model.OrderPaywayWechat:
	// 	使用微信支付
	// str, err = payment.NewWechat().AppPay(prepay)
	// if err != nil {
	// 	return nil, err
	// }
	case model.OrderPaywayAlipayAuthFreeze:
		// 使用支付宝预授权支付
		str, err = payment.NewAlipay().FandAuthFreeze(prepay)
		if err != nil {
			zap.L().Error("支付宝预授权支付请求失败", zap.Error(err))
			return nil, err
		}
	case model.OrderPaywayWechatDeposit:
		// TODO
	default:
		snag.Panic("未知的支付请求")
	}
	result.Prepay = str
	return
}

// DoPayment 处理预支付或者预支付转支付
func (s *orderBiz) DoPayment(pc *model.PaymentCache) {
	if pc == nil {
		return
	}
	switch pc.CacheType {
	case model.PaymentCacheTypeAlipayAuthFreeze:
		// 预支付
		s.OrderPaid(pc.Subscribe)
	case model.PaymentCacheTypeDeposit:
		service.NewOrder().DepositPay(pc.DepositCredit)
	default:
		return
	}
}

// CancelDeposit 取消押金订单 订单为以下状态时可以取消订单：INIT（初始化）、AUTHORIZED（已创建）（此时一般为用户取消服务时使用）
func (s *orderBiz) CancelDeposit(req *definition.OrderDepositCancelReq) error {
	// 查询订单是否存在
	o, _ := s.orm.Query().Where(order.OutOrderNo(req.OutOrderNo)).First(s.ctx)
	if o == nil {
		return errors.New("订单不存在")
	}
	// 查询在支付宝的订单状态
	var detailQuery *alipay.FundAuthOperationDetailQueryRsp
	detailQuery, err := payment.NewAlipay().AlipayFundAuthOperationDetailQuery(
		definition.FundAuthOperationDetailReq{
			OutOrderNo:   o.OutOrderNo,
			OutRequestNo: o.OutRequestNo,
		},
	)
	if err != nil {
		zap.L().Error("查询支付宝订单状态失败", zap.Error(err))
		return err
	}
	// 订单状态为初始化或已创建时可以取消订单
	if detailQuery.OrderStatus == alipay.OrderStatusInit || detailQuery.OrderStatus == alipay.OrderStatusAuthorized {
		// 取消订单
		_, err = payment.NewAlipay().AlipayFundAuthOperationCancel(o.AuthNo, o.OutRequestNo)
		if err != nil {
			zap.L().Error("取消订单失败", zap.Error(err))
			return err
		}
		// 更新订单状态
		_, err = o.Update().SetStatus(model.OrderStatusCanceled).Save(s.ctx)
		if err != nil {
			zap.L().Error("更新订单状态失败", zap.Error(err), log.Payload(req))
			return err
		}
	}
	return nil
}

// DoPaymentFreezeToPay 处理冻结资金转支付
func (s *orderBiz) DoPaymentFreezeToPay(req *definition.OrderDepositFreezeToPayRes) error {
	// 查询订单是否存在
	o, _ := s.orm.Query().Where(order.OutTradeNo(req.OutTradeNo)).First(s.ctx)
	if o == nil {
		return errors.New("订单不存在")
	}
	// 更新订单状态
	_, err := o.Update().SetTradeNo(req.TradeNo).
		SetOutTradeNo(req.OutTradeNo).
		SetTradePayAt(time.Now()).Save(s.ctx)
	if err != nil {
		zap.L().Error("订单支付失败", zap.Error(err))
		return err
	}

	// 查询订单是否有押金订单 有押金订单转为拒绝 并备注(拒接退款)
	orderRefund, _ := ent.Database.OrderRefund.QueryNotDeleted().Where(orderrefund.OrderID(o.ID)).First(s.ctx)
	if orderRefund != nil {
		_, err = orderRefund.Update().SetStatus(model.RefundStatusRefused).SetRemark("拒绝退款").Save(s.ctx)
		if err != nil {
			zap.L().Error("订单退款失败", zap.Error(err))
			return err
		}
	}
	return nil
}

// Create 订单创建
func (s *orderBiz) Create(r *ent.Rider, req *definition.OrderCreateReq) (result *model.OrderCreateRes, err error) {
	if req.OrderType == model.OrderTypeFee {
		// 逾期费用不能使用预授权
		if req.Payway == model.OrderPaywayAlipayAuthFreeze {
			return nil, errors.New("逾期费用不能使用支付宝预授权")
		}
		return service.NewOrder().CreateFee(r.ID, req.Payway), nil
	}

	result = new(model.OrderCreateRes)

	if req.PlanID == 0 {
		return nil, errors.New("套餐ID不能为空")
	}

	// 查询套餐是否存在
	p := service.NewPlan().QueryEffectiveWithID(req.PlanID)

	if p.BrandID != nil && p.Edges.Brand == nil {
		return nil, errors.New("骑士卡错误")
	}

	// 查询是否企业骑手
	if r.EnterpriseID != nil {
		return nil, errors.New("企业骑手无法购买骑士卡")
	}

	// 查询是否有退款中的订单
	if exists, _ := ent.Database.Order.QueryNotDeleted().Where(
		order.RiderID(r.ID),
		order.Status(model.OrderStatusRefundPending),
	).Exist(s.ctx); exists {
		return nil, errors.New("存在退款中的订单")
	}

	var past *int
	sub := service.NewSubscribe().Recent(r.ID, *r.PersonID)
	// 判定类型条件
	var subID, orderID *uint64
	t := req.OrderType
	switch t {
	case model.OrderTypeNewly, model.OrderTypeAgain:
		// 新签/重签判定
		if req.CityID == 0 {
			return nil, errors.New("城市ID不能为空")
		}
		t, past = service.NewOrder().PreconditionNewly(sub)
	case model.OrderTypeRenewal:
		// 续签判定
		service.NewOrder().PreconditionRenewal(sub)
		if sub.Remaining < 0 && int(p.Days)+sub.Remaining < 0 {
			return nil, errors.New("无法继续, 逾期天数大于骑士卡天数")
		}
		subID = silk.UInt64(sub.ID)
		orderID = silk.UInt64(sub.InitialOrderID)
		req.CityID = sub.CityID
	default:
		return nil, errors.New("订单类型错误")
	}

	// 判定押金是否需要支付
	var deposit float64

	// 只有新签和重签才需要支付押金
	if t == model.OrderTypeNewly || t == model.OrderTypeAgain {
		// 车电套餐需要传入门店ID
		if p.Type == model.PlanTypeEbikeWithBattery.Value() && req.StoreID == nil {
			return nil, errors.New("车电套餐需要选择门店")
		}

		// DepositOrderNo 表示如果是微信免押 则需要传入押金订单号 不需要生成新的押金订单
		if req.DepositType != nil && req.Payway != model.OrderPaywayWechatDeposit && p.Deposit && req.DepositOrderNo == nil {
			// 只有套餐支持押金支付或者 芝麻信用免押时才需要支付押金
			if (*req.DepositType == model.DepositTypeAlipayAuthFreeze || *req.DepositType == model.DepositTypePay) && (p.DepositAlipayAuthFreeze || p.DepositPay) {
				deposit = p.DepositAmount
			}
		}
	}

	no := tools.NewUnique().NewSN28()

	price := p.Price

	// 计算新签优惠
	var ramount float64
	if t == model.OrderTypeNewly && p.DiscountNewly > 0 {
		ramount = p.DiscountNewly
		price = tools.NewDecimal().Sub(price, ramount)
	}

	// 获取优惠券
	var camount float64
	now := time.Now()
	if len(req.Coupons) > 0 {
		coupons := service.NewCoupon().QueryIDs(req.Coupons)
		if len(req.Coupons) != len(coupons) {
			return nil, errors.New("优惠券不存在")
		}
		var isExclusive bool
		cm := make(map[uint64]uint64)
		for _, c := range coupons {
			// 校验有效期
			if c.ExpiresAt.Before(now) {
				return nil, errors.New("优惠券已过期")
			}

			// 是否互斥
			if c.Rule == model.CouponRuleExclusive.Value() {
				isExclusive = true
			}

			// 是否叠加
			if _, ok := cm[c.TemplateID]; ok {
				return nil, errors.New("优惠券不可叠加")
			}

			// 是否限制城市
			if len(c.Cities) > 0 {
				cityUseable := false
				for _, cc := range c.Cities {
					if cc.ID == req.CityID {
						cityUseable = true
					}
				}
				if !cityUseable {
					return nil, errors.New("当前城市无法使用优惠券")
				}
			}

			// 是否限制骑士卡
			if len(c.Plans) > 0 {
				planUsable := false
				for _, cp := range c.Plans {
					if cp.ID == req.PlanID {
						planUsable = true
					}
				}
				if !planUsable {
					return nil, errors.New("当前骑士卡无法使用优惠券")
				}
			}

			cm[c.TemplateID] = c.ID

			// 累加优惠券金额
			camount += c.Amount
		}
		if isExclusive && len(req.Coupons) > 1 {
			return nil, errors.New("所选优惠券互斥")
		}
		price = tools.NewDecimal().Sub(price, camount)
		// 暂时处理成支付一分钱
		if price <= 0 {
			price = 0.01
		}
	}

	var points int64
	// 积分抵扣
	price, points, err = s.DeductPoints(r, req.Point, req.PointNum, price)
	if err != nil {
		return nil, err
	}

	if price < 0 {
		return nil, errors.New("订单金额错误")
	}

	// Development模式支付一分钱
	if ar.Config.Environment.IsDevelopment() {
		price = 0.01
		if deposit > 0 {
			deposit = 0.01
		}
	}

	// 总计支付金额
	total := tools.NewDecimal().Sum(price, deposit)

	paySubscribe := &model.PaymentSubscribe{
		CityID:         req.CityID,
		OrderType:      t,
		RiderID:        r.ID,
		Name:           "购买" + p.Name,
		Amount:         total,
		Payway:         req.Payway,
		Deposit:        deposit,
		PastDays:       past,
		Commission:     p.Commission,
		Model:          p.Model,
		Days:           p.Days,
		OrderID:        orderID,
		SubscribeID:    subID,
		Points:         points,
		PointRatio:     model.PointRatio,
		CouponAmount:   camount,
		Coupons:        req.Coupons,
		DiscountNewly:  ramount,
		EbikeBrandID:   p.BrandID,
		Plan:           p.BasicInfo(),
		StoreID:        req.StoreID,
		DepositOrderNo: req.DepositOrderNo,
		AgreementHash:  req.AgreementHash,
		DepositType:    req.DepositType,
	}

	prepay := &model.PaymentCache{
		Subscribe: paySubscribe,
	}

	// 订单缓存
	switch req.Payway {
	case model.OrderPaywayAlipayAuthFreeze:
		prepay.CacheType = model.PaymentCacheTypeAlipayAuthFreeze
		paySubscribe.OutOrderNo = no
	default:
		prepay.CacheType = model.PaymentCacheTypePlan
		paySubscribe.OutTradeNo = no
	}

	service.NewOrder().Prepay(req.Payway, no, prepay, result)

	return
}

// OrderPaid 预授权支付成功
func (s *orderBiz) OrderPaid(trade *model.PaymentSubscribe) {

	no := trade.OutOrderNo
	q := ent.Database.Order.Query().Where(order.OutOrderNo(no))

	// 查询订单是否已存在
	if exists, err := q.Exist(s.ctx); err == nil && exists {
		return
	}

	zap.L().Info("订单支付回调: "+no, log.JsonData(trade))

	var sub *ent.Subscribe
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		var o *ent.Order
		// 创建订单
		oc := tx.Order.Create().
			SetPayway(trade.Payway).
			SetPlanID(trade.Plan.ID).
			SetRiderID(trade.RiderID).
			SetAmount(tools.NewDecimal().Sub(trade.Amount, trade.Deposit)).
			SetTotal(trade.Amount).
			SetStatus(model.OrderStatusPaid).
			SetType(trade.OrderType).
			SetCityID(trade.CityID).
			SetInitialDays(int(trade.Days)).
			SetNillableParentID(trade.OrderID).
			SetNillableSubscribeID(trade.SubscribeID).
			SetNillablePastDays(trade.PastDays).
			SetPoints(trade.Points).
			SetPointRatio(trade.PointRatio).
			SetCouponAmount(trade.CouponAmount).
			SetDiscountNewly(trade.DiscountNewly).
			SetNillableBrandID(trade.EbikeBrandID).
			SetOutOrderNo(trade.OutOrderNo).
			SetAuthNo(trade.AuthNo).
			SetOutRequestNo(trade.OutRequestNo)
		if len(trade.Coupons) > 0 {
			oc.AddCouponIDs(trade.Coupons...)
			// 更新优惠券使用状态
			for _, couponID := range trade.Coupons {
				err = tx.Coupon.UpdateOneID(couponID).SetPlanID(trade.Plan.ID).SetUsedAt(time.Now()).Exec(s.ctx)
				if err != nil {
					zap.L().Error("订单已支付, 但优惠券更新失败: "+no+", couponID="+strconv.FormatUint(couponID, 10), zap.Error(err))
				}
			}
		}
		o, err = oc.Save(s.ctx)
		if err != nil {
			zap.L().Error("订单已支付, 但订单创建失败: "+no, zap.Error(err))
			return
		}

		// 更新积分
		if trade.Points > 0 {
			// 赠送积分
			gift, proportion := service.NewPoint().CalculateGift(trade.Amount, trade.CityID)

			var r *ent.Rider
			r, _ = tx.Rider.QueryNotDeleted().Where(rider.ID(trade.RiderID)).First(s.ctx)
			if r == nil {
				zap.L().Error("订单已支付查询骑手失败: "+no, zap.Error(err))
				return
			}

			before := r.Points
			_, err = r.Update().AddPoints(-trade.Points + gift).Save(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但积分更新失败: "+no, zap.Error(err))
				return
			}

			err = tx.PointLog.Create().
				SetPoints(-trade.Points).
				SetRiderID(trade.RiderID).
				SetReason("订阅骑士卡").
				SetType(model.PointLogTypeConsume.Value()).
				SetAfter(before - trade.Points).
				SetAttach(&model.PointLogAttach{Plan: trade.Plan}).
				SetOrder(o).
				Exec(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但积分消费记录创建失败: "+no, zap.Error(err))
			}

			service.NewPoint().RemovePreConsume(r, trade.Points)

			// 存储赠送积分
			if gift > 0 {
				err = tx.PointLog.Create().
					SetPoints(gift).
					SetRiderID(trade.RiderID).
					SetReason("消费赠送").
					SetType(model.PointLogTypeAward.Value()).
					SetAfter(r.Points).
					SetAttach(&model.PointLogAttach{PointGift: &model.PointGift{
						Amount:     trade.Amount,
						Proportion: proportion,
					}}).
					SetOrder(o).
					Exec(s.ctx)
				if err != nil {
					zap.L().Error("订单已支付, 但积分赠送记录创建失败: "+no, zap.Error(err))
				}
			}
		}

		// 如果有押金, 创建押金订单
		var do *ent.Order
		if trade.Deposit > 0 {
			depositOrder := tx.Order.Create().
				SetStatus(model.OrderStatusPaid).
				SetPayway(trade.Payway).
				SetType(model.OrderTypeDeposit).
				SetAmount(trade.Deposit).
				SetTotal(trade.Amount).
				SetCityID(trade.CityID).
				SetRiderID(trade.RiderID).
				SetParentID(o.ID).
				SetOutOrderNo(trade.OutOrderNo).
				SetAuthNo(trade.AuthNo).
				SetOutRequestNo(trade.OutRequestNo)
			do, err = depositOrder.Save(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但押金订单创建失败: "+no, zap.Error(err))
				return
			}
		}

		// 如果有押金订单编号, 更新订单关联
		if trade.DepositOrderNo != nil {
			do, err = tx.Order.QueryNotDeleted().Where(order.Status(model.OrderStatusPaid), order.Type(model.OrderTypeDeposit)).
				Where(
					order.Or(
						order.OutOrderNo(*trade.DepositOrderNo),
						order.OutTradeNo(*trade.DepositOrderNo),
					)).First(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但押金订单查询失败: "+trade.OutTradeNo, zap.Error(err))
				return
			}
			_, err = do.Update().SetParentID(o.ID).Save(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但押金订单更新失败: "+trade.OutTradeNo, zap.Error(err))
				return
			}
		}

		// 实际剩余天数应减去当天
		remaining := int(trade.Days) - 1

		// 创建或更新subscribe
		// 新签或重签
		if trade.OrderType == model.OrderTypeNewly || trade.OrderType == model.OrderTypeAgain {
			// 创建subscribe
			sq := tx.Subscribe.Create().
				SetType(trade.OrderType).
				SetRiderID(trade.RiderID).
				SetModel(trade.Model).
				SetRemaining(remaining).
				SetInitialDays(int(trade.Days)).
				SetStatus(model.SubscribeStatusInactive).
				SetPlanID(trade.Plan.ID).
				SetCityID(trade.CityID).
				SetInitialOrderID(o.ID).
				AddOrders(o).
				SetNillableBrandID(trade.EbikeBrandID).
				SetIntelligent(trade.Plan.Intelligent).
				SetNillableStoreID(trade.StoreID).
				SetNillableAgreementHash(trade.AgreementHash)
			// 根据用户选择是否需要签约 默认不需要签约
			if trade.DepositType != nil {
				sq.SetDepositType(trade.DepositType.Value())
				sq.SetNeedContract(false)
				if *trade.DepositType == model.DepositTypeContract {
					sq.SetNeedContract(true)
				}
			}

			if do != nil {
				sq.AddOrders(do)
			}
			sub, err = sq.Save(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但新签或重签订阅创建失败: "+no, zap.Error(err))
				return
			}

			// 更新推广中的订阅
			err = tx.PromotionMember.Update().Where(promotionmember.RiderID(trade.RiderID)).SetSubscribe(sub).Exec(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但推广订阅更新失败: "+no, zap.Error(err))
				return
			}
			err = tx.PromotionReferrals.Update().Where(promotionreferrals.ReferredMemberID(trade.RiderID)).SetSubscribe(sub).Exec(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但推广订阅更新失败: "+no, zap.Error(err))
				return
			}
		}

		// 续签
		if trade.OrderType == model.OrderTypeRenewal {
			status := model.SubscribeStatusUsing
			// 查询状态
			sub, _ = service.NewSubscribe().Query(*trade.SubscribeID)
			if sub.Status == model.SubscribeStatusPaused {
				status = model.SubscribeStatusPaused
			}

			sub, err = tx.Subscribe.UpdateOneID(*trade.SubscribeID).
				AddRenewalDays(int(trade.Days)).
				AddRemaining(remaining).
				SetStatus(status).
				Save(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但续签订阅更新失败: "+no, zap.Error(err))
				return
			}

			entAt := tools.NewTime().WillEnd(sub.StartAt.Local(), sub.InitialDays)
			// 更新订单中的到期订阅到期时间(新增这个时间是为了记录预支付转支付时间)
			_, err = o.Update().SetSubscribeEndAt(entAt).Save(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付,  但订单更新失败: "+no, zap.Error(err))
				return
			}
		}

		// 当新签和重签的时候有提成
		if trade.OrderType == model.OrderTypeNewly && trade.Commission > 0 && sub != nil {
			// 创建提成
			_, err = tx.Commission.Create().SetOrderID(o.ID).SetPlanID(trade.Plan.ID).SetAmount(trade.Commission).SetStatus(model.CommissionStatusPending).SetRiderID(sub.RiderID).SetSubscribe(sub).Save(s.ctx)
			if err != nil {
				zap.L().Error("订单已支付, 但提成创建失败: "+no, zap.Error(err))
				return
			}
		}

		// 计算推广返佣 续签反佣 新签和重签在激活骑手返佣
		if trade.OrderType == model.OrderTypeRenewal {

			err = service.NewPromotionCommissionService().CommissionCalculation(tx, &promotion.CommissionCalculation{
				RiderID:      trade.RiderID,
				Type:         promotion.CommissionTypeRenewal,
				OrderID:      o.ID,
				ActualAmount: o.Total,
				PlanID:       trade.Plan.ID,
			})
			if err != nil {
				zap.L().Error("订单已支付, 续费返佣失败: "+no, zap.Error(err))
				return
			}
		}

		return
	})

	// 删除缓存
	cache.Del(context.Background(), no)

	if sub != nil && trade.OrderType == model.OrderTypeRenewal {
		go func() {
			_ = service.NewSubscribe().UpdateStatus(sub, false)
		}()
	}
}

// DeductPoints 积分抵扣
func (s *orderBiz) DeductPoints(r *ent.Rider, point bool, pointNum *int64, price float64) (float64, int64, error) {

	if !point || price <= 0.01 {
		// 不使用积分或订单金额太低不支持积分抵扣
		return price, 0, nil
	}

	pointServ := service.NewPoint()
	realPoints := pointServ.Real(r)
	if realPoints <= 0 {
		// 没有可用积分
		return price, 0, nil
	}

	var points int64
	if pointNum != nil && *pointNum > 0 {
		// 指定了积分抵扣数量
		if *pointNum > realPoints {
			return 0, 0, errors.New("积分不足")
		}
		if *pointNum > int64(price/model.PointRatio) {
			return 0, 0, fmt.Errorf("积分抵扣金额超过订单金额，最多抵扣 %d 积分", int64(price/model.PointRatio))
		}
		points = *pointNum
	} else {
		// 自动计算积分抵扣数量
		cents := int64(price / model.PointRatio)
		if realPoints < cents {
			// 若积分小于所需积分, 则全部扣除
			points = realPoints
		} else {
			// 若剩余积分大于所需金额, 则扣除剩余金额积分数量
			points = cents
		}
	}

	price = tools.NewDecimal().Sub(price, float64(points)*model.PointRatio)
	_, err := pointServ.PreConsume(r, points)
	if err != nil {
		return 0, 0, err
	}

	// 确保价格不小于最低金额
	if price <= 0 {
		price = 0.01
	}

	return price, points, nil
}
