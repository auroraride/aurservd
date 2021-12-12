// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func RandomIntMaxMin(min, max int) int {
    return rand.Intn(max - min) + min
}
