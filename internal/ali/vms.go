// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/16
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
    openapi "github.com/alibabacloud-go/darabonba-openapi/client"
    dyvmsapi20170525 "github.com/alibabacloud-go/dyvmsapi-20170525/v2/client"
    "github.com/alibabacloud-go/tea/tea"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/snag"
    log "github.com/sirupsen/logrus"
)

type vms struct {
    client *dyvmsapi20170525.Client
}

func NewVms() *vms {
    cfg := ar.Config.Aliyun.Vms
    config := &openapi.Config{
        AccessKeyId:     &cfg.AccessKeyId,
        AccessKeySecret: &cfg.AccessKeySecret,
        Endpoint:        &cfg.Endpoint,
    }
    client, err := dyvmsapi20170525.NewClient(config)
    if err != nil {
        snag.Panic(err)
    }
    return &vms{client}
}

// SendVoiceMessageByTts 发送语音通知
func (v *vms) SendVoiceMessageByTts(phone, param string, tel, tmplate *string) {
    req := &dyvmsapi20170525.SingleCallByTtsRequest{
        CalledShowNumber: tel,
        TtsCode:          tmplate,
        CalledNumber:     &phone,
        TtsParam:         &param,
        Speed:            tea.Int32(-5),
    }
    res, err := v.client.SingleCallByTts(req)
    if err != nil {
        snag.Panic(err)
    }
    log.Infof("%s, %s 发送语音通知结果: %v", phone, param, res)
}
