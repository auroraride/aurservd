// Copyright (C) liasica. 2024-present.
//
// Created at 2024-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import (
	"encoding/pem"

	"go.uber.org/zap"

	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/pkg/tools"
)

var (
	personRsa *tools.Rsa
)

func PersonRsa() *tools.Rsa {
	if personRsa == nil {
		LoadRsa()
	}

	return personRsa
}

func LoadRsa() {
	privBlock, _ := pem.Decode(assets.PersonPrivateKey)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		zap.L().Fatal("无效的私钥")
	}

	var err error
	personRsa, err = tools.NewRsa(privBlock.Bytes)
	if err != nil {
		zap.L().Fatal("RSA初始化失败", zap.Error(err))
	}
}
