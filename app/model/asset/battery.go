package asset

import (
	"github.com/auroraride/aurservd/app/model"
)

const (
	BatteryBrandXC  = "XC" // 星创电池
	BatteryBrandTB  = "TB" // 拓邦电池
	BatteryXcLength = 16   // 星创电池长度
	BatteryTbLength = 24   // 拓邦电池长度
)

// BatteryAssetStatus 电池资产状态
type BatteryAssetStatus uint8

const (
	BatteryAssetStatusPending    BatteryAssetStatus = iota // 待入库
	BatteryAssetStatusStock                                // 库存中
	BatteryAssetStatusDelivering                           // 配送中
	BatteryAssetStatusUsing                                // 使用中
	BatteryAssetStatusFault                                // 故障
	BatteryAssetStatusScrap                                // 报废
)

func (s BatteryAssetStatus) String() string {
	switch s {
	case BatteryAssetStatusPending:
		return "待入库"
	case BatteryAssetStatusStock:
		return "库存中"
	case BatteryAssetStatusDelivering:
		return "配送中"
	case BatteryAssetStatusUsing:
		return "使用中"
	case BatteryAssetStatusFault:
		return "故障"
	case BatteryAssetStatusScrap:
		return "报废"
	default:
		return "未知"
	}
}

func (s BatteryAssetStatus) Value() uint8 {
	return uint8(s)
}

type ScrapReasonType uint8

const (
	ScrapReasonLost   ScrapReasonType = iota + 1 // 丢失
	ScrapReasonDamage                            // 损坏
	ScrapReasonOther                             // 其他
)

func (s ScrapReasonType) String() string {
	switch s {
	case ScrapReasonLost:
		return "丢失"
	case ScrapReasonDamage:
		return "损坏"
	case ScrapReasonOther:
		return "其他"
	default:
		return "未知"
	}
}

func (s ScrapReasonType) Value() uint8 {
	return uint8(s)
}

// BatteryAssetLocationType 电池资产位置类型
type BatteryAssetLocationType uint8

const (
	BatteryAssetLocationTypeWarehouse BatteryAssetLocationType = iota + 1 // 仓库
	BatteryAssetLocationTypeStore                                         // 门店
	BatteryAssetLocationTypeCabinet                                       // 电柜
	BatteryAssetLocationTypeStation                                       // 站点
	BatteryAssetLocationTypeRider                                         // 骑手
	BatteryAssetLocationTypeOperation                                     // 运维
)

func (s BatteryAssetLocationType) String() string {
	switch s {
	case BatteryAssetLocationTypeWarehouse:
		return "仓库"
	case BatteryAssetLocationTypeStore:
		return "门店"
	case BatteryAssetLocationTypeCabinet:
		return "电柜"
	case BatteryAssetLocationTypeStation:
		return "站点"
	case BatteryAssetLocationTypeRider:
		return "骑手"
	case BatteryAssetLocationTypeOperation:
		return "运维"
	default:
		return "未知"
	}
}

func (s BatteryAssetLocationType) Value() uint8 {
	return uint8(s)
}

// BatteryCreateReq 创建电池请求
type BatteryCreateReq struct {
	SN     string `json:"sn" validate:"required" trans:"电池编号"`
	CityID uint64 `json:"cityId" validate:"required" trans:"城市"`
}

// BatteryModifyReq 修改电池请求
type BatteryModifyReq struct {
	ID              uint64           `json:"id" param:"id" validate:"required" trans:"电池ID"`
	Enable          *bool            `json:"enable"`      // 是否启用
	CityID          *uint64          `json:"cityId"`      // 城市ID
	ScrapReasonType *ScrapReasonType `json:"scrapReason"` // 报废原因
	Remark          *string          `json:"remark"`      // 备注
}

// BatteryFilter 电池筛选条件
type BatteryFilter struct {
	SN                    *string `json:"sn" query:"sn"`                                                     // 编号
	Model                 *string `json:"model" query:"model"`                                               // 型号
	CityID                *uint64 `json:"cityId" query:"cityId"`                                             // 城市
	OwnerType             *uint8  `json:"ownerType" query:"ownerType" enums:"1,2"`                           // 归属类型   1:平台 2:代理商
	AssetLocationsType    *uint8  `json:"assetLocationsType" query:"assetLocationsType" enums:"1,2,3,4,5,6"` // 资产位置类型 1:仓库 2:门店 3:电柜 4:站点 5:骑手 6:运维
	AssetLocationsKeywork *bool   `json:"assetKeywork" query:"assetKeywork"`                                 // 资产位置关键词
	AssetStatus           *uint8  `json:"assetStatus" query:"assetStatus" enums:"1,2,3,4,5"`                 // 资产状态 1:待入库 2:库存中 3:配送中 4:使用中 5:故障
	Enable                *bool   `json:"enable" query:"enable"`                                             // 是否启用
}

// BatteryListReq 电池列表请求
type BatteryListReq struct {
	model.PaginationReq
	BatteryFilter
}

// BatteryListRes 电池列表返回
type BatteryListRes struct {
	ID             uint64 `json:"id"`             // 电池ID
	CityName       string `json:"cityName"`       // 城市
	Belong         string `json:"belong"`         // 归属
	AssetLocations string `json:"assetLocations"` // 资产位置
	Brand          string `json:"brand"`          // 品牌
	Model          string `json:"model"`          // 电池型号
	SN             string `json:"sn"`             // 编号
	AssetStatus    string `json:"assetStatus"`    // 资产状态
	Enable         bool   `json:"enable"`         // 是否启用
	Remark         string `json:"remark"`
}

type Battery struct {
	ID    uint64 `json:"id"`
	SN    string `json:"sn"`    // 编号
	Model string `json:"model"` // 型号
	Brand string `json:"brand"` // 品牌
}

type BatterySearchReq struct {
	Serial       string  `json:"serial" query:"serial" trans:"流水号" validate:"required,min=4"`
	EnterpriseID *uint64 `json:"enterpriseId" query:"enterpriseId"` // 团签ID: 0为查询非团签电池; 不携带为全部数据
	StationID    *uint64 `json:"stationId" query:"stationId"`       // 站点ID: 0为查询非站点电池; 不携带为全部数据
}

type BatteryBind struct {
	RiderID   uint64 `json:"riderId" validate:"required"`   // 骑手ID
	BatteryID uint64 `json:"batteryId" validate:"required"` // 电池ID
}

// BatteryDetail 电池信息
type BatteryDetail struct {
	ID    uint64  `json:"id"`    // 电池ID
	Model string  `json:"model"` // 电池型号
	SN    string  `json:"sn"`    // 电池编码
	Soc   float64 `json:"soc"`   // 当前电量, 暂时隐藏
}

type BatteryInCabinet struct {
	CabinetID uint64 `json:"cabinetId"` // 所在电柜ID
	Ordinal   int    `json:"ordinal"`   // 仓位序号
}

type BatteryUnbindRequest struct {
	RiderID uint64 `json:"riderId" validate:"required"` // 骑手ID
}

type BatteryBatchQueryRequest struct {
	IDs []uint64 `json:"ids" validate:"required,min=1"`
}

type BatteryQueryRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// BatteryEnterpriseTransfer 代理商电池转移信息
type BatteryEnterpriseTransfer struct {
	Sn           string  `json:"sn"`           // 电池编码
	StationID    *uint64 `json:"stationId"`    // 站点ID
	EnterpriseID *uint64 `json:"enterpriseId"` // 团签ID
}

var (
	BatteryModelXC = map[string]string{
		"08": "72V30AH",
		"11": "72V35AH",
		"12": "60V30AH",
		"16": "72V35AH",
	}
)

// BatteryFlowDetail 电池流转明细
type BatteryFlowDetail struct {
	ID       uint64 `json:"id"`       // 电池ID
	SN       string `json:"sn"`       // 电池编号
	Model    string `json:"model"`    // 电池型号
	Brand    string `json:"brand"`    // 电池品牌
	FlowID   uint64 `json:"flowId"`   // 流水ID
	FlowType uint8  `json:"flowType"` // 流水类型
}

// BatteryScrapListReq 电池报废列表请求
type BatteryScrapListReq struct {
	model.PaginationReq
	ScrapFilter
}
type ScrapFilter struct {
	SN              *string          `json:"sn" query:"sn"`                                 // 电池编号
	Model           *string          `json:"model" query:"model"`                           // 电池型号
	ScrapReasonType *ScrapReasonType `json:"scrapReason" query:"scrapReason" enums:"1,2,3"` // 报废原因 1:丢失 2:损坏 3:其他
	OperateID       *uint64          `json:"operateId" query:"operateId"`                   // 操作人ID
	Start           *string          `json:"start" query:"start"`                           // 开始时间
	End             *string          `json:"end" query:"end"`                               // 结束时间
}

// BatteryScrapListRes 电池报废列表返回
type BatteryScrapListRes struct {
	ID          uint64          `json:"id"`          // 电池ID
	SN          string          `json:"sn"`          // 电池编号
	Model       string          `json:"model"`       // 电池型号
	ScrapReason ScrapReasonType `json:"scrapReason"` // 报废原因
	Brand       string          `json:"brand"`       // 电池品牌
	Operate     string          `json:"operate"`     // 操作人
	Remark      string          `json:"remark"`      // 备注
	ScrapAt     string          `json:"scrapAt"`     // 报废时间
	CreatedAt   string          `json:"createdAt"`   // 创建时间
}
