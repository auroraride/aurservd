// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    jsoniter "github.com/json-iterator/go"
)

// RiderCabinetOperateStep 换电步骤
type RiderCabinetOperateStep uint8

const (
    RiderCabinetOperateStepOpenEmpty RiderCabinetOperateStep = iota + 1 // 第一步, 开启空电仓
    RiderCabinetOperateStepPutInto                                      // 第二步, 放入旧电池并关闭仓门
    RiderCabinetOperateStepOpenFull                                     // 第三步, 开启满电仓
    RiderCabinetOperateStepClose                                        // 第四步, 取出新电池并关闭仓门
)

func (ros RiderCabinetOperateStep) String() string {
    switch ros {
    case RiderCabinetOperateStepOpenEmpty:
        return "第一步, 开启空电仓"
    case RiderCabinetOperateStepPutInto:
        return "第二步, 放入旧电池并关闭仓门"
    case RiderCabinetOperateStepOpenFull:
        return "第三步, 开启满电仓"
    case RiderCabinetOperateStepClose:
        return "第四步, 取出新电池并关闭仓门"
    }
    return "未知"
}

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
    EmptyBin    *CabinetBinBasicInfo `json:"emptyBin,omitempty"`    // 空仓位
    FullBin     *CabinetBinBasicInfo `json:"fullBin,omitempty"`     // 满电仓位
    Alternative *CabinetBinBasicInfo `json:"alternative,omitempty"` // 备选方案
}

// RiderCabinetInfo 待换电信息
type RiderCabinetInfo struct {
    ID             uint64       `json:"id"`             // 电柜ID
    UUID           string       `json:"uuid"`           // 操作ID
    Full           bool         `json:"full"`           // 是否有满电电池
    Name           string       `json:"name"`           // 电柜名称
    Health         uint8        `json:"health"`         // 电柜健康状态 0离线 1正常 2故障
    Serial         string       `json:"serial"`         // 电柜编码
    Doors          uint         `json:"doors"`          // 总仓位
    BatteryNum     uint         `json:"batteryNum"`     // 总电池数
    BatteryFullNum uint         `json:"batteryFullNum"` // 总满电电池数
    Brand          CabinetBrand `json:"brand"`          // 电柜型号
    Voltage        float64      `json:"voltage"`        // 电池电压
    CityID         uint64       `json:"cityId"`         // 城市ID
    RiderCabinetOperateProcess
}

func (c *RiderCabinetInfo) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(c)
}

func (c *RiderCabinetInfo) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, c)
}

// RiderCabinetOperateReq 请求换电
type RiderCabinetOperateReq struct {
    UUID        *string `json:"uuid" validate:"required"` // 操作ID
    Alternative *bool   `json:"alternative"`              // 是否使用备选方案
}

// RiderCabinetOperating 电柜处理
type RiderCabinetOperating struct {
    UUID        string             `json:"uuid"`
    Serial      string             `json:"serial"`
    ID          uint64             `json:"id"`          // 电柜ID
    EmptyIndex  int                `json:"emptyIndex"`  // 空店仓
    FullIndex   int                `json:"fullIndex"`   // 满电仓
    Electricity BatteryElectricity `json:"electricity"` // 满电电池电量
}

type RiderCabinetOperateStatus uint8

const (
    RiderCabinetOperateStatusProcessing RiderCabinetOperateStatus = iota + 1 // 处理中
    RiderCabinetOperateStatusSuccess                                         // 成功
    RiderCabinetOperateStatusFail                                            // 失败
)

// RiderCabinetOperateRes 换电操作步骤返回
type RiderCabinetOperateRes struct {
    Step    RiderCabinetOperateStep   `json:"step"`    // 操作步骤 1:开空电仓 2:放旧电池 3:开满电仓 4:取新电池
    Status  RiderCabinetOperateStatus `json:"status"`  // 状态 1:处理中 2:成功 3:失败
    Message string                    `json:"message"` // 消息
    Stop    bool                      `json:"stop"`    // 步骤是否终止
}

func (c *RiderCabinetOperateRes) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(c)
}

func (c *RiderCabinetOperateRes) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, c)
}

// RiderCabinetOperateStatusReq 获取操作状态
type RiderCabinetOperateStatusReq struct {
    UUID *string `json:"uuid" query:"uuid" trans:"操作ID"`
}
