// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-28
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "fmt"
    "github.com/sony/sonyflake"
)

type unique struct {
}

func NewUnique() *unique {
    return &unique{}
}

func (*unique) NewSonyflakeID() string {
    sid, _ := sonyflake.NewSonyflake(sonyflake.Settings{}).NextID()
    return fmt.Sprintf("%d", sid)
}
