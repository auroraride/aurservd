// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type SubscribeAllocate struct {
    EbikeID     *uint64 `json:"ebikeId"` // 电车ID
    SubscribeID uint64  `json:"subscribeId" validate:"required" trans:"订阅ID"`
}
