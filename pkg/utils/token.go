// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	b64 "encoding/base64"
)

func NewEcdsaToken() string {
	curve := elliptic.P256()
	key, _ := ecdsa.GenerateKey(curve, rand.Reader)
	b, _ := x509.MarshalECPrivateKey(key)
	pub := key.PublicKey
	token := b64.RawURLEncoding.EncodeToString(elliptic.MarshalCompressed(pub.Curve, pub.X, pub.Y)) + "/" + b64.RawURLEncoding.EncodeToString(b)
	return token
}
