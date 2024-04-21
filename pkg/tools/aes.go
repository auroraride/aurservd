// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var (
	ErrInvalidBlockSize    = errors.New("invalid blocksize")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data (empty or not padded)")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type AesCrypto struct {
	iv  []byte
	key []byte
}

func NewAesCrypto(iv, key []byte) *AesCrypto {
	return &AesCrypto{
		iv:  iv,
		key: key,
	}
}

// CBCEncrypt CBC加密
func (t *AesCrypto) CBCEncrypt(plaintext []byte) (ciphertext []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(t.key)
	if err != nil {
		return
	}

	plaintext, _ = pkcs7Pad(plaintext, block.BlockSize())

	// 加密字符串
	cbc := cipher.NewCBCEncrypter(block, t.iv)
	ciphertext = make([]byte, len(plaintext))
	cbc.CryptBlocks(ciphertext, plaintext)
	return
}

// CBCDecrypt CBC解密
func (t *AesCrypto) CBCDecrypt(ciphertext []byte) (plaintext []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(t.key)
	if err != nil {
		return
	}

	// 解密字符串
	cbc := cipher.NewCBCDecrypter(block, t.iv)
	plaintext = make([]byte, len(ciphertext))
	cbc.CryptBlocks(plaintext, ciphertext)
	plaintext, _ = pkcs7Unpad(plaintext, block.BlockSize())
	return
}

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

// CFBEncrypt CFB加密
func (t *AesCrypto) CFBEncrypt(b []byte) (s string, err error) {
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

// CFBDecrypt CFB解密
func (t *AesCrypto) CFBDecrypt(b []byte) (s string, err error) {
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
