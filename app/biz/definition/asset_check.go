package definition

import "github.com/auroraride/aurservd/app/model"

// AssetCheckByAssetSnReq 通过sn查询资产请求
type AssetCheckByAssetSnReq struct {
	SN string `json:"sn" query:"sn" param:"sn"` // 资产编号
}

// AssetCheckCreateReq 创建资产盘点请求
type AssetCheckCreateReq struct {
	AssetCheckCreateDetail []model.AssetCheckCreateDetail `json:"details" validate:"required,dive,required"` // 资产盘点请求详情
	StartAt                string                         `json:"startAt" validate:"required"`               // 盘点开始时间
	EndAt                  string                         `json:"endAt" validate:"required"`                 // 盘点结束时间
}

// AssetCheckCreateRes 创建资产盘点返回
type AssetCheckCreateRes struct {
	ID uint64 `json:"id"` // 资产盘点ID
}
