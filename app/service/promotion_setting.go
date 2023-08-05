package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotionsetting"
	"github.com/auroraride/aurservd/pkg/snag"
)

type promotionSettingService struct {
	ctx context.Context
	*BaseService
}

func NewPromotionSettingService(params ...any) *promotionSettingService {
	return &promotionSettingService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
}

func (s *promotionSettingService) Initialize() {

	for k, v := range promotion.Settings {
		ms, _ := ent.Database.PromotionSetting.Query().Where(promotionsetting.Key(k.Value())).First(s.ctx)
		if ms == nil {
			_, err := ent.Database.PromotionSetting.Create().SetKey(k.Value()).SetTitle(v.Title).SetContent(v.Context).Save(s.ctx)
			if err != nil {
				zap.L().Fatal("营销设置初始化失败", zap.Error(err))
			}
		}
	}

	// 查询是否有推广返佣(全局)的配置，没有则创建
	com, _ := ent.Database.PromotionCommission.QueryNotDeleted().Where(promotioncommission.Type(0)).First(s.ctx)
	if com == nil {
		ent.Database.PromotionCommission.Create().
			SetType(0).
			SetName("推广返佣(全局)").
			SetRule(&promotion.CommissionRule{}).
			SetEnable(true).
			SetStartAt(time.Now()).
			SaveX(s.ctx)
	}

	// 初始化推广等级任务
	existingKeys := make(map[string]bool)

	lt := ent.Database.PromotionLevelTask.Query().AllX(s.ctx)

	for _, v := range lt {
		existingKeys[v.Key] = true
	}

	for k := range promotion.CommissionRuleKeyNames {
		if _, exists := existingKeys[k.Value()]; !exists {
			q := ent.Database.PromotionLevelTask.Create().
				SetKey(k.Value()).
				SetName(promotion.CommissionRuleKeyNames[k]).
				SetDescription(promotion.CommissionRuleKeyNames[k])

			if k == promotion.FirstLevelNewSubscribeKey || k == promotion.SecondLevelNewSubscribeKey {
				q.SetType(promotion.LevelTaskTypeSign.Value())
			} else {
				q.SetType(promotion.LevelTaskTypeRenew.Value())
			}

			q.SaveX(s.ctx)
		}
	}
}

// Setting 获取会员设置
func (s *promotionSettingService) Setting(req *promotion.SettingReq) *promotion.Setting {
	item, _ := ent.Database.PromotionSetting.Query().Where(promotionsetting.Key(req.Key.Value())).First(s.ctx)
	if item == nil {
		return nil
	}
	return &promotion.Setting{
		Key:     promotion.SettingKey(item.Key),
		Title:   item.Title,
		Context: item.Content,
	}
}

// Update 修改会员设置
func (s *promotionSettingService) Update(req *promotion.Setting) {
	_, err := ent.Database.PromotionSetting.Update().
		Where(promotionsetting.Key(req.Key.Value())).
		SetTitle(req.Title).
		SetContent(req.Context).
		Save(s.ctx)
	if err != nil {
		snag.Panic("修改会员设置失败")
	}
}
