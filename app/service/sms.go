// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/utils"
    "time"
)

type sms struct {
    debug map[string]struct{}
}

func NewSms() *sms {
    return &sms{
        debug: map[string]struct{}{
            "00000000": {},
        },
    }
}

// SendCode 发送验证码
func (s *sms) SendCode(phone string) (string, error) {
    c, err := ali.NewSmsClient()
    if err != nil {
        return "", err
    }
    cfg := ar.Config.Aliyun.Sms
    code := fmt.Sprintf("%06d", utils.RandomIntMaxMin(1000, 999999))
    var id string
    id, err = c.SetTmplate(cfg.Template.General.Code).SetParam(map[string]string{"code": code}).SendCode(phone)
    if err != nil {
        return "", err
    }
    ar.Cache.Set(context.Background(), id, code, 5*time.Minute)
    return id, nil
}

// VerifyCode 校验短信验证码
func (s *sms) VerifyCode(id, code string) bool {
    ctx := context.Background()
    isValid := ar.Cache.Get(ctx, id).Val() == code
    if isValid {
        ar.Cache.Del(ctx, id)
    }
    return isValid
}
