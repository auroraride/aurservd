// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
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

// EncryptAES 加密函数
func EncryptAES(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充原文以符合块大小
	blockSize := block.BlockSize()
	originalData := []byte(plaintext)
	paddedData := pad(originalData, blockSize)

	// 初始化向量IV必须是唯一，但不需要保密
	ciphertext := make([]byte, blockSize+len(paddedData))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], paddedData)

	// 返回十六进制编码的字符串
	return hex.EncodeToString(ciphertext), nil
}

// PKCS#7填充
func pad(buf []byte, blockSize int) []byte {
	padding := blockSize - len(buf)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(buf, padText...)
}

// DecryptAES 解密函数
func DecryptAES(key []byte, encryptedHex string) (string, error) {
	encrypted, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 获取IV
	blockSize := block.BlockSize()
	iv := encrypted[:blockSize]
	encrypted = encrypted[blockSize:]

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encrypted, encrypted)

	// 去除填充
	unpaddedData, err := unpad(encrypted, blockSize)
	if err != nil {
		return "", err
	}

	return string(unpaddedData), nil
}

// 去除PKCS#7填充
func unpad(buf []byte, blockSize int) ([]byte, error) {
	length := len(buf)
	padLen := int(buf[length-1])
	if padLen > blockSize || padLen > length {
		return nil, errors.New("invalid padding")
	}
	return buf[:length-padLen], nil
}
