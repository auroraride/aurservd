// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "github.com/auroraride/aurservd/internal/ar"
)

func Demo() {
    setStore()
    setCabinet()
}

func setStore() {
    ctx := context.Background()
    items, _ := ar.Ent.Store.QueryNotDeleted().WithBranch().All(ctx)
    for _, item := range items {
        b := item.Edges.Branch
        if b == nil {
            continue
        }
        item.Update().SetLng(b.Lng).SetLat(b.Lat).SetAddress(b.Address).SaveX(ctx)
    }
}

func setCabinet() {
    ctx := context.Background()
    items, _ := ar.Ent.Cabinet.QueryNotDeleted().WithBranch().All(ctx)
    for _, item := range items {
        b := item.Edges.Branch
        if b == nil {
            continue
        }
        item.Update().SetLng(b.Lng).SetLat(b.Lat).SetAddress(b.Address).SaveX(ctx)
    }
}
