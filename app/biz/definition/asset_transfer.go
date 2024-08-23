package definition

import "github.com/auroraride/aurservd/app/model"

// AssetTransferCreateReq 资产调拨请求
type AssetTransferCreateReq struct {
	FromLocationID *uint64                           `json:"fromLocationID"`                     // 调出仓库/门店ID
	ToLocationType model.AssetLocationsType          `json:"toLocationType" validate:"required"` // 调拨后位置类型  1:仓库 2:门店 3:站点 4:运维
	ToLocationID   uint64                            `json:"toLocationID" validate:"required"`   // 调拨后位置ID
	Details        []model.AssetTransferCreateDetail `json:"details"`                            // 资产调拨详情
	Reason         string                            `json:"reason" validate:"required"`         // 调拨事由
}

// AssetTransferDetailListReq 资产出入库明细请求
type AssetTransferDetailListReq struct {
	model.PaginationReq
	AssetTransferType *model.AssetTransferType `json:"assetTransferType" query:"assetTransferType" enums:"1,2,3,4,5,6"` // 调拨类型 1:初始入库 2:调拨 3:激活 4:寄存 5:取消寄存 6:退租
	Start             *string                  `json:"start" query:"start"`                                             // 开始时间
	End               *string                  `json:"end" query:"end"`                                                 // 结束时间
	AssetType         *model.AssetType         `json:"assetType" query:"assetType" enums:"1,2,3,4,5,6"`                 // 资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它
	Keyword           *string                  `json:"keyword" query:"keyword"`                                         // 关键字
}
