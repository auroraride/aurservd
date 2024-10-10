package model

// AssetEbikeAttributesColumns 资产电车导入固定列
var AssetEbikeAttributesColumns = []any{"型号", "车架号", "仓库"}

// AssetCreateReq 创建资产请求
type AssetCreateReq struct {
	AssetType     AssetType              `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	SN            *string                `json:"sn"`                            // 资产编号
	CityID        *uint64                `json:"cityId"`                        // 城市ID(AssetType为 2:智能电池 需要填写)
	LocationsType AssetLocationsType     `json:"locationsType" enums:"1"`       // 资产位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	LocationsID   uint64                 `json:"locationsId"`                   // 资产位置ID
	Attribute     []AssetAttributeCreate `json:"attribute"`                     // 属性
	Enable        *bool                  `json:"enable"`                        // 是否启用
	BrandID       *uint64                `json:"brandId"`                       // 品牌ID(AssetType为 1:电车 需要填写)
}

// AssetModifyReq 修改资产请求
type AssetModifyReq struct {
	ID        uint64                 `json:"id" param:"id" validate:"required"` // 资产ID
	Enable    *bool                  `json:"enable"`                            // 是否启用
	CityID    *uint64                `json:"cityId"`                            // 城市ID
	Remark    *string                `json:"remark"`                            // 备注
	BrandID   *uint64                `json:"brandId"`                           // 品牌ID
	Attribute []AssetAttributeUpdate `json:"attribute"`                         // 属性
}

// AssetFilter 资产筛选条件
type AssetFilter struct {
	SN               *string             `json:"sn" query:"sn"`                                           // 编号
	ModelID          *uint64             `json:"modelId" query:"modelId"`                                 // 型号ID
	CityID           *uint64             `json:"cityId" query:"cityId"`                                   // 城市
	OwnerType        *uint8              `json:"ownerType" query:"ownerType" enums:"1,2"`                 // 归属类型   1:平台 2:团签
	EnterpriseID     *uint64             `json:"enterpriseId" query:"enterpriseId"`                       // 团签企业ID
	StationID        *uint64             `json:"stationId" query:"stationId"`                             // 站点ID
	LocationsType    *AssetLocationsType `json:"locationsType" query:"locationsType" enums:"1,2,3,4,5,6"` // 资产位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	LocationsID      *uint64             `json:"locationsId" query:"locationsId"`                         // 资产位置ID
	LocationsKeyword *string             `json:"locationsKeyword" query:"locationsKeyword"`               // 资产位置关键词 只有LocationsType =（5:电柜 6:骑手）有效
	Status           *AssetStatus        `json:"status" query:"status" enums:"1,2,3,4,5"`                 // 资产状态 0:待入库 1:库存中 2:配送中 3:使用中 4:故障 5:报废
	Enable           *bool               `json:"enable" query:"enable"`                                   // 是否启用
	AssetType        *AssetType          `json:"assetType" query:"assetType" enums:"1,2,3,4,5,6"`         // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	BrandID          *uint64             `json:"brandId" query:"brandId"`                                 // 电车品牌ID
	Rto              *bool               `json:"rto" query:"rto"`                                         // 电车是否赠送
	Attribute        *string             `json:"attribute" query:"attribute"`                             // 属性 id:value
	AssetManagerID   uint64              `json:"assetManagerId" query:"assetManagerId"`                   // 仓管人员ID
	EmployeeID       uint64              `json:"employeeId" query:"employeeId"`                           // 门店店员ID
	Battery          *bool               `json:"battery" query:"battery"`                                 // 电池是否统一查询
	MaterialID       *uint64             `json:"materialId" query:"materialId"`                           // 其他物资ID
	Keyword          *string             `json:"keyword" query:"keyword"`                                 // 电车关键字
}

// AssetListReq 资产列表请求
type AssetListReq struct {
	PaginationReq
	AssetFilter
}

// AssetListRes 资产列表返回
type AssetListRes struct {
	ID             uint64                    `json:"id"`               // 资产ID
	CityName       string                    `json:"cityName"`         // 城市
	CityID         uint64                    `json:"cityId,omitempty"` // 城市ID
	Belong         string                    `json:"belong"`           // 归属
	AssetLocations string                    `json:"assetLocations"`   // 资产位置
	LocationsID    uint64                    `json:"locationsId"`      // 资产位置ID
	Brand          string                    `json:"brand"`            // 品牌
	BrandID        uint64                    `json:"brandId"`          // 品牌ID
	Model          string                    `json:"model"`            // 资产型号
	SN             string                    `json:"sn"`               // 编号
	AssetStatus    string                    `json:"assetStatus"`      // 资产状态(文字)
	Enable         bool                      `json:"enable"`           // 是否启用
	Remark         string                    `json:"remark"`           // 备注
	Attribute      map[uint64]AssetAttribute `json:"attribute"`        // 属性
	Status         AssetStatus               `json:"status"`           // 资产状态
	Rto            string                    `json:"rto"`              // 电车是否赠送
}

// AssetBatchCreateReq 批量创建资产请求
type AssetBatchCreateReq struct {
	AssetType AssetType `json:"assetType" validate:"required" form:"assetType" query:"assetType"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它

}

// AssetExportTemplateReq 导出模版请求
type AssetExportTemplateReq struct {
	AssetType AssetType `json:"assetType" validate:"required" query:"assetType"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
}

// AssetNumRes 资产有效数量返回
type AssetNumRes struct {
	Num       int        `json:"num"`                 // 有效数量
	AssetID   *uint64    `json:"assetId,omitempty"`   // 一个有效的资产ID
	AssetType *AssetType `json:"assetType,omitempty"` // 资产类型
}

// QueryAssetBatteryReq 查询资产请求
type QueryAssetBatteryReq struct {
	LocationsType *AssetLocationsType `json:"locationsType" query:"locationsType" enums:"1,2,3,4,5,6"` // 资产位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	LocationsID   *uint64             `json:"locationsId" query:"locationsId"`                         // 资产位置ID
	ModelID       uint64              `json:"modelId" query:"modelId"`                                 // 型号ID
}

// QueryAssetReq 查询资产请求
type QueryAssetReq struct {
	LocationsType AssetLocationsType `json:"locationsType" query:"locationsType" enums:"1,2,3,4,5,6"` // 资产位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	LocationsID   uint64             `json:"locationsId" query:"locationsId"`                         // 资产位置ID
	ID            *uint64            `json:"id" query:"id"`                                           // 资产ID
	Sn            *string            `json:"sn" query:"sn"`                                           // 资产编号
}
