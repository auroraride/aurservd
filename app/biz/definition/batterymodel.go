// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-12, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

type BatteryModelType uint8

const (
	BatteryModelTypeIntelligent    BatteryModelType = iota + 1 // 智能电池
	BatteryModelTypeNonIntelligent                             // 非智能电池
)

func (t BatteryModelType) Value() uint8 {
	return uint8(t)
}

type BatteryModelListReq struct {
	Type *BatteryModelType `json:"type" query:"type"` // 电池类型
}

// BatteryModelDetail 电池型号信息
type BatteryModelDetail struct {
	ID       uint64           `json:"id"`       // 电池型号ID
	Model    string           `json:"model"`    // 电池型号
	Type     BatteryModelType `json:"type"`     // 电池类型
	Voltage  uint             `json:"voltage"`  // 电压
	Capacity uint             `json:"capacity"` // 容量
}

// BatteryModelCreateReq 创建
type BatteryModelCreateReq struct {
	Type     BatteryModelType `json:"type" validate:"required" enums:"1,2" trans:"电池类型"` // 电池类型
	Voltage  uint             `json:"voltage" validate:"required" trans:"电池电压"`          // 电池电压
	Capacity uint             `json:"capacity" validate:"required" trans:"电池容量"`         // 电池容量
}

// BatteryModelModifyReq 修改
type BatteryModelModifyReq struct {
	model.IDParamReq
	BatteryModelCreateReq
}
