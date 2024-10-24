// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

// EnterpriseAssetListReq 团签资产列表请求
type EnterpriseAssetListReq struct {
	model.PaginationReq
	EnterpriseAssetListFilter
}

type EnterpriseAssetListFilter struct {
	CityID       *uint64 `json:"cityId" query:"cityId"`             // 城市ID
	StationID    *uint64 `json:"stationId" query:"stationId"`       // 站点ID
	ModelID      *uint64 `json:"modelID" query:"modelID"`           // 电池型号ID
	BrandId      *uint64 `json:"brandId" query:"brandId"`           // 电车型号ID
	OtherName    *string `json:"otherName" query:"otherName"`       // 其他物资名称
	Start        *string `json:"start" query:"start"`               // 开始时间
	End          *string `json:"end" query:"end"`                   // 结束时间
	EnterpriseID *uint64 `json:"enterpriseID" query:"enterpriseID"` // 团签ID
}

// EnterpriseAssetDetail 团签资产信息
type EnterpriseAssetDetail struct {
	ID         uint64           `json:"id"`         // 团签ID
	Name       string           `json:"name"`       // 团签名称
	Enterprise EnterpriseDetail `json:"enterprise"` // 团签企业
	City       model.City       `json:"city"`       // 城市
	Total      CommonAssetTotal `json:"total"`      // 资产统计
}

type EnterpriseDetail struct {
	ID   uint64 `json:"id"`   // 企业ID
	Name string `json:"name"` // 企业名称
}
