// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "fmt"
    "github.com/golang-module/carbon/v2"
    "time"
)

type TaskStatus uint8

const (
    TaskStatusNotStart   TaskStatus = iota // 未开始
    TaskStatusProcessing                   // 处理中
    TaskStatusSuccess                      // 成功
    TaskStatusFail                         // 失败
)

func (ts TaskStatus) String() string {
    switch ts {
    case TaskStatusNotStart:
        return "未开始"
    case TaskStatusSuccess:
        return "成功"
    case TaskStatusFail:
        return "失败"
    default:
        return "处理中"
    }
}

// IsSuccess 是否成功
func (ts TaskStatus) IsSuccess() bool {
    return ts == TaskStatusSuccess
}

// ExchangeStep 换电步骤
// RiderCabinetOperateStep
type ExchangeStep uint8

const (
    ExchangeStepOpenEmpty ExchangeStep = iota + 1 // 第一步, 开启空电仓
    ExchangeStepPutInto                           // 第二步, 放入旧电池并关闭仓门
    ExchangeStepOpenFull                          // 第三步, 开启满电仓
    ExchangeStepPutOut                            // 第四步, 取出新电池并关闭仓门
)

func (es ExchangeStep) Is(step ExchangeStep) bool {
    return es == step
}

func (es ExchangeStep) EqualInt(n int) bool {
    return es == ExchangeStep(n)
}

func (es ExchangeStep) Int() int {
    return int(es)
}

func (es ExchangeStep) String() string {
    switch es {
    case ExchangeStepOpenEmpty:
        return "第一步, 开启空电仓"
    case ExchangeStepPutInto:
        return "第二步, 放入旧电池并关闭仓门"
    case ExchangeStepOpenFull:
        return "第三步, 开启满电仓"
    case ExchangeStepPutOut:
        return "第四步, 取出新电池并关闭仓门"
    }
    return "未知"
}

// Next 获取下个步骤
func (es ExchangeStep) Next() ExchangeStep {
    return ExchangeStep(uint8(es) + 1)
}

// IsLast 是否最后一步
func (es ExchangeStep) IsLast() bool {
    return es == ExchangeStepPutOut
}

// ExchangeCabinetInfo 任务电柜设备信息
type ExchangeCabinetInfo struct {
    Health         uint8 `json:"health,omitempty"`         // 电柜健康状态 0离线 1正常 2故障
    Doors          int   `json:"doors,omitempty"`          // 总仓位
    BatteryNum     int   `json:"batteryNum,omitempty"`     // 总电池数
    BatteryFullNum int   `json:"batteryFullNum,omitempty"` // 总满电电池数
}

// ExchangeInfo 换电详情
type ExchangeInfo struct {
    Cabinet *ExchangeCabinetInfo `json:"cabinet,omitempty"` // 电柜信息
    Empty   *BinInfo             `json:"empty,omitempty"`   // 空仓位
    Fully   *BinInfo             `json:"fully,omitempty"`   // 满电仓位
    Steps   []*ExchangeStepInfo  `json:"steps,omitempty"`   // 步骤信息
    Message string               `json:"message,omitempty"` // 错误信息
}

type ExchangeStepInfo struct {
    Step   ExchangeStep `json:"step,omitempty"`   // 操作步骤 1:开空电仓 2:放旧电池 3:开满电仓 4:取新电池
    Status TaskStatus   `json:"status,omitempty"` // 状态 1:处理中 2:成功 3:失败
    Time   time.Time    `json:"time,omitempty"`   // 时间
}

func (si *ExchangeStepInfo) String() string {
    t := si.Time
    if t.IsZero() {
        t = time.Now()
    }
    return fmt.Sprintf("{ %s -> %s }: %s", t.Format(carbon.DateTimeLayout), si.Step, si.Status)
}

// IsSuccess 步骤是否成功
func (si *ExchangeStepInfo) IsSuccess() bool {
    return si.Status.IsSuccess()
}
