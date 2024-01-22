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
	"github.com/golang-module/carbon/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/utils"
)

const (
	urlAccessToken      = "https://kyc.qcloud.com/api/oauth2/access_token"
	urlTicket           = "https://kyc.qcloud.com/api/oauth2/api_ticket"
	urlOcrResult        = "https://miniprogram-kyc.tencentcloudapi.com/api/server/getOcrResult"
	urlFaceId           = "https://miniprogram-kyc.tencentcloudapi.com/api/server/getfaceid"
	urlFaceVerifyResult = "https://miniprogram-kyc.tencentcloudapi.com/api/v2/base/queryfacerecord"

	cacheKeyAccessToken     = "AURORARIDE:WB_FACE:ACCESS_TOKEN"
	cacheKeySignTicket      = "AURORARIDE:WB_FACE:SIGN_TICKET"
	cacheKeyFaceVerifyTimes = "AURORARIDE:WB_FACE:TIMES:"
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
	Code            string `json:"code,omitempty"`
	Msg             string `json:"msg,omitempty"`
	TransactionTime string `json:"transactionTime,omitempty"`
	AccessToken     string `json:"access_token,omitempty"`
	ExpireTime      string `json:"expire_time,omitempty"`
	ExpireIn        int    `json:"expire_in,omitempty"`
}

type TicketRes struct {
	Code            string `json:"code,omitempty"`
	Msg             string `json:"msg,omitempty"`
	TransactionTime string `json:"transactionTime,omitempty"`
	Tickets         []struct {
		Value      string `json:"value,omitempty"`
		ExpireIn   int    `json:"expire_in,omitempty"`
		ExpireTime string `json:"expire_time,omitempty"`
	} `json:"tickets,omitempty"`
}

// OcrResultRes OCR结果
// 身份证OCR错误码: https://cloud.tencent.com/document/product/1007/47902
type OcrResultRes struct {
	Code            string    `json:"code,omitempty"`
	Msg             string    `json:"msg,omitempty"`
	BizSeqNo        string    `json:"bizSeqNo,omitempty"`
	Result          OcrResult `json:"result,omitempty"`
	TransactionTime string    `json:"transactionTime,omitempty"`
}

type OcrResult struct {
	FrontCode         string `json:"frontCode,omitempty"`         // 0 说明人像面识别成功
	BackCode          string `json:"backCode,omitempty"`          // 0 说明国徽面识别成功
	OrderNo           string `json:"orderNo,omitempty"`           // 订单编号
	Name              string `json:"name,omitempty"`              // frontCode 为 0 返回：证件姓名
	Sex               string `json:"sex,omitempty"`               // frontCode 为 0 返回：性别
	Nation            string `json:"nation,omitempty"`            // frontCode 为 0 返回：民族
	Birth             string `json:"birth,omitempty"`             // frontCode 为 0 返回：出生日期（例：19920320）
	Address           string `json:"address,omitempty"`           // frontCode 为 0 返回：地址
	Idcard            string `json:"idcard,omitempty"`            // frontCode 为 0 返回：身份证号
	ValidDate         string `json:"validDate,omitempty"`         // backCode 为 0 返回：证件的有效期（例：20160725-20260725）
	Authority         string `json:"authority,omitempty"`         // backCode 为 0 返回：发证机关
	FrontPhoto        string `json:"frontPhoto,omitempty"`        // 人像面照片，转换后为 JPG 格式
	BackPhoto         string `json:"backPhoto,omitempty"`         // 国徽面照片，转换后为 JPG 格式
	FrontCrop         string `json:"frontCrop,omitempty"`         // 人像面切边照片，切边图在识别原图少边或者存在遮挡的情况有小概率可能会导致切图失败，该字段会返回空；如切边图为空时建议可使用原图替代
	BackCrop          string `json:"backCrop,omitempty"`          // 国徽面切边照片，切边图在识别原图少边或者存在遮挡的情况有小概率可能会导致切图失败，该字段会返回空；如切边图为空时建议可使用原图替代
	HeadPhoto         string `json:"headPhoto,omitempty"`         // 身份证头像照片
	FrontWarnCode     string `json:"frontWarnCode,omitempty"`     // 人像面告警码，在身份证有遮挡、缺失、信息不全时会返回告警码；当 frontCode 为0时才会出现告警码，告警码的含义请参考 身份证 OCR 错误码
	BackWarnCode      string `json:"backWarnCode,omitempty"`      // 国徽面告警码，在身份证有遮挡、缺失、信息不全时会返回告警码；当 backCode 为0时才会出现告警码，告警码的含义请参考 身份证 OCR 错误码
	OperateTime       string `json:"operateTime,omitempty"`       // 做 OCR 的操作时间（例：2020-02-27 17:08:03）
	FrontMultiWarning string `json:"frontMultiWarning,omitempty"` // 正面多重告警码，含义请参考 身份证 OCR 错误码
	BackMultiWarning  string `json:"backMultiWarning,omitempty"`  // 反面多重告警码，含义请参考 身份证 OCR 错误码
	FrontClarity      string `json:"frontClarity,omitempty"`      // 正面图片清晰度
	BackClarity       string `json:"backClarity,omitempty"`       // 反面图片清晰度
	Success           bool   `json:"success,omitempty"`
	BizSeqNo          string `json:"bizSeqNo,omitempty"`
	TransactionTime   string `json:"transactionTime,omitempty"`
}

type FaceIdReq struct {
	OrderNo string `json:"orderNo,omitempty"` // 订单号，字母/数字组成的字符串，由合作方上传，每次唯一，不能超过32位
	Name    string `json:"name,omitempty"`    // 姓名
	IdNo    string `json:"idNo,omitempty"`    // 证件号码
	UserId  string `json:"userId,omitempty"`  // 用户 ID ，用户的唯一标识
}

type FaceIdRes struct {
	Code     string `json:"code,omitempty"`
	Msg      string `json:"msg,omitempty"`
	BizSeqNo string `json:"bizSeqNo,omitempty"`
	Result   struct {
		BizSeqNo        string `json:"bizSeqNo,omitempty"`
		TransactionTime string `json:"transactionTime,omitempty"`
		OrderNo         string `json:"orderNo,omitempty"`
		FaceId          string `json:"faceId,omitempty"`
		Success         bool   `json:"success,omitempty"`
	} `json:"result,omitempty"`
	TransactionTime string `json:"transactionTime,omitempty"`
}

type FaceVerifyRes struct {
	Code            string            `json:"code,omitempty"`
	Msg             string            `json:"msg,omitempty"`
	BizSeqNo        string            `json:"bizSeqNo,omitempty"`
	Result          *FaceVerifyResult `json:"result,omitempty"`
	TransactionTime string            `json:"transactionTime,omitempty"`
}

type FaceVerifyResult struct {
	OrderNo      string `json:"orderNo,omitempty"`
	LiveRate     string `json:"liveRate,omitempty"`
	Similarity   string `json:"similarity,omitempty"`
	OccurredTime string `json:"occurredTime,omitempty"`
	AppId        string `json:"appId,omitempty"`
	Photo        string `json:"photo,omitempty"`
	Video        string `json:"video,omitempty"`
	BizSeqNo     string `json:"bizSeqNo,omitempty"`
	SdkVersion   string `json:"sdkVersion,omitempty"`
	TrtcFlag     string `json:"trtcFlag,omitempty"`
}

func NewWbFace() *wbface {
	return _wbface
}

func (w *wbface) GetTimes(idCardNumber string) (times int) {
	times, _ = w.cache.Get(context.Background(), cacheKeyFaceVerifyTimes+idCardNumber).Int()
	return
}

func (w *wbface) UpdateTimes(idCardNumber string) {
	diff := time.Duration(carbon.Now().EndOfDay().DiffAbsInSeconds(carbon.Now())) * time.Second
	w.cache.Set(context.Background(), cacheKeyFaceVerifyTimes+idCardNumber, w.GetTimes(idCardNumber)+1, diff)
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

func (w *wbface) Licence() string {
	return w.licence
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

// SignTicket 获取SIGN ticket
// TODO: 当获取失败的时候刷新
func (w *wbface) SignTicket() string {
	_, ticket := w.GetCached()
	return ticket
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

// OcrResult 获取OCR结果
// https://cloud.tencent.com/document/product/1007/35864
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

// GetFaceId 上传用户信息获取人脸核身参数
// https://cloud.tencent.com/document/product/1007/35866
func (w *wbface) GetFaceId(req *FaceIdReq) (faceId, sign, nonce string, err error) {
	sign, nonce = w.Sign(w.SignTicket(), req.UserId)
	body := map[string]string{
		"appId":   w.appid,
		"orderNo": req.OrderNo,
		"name":    req.Name,
		"idNo":    req.IdNo,
		"userId":  req.UserId,
		"version": w.version,
		"sign":    sign,
		"nonce":   nonce,
	}
	var (
		r      *resty.Response
		result FaceIdRes
	)

	r, err = resty.New().R().
		SetBody(body).
		SetResult(&result).
		SetQueryParam("orderNo", req.OrderNo).
		Post(urlFaceId)

	zap.L().Info("Tencent - Face - GetFaceId", log.ResponseBody(r.Body()))

	if err != nil {
		return
	}

	if result.Code != "0" {
		return "", "", "", errors.New(result.Msg)
	}

	return result.Result.FaceId, sign, nonce, nil
}

// FaceVerifyResult 获取人脸核身结果
// https://cloud.tencent.com/document/product/1007/35880
func (w *wbface) FaceVerifyResult(orderNo string) (data *FaceVerifyResult, err error) {
	ticket := w.SignTicket()
	sign, nonce := w.Sign(ticket, orderNo)

	var (
		r      *resty.Response
		result FaceVerifyRes
	)

	r, err = resty.New().R().
		SetBody(map[string]string{
			"appId":   w.appid,
			"version": w.version,
			"nonce":   nonce,
			"orderNo": orderNo,
			"sign":    sign,
			"getFile": "1",
		}).
		SetResult(&result).
		SetQueryParam("orderNo", orderNo).
		Post(urlFaceVerifyResult)

	if err != nil {
		return
	}

	if result.Code != "0" {
		zap.L().Info("Tencent - Face - FaceVerifyResult", log.ResponseBody(r.Body()))
		err = errors.New(result.Msg)
		return
	}

	data = result.Result

	return
}
