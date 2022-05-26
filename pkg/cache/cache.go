// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package cache

import (
    "context"
    "github.com/go-redis/redis/v8"
    "time"
)

type SetFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd

var (
    client *redis.Client
    Set    func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
    Get    func(ctx context.Context, key string) *redis.StringCmd
    Del    func(ctx context.Context, keys ...string) *redis.IntCmd
)

func CreateClient(addr, password string, db int) {
    client = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    Set = client.Set
    Get = client.Get
    Del = client.Del
}
