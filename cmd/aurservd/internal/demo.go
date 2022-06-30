// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "github.com/auroraride/aurservd/internal/ent"
)

func Demo() {
    ctx := context.Background()
    cab, _ := ent.Database.Cabinet.QueryNotDeleted().First(ctx)
    bins := cab.Bin
    bins[1].Battery = true
    cab.Update().SetBin(bins).SaveX(ctx)
}
