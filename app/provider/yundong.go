// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
    log "github.com/sirupsen/logrus"
    "strconv"
    "time"
)

// TODO 维护一个全局的token
type yundong struct {
    logger          *Logger
    url             string
    appid           string
    appkey          string
    tokenRetryTimes int // token获取重试次数
    retryTimes      int
}

// YDRes 云动通用返回
type YDRes struct {
    Code int    `json:"code,omitempty"`
    Msg  string `json:"msg,omitempty"`
}

func (p *yundong) Logger() *Logger {
    return p.logger
}

type Yundongurl string

const (
    yundongTokenKey = "YUNDONGTOKEN"

    yundongTokenUrl   Yundongurl = "/token"
    yundongStatusUrl  Yundongurl = "/cabinet/status"
    yundongControlUrl Yundongurl = "/cabinet/control"
)

func NewYundong() *yundong {
    cfg := ar.Config.ThirdParty.Yundong

    // 获取所有项目
    return &yundong{
        logger: NewLogger("yundong"),
        appid:  cfg.Appid,
        appkey: cfg.Appkey,
        url:    cfg.Url,
    }
}

func (p *yundong) PrepareRequest() {
    p.tokenRetryTimes = 0
    p.retryTimes = 0
}

func (p *yundong) GetUrl(path Yundongurl) string {
    return p.url + string(path)
}

// YDTokenRes token请求返回
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
    // 如果需要刷新token则删除缓存token
    if p.retryTimes > 0 {
        ar.Cache.Del(context.Background(), yundongTokenKey)
    }
    token = ar.Cache.Get(context.Background(), yundongTokenKey).Val()
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
            ar.Cache.Set(context.Background(), yundongTokenKey, token, time.Duration(int64(res.Expirets)-time.Now().Unix())*time.Second)
        }
    }
    return
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

type YDBin struct {
    Doorno       int     `json:"doorno,omitempty"`
    Doorstatus   int     `json:"doorstatus,omitempty"`
    Totalv       int     `json:"totalv,omitempty"`
    Chargei      int     `json:"chargei,omitempty"`
    Soc          float64 `json:"soc,omitempty"`
    HealthStatus int     `json:"healthStatus,omitempty"`
    IsEnable     bool    `json:"isEnable,omitempty"`
    BatterySN    string  `json:"batterySN,omitempty"`
    Putbattery   int     `json:"putbattery,omitempty"`
    Batterytype  string  `json:"batterytype,omitempty"`
}

type YDStatusRes struct {
    Code int    `json:"code,omitempty"`
    Msg  string `json:"msg,omitempty"`
    Data struct {
        CabinetSN        string  `json:"cabinetSN,omitempty"`
        Isonline         int     `json:"isonline,omitempty"`
        NumOfBattery     int     `json:"numOfBattery,omitempty"`
        Allowexchangenum int     `json:"allowexchangenum,omitempty"`
        Doorstatus       []YDBin `json:"doorstatus,omitempty"`
    } `json:"data,omitempty"`
}

func (r *YDStatusRes) String() string {
    return fmt.Sprintf("Code: %d, Message: %s", r.Code, r.Msg)
}

func (r *YDStatusRes) GetBins() (bins []YDBin) {
    bins = make([]YDBin, len(r.Data.Doorstatus))
    for _, ds := range r.Data.Doorstatus {
        bins[ds.Doorno] = ds
    }
    return
}

func (p *yundong) Cabinets() ([]*ent.Cabinet, error) {
    return ar.Ent.Cabinet.Query().Where(cabinet.Brand(model.CabinetBrandYundong.Value())).All(context.Background())
}

func (p *yundong) Brand() string {
    return "云动"
}

func (p *yundong) UpdateStatus(up *ent.CabinetUpdateOne, item *ent.Cabinet) any {
    res := new(YDStatusRes)
    _, err := p.RequestClient(false).
        SetResult(res).
        Get(p.GetUrl(yundongStatusUrl) + "?cabinetNo=" + item.Serial)
    // token 请求失败, 重新请求token后重试
    if res.Code == 1000 && p.retryTimes < 1 {
        p.retryTimes += 1
        return p.UpdateStatus(up, item)
    }
    if err != nil || res.Code != 0 {
        msg := fmt.Sprintf("云动状态获取失败, serial: %s, err: %#v, res: %s", item.Serial, err, res)
        log.Error(msg)
        return msg
    }

    // 仓位信息
    if res.Code == 0 {
        var full uint = 0
        var num uint = 0
        // 设置仓位状态
        bins := make([]model.CabinetBin, len(res.Data.Doorstatus))
        doors := res.GetBins()
        for index, ds := range doors {
            e := model.NewBatteryElectricity(utils.NewNumber().Decimal(ds.Soc))
            hasBattery := ds.Putbattery == 1
            if hasBattery {
                num += 1
            }
            errs := make([]string, 0)
            bin := model.CabinetBin{
                Name:        fmt.Sprintf("%d号仓", index+1),
                BatterySN:   ds.BatterySN,
                Full:        e.IsBatteryFull(),
                Battery:     hasBattery,
                Electricity: e,
                OpenStatus:  ds.Doorstatus == 1,
                DoorHealth:  ds.HealthStatus == 0,
                Current:     utils.NewNumber().Decimal(float64(ds.Chargei) / 1000),
                Voltage:     utils.NewNumber().Decimal(float64(ds.Totalv) / 1000),
            }
            if bin.Full {
                full += 1
            }
            if hasBattery && bin.Voltage == 0 && bin.Current == 0 && bin.Electricity == 0 {
                errs = append(errs, model.CabinetBinBatteryFault)
            }
            bin.ChargerErrors = errs
            bins[index] = bin
            if len(item.Bin) > index {
                bins[index].Remark = item.Bin[index].Remark
            }
        }
        up.SetBin(bins).
            SetBatteryNum(num).
            SetBatteryFullNum(full).
            SetHealth(uint(res.Data.Isonline)).
            SetDoors(uint(len(doors)))
    }
    return res
}

// DoorOperate 云动柜门操作
func (p *yundong) DoorOperate(user, serial, operation string, door int) (state bool) {
    type params struct {
        Doorno []int `json:"doorno"`
    }
    type body struct {
        CabinetSN  string `json:"cabinetSN"`
        EmployCode string `json:"employCode"`
        Params     params `json:"params"`
        Action     string `json:"action"`
    }

    res := new(YDRes)
    _, err := p.RequestClient(false).
        SetResult(res).
        SetBody(body{
            CabinetSN:  serial,
            EmployCode: user,
            Params:     params{Doorno: []int{door}},
            Action:     operation,
        }).
        Post(p.GetUrl(yundongControlUrl))
    // token 请求失败, 重新请求token后重试
    if res.Code == 1000 && p.retryTimes < 1 {
        p.retryTimes += 1
        return p.DoorOperate(user, serial, operation, door)
    }
    if err != nil {
        log.Error(err)
        return
    }
    return res.Code == 0
}
