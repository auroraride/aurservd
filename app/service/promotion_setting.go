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

	// ent.Database.PromotionLevelTask.Delete().ExecX(s.ctx)
	// // 不应该写这里临时用一下
	// _, err := ent.Database.ExecContext(context.Background(), `
	// INSERT INTO "public"."promotion_level_task" ("id", "created_at", "updated_at", "name", "description", "growth_value", "type","key") VALUES (313532612608, '2023-07-21 10:33:03+08', '2023-07-21 10:33:05+08', '一级团员新签', '一级下线用户注册并签约', 4, 1,'firstLevelNewSubscribe')`)
	// if err != nil {
	// 	return
	// }
	// _, err = ent.Database.ExecContext(context.Background(), `
	// INSERT INTO "public"."promotion_level_task" ("id", "created_at", "updated_at", "name", "description", "growth_value", "type","key") VALUES (313532612609, '2023-07-21 10:33:34+08', '2023-07-21 10:33:36+08', '一级团员续费', '一级下线用户续费', 2, 2,'firstLevelRenewalSubscribe')`)
	// if err != nil {
	// 	return
	// }
	// _, err = ent.Database.ExecContext(context.Background(), `
	// INSERT INTO "public"."promotion_level_task" ("id", "created_at", "updated_at", "name", "description", "growth_value", "type","key") VALUES (313532612610, '2023-07-21 10:33:56+08', '2023-07-21 10:33:57+08', '二级团员新签', '二级下线用户注册并签约', 2, 1,'secondLevelNewSubscribe')`)
	// if err != nil {
	// 	return
	// }
	// _, err = ent.Database.ExecContext(context.Background(), `
	// INSERT INTO "public"."promotion_level_task" ("id", "created_at", "updated_at", "name", "description", "growth_value", "type","key") VALUES (313532612612, '2023-07-21 10:34:19+08', '2023-07-21 10:34:20+08', '二级团员续费', '二级下线用户续费', 1, 2,'secondLevelRenewalSubscribe')`)
	// if err != nil {
	// 	return
	// }
	//

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

// UpdateSetting 修改会员设置
func (s *promotionSettingService) UpdateSetting(req *promotion.Setting) {
	_, err := ent.Database.PromotionSetting.Update().
		Where(promotionsetting.Key(req.Key.Value())).
		SetTitle(req.Title).
		SetContent(req.Context).
		Save(s.ctx)
	if err != nil {
		snag.Panic("修改会员设置失败")
	}
}
