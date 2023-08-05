package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionreferrals"
)

type promotionReferralsService struct {
	ctx context.Context
}

func NewPromotionReferralsService() *promotionReferralsService {
	return &promotionReferralsService{
		ctx: context.Background(),
	}
}

// MemberReferrals 记录会员推荐关系
func (s *promotionReferralsService) MemberReferrals(tx *ent.Tx, req promotion.Referrals) {
	var ref *ent.PromotionReferrals
	if req.ReferringMemberId != nil {
		ref, _ = ent.Database.PromotionReferrals.Query().Where(promotionreferrals.ReferredMemberID(*req.ReferringMemberId)).First(s.ctx)
	}

	// 记录推荐关系
	q := tx.PromotionReferrals.Create().
		SetNillableReferringMemberID(req.ReferringMemberId).
		SetReferredMemberID(req.ReferredMemberId).
		SetNillableRiderID(req.RiderID)
	// 如果有推荐人，设置父级id
	if ref != nil {
		q.SetParentID(ref.ID)
	}
	q.ExecX(s.ctx)
}
