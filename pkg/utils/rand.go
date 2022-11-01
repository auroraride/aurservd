// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
    b64 "encoding/base64"
    "math/rand"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

// RandomIntMaxMin 获取范围内随机整数
func RandomIntMaxMin(min, max int) int {
    return rand.Intn(max-min) + min
}

// RandTokenString 获取随机字符串
func RandTokenString() string {
    b := make([]byte, 128)
    rand.Read(b)
    return b64.RawURLEncoding.EncodeToString(b)
}
