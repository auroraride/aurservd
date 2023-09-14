// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-16
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	sls "github.com/aliyun/aliyun-log-go-sdk"

	"github.com/auroraride/aurservd/internal/ar"
)

type slsClient struct {
	sls.ClientInterface
	project string
}

func NewSls() *slsClient {
	cfg := ar.Config.Aliyun.Sls
	config := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyId),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
	}
	// 访问的域名
	config.Endpoint = tea.String(cfg.Endpoint)
	result := sls.CreateNormalInterfaceV2(cfg.Endpoint, sls.NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.AccessKeySecret, ""))

	client := &slsClient{
		ClientInterface: result,
		project:         cfg.Project,
	}

	return client
}
