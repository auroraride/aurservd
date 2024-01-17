// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-15
// Based on aurservd by liasica, magicrolan@qq.com.

package tencent

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/utils"
)

const (
	urlAccessToken = "https://kyc.qcloud.com/api/oauth2/access_token"
	urlTicket      = "https://kyc.qcloud.com/api/oauth2/api_ticket"

	cacheKeyAccessToken = "AURORARIDE:WB_FACE:ACCESS_TOKEN"
	cacheKeySignTicket  = "AURORARIDE:WB_FACE:SIGN_TICKET"
)

// 缓存时间
var cacheDuration = time.Minute * 20

var _wbface *wbface

type wbface struct {
	appid   string
	secret  string
	licence string
	cache   *redis.Client
	version string
}

type AccessTokenRes struct {
	Code            string `json:"code"`
	Msg             string `json:"msg"`
	TransactionTime string `json:"transactionTime"`
	AccessToken     string `json:"access_token"`
	ExpireTime      string `json:"expire_time"`
	ExpireIn        int    `json:"expire_in"`
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

func NewWbFace() *wbface {
	return _wbface
}

func BootWbFace(cache *redis.Client) {
	if _wbface != nil {
		return
	}
	cfg := ar.Config.WbFace
	_wbface = &wbface{
		appid:   cfg.AppId,
		secret:  cfg.Secret,
		licence: cfg.Licence,
		cache:   cache,
		version: "1.0.0",
	}

	go func() {
		// 判定是否需要刷新
		if cache.Get(context.Background(), cacheKeyAccessToken).Val() == "" {
			_ = _wbface.Refresh()
		}

		// 定时刷新
		ticker := time.NewTicker(cacheDuration - time.Minute)
		for range ticker.C {
			_ = _wbface.Refresh()
		}
	}()
}

func (w *wbface) AppId() string {
	return w.appid
}

func (w *wbface) Version() string {
	return w.version
}

func (w *wbface) nonce() string {
	return utils.RandStr(32)
}

func (w *wbface) Refresh() (err error) {
	var token string
	token, err = w.RefreshAccessToken()
	if err != nil {
		return
	}

	return w.RefreshSignTicket(token)
}

func (w *wbface) ClearCached() {
	w.cache.Del(context.Background(), cacheKeyAccessToken)
	w.cache.Del(context.Background(), cacheKeySignTicket)
}

func (w *wbface) GetCached(retried ...bool) (token string, ticket string) {
	token = w.cache.Get(context.Background(), cacheKeyAccessToken).Val()
	ticket = w.cache.Get(context.Background(), cacheKeySignTicket).Val()

	if token == "" && len(retried) == 0 {
		_ = w.Refresh()
		return w.GetCached(true)
	}
	if ticket == "" {
		_ = w.RefreshSignTicket(token)
	}
	return
}

func (w *wbface) RefreshAccessToken() (token string, err error) {
	var (
		result AccessTokenRes
		r      *resty.Response
	)

	ctx := context.Background()
	r, err = resty.New().R().
		SetQueryParams(map[string]string{
			"app_id":     w.appid,
			"secret":     w.secret,
			"grant_type": "client_credential",
			"version":    "1.0.0",
		}).
		SetResult(&result).
		Get(urlAccessToken)

	zap.L().Info("Tencent - RefreshAccessToken", log.ResponseBody(r.Body()))

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

	w.cache.Set(ctx, cacheKeyAccessToken, token, cacheDuration)
	return
}

func (w *wbface) RefreshSignTicket(token string) (err error) {
	var (
		r      *resty.Response
		result TicketRes
	)

	r, err = resty.New().R().
		SetQueryParams(map[string]string{
			"app_id":       w.appid,
			"type":         "SIGN",
			"version":      "1.0.0",
			"access_token": token,
		}).
		SetResult(&result).
		Get(urlTicket)

	zap.L().Info("Tencent - RefreshSignTicket", log.ResponseBody(r.Body()))

	if result.Code != "0" && err == nil {
		err = errors.New(result.Msg)
	}

	if err != nil {
		return
	}

	if len(result.Tickets) > 0 {
		w.cache.Set(context.Background(), cacheKeySignTicket, result.Tickets[0].Value, cacheDuration)
	}
	return
}

func (w *wbface) NonceTicket(userId string, retried ...bool) (ticket string, err error) {
	token, _ := w.GetCached()

	var (
		r      *resty.Response
		result TicketRes
	)
	r, err = resty.New().R().
		SetQueryParams(map[string]string{
			"app_id":       w.appid,
			"type":         "NONCE",
			"version":      "1.0.0",
			"user_id":      userId,
			"access_token": token,
		}).
		SetResult(&result).
		Get(urlTicket)

	zap.L().Info("Tencent - NonceTicket", log.ResponseBody(r.Body()))

	retry := err != nil || result.Code == "15" || result.Code == "400101"
	if retry && len(retried) == 0 {
		w.ClearCached()
		return w.NonceTicket(userId, true)
	}

	if len(result.Tickets) > 0 {
		ticket = result.Tickets[0].Value
	}

	if result.Code != "0" && err == nil {
		err = errors.New(result.Msg)
	}
	return
}

func (w *wbface) SignTicket() string {
	return w.cache.Get(context.Background(), cacheKeySignTicket).Val()
}

func (w *wbface) Sign(userId, ticket string) (sign string, nonce string) {
	nonce = w.nonce()
	appid := w.appid
	list := []string{appid, userId, w.version, ticket, nonce}
	sort.Strings(list)

	s := strings.Join(list, "")

	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	sign = strings.ToUpper(hex.EncodeToString(bs))
	return
}
