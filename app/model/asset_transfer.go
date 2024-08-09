package model

// AssetTransferStatus 调拨状态
type AssetTransferStatus uint8

const (
	AssetTransferStatusDelivering AssetTransferStatus = iota + 1 //  配送中
	AssetTransferStatusPending                                   //  待入库
	AssetTransferStatusStock                                     //  已入库
	AssetTransferStatusCancel                                    //  已取消
)

func (s AssetTransferStatus) String() string {
	switch s {
	case AssetTransferStatusPending:
		return "待入库"
	case AssetTransferStatusDelivering:
		return "配送中"
	case AssetTransferStatusStock:
		return "已入库"
	case AssetTransferStatusCancel:
		return "已取消"
	default:
		return "未知"
	}
}

func (s AssetTransferStatus) Value() uint8 {
	return uint8(s)
}

// AssetTransferType 调拨类型
type AssetTransferType uint8

const (
	AssetTransferTypeInitial     AssetTransferType = iota + 1 // 初始入库
	AssetTransferTypeTransfer                                 // 调拨
	AssetTransferTypeActive                                   // 激活
	AssetTransferTypePause                                    // 寄存
	AssetTransferTypeContinue                                 // 取消寄存
	AssetTransferTypeUnSubscribe                              // 退租
	AssetTransferTypeExchange                                 // 换电
)

func (s AssetTransferType) String() string {
	switch s {
	case AssetTransferTypeInitial:
		return "初始入库"
	case AssetTransferTypeTransfer:
		return "调拨"
	case AssetTransferTypeActive:
		return "激活"
	case AssetTransferTypePause:
		return "寄存"
	case AssetTransferTypeContinue:
		return "取消寄存"
	case AssetTransferTypeUnSubscribe:
		return "退租"
	case AssetTransferTypeExchange:
		return "换电"
	default:
		return "未知"
	}
}

func (s AssetTransferType) Value() uint8 {
	return uint8(s)
}

// AssetTransferCreateReq 资产调拨请求
type AssetTransferCreateReq struct {
	FromLocationType  *AssetLocationsType         `json:"fromLocationType"`                          // 调拨前位置类型  1:仓库 2:门店 3:站点 4:运维 (初始调拨此字段不填写)
	FromLocationID    *uint64                     `json:"fromLocationID"`                            // 调拨前位置ID (初始调拨此字段不填写)
	ToLocationType    AssetLocationsType          `json:"toLocationType" validate:"required"`        // 调拨后位置类型  1:仓库 2:门店 3:站点 4:运维
	ToLocationID      uint64                      `json:"toLocationID" validate:"required"`          // 调拨后位置ID
	Details           []AssetTransferCreateDetail `json:"details" validate:"required,dive,required"` // 资产调拨详情
	Reason            string                      `json:"reason" validate:"required"`                // 调拨事由
	AssetTransferType *AssetTransferType          `json:"assetTransferType" enums:"1,2,3,4,5,6"`     // 调拨类型 1:初始入库 2:调拨 3:激活 4:寄存 5:取消寄存 6:退租
}

// AssetTransferCreateDetail 资产调拨详情
type AssetTransferCreateDetail struct {
	AssetType  AssetType `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	SN         *string   `json:"sn"`                            // 资产编号
	Name       *string   `json:"name"`                          // 资产名称
	Num        *uint     `json:"num"`                           // 调拨数量
	MaterialID *uint64   `json:"materialId"`                    // 其它物资分类ID
}

// AssetTransferListReq 资产调拨列表请求
type AssetTransferListReq struct {
	PaginationReq
	AssetTransferFilter
}

// AssetTransferFilter 资产调拨筛选条件
type AssetTransferFilter struct {
	FromLocationType *AssetLocationsType  `json:"fromLocationType" query:"fromLocationType" enums:"1,2,3,4"` // 调拨前位置类型  1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	FromLocationID   *uint64              `json:"fromLocationID" query:"fromLocationID"`                     // 调拨前位置ID
	ToLocationType   *AssetLocationsType  `json:"toLocationType" query:"toLocationType" enums:"1,2,3,4"`     // 调拨后位置类型  1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	ToLocationID     *uint64              `json:"toLocationID" query:"toLocationID"`                         // 调拨后位置ID
	Status           *AssetTransferStatus `json:"status" query:"status" enums:"1,2,3,4"`                     // 调拨状态 1:配送中 2:待入库 3:已入库 4:已取消
	OutStart         *string              `json:"outStart" query:"outStart"`                                 // 出库开始时间
	OutEnd           *string              `json:"outEnd" query:"outEnd"`                                     // 出库结束时间
	InStart          *string              `json:"inStart" query:"inStart"`                                   // 入库开始时间
	InEnd            *string              `json:"inEnd" query:"inEnd"`                                       // 入库结束时间
	Keyword          *string              `json:"keyword" query:"keyword"`                                   // 关键字 (调拨单号，调拨事由、出库人、接收人)
}

// AssetTransferListRes 资产调拨列表响应
type AssetTransferListRes struct {
	ID               uint64 `json:"id"`               // 调拨ID
	SN               string `json:"sn"`               // 调拨单号
	Reason           string `json:"reason"`           // 调拨事由
	FromLocationName string `json:"fromLocationName"` // 调出目标名称
	FromLocationType uint8  `json:"fromLocationType"` // 调出目标类型
	FromLocationID   uint64 `json:"fromLocationID"`   // 调出目标ID
	ToLocationType   uint8  `json:"toLocationType"`   // 调入目标类型
	ToLocationName   string `json:"toLocationName"`   // 调入目标名称
	ToLocationID     uint64 `json:"toLocationID"`     // 调入目标ID
	OutOperateName   string `json:"outOperateName"`   // 出库操作人
	OutNum           uint   `json:"outNum"`           // 出库数量
	InNum            uint   `json:"inNum"`            // 入库数量
	OutTimeAt        string `json:"outTimeAt"`        // 出库时间
	Status           string `json:"status"`           // 调拨状态
	Remark           string `json:"remark"`           // 备注
}

// AssetTransferDetailReq 资产调拨详情请求
type AssetTransferDetailReq struct {
	ID uint64 `json:"id" param:"id" validate:"required"` // 调拨ID
}

// AssetTransferDetail 资产调拨详情
type AssetTransferDetail struct {
	ID        uint64    `json:"id"`        // 资产ID
	AssetType AssetType `json:"assetType"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	SN        string    `json:"sn"`        // 资产编号
	Name      string    `json:"name"`      // 资产名称
	OutNum    uint      `json:"out"`       // 出库数量
	InNum     uint      `json:"in"`        // 入库数量
}

// AssetTransferReceiveReq 接收资产调拨
type AssetTransferReceiveReq struct {
	ID     uint64                       `json:"id" param:"id" validate:"required"`        // 调拨ID
	Detail []AssetTransferReceiveDetail `json:"detail" validate:"required,dive,required"` // 资产明细
	Remark *string                      `json:"remark"`                                   // 备注
}

// AssetTransferReceiveBatchReq 批量接收资产
type AssetTransferReceiveBatchReq struct {
	OperateType          AssetOperateRoleType      `json:"operateType" validate:"required" enums:"1,2,3,4,5,6"` // 操作类型 1:资产后台 2:门店 3:代理 4:运维 5:电柜 6:骑手
	AssetTransferReceive []AssetTransferReceiveReq `json:"assetTransferReceive" validate:"required,dive,required"`
}

// AssetTransferReceiveDetail 接收资产明细
type AssetTransferReceiveDetail struct {
	AssetType  AssetType `json:"assetType" validate:"required"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	SN         *string   `json:"sn"`                            // 资产编号
	Num        *uint     `json:"num"`                           // 调拨数量
	MaterialID *uint64   `json:"materialId"`                    // 其它物资分类ID
}

// GetTransferBySNReq 根据调拨单号获取调拨请求
type GetTransferBySNReq struct {
	SN string `json:"sn" param:"sn" validate:"required"`
}

// AssetTransferFlowReq 资产流转明细请求
type AssetTransferFlowReq struct {
	SN        string     `json:"sn" param:"sn" validate:"required" query:"sn"`    // 资产编号
	Start     *string    `json:"start" query:"start"`                             // 开始时间
	End       *string    `json:"end" query:"end"`                                 // 结束时间
	AssetType *AssetType `json:"assetType" query:"assetType" enums:"1,2,3,4,5,6"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
}

// AssetTransferFlow 资产流转明细
type AssetTransferFlow struct {
	Out AssetTransferFlowDetail `json:"out"` // 出库
	In  AssetTransferFlowDetail `json:"in"`  // 入库
}

// AssetTransferFlowDetail 资产流转明细详情
type AssetTransferFlowDetail struct {
	LocationsType    string `json:"locationsType"`    // 位置类型 1:仓库 2:门店 3:站点 4:运维 5:电柜 6:骑手
	LocationsName    string `json:"locationsName"`    // 位置名称
	TransferTypeName string `json:"transferTypeName"` // 调拨类型 1:初始入库 2:调拨 3:激活 4:寄存 5:取消寄存 6:退租
	TimeAt           string `json:"timeAt"`           // 时间
	OperatorName     string `json:"operatorName"`     // 操作人
}
