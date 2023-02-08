// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import "github.com/go-redis/redis/v9"

var (
    Quit                           chan bool
    Redis                          *redis.Client
    CabinetNameCacheKey            string
    RiderExchangeTimeLimitCacheKey string
)

func init() {
    Quit = make(chan bool)
}
