// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "math"
    "time"
)

type timeTool struct {
}

func NewTime() *timeTool {
    return &timeTool{}
}

// SubDays 计算天数差 after - before
// TODO 计算天数规则: 现行进一法
func (*timeTool) SubDays(after, before time.Time) int {
    return int(math.Ceil(after.Sub(before).Hours() / 24))
}

// SubDaysToNow 距离现在多少天
func (t *timeTool) SubDaysToNow(before time.Time) int {
    return t.SubDays(time.Now(), before)
}

func (t *timeTool) SubDaysToNowString(before string) int {
    return t.SubDaysToNow(carbon.Parse(before).Carbon2Time())
}

// ParseDateString 格式化日期字符串
func (*timeTool) ParseDateString(str string) (time.Time, error) {
    res := carbon.ParseByLayout(str, carbon.DateLayout)
    if res.Error != nil {
        return time.Time{}, res.Error
    }
    return res.Carbon2Time(), nil
}

func (t *timeTool) ParseDateStringX(str string) time.Time {
    res, err := t.ParseDateString(str)
    if err != nil {
        log.Error(err)
        snag.Panic("日期格式错误")
    }
    return res
}
