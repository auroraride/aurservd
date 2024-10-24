// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
	"github.com/auroraride/adapter"
	jsoniter "github.com/json-iterator/go"
)

const (
	RiderCabinetOperateReasonEmpty = "开空电仓"
	RiderCabinetOperateReasonFull  = "开满电仓"
)

// RiderCabinetOperateInfoReq 骑手请求换电信息
type RiderCabinetOperateInfoReq struct {
	Serial string `json:"serial" validate:"required" param:"serial" trans:"机器码"`
}

// RiderCabinetOperateProcess 换电流程获取仓门属性
type RiderCabinetOperateProcess struct {
	EmptyBin    *BinInfo `json:"emptyBin,omitempty"`    // 空仓位
	FullBin     *BinInfo `json:"fullBin,omitempty"`     // 满电仓位
	Alternative *BinInfo `json:"alternative,omitempty"` // 备选方案
}

// RiderExchangeInfo 待换电信息
type RiderExchangeInfo struct {
	ID             uint64               `json:"id"`             // 电柜ID
	UUID           string               `json:"uuid"`           // 操作ID
	Full           bool                 `json:"full"`           // 是否有满电电池
	Name           string               `json:"name"`           // 电柜名称
	Health         uint8                `json:"health"`         // 电柜健康状态 0离线 1正常 2故障
	Serial         string               `json:"serial"`         // 电柜编码
	Doors          int                  `json:"doors"`          // 总仓位
	BatteryNum     int                  `json:"batteryNum"`     // 总电池数
	BatteryFullNum int                  `json:"batteryFullNum"` // 总满电电池数
	Brand          adapter.CabinetBrand `json:"brand"`          // 电柜型号
	Model          string               `json:"model"`          // 电池型号
	CityID         uint64               `json:"cityId"`         // 城市ID
	*RiderCabinetOperateProcess

	RetryToken string `json:"retryToken"` // 重试令牌，用于重试操作，重试操作时需要携带，不可泄露，10分钟有效期
}

func (c *RiderExchangeInfo) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(c)
}

func (c *RiderExchangeInfo) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, c)
}

// RiderExchangeProcessReq 请求换电
type RiderExchangeProcessReq struct {
	UUID        string `json:"uuid" validate:"required"` // 操作ID
	Alternative bool   `json:"alternative"`              // 是否使用备选方案
}

// RiderExchangeProcessRes 换电操作步骤返回
type RiderExchangeProcessRes struct {
	Step    uint8      `json:"step"`    // 操作步骤 1:开空电仓 2:放旧电池 3:开满电仓 4:取新电池
	Status  TaskStatus `json:"status"`  // 状态 0:未开始 1:处理中 2:成功 3:失败
	Message string     `json:"message"` // 消息
	Stop    bool       `json:"stop"`    // 步骤是否终止，当该参数为true时，表示换电流程已结束
}

func (r *RiderExchangeProcessRes) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(r)
}

func (r *RiderExchangeProcessRes) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, r)
}

// RiderExchangeProcessStatusReq 获取操作状态
type RiderExchangeProcessStatusReq struct {
	UUID string `json:"uuid" query:"uuid" trans:"操作ID"`
}
