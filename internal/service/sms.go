// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/utils"
)

type sms struct {
    phone string
}

func NewSms(phone string) *sms {
    return &sms{
        phone: phone,
    }
}

// SendCode 发送验证码
func (s *sms) SendCode() (string, error) {
    c, err := ali.NewSmsClient()
    if err != nil {
        return "", err
    }
    cfg := ar.Config.Aliyun.Sms
    code := fmt.Sprintf("%06d", utils.RandomIntMaxMin(1000, 999999))
    return c.SetTmplate(cfg.Template.General.Code).SetParam(map[string]string{"code": code}).SendCode(s.phone)
}
