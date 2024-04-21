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
	"unsafe"
)

// RandIntMaxMin 获取范围内随机整数
func RandIntMaxMin(min, max int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(max-min) + min
}

// RandTokenString 获取随机字符串
func RandTokenString() string {
	b := make([]byte, 512)
	_, _ = crand.Read(b)
	return b64.StdEncoding.EncodeToString(b)
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

// RandStr 生成随机字符串
// 参考: https://xie.infoq.cn/article/f274571178f1bbe6ff8d974f3
func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
