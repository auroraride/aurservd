package service

import (
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionleveltask"
)

type promotionLevelTaskService struct {
	*BaseService
}

func NewPromotionLevelTaskService(params ...any) *promotionLevelTaskService {
	return &promotionLevelTaskService{
		BaseService: newService(params...),
	}
}

// List 会员任务列表
func (s *promotionLevelTaskService) List() []promotion.LevelTask {
	item, _ := ent.Database.PromotionLevelTask.Query().Order(ent.Asc(promotionleveltask.FieldCreatedAt)).All(s.ctx)
	res := make([]promotion.LevelTask, 0, len(item))
	for _, v := range item {
		res = append(res, promotion.LevelTask{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Type:        promotion.LevelTaskType(v.Type),
			GrowthValue: v.GrowthValue,
		})
	}
	return res
}

// Update 修改任务列表
func (s *promotionLevelTaskService) Update(req *promotion.LevelTask) {
	ent.Database.PromotionLevelTask.UpdateOneID(req.ID).
		SetGrowthValue(req.GrowthValue).
		SaveX(s.ctx)
}

// LevelTaskSelect 任务选择
func (s *promotionLevelTaskService) LevelTaskSelect() []*promotion.LevelTaskSelect {
	res := make([]*promotion.LevelTaskSelect, 0)
	lt := ent.Database.PromotionLevelTask.Query().AllX(s.ctx)
	for _, v := range lt {
		res = append(res, &promotion.LevelTaskSelect{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return res
}

// QueryByKey 通过key查询成长值
func (s *promotionLevelTaskService) QueryByKey(key string) (*ent.PromotionLevelTask, error) {
	return ent.Database.PromotionLevelTask.Query().Where(promotionleveltask.Key(key)).First(s.ctx)
}
