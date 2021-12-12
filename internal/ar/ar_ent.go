// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import "github.com/auroraride/aurservd/internal/ent"

var Ent *ent.Client

type orm struct {
    *ent.Client
}

func OpenDatabase() *orm {
    Ent = ent.OpenPgx(Config.Database.Postgres.Dsn)
    return &orm{Ent}
}
