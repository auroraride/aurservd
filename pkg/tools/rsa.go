// Copyright (C) liasica. 2024-present.
//
// Created at 2024-03-02
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
)

type Rsa struct {
	key *rsa.PrivateKey
}

func NewRsa(der []byte) (*Rsa, error) {
	private, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return nil, err
	}
	return &Rsa{key: private}, nil
}

func (r *Rsa) PublicKey() *rsa.PublicKey {
	return &r.key.PublicKey
}

func (r *Rsa) PrivateKey() *rsa.PrivateKey {
	return r.key
}

func (r *Rsa) Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha1.New(), rand.Reader, r.PublicKey(), data, nil)
}

func (r *Rsa) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha1.New(), rand.Reader, r.PrivateKey(), data, nil)
}
