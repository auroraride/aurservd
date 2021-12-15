// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import (
    "bytes"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
    jsoniter "github.com/json-iterator/go"
    "strconv"
    "time"
)

const (
    resSuccess = 0
    resExists  = 53000000
)

const (
    EnvSandbox = "sandbox"
    EnvOnline  = "online"
)

const (
    headerContentType = "application/json;charset=UTF-8"
    headerAccept      = "*/*"
    headerAuthMode    = "Signature"
)

const (
    // createPersonAccountUrl 创建个人签署账号
    createPersonAccountUrl = `/v1/accounts/createByThirdPartyUserId`

    // docTemplateUrl 查询模板文件详情
    docTemplateUrl = `/v1/docTemplates/%s`
)

type esign struct {
    config  ar.EsignConfig
    headers map[string]string
    client  *resty.Request
}

type commonRes struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func New() *esign {
    var config ar.EsignConfig
    cfg := ar.Config.Esign
    switch cfg.Target {
    case EnvSandbox:
        config = cfg.Sandbox
    case EnvOnline:
        config = cfg.Online
    default:
        panic(response.NewError("环境设置失败"))
    }
    return &esign{
        config: config,
        headers: map[string]string{
            "X-Tsign-Open-App-Id":       config.Appid,
            "Content-Type":              headerContentType,
            "Accept":                    headerAccept,
            "X-Tsign-Open-Auth-Mode":    headerAuthMode,
            "X-Tsign-Open-Ca-Timestamp": strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
        },
    }
}

// getSign 获取签名
func (e *esign) getSign(api, method, md5 string) string {
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
    return utils.Sha256Base64String(str, e.config.Secret)
}

// request 请求
func (e *esign) request(api, method string, body interface{}, data interface{}) interface{} {
    res := new(commonRes)
    res.Data = data
    s, _ := jsoniter.MarshalToString(body)
    md5 := utils.Md5String(s)
    req := resty.New().
        R().
        SetResult(res).
        SetBody(body).
        SetHeaders(e.headers).
        SetHeader("Content-MD5", md5).
        SetHeader("X-Tsign-Open-Ca-Signature", e.getSign(api, method, md5))
    var err error
    switch method {
    case "POST":
        _, err = req.SetBody(body).Post(e.config.BaseUrl + api)
    case "GET":
        _, err = req.Get(e.config.BaseUrl + api)
    }
    if err != nil {
        panic(response.NewError(err))
    }
    e.isResSuccess(res)
    return res.Data
}

// isResSuccess 返回是否成功
func (e *esign) isResSuccess(res *commonRes) {
    switch res.Code {
    case resSuccess, resExists:
        return
    }
    panic(response.NewError(res.Message))
}
