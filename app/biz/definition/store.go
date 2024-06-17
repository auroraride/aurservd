package definition

import "github.com/auroraride/aurservd/app/model"

type StoreListReq struct {
	Lng          float64                  `json:"lng" query:"lng" validate:"required"` // 经度
	Lat          float64                  `json:"lat" query:"lat" validate:"required"` // 纬度
	Distance     *float64                 `json:"distance" query:"distance" `          // 距离
	CityID       *uint64                  `json:"cityId" query:"cityId"`               // 城市ID
	PlanID       *uint64                  `json:"planId" query:"planId"`               // 套餐ID
	Status       *model.StoreStatus       `json:"status" query:"status"`               // 门店状态 0维护 1营业 2休息 3隐藏
	BusinessType *model.StoreBusinessType `json:"businessType" query:"businessType"`   // 门店业务 1租车 2修车 3买车 4驿站
	Keyword      *string                  `json:"keyword" query:"keyword"`             // 关键字
}

// StoreDetail 门店信息
type StoreDetail struct {
	ID            uint64                        `json:"id"`                  // 门店ID
	Name          string                        `json:"name"`                // 门店名称
	Status        model.StoreStatus             `json:"status"`              // 门店状态 0维护 1营业 2休息 3隐藏
	City          model.City                    `json:"city"`                // 城市
	Employee      *model.Employee               `json:"employee,omitempty"`  // 店员, 有可能不存在
	EbikeObtain   bool                          `json:"ebikeObtain"`         // 是否可以领取车辆
	EbikeRepair   bool                          `json:"ebikeRepair"`         // 是否可以维修车辆
	EbikeSale     bool                          `json:"ebikeSale"`           // 是否可以买车
	BusinessHours string                        `json:"businessHours"`       // 营业时间
	Lng           float64                       `json:"lng"`                 // 经度
	Lat           float64                       `json:"lat"`                 // 纬度
	Distance      float64                       `json:"distance"`            // 距离(米)
	Address       string                        `json:"address"`             // 地址
	Rest          bool                          `json:"rest"`                // 是否拥有驿站
	Photos        []string                      `json:"photos"`              // 门店照片
	Brands        []*model.PlanEbikeBrandOption `json:"brands,omitempty"`    // 租车车电选项
	SaleGoods     []*GoodsDetail                `json:"saleGoods,omitempty"` // 买车列表
	RtoBrands     []*model.PlanEbikeBrandOption `json:"rtoBrands,omitempty"` // 以租代购车电选项
	Phone         string                        `json:"phone"`               // 门店电话
}

// StoreDetailReq 门店详情请求
type StoreDetailReq struct {
	model.IDParamReq
	Lng float64 `json:"lng" validate:"required" query:"lng"` // 经度
	Lat float64 `json:"lat" validate:"required" query:"lat"` // 纬度
}
