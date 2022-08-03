// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/actuator"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    jsoniter "github.com/json-iterator/go"
)

type OldDetail struct {
    Info struct {
        Id               int64  `json:"id,omitempty"`
        Name             string `json:"name,omitempty"`
        Uuid             string `json:"uuid,omitempty"`
        Model            string `json:"model,omitempty"`
        Serial           string `json:"serial,omitempty"`
        FullIndex        int    `json:"fullIndex,omitempty"`
        PutInDoor        int    `json:"putInDoor,omitempty"`
        EmptyIndex       int    `json:"emptyIndex,omitempty"`
        PutOutDoor       int    `json:"putOutDoor,omitempty"`
        Electricity      int    `json:"electricity,omitempty"`
        RiderElectricity int    `json:"riderElectricity,omitempty"`
    } `json:"info,omitempty"`
    Result struct {
        Step    int    `json:"step,omitempty"`
        Stop    bool   `json:"stop,omitempty"`
        Status  int    `json:"status,omitempty"`
        Message string `json:"message,omitempty"`
    } `json:"result,omitempty"`
    Alternative bool `json:"alternative,omitempty"`
}

func ExchangeDemo() {
    ctx := context.Background()
    orm := ent.Database.Exchange
    items, _ := orm.QueryNotDeleted().WithCabinet().WithSubscribe().Where(
        exchange.CabinetIDNotNil(),
        exchange.DetailNotNil(),
    ).All(ctx)
    for _, item := range items {
        var detail OldDetail
        _ = jsoniter.Unmarshal(item.Detail, &detail)

        cab := item.Edges.Cabinet
        sub := item.Edges.Subscribe

        if cab == nil || sub == nil {
            fmt.Printf("%d \n", item.ID)
            continue
        }

        start := item.StartAt
        if start.IsZero() {
            start = item.CreatedAt
        }

        stop := item.FinishAt
        if stop.IsZero() {
            stop = item.UpdatedAt
        }

        status := actuator.TaskStatusSuccess
        if !item.Success {
            status = actuator.TaskStatusFail
        }

        var steps []*actuator.ExchangeStepInfo

        ss := actuator.ExchangeStep(detail.Result.Step)
        for i := actuator.ExchangeStepOpenEmpty; i <= ss; i++ {
            step := &actuator.ExchangeStepInfo{
                Step: i,
            }
            if i == actuator.ExchangeStepOpenEmpty {
                step.Time = start
            }
            if i != ss {
                step.Status = actuator.TaskStatusSuccess
            } else {
                step.Time = stop
                step.Status = status
            }
            steps = append(steps, step)
        }

        info := &actuator.ExchangeInfo{
            Message: detail.Result.Message,
            Cabinet: actuator.Cabinet{
                Health:         model.CabinetHealthStatusOnline,
                Doors:          cab.Doors,
                BatteryNum:     cab.BatteryNum,
                BatteryFullNum: cab.BatteryFullNum,
            },
            Exchange: &actuator.Exchange{
                Alternative: detail.Alternative,
                Model:       sub.Model,
                Empty: &actuator.BinInfo{
                    Index:       detail.Info.EmptyIndex,
                    Electricity: model.BatteryElectricity(detail.Info.RiderElectricity),
                    Voltage:     -1,
                },
                Fully: &actuator.BinInfo{
                    Index:       detail.Info.FullIndex,
                    Electricity: model.BatteryElectricity(detail.Info.Electricity),
                    Voltage:     -1,
                },
                Steps: steps,
            },
        }

        _, _ = item.Update().SetInfo(info).Save(ctx)
    }
}
