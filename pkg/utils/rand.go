// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
	crand "crypto/rand"
	b64 "encoding/base64"
	"math/rand"
	"time"
)

// RandomIntMaxMin 获取范围内随机整数
func RandomIntMaxMin(min, max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min
}

// RandTokenString 获取随机字符串
func RandTokenString() string {
	b := make([]byte, 512)
	_, _ = crand.Read(b)
	return b64.StdEncoding.EncodeToString(b)
}
