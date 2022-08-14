// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type ReserveStatus uint8

const (
    ReserveStatusPending    ReserveStatus = iota // 已预约
    ReserveStatusProcessing                      // 进行中
    ReserveStatusDone                            // 已完成
    ReserveStatusOverdue                         // 已超时
    ReserveStatusCancel                          // 已取消
)

func (rs ReserveStatus) Value() uint8 {
    return uint8(rs)
}

func (rs ReserveStatus) String() string {
    switch rs {
    default:
        return "已预约"
    case ReserveStatusDone:
        return "已完成"
    case ReserveStatusProcessing:
        return "进行中"
    case ReserveStatusOverdue:
        return "已超时"
    case ReserveStatusCancel:
        return "已取消"
    }
}
