// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-15
// Based on aurservd by liasica, magicrolan@qq.com.

package tencent

import (
	"context"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"

	"github.com/auroraride/aurservd/internal/ar"
)

const (
	urlAccessToken = "https://kyc.qcloud.com/api/oauth2/access_token"
	urlTicket      = "https://kyc.qcloud.com/api/oauth2/api_ticket"

	cacheKeyAccessToken = "AURORARIDE_WB_FACE_ACCESS_TOKEN_KEY"
)

type wbface struct {
	appid   string
	secret  string
	licence string
	cache   *redis.Client
}

func NewWbFace(cache *redis.Client) *wbface {
	cfg := ar.Config.WbFace
	return &wbface{
		appid:   cfg.AppId,
		secret:  cfg.Secret,
		licence: cfg.Licence,
		cache:   cache,
	}
}

type AccessTokenRes struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	TransactionTime string `json:"transactionTime"`
	AccessToken     string `json:"access_token"`
	ExpireTime      string `json:"expire_time"`
	ExpireIn        string `json:"expire_in"`
}

func (w *wbface) AccessToken(params ...bool) (token string, err error) {
	var result AccessTokenRes

	ctx := context.Background()

	if len(params) == 0 || !params[0] {
		token = w.cache.Get(ctx, cacheKeyAccessToken).String()
	}

	if token == "" {
		_, err = resty.New().R().
			SetQueryParams(map[string]string{
				"app_id":     w.appid,
				"secret":     w.secret,
				"grant_type": "client_credential",
				"version":    "1.0.0",
			}).
			SetResult(&result).
			Get(urlAccessToken)

		if result.Code != "0" {
			msg := result.Msg
			if msg == "" {
				msg = "access_token 请求失败"
			}
			err = errors.New(msg)
		}

		if err != nil {
			return
		}

		w.cache.Set(ctx, cacheKeyAccessToken, token, time.Minute*15)
	}

	return
}
