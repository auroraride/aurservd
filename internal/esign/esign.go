// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import (
	"bytes"
	"strconv"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

const (
	resSuccess = 0
	resExists  = 53000000
	methodPost = "POST"
	methodGet  = "GET"
	methodPut  = "PUT"
)

const (
	EnvSandbox = "sandbox"
	EnvOnline  = "online"
)

const (
	headerContentType = "application/json; charset=UTF-8"
	headerAccept      = "*/*"
	headerAuthMode    = "Signature"
)

const (
	// createPersonAccountUrl 创建个人签署账号
	createPersonAccountUrl = `/v1/accounts/createByThirdPartyUserId`

	// docTemplateUrl 查询模板文件详情
	docTemplateUrl = `/v1/docTemplates/%s`

	// createByTemplateUrl 填充内容生成PDF
	createByTemplateUrl = `/v1/files/createByTemplate`

	// createFlowOneStepUrl 一步发起签署
	createFlowOneStepUrl = `/api/v2/signflows/createFlowOneStep`

	// executeUrl 获取签署链接
	executeUrl = `/v1/signflows/%s/executeUrl?accountId=%s&appScheme=%s&urlType=0`

	// signResultUrl 签署结果查询
	signResultUrl = `/v1/signflows/%s`

	// documentUrl 流程文档下载
	documentUrl = `/v1/signflows/%s/documents`

	// 修改个人签署账号
	accountUrl = `/v1/accounts/%s`
)

type Esign struct {
	Config        ar.EsignConfig
	headers       map[string]string
	redirect      string
	serialization jsoniter.Config
	sn            string
}

type commonRes struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Success *bool       `json:"success,omitempty"`
}

func New() *Esign {
	var config ar.EsignConfig
	cfg := ar.Config.Esign
	switch cfg.Target {
	case EnvSandbox:
		config = cfg.Sandbox
	case EnvOnline:
		config = cfg.Online
	default:
		snag.Panic("环境设置失败")
	}
	return &Esign{
		serialization: jsoniter.Config{SortMapKeys: true},
		Config:        config,
		headers: map[string]string{
			"X-Tsign-Open-App-Id":    config.Appid,
			"Content-Type":           headerContentType,
			"Accept":                 headerAccept,
			"X-Tsign-Open-Auth-Mode": headerAuthMode,
		},
	}
}

func (e *Esign) SetSn(sn string) {
	e.sn = sn
}

// getSign 获取签名
func (e *Esign) getSign(api, method, md5 string) (string, string) {
	var buffer bytes.Buffer
	buffer.WriteString(method)
	buffer.WriteString("\n")
	buffer.WriteString(headerAccept)
	buffer.WriteString("\n")
	buffer.WriteString(md5)
	buffer.WriteString("\n")
	buffer.WriteString(headerContentType)
	buffer.WriteString("\n")
	buffer.WriteString("")
	buffer.WriteString("\n")
	buffer.WriteString(api)
	str := buffer.String()
	return utils.Sha256Base64String(str, e.Config.Secret), str
}

type reqLog struct {
	Api      string      `json:"api,omitempty"`
	Method   string      `json:"method,omitempty"`
	Body     interface{} `json:"body,omitempty"`
	Res      interface{} `json:"res,omitempty"`
	MD5      string      `json:"MD5,omitempty"`
	Secret   string      `json:"secret,omitempty"`
	Sign     string      `json:"sign,omitempty"`
	Raw      string      `json:"raw,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

// request 请求
// TODO 优化请求
func (e *Esign) request(api, method string, body interface{}, data interface{}) interface{} {
	var (
		err        error
		r          *resty.Response
		md5        string
		bodyString string
		singnature string
	)
	res := new(commonRes)

	defer func() {
		// 记录请求日志
		logdata := reqLog{
			Api:    api,
			Method: method,
			Body:   bodyString,
			Res:    res,
			MD5:    md5,
			Secret: e.Config.Secret,
			Sign:   singnature,
			// Raw:      r.Body(),
			Response: string(r.Body()),
		}

		zap.L().Info("E签宝请求结果:", log.JsonData(logdata))
	}()

	res.Data = data
	if body != nil {
		bodyString, _ = e.serialization.Froze().MarshalToString(body)
		md5 = utils.Md5Base64String(bodyString)
	}
	singnature, _ = e.getSign(api, method, md5)
	req := resty.New().
		R().
		SetResult(res).
		SetBody(bodyString).
		SetHeaders(e.headers).
		SetHeader("Content-MD5", md5).
		SetHeader("X-Tsign-Open-Ca-Signature", singnature).
		SetHeader("X-Tsign-Open-Ca-Timestamp", strconv.FormatInt(time.Now().UnixNano()/1e6, 10))

	switch method {
	case methodPost:
		r, err = req.SetBody(body).Post(e.Config.BaseUrl + api)
	case methodPut:
		r, err = req.SetBody(body).Put(e.Config.BaseUrl + api)
	case methodGet:
		r, err = req.Get(e.Config.BaseUrl + api)
	}
	if err != nil {
		snag.Panic(err)
	}

	e.isResSuccess(r, res)
	return res.Data
}

// isResSuccess 返回是否成功
func (e *Esign) isResSuccess(r *resty.Response, res *commonRes) {
	if r.IsSuccess() {
		switch res.Code {
		case resSuccess, resExists:
			return
		}
	}
	snag.Panic(res.Message)
}
