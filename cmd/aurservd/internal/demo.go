// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
)

func Demo() {
    service.NewCabinet().Data(&model.CabinetDataReq{Status: 3})
}
