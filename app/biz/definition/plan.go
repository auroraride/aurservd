package definition

import "github.com/auroraride/aurservd/app/model"

type PlanNewlyRes struct {
	Brands    []*model.PlanEbikeBrandOption `json:"brands,omitempty"`    // 车电选项
	Models    []*model.PlanModelOption      `json:"models,omitempty"`    // 单电选项
	Configure *model.PaymentConfigure       `json:"configure,omitempty"` // 支付配置

	BatteryDescription model.SettingPlanDescription `json:"batteryDescription"` // 单电介绍
	EbikeDescription   model.SettingPlanDescription `json:"ebikeDescription"`   // 车电介绍

	RtoBrands []*model.PlanEbikeBrandOption `json:"rtoBrands,omitempty"` // 以租代购车电选项
}
