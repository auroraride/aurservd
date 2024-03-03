// Copyright (C) liasica. 2024-present.
//
// Created at 2024-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAesCBC(t *testing.T) {
	key, _ := base64.StdEncoding.DecodeString("I8+6cGLzBVliTnuvyc5MgrBFtjHGXCOuemF7bF68dBE=")
	iv, _ := base64.StdEncoding.DecodeString("0NQ1fPijzlMew44RtFZ6jA==")
	message := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit")

	c := NewAesCrypto(iv, key)
	encrypted, err := c.CBCEncrypt(message)
	require.NoError(t, err)

	var decrypted []byte
	decrypted, err = c.CBCDecrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, message, decrypted)

	encrypted, _ = base64.StdEncoding.DecodeString("nblj0msZP1O2mUXX+WGdgwEZXeYRSO2HVRgxJVnPua9jD/LPwLrSK8YyR0YHM9uaCGgWeLTHNbGeNJQ6Vr5zkw==")
	decrypted, err = c.CBCDecrypt(encrypted)
	require.NoError(t, err)
	require.Equal(t, message, decrypted)
}
