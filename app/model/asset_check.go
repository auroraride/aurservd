package model

type AssetCheckDetailsStatus uint8

const (
	AssetCheckDetailsStatusUntreated AssetCheckDetailsStatus = iota // 未处理
	AssetCheckDetailsStatusIn                                       // 已入库
	AssetCheckDetailsStatusOut                                      // 已出库
	AssetCheckDetailsStatusScrap                                    // 已报废
)

func (a AssetCheckDetailsStatus) String() string {
	switch a {
	case AssetCheckDetailsStatusUntreated:
		return "未处理"
	case AssetCheckDetailsStatusIn:
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

type AssetCheckStatus uint8

const (
	AssetCheckStatusPending   AssetCheckStatus = iota // 待处理
	AssetCheckStatusProcessed                         // 已处理
)

func (a AssetCheckStatus) String() string {
	switch a {
	case AssetCheckStatusPending:
		return "待处理"
	case AssetCheckStatusProcessed:
		return "已处理"
	default:
		return ""
	}
}

func (a AssetCheckStatus) Value() uint8 {
	return uint8(a)
}

// AssetCheckResult 盘点结果
type AssetCheckResult uint8

const (
	AssetCheckResultUntreated AssetCheckResult = iota // 未盘点
	AssetCheckResultNormal                            // 正常
	AssetCheckResultLoss                              // 亏
	AssetCheckResultSurplus                           // 盈
)

func (a AssetCheckResult) String() string {
	switch a {
	case AssetCheckResultUntreated:
		return "未盘点"
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
	OperatorID             uint64                   `json:"operatorId"`                                // 操作人ID
	OperatorType           OperatorType             `json:"operatorType"`                              // 操作人类型  2:门店 3:代理 6:资产后台(仓库)
	AssetCheckCreateDetail []AssetCheckCreateDetail `json:"details" validate:"required,dive,required"` // 资产盘点请求详情
	StartAt                string                   `json:"startAt" validate:"required"`               // 盘点开始时间
	EndAt                  string                   `json:"endAt" validate:"required"`                 // 盘点结束时间
}

type AssetCheckCreateDetail struct {
	AssetID   uint64    `json:"assetId" validate:"required"`   // 资产ID
	AssetType AssetType `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池
}

// AssetCheckDetailListReq 资产盘点请求详情
type AssetCheckDetailListReq struct {
	PaginationReq
	ID        uint64    `json:"id" validate:"required" param:"id"` // 盘点ID
	AssetType AssetType `json:"assetType" query:"assetType"`       // 资产类型 1:电车 2:智能电池
	RealCheck bool      `json:"realCheck" query:"realCheck"`       // 是否实际盘点 true:实际盘点 false:应盘点
	SN        *string   `json:"sn" query:"sn"`                     // 资产编号
	Attribute *string   `json:"attribute" query:"attribute"`       // 属性 id:value
	BrandID   *uint64   `json:"brandId" query:"brandId"`           // 品牌
	ModelID   *uint64   `json:"modelId" query:"modelId"`           // 型号ID
}

// AssetCheckListReq 获取资产盘点请求
type AssetCheckListReq struct {
	PaginationReq
	LocationsID   *uint64             `json:"locationsId" query:"locationsId"`     // 位置ID
	LocationsType *AssetLocationsType `json:"locationsType" query:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	Keyword       *string             `json:"keyword" query:"keyword"`             // 关键字
	StartAt       *string             `json:"startAt" query:"startAt"`             // 开始时间
	EndAt         *string             `json:"endAt" query:"endAt"`                 // 结束时间
	CheckResult   *bool               `json:"checkResult" query:"checkResult"`     // 盘点结果 true:正常 false:异常
	LocationsIds  []uint64            `json:"locationsIds" query:"locationsIds"`   // 位置IDs
}

// AssetCheckListRes 获取资产盘点返回
type AssetCheckListRes struct {
	ID             uint64                `json:"id"`             // 盘点ID
	StartAt        string                `json:"startAt"`        // 盘点开始时间
	EndAt          string                `json:"endAt"`          // 盘点结束时间
	OperatorID     uint64                `json:"operatorId"`     // 操作人ID
	OperatorName   string                `json:"OperatorName"`   // 操作人名称
	BatteryNum     uint                  `json:"batteryNum"`     // 应盘点电池数量
	BatteryNumReal uint                  `json:"batteryNumReal"` // 实盘电池数量
	EbikeNum       uint                  `json:"ebikeNum"`       // 应盘点电车数量
	EbikeNumReal   uint                  `json:"ebikeNumReal"`   // 实盘电车数量
	LocationsID    uint64                `json:"locationsId"`    // 位置ID
	LocationsType  uint8                 `json:"locationsType"`  // 位置类型
	LocationsName  string                `json:"locationsName"`  // 位置名称
	CheckResult    bool                  `json:"checkResult"`    // 盘点结果 true:正常 false:异常
	Status         string                `json:"status"`         // 状态 1:待处理 2:已处理
	Abnormals      []*AssetCheckAbnormal `json:"abnormals"`      // 盘点异常资产
}

// AssetCheckAbnormal 异常资产
type AssetCheckAbnormal struct {
	ID                uint64           `json:"id"`                // 盘点详情ID
	AssetID           uint64           `json:"assetId"`           // 资产ID
	Name              string           `json:"name"`              // 名称
	Model             string           `json:"model"`             // 型号
	Brand             string           `json:"brand"`             // 品牌
	Result            AssetCheckResult `json:"result"`            // 盘点结果 0:正常 1:亏 2:盈
	SN                string           `json:"sn"`                // 资产编号
	LocationsName     string           `json:"locationsName"`     // 理论位置名称
	RealLocationsName string           `json:"realLocationsName"` // 实际位置名称
	Status            string           `json:"status"`            // 处理状态 0:未处理 1:已入库 2:已出库 3:已报废
	OperatorName      string           `json:"operatorName"`      // 操作人名称
	OperatorAt        string           `json:"operatorAt"`        // 操作时间
	AssetType         AssetType        `json:"assetType"`         // 资产类型 1:电车 2:电池
}

// AssetCheckDetail 资产明细
type AssetCheckDetail struct {
	ID                uint64                    `json:"id"`                // 盘点详情ID
	AssetID           uint64                    `json:"assetId"`           // 资产ID
	AssetSN           string                    `json:"assetSN"`           // 资产编号
	Model             string                    `json:"model"`             // 型号
	BrandName         string                    `json:"brandName"`         // 品牌
	LocationsName     string                    `json:"locationsName"`     // 位置名称
	RealLocationsName string                    `json:"realLocationsName"` // 实际位置名称
	AssetStatus       uint8                     `json:"assetStatus"`       // 资产状态  0:待入库 1:库存中 2:配送中 3:使用中 4:故障 5:报废
	Attribute         map[uint64]AssetAttribute `json:"attribute"`         // 属性
	AssetType         uint8                     `json:"assetType"`         // 资产类型 1:电车 2:电池
}

// AssetCheckByAssetSnReq 通过sn查询资产请求
type AssetCheckByAssetSnReq struct {
	SN            string             `json:"sn" query:"sn" param:"sn"` // 资产编号
	OperatorID    uint64             `json:"operatorId"`               // 操作人ID
	OperatorType  OperatorType       `json:"OperatorType"`             // 操作人类型 1:资产后台(仓库) 2:门店 3:代理
	LocationsType AssetLocationsType `json:"locationsType"`            // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64             `json:"locationsId"`              // 位置ID
}

// AssetCheckByAssetSnRes 通过sn查询资产返回
type AssetCheckByAssetSnRes struct {
	AssetID       uint64                    `json:"assetId"`       // 资产ID
	AssetSN       string                    `json:"assetSN"`       // 资产编号
	Model         string                    `json:"model"`         // 型号
	BrandName     string                    `json:"brandName"`     // 品牌
	AssetType     AssetType                 `json:"assetType"`     // 资产类型 1:电车 2:智能电池
	LocationsType AssetLocationsType        `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	LocationsID   uint64                    `json:"locationsId"`   // 位置ID
	LocationsName string                    `json:"locationsName"` // 位置名称
	Attribute     map[uint64]AssetAttribute `json:"attribute"`     // 属性
}

type CheckAssetCheckOwnerReq struct {
	AssetID       uint64             `json:"assetId"`       // 资产ID
	AssetType     AssetType          `json:"assetType"`     // 资产类型
	OperatorType  OperatorType       `json:"OperatorType"`  // 操作人类型  0:业务管理员 1:门店 2:电柜 3:代理 4:运维 5:骑手 6:资产管理员
	OperatorID    uint64             `json:"operatorId"`    // 操作人ID
	LocationsType AssetLocationsType `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64             `json:"locationsId"`   // 位置ID
}

type GetCheckAssetReq struct {
	OperatorType  OperatorType       `json:"OperatorType"`  // 操作人类型 0:业务管理员 1:门店 2:电柜 3:代理 4:运维 5:骑手 6:资产管理员
	OperatorID    uint64             `json:"operatorId"`    // 操作人ID
	LocationsType AssetLocationsType `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64             `json:"locationsId"`   // 位置ID
	AssetType     AssetType          `json:"assetType"`     // 资产类型 1:电车 2:电池
}

type GetAssetByOperateRole struct {
	OperatorType  OperatorType       `json:"OperatorType"`  // 操作人类型 0:业务管理员 1:门店 2:电柜 3:代理 4:运维 5:骑手 6:资产管理员
	OperatorID    uint64             `json:"operatorId"`    // 操作人ID
	LocationsType AssetLocationsType `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64             `json:"locationsId"`   // 位置ID
}

type MarkStartOrEndCheckReq struct {
	OperatorType  OperatorType       `json:"OperatorType"`  // 操作人类型 0:业务管理员 1:门店 2:电柜 3:代理 4:运维 5:骑手 6:资产管理员
	OperatorID    uint64             `json:"operatorId"`    // 操作人ID
	LocationsType AssetLocationsType `json:"locationsType"` // 位置类型 1:仓库 2:门店 3:站点
	LocationsID   uint64             `json:"locationsId"`   // 位置ID
	Enable        bool               `json:"enable"`        // 是否开始盘点 true:开始 false:结束
}

// AssetCheckListAbnormalReq 获取盘点异常资产请求
type AssetCheckListAbnormalReq struct {
	ID uint64 `json:"id" validate:"required" param:"id"` // 盘点ID
}

// AssetCheckAbnormalOperateReq 异常资产操作
type AssetCheckAbnormalOperateReq struct {
	ID uint64 `json:"id" validate:"required" param:"id"` // 盘点异常id
}
