// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-30
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    "fmt"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/golang-module/carbon/v2"
    "strings"
    "time"
)

// BatteryLog 电池变动日志
type BatteryLog struct {
    UUID   string              `json:"uuid" sls:"操作ID" index:"doc"`
    Status model.CabinetStatus `json:"status" sls:"电柜状态" index:"doc"`
    Brand  string              `json:"brand" sls:"品牌" index:"doc"`
    Serial string              `json:"serial" sls:"编码" index:"doc"`
    Number int                 `json:"number" sls:"变动数量" index:"doc"`
    Time   string              `json:"time" sls:"时间" index:"doc"`
    From   int                 `json:"from" sls:"变动前数量" index:"doc"`
    To     int                 `json:"to" sls:"变动后数量" index:"doc"`
    Task   string              `json:"task" sls:"电柜任务"`
    Info   string              `json:"info" sls:"电柜信息"`
}

func (l *BatteryLog) GetLogstoreName() string {
    return ar.Config.Aliyun.Sls.BatteryLog
}

func (l *BatteryLog) Send() {
    PutLog(l)
}

func NewBatteryLog(brand, serial string, from, to int, updatedAt time.Time) *BatteryLog {
    return &BatteryLog{
        Brand:  model.CabinetBrand(brand).String(),
        Serial: serial,
        Time:   updatedAt.Format(carbon.DateTimeLayout),
        Number: to - from,
        From:   from,
        To:     to,
    }
}

// SetTask 设置电柜任务
func (l *BatteryLog) SetTask(task *ec.Task) *BatteryLog {
    if task == nil {
        return l
    }
    l.UUID = task.ID
    l.Task = task.Job.Label()
    l.Info = task.String()
    return l
}

func (l *BatteryLog) SetStatus(status model.CabinetStatus) *BatteryLog {
    l.Status = status
    return l
}

func (l *BatteryLog) SetBin(oldBins, bins model.CabinetBins) *BatteryLog {
    oldMap := make(map[int]*model.CabinetBin)

    for _, bin := range oldBins {
        oldMap[bin.Index] = bin
    }

    var diff []string
    for _, bin := range bins {
        old, ok := oldMap[bin.Index]
        if !ok {
            continue
        }

        if old.Battery == bin.Battery {
            continue
        }

        x := 1
        if old.Battery {
            x = -1
        }

        diff = append(diff, fmt.Sprintf("%s: %d", bin.Name, x))
    }

    if len(l.Info) > 0 {
        l.Info += "\n"
    }
    l.Info += fmt.Sprintf("仓位变化: %s", strings.Join(diff, ";"))

    return l
}
