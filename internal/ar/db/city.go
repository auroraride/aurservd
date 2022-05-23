// Copyright (C) liasica. 2022-present.
//
// Created at 2022-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package db

import (
    "context"
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/internal/ent"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
)

func insertCity(client *ent.Client) {
    type R struct {
        Adcode   uint64  `json:"adcode"`
        Name     string  `json:"name"`
        Code     string  `json:"code"`
        Lng      float64 `json:"lng,omitempty"`
        Lat      float64 `json:"lat,omitempty"`
        Children []R     `json:"children,omitempty"`
    }

    ctx := context.Background()
    tx, _ := client.Tx(ctx)
    // 导入城市
    log.Println("导入城市")
    var items []R
    err := jsoniter.Unmarshal(assets.City, &items)
    if err == nil {
        for _, item := range items {
            parent := client.City.Create().
                SetID(item.Adcode).
                SetName(item.Name).
                SetCode(item.Code).
                SaveX(ctx)
            for _, child := range item.Children {
                client.City.Create().
                    SetID(child.Adcode).
                    SetName(child.Name).
                    SetCode(child.Code).
                    SetOpen(false).
                    SetParent(parent).
                    SetLat(child.Lat).
                    SetLng(child.Lng).
                    SaveX(ctx)
            }
        }
    }
    _ = tx.Commit()
}
