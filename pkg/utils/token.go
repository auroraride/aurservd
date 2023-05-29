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

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func NewEcdsaToken() string {
	curve := elliptic.P256()
	key, _ := ecdsa.GenerateKey(curve, rand.Reader)
	priv := ecies.ImportECDSA(key)
	pub := priv.PublicKey
	b, _ := x509.MarshalECPrivateKey(priv.ExportECDSA())
	token := b64.RawURLEncoding.EncodeToString(elliptic.MarshalCompressed(elliptic.P256(), pub.X, pub.Y)) + "/" + b64.RawURLEncoding.EncodeToString(b)
	return token
}
