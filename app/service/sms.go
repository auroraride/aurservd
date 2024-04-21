// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
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
	t, _ := cache.Get(context.Background(), phone).Int64()
	if t-time.Now().Unix() > 0 {
		return "", errors.New("请勿频繁请求")
	}
	c, err := ali.NewSms()
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", utils.RandIntMaxMin(1000, 999999))
	var id string
	// log.Info(*c.Endpoint)
	// id = "test"
	cfg := ar.Config.Aliyun.Sms
	id, err = c.SetTemplate(cfg.Template.General).SetParam(map[string]string{"code": code}).SendCode(phone)
	if err != nil {
		return "", err
	}
	// 保存下次请求有效期
	cache.Set(context.Background(), phone, time.Now().Unix()+59, 1*time.Minute)
	cache.Set(context.Background(), id, code, 5*time.Minute)

	zap.L().Info("SMS", zap.String("phone", phone), zap.String("code", code))
	return id, nil
}

// VerifyCode 校验短信验证码
func (s *sms) VerifyCode(phone, id, code string) bool {
	if _, ok := ar.Config.App.Debug.Phone[phone]; ok {
		return true
	}

	ctx := context.Background()
	isValid := cache.Get(ctx, id).Val() == code
	if isValid {
		cache.Del(ctx, id)
	}

	return isValid
}

func (s *sms) VerifyCodeX(phone, id, code string) {
	if s.VerifyCode(phone, id, code) {
		return
	}
	snag.Panic("短信验证码校验失败")
}

// SendSignSuccess 发送合同通知
func (s *sms) SendSignSuccess(t time.Time, name, hash, phone string) (id string, err error) {
	c, err := ali.NewSms()
	if err != nil {
		return "", err
	}
	cfg := ar.Config.Aliyun.Sms
	id, err = c.SetTemplate(cfg.Template.ContractSuccess).
		SetParam(map[string]string{
			"time": t.Format("2006-01-02 15:04:05"),
			"name": name,
			"hash": hash,
		}).SendCode(phone)
	if err != nil {
		return "", err
	}
	return id, nil
}
