package model

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

// ScrapReasonTypeMap 报废原因类型映射
var ScrapReasonTypeMap = map[ScrapReasonType]string{
	ScrapReasonLost:   "丢失",
	ScrapReasonDamage: "损坏",
	ScrapReasonOther:  "其他",
}

func (s ScrapReasonType) Value() uint8 {
	return uint8(s)
}

// AssetScrapReq 报废资产请求
type AssetScrapReq struct {
	ScrapReasonType ScrapReasonType     `json:"scrapReasonType" validate:"required"`       // 报废原因
	Remark          *string             `json:"remark"`                                    // 备注
	Details         []AssetScrapDetails `json:"details" validate:"required,dive,required"` // 报废明细
	WarehouseID     *uint64             `json:"warehouseId"`                               // 仓库ID
}

type AssetScrapDetails struct {
	AssetType  AssetType `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	Sn         *string   `json:"sn"`                            // 资产Sn (当AssetType = 1:电车 2:智能电池  必传)
	Num        *uint     `json:"num"`                           // 报废数量  （当AssetType = 3:非智能电池 4:电柜配件 5:电车配件 6:其它  必传）
	MaterialID *uint64   `json:"materialId"`                    // 其它物资分类ID （当AssetType = 4:电柜配件 5:电车配件 6:其它  必传）
	ModelID    *uint64   `json:"modelId"`                       // 电池型号ID（当AssetType = 3:非智能电池  必传）
}

// AssetScrapListReq 资产报废列表请求
type AssetScrapListReq struct {
	PaginationReq
	ScrapFilter
}

// ScrapFilter 报废筛选条件
type ScrapFilter struct {
	AssetType       *AssetType       `json:"assetType" query:"assetType" enums:"1,2,3"`     // 资产类型 1:电车 2:智能电池 3:其它
	SN              *string          `json:"sn" query:"sn"`                                 // 资产编号
	ModelID         *uint64          `json:"modelId" query:"modelId"`                       // 电池型号
	BrandID         *uint64          `json:"brandID" query:"brandID"`                       // 电车型号
	ScrapReasonType *ScrapReasonType `json:"scrapReason" query:"scrapReason" enums:"1,2,3"` // 报废原因 1:丢失 2:损坏 3:其他
	OperateName     *string          `json:"operateName" query:"operateName"`               // 操作人
	Start           *string          `json:"start" query:"start"`                           // 开始时间
	End             *string          `json:"end" query:"end"`                               // 结束时间
	Attribute       *string          `json:"attribute" query:"attribute"`                   // 属性 id:value
	AssetName       *string          `json:"assetName" query:"assetName"`                   // 资产名称
}

// AssetScrapListRes 资产报废列表返回
type AssetScrapListRes struct {
	ID          uint64                    `json:"id"`              // 报废ID
	ScrapReason string                    `json:"scrapReason"`     // 报废原因
	OperateName string                    `json:"operateName"`     // 操作人
	Remark      string                    `json:"remark"`          // 备注
	ScrapAt     string                    `json:"scrapAt"`         // 报废时间
	AssetID     uint64                    `json:"assetID"`         // 资产ID
	SN          string                    `json:"sn"`              // 资产编号
	Model       string                    `json:"model,omitempty"` // 资产型号
	Brand       string                    `json:"brand"`           // 资产品牌
	InTimeAt    string                    `json:"inTimeAt"`        // 入库时间
	Attribute   map[uint64]AssetAttribute `json:"attribute"`       // 属性
	Name        string                    `json:"name"`            // 资产名称
	AssetType   string                    `json:"assetType"`       // 资产类型
	Num         uint                      `json:"num"`             // 报废数量
}

// AssetScrapDetailRes 报废详情
type AssetScrapDetailRes struct {
}

// AssetScrapBatchRestoreReq 批量恢复资产请求
type AssetScrapBatchRestoreReq struct {
	IDs []uint64 `json:"ids" validate:"required" query:"ids"` // 资产ID
}
