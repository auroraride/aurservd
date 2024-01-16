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

func (w *wbface) AppId() string {
	return w.appid
}

type AccessTokenRes struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	TransactionTime string `json:"transactionTime"`
	AccessToken     string `json:"access_token"`
	ExpireTime      string `json:"expire_time"`
	ExpireIn        int    `json:"expire_in"`
}

func (w *wbface) AccessToken(params ...bool) (token string, err error) {
	var result AccessTokenRes

	ctx := context.Background()

	var force bool
	if len(params) > 0 {
		force = params[0]
	}

	if !force {
		token = w.cache.Get(ctx, cacheKeyAccessToken).Val()
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

		token = result.AccessToken

		w.cache.Set(ctx, cacheKeyAccessToken, token, time.Second*time.Duration(result.ExpireIn-60))
	}

	return
}

type TicketRes struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	TransactionTime string `json:"transactionTime"`
	Tickets         []struct {
		Value      string `json:"value"`
		ExpireIn   int    `json:"expire_in"`
		ExpireTime string `json:"expire_time"`
	} `json:"tickets"`
}

func (w *wbface) Ticket(userId string, params ...bool) (ticket string, err error) {
	var result TicketRes

	accessToken, _ := w.AccessToken(params...)

	_, err = resty.New().R().
		SetQueryParams(map[string]string{
			"app_id":       w.appid,
			"type":         "NONCE",
			"version":      "1.0.0",
			"user_id":      userId,
			"access_token": accessToken,
		}).
		SetResult(&result).
		Get(urlTicket)

	retry := err != nil || result.Code == "15" || result.Code == "400101"
	if retry && len(params) == 0 {
		return w.Ticket(userId, true)
	}

	if len(result.Tickets) > 0 {
		ticket = result.Tickets[0].Value
	}

	if result.Code != "0" && err == nil {
		err = errors.New(result.Msg)
	}
	return
}
