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
    "go.uber.org/zap"
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

func (v *vms) log(phone, param *string, res *dyvmsapi20170525.SingleCallByTtsResponse, err error) {
    p := "<nil>"
    pa := "<nil>"
    if phone != nil {
        p = *phone
    }
    if param != nil {
        pa = *param
    }
    zap.L().Info(p+", "+pa+" -> 语音发送结果: "+res.String(), zap.Error(err))
    return
}

// SendVoiceMessageByTts 发送语音通知
func (v *vms) SendVoiceMessageByTts(phone, param *string, tel, tmplate *string) bool {
    req := &dyvmsapi20170525.SingleCallByTtsRequest{
        CalledShowNumber: tel,
        TtsCode:          tmplate,
        CalledNumber:     phone,
        TtsParam:         param,
        Speed:            tea.Int32(-200),
    }
    res, err := v.client.SingleCallByTts(req)
    v.log(phone, param, res, err)
    return err == nil && *res.Body.Code == "OK"
}
