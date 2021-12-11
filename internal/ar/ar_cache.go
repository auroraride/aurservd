// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import "github.com/go-redis/redis/v8"

var Cache *cache

type cache struct {
    *redis.Client
}

func NewCache() {
    cfg := Config.Database.Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.Addr,
        Password: cfg.Password,
        DB:       cfg.DB,
    })
    Cache = &cache{rdb}
}
