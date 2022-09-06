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

// DiffDays 按照标准时间计算天数 b hh:mm:ss -> a hh:mm:ss
// 用来计算[寄存时间]
func (t *timeTool) DiffDays(after, before time.Time) int {
    return int(carbon.Time2Carbon(before).DiffInDays(carbon.Time2Carbon(after)))
}

func (t *timeTool) DiffDaysToNow(before time.Time) int {
    return t.DiffDays(time.Now(), before)
}

func (t *timeTool) DiffDaysToNowString(before string) int {
    return int(carbon.Parse(before).DiffInDays(carbon.Time2Carbon(time.Now())))
}

// LastDays 获取两个日期相差天数(a - b), b 00:00:00 -> a 00:00:00
// 用来计算[剩余时间]
func (t *timeTool) LastDays(after, before time.Time) int {
    return int(carbon.Time2Carbon(before).StartOfDay().DiffInDays(carbon.Time2Carbon(after).StartOfDay()))
}

func (t *timeTool) LastDaysToNow(after time.Time) int {
    return int(carbon.Now().StartOfDay().DiffInDays(carbon.Time2Carbon(after).StartOfDay()))
}

func (t *timeTool) LastDaysToNowString(after string) int {
    return int(carbon.Now().StartOfDay().DiffInDays(carbon.Parse(after).StartOfDay()))
}

// UsedDays 获取两个日期相差天数(a - b), b 00:00:00 -> a+1 00:00:00
// 用来计算[已使用时间]
func (t *timeTool) UsedDays(after, before time.Time) int {
    return int(carbon.Time2Carbon(before).StartOfDay().DiffInDays(carbon.Time2Carbon(after).StartOfDay().AddDay()))
}

func (t *timeTool) UsedDaysToNow(before time.Time) int {
    return int(carbon.Time2Carbon(before).StartOfDay().DiffInDays(carbon.Now().StartOfDay().AddDay()))
}

func (t *timeTool) UsedDaysToNowString(before string) int {
    return int(carbon.Parse(before).StartOfDay().DiffInDays(carbon.Now().StartOfDay().AddDay()))
}

// WillEnd 计算到期时间
// 从今日算起, 天数应该减一
// 若从明日算起, 天数不应减一
// params: 0 是否从明日算起, 默认从今天算
func (t *timeTool) WillEnd(before time.Time, days int, params ...bool) time.Time {
    // 是否从明日开始
    if len(params) == 0 || !params[0] {
        days -= 1
    }
    return before.AddDate(0, 0, days)
}

// ParseDateString 格式化日期字符串
func (*timeTool) ParseDateString(str string) (time.Time, error) {
    res := carbon.ParseByLayout(str, carbon.DateLayout)
    if res.Error != nil {
        return time.Time{}, res.Error
    }
    return res.StartOfDay().Carbon2Time(), nil
}

func (t *timeTool) ParseDateStringX(str string) time.Time {
    res, err := t.ParseDateString(str)
    if err != nil {
        log.Error(err)
        snag.Panic("日期格式错误")
    }
    return res
}

// ParseNextDateString 格式化日期字符串到下一天
func (t *timeTool) ParseNextDateString(str string) (time.Time, error) {
    res := carbon.ParseByLayout(str, carbon.DateLayout)
    if res.Error != nil {
        return time.Time{}, res.Error
    }
    return res.StartOfDay().AddDay().Carbon2Time(), nil
}

func (t *timeTool) ParseNextDateStringX(str string) time.Time {
    res, err := t.ParseNextDateString(str)
    if err != nil {
        log.Error(err)
        snag.Panic("日期格式错误")
    }
    return res
}

// PauseBeginning 暂停或寄存开始日期计算
// start 暂停或寄存开始时间
func (t *timeTool) PauseBeginning(start time.Time) time.Time {
    startDay := carbon.Time2Carbon(start).StartOfDay().Carbon2Time()
    // 判定开始时间是否0点
    if startDay.Equal(start) {
        return startDay
    }
    return carbon.Time2Carbon(startDay).Tomorrow().Carbon2Time()
}
