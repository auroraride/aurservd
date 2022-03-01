// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package db

import (
    "context"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/utils"
    log "github.com/sirupsen/logrus"
)

func insertManager(client *ent.Client) {
    log.Println("插入超级管理员")
    password, _ := utils.PasswordGenerate("AuroraAdmin@2022#!")
    client.Manager.Create().
        SetName("超级管理员").
        SetPhone("18888888888").
        SetPassword(password).
        ExecX(context.Background())
}
