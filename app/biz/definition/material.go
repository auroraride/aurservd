// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-15, by Jorjan

package definition

import "github.com/auroraride/aurservd/app/model"

// MaterialListReq 其他物资列表请求
type MaterialListReq struct {
	model.PaginationReq
	Keyword *string          `json:"keyword" query:"keyword"` // 关键字
	Type    *model.AssetType `json:"type" query:"type"`       // 其他物资类型 4:电柜配件 5:电车配件 6:其它
}

// MaterialDetail 其他物资详情信息
type MaterialDetail struct {
	ID        uint64          `json:"id"`        // ID
	Name      string          `json:"name"`      // 名称
	Type      model.AssetType `json:"type"`      // 类型
	Statement string          `json:"statement"` // 说明
}

// MaterialCreateReq 其他物资创建
type MaterialCreateReq struct {
	Name      string          `json:"name" validate:"required,max=30" trans:"其他物资名称"`
	Type      model.AssetType `json:"type" validate:"required" enums:"4,5,6" trans:"其他物资类型"`
	Statement string          `json:"statement" validate:"max=50"` // 物资说明
}

// MaterialModifyReq 其他物资修改
type MaterialModifyReq struct {
	model.IDParamReq
	MaterialCreateReq
}
