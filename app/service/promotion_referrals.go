package service

import (
	"github.com/auroraride/aurservd/internal/ent/promotionreferralsprogress"
	"github.com/auroraride/aurservd/pkg/tools"

	"github.com/auroraride/aurservd/app/model"
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

// CreateReferrals 处理推荐关系
func (s *promotionReferralsService) CreateReferrals(req *promotion.Referrals) {
	if req.ReferringMemberId != nil {
		// 获取推荐会员信息
		pinfo, _ := NewPromotionMemberService().GetMemberById(*req.ReferringMemberId)
		if pinfo == nil {
			// 上级会员不存在，设置上级会员为nil
			req.ReferringMemberId = nil
		}
	}

	// 如果推荐人是自己，设置推荐人为nil
	if req.ReferringMemberId != nil && req.ReferringMemberId == req.ReferredMemberId {
		req.ReferringMemberId = nil
	}

	// 创建推荐关系
	NewPromotionReferralsService().MemberReferrals(&promotion.Referrals{
		ReferringMemberId: req.ReferringMemberId,
		ReferredMemberId:  req.ReferredMemberId,
		RiderID:           req.RiderID,
	})
}

// MemberReferrals 记录会员推荐关系
func (s *promotionReferralsService) MemberReferrals(req *promotion.Referrals) {
	ent.Database.PromotionReferrals.Create().
		SetNillableReferringMemberID(req.ReferringMemberId).
		SetReferredMemberID(*req.ReferredMemberId).
		SetNillableRiderID(req.RiderID).
		ExecX(s.ctx)
}

// MemberReferralsProgress 记录会员关系进度
func (s *promotionReferralsService) MemberReferralsProgress(req *promotion.Referrals) {
	ent.Database.PromotionReferralsProgress.Create().
		SetNillableRiderID(req.RiderID).
		SetNillableReferringMemberID(req.ReferringMemberId).
		SetReferredMemberID(*req.ReferredMemberId).
		SetStatus(req.Status.Value()).
		ExecX(s.ctx)
}

func (s *promotionReferralsService) UpdatedReferralsStatus(req *promotion.Referrals) {
	ent.Database.PromotionReferralsProgress.Update().
		Where(
			promotionreferralsprogress.ReferredMemberID(*req.ReferredMemberId),
		).SetStatus(req.Status.Value()).
		SetRemark(req.Remark).
		ExecX(s.ctx)
}

// ReferralsProgressList 推荐关系进度列表查询
func (s *promotionReferralsService) ReferralsProgressList(mem *ent.PromotionMember, req *promotion.ReferralsProgressReq) *model.PaginationRes {
	q := ent.Database.PromotionReferralsProgress.Query().
		Where(promotionreferralsprogress.ReferringMemberID(mem.ID))

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(promotionreferralsprogress.CreatedAtGTE(start), promotionreferralsprogress.CreatedAtLTE(end))
	}

	if req.Status != nil {
		q.Where(promotionreferralsprogress.Status(req.Status.Value()))
	}

	q.WithRider().Order(ent.Desc(promotionreferralsprogress.FieldCreatedAt))

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.PromotionReferralsProgress) (res promotion.ReferralsProgressRes) {
			res = promotion.ReferralsProgressRes{
				Name:      item.Name,
				Status:    promotion.ReferralsStatus(*item.Status).String(),
				Remark:    item.Remark,
				CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			if item.Edges.Rider != nil {
				res.Phone = NewPromotionMemberService().MaskSensitiveInfo(item.Edges.Rider.Phone, 3, 4)
			}
			return res
		},
	)
}
