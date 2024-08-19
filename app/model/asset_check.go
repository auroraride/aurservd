package model

type AssetCheckDetailsStatus uint8

const (
	AssetCheckDetailsStatusUntreated AssetCheckDetailsStatus = iota // 未处理
	AssetCheckDetailsStatusStock                                    // 已入库
	AssetCheckDetailsStatusOut                                      // 已出库
	AssetCheckDetailsStatusScrap                                    // 已报废
)

func (a AssetCheckDetailsStatus) String() string {
	switch a {
	case AssetCheckDetailsStatusUntreated:
		return "未处理"
	case AssetCheckDetailsStatusStock:
		return "已入库"
	case AssetCheckDetailsStatusOut:
		return "已出库"
	case AssetCheckDetailsStatusScrap:
		return "已报废"
	default:
		return ""
	}
}

func (a AssetCheckDetailsStatus) Value() uint8 {
	return uint8(a)
}

// AssetCheckResult 盘点结果
type AssetCheckResult uint8

const (
	AssetCheckResultNormal  AssetCheckResult = iota // 正常
	AssetCheckResultLoss                            // 亏
	AssetCheckResultSurplus                         // 盈
)

func (a AssetCheckResult) String() string {
	switch a {
	case AssetCheckResultNormal:
		return "正常"
	case AssetCheckResultLoss:
		return "亏"
	case AssetCheckResultSurplus:
		return "盈"
	default:
		return ""
	}
}

func (a AssetCheckResult) Value() uint8 {
	return uint8(a)
}

// AssetCheckCreateReq 创建资产盘点请求
type AssetCheckCreateReq struct {
	LocationsType          AssetLocationsType       `json:"locationsType" validate:"required"`         // 位置类型 1:仓库 2:门店 3:站点
	LocationsID            uint64                   `json:"locationsId" validate:"required"`           // 位置ID
	OpratorID              uint64                   `json:"opratorId"`                                 // 操作人ID
	OpratorType            AssetOperateRoleType     `json:"opratorType"`                               // 操作人类型 1:资产后台(仓库) 2:门店 3:代理
	AssetCheckCreateDetail []AssetCheckCreateDetail `json:"details" validate:"required,dive,required"` // 资产盘点请求详情
	StartAt                string                   `json:"startAt" validate:"required"`               // 盘点开始时间
	EndAt                  string                   `json:"endAt" validate:"required"`                 // 盘点结束时间
}

// AssetCheckCreateDetail 资产盘点请求详情
type AssetCheckCreateDetail struct {
	AssetID   uint64    `json:"assetId" validate:"required"`   // 资产ID
	AssetType AssetType `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池
}

// AssetCheckListReq 获取资产盘点请求
type AssetCheckListReq struct {
	PaginationReq
	LocationsID   *uint64             `json:"locationsId" validate:"required"`   // 位置ID
	LocationsType *AssetLocationsType `json:"locationsType" validate:"required"` // 位置类型 1:仓库 2:门店 3:站点
	Keyword       *string             `json:"keyword" query:"keyword"`           // 关键字
	StartAt       *string             `json:"startAt" query:"startAt"`           // 开始时间
	EndAt         *string             `json:"endAt" query:"endAt"`               // 结束时间
}

// AssetCheckListRes 获取资产盘点返回
type AssetCheckListRes struct {
	StartAt        string `json:"startAt"`        // 盘点开始时间
	EndAt          string `json:"endAt"`          // 盘点结束时间
	OpratorID      uint64 `json:"opratorId"`      // 操作人ID
	OpratorName    string `json:"opratorName"`    // 操作人名称
	BatteryNum     uint   `json:"batteryNum"`     // 应盘点电池数量
	BatteryNumReal uint   `json:"batteryNumReal"` // 实盘电池数量
	EbikeNum       uint   `json:"ebikeNum"`       // 应盘点电车数量
	EbikeNumReal   uint   `json:"ebikeNumReal"`   // 实盘电车数量
	LocationsID    uint64 `json:"locationsId"`    // 位置ID
	LocationsType  uint8  `json:"locationsType"`  // 位置类型
	LocationsName  string `json:"locationsName"`  // 位置名称
	CheckResult    bool   `json:"checkResult"`    // 盘点结果 true:正常 false:异常
}

// AssetCheckAbnormal 异常资产
type AssetCheckAbnormal struct {
	AssetID           uint64 `json:"assetId"`           // 资产ID
	Name              string `json:"name"`              // 名称
	Model             string `json:"model"`             // 型号
	Brand             string `json:"brand"`             // 品牌
	Result            string `json:"result"`            // 盘点结果 0:正常 1:亏 2:盈
	SN                string `json:"sn"`                // 资产编号
	LocationsName     string `json:"locationsName"`     // 理论位置名称
	RealLocationsName string `json:"realLocationsName"` // 实际位置名称
	Status            string `json:"status"`            // 处理状态 0:未处理 1:已入库 2:已出库 3:已报废
	OpratorName       string `json:"opratorName"`       // 操作人名称
	OpratorAt         string `json:"opratorAt"`         // 操作时间
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

type MarkStartOrEndCheckReq struct {
	OpratorType   AssetOperateRoleType `json:"opratorType"`   // 操作人类型
	OpratorID     uint64               `json:"opratorId"`     // 操作人ID
	LocationsType AssetLocationsType   `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64               `json:"locationsId"`   // 位置ID
	Enable        bool                 `json:"enable"`        // 是否开始盘点 true:开始 false:结束
}

// AssetCheckListAbnormalReq 获取盘点异常资产请求
type AssetCheckListAbnormalReq struct {
	AssetCheckID uint64 `json:"assetCheckId" validate:"required"` // 盘点ID
}
