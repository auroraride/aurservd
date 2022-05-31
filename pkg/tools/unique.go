// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-28
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "fmt"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    "github.com/sony/sonyflake"
    "strings"
    "time"
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

// NewSN28 生成28位字符串
func (u *unique) NewSN28() string {
    return fmt.Sprintf(
        "%s%d",
        strings.ReplaceAll(time.Now().Format(carbon.ShortDateTimeMicroLayout), ".", ""),
        utils.RandomIntMaxMin(10000000, 99999999),
    )
}
