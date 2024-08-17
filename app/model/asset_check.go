package model

// AssetCheckCreateReq 创建资产盘点请求
type AssetCheckCreateReq struct {
	LocationsType          AssetLocationsType       `json:"locationsType"`                             // 位置类型 1:仓库 2:门店 3:站点
	LocationsID            uint64                   `json:"locationsId"`                               // 位置ID
	OpratorID              uint64                   `json:"opratorId"`                                 // 操作人ID
	OpratorType            AssetOperateRoleType     `json:"opratorType"`                               // 操作人类型 1:资产后台(仓库) 2:门店 3:代理
	AssetCheckCreateDetail []AssetCheckCreateDetail `json:"details" validate:"required,dive,required"` // 资产盘点请求详情
}

// AssetCheckCreateDetail 资产盘点请求详情
type AssetCheckCreateDetail struct {
	AssetID   uint64    `json:"assetId" validate:"required"`   // 资产ID
	AssetType AssetType `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池
}

// AssetCheckGetReq 获取资产盘点请求
type AssetCheckGetReq struct {
	SN            string               `json:"sn" query:"sn" validate:"required"` // 资产编号
	LocationsID   uint64               `json:"locationsId" validate:"required"`   // 位置ID
	LocationsType AssetLocationsType   `json:"locationsType" validate:"required"` // 位置类型 1:仓库 2:门店 3:站点
	OpratorID     uint64               `json:"opratorId"`                         // 操作人ID
	OpratorType   AssetOperateRoleType `json:"opratorType"`                       // 操作人类型 1:资产后台(仓库) 2:门店 3:代理
}

// AssetCheckGetRes 获取资产盘点返回
type AssetCheckGetRes struct {
	StartAt        string               `json:"startAt"`        // 盘点开始时间
	EndAt          string               `json:"endAt"`          // 盘点结束时间
	OpratorID      uint64               `json:"opratorId"`      // 操作人ID
	OpratorName    string               `json:"opratorName"`    // 操作人名称
	BatteryNum     uint                 `json:"batteryNum"`     // 应盘点电池数量
	BatteryNumReal uint                 `json:"batteryNumReal"` // 实盘电池数量
	EbikeNum       uint                 `json:"ebikeNum"`       // 应盘点电车数量
	EbikeNumReal   uint                 `json:"ebikeNumReal"`   // 实盘电车数量
	LocationsID    uint64               `json:"locationsId"`    // 位置ID
	LocationsType  uint8                `json:"locationsType"`  // 位置类型
	CheckResult    uint8                `json:"checkResult"`    // 盘点结果 1:正常 2:异常
	Abnormal       []AssetCheckAbnormal `json:"abnormal"`       // 异常资产
}

// AssetCheckAbnormal 异常资产
type AssetCheckAbnormal struct {
	AssetID   uint64 `json:"assetId"`   // 资产ID
	AssetType uint8  `json:"assetType"` // 资产类型 1:电车 2:电池
	Name      string `json:"name"`      // 名称
	Model     string `json:"model"`     // 型号
	Brand     string `json:"brand"`     // 品牌
	Loss      bool   `json:"loss"`      // 是否丢失 true:是 false:否
}

// AssetCheckDetailLocations 资产明细门店站点信息
type AssetCheckDetailLocations struct {
	Name    string             `json:"name"`    // 名称
	Sum     uint               `json:"sum"`     // 数量
	Details []AssetCheckDetail `json:"details"` // 资产明细
}

// AssetCheckDetail 资产明细
type AssetCheckDetail struct {
	AssetID       uint64 `json:"assetId"`       // 资产ID
	AssetSN       string `json:"assetSN"`       // 资产编号
	Model         string `json:"model"`         // 型号
	BrandName     string `json:"brandName"`     // 品牌
	LocationsName string `json:"locationsName"` // 位置名称
	AssetStatus   uint8  `json:"assetStatus"`   // 资产状态
}

// AssetCheckByAssetSnReq 通过sn查询资产请求
type AssetCheckByAssetSnReq struct {
	SN            string               `json:"sn" query:"sn" param:"sn"` // 资产编号
	OpratorID     uint64               `json:"opratorId"`                // 操作人ID
	OpratorType   AssetOperateRoleType `json:"opratorType"`              // 操作人类型 1:资产后台(仓库) 2:门店 3:代理
	LocationsType AssetLocationsType   `json:"locationsType"`            // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64               `json:"locationsId"`              // 位置ID
}

// AssetCheckByAssetSnRes 通过sn查询资产返回
type AssetCheckByAssetSnRes struct {
	AssetID   uint64    `json:"assetId"`   // 资产ID
	AssetSN   string    `json:"assetSN"`   // 资产编号
	Model     string    `json:"model"`     // 型号
	BrandName string    `json:"brandName"` // 品牌
	AssetType AssetType `json:"assetType"` // 资产类型 1:电车 2:智能电池
}

type CheckAssetCheckOwnerReq struct {
	AssetID       uint64               `json:"assetId"`       // 资产ID
	AssetType     AssetType            `json:"assetType"`     // 资产类型
	OpratorType   AssetOperateRoleType `json:"opratorType"`   // 操作人类型
	OpratorID     uint64               `json:"opratorId"`     // 操作人ID
	LocationsType AssetLocationsType   `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64               `json:"locationsId"`   // 位置ID
}

type GetCheckAssetReq struct {
	OpratorType   AssetOperateRoleType `json:"opratorType"`   // 操作人类型
	OpratorID     uint64               `json:"opratorId"`     // 操作人ID
	LocationsType AssetLocationsType   `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64               `json:"locationsId"`   // 位置ID
	AssetType     AssetType            `json:"assetType"`     // 资产类型 1:电车 2:电池
}

type GetAssetByOperateRole struct {
	OpratorType   AssetOperateRoleType `json:"opratorType"`   // 操作人类型
	OpratorID     uint64               `json:"opratorId"`     // 操作人ID
	LocationsType AssetLocationsType   `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64               `json:"locationsId"`   // 位置ID
}
