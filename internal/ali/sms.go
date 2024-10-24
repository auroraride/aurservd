// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package ali

import (
	"errors"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ar"
)

type smsClient struct {
	*dysmsapi.Client

	sign string
	tmpl string
	data string
}

// NewSms 创建sms客户端
func NewSms() (c *smsClient, err error) {
	cfg := ar.Config.Aliyun.Sms
	config := &openapi.Config{
		AccessKeyId:     &cfg.AccessKeyId,
		AccessKeySecret: &cfg.AccessKeySecret,
		Endpoint:        &cfg.Endpoint,
	}
	var client *dysmsapi.Client
	client, err = dysmsapi.NewClient(config)
	if err != nil {
		zap.L().Error("阿里云sms初始化失败", zap.Error(err))
		return
	}
	c = &smsClient{
		Client: client,
		sign:   cfg.Sign,
	}
	return
}

// SetTemplate 设置模板
func (c *smsClient) SetTemplate(tmpl string) *smsClient {
	c.tmpl = tmpl
	return c
}

// SetParam 设置参数
func (c *smsClient) SetParam(items map[string]string) *smsClient {
	b, _ := jsoniter.Marshal(items)
	c.data = string(b)
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
	zap.L().Info(phone + " -> 短信发送结果: " + res.String())
	if err != nil {
		return
	}
	if res == nil || res.Body == nil {
		return "", errors.New("短信发送失败")
	}
	if strings.ToUpper(*res.Body.Code) != "OK" {
		return "", errors.New(*res.Body.Message)
	}
	id = *res.Body.BizId
	return
}
