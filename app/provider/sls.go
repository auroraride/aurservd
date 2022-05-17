// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-17
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "time"
)

type CabinetLog struct {
    Times int `json:"times"` // 第%d轮
}

func GenerateSlsLogGroup(cabinet *ent.Cabinet) (lg *sls.LogGroup) {
    lg = &sls.LogGroup{
        Source: tea.String(cabinet.Serial),
        Topic:  tea.String(model.CabinetBrand(cabinet.Brand).Name()),
    }
    logs := []*sls.Log{
        {
            Time: tea.Uint32(uint32(time.Now().Unix())),
            Contents: []*sls.LogContent{
                {

                },
                {
                    Key:   tea.String("name"),
                    Value: tea.String("1号仓"),
                },
                {
                    Key:   tea.String("battery"),
                    Value: tea.String("true"),
                },
            },
        },
    }
    lg.Logs = logs
    return
}
