// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package cache

import (
    "context"
    "github.com/go-redis/redis/v9"
    "time"
)

const (
    CabinetNameCacheKey = "__CABINET_NAMES__"
)

type SetFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd

var (
    client *redis.Client
    Set    func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
    Get    func(ctx context.Context, key string) *redis.StringCmd
    Del    func(ctx context.Context, keys ...string) *redis.IntCmd

    HSet    func(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
    HDel    func(ctx context.Context, key string, fields ...string) *redis.IntCmd
    HGet    func(ctx context.Context, key, field string) *redis.StringCmd
    HGetAll func(ctx context.Context, key string) *redis.MapStringStringCmd

    Rpush func(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
    BRPop func(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd

    Lpush func(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
    BLPop func(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd

    Expire func(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd

    TTL func(ctx context.Context, key string) *redis.DurationCmd
)

func CreateClient(c *redis.Client) {
    client = c

    Set = client.Set
    Get = client.Get
    Del = client.Del

    HSet = client.HSet
    HDel = client.HDel
    HGet = client.HGet
    HGetAll = client.HGetAll

    Rpush = client.RPush
    BRPop = client.BRPop

    Lpush = client.LPush
    BLPop = client.BLPop

    Expire = client.Expire

    TTL = client.TTL
}

func Float64(key string) float64 {
    var f float64
    saved, err := client.Get(context.Background(), key).Float64()
    if err == nil {
        f = saved
    }
    return f
}

func Int(key string) int {
    var res int
    saved, err := client.Get(context.Background(), key).Int()
    if err == nil {
        res = saved
    }
    return res
}
