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
	EbikeTotal            int              `json:"ebikeTotal"`            // 电车总数
	Ebikes                []*AssetMaterial `json:"ebikes"`                // 电车物资详情
	SmartBatteryTotal     int              `json:"smartBatteryTotal"`     // 智能电池总数
	SmartBatteries        []*AssetMaterial `json:"smartBatteries"`        // 智能电池物资详情
	NonSmartBatteryTotal  int              `json:"nonSmartBatteryTotal"`  // 非智能电池总数
	NonSmartBatteries     []*AssetMaterial `json:"nonSmartBatteries"`     // 非智能电池物资详情
	CabinetAccessoryTotal int              `json:"cabinetAccessoryTotal"` // 电柜配件总数
	CabinetAccessories    []*AssetMaterial `json:"cabinetAccessories"`    // 电柜配件物资详情
	EbikeAccessoryTotal   int              `json:"ebikeAccessoryTotal"`   // 电车配件总数
	EbikeAccessories      []*AssetMaterial `json:"ebikeAccessories"`      // 电车配件物资详情
	OtherAssetTotal       int              `json:"otherAssetTotal"`       // 其他物资总数
	OtherAssets           []*AssetMaterial `json:"otherAssets"`           // 其他物资详情
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

type WarehousePeopleSigninReq struct {
	Phone      string `json:"phone,omitempty" validate:"required_if=SigninType 1" trans:"电话"`
	SmsId      string `json:"smsId,omitempty" validate:"required_if=SigninType 1" trans:"短信ID"`
	Code       string `json:"code,omitempty" validate:"required_if=SigninType 1,required_if=SigninType 2" trans:"验证码"`
	SigninType uint64 `json:"signinType" validate:"required,oneof=1 2"`
}

type WarehousePeopleSigninRes struct {
	Profile WarehousePeopleProfile `json:"profile"`
	Token   string                 `json:"token"`
}

type OpenidReq struct {
	Code string `json:"code" query:"code"`
}

type OpenidRes struct {
	Openid string `json:"openid"`
}

type WarehousePeopleProfile struct {
	ID    uint64 `json:"id"`
	Phone string `json:"phone"` // 手机号
	Name  string `json:"name"`  // 姓名
}

// TransferListReq 调拨记录筛选条件
type TransferListReq struct {
	model.PaginationReq
	FromLocationType *model.AssetLocationsType  `json:"fromLocationType" query:"fromLocationType" enums:"1,2,3,4"` // 调拨前位置类型  1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	FromLocationID   *uint64                    `json:"fromLocationID" query:"fromLocationID"`                     // 调拨前位置ID
	ToLocationType   *model.AssetLocationsType  `json:"toLocationType" query:"toLocationType" enums:"1,2,3,4"`     // 调拨后位置类型  1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	ToLocationID     *uint64                    `json:"toLocationID" query:"toLocationID"`                         // 调拨后位置ID
	Status           *model.AssetTransferStatus `json:"status" query:"status" enums:"1,2,3,4"`                     // 调拨状态 1:配送中 2:待入库 3:已入库 4:已取消
	OutStart         *string                    `json:"outStart" query:"outStart"`                                 // 出库开始时间
	OutEnd           *string                    `json:"outEnd" query:"outEnd"`                                     // 出库结束时间
	Keyword          *string                    `json:"keyword" query:"keyword"`                                   // 关键字 (调拨单号，调拨事由、出库人、接收人)
}
