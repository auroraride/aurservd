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
	urlOcrResult   = "https://miniprogram-kyc.tencentcloudapi.com/api/server/getOcrResult"

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

// OcrResultRes OCR结果
// https://cloud.tencent.com/document/product/1007/35864
// 身份证OCR错误码: https://cloud.tencent.com/document/product/1007/47902
type OcrResultRes struct {
	Code            string    `json:"code"`
	Msg             string    `json:"msg"`
	BizSeqNo        string    `json:"bizSeqNo"`
	Result          OcrResult `json:"result"`
	TransactionTime string    `json:"transactionTime"`
}

type OcrResult struct {
	FrontCode         string `json:"frontCode"`         // 0 说明人像面识别成功
	BackCode          string `json:"backCode"`          // 0 说明国徽面识别成功
	OrderNo           string `json:"orderNo"`           // 订单编号
	Name              string `json:"name"`              // frontCode 为 0 返回：证件姓名
	Sex               string `json:"sex"`               // frontCode 为 0 返回：性别
	Nation            string `json:"nation"`            // frontCode 为 0 返回：民族
	Birth             string `json:"birth"`             // frontCode 为 0 返回：出生日期（例：19920320）
	Address           string `json:"address"`           // frontCode 为 0 返回：地址
	Idcard            string `json:"idcard"`            // frontCode 为 0 返回：身份证号
	ValidDate         string `json:"validDate"`         // backCode 为 0 返回：证件的有效期（例：20160725-20260725）
	Authority         string `json:"authority"`         // backCode 为 0 返回：发证机关
	FrontPhoto        string `json:"frontPhoto"`        // 人像面照片，转换后为 JPG 格式
	BackPhoto         string `json:"backPhoto"`         // 国徽面照片，转换后为 JPG 格式
	FrontCrop         string `json:"frontCrop"`         // 人像面切边照片，切边图在识别原图少边或者存在遮挡的情况有小概率可能会导致切图失败，该字段会返回空；如切边图为空时建议可使用原图替代
	BackCrop          string `json:"backCrop"`          // 国徽面切边照片，切边图在识别原图少边或者存在遮挡的情况有小概率可能会导致切图失败，该字段会返回空；如切边图为空时建议可使用原图替代
	HeadPhoto         string `json:"headPhoto"`         // 身份证头像照片
	FrontWarnCode     string `json:"frontWarnCode"`     // 人像面告警码，在身份证有遮挡、缺失、信息不全时会返回告警码；当 frontCode 为0时才会出现告警码，告警码的含义请参考 身份证 OCR 错误码
	BackWarnCode      string `json:"backWarnCode"`      // 国徽面告警码，在身份证有遮挡、缺失、信息不全时会返回告警码；当 backCode 为0时才会出现告警码，告警码的含义请参考 身份证 OCR 错误码
	OperateTime       string `json:"operateTime"`       // 做 OCR 的操作时间（例：2020-02-27 17:08:03）
	FrontMultiWarning string `json:"frontMultiWarning"` // 正面多重告警码，含义请参考 身份证 OCR 错误码
	BackMultiWarning  string `json:"backMultiWarning"`  // 反面多重告警码，含义请参考 身份证 OCR 错误码
	FrontClarity      string `json:"frontClarity"`      // 正面图片清晰度
	BackClarity       string `json:"backClarity"`       // 反面图片清晰度
	Success           bool   `json:"success"`
	BizSeqNo          string `json:"bizSeqNo"`
	TransactionTime   string `json:"transactionTime"`
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

	zap.L().Info("Tencent - Face - RefreshAccessToken", log.ResponseBody(r.Body()))

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

	zap.L().Info("Tencent - Face - RefreshSignTicket", log.ResponseBody(r.Body()))

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

	zap.L().Info("Tencent - Face - NonceTicket", log.ResponseBody(r.Body()))

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

func (w *wbface) Sign(ticket string, params ...string) (sign string, nonce string) {
	nonce = w.nonce()
	appid := w.appid
	list := []string{appid, w.version, ticket, nonce}
	list = append(list, params...)
	sort.Strings(list)

	s := strings.Join(list, "")

	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	sign = strings.ToUpper(hex.EncodeToString(bs))
	return
}

func (w *wbface) OcrResult(orderNo string) (err error, ocrResult *OcrResult) {
	ticket := w.SignTicket()
	sign, nonce := w.Sign(ticket, orderNo)

	var (
		r      *resty.Response
		result OcrResultRes
	)

	params := map[string]string{
		"app_id":   w.appid,
		"order_no": orderNo,
		"get_file": "1",
		"nonce":    nonce,
		"version":  w.version,
		"sign":     sign,
	}

	r, err = resty.New().R().
		SetQueryParams(params).
		SetResult(&result).
		Get(urlOcrResult)

	if err != nil {
		return
	}

	if result.Code != "0" {
		zap.L().Info("Tencent - Face - OcrResult", log.ResponseBody(r.Body()))
		err = errors.New(result.Msg)
		return
	}

	ocrResult = &result.Result

	// 判定是否成功获取
	if ocrResult.FrontCode != "0" || ocrResult.BackCode != "0" {
		err = errors.New("身份证识别失败，请重试")
		return
	}

	// 判定人像面切图
	if ocrResult.FrontCrop == "" {
		ocrResult.FrontCrop = ocrResult.FrontPhoto
	}

	// 判定国徽面切图
	if ocrResult.BackCrop == "" {
		ocrResult.BackCode = ocrResult.BackPhoto
	}

	return
}
