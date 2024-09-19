package model

type AssetMaintenanceStatus uint8

const (
	AssetMaintenanceStatusUnder   AssetMaintenanceStatus = iota + 1 // 维修中
	AssetMaintenanceStatusSuccess                                   // 已维修
	AssetMaintenanceStatusFail                                      // 维修失败
	AssetMaintenanceStatusCancel                                    // 已取消
	AssetMaintenanceStatusPause                                     // 已暂停
)

func (a AssetMaintenanceStatus) String() string {
	switch a {
	case AssetMaintenanceStatusUnder:
		return "维修中"
	case AssetMaintenanceStatusSuccess:
		return "已维修"
	case AssetMaintenanceStatusFail:
		return "维修失败"
	case AssetMaintenanceStatusCancel:
		return "已取消"
	case AssetMaintenanceStatusPause:
		return "已暂停"
	default:
		return "未知"
	}
}

func (a AssetMaintenanceStatus) Value() uint8 {
	return uint8(a)
}

// AssetMaintenanceListReq 维修记录列表请求
type AssetMaintenanceListReq struct {
	PaginationReq
	Keyword        *string `json:"keyword"  query:"keyword"`
	Status         *uint8  `json:"status" query:"status"`                 // 状态 2:已维修 3:维修失败
	IsUseAccessory *bool   `json:"isUseAccessory" query:"isUseAccessory"` // 是否使用配件 true:已使用  false:未使用
	Start          *string `json:"start" query:"start"`                   // 开始时间
	End            *string `json:"end" query:"end"`                       // 结束时间
}

// AssetMaintenanceListRes 维修记录列表响应
type AssetMaintenanceListRes struct {
	ID            uint64                   `json:"id"`            // ID
	CabinetName   string                   `json:"cabinetName"`   // 电柜名称
	CabinetSn     string                   `json:"cabinetSn"`     // 电柜编号
	Reason        string                   `json:"reason"`        // 维护理由
	Content       string                   `json:"content"`       // 维护内容
	OperatorName  string                   `json:"operatorName"`  // 维护人
	CreatedAt     string                   `json:"createdAt"`     // 维护时间
	OperatorPhone string                   `json:"operatorPhone"` // 维护人电话
	Status        string                   `json:"status"`        // 维修状态  1:维修中 2:已维修 3:维修失败 4:已取消 5:暂停维护
	Details       []AssetMaintenanceDetail `json:"details"`       // 维护详情

}

// AssetMaintenanceDetail 维修详情
type AssetMaintenanceDetail struct {
	AssetName string `json:"assetName" validate:"required"` // 资产名称
	AssetType string `json:"assetType" validate:"required"` // 资产类型
	Num       uint8  `json:"num" validate:"required"`       // 数量
}

// AssetMaintenanceCreateReq 创建维修记录请求
type AssetMaintenanceCreateReq struct {
	CabinetID       uint64 `json:"cabinetId" validate:"required"` // 电柜ID
	OperatorID      uint64 `json:"operatorId"`                    // 维护人ID
	OperateRoleType uint8  `json:"operateRoleType"`               // 维护人角色类型
}

type AssetMaintenanceCreateDetail struct {
	MaterialID uint64 `json:"materialId" validate:"required"` // 其它物资分类ID
	Num        uint8  `json:"num" validate:"required"`        // 数量
}

// AssetMaintenanceModifyReq 修改维修记录请求
type AssetMaintenanceModifyReq struct {
	ID      uint64                         `json:"id" validate:"required" param:"id"` // ID
	Reason  string                         `json:"reason" validate:"required"`        // 维护理由
	Content string                         `json:"content" validate:"required"`       // 维护内容
	Status  AssetMaintenanceStatus         `json:"status" validate:"required"`        // 维修状态1:维修中 2:已维修 3:维修失败 4:已取消 5:暂停维护
	Details []AssetMaintenanceCreateDetail `json:"details"`                           // 维护详情
}

// AssetMaintenanceRes 维修记录信息
type AssetMaintenanceRes struct {
	ID     uint64                 `json:"id"`     // 维保ID
	Status AssetMaintenanceStatus `json:"status"` // 维修状态 1:维修中 2:已维修 3:维修失败 4:已取消 5:暂停维护
}
