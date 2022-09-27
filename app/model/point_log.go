// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "golang.org/x/exp/slices"

type PointLogType uint8

const (
    PointLogTypeConsume PointLogType = 1 + iota // 消费
    PointLogTypeAward                           // 奖励
)

var (
    PointLogTypes = []PointLogType{PointLogTypeConsume, PointLogTypeAward}
)

func (t PointLogType) IsValid() bool {
    return slices.Contains(PointLogTypes, t)
}

func (t PointLogType) Value() uint8 {
    return uint8(t)
}

type PointModifyReq struct {
    RiderID uint64       `json:"riderId" validate:"required" trans:"骑手"`
    Points  int64        `json:"points" validate:"required" trans:"积分"`
    Reason  string       `json:"reason" validate:"required" trans:"原因"`
    Type    PointLogType `json:"type" validate:"required,enum" trans:"类别" enums:"1,2"` // 1:消费 2:奖励
}
