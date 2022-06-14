// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
)

type aescrypto struct {
    iv  []byte
    key []byte
}

func NewAESCrypto() *aescrypto {
    return &aescrypto{
        iv:  []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
        key: []byte("OSEriONseMISicolADhIcrinGtOPmoSp"),
    }
}

func (t *aescrypto) Encrypt(b []byte) (s string, err error) {
    // 创建加密算法aes
    c, err := aes.NewCipher(t.key)
    if err != nil {
        return "", err
    }
    // 加密字符串
    cfb := cipher.NewCFBEncrypter(c, t.iv)
    cb := make([]byte, len(b))
    cfb.XORKeyStream(cb, b)
    s = base64.StdEncoding.EncodeToString(cb)
    return
}

func (t *aescrypto) Decrypt(b []byte) (s string, err error) {
    // 创建加密算法aes
    c, err := aes.NewCipher(t.key)
    if err != nil {
        return "", err
    }
    // 解密字符串
    cfbdec := cipher.NewCFBDecrypter(c, t.iv)
    plainByte := make([]byte, len(b))
    cfbdec.XORKeyStream(plainByte, b)
    s = string(plainByte)
    return
}
