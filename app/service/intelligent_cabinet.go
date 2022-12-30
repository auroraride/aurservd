// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    am "github.com/auroraride/adapter/model"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/goccy/go-json"
    log "github.com/sirupsen/logrus"
    "net/http"
)

type intelligentCabinetService struct {
    *BaseService
}

func NewIntelligentCabinet(params ...any) *intelligentCabinetService {
    return &intelligentCabinetService{
        BaseService: newService(params...),
    }
}

func (s *intelligentCabinetService) ExchangeUsable(serial string, br model.CabinetBrand) (uid string, info *model.RiderCabinetOperateProcess) {
    var url string

    switch br {
    case model.CabinetBrandKaixin:
        url = ar.Config.Adapter.Kaixin.Api
    default:
        snag.Panic("电柜错误")
    }

    r, err := NewAdapter(url, s.rider).Post("/exchange/usable", &am.ExchangeUsableRequest{
        Serial: serial,
        Minsoc: cache.Float64(model.SettingExchangeMinBattery),
        Lock:   10,
    })
    if err != nil {
        log.Errorf("换电信息请求失败: %v", err)
        snag.Panic("请求失败")
    }

    b := r.Body()
    log.Printf("换电请求成功: %s", string(b))

    // 解析换电信息
    res := new(am.ResponseStuff[am.ExchangeUsableResponse])
    err = json.Unmarshal(b, res)
    if err != nil {
        log.Errorf("换电信息结果解析失败: %v", err)
        snag.Panic("请求失败")
    }

    if res.Code != http.StatusOK {
        snag.Panic(res.Message)
    }

    v := res.Data
    info = &model.RiderCabinetOperateProcess{
        EmptyBin: &model.BinInfo{
            Index: v.Empty.Ordinal - 1,
        },
    }

    fully := &model.BinInfo{
        Index:       v.Fully.Ordinal - 1,
        Voltage:     v.Fully.Voltage,
        Electricity: model.BatteryElectricity(v.Fully.Soc),
    }

    if v.Fully.Soc >= model.IntelligentBatteryFullSoc {
        info.FullBin = fully
    } else {
        info.Alternative = fully
    }

    uid = v.UUID

    return
}
