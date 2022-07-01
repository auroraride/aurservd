// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "net/url"
    "strconv"
    "time"
)

type yundong struct {
    logger            *Logger
    url               string
    appid             string
    appkey            string
    tokenRetryTimes   int // token获取重试次数
    retryTimes        int
    operateRetryTimes int // 获取操作日志判定操作是否成功重试次数 30次为上限
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

    yundongTokenUrl    Yundongurl = "/token"
    yundongStatusUrl   Yundongurl = "/cabinet/status"
    yundongControlUrl  Yundongurl = "/cabinet/control"
    yundongOperatedUrl Yundongurl = "/cabinet/operated"
    yundongOperatorlog Yundongurl = "/cabinet/operatorlog"
    yundongBasicinfo   Yundongurl = "/zhangfei/cabinet/basicinfo"
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
    p.operateRetryTimes = 0
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
        cache.Del(context.Background(), yundongTokenKey)
    }
    token = cache.Get(context.Background(), yundongTokenKey).Val()
    if token == "" {
        client := p.RequestClient(true)
        res := new(YDTokenRes)
        r, err := client.SetResult(res).Post(p.GetUrl(yundongTokenUrl))

        log.Info(string(r.Body()))
        if err != nil || res.Code != 0 {
            p.tokenRetryTimes += 1
            if p.tokenRetryTimes < 2 {
                return p.FetchToken(true)
            }
        } else {
            token = res.Token
            cache.Set(context.Background(), yundongTokenKey, token, time.Duration(int64(res.Expirets)-time.Now().Unix())*time.Second)
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
    return ent.Database.Cabinet.Query().Where(cabinet.Brand(model.CabinetBrandYundong.Value()), cabinet.Status(model.CabinetStatusNormal)).All(context.Background())
}

func (p *yundong) Brand() string {
    return "云动"
}

func (p *yundong) UpdateStatus(item *ent.Cabinet, params ...any) error {
    res := new(YDStatusRes)
    _, err := p.RequestClient(false).
        SetResult(res).
        Get(p.GetUrl(yundongStatusUrl) + "?cabinetNo=" + item.Serial)
    // token 请求失败, 重新请求token后重试
    if res.Code == 1000 && p.retryTimes < 1 {
        p.retryTimes += 1
        return p.UpdateStatus(item)
    }

    up := item.Update().SetHealth(model.CabinetHealthStatusOffline)

    oldHealth := item.Health
    oldBins := item.Bin
    oldNum := item.BatteryNum

    defer func(up *ent.CabinetUpdateOne, ctx context.Context) {
        v, _ := up.Save(ctx)
        *item = *v
        monitor(oldBins, oldHealth, oldNum, item)
    }(up, context.Background())

    // log.Infof("云动状态获取结果：%s", string(r.Body()))
    if err != nil {
        p.logger.Write(fmt.Sprintf("云动状态获取失败, serial: %s, err: %s, res: %s\n", item.Serial, err.Error(), res))
        return err
    }
    p.logger.Write(res)

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
            isFull := e.IsBatteryFull()
            voltage := utils.NewNumber().Decimal(float64(ds.Totalv) / 1000)
            current := utils.NewNumber().Decimal(float64(ds.Chargei) / 1000)
            if hasBattery {
                num += 1
            }
            if isFull {
                full += 1
            }
            errs := make([]string, 0)
            if hasBattery && voltage == 0 && current == 0 && e == 0 {
                errs = append(errs, model.CabinetBinBatteryFault)
            }

            bin := model.CabinetBin{
                Index:         index,
                Name:          fmt.Sprintf("%d号仓", index+1),
                BatterySN:     ds.BatterySN,
                Full:          isFull,
                Battery:       hasBattery,
                Electricity:   e,
                OpenStatus:    ds.Doorstatus == 1,
                DoorHealth:    ds.IsEnable,
                Current:       current,
                Voltage:       voltage,
                ChargerErrors: errs,
            }

            if len(item.Bin) > index {
                bin.Remark = item.Bin[index].Remark
            }

            bins[index] = bin
        }

        // 判断是否处于换电过程中, 如果处于换电过程中则不保存电池数量, 以避免电池变动数量大的情况出现
        if len(params) > 0 {
            switch params[0].(type) {
            case bool:
                if !params[0].(bool) {
                    up.SetBatteryNum(num)
                }
                break
            }
        } else {
            up.SetBatteryNum(num)
        }

        up.SetBin(bins).
            SetBatteryFullNum(full).
            SetHealth(uint8(res.Data.Isonline)).
            SetDoors(uint(len(doors)))
        return nil
    }
    return errors.New("云动状态获取失败")
}

// DoorOperate 云动柜门操作
// user携带操作ID，比对操作日志实时获取状态
func (p *yundong) DoorOperate(code, serial, operation string, door int) (state bool) {
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
    now := time.Now()

    r, err := p.RequestClient(false).
        SetResult(res).
        SetBody(body{
            CabinetSN:  serial,
            EmployCode: code,
            Params:     params{Doorno: []int{door}},
            Action:     operation,
        }).
        Post(p.GetUrl(yundongControlUrl))
    if t, err := time.Parse(time.RFC1123, r.Header().Get("Date")); err == nil {
        now = t.Add(8 * time.Hour)
    }
    log.Info(string(r.Body()))
    // token 请求失败, 重新请求token后重试
    if res.Code == 1000 && p.retryTimes < 1 {
        p.retryTimes += 1
        return p.DoorOperate(code, serial, operation, door)
    }
    if err != nil {
        log.Error(err)
        return
    }
    start := now.Add(-10 * time.Second).Format(carbon.DateTimeLayout)
    end := now.Add(10 * time.Second).Format(carbon.DateTimeLayout)

    return res.Code == 0 && p.GetOperateState(code, serial, start, end)
}

func (p *yundong) Reboot(code string, serial string) (state bool) {
    type body struct {
        CabinetSN  string `json:"cabinetSN"`
        EmployCode string `json:"employCode"`
        Action     string `json:"action"`
    }

    res := new(YDRes)
    r, err := p.RequestClient(false).
        SetResult(res).
        SetBody(body{
            CabinetSN:  serial,
            EmployCode: code,
            Action:     "rebootCabinet",
        }).
        Post(p.GetUrl(yundongOperatedUrl))
    log.Info(string(r.Body()))
    // token 请求失败, 重新请求token后重试
    if res.Code == 1000 && p.retryTimes < 1 {
        p.retryTimes += 1
        return p.Reboot(code, serial)
    }
    if err != nil {
        log.Error(err)
        return
    }
    return res.Code == 0
}

type OperatorlogData struct {
    Actiontime time.Time `json:"actiontime"`
    EmployCode string    `json:"employCode"`
    Params     struct {
        Doorno []interface{} `json:"doorno,omitempty"`
    } `json:"params"`
    Action string `json:"action"`
    Result string `json:"result"`
}

type OperatorlogRes struct {
    Code  int               `json:"code"`
    Msg   string            `json:"msg"`
    Total int               `json:"total"`
    Data  []OperatorlogData `json:"data"`
}

// GetOperateState 获取操作日志判定操作是否成功
func (p *yundong) GetOperateState(opId, serial, start, end string) (state bool) {
    res := new(OperatorlogRes)
    r, err := p.RequestClient(false).
        SetResult(res).
        Get(fmt.Sprintf("%s?cabinetSN=%s&starttime=%s&endtime=%s&pageNo=0&pageNum=50", p.GetUrl(yundongOperatorlog), url.QueryEscape(serial), url.QueryEscape(start), url.QueryEscape(end)))
    log.Info(string(r.Body()))
    // token 请求失败, 重新请求token后重试
    if res.Code == 1000 && p.retryTimes < 1 {
        p.retryTimes += 1
        return p.GetOperateState(opId, serial, start, end)
    }
    if err != nil {
        log.Error(err)
        return
    }
    for _, d := range res.Data {
        if d.EmployCode == opId && d.Result == "succ" {
            return true
        }
    }
    // 重复30次查询
    if !state {
        time.Sleep(2 * time.Second)
        p.operateRetryTimes += 1
        return p.GetOperateState(opId, serial, start, end)
    }
    return
}

// UpdateBasicInfo 更新云动电柜信息(投产)
func (p *yundong) UpdateBasicInfo(req model.YundongDeployInfo) {
    type info struct {
        Type        string                  `json:"type"`
        CabinetSN   string                  `json:"cabinetSN"`
        AgentId     string                  `json:"agentId"`
        WarehouseId string                  `json:"warehouseId"`
        AreaCode    string                  `json:"areaCode"`
        UpdateInfo  model.YundongDeployInfo `json:"updateInfo"`
    }

    data := info{
        Type:        "normal",
        CabinetSN:   req.SN,
        AgentId:     "90",
        WarehouseId: "10",
        AreaCode:    req.AreaCode,
        UpdateInfo:  req,
    }

    res := new(YDRes)
    r, err := p.RequestClient(false).
        SetResult(res).
        SetBody(data).
        Post(p.GetUrl(yundongBasicinfo))

    log.Info(string(r.Body()))

    if err != nil {
        snag.Panic(err)
    }
}
