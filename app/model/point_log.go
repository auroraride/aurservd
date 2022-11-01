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

func (t PointLogType) String() string {
    switch t {
    case PointLogTypeConsume:
        return "消费"
    case PointLogTypeAward:
        return "奖励"
    }
    return " - "
}

type PointLogAttach struct {
    Plan      *Plan      `json:"plan,omitempty"`
    PointGift *PointGift `json:"pointGift,omitempty"`
}

type PointModifyReq struct {
    RiderID uint64       `json:"riderId" validate:"required" trans:"骑手"`
    Points  int64        `json:"points" validate:"required" trans:"积分"`
    Reason  string       `json:"reason" validate:"required" trans:"原因"`
    Type    PointLogType `json:"type" validate:"required,enum" trans:"类别" enums:"1,2"` // 1:消费 2:奖励
}

type PointLogListReq struct {
    PaginationReq
    RiderID uint64       `json:"riderId" query:"riderId"`       // 骑手
    Keyword string       `json:"keyword" query:"keyword"`       // 骑手姓名或电话
    Type    PointLogType `json:"type" enums:"1,2" query:"type"` // 类别, 1:消费 2:奖励
}

type PointLogListRes struct {
    ID       uint64  `json:"id"`
    Type     string  `json:"type"`               // 类别
    Plan     string  `json:"plan,omitempty"`     // 订单骑士卡
    Points   int64   `json:"points"`             // 积分
    After    int64   `json:"after"`              // 变动后
    Reason   *string `json:"reason,omitempty"`   // 原因
    Modifier string  `json:"modifier,omitempty"` // 操作人
    Time     string  `json:"time"`               // 时间
}

type PointBatchReq struct {
    Phones []string     `json:"phones" validate:"required,min=1" trans:"电话"`
    Points int64        `json:"points" validate:"required" trans:"积分"`
    Reason string       `json:"reason" validate:"required" trans:"原因"`
    Type   PointLogType `json:"type" validate:"required,enum" trans:"类别" enums:"1,2"` // 1:消费 2:奖励
}

type PointGift struct {
    Amount     float64 `json:"amount"`     // 消费金额
    Proportion float64 `json:"proportion"` // 赠送比例
}
