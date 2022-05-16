// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
    log "github.com/sirupsen/logrus"
    "strconv"
    "time"
)

var (
    Yundong *yundong
)

type yundong struct {
    logger          *Logger
    url             string
    appid           string
    appkey          string
    tokenRetryTimes int // token获取重试次数
    items           []*ent.Cabinet
}

type Yundongurl string

const (
    yundongTokenKey = "YUNDONGTOKEN"

    yundongTokenUrl  Yundongurl = "/token"
    yundongStatusUrl Yundongurl = "/cabinet/status"
)

func NewYundong() *yundong {
    cfg := ar.Config.ThirdParty.Yundong
    // 获取所有项目
    return &yundong{
        logger: NewLogger("yundong"),
        appid:  cfg.Appid,
        appkey: cfg.Appkey,
        url:    cfg.Url,
        items:  ar.Ent.Cabinet.Query().AllX(context.Background()),
    }
}

func (p *yundong) GetUrl(path Yundongurl) string {
    return p.url + string(path)
}

// YDTokenRes token 请求返回
type YDTokenRes struct {
    Code     int    `json:"code,omitempty"`
    Msg      string `json:"msg,omitempty"`
    Expirets int    `json:"expirets,omitempty"`
    Token    string `json:"token,omitempty"`
}

// FetchToken 获取token
func (p *yundong) FetchToken(tokenRequest bool) (token string) {
    if tokenRequest {
        return
    }
    token = ar.Cache.Get(context.Background(), yundongTokenKey).String()
    if token == "" {
        r := p.RequestClient(true)
        res := new(YDTokenRes)
        _, err := r.SetResult(res).Post(p.GetUrl(yundongTokenUrl))
        if err != nil || res.Code != 0 {
            p.tokenRetryTimes += 1
            if p.tokenRetryTimes < 2 {
                return p.FetchToken(true)
            }
        } else {
            token = res.Token
        }
    }
    return
}

func (p *yundong) PrepareRequest() {
    p.tokenRetryTimes = 0
}

func (p *yundong) RequestClient(tokenRequest bool) *resty.Request {
    ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
    token := ""
    token = p.FetchToken(tokenRequest)
    if !tokenRequest && token == "" {
        return nil
    }
    r := resty.New().SetTimeout(20*time.Second).R().
        SetHeader("appid", p.appid).
        SetHeader("ts", ts).
        SetHeader("auth", utils.HmacSha1Hexadecimal(fmt.Sprintf("%s%s%s", p.appid, token, ts), p.appkey)).
        SetHeader("token", token)
    return r
}

type YDDoorstatus struct {
    Doorno       int    `json:"doorno,omitempty"`
    Doorstatus   int    `json:"doorstatus,omitempty"`
    Totalv       int    `json:"totalv,omitempty"`
    Chargei      int    `json:"chargei,omitempty"`
    Soc          int    `json:"soc,omitempty"`
    HealthStatus int    `json:"healthStatus,omitempty"`
    IsEnable     bool   `json:"isEnable,omitempty"`
    BatterySN    string `json:"batterySN,omitempty"`
    Putbattery   int    `json:"putbattery,omitempty"`
    Batterytype  string `json:"batterytype,omitempty"`
}

type YDStatusRes struct {
    Code int    `json:"code,omitempty"`
    Msg  string `json:"msg,omitempty"`
    Data struct {
        CabinetSN        string         `json:"cabinetSN,omitempty"`
        Isonline         int            `json:"isonline,omitempty"`
        NumOfBattery     int            `json:"numOfBattery,omitempty"`
        Allowexchangenum int            `json:"allowexchangenum,omitempty"`
        Doorstatus       []YDDoorstatus `json:"doorstatus,omitempty"`
    } `json:"data,omitempty"`
}

func (r YDStatusRes) String() string {
    return fmt.Sprintf("Code: %d, Message: %s", r.Code, r.Msg)
}

func (p *yundong) UpdateStatus() {
    res := new(YDStatusRes)
    _, err := p.RequestClient(false).
        SetResult(res).
        Get(p.GetUrl(yundongStatusUrl))
    if err != nil || res.Code != 0 {
        log.Errorf("云动状态获取失败, err: %s, res: %s", err, res)
        return
    }
    // 更新电柜状态
    p.logger.Write(res)
}
