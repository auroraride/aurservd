// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/manager"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
)

func DatabaseInitial() {
    cityInitial()
    managerInitial()
}

func managerInitial() {
    client := ent.Database.Manager
    p := "18888888888"

    if e, _ := client.QueryNotDeleted().Where(manager.Phone(p)).Exist(context.Background()); e {
        return
    }

    password, _ := utils.PasswordGenerate("AuroraAdmin@2022#!")
    client.Create().
        SetName("超级管理员").
        SetPhone(p).
        SetPassword(password).
        ExecX(context.Background())
}

func cityInitial() {
    type R struct {
        Adcode   uint64  `json:"adcode"`
        Name     string  `json:"name"`
        Code     string  `json:"code"`
        Lng      float64 `json:"lng,omitempty"`
        Lat      float64 `json:"lat,omitempty"`
        Children []R     `json:"children,omitempty"`
    }

    ctx := context.Background()

    if c, _ := ent.Database.City.Query().Count(context.Background()); c > 0 {
        return
    }

    tx, _ := ent.Database.Tx(ctx)
    // 导入城市
    log.Println("导入城市")
    var items []R
    err := jsoniter.Unmarshal(assets.City, &items)
    if err == nil {
        for _, item := range items {
            parent, err := tx.City.Create().
                SetID(item.Adcode).
                SetName(item.Name).
                SetCode(item.Code).
                Save(ctx)
            snag.PanicIfErrorX(err, tx.Rollback)
            for _, child := range item.Children {
                _, err = tx.City.Create().
                    SetID(child.Adcode).
                    SetName(child.Name).
                    SetCode(child.Code).
                    SetOpen(false).
                    SetParent(parent).
                    SetLat(child.Lat).
                    SetLng(child.Lng).
                    Save(ctx)
                snag.PanicIfErrorX(err, tx.Rollback)
            }
        }
    }
    _ = tx.Commit()
}
