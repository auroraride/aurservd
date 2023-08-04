package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionleveltask"
)

type promotionLevelTaskService struct {
	ctx context.Context
	*BaseService
}

func NewPromotionLevelTaskService(params ...any) *promotionLevelTaskService {
	return &promotionLevelTaskService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
}

// TaskList 会员任务列表
func (s *promotionLevelTaskService) TaskList() []promotion.LevelTask {
	item := ent.Database.PromotionLevelTask.Query().Order(ent.Desc(promotionleveltask.FieldCreatedAt)).AllX(s.ctx)
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

// UpdateTask 修改任务列表
func (s *promotionLevelTaskService) UpdateTask(req *promotion.LevelTask) {
	ent.Database.PromotionLevelTask.UpdateOneID(req.ID).
		SetGrowthValue(req.GrowthValue).
		SaveX(s.ctx)
}

// LevelTaskSelect 任务选择
func (s *promotionLevelTaskService) LevelTaskSelect() []*promotion.LevelTaskSelect {
	ent.Database.PromotionLevelTask.Query().AllX(s.ctx)
	return nil
}

// QueryByKey 通过key查询成长值
func (s *promotionLevelTaskService) QueryByKey(key string) (*ent.PromotionLevelTask, error) {
	return ent.Database.PromotionLevelTask.Query().Where(promotionleveltask.Key(key)).First(s.ctx)
}
