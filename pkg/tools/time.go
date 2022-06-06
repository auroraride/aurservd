// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type timeTool struct {
}

func NewTime() *timeTool {
    return &timeTool{}
}

// DiffDaysOfStart 获取两个日期相差天数(a - b), 计算方式为a / b都转换为当日的0天(不将当日计算在内)
func (t *timeTool) DiffDaysOfStart(after, before time.Time) int {
    return int(carbon.Time2Carbon(before).StartOfDay().DiffInDays(carbon.Time2Carbon(after).StartOfDay()))
}

func (t *timeTool) DiffDaysOfStartToNow(before time.Time) int {
    return int(carbon.Time2Carbon(before).StartOfDay().DiffInDays(carbon.Now().StartOfDay()))
}

func (t *timeTool) DiffDaysOfStartToNowString(before string) int {
    return int(carbon.Parse(before).StartOfDay().DiffInDays(carbon.Now().StartOfDay()))
}

// DiffDaysOfNextDay 获取两个日期相差天数(a - b), 计算方式为a转换为当日的0天, b转换为第二天的0天(将当日计算在内)
func (t *timeTool) DiffDaysOfNextDay(after, before time.Time) int {
    return int(carbon.Time2Carbon(before).StartOfDay().DiffInDays(carbon.Time2Carbon(after).StartOfDay().AddDay()))
}

func (t *timeTool) DiffDaysOfNextDayToNow(before time.Time) int {
    return int(carbon.Time2Carbon(before).StartOfDay().DiffInDays(carbon.Now().StartOfDay().AddDay()))
}

func (t *timeTool) DiffDaysOfNextDayToNowString(before string) int {
    return int(carbon.Parse(before).StartOfDay().DiffInDays(carbon.Now().StartOfDay().AddDay()))
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
