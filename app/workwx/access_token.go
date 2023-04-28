// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

import (
	"context"
	"fmt"
	"time"

	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/go-resty/resty/v2"
)

const (
	cacheKeyAccessToken = "WORKWX_ACCESS_TOKEN"
)

type tokenResponse struct {
	baseResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// getAccessToken 获取accesstoken
func (w *Client) getAccessToken(params ...bool) string {
	if len(params) == 0 || !params[0] {
		token := cache.Get(context.Background(), cacheKeyAccessToken).Val()
		if token != "" {
			return token
		}
	}

	var res tokenResponse
	_, err := resty.New().
		R().
		SetResult(&res).
		Get(fmt.Sprintf(`%s/gettoken?corpid=%s&corpsecret=%s`, baseURL, w.corpID, w.corpSecret))
	if err != nil {
		snag.Panic(err)
		return ""
	}
	if res.Errcode != 0 {
		snag.Panic(res.Errmsg)
	}
	cache.Set(context.Background(), cacheKeyAccessToken, res.AccessToken, time.Duration(res.ExpiresIn)*time.Second)
	return res.AccessToken
}
