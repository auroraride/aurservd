// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "context"
    "github.com/auroraride/aurservd/pkg/cache"
    jsoniter "github.com/json-iterator/go"
    "go.mongodb.org/mongo-driver/bson/primitive"
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
    Model          string       `json:"model"`          // 电池型号
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
    UUID        primitive.ObjectID `json:"uuid" validate:"required"` // 操作ID
    Alternative bool               `json:"alternative"`              // 是否使用备选方案
}

type RiderCabinetOperateStatus uint8

const (
    RiderCabinetOperateStatusProcessing RiderCabinetOperateStatus = iota + 1 // 处理中
    RiderCabinetOperateStatusSuccess                                         // 成功
    RiderCabinetOperateStatusFail                                            // 失败
)

func (s RiderCabinetOperateStatus) String() string {
    switch s {
    case RiderCabinetOperateStatusSuccess:
        return "成功"
    case RiderCabinetOperateStatusFail:
        return "失败"
    default:
        return "处理中"
    }
}

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

// CabinetProcessJob 获取电柜当前换电任务信息
func CabinetProcessJob(serial string) (*CabinetExchangeProcess, bool) {
    info := new(CabinetExchangeProcess)
    err := cache.Get(context.Background(), serial).Scan(info)
    if err != nil {
        return nil, false
    }
    return info, info != nil && info.Step > 0 && info.Step <= 4
}

// CabinetBusying 查询电柜是否正在业务中
func CabinetBusying(serial string) bool {
    _, busy := CabinetProcessJob(serial)
    return busy
}
