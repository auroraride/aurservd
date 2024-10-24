package service

import (
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionlevel"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"

	"github.com/auroraride/aurservd/pkg/snag"
)

type promotionLevelService struct {
	*BaseService
}

func NewPromotionLevelService(params ...any) *promotionLevelService {
	return &promotionLevelService{
		BaseService: newService(params...),
	}
}

// Level 会员等级列表
func (s *promotionLevelService) Level() []promotion.Level {
	item := ent.Database.PromotionLevel.QueryNotDeleted().Order(ent.Asc(promotionlevel.FieldLevel)).AllX(s.ctx)
	res := make([]promotion.Level, 0, len(item))
	for _, v := range item {
		res = append(res, promotion.Level{
			ID:              v.ID,
			Level:           &v.Level,
			GrowthValue:     &v.GrowthValue,
			CommissionRatio: &v.CommissionRatio,
		})
	}
	return res
}

// Create 新增会员等级
func (s *promotionLevelService) Create(req *promotion.Level) {
	ent.Database.PromotionLevel.Create().
		SetLevel(*req.Level).
		SetGrowthValue(*req.GrowthValue).
		SetCommissionRatio(*req.CommissionRatio).
		SaveX(s.ctx)
}

// Update 修改会员等级
func (s *promotionLevelService) Update(req *promotion.Level) {
	ent.Database.PromotionLevel.UpdateOneID(req.ID).
		SetLevel(*req.Level).
		SetGrowthValue(*req.GrowthValue).
		SetCommissionRatio(*req.CommissionRatio).
		SaveX(s.ctx)
}

// Delete 删除会员等级
func (s *promotionLevelService) Delete(id uint64) {
	// 查询是否有会员使用该等级
	count, _ := ent.Database.PromotionMember.Query().Where(promotionmember.LevelID(id)).Count(s.ctx)
	if count > 0 {
		snag.Panic("该会员等级已被使用，无法删除")
	}
	ent.Database.PromotionLevel.SoftDeleteOneID(id).SaveX(s.ctx)
}

// LevelSelection 会员等级选择
func (s *promotionLevelService) LevelSelection() []*promotion.LevelSelection {
	levels, _ := ent.Database.PromotionLevel.QueryNotDeleted().Order(ent.Asc(promotionlevel.FieldLevel)).All(s.ctx)
	var selections []*promotion.LevelSelection
	for _, level := range levels {
		selections = append(selections, &promotion.LevelSelection{
			ID:    level.ID,
			Level: level.Level,
		})
	}
	return selections
}
