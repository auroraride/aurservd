// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    PointRatio = 0.01 // 积分兑换比例, 1:100
)

type PointRes struct {
    Points int64 `json:"points"` // 剩余积分
    Locked int64 `json:"locked"` // 锁定积分
}
