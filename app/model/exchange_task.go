// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "fmt"
    "github.com/golang-module/carbon/v2"
    "time"
)

// DoorStatus 柜门状态(处理换电用)
type DoorStatus uint8

const (
    DoorStatusUnknown            DoorStatus = iota // 未知
    DoorStatusClose                                // 关闭
    DoorStatusOpen                                 // 开启
    DoorStatusFail                                 // 故障
    DoorStatusBatteryFull                          // 电池未取出
    DoorStatusBatteryEmpty                         // 电池未放入
    DoorStatusCabinetStatusError                   // 电柜状态获取失败
)

func (bds DoorStatus) String() string {
    switch bds {
    case DoorStatusClose:
        return "关闭"
    case DoorStatusOpen:
        return "开启"
    case DoorStatusFail:
        return "故障"
    case DoorStatusBatteryFull:
        return "电池未取出"
    case DoorStatusBatteryEmpty:
        return "电池未放入"
    case DoorStatusCabinetStatusError:
        return "电柜状态获取失败"
    }
    return "未知"
}

// DoorError 换电故障
var DoorError = map[DoorStatus]string{
    DoorStatusUnknown:      "仓门状态未知",
    DoorStatusFail:         "仓门故障",
    DoorStatusBatteryFull:  "电池未取出",
    DoorStatusBatteryEmpty: "电池未放入",
}

// ExchangeTaskCabinet 任务电柜设备信息
type ExchangeTaskCabinet struct {
    Health         uint8       `json:"health" bson:"health"`                 // 电柜健康状态 0离线 1正常 2故障
    Doors          int         `json:"doors" bson:"doors"`                   // 总仓位
    BatteryNum     int         `json:"batteryNum" bson:"batteryNum"`         // 总电池数
    BatteryFullNum int         `json:"batteryFullNum" bson:"batteryFullNum"` // 总满电电池数
    Bins           CabinetBins `json:"bins" bson:"bins"`                     // 仓位信息
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

type ExchangeTaskInfo struct {
    ExchangeID  uint64              `json:"exchangeId" bson:"exchangeId"`   // 数据库换电ID
    Alternative bool                `json:"alternative" bson:"alternative"` // 是否备选方案
    Model       string              `json:"model" bson:"model"`             // 电池型号
    Empty       *BinInfo            `json:"empty" bson:"empty"`             // 空仓位
    Fully       *BinInfo            `json:"fully" bson:"fully"`             // 满电仓位
    Steps       []*ExchangeStepInfo `json:"steps" bson:"steps"`             // 步骤信息
}

// CurrentStep 获取当前步骤
func (e *ExchangeTaskInfo) CurrentStep() *ExchangeStepInfo {
    return e.Steps[len(e.Steps)-1]
}

// IsLastStep 是否最后一步
func (e *ExchangeTaskInfo) IsLastStep() bool {
    return e.CurrentStep().Step.IsLast()
}

// StartNextStep 开始下一个换电步骤
func (e *ExchangeTaskInfo) StartNextStep() {
    if len(e.Steps) == 0 {
        return
    }

    // 标记上个步骤为成功
    e.Steps[len(e.Steps)-1].Status = TaskStatusSuccess

    // 判断是否最终步骤, 并加入下一个步骤信息
    if len(e.Steps) < int(ExchangeStepPutOut) {
        e.Steps = append(e.Steps, &ExchangeStepInfo{
            Step:   e.CurrentStep().Step.Next(),
            Status: TaskStatusProcessing,
        })
    }
}

// CurrentBin 获取当前操作仓位信息
func (e *ExchangeTaskInfo) CurrentBin() *BinInfo {
    step := e.CurrentStep().Step
    if step < ExchangeStepOpenFull {
        return e.Empty
    }
    return e.Fully
}

// IsSuccess 换电是否成功
func (e *ExchangeTaskInfo) IsSuccess() bool {
    s := e.CurrentStep()
    return s.Step.IsLast() && s.Status.IsSuccess()
}

// StepResult 步骤结果
func (e *ExchangeTaskInfo) StepResult(step ExchangeStep) *ExchangeStepInfo {
    return e.Steps[step-1]
}

type ExchangeInfo struct {
    Cabinet  *ExchangeTaskCabinet `json:"cabinet"`           // 电柜信息
    Exchange *ExchangeTaskInfo    `json:"exchange"`          // 换电信息
    Message  string               `json:"message,omitempty"` // 错误信息
}

type ExchangeStepInfo struct {
    Step   ExchangeStep `json:"step" bson:"step"`     // 操作步骤 1:开空电仓 2:放旧电池 3:开满电仓 4:取新电池
    Status TaskStatus   `json:"status" bson:"status"` // 状态 1:处理中 2:成功 3:失败
    Time   time.Time    `json:"time" bson:"time"`     // 时间
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
