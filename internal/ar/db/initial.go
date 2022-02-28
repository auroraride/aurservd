// Copyright (C) liasica. 2022-present.
//
// Created at 2022-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package db

import "github.com/auroraride/aurservd/internal/ent"

type Options struct {
    City bool
}

func Database(client *ent.Client, o *Options) {
    if o.City {
        insertCities(client)
    }
}