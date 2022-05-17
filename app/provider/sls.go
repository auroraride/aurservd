// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-17
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "fmt"
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/jinzhu/copier"
    "reflect"
    "strings"
    "time"
)

type CabinetLog struct {
    Serial      string                   `sls:"编码"`
    Name        string                   `sls:"仓位"`
    BatterySN   string                   `sls:"电池序列号"`
    Full        bool                     `sls:"是否满电"`
    Battery     bool                     `sls:"是否有电池"`
    Electricity model.BatteryElectricity `sls:"当前电量"`
    OpenStatus  bool                     `sls:"是否开门"`
    DoorHealth  bool                     `sls:"仓门是否正常"`
    Current     float64                  `sls:"充电电流(A)"`
    Voltage     float64                  `sls:"电压(V)"`
    Errors      string                   `sls:"故障信息"`
    Remark      string                   `sls:"备注"`
}

type OperationLog struct {
    User      string                   `sls:"操作人"`
    UserID    uint64                   `sls:"操作人ID"`
    Phone     string                   `sls:"操作人电话"`
    Serial    string                   `sls:"编码"`
    Name      string                   `sls:"仓位"`
    Operation model.CabinetDoorOperate `sls:"操作"`
    Success   bool                     `sls:"是否成功"`
    Remark    string                   `sls:"备注"`
}

// ParseLogContent 转换为sls日志
func ParseLogContent(pointer any) (contents []*sls.LogContent) {
    t := reflect.TypeOf(pointer).Elem()
    n := t.NumField()
    value := reflect.ValueOf(pointer).Elem()

    contents = make([]*sls.LogContent, n)
    for i := 0; i < n; i++ {
        tag, _ := t.Field(i).Tag.Lookup("sls")
        v := value.Field(i)
        cv := ""
        if v.Type().Kind() == reflect.Bool {
            cv = "否"
            if v.Bool() {
                cv = "是"
            }
        } else {
            cv = fmt.Sprintf("%v", v.Interface())
        }
        contents[i] = &sls.LogContent{
            Key:   tea.String(tag),
            Value: tea.String(cv),
        }
    }
    return
}

// GenerateSlsStatusLogGroup 生成status log日志
func GenerateSlsStatusLogGroup(cabinet *ent.Cabinet) (lg *sls.LogGroup) {
    t := tea.Uint32(uint32(time.Now().Unix()))
    lg = &sls.LogGroup{
        Source: tea.String(cabinet.Serial),
        Topic:  tea.String(model.CabinetBrand(cabinet.Brand).String()),
    }
    logs := make([]*sls.Log, len(cabinet.Bin))
    for i, bin := range cabinet.Bin {
        c := new(CabinetLog)
        _ = copier.Copy(c, bin)
        c.Serial = cabinet.Serial
        c.Errors = strings.Join(bin.ChargerErrors, ",")
        logs[i] = &sls.Log{
            Time:     t,
            Contents: ParseLogContent(c),
        }
    }
    lg.Logs = logs
    return
}
