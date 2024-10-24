package service

import (
	"github.com/labstack/echo/v4"

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
	if req.ReferringMemberID != nil {
		// 获取推荐会员信息
		pinfo, _ := NewPromotionMemberService().GetMemberById(*req.ReferringMemberID)
		if pinfo == nil {
			// 上级会员不存在，设置上级会员为nil
			req.ReferringMemberID = nil
		}
	}

	// 如果推荐人是自己，设置推荐人为nil
	if req.ReferringMemberID != nil && req.ReferringMemberID == req.ReferredMemberID {
		req.ReferringMemberID = nil
	}

	// 创建推荐关系
	NewPromotionReferralsService().MemberReferrals(&promotion.Referrals{
		ReferringMemberID: req.ReferringMemberID,
		ReferredMemberID:  req.ReferredMemberID,
		RiderID:           req.RiderID,
	})
}

// MemberReferrals 记录会员推荐关系
func (s *promotionReferralsService) MemberReferrals(req *promotion.Referrals) {
	ent.Database.PromotionReferrals.Create().
		SetNillableReferringMemberID(req.ReferringMemberID).
		SetReferredMemberID(*req.ReferredMemberID).
		SetNillableRiderID(req.RiderID).
		ExecX(s.ctx)
}

// MemberReferralsProgress 记录会员关系进度
func (s *promotionReferralsService) MemberReferralsProgress(req *promotion.Referrals) {
	rp, _ := ent.Database.PromotionReferralsProgress.Query().Where(
		promotionreferralsprogress.ReferredMemberID(*req.ReferredMemberID),
		// 这里不用查询已被绑定的状态 因为已被绑定的状态不会再次进入这个方法
		promotionreferralsprogress.Status(promotion.ReferralsStatusInviting.Value()),
	).First(s.ctx)

	if rp != nil {
		// 如果骑手重复绑定同一个邀请人 不重新记录
		if rp.ReferringMemberID != *req.ReferringMemberID {
			// 假如有关系待绑定 修改原有绑定关系为失效  重新新增一条为现有绑定关系
			rp.Update().Where(
				promotionreferralsprogress.ReferredMemberID(*req.ReferredMemberID),
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
		SetNillableReferringMemberID(req.ReferringMemberID).
		SetNillableReferredMemberID(req.ReferredMemberID).
		SetName(req.Name).
		SetStatus(req.Status.Value()).
		SetRemark(req.Remark).
		ExecX(s.ctx)
}

func (s *promotionReferralsService) UpdatedReferralsStatus(req *promotion.Referrals) {
	ent.Database.PromotionReferralsProgress.Update().
		Where(
			promotionreferralsprogress.ReferredMemberID(*req.ReferredMemberID),
			promotionreferralsprogress.Status(promotion.ReferralsStatusInviting.Value()),
		).SetStatus(req.Status.Value()).
		SetRemark(req.Remark).
		ExecX(s.ctx)
}

// ReferralsProgressList 推荐关系进度列表查询
func (s *promotionReferralsService) ReferralsProgressList(ctx echo.Context, req *promotion.ReferralsProgressReq) *model.PaginationRes {
	q := ent.Database.PromotionReferralsProgress.Query()

	if req.MemberID != nil {
		q.Where(promotionreferralsprogress.ReferringMemberID(*req.MemberID))
	}

	if req.Keyword != nil {
		q.Where(promotionreferralsprogress.Or(
			promotionreferralsprogress.HasRiderWith(rider.NameContains(*req.Keyword)),
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

			name := item.Name
			if item.Edges.Rider != nil && item.Edges.Rider.Name != "" {
				name = item.Edges.Rider.Name
			}

			var phone string
			if item.Edges.Rider != nil {
				phone = item.Edges.Rider.Phone
			}

			if ctx.Path() != "/manager/v1/promotion/progress/list/:id" {
				phone = NewPromotionMemberService().MaskSensitiveInfo(phone, 3, 4)
				name = NewPromotionMemberService().MaskName(name)
			}
			res.Name = name
			res.Phone = phone
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
			ReferredMemberID: &memId,
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
					ReferringMemberID: &re.ReferringMemberID,
					ReferredMemberID:  &re.ReferredMemberID,
					RiderID:           &ri.ID,
				})
			}

		}
		// 更新邀请进度
		NewPromotionReferralsService().UpdatedReferralsStatus(referralsProgress)
	}
}
