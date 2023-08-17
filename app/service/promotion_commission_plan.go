package service

import (
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotioncommissionplan"
)

type promotionCommissionPlanService struct {
	*BaseService
}

func NewPromotionCommissionPlanService(params ...any) *promotionCommissionPlanService {
	return &promotionCommissionPlanService{
		BaseService: newService(params...),
	}
}

func (s *promotionCommissionPlanService) GetCommissionPlan(planID uint64, commissionID []uint64, isCommissionCustom bool) *ent.PromotionCommissionPlan {
	q := ent.Database.PromotionCommissionPlan.QueryNotDeleted().
		Where(
			promotioncommissionplan.PlanID(planID),
			promotioncommissionplan.CommissionIDIn(commissionID...),
		)
	if isCommissionCustom {
		q.Where(
			promotioncommissionplan.HasPromotionCommissionWith(
				promotioncommission.Type(promotion.CommissionCustom.Value()),
				promotioncommission.Enable(true),
			),
		)
	} else {
		q.Where(
			promotioncommissionplan.HasPromotionCommissionWith(
				promotioncommission.TypeNEQ(promotion.CommissionCustom.Value()),
				promotioncommission.Enable(true),
			),
		)
	}
	commission, _ := q.WithPromotionCommission().First(s.ctx)
	return commission
}
