// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import (
    "bytes"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "strconv"
    "time"
)

const (
    resSuccess = 0
    resExists  = 53000000
    methodPost = "POST"
    methodGet  = "GET"
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
)

type Esign struct {
    config        ar.EsignConfig
    headers       map[string]string
    client        *resty.Request
    serialization jsoniter.Config
}

type commonRes struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
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
        config:        config,
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
    return utils.Sha256Base64String(str, e.config.Secret), str
}

type reqLog struct {
    Api    string
    Method string
    Body   interface{}
    Res    interface{}
    MD5    string
    Secret string
    Sign   string
    Raw    string
}

// request 请求
func (e *Esign) request(api, method string, body interface{}, data interface{}) interface{} {
    var md5 string
    res := new(commonRes)
    res.Data = data
    bodyString, _ := e.serialization.Froze().MarshalToString(body)
    md5 = utils.Md5String(bodyString)
    singnature, raw := e.getSign(api, method, md5)
    req := resty.New().
        R().
        SetResult(res).
        SetBody(bodyString).
        SetHeaders(e.headers).
        SetHeader("Content-MD5", md5).
        SetHeader("X-Tsign-Open-Ca-Signature", singnature)
    var err error
    switch method {
    case methodPost:
        _, err = req.SetBody(body).Post(e.config.BaseUrl + api)
    case methodGet:
        _, err = req.Get(e.config.BaseUrl + api)
    }
    if err != nil {
        snag.Panic(err)
    }
    // 记录请求日志
    if e.config.Log {
        logdata := reqLog{
            Api:    api,
            Method: method,
            Body:   bodyString,
            Res:    res,
            MD5:    md5,
            Secret: e.config.Secret,
            Sign:   singnature,
            Raw:    raw,
        }
        logstr, _ := e.serialization.Froze().MarshalToString(logdata)
        log.Info(logstr)
    }
    e.isResSuccess(res)
    return res.Data
}

// isResSuccess 返回是否成功
func (e *Esign) isResSuccess(res *commonRes) {
    switch res.Code {
    case resSuccess, resExists:
        return
    }
    snag.Panic(res.Message)
}
