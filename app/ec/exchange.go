// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package ec

import (
    "github.com/auroraride/aurservd/app/model"
)

// Exchange 换电信息
type Exchange struct {
    ExchangeID  uint64                    `json:"exchangeId" bson:"exchangeId"`   // 数据库换电ID
    Alternative bool                      `json:"alternative" bson:"alternative"` // 是否备选方案
    Model       string                    `json:"model" bson:"model"`             // 电池型号
    Empty       *model.BinInfo            `json:"empty" bson:"empty"`             // 空仓位
    Fully       *model.BinInfo            `json:"fully" bson:"fully"`             // 满电仓位
    Steps       []*model.ExchangeStepInfo `json:"steps" bson:"steps"`             // 步骤信息
}

// CurrentStep 获取当前步骤
func (e *Exchange) CurrentStep() *model.ExchangeStepInfo {
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
    e.Steps[len(e.Steps)-1].Status = model.TaskStatusSuccess

    // 判断是否最终步骤, 并加入下一个步骤信息
    if len(e.Steps) < int(model.ExchangeStepPutOut) {
        e.Steps = append(e.Steps, &model.ExchangeStepInfo{
            Step:   e.CurrentStep().Step.Next(),
            Status: model.TaskStatusProcessing,
        })
    }
}

// CurrentBin 获取当前操作仓位信息
func (e *Exchange) CurrentBin() *model.BinInfo {
    step := e.CurrentStep().Step
    if step < model.ExchangeStepOpenFull {
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
func (e *Exchange) StepResult(step model.ExchangeStep) *model.ExchangeStepInfo {
    return e.Steps[step-1]
}

type ExchangeInfo struct {
    Cabinet  *Cabinet  `json:"cabinet"`           // 电柜信息
    Exchange *Exchange `json:"exchange"`          // 换电信息
    Message  string    `json:"message,omitempty"` // 错误信息
}
