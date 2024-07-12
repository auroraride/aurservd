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

// ScrapOperateRoleType 操作资产报废角色类型
type ScrapOperateRoleType uint8

const (
	ScrapOperateRoleAdmin     ScrapOperateRoleType = iota + 1 // 管理员
	ScrapOperateRoleStore                                     // 门店管理员
	ScrapOperateRoleOperation                                 // 运维
	ScrapOperateRoleMaterial                                  // 物资管理员
	ScrapOperateRoleAgent                                     // 代理管理员
)

func (s ScrapOperateRoleType) String() string {
	switch s {
	case ScrapOperateRoleAdmin:
		return "管理员"
	case ScrapOperateRoleStore:
		return "门店管理员"
	case ScrapOperateRoleOperation:
		return "运维"
	case ScrapOperateRoleMaterial:
		return "物资管理员"
	case ScrapOperateRoleAgent:
		return "代理管理员"
	default:
		return "未知"
	}
}

func (s ScrapOperateRoleType) Value() uint8 {
	return uint8(s)
}

// AssetScrapReq 报废资产请求
type AssetScrapReq struct {
	AssetType       AssetType       `json:"assetType" validate:"required"`       // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	ID              uint64          `json:"id" param:"id" validate:"required"`   // 资产ID
	ScrapReasonType ScrapReasonType `json:"scrapReasonType" validate:"required"` // 报废原因
	Remark          *string         `json:"remark"`                              // 备注
}

// AssetScrapListReq 资产报废列表请求
type AssetScrapListReq struct {
	PaginationReq
	ScrapFilter
}
type ScrapFilter struct {
	AssetType       *AssetType       `json:"assetType" query:"assetType" enums:"1,2,3,4,5,6"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	SN              *string          `json:"sn" query:"sn"`                                   // 资产编号
	ModelID         *uint64          `json:"modelId" query:"modelId"`                         // 资产型号
	ScrapReasonType *ScrapReasonType `json:"scrapReason" query:"scrapReason" enums:"1,2,3"`   // 报废原因 1:丢失 2:损坏 3:其他
	OperateName     *string          `json:"operateName" query:"operateName"`                 // 操作人
	Start           *string          `json:"start" query:"start"`                             // 开始时间
	End             *string          `json:"end" query:"end"`                                 // 结束时间
}

// AssetScrapListRes 电池报废列表返回
type AssetScrapListRes struct {
	ID          uint64 `json:"id"`              // 资产ID
	SN          string `json:"sn"`              // 资产编号
	Model       string `json:"model,omitempty"` // 资产型号
	ScrapReason string `json:"scrapReason"`     // 报废原因
	Brand       string `json:"brand"`           // 资产品牌
	OperateName string `json:"operateName"`     // 操作人
	Remark      string `json:"remark"`          // 备注
	ScrapAt     string `json:"scrapAt"`         // 报废时间
	CreatedAt   string `json:"createdAt"`       // 创建时间
}

// AssetScrapBatchRestoreReq 批量恢复资产请求
type AssetScrapBatchRestoreReq struct {
	IDs []uint64 `json:"ids" validate:"required" query:"ids"` // 资产ID
}
