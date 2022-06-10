// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Operate uint

const (
    OperatePersonBan            = iota // 封禁身份
    OperatePersonUnBan                 // 解封身份
    OperateRiderBLock                  // 封禁账户
    OperateRiderUnBLock                // 解封账户
    OperateSubscribeAlter              // 修改订阅时间
    OperateEnterprisePrepayment        // 企业预储值
    OperateSubscribePause              // 暂停计费
    OperateSubscribeContinue           // 继续计费
    OperateDeposit                     // 调整押金
)

func (o Operate) String() string {
    switch o {
    case OperatePersonBan:
        return "封禁用户"
    case OperatePersonUnBan:
        return "解封用户"
    case OperateRiderBLock:
        return "封禁账户"
    case OperateRiderUnBLock:
        return "解封账户"
    case OperateSubscribeAlter:
        return "修改时间"
    case OperateEnterprisePrepayment:
        return "企业预储值"
    case OperateSubscribePause:
        return "暂停计费"
    case OperateSubscribeContinue:
        return "继续计费"
    case OperateDeposit:
        return "调整押金"
    default:
        return "未知操作"
    }
}

// LogOperate 操作日志
type LogOperate struct {
    Operate     string `json:"operate"`     // 操作类别
    Before      string `json:"before"`      // 操作之前
    After       string `json:"after"`       // 操作之后
    ManagerName string `json:"managerName"` // 操作人姓名
    Phone       string `json:"phone"`       // 操作人电话
    Time        string `json:"time"`        // 操作时间
}
