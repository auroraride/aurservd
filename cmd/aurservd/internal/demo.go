// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    jsoniter "github.com/json-iterator/go"
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
                    steps := make([]*model.ExchangeStepInfo, len(ie.Steps))
                    for i, st := range ie.Steps {
                        steps[i] = &model.ExchangeStepInfo{
                            Step:   model.ExchangeStep(st.Step.Int()),
                            Status: model.TaskStatus(uint8(st.Status)),
                            Time:   st.Time,
                        }
                    }
                    updater.SetEmpty(ie.Empty).SetFully(ie.Fully).SetSteps(steps)
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

    if !ar.Config.Debug {
        return
    }
    // demoOrder()
}

func demoOrder() {
    ctx := context.Background()
    fmt.Println(ent.Database.ExecContext(ctx, `TRUNCATE TABLE "subscribe"; TRUNCATE TABLE "order"; TRUNCATE TABLE "contract"; TRUNCATE TABLE "allocate"; UPDATE ebike SET status = 0, rider_id = NULL WHERE status IS NOT NULL;`))
    var result model.PaymentCache
    _ = jsoniter.Unmarshal(demoEbikeOrder(), &result)
    service.NewOrder().OrderPaid(result.Subscribe)
}

func demoBatteryOrder() []byte {
    return []byte(`{"cacheType":1,"create":{"cityId":610100,"orderType":1,"outTradeNo":"202210151653590733151174430","riderId":98784248576,"name":"购买车电骑士卡","amount":0.01,"payway":2,"plan":{"id":94489280549,"name":"测试骑行卡"},"deposit":0,"pastDays":null,"commission":20,"model":"60V26AH","days":30,"orderId":null,"subscribeId":null,"points":0,"pointRatio":0.01,"couponAmount":0,"coupons":null,"reliefNewly":30}}`)
}

func demoEbikeOrder() []byte {
    return []byte(`{"cacheType":1,"create":{"cityId":610100,"orderType":1,"outTradeNo":"202210151653590733151174430","riderId":98784248576,"name":"购买车电骑士卡","amount":0.01,"payway":2,"plan":{"id":94489280549,"name":"测试骑行卡"},"deposit":0,"pastDays":null,"commission":20,"model":"60V26AH","days":30,"orderId":null,"subscribeId":null,"points":0,"pointRatio":0.01,"couponAmount":0,"coupons":null,"reliefNewly":30,"ebikeBrandId":210453397505}}`)
}
