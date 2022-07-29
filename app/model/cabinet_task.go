// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-28
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "context"
    "github.com/auroraride/aurservd/pkg/cache"
    jsoniter "github.com/json-iterator/go"
    "time"
)

const (
    CabinetTaskCacheKey = "CABINET_TASK"
)

// CabinetJob 电柜任务
type CabinetJob string

const (
    CabinetJobRiderExchange    CabinetJob = "RDR_EXCHANGE"    // 骑手-换电
    CabinetJobRiderActive                 = "RDR_ACTIVE"      // 骑手-激活
    CabinetJobRiderUnSubscribe            = "RDR_UNSUBSCRIBE" // 骑手-退租
    CabinetJobPause                       = "RDR_PAUSE"       // 骑手-寄存
    CabinetJobContinue                    = "RDR_CONTINUE"    // 骑手-取消寄存
    CabinetJobManagerOpen                 = "MGR_OPEN"        // 管理-开门
)

type CabinetTaskSerial struct {
    Serial string `json:"serial"` // 电柜编码
}

// CabinetTask 电柜任务详情
type CabinetTask struct {
    Task       CabinetJob `json:"task"`       // 任务类别
    Start      time.Time  `json:"start"`      // 开始时间
    Expiration time.Time  `json:"expiration"` // 过期时间
    Finish     bool       `json:"finish"`     // 是否已结束

    Exchange *CabinetTaskExchange `json:"exchange"` // 换电信息
}

func (c *CabinetTask) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(c)
}

func (c *CabinetTask) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, c)
}

func (c *CabinetTask) Save(serial string, maxTime int) {
    cache.HSet(context.Background(), CabinetTaskCacheKey, serial, c)
}

// CabinetTaskDevice 任务电柜设备信息
type CabinetTaskDevice struct {
    Health         uint8 `json:"health"`         // 电柜健康状态 0离线 1正常 2故障
    Doors          uint  `json:"doors"`          // 总仓位
    BatteryNum     uint  `json:"batteryNum"`     // 总电池数
    BatteryFullNum uint  `json:"batteryFullNum"` // 总满电电池数
}

// CabinetTaskExchange 换电信息
type CabinetTaskExchange struct {
    EmptyBin    CabinetBinBasicInfo     `json:"emptyBin"`    // 空仓位
    FullBin     CabinetBinBasicInfo     `json:"fullBin"`     // 满电仓位
    Alternative bool                    `json:"alternative"` // 是否备选方案
    Success     bool                    `json:"success"`     // 是否成功
    Step        RiderCabinetOperateStep `json:"step"`        // 当前步骤
}

// CabinetTaskExchangeBin 任务电柜仓位信息
type CabinetTaskExchangeBin struct {
    Index       int                `json:"index"`          // 仓位index
    Electricity BatteryElectricity `json:"electricity"`    // 电量
    Time        *time.Time         `json:"time,omitempty"` // 操作时间
}
