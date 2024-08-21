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
