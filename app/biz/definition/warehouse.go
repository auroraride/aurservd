package definition

import "github.com/auroraride/aurservd/app/model"

type WareHouseListReq struct {
	Keyword *string `json:"keyword" query:"keyword"` // 关键字
	CityID  *uint64 `json:"cityId" query:"cityId"`   // 城市ID
}

// WarehouseDetail 仓库信息
type WarehouseDetail struct {
	ID      uint64     `json:"id"`      // 门店ID
	Name    string     `json:"name"`    // 门店名称
	City    model.City `json:"city"`    // 城市
	Lng     float64    `json:"lng"`     // 经度
	Lat     float64    `json:"lat"`     // 纬度
	Address string     `json:"address"` // 地址
	QRCode  string     `json:"qrcode"`  // 仓库二维码
}

// WarehouseCreateReq 创建
type WarehouseCreateReq struct {
	Name    string  `json:"name" validate:"required,max=30" trans:"仓库名称"` // 商品名称
	CityID  uint64  `json:"cityId" validate:"required" trans:"城市ID"`      // 城市ID
	Lat     float64 `json:"lat" validate:"required" trans:"城市纬度"`         // 纬度
	Lng     float64 `json:"lng" validate:"required" trans:"城市经度"`         // 经度
	Address string  `json:"address" validate:"max=50"`                    // 详细地址
	Remark  string  `json:"remark" validate:"max=50"`                     // 备注
}

// WarehouseModifyReq 修改
type WarehouseModifyReq struct {
	model.IDParamReq
	WarehouseCreateReq
}
