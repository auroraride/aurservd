// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/pkg/tools"
)

func Demo() {
    service.NewOrder().OrderPaid(&model.PaymentSubscribe{
        CityID:      130100,
        OrderType:   2,
        OutTradeNo:  "xxxxxxxxxx",
        RiderID:     98784248996,
        Name:        "testname",
        TradeNo:     "xxxxaaaaa",
        PlanID:      94489280540,
        Model:       "60V26AH",
        Days:        10,
        SubscribeID: tools.NewPointer().UInt64(115964118051),
    })
}
