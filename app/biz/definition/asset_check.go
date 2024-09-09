package definition

import "github.com/auroraride/aurservd/app/model"

// AssetCheckByAssetSnReq 通过sn查询资产请求
type AssetCheckByAssetSnReq struct {
	SN        string  `json:"sn" query:"sn" param:"sn"`    // 资产编号
	StationID *uint64 `json:"stationID" query:"stationID"` // 站点ID（扫码时选择站点必填）
}

// AssetCheckCreateReq 创建资产盘点请求
type AssetCheckCreateReq struct {
	AssetCheckCreateDetail []model.AssetCheckCreateDetail `json:"details" validate:"required,dive,required"` // 资产盘点请求详情
	StartAt                string                         `json:"startAt" validate:"required"`               // 盘点开始时间
	EndAt                  string                         `json:"endAt" validate:"required"`                 // 盘点结束时间
	StationID              *uint64                        `json:"stationId"`                                 // 盘点站点ID
}

// AssetCheckCreateRes 创建资产盘点返回
type AssetCheckCreateRes struct {
	ID uint64 `json:"id"` // 资产盘点ID
}

// AssetCheckListReq 盘点记录请求
type AssetCheckListReq struct {
	model.PaginationReq
	Keyword     *string `json:"keyword" query:"keyword"`         // 关键字
	StartAt     *string `json:"startAt" query:"startAt"`         // 开始时间
	EndAt       *string `json:"endAt" query:"endAt"`             // 结束时间
	CheckResult *bool   `json:"checkResult" query:"checkResult"` // 盘点结果 true:正常 false:异常
}
