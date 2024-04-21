// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-28
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/sony/sonyflake"

	"github.com/auroraride/aurservd/pkg/utils"
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
	str := strings.ReplaceAll(time.Now().Format(carbon.ShortDateTimeMicroLayout), ".", "")
	n := 20 - len(str)
	if n > 0 {
		str += strings.Repeat("0", n)
	}
	return str + strconv.FormatUint(uint64(utils.RandIntMaxMin(10000000, 99999999)), 10)
}

// Rand 生成指定长度字符串
func (u *unique) Rand(n int) string {
	str := strings.ReplaceAll(time.Now().Format(carbon.ShortDateTimeMicroLayout), ".", "")
	x := n - len(str)

	maxLimit := int(math.Pow10(x)) - 1
	lowLimit := int(math.Pow10(x - 1))

	str += strconv.FormatUint(uint64(utils.RandIntMaxMin(lowLimit, maxLimit)), 10)
	return str
}

func (u *unique) NewSN() string {
	str := strings.ReplaceAll(time.Now().Format(carbon.ShortDateTimeMilliLayout), ".", "")
	n := 17 - len(str)
	if n > 0 {
		str += strings.Repeat("0", n)
	}
	return str
}
