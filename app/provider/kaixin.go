// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-16
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "fmt"
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/log"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/go-resty/resty/v2"
    jsoniter "github.com/json-iterator/go"
    "go.uber.org/zap"
    "regexp"
    "strconv"
    "strings"
)

type kaixin struct {
    key    string
    url    string
    user   string
    logger *Logger
    errors map[string]string
}

func (p *kaixin) Reboot(code string, serial string) bool {
    return false
}

func (p *kaixin) Cabinets() ([]*ent.Cabinet, error) {
    return ent.Database.Cabinet.QueryNotDeleted().
        Where(
            cabinet.Brand(adapter.CabinetBrandKaixin),
            cabinet.Status(model.CabinetStatusNormal.Value()),
            cabinet.Intelligent(false),
        ).
        All(context.Background())
}

func (p *kaixin) Brand() string {
    return "凯信"
}

func (p *kaixin) Logger() *Logger {
    return p.logger
}

type kaixinUrl string

const (
    kaixinUrlDetailData    kaixinUrl = "/batteryCupboard/detailData/"
    kaixinUrlDoorOperation kaixinUrl = "/batteryCupboard/DoorOperation/"
)

func NewKaixin() *kaixin {
    cfg := ar.Config.ThirdParty.Kaixin
    return &kaixin{
        key:    cfg.Key,
        url:    cfg.Url,
        user:   "AURORARIDE",
        logger: NewLogger("kaixin"),
        errors: map[string]string{
            "1":  "电池充电过慢",
            "2":  "电池充电过快",
            "3":  "220V 丢失/充电器损坏",
            "4":  "充电器状态错误",
            "5":  "电池未连接到充电器",
            "6":  "行程开关故障",
            "7":  "充电触点接触不良",
            "8":  "电池无法充满",
            "9":  "电池无法充电",
            "10": "充电器通讯故障",
            "11": "行程开关接触不良",
            "12": "已取出，未解绑",
        },
    }
}

func (p *kaixin) Detailcode(serial string) string {
    return utils.Md5String(utils.Md5String(p.user+serial) + p.key)
}

func (p *kaixin) GetUrl(path kaixinUrl) string {
    return p.url + string(path)
}

type KXRes[T any] struct {
    Cupboard string `json:"cupboard"`
    Data     T      `json:"data,omitempty"`
    Msg      string `json:"msg"`
    State    string `json:"state"`
}

type KXBin struct {
    Bcd string  `json:"bcd"`
    Bci float64 `json:"bci"`
    Bcu float64 `json:"bcu"`
    Bex int     `json:"bex"`
    Bfl int     `json:"bfl"`
    Bft int     `json:"bft"`
    Cpg float64 `json:"cpg"`
    Cwk int     `json:"cwk"`
    Dft int     `json:"dft"`
    Dnm string  `json:"dnm"`
    Dst int     `json:"dst"`
}

type KXStatusData []struct {
    Col     string  `json:"col"`
    Content []KXBin `json:"content"`
}
type KXStatusRes KXRes[KXStatusData]

func (r *KXRes[T]) String() string {
    return fmt.Sprintf("Code: %s, Message: %s", r.State, r.Msg)
}

// CountBins 获取仓位数量
func (r *KXStatusRes) CountBins() (n int) {
    for _, d := range r.Data {
        n += len(d.Content)
    }
    return
}

// GetBins 获取仓位
func (r *KXStatusRes) GetBins() (bins []KXBin) {
    bins = make([]KXBin, r.CountBins())
    for _, d := range r.Data {
        for _, bin := range d.Content {
            i, _ := strconv.Atoi(bin.Dnm)
            bins[i-1] = bin
        }
    }
    return
}

func (p *kaixin) GetChargerErrors(n int) (errors []string) {
    errors = make([]string, 0)
    if n == 0 {
        return
    }
    // 分隔字符串
    for _, s := range strings.Split(strconv.Itoa(n), "") {
        v, ok := p.errors[s]
        if ok {
            errors = append(errors, v)
        } else {
            errors = append(errors, s)
        }

    }
    return
}

func (p *kaixin) FetchStatus(serial string) (online bool, bins model.CabinetBins, err error) {
    res := new(KXStatusRes)
    url := p.GetUrl(kaixinUrlDetailData)
    client := resty.New().R().
        SetFormData(map[string]string{
            "user":      p.user,
            "cupboard":  serial,
            "checkcode": p.Detailcode(serial),
        })
    var r *resty.Response
    r, err = client.Post(url)

    if err != nil {
        p.logger.Write(fmt.Sprintf("凯信状态获取失败, serial: %s, err: %s\n", serial, err.Error()))
        return
    }

    // regexp.MustCompile(`(?m)({.*})(.*)`)
    err = jsoniter.Unmarshal(r.Body(), res)
    if err != nil {
        p.logger.Write(fmt.Sprintf("凯信状态解析失败, serial: %s, body: %s\n", serial, r.Body()))
        return
    }
    p.logger.Write(r.Body())

    if res.State == "ok" {
        doors := res.GetBins()
        bins = make(model.CabinetBins, len(doors))
        for index, ds := range doors {
            e := model.NewBatterySoc(utils.NewNumber().Decimal(ds.Cpg))
            hasBattery := ds.Bex == 2
            current := utils.NewNumber().Decimal(ds.Bci)
            isFull := e.IsBatteryFull()

            // 错误列表
            errs := p.GetChargerErrors(ds.Bft)

            // 2022-08-18 15:00 姓曹的说:
            // 仓位是否锁仓要加入电柜错误判定
            doorHealth := ds.Dft == 1 && len(errs) == 0

            voltage := utils.NewNumber().Decimal(ds.Bcu)
            if voltage == 0 && hasBattery {
                errs = append(errs, model.CabinetBinBatteryFault)
            }

            bin := &model.CabinetBin{
                Index:         index,
                Name:          fmt.Sprintf("%d号仓", index+1),
                BatterySN:     ds.Bcd,
                Full:          isFull,
                Battery:       hasBattery,
                Electricity:   e,
                OpenStatus:    ds.Dst == 1,
                DoorHealth:    doorHealth,
                Current:       current,
                Voltage:       voltage,
                ChargerErrors: errs,
            }

            bins[index] = bin
        }
        online = true
        return
    }

    // setOfflineTime(item.Serial, true)
    err = fmt.Errorf("[%s]凯信状态获取失败", serial)
    return
}

type KXOperationRes KXRes[any]

// DoorOperate 操作柜门
func (p *kaixin) DoorOperate(code, serial, operation string, door int) (state bool) {
    return p.doDoorOperate(code, serial, operation, door, "")
}

// BatteryBind 绑定电池
func (p *kaixin) BatteryBind(code, serial, model string, door int) (state bool) {
    re := regexp.MustCompile(`(?m)(\d+)V(\d+)AH`)
    battery := re.ReplaceAllString(strings.ToUpper(model), `JG10${1}${2}`)
    battery += fmt.Sprintf("%02d%02d", door+1, utils.RandomIntMaxMin(0, 99))
    return p.doDoorOperate(code, serial, "6", door, battery)
}

func (p *kaixin) doDoorOperate(code, serial, operation string, door int, battery string) (state bool) {
    res := new(KXOperationRes)
    url := p.GetUrl(kaixinUrlDoorOperation)
    // 凯信操作柜门index从1开始
    d := strconv.Itoa(door + 1)
    data := map[string]string{
        "user":      code,
        "cupboard":  serial,
        "checkcode": utils.Md5String(utils.Md5String(code+serial+d+operation) + p.key),
        "door":      d,
        "operation": operation,
    }
    if battery != "" {
        data["battery"] = battery
    }

    client := resty.New().R().SetFormData(data)

    r, err := client.Post(url)
    var b []byte
    if r != nil {
        b = r.Body()
    }

    zap.L().Info("发送仓门指令", log.Payload(data), zap.ByteString("response", b), zap.Error(err))

    if err != nil {
        return
    }

    err = jsoniter.Unmarshal(b, res)
    if err != nil {
        // log.Error(err)
        return
    }

    return res.State == "ok"
}
