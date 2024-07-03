package definition

import "github.com/auroraride/aurservd/app/model"

// EbikeDetailReq 车电详情请求
type EbikeDetailReq struct {
	model.IDParamReq
	StoreID uint64  `json:"storeId" validate:"required" query:"storeId"` // 门店ID
	PlanID  uint64  `json:"planId" query:"planId" validate:"required"`   // 套餐ID
	Lat     float64 `json:"lat" validate:"required" query:"lat"`         // 纬度
	Lng     float64 `json:"lng" validate:"required" query:"lng"`         // 经度
}

// EbikeDetailRes 车电详情返回
type EbikeDetailRes struct {
	Plan  model.Plan `json:"plan"` // 套餐
	Store struct {
		model.StoreLngLat
		Address  string  `json:"address"`  // 地址
		Distance float64 `json:"distance"` // 距离
	} `json:"store"` // 门店
	Brand model.EbikeBrand `json:"brand"` // 车电品牌
}

// EbikeBrandDeleteReq 删除车电品牌
type EbikeBrandDeleteReq struct {
	ID uint64 `json:"id" validate:"required" param:"id"` // 车电品牌ID
}

// EbikeBatchModifyReq 批量修改电车型号
type EbikeBatchModifyReq struct {
	SN      []string `json:"sn"  validate:"required"`     // 车架号
	BrandID uint64   `json:"brandId" validate:"required"` // 品牌ID
}
