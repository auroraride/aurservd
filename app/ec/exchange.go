// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package ec

import (
    "fmt"
    "github.com/golang-module/carbon/v2"
    "time"
)

// ExchangeStep 换电步骤
// RiderCabinetOperateStep
type ExchangeStep uint8

const (
    ExchangeStepOpenEmpty ExchangeStep = iota + 1 // 第一步, 开启空电仓
    ExchangeStepPutInto                           // 第二步, 放入旧电池并关闭仓门
    ExchangeStepOpenFull                          // 第三步, 开启满电仓
    ExchangeStepPutOut                            // 第四步, 取出新电池并关闭仓门
)

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

// Exchange 换电信息
type Exchange struct {
    ExchangeID  uint64              `json:"exchangeId" bson:"exchangeId"`   // 数据库换电ID
    Alternative bool                `json:"alternative" bson:"alternative"` // 是否备选方案
    Model       string              `json:"model" bson:"model"`             // 电池型号
    Empty       *BinInfo            `json:"empty" bson:"empty"`             // 空仓位
    Fully       *BinInfo            `json:"fully" bson:"fully"`             // 满电仓位
    Steps       []*ExchangeStepInfo `json:"steps" bson:"steps"`             // 步骤信息
}

// CurrentStep 获取当前步骤
func (e *Exchange) CurrentStep() *ExchangeStepInfo {
    return e.Steps[len(e.Steps)-1]
}

// IsLastStep 是否最后一步
func (e *Exchange) IsLastStep() bool {
    return e.CurrentStep().Step.IsLast()
}

// StartNextStep 开始下一个换电步骤
func (e *Exchange) StartNextStep() {
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
func (e *Exchange) CurrentBin() *BinInfo {
    step := e.CurrentStep().Step
    if step < ExchangeStepOpenFull {
        return e.Empty
    }
    return e.Fully
}

// IsSuccess 换电是否成功
func (e *Exchange) IsSuccess() bool {
    s := e.CurrentStep()
    return s.Step.IsLast() && s.Status.IsSuccess()
}

// StepResult 步骤结果
func (e *Exchange) StepResult(step ExchangeStep) *ExchangeStepInfo {
    return e.Steps[step-1]
}

type ExchangeInfo struct {
    Cabinet  Cabinet   `json:"cabinet"`           // 电柜信息
    Exchange *Exchange `json:"exchange"`          // 换电信息
    Message  string    `json:"message,omitempty"` // 错误信息
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
