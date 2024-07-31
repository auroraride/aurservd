package model

// AssetTransferStatus 调拨状态
type AssetTransferStatus uint8

const (
	AssetTransferStatusDelivering AssetTransferStatus = iota //  配送中
	AssetTransferStatusStock                                 //  已入库
	AssetTransferStatusCancel                                //  已取消
)

func (s AssetTransferStatus) String() string {
	switch s {
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
	AssetTransferTypeInitial   AssetTransferType = iota + 1 // 初始入库
	AssetTransferTypePlatform                               // 平台调拨
	AssetTransferTypeStore                                  // 门店调拨
	AssetTransferTypeAgent                                  // 代理调拨
	AssetTransferTypeOperation                              // 运维调拨
	AssetTransferTypeSystem                                 // 系统业务自动调拨
)

func (s AssetTransferType) String() string {
	switch s {
	case AssetTransferTypeInitial:
		return "初始入库"
	case AssetTransferTypePlatform:
		return "平台调拨"
	case AssetTransferTypeStore:
		return "门店调拨"
	case AssetTransferTypeAgent:
		return "代理调拨"
	case AssetTransferTypeOperation:
		return "运维调拨"
	case AssetTransferTypeSystem:
		return "系统业务自动调拨"
	default:
		return "未知"
	}
}

func (s AssetTransferType) Value() uint8 {
	return uint8(s)
}

// AssetTransferCreateReq 资产调拨请求
type AssetTransferCreateReq struct {
	FromLocationType *AssetLocationsType         `json:"from_location_type" `                       // 调拨前位置类型  1:仓库 2:门店 3:站点 4:运维
	FromLocationID   *uint64                     `json:"from_location_id" validate:"required"`      // 调拨前位置ID
	ToLocationType   AssetLocationsType          `json:"to_location_type" validate:"required"`      // 调拨后位置类型  1:仓库 2:门店 3:站点 4:运维
	ToLocationID     uint64                      `json:"to_location_id" validate:"required"`        // 调拨后位置ID
	Details          []AssetTransferCreateDetail `json:"details" validate:"required,dive,required"` // 资产调拨详情
}

// AssetTransferCreateDetail 资产调拨详情
type AssetTransferCreateDetail struct {
	AssetType  AssetType `json:"asset_type" validate:"required"` // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	SN         *string   `json:"sn"`                             // 资产编号
	Name       *string   `json:"name"`                           // 资产名称
	Num        *uint     `json:"num"`                            // 调拨数量
	MaterialID *uint64   `json:"materialId"`                     // 其它物资分类ID
}
