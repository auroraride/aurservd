// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

// MaintainerAssetListReq 运维资产列表请求
type MaintainerAssetListReq struct {
	model.PaginationReq
	Keyword   *string `json:"keyword" query:"keyword"`     // 关键字 姓名，手机号
	ModelID   *uint64 `json:"modelId" query:"modelId"`     // 电池型号ID
	BrandID   *uint64 `json:"brandId" query:"brandId"`     // 电车型号ID
	OtherName *string `json:"otherName" query:"otherName"` // 其他物资名称
}

// MaintainerAssetDetail 运维资产信息
type MaintainerAssetDetail struct {
	ID    uint64           `json:"id"`    // 运维ID
	Name  string           `json:"name"`  // 运维名称
	Phone string           `json:"phone"` // 运维电话
	Total CommonAssetTotal `json:"total"` // 资产统计
}
