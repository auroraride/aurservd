// Copyright (C) liasica. 2022-present.
//
// Created at 2022-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package db

import (
    "context"
    "github.com/auroraride/aurservd/internal/ent"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "io/ioutil"
)

func insertCity(client *ent.Client) {
    type R struct {
        ID       uint64 `json:"id"`
        Adcode   string `json:"adcode,omitempty"`
        Name     string `json:"name,omitempty"`
        Code     string `json:"code,omitempty"`
        ParentId uint64 `json:"parent_id,omitempty"`
        Open     bool   `json:"open,omitempty"`
        Children []R    `json:"children,omitempty"`
    }

    ctx := context.Background()
    tx, _ := client.Tx(ctx)
    // 导入城市
    b, err := ioutil.ReadFile("city.json")
    log.Println("导入城市")
    if err == nil {
        var items []R
        err = jsoniter.Unmarshal(b, &items)
        if err == nil {
            for _, item := range items {
                parent := client.City.Create().
                    SetName(item.Name).
                    SetAdcode(item.Adcode).
                    SetCode(item.Code).
                    SaveX(ctx)
                for _, child := range item.Children {
                    client.City.Create().
                        SetName(child.Name).
                        SetAdcode(child.Adcode).
                        SetCode(child.Code).
                        SetOpen(false).
                        SetParent(parent).
                        SaveX(ctx)
                }
            }
        }
    }
    _ = tx.Commit()
}
