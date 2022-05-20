// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type UserPlanItem struct {
    Days  uint         `json:"days"`  // 剩余天数
    Model BatteryModel `json:"model"` // 电池型号
}
