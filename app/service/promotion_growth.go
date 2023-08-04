package service

import (
	"context"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotiongrowth"
	"github.com/auroraride/aurservd/internal/ent/promotionleveltask"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/pkg/tools"
)

type promotionGrowthService struct {
	ctx context.Context
}

func NewPromotionGrowthService() *promotionGrowthService {
	return &promotionGrowthService{
		ctx: context.Background(),
	}
}

// List 会员成长值列表
func (s *promotionGrowthService) List(req *promotion.GrowthReq) *model.PaginationRes {
	q := ent.Database.PromotionGrowth.Query().WithMember().WithTask().Order(ent.Desc(promotiongrowth.FieldCreatedAt))

	if req.ID != nil {
		q.Where(promotiongrowth.MemberID(*req.ID))
	}

	if req.Keyword != nil {
		q.Where(
			promotiongrowth.HasMemberWith(
				promotionmember.Or(
					promotionmember.PhoneContainsFold(*req.Keyword),
					promotionmember.NameContainsFold(*req.Keyword),
				),
			),
		)
	}

	if req.Status != nil {
		q.Where(promotiongrowth.Status(*req.Status))
	}

	if req.LevelTaskID != nil {
		q.Where(promotiongrowth.HasTaskWith(promotionleveltask.ID(*req.LevelTaskID)))
	}

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			promotiongrowth.CreatedAtGTE(start),
			promotiongrowth.CreatedAtLTE(end),
		)
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.PromotionGrowth) (res promotion.GrowthRes) {
			res = promotion.GrowthRes{
				GrowthDetail: promotion.GrowthDetail{
					ID:          item.ID,
					Status:      item.Status,
					GrowthValue: item.GrowthValue,
					CreatedAt:   item.CreatedAt.Format(carbon.DateTimeLayout),
					Remark:      item.Remark,
				}}
			if item.Edges.Member != nil {
				res.Phone = item.Edges.Member.Phone
				res.Name = item.Edges.Member.Name
			}
			if item.Edges.Task != nil {
				res.LevelTaskName = item.Edges.Task.Name
			}
			return
		})
}

// Create 创建会员成长值
func (s *promotionGrowthService) Create(tx *ent.Tx, req *promotion.GrowthCreateReq) error {
	return tx.PromotionGrowth.Create().
		SetMemberID(req.MemberID).
		SetTaskID(req.TaksID).
		SetGrowthValue(req.GrowthValue).
		SetStatus(req.Status).
		Exec(s.ctx)
}

// Update 编辑会员成长值
// func (s *promotionGrowthService) Update(req *promotion.PromotionGrowthEditReq) {
// 	mg, _ := ent.Database.PromotionGrowth.Query().First(s.ctx)
// 	if mg == nil {
// 		snag.Panic("会员成长值不存在")
// 	}
// 	// 查询会员信息
// 	memberInfo, _ := NewPromotionService().GetMemberById(*mg.MemberID)
// 	if memberInfo != nil {
// 		snag.Panic("会员不存在")
// 	}
// 	currentLevelInfo := memberInfo.Edges.TemaLevel
// 	if currentLevelInfo == nil {
// 		snag.Panic("会员等级不存在")
// 	}
//
// 	deltaGrowthValue := req.GrowthValue
// 	if req.Status == promotion.PromotionGrowthStatusInvalid.Value() {
// 		deltaGrowthValue = -deltaGrowthValue
// 	}
//
// 	// 查询当前等级上下一级的升级条件
// 	var preLevel, nextLevel = currentLevelInfo.TemaLevel - 1, currentLevelInfo.TemaLevel + 1
//
// 	levelInfo, _ := ent.Database.Level.Query().
// 		Where(
// 			memberlevel.Or(
// 				memberlevel.TemaLevel(preLevel),
// 				memberlevel.TemaLevel(nextLevel),
// 			),
// 		).All(s.ctx)
// 	var levelID = currentLevelInfo.ID
// 	// TODO 这里会员等级达成条件可以修改 会有很多问题
// 	for _, level := range levelInfo {
// 		// 判断 会员成长值是否满足升降级条件  满足则升降级 否则级别不变
// 		if level.TemaLevel == preLevel && memberInfo.TotalGrowthValue+deltaGrowthValue < level.GrowthValue {
// 			levelID = level.ID
// 		} else if level.TemaLevel == nextLevel && memberInfo.TotalGrowthValue+deltaGrowthValue >= level.GrowthValue {
// 			levelID = level.ID
// 		}
// 	}
//
// 	// 更新会员信息
// 	_, err := ent.Database.Member.Update().Where(member.ID(*mg.MemberID)).
// 		AddTotalGrowthValue(int64(deltaGrowthValue)).
// 		SetLevelID(levelID).
// 		Save(s.ctx)
// 	if err != nil {
// 		snag.Panic("会员成长值更新失败")
// 	}
// 	_, err = ent.Database.MemberGrowth.UpdateOneID(mg.ID).SetStatus(req.Status).Save(s.ctx)
// 	if err != nil {
// 		snag.Panic("成长值状态更新失败")
// 	}
// }
