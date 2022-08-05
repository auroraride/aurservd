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
    "time"
)

// BatteryLog 电池变动日志
type BatteryLog struct {
    UUID     string              `json:"uuid" sls:"操作ID" index:"doc"`
    Brand    string              `json:"brand" sls:"品牌" index:"doc"`
    Serial   string              `json:"serial" sls:"编码" index:"doc"`
    Number   int                 `json:"number" sls:"变动数量" index:"doc"`
    Exchange string              `json:"exchange" sls:"换电信息" index:"doc"`
    Time     string              `json:"time" sls:"时间" index:"doc"`
    From     int                 `json:"from" sls:"变动前数量" index:"doc"`
    To       int                 `json:"to" sls:"变动后数量" index:"doc"`
    Status   model.CabinetStatus `json:"status" sls:"电柜状态" index:"doc"`
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

// SetExchangeProcess 设置换电信息
func (l *BatteryLog) SetExchangeProcess(task *ec.Task) *BatteryLog {
    if task == nil {
        return l
    }
    l.UUID = task.ID.Hex()
    l.Exchange = fmt.Sprintf(
        "ID: %d, 电话: %s, 名字: %s, 步骤: %s, 空: %d仓, 满: %d仓",
        task.Rider.ID,
        task.Rider.Phone,
        task.Rider.Name,
        task.Exchange.CurrentStep().Step,
        task.Exchange.Empty.Index+1,
        task.Exchange.Fully.Index+1,
    )
    return l
}

func (l *BatteryLog) SetStatus(status model.CabinetStatus) *BatteryLog {
    l.Status = status
    return l
}

func (l *BatteryLog) SetBin(old, value model.CabinetBins) *BatteryLog {
    om := make(map[int]*model.CabinetBin)

    for _, bin := range old {
        om[bin.Index] = bin
    }

    var diff []string
    for _, bin := range value {
        o, ok := om[bin.Index]
        if ok && o.Battery == bin.Battery {
            continue
        }

        x := 1
        if o.Battery {
            x = -1
        }

        diff = append(diff, fmt.Sprintf("%s: %d", bin.Name, x))
    }

    return l
}
