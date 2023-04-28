// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Md5String hashes using md5 algorithm
func Md5String(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Md5Base64String 字符串加密为md5(base64)
func Md5Base64String(s string) string {
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

// HmacSha1Hexadecimal 获取Hmac-sha1 Hexadecimal加密字符串
func HmacSha1Hexadecimal(input, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
