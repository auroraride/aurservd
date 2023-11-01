package service

import (
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionreferralsprogress"
	"github.com/auroraride/aurservd/internal/ent/rider"
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
	rp, _ := ent.Database.PromotionReferralsProgress.Query().Where(
		promotionreferralsprogress.ReferredMemberID(*req.ReferredMemberId),
		// 这里不用查询已被绑定的状态 因为已被绑定的状态不会再次进入这个方法
		promotionreferralsprogress.Status(promotion.ReferralsStatusInviting.Value()),
	).First(s.ctx)

	if rp != nil {
		// 如果骑手重复绑定同一个邀请人 不重新记录
		if rp.ReferringMemberID != *req.ReferringMemberId {
			// 假如有关系待绑定 修改原有绑定关系为失效  重新新增一条为现有绑定关系
			rp.Update().Where(
				promotionreferralsprogress.ReferredMemberID(*req.ReferredMemberId),
				promotionreferralsprogress.Status(promotion.ReferralsStatusInviting.Value()),
			).
				SetStatus(promotion.ReferralsStatusFail.Value()).
				SetRemark("骑手已被他人邀请").
				ExecX(s.ctx)
			s.CreateReferralsProgress(req)
		}
		return
	}
	// 记录骑手邀请进度
	s.CreateReferralsProgress(req)
}

func (s *promotionReferralsService) CreateReferralsProgress(req *promotion.Referrals) {
	// 记录骑手邀请进度
	ent.Database.PromotionReferralsProgress.Create().
		SetNillableRiderID(req.RiderID).
		SetNillableReferringMemberID(req.ReferringMemberId).
		SetNillableReferredMemberID(req.ReferredMemberId).
		SetName(req.Name).
		SetStatus(req.Status.Value()).
		SetRemark(req.Remark).
		ExecX(s.ctx)
}

func (s *promotionReferralsService) UpdatedReferralsStatus(req *promotion.Referrals) {
	ent.Database.PromotionReferralsProgress.Update().
		Where(
			promotionreferralsprogress.ReferredMemberID(*req.ReferredMemberId),
			promotionreferralsprogress.Status(promotion.ReferralsStatusInviting.Value()),
		).SetStatus(req.Status.Value()).
		SetRemark(req.Remark).
		ExecX(s.ctx)
}

// ReferralsProgressList 推荐关系进度列表查询
func (s *promotionReferralsService) ReferralsProgressList(mem *ent.PromotionMember, req *promotion.ReferralsProgressReq) *model.PaginationRes {
	q := ent.Database.PromotionReferralsProgress.Query()

	if mem != nil {
		q.Where(promotionreferralsprogress.ReferringMemberID(mem.ID))
	}

	if req.Keyword != nil {
		q.Where(promotionreferralsprogress.Or(
			promotionreferralsprogress.NameContains(*req.Keyword),
			promotionreferralsprogress.HasRiderWith(rider.PhoneContains(*req.Keyword)),
		))
	}

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
				Status:    *item.Status,
				Remark:    item.Remark,
				CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			if item.Edges.Rider != nil {
				res.Phone = NewPromotionMemberService().MaskSensitiveInfo(item.Edges.Rider.Phone, 3, 4)
			}
			name := item.Name
			if item.Edges.Rider != nil && item.Edges.Rider.Name != "" {
				name = item.Edges.Rider.Name
			}
			res.Name = NewPromotionMemberService().MaskName(name)
			return res
		},
	)
}

// RiderBindReferrals 实名认证成功后，判断骑手是否可以绑定
func (s *promotionReferralsService) RiderBindReferrals(ri *ent.Rider) {
	// 推广小程序 判断骑手是否可以绑定 需要实名认证
	memId, _ := ent.Database.PromotionMember.Query().
		Where(promotionmember.RiderID(ri.ID)).
		Where(
			promotionmember.HasRiderWith(
				rider.PersonIDNotNil(),
			),
		).FirstID(s.ctx)
	if memId != 0 {
		referralsProgress := &promotion.Referrals{
			ReferredMemberId: &memId,
			Name:             ri.Name,
		}
		if ir := NewPromotionMemberService().IsRiderCanBind(ri); ir != promotion.MemberAllowBind {
			// 邀请失败
			referralsProgress.Remark = ir.String()
			referralsProgress.Status = promotion.ReferralsStatusFail
		} else {
			re, _ := ent.Database.PromotionReferralsProgress.Query().Where(promotionreferralsprogress.ReferredMemberID(memId)).First(s.ctx)
			if re != nil {
				// 邀请成功
				referralsProgress.Status = promotion.ReferralsStatusSuccess
				referralsProgress.Remark = "邀请成功"
				// 绑定关系
				NewPromotionReferralsService().CreateReferrals(&promotion.Referrals{
					ReferringMemberId: &re.ReferringMemberID,
					ReferredMemberId:  &re.ReferredMemberID,
					RiderID:           &ri.ID,
				})
			}

		}
		// 更新邀请进度
		NewPromotionReferralsService().UpdatedReferralsStatus(referralsProgress)
	}
}
