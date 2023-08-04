package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
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
	// 设置默认佣金配置
	for k, v := range promotion.Settings {
		ms, _ := ent.Database.PromotionSetting.Query().Where(promotionsetting.Key(k.Value())).First(s.ctx)
		if ms == nil {
			_, err := ent.Database.PromotionSetting.Create().SetKey(k.Value()).SetTitle(v.Title).SetContent(v.Context).Save(s.ctx)
			if err != nil {
				zap.L().Fatal("营销设置初始化失败", zap.Error(err))
			}
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
	// _, err = ent.Database.ExecContext(context.Background(), `
	// INSERT INTO "public"."promotion_commission" ("id", "created_at", "updated_at", "deleted_at", "creator", "last_modifier", "remark", "type", "name", "rule", "enable", "use_count", "amount_sum", "desc", "history_id", "member_id") VALUES (317827579906, '2023-07-19 18:03:03+08', '2023-07-20 10:26:12.810352+08', NULL, NULL, NULL, NULL, 0, '推广返佣(全局)', '{"newUserCommission": {"firstLevelNewSubscribe": {"desc": "邀请好友注册极光出行", "name": "一级团员新签任务", "ratio": [20], "limitedType": 0}, "secondLevelNewSubscribe": {"desc": "好友邀请的二级团员注册极光出行", "name": "二级团员新签任务", "ratio": [5], "limitedType": 0}}, "renewalCommission": {"firstLevelRenewalSubscribe": {"desc": "已激活的好友再次续费", "name": "一级团员续费任务", "ratio": [5], "limitedType": 2}, "secondLevelRenewalSubscribe": {"desc": "已激活的二级团员再次续费", "name": "二级团员续费任务", "ratio": [1, 1, 1], "limitedType": 1}}}', 't', 3, 0, '', NULL, NULL)`)
	// if err != nil {
	// 	return
	// }

}

// Setting 获取会员设置
func (s *promotionSettingService) Setting(req *promotion.SettingReq) *promotion.Setting {
	item := ent.Database.PromotionSetting.Query().Where(promotionsetting.Key(req.Key.Value())).FirstX(s.ctx)
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
