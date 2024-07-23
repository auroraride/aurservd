// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-15, by Jorjan

package definition

import "github.com/auroraride/aurservd/app/model"

type MaterialType uint8

const (
	MaterialTypeBattery MaterialType = iota + 1 // 电池配件物资
	MaterialTypeEbike                           // 电车配件物资
	MaterialTypeOther                           // 其他物资
)

func (t MaterialType) Value() uint8 {
	return uint8(t)
}

type MaterialListReq struct {
	Keyword *string       `json:"keyword" query:"keyword"` // 关键字
	Type    *MaterialType `json:"type" query:"type"`       // 物资类型
}

// MaterialDetail 物资信息
type MaterialDetail struct {
	ID        uint64       `json:"id"`        // ID
	Name      string       `json:"name"`      // 名称
	Type      MaterialType `json:"type"`      // 类型
	Statement string       `json:"statement"` // 说明
	Allot     bool         `json:"allot"`     // 是否可调拨
}

// MaterialCreateReq 创建
type MaterialCreateReq struct {
	Name      string       `json:"name" validate:"required,max=30" trans:"物资名称"`        // 物资名称
	Type      MaterialType `json:"type" validate:"required" enums:"1,2,3" trans:"物资类型"` // 物资类型
	Statement string       `json:"statement" validate:"max=50"`                         // 物资说明
	Allot     bool         `json:"allot"`                                               // 备注
}

// MaterialModifyReq 修改
type MaterialModifyReq struct {
	model.IDParamReq
	MaterialCreateReq
}
