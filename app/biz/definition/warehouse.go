package definition

import "github.com/auroraride/aurservd/app/model"

type WareHouseListReq struct {
	model.PaginationReq
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
	EbikeTotal            int                    `json:"ebikeTotal"`            // 电车总数
	Ebikes                []*model.StockMaterial `json:"ebikes"`                // 电车物资详情
	SmartBatteryTotal     int                    `json:"smartBatteryTotal"`     // 智能电池总数
	SmartBatteries        []*model.StockMaterial `json:"smartBatteries"`        // 智能电池物资详情
	NonSmartBatteryTotal  int                    `json:"nonSmartBatteryTotal"`  // 非智能电池总数
	NonSmartBatteries     []*model.StockMaterial `json:"nonSmartBatteries"`     // 非智能电池物资详情
	CabinetAccessoryTotal int                    `json:"cabinetAccessoryTotal"` // 电柜配件总数
	CabinetAccessories    []*model.StockMaterial `json:"cabinetAccessories"`    // 电柜配件物资详情
	EbikeAccessoryTotal   int                    `json:"ebikeAccessoryTotal"`   // 电车配件总数
	EbikeAccessories      []*model.StockMaterial `json:"ebikeAccessories"`      // 电车配件物资详情
	OtherAssetTotal       int                    `json:"otherAssetTotal"`       // 其他物资总数
	OtherAssets           []*model.StockMaterial `json:"otherAssets"`           // 其他物资详情
}
