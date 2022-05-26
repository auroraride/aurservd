// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package cache

import "context"

func BatteryFull(key string) float64 {
    f := float64(80)
    saved, err := client.Get(context.Background(), key).Float64()
    if err == nil {
        f = saved
    }
    return f
}
