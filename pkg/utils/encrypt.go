// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha256"
    "encoding/base64"
)

// Md5String 字符串加密为md5
func Md5String(s string) string {
    hash := md5.New()
    hash.Write([]byte(s))
    md5Data := hash.Sum(nil)
    return base64.StdEncoding.EncodeToString(md5Data)
}

// Sha256Base64String 字符串加密为Sha256Base64
func Sha256Base64String(s string, secret string) string {
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(s))
    buf := h.Sum(nil)
    return base64.StdEncoding.EncodeToString(buf)
}
