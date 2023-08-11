package service

import (
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotionearnings"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/rider"

	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type promotionEarningsService struct {
	*BaseService
}

func NewPromotionEarningsService(params ...any) *promotionEarningsService {
	return &promotionEarningsService{
		BaseService: newService(params...),
	}
}

// List 会员收益列表
func (s *promotionEarningsService) List(req *promotion.EarningsReq) *model.PaginationRes {
	q := ent.Database.PromotionEarnings.Query().WithRider().WithCommission().Order(ent.Desc(promotionearnings.FieldCreatedAt))

	if req.ID != nil {
		q.Where(promotionearnings.MemberID(*req.ID))
	}

	if req.Keyword != nil {
		q.Where(
			promotionearnings.HasRiderWith(
				rider.Or(
					rider.PhoneContainsFold(*req.Keyword),
					rider.NameContainsFold(*req.Keyword),
				),
			),
		)
	}
	if req.Status != nil {
		q.Where(promotionearnings.Status(*req.Status))
	}
	if req.CommissionRuleKey != nil {
		q.Where(promotionearnings.CommissionRuleKey(*req.CommissionRuleKey))
	}

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			promotionearnings.CreatedAtGTE(start),
			promotionearnings.CreatedAtLTE(end),
		)
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.PromotionEarnings) (res promotion.EarningsRes) {
			res = promotion.EarningsRes{
				EarningsDetail: promotion.EarningsDetail{
					ID:                 item.ID,
					MemberID:           item.MemberID,
					CommissionRuleName: promotion.CommissionRuleKey(item.CommissionRuleKey).String(),
					CommissionID:       item.CommissionID,
					Amount:             item.Amount,
					CreateTime:         item.CreatedAt.Format(carbon.DateTimeLayout),
					Remark:             item.Remark,
					Status:             item.Status,
				},
			}
			if item.Edges.Rider != nil {
				res.Phone = item.Edges.Rider.Phone
				res.Name = item.Edges.Rider.Name
			}
			return
		},
	)
}

// Create 创建会员收益
func (s *promotionEarningsService) Create(tx *ent.Tx, req *promotion.EarningsCreateReq) error {
	return tx.PromotionEarnings.Create().
		SetMemberID(req.MemberID).
		SetRiderID(req.RiderID).
		SetCommissionID(req.CommissionID).
		SetCommissionRuleKey(string(req.CommissionRuleKey)).
		SetAmount(req.Amount).
		SetOrderID(req.OrderID).
		Exec(s.ctx)
}

// Cancel 取消收益
func (s *promotionEarningsService) Cancel(req *promotion.EarningsCancelReq) {
	// 查询收益
	earning, _ := ent.Database.PromotionEarnings.Query().Where(promotionearnings.ID(req.ID)).First(s.ctx)
	if earning == nil {
		snag.Panic("收益不存在")
	}
	if earning.Status == promotion.EarningsStatusCanceled.Value() {
		snag.Panic("收益已取消")
	}
	// 查询会员
	m, _ := ent.Database.PromotionMember.Query().Where(promotionmember.ID(earning.MemberID)).First(s.ctx)
	if m == nil {
		snag.Panic("会员不存在")
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 取消后优先从用户可提现余额中扣除
		// 可提现余额不足部分，从冻结余额扣除
		// 冻结余额不足，将冻结余额扣为负数
		balance, frozen := s.cancelIncome(m.Balance, m.Frozen, earning.Amount)

		// 更新会员余额
		_, err = tx.PromotionMember.UpdateOneID(earning.MemberID).SetBalance(balance).SetFrozen(frozen).Save(s.ctx)
		if err != nil {
			snag.Panic("更新会员余额失败")
		}

		// 更新返佣总收益
		com, _ := ent.Database.PromotionCommission.Query().Where(promotioncommission.ID(earning.CommissionID)).First(s.ctx)
		if com == nil {
			snag.Panic("返佣不存在")
		}
		_, err = tx.PromotionCommission.UpdateOneID(earning.CommissionID).SetAmountSum(tools.NewDecimal().Sub(com.AmountSum, earning.Amount)).Save(s.ctx)
		if err != nil {
			snag.Panic("更新返佣总收益失败")
		}

		// 更新收益状态
		_, err = tx.PromotionEarnings.UpdateOneID(req.ID).SetStatus(promotion.EarningsStatusCanceled.Value()).SetNillableRemark(req.Remark).Save(s.ctx)
		if err != nil {
			snag.Panic("更新收益状态失败")
		}

		return
	})
}

// CancelIncome 取消收益方法
func (s *promotionEarningsService) cancelIncome(balance float64, frozen float64, amount float64) (float64, float64) {
	dl := tools.NewDecimal()

	if dl.Sub(balance, amount) >= 0 {
		// 优先从余额中扣除取消金额
		balance = dl.Sub(balance, amount)
	} else {
		// 余额不足，从冻结金额中扣除 可为负数
		frozen = dl.Sub(frozen, dl.Sub(amount, balance))
		balance = 0
	}

	return balance, frozen
}

// CountCommission 查询返佣次数
func (s *promotionEarningsService) CountCommission(memberID, riderID uint64) (count int, err error) {
	return ent.Database.PromotionEarnings.Query().Where(
		promotionearnings.MemberID(memberID),
		promotionearnings.RiderID(riderID),
	).Where(promotionearnings.Or(
		promotionearnings.CommissionRuleKeyEQ(promotion.FirstLevelRenewalSubscribeKey.Value()),
		promotionearnings.CommissionRuleKeyEQ(promotion.SecondLevelRenewalSubscribeKey.Value()),
	)).Count(s.ctx)
}
