// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
	cloudauth "github.com/alibabacloud-go/cloudauth-20190307/v3/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"

	"github.com/auroraride/aurservd/internal/ar"
)

type Cloudauth struct {
	*cloudauth.Client
}

func NewCloudauth() (ca *Cloudauth, err error) {
	cfg := ar.Config.Aliyun.Cloudauth

	var credential credentials.Credential
	credential, err = credentials.NewCredential(nil)
	if err != nil {
		return
	}

	// 初始化Client。
	config := &openapi.Config{
		Credential: credential,
		Endpoint:   tea.String(cfg.Endpoint),
	}

	var client *cloudauth.Client
	client, err = cloudauth.NewClient(config)
	if err != nil {
		return
	}

	return &Cloudauth{client}, nil
}
