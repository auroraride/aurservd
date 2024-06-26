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

// PlanDetailReq 套餐详情
type PlanDetailReq struct {
	model.IDParamReq
}

// PlanDetailRes 套餐详情
type PlanDetailRes struct {
	model.Plan
	Notes []string `json:"notes"` // 购买须知
}

var DefaultMaxDistance = 50000.0

type StorePlanSortType uint8

const (
	StorePlanSortTypeDistance StorePlanSortType = iota + 1 // 距离排序
	StorePlanSortTypePrice                                 // 价格排序
)

func (t StorePlanSortType) Value() uint8 {
	return uint8(t)
}

// StorePlanReq 门店套餐请求
type StorePlanReq struct {
	CityId   uint64             `json:"cityId" query:"cityId" validate:"required" trans:"城市"`
	Lng      float64            `json:"lng" query:"lng" validate:"required" trans:"经度"`
	Lat      float64            `json:"lat" query:"lat" validate:"required" trans:"纬度"`
	Distance *float64           `json:"distance" query:"distance"`             // 距离
	SortType *StorePlanSortType `json:"sortType" query:"sortType" enums:"1,2"` // 排序方式 1-距离优先 2-价格优先
	BrandId  *uint64            `json:"brandId"`                               // 电车品牌ID
}

type ListByStoreRes struct {
	StorePlan []*StoreEbikePlan `json:"storePlan"` // 门店车电套餐
}

type StoreEbikePlan struct {
	StoreId    uint64  `json:"storeId"`         // 门店ID
	StoreName  string  `json:"storeName"`       // 门店名称
	Distance   float64 `json:"distance"`        // 门店距离
	PlanId     uint64  `json:"planId"`          // 套餐ID
	BrandId    uint64  `json:"brandId"`         // 品牌ID
	BrandName  string  `json:"brandName"`       // 品牌名称
	Cover      string  `json:"cover,omitempty"` // 品牌封面图
	Rto        bool    `json:"rto"`             // 是否以租代购套餐
	Daily      bool    `json:"daily"`           // 是否日租套餐
	DailyPrice float64 `json:"dailyPrice"`      // 套餐日租价格
	MonthPrice float64 `json:"monthPrice"`      // 套餐月租价格
}

// StorePlanDetailReq 门店套餐详情请求
type StorePlanDetailReq struct {
	StoreId uint64 `json:"storeId" query:"storeId" validate:"required" trans:"门店ID"`
	PlanId  uint64 `json:"planId" query:"planId" validate:"required" trans:"套餐ID"`
}

// StorePlanDetail 门店电车套餐详情
type StorePlanDetail struct {
	Children         *model.PlanModelOptions      `json:"children"`            // 子项
	Name             string                       `json:"name"`                // 名称
	Cover            string                       `json:"cover,omitempty"`     // 封面图
	Configure        *model.PaymentConfigure      `json:"configure,omitempty"` // 支付配置
	EbikeDescription model.SettingPlanDescription `json:"ebikeDescription"`    // 车电介绍
}
