// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
    openapi "github.com/alibabacloud-go/darabonba-openapi/client"
    dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
    "github.com/auroraride/aurservd/internal/ar"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
)

type smsClient struct {
    *dysmsapi.Client

    sign string
    tmpl string
    data string
}

// NewSmsClient 创建sms客户端
func NewSmsClient() (c *smsClient, err error) {
    cfg := ar.Config.Aliyun.Sms
    config := &openapi.Config{
        AccessKeyId:     &cfg.AccessId,
        AccessKeySecret: &cfg.AccessSecret,
        Endpoint:        &cfg.Endpoint,
    }
    var client *dysmsapi.Client
    client, err = dysmsapi.NewClient(config)
    if err != nil {
        log.Errorf("阿里云sms初始化失败: %v", err)
        return
    }
    c = &smsClient{
        Client: client,
        sign:   cfg.Sign,
    }
    return
}

// SetTmplate 设置模板
func (c *smsClient) SetTmplate(tmpl string) *smsClient {
    c.tmpl = tmpl
    return c
}

// SetParam 设置参数
func (c *smsClient) SetParam(items map[string]string) *smsClient {
    c.data, _ = jsoniter.MarshalToString(items)
    return c
}

// SendCode 发送短信验证码
func (c *smsClient) SendCode(phone string) (id string, err error) {
    req := &dysmsapi.SendSmsRequest{
        PhoneNumbers:  &phone,
        SignName:      &c.sign,
        TemplateCode:  &c.tmpl,
        TemplateParam: &c.data,
    }
    res, err := c.SendSms(req)
    if err != nil {
        return
    }
    id = *res.Body.BizId
    return
}
