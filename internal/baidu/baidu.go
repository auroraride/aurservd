// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package baidu

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/go-resty/resty/v2"
    "time"
)

const (
    accessTokenKey = "BAIDU_ACCESS_TOKEN"

    // accessTokenUrl 公共access_token接口URL
    accessTokenUrl = `https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s`
)

type baiduClient struct {
    apiKey      string
    secretKey   string
    accessToken string
    failTimes   int
}

// accessTokenResp 公共access_token返回体
type accessTokenResp struct {
    Error            string `json:"error,omitempty"`
    ErrorDescription string `json:"error_description,omitempty"`
    RefreshToken     string `json:"refresh_token,omitempty"`
    ExpiresIn        int    `json:"expires_in,omitempty"`
    Scope            string `json:"scope,omitempty"`
    SessionKey       string `json:"session_key,omitempty"`
    AccessToken      string `json:"access_token,omitempty"`
    SessionSecret    string `json:"session_secret,omitempty"`
}

// New 初始化百度客户端
func New() *baiduClient {
    cfg := ar.Config.Baidu.Face
    b := &baiduClient{
        apiKey:    cfg.ApiKey,
        secretKey: cfg.SecretKey,
    }
    b.getAccessToken()
    return b
}

// requestAccessToken 从服务器请求百度 access_token
func (b *baiduClient) requestAccessToken() {
    var err error
    res := new(accessTokenResp)
    url := fmt.Sprintf(accessTokenUrl, b.apiKey, b.secretKey)
    _, err = resty.New().R().
        SetResult(res).
        Post(url)
    if err != nil {
        panic(response.NewError(err))
    }
    if res.Error != "" {
        panic(response.NewError(res.ErrorDescription))
    }
    ar.Cache.Set(context.Background(), accessTokenKey, res.AccessToken, time.Second*time.Duration(res.ExpiresIn-120))
    b.accessToken = res.AccessToken
    return
}

// getAccessToken 获取 access_token
func (b *baiduClient) getAccessToken() {
    t := ar.Cache.Get(context.Background(), accessTokenKey).Val()
    if t == "" {
        b.requestAccessToken()
    }
    b.accessToken = t
}
