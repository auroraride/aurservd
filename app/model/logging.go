// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Operate uint

type OperatorType uint8

const (
    OperatorTypeManager  OperatorType = iota // 管理员
    OperatorTypeEmployee                     // 店员
    OperatorTypeCabinet                      // 电柜
    OperatorTypeAgent                        // 代理
)

func (ot OperatorType) String() string {
    return []string{"管理员", "店员", "电柜", "代理"}[int(ot)]
}

const (
    OperatePersonBan            Operate = iota // 封禁身份
    OperatePersonUnBan                         // 解封身份
    OperateRiderBLock                          // 封禁账户
    OperateRiderUnBLock                        // 解封账户
    OperateSubscribeAlter                      // 修改订阅时间
    OperateEnterprisePrepayment                // 企业预储值
    OperateSubscribePause                      // 寄存
    OperateSubscribeContinue                   // 取消寄存
    OperateDeposit                             // 调整押金
    OperateProfile                             // 修改资料
    OperateRefund                              // 处理退款
    OperateUnsubscribe                         // 退租
    OperateAssistanceAllocate                  // 救援分配
    OperateAssistanceFree                      // 救援免费
    OperateAssistanceRefuse                    // 救援拒绝
    OperateActive                              // 激活订阅
    OperateCabinetMaintain                     // 电柜维护
    OperateSuspend                             // 暂停扣费
    OperateResume                              // 继续扣费
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
        return "寄存"
    case OperateSubscribeContinue:
        return "取消寄存"
    case OperateDeposit:
        return "调整押金"
    case OperateProfile:
        return "修改资料"
    case OperateRefund:
        return "处理退款"
    case OperateUnsubscribe:
        return "骑手门店退租"
    case OperateAssistanceAllocate:
        return "救援分配"
    case OperateAssistanceFree:
        return "救援免费"
    case OperateAssistanceRefuse:
        return "救援拒绝"
    case OperateActive:
        return "激活订阅"
    case OperateCabinetMaintain:
        return "电柜维护"
    case OperateSuspend:
        return "暂停扣费"
    case OperateResume:
        return "继续扣费"
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
