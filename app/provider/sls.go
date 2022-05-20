// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-17
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "strings"
    "time"
)

// GenerateSlsStatusLogGroup 生成status log日志
func GenerateSlsStatusLogGroup(cabinet *ent.Cabinet) (lg *sls.LogGroup) {
    t := tea.Uint32(uint32(time.Now().Unix()))
    lg = &sls.LogGroup{}
    logs := make([]*sls.Log, len(cabinet.Bin))
    for i, bin := range cabinet.Bin {
        c := new(logging.CabinetLog)
        _ = copier.Copy(c, bin)
        c.Serial = cabinet.Serial
        c.Errors = strings.Join(bin.ChargerErrors, ",")
        c.Brand = model.CabinetBrand(cabinet.Brand).String()
        c.Time = time.Now().Format(carbon.DateTimeLayout)
        logs[i] = &sls.Log{
            Time:     t,
            Contents: logging.ParseLogContent(c),
        }
    }
    lg.Logs = logs
    return
}
