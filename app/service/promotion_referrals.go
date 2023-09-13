package service

import (
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
)

type promotionReferralsService struct {
	*BaseService
}

func NewPromotionReferralsService() *promotionReferralsService {
	return &promotionReferralsService{
		BaseService: newService(),
	}
}

// MemberReferrals 记录会员推荐关系
func (s *promotionReferralsService) MemberReferrals(tx *ent.Tx, req promotion.Referrals) {
	tx.PromotionReferrals.Create().
		SetNillableReferringMemberID(req.ReferringMemberId).
		SetReferredMemberID(req.ReferredMemberId).
		SetNillableRiderID(req.RiderID).
		SetNillableReferralTime(req.ReferralTime).
		ExecX(s.ctx)
}
