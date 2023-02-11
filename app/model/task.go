// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-10
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "database/sql/driver"
)

// TaskJob 电柜任务
type TaskJob string

func (j TaskJob) Values() []string {
    return []string{
        TaskJobExchange.String(),
        TaskJobRiderActive.String(),
        TaskJobRiderUnSubscribe.String(),
        TaskJobPause.String(),
        TaskJobContinue.String(),
        TaskJobManagerOpen.String(),
        TaskJobManagerLock.String(),
        TaskJobManagerUnLock.String(),
        TaskJobManagerReboot.String(),
        TaskJobManagerExchange.String(),
    }
}

func (j TaskJob) String() string {
    return string(j)
}

const (
    TaskJobExchange         TaskJob = "RDR_EXCHANGE"    // 骑手-换电
    TaskJobRiderActive      TaskJob = "RDR_ACTIVE"      // 骑手-激活
    TaskJobRiderUnSubscribe TaskJob = "RDR_UNSUBSCRIBE" // 骑手-退租
    TaskJobPause            TaskJob = "RDR_PAUSE"       // 骑手-寄存
    TaskJobContinue         TaskJob = "RDR_CONTINUE"    // 骑手-继续
    TaskJobManagerOpen      TaskJob = "MGR_OPEN"        // 管理-开门
    TaskJobManagerLock      TaskJob = "MGR_LOCK"        // 管理-锁仓
    TaskJobManagerUnLock    TaskJob = "MGR_UNLOCK"      // 管理-解锁
    TaskJobManagerReboot    TaskJob = "MGR_REBOOT"      // 管理-重启
    TaskJobManagerExchange  TaskJob = "MGR_EXCHANGE"    // 管理-换电
)

func (j TaskJob) Label() string {
    switch j {
    case TaskJobExchange:
        return "骑手换电"
    case TaskJobRiderActive:
        return "骑手激活"
    case TaskJobRiderUnSubscribe:
        return "骑手退租"
    case TaskJobPause:
        return "骑手寄存"
    case TaskJobContinue:
        return "骑手继续"
    case TaskJobManagerOpen:
        return "管理开门"
    case TaskJobManagerLock:
        return "管理锁仓"
    case TaskJobManagerUnLock:
        return "管理解锁"
    case TaskJobManagerReboot:
        return "管理重启"
    case TaskJobManagerExchange:
        return "管理换电"
    }
    return "未知任务"
}

type TaskStatus int

func (ts *TaskStatus) Scan(t interface{}) (err error) {
    *ts = TaskStatus(t.(int64))
    return
}

func (ts TaskStatus) Value() (driver.Value, error) {
    return int64(ts), nil
}

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
