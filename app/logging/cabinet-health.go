// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-29
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    "github.com/auroraride/adapter"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/golang-module/carbon/v2"
    "time"
)

type HealthLog struct {
    Brand  string `json:"brand" sls:"品牌" index:"doc"`
    Serial string `json:"serial" sls:"编码" index:"doc"`
    From   string `json:"from" sls:"旧状态" index:"doc"`
    To     string `json:"to" sls:"新状态" index:"doc"` // 0离线 1在线 2故障
    Time   string `json:"time" sls:"时间" index:"doc"`
}

func (l *HealthLog) GetLogstoreName() string {
    return ar.Config.Aliyun.Sls.HealthLog
}

func (l *HealthLog) Send() {
    PutLog(l)
}

func NewHealthLog(brand adapter.CabinetBrand, serial string, updatedAt time.Time) *HealthLog {
    return &HealthLog{
        Brand:  brand.String(),
        Serial: serial,
        Time:   updatedAt.Format(carbon.DateTimeLayout),
    }
}

func (l *HealthLog) GetStatus(status uint8) string {
    switch status {
    case model.CabinetHealthStatusOnline:
        return "在线"
    case model.CabinetHealthStatusFault:
        return "在线"
    default:
        return "离线"
    }
}

func (l *HealthLog) SetStatus(from, to uint8) *HealthLog {
    l.From = l.GetStatus(from)
    l.To = l.GetStatus(to)
    return l
}
