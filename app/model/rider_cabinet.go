// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// RiderCabinetOperateStep 换电步骤
type RiderCabinetOperateStep uint8

const (
    RiderCabinetOperateStepOpen    RiderCabinetOperateStep = iota + 1 // 第一步, 开仓门
    RiderCabinetOperateStepPutInto                                    // 第二步, 放入旧电池并关闭仓门
    RiderCabinetOperateStepTakeOut                                    // 第三步, 取出电池
)

// RiderCabinetOperateReq 骑手请求换电信息
type RiderCabinetOperateReq struct {
    Serial string `json:"serial" validate:"required" param:"serial" trans:"机器码"`
}

// RiderCabinetOperateProcess 换电流程获取仓门属性
type RiderCabinetOperateProcess struct {
    EmptyBin    *CabinetBinBasicInfo `json:"emptyBin,omitempty"`    // 空仓位
    FullBin     *CabinetBinBasicInfo `json:"fullBin,omitempty"`     // 满电仓位
    Alternative *CabinetBinBasicInfo `json:"alternative,omitempty"` // 备选方案
}
