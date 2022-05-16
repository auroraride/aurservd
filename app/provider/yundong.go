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
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
    log "github.com/sirupsen/logrus"
    "strconv"
    "sync"
    "time"
)

type yundong struct {
    logger          *Logger
    url             string
    appid           string
    appkey          string
    tokenRetryTimes int // token获取重试次数
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
    }
}

func (p *yundong) PrepareRequest() {
    p.tokenRetryTimes = 0
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

func (p *yundong) UpdateStatus() {
    log.Info("开始轮询获取云动电柜状态")
    start := time.Now()

    var wg sync.WaitGroup

    items, err := ar.Ent.Cabinet.Query().Where(cabinet.Brand(model.CabinetBrandYundong.String())).All(context.Background())
    if err != nil {
        log.Errorf("云动电柜查询失败: %#v", err)
        return
    }

    length := len(items)
    wg.Add(length)
    logs := make([]*YDStatusRes, length)

    go func() {
        for i, item := range items {
            up := ar.Ent.Cabinet.UpdateOne(item)
            res := new(YDStatusRes)
            _, err = p.RequestClient(false).
                SetResult(res).
                Get(p.GetUrl(yundongStatusUrl))
            if err != nil || res.Code != 0 {
                log.Errorf("云动状态获取失败, serial: %s, err: %#v, res: %s", item.Serial, err, res)
                goto Save
            }
            logs[i] = res

            // 仓位信息
            if res.Code == 0 {
                var full uint = 0
                var num uint = 0
                // 设置仓位状态
                bins := make([]model.CabinetBin, len(res.Data.Doorstatus))
                doors := res.GetBins()
                for index, ds := range doors {
                    e := model.NewBatteryElectricity(ds.Soc)
                    if e.IsBatteryFull() {
                        full += 1
                    }
                    hasBattery := ds.Putbattery == 1
                    if hasBattery {
                        num += 1
                    }
                    errs := make([]string, 0)
                    bin := model.CabinetBin{
                        Name:          fmt.Sprintf("%d号仓", index+1),
                        BatterySN:     ds.BatterySN,
                        Full:          e.IsBatteryFull(),
                        Battery:       hasBattery,
                        Electricity:   e,
                        OpenStatus:    ds.Doorstatus == 1,
                        DoorHealth:    ds.HealthStatus == 0,
                        Current:       float64(ds.Chargei) / 1000,
                        Voltage:       float64(ds.Totalv) / 1000,
                    }
                    if bin.Voltage == 0 && hasBattery {
                        errs = append(errs, "有电池无电压")
                    }
                    bin.ChargerErrors = errs
                    bins[index] = bin
                    if len(item.Bin) > index {
                        bins[index].Remark = item.Bin[index].Remark
                        bins[index].Locked = item.Bin[index].Locked
                    }
                }
                up.SetHealth(uint(res.Data.Isonline)).SetBin(bins).SetBatteryNum(num).SetBin(bins)
            } else {
                // 未获取到电柜状态设置为离线
                up.SetHealth(model.CabinetHealthStatusOffline)
            }

            Save:
            // 存储电柜信息
            up.SaveX(context.Background())

            wg.Done()
        }
    }()

    wg.Wait()

    // 写入电柜日志
    p.logger.Write(logs)

    log.Infof("云动电柜状态轮询完成, 耗时%.2fs", time.Now().Sub(start).Seconds())
}
