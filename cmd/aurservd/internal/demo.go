// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "time"
)

func Demo() {
    // for {
    //     if rpc.IsXcBmsClientExists() {
    //         sns := []string{"XCB0862022110268", "XCB0862022110339", "XCB0862022110190", "XCB0862022110130", "XCB0862022110045"}
    //         for _, sn := range sns {
    //             res := service.NewBatteryXc().Position(&model.XcBatteryPositionReq{SN: sn})
    //             b, _ := jsoniter.MarshalIndent(res, "", "  ")
    //             _ = os.WriteFile("./runtime/"+sn+".json", b, 0755)
    //         }
    //
    //         return
    //     }
    // }

    go func() {
        start := time.Now()
        orm := ent.Database.Exchange
        ctx := context.Background()
        q := orm.Query().Where(exchange.InfoNotNil(), exchange.CabinetInfoIsNil(), exchange.StoreIDIsNil())
        for e, err := q.Clone().Limit(1).Exist(ctx); e && err == nil; {
            for _, ex := range q.Clone().Limit(100).AllX(ctx) {
                info := ex.Info
                if info == nil {
                    continue
                }
                var save bool
                updater := ex.Update().SetMessage(info.Message)
                ic := info.Cabinet
                if ic != nil {
                    save = true
                    updater.SetCabinetInfo(&model.ExchangeCabinetInfo{
                        Health:         ic.Health,
                        Doors:          ic.Doors,
                        BatteryNum:     ic.BatteryNum,
                        BatteryFullNum: ic.BatteryFullNum,
                    })
                }

                ie := info.Exchange
                if ie != nil {
                    save = true
                    updater.SetEmpty(ie.Empty).SetFully(ie.Fully).SetSteps(ie.Steps)
                }

                if save {
                    err = updater.Exec(ctx)
                    if err != nil {
                        fmt.Println(ex.ID, err)
                    }
                }
            }
        }
        fmt.Println("所有处理均已完成, 总耗时: ", time.Now().Sub(start))
    }()
}
