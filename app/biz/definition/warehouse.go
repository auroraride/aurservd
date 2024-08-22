package definition

import "github.com/auroraride/aurservd/app/model"

type WareHouseListReq struct {
	model.PaginationReq
	Keyword *string `json:"keyword" query:"keyword"` // 关键字
	CityID  *uint64 `json:"cityId" query:"cityId"`   // 城市ID
}

// WarehouseDetail 仓库信息
type WarehouseDetail struct {
	ID       uint64     `json:"id"`       // 门店ID
	Name     string     `json:"name"`     // 门店名称
	City     model.City `json:"city"`     // 城市
	Lng      float64    `json:"lng"`      // 经度
	Lat      float64    `json:"lat"`      // 纬度
	Address  string     `json:"address"`  // 地址
	QRCode   string     `json:"qrcode"`   // 仓库二维码
	CityID   uint64     `json:"cityId"`   // 城市ID
	CityName string     `json:"cityName"` // 城市名称
}

// WarehouseCreateReq 创建
type WarehouseCreateReq struct {
	Name    string `json:"name" validate:"required,max=30" trans:"仓库名称"` // 商品名称
	CityID  uint64 `json:"cityId" validate:"required" trans:"城市ID"`      // 城市ID
	Address string `json:"address" validate:"max=50"`                    // 详细地址
	Remark  string `json:"remark" validate:"max=50"`                     // 备注
}

// WarehouseModifyReq 修改
type WarehouseModifyReq struct {
	model.IDParamReq
	WarehouseCreateReq
}

// WareHouseAssetListReq 仓库资产列表请求
type WareHouseAssetListReq struct {
	model.PaginationReq
	Name      *string `json:"name" query:"name"`           // 仓库名称关键字
	CityID    *uint64 `json:"cityId" query:"cityId"`       // 城市ID
	ModelID   *uint64 `json:"modelID" query:"modelID"`     // 电池型号ID
	BrandId   *uint64 `json:"brandId" query:"brandId"`     // 电车型号ID
	OtherName *string `json:"otherName" query:"otherName"` // 其他物资名称
	Start     *string `json:"start" query:"start"`         // 开始时间
	End       *string `json:"end" query:"end"`             // 结束时间
}

// WareHouseAssetDetail 仓库资产信息
type WareHouseAssetDetail struct {
	ID             uint64         `json:"id"`             // 仓库ID
	Name           string         `json:"name"`           // 仓库名称
	City           model.City     `json:"city"`           // 城市
	Lng            float64        `json:"lng"`            // 经度
	Lat            float64        `json:"lat"`            // 纬度
	WarehouseAsset WarehouseAsset `json:"warehouseAsset"` // 仓库资产
}

// WarehouseAsset 仓库资产
type WarehouseAsset struct {
	CommonAssetDetail
}

// WarehouseByCityRes 城市仓库信息
type WarehouseByCityRes struct {
	City          model.City               `json:"city"`          // 城市
	WarehouseList []*WarehouseByCityDetail `json:"warehouseList"` // 仓库信息
}

// WarehouseByCityDetail 城市仓库信息详情
type WarehouseByCityDetail struct {
	ID   uint64 `json:"id"`   // 仓库ID
	Name string `json:"name"` // 仓库名称
}
