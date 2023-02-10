// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// TaskJob 电柜任务
type TaskJob string

const (
    JobExchange         TaskJob = "RDR_EXCHANGE"    // 骑手-换电
    JobRiderActive      TaskJob = "RDR_ACTIVE"      // 骑手-激活
    JobRiderUnSubscribe TaskJob = "RDR_UNSUBSCRIBE" // 骑手-退租
    JobPause            TaskJob = "RDR_PAUSE"       // 骑手-寄存
    JobContinue         TaskJob = "RDR_CONTINUE"    // 骑手-继续
    JobManagerOpen      TaskJob = "MGR_OPEN"        // 管理-开门
    JobManagerLock      TaskJob = "MGR_LOCK"        // 管理-锁仓
    JobManagerUnLock    TaskJob = "MGR_UNLOCK"      // 管理-解锁
    JobManagerReboot    TaskJob = "MGR_REBOOT"      // 管理-重启
    JobManagerExchange  TaskJob = "MGR_EXCHANGE"    // 管理-换电
)

func (j TaskJob) Label() string {
    switch j {
    case JobExchange:
        return "骑手换电"
    case JobRiderActive:
        return "骑手激活"
    case JobRiderUnSubscribe:
        return "骑手退租"
    case JobPause:
        return "骑手寄存"
    case JobContinue:
        return "骑手继续"
    case JobManagerOpen:
        return "管理开门"
    case JobManagerLock:
        return "管理锁仓"
    case JobManagerUnLock:
        return "管理解锁"
    case JobManagerReboot:
        return "管理重启"
    case JobManagerExchange:
        return "管理换电"
    }
    return "未知任务"
}

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
