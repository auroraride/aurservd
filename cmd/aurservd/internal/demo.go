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
    jsoniter "github.com/json-iterator/go"
)

func Demo() {
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
