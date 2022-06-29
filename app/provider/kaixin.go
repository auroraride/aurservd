// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-16
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
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
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
    return ar.Ent.Cabinet.Query().Where(cabinet.Brand(model.CabinetBrandKaixin.Value()), cabinet.Status(model.CabinetStatusNormal)).All(context.Background())
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

func (p *kaixin) PrepareRequest() {}

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

func (p *kaixin) UpdateStatus(up *ent.CabinetUpdateOne, item *ent.Cabinet) {
    res := new(KXStatusRes)
    url := p.GetUrl(kaixinUrlDetailData)
    client := resty.New().R().
        SetFormData(map[string]string{
            "user":      p.user,
            "cupboard":  item.Serial,
            "checkcode": p.Detailcode(item.Serial),
        })
    r, err := client.Post(url)

    if err != nil {
        p.logger.Write(fmt.Sprintf("凯信状态获取失败, serial: %s, err: %s", item.Serial, err.Error()))
        return
    }

    err = jsoniter.Unmarshal(r.Body(), res)
    if err != nil {
        p.logger.Write(fmt.Sprintf("凯信状态解析失败, serial: %s, body: %s", item.Serial, r.Body()))
        return
    }

    if res.State == "ok" {
        doors := res.GetBins()
        bins := make([]model.CabinetBin, len(doors))
        var full uint = 0
        var num uint = 0
        for index, ds := range doors {
            e := model.NewBatteryElectricity(utils.NewNumber().Decimal(ds.Cpg))
            hasBattery := ds.Bex == 2
            current := utils.NewNumber().Decimal(ds.Bci)
            isFull := e.IsBatteryFull()
            if hasBattery {
                num += 1
            }
            if isFull {
                full += 1
            }

            // 错误列表
            errs := p.GetChargerErrors(ds.Bft)

            voltage := utils.NewNumber().Decimal(ds.Bcu)
            if voltage == 0 && hasBattery {
                errs = append(errs, model.CabinetBinBatteryFault)
            }

            bin := model.CabinetBin{
                Index:         index,
                Name:          fmt.Sprintf("%d号仓", index+1),
                BatterySN:     ds.Bcd,
                Full:          isFull,
                Battery:       hasBattery,
                Electricity:   e,
                OpenStatus:    ds.Dst == 1,
                DoorHealth:    ds.Dft == 1,
                Current:       current,
                Voltage:       voltage,
                ChargerErrors: errs,
            }

            if len(item.Bin) > index {
                bin.Remark = item.Bin[index].Remark
            }

            bins[index] = bin
        }
        v, _ := up.SetBatteryFullNum(full).
            SetBatteryNum(num).
            SetBin(bins).
            SetHealth(model.CabinetHealthStatusOnline).
            SetDoors(uint(len(doors))).
            Save(context.Background())
        *item = *v
    }
    p.logger.Write(res)
    return
}

type KXOperationRes KXRes[any]

// DoorOperate 操作柜门
func (p *kaixin) DoorOperate(code, serial, operation string, door int) (state bool) {
    res := new(KXOperationRes)
    url := p.GetUrl(kaixinUrlDoorOperation)
    // 凯信操作柜门index从1开始
    d := strconv.Itoa(door + 1)
    client := resty.New().R().
        SetFormData(map[string]string{
            "user":      code,
            "cupboard":  serial,
            "checkcode": utils.Md5String(utils.Md5String(code+serial+d+operation) + p.key),
            "door":      d,
            "operation": operation,
        })
    r, err := client.Post(url)
    log.Info(string(r.Body()))

    if err != nil {
        log.Error(err)
        return
    }

    err = jsoniter.Unmarshal(r.Body(), res)
    if err != nil {
        log.Error(err)
        return
    }

    return res.State == "ok"
}
