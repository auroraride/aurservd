// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-01
// Based on aurservd by liasica, magicrolan@qq.com.

package mgo

import (
    "context"
    "github.com/qiniu/qmgo"
    "github.com/qiniu/qmgo/options"
    log "github.com/sirupsen/logrus"
)

var (
    Client        *qmgo.Client
    DB            *qmgo.Database
    CabinetTask   *qmgo.Collection
    EbikeAllocate *qmgo.Collection
)

func Connect(url, db string) {
    ctx := context.Background()
    var err error
    Client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: url})
    if err != nil {
        log.Fatalln(err)
    }
    DB = Client.Database(db)

    CabinetTask = DB.Collection("cabinet_task")
    err = CabinetTask.CreateIndexes(ctx, []options.IndexModel{
        {
            Key: []string{
                "serial",
            },
        },
        {
            Key: []string{
                "job",
            },
        },
        {
            Key: []string{
                "deactivated",
            },
        },
        {
            Key: []string{
                "cabinetId",
            },
        },
    })
    if err != nil {
        log.Fatalln(err)
    }

    EbikeAllocate = DB.Collection("ebike_allocate")
    err = EbikeAllocate.CreateIndexes(ctx, []options.IndexModel{
        {
            Key: []string{"ebike.id"},
        },
        {
            Key: []string{"rider.id"},
        },
        {
            Key: []string{"subscribeId"},
        },
        {
            Key: []string{"status"},
        },
        {
            Key: []string{"createdAt"},
        },
        {
            Key: []string{"employeeId"},
        },
        {
            Key: []string{"storeId"},
        },
        {
            Key: []string{"sn"},
        },
    })
}
