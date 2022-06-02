// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// LogOperate 操作日志
type LogOperate struct {
    Operate     string `json:"operate"`     // 操作类别
    Before      string `json:"before"`      // 操作之前
    After       string `json:"after"`       // 操作之后
    ManagerName string `json:"managerName"` // 操作人姓名
    Phone       string `json:"phone"`       // 操作人电话
    Time        string `json:"time"`        // 操作时间
}
