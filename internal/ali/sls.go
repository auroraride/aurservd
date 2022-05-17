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
    log "github.com/sirupsen/logrus"
    "time"
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
    result := sls.CreateNormalInterface(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret, "")

    client := &slsClient{
        ClientInterface: result,
        project:         cfg.Project,
    }

    return client
}

func (c *slsClient) PutLog() {
    logs := []*sls.Log{
        {
            Time: tea.Uint32(uint32(time.Now().Unix())),
            Contents: []*sls.LogContent{
                {
                    Key:   tea.String("name"),
                    Value: tea.String("1号仓"),
                },
                {
                    Key:   tea.String("battery"),
                    Value: tea.String("true"),
                },
            },
        },
    }
    err := c.PutLogs(c.project, "cabinet-dev", &sls.LogGroup{
        Logs: logs,
        Category: tea.String("YUNDONG"),
        Topic: tea.String("NCAWDFA0L75N027"),
    })
    if err != nil {
        log.Error(err)
        return
    }
}
