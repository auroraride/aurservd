// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/adapter"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/go-resty/resty/v2"
    "github.com/goccy/go-json"
    log "github.com/sirupsen/logrus"
    "net/http"
    "time"
)

type intelligentCabinetService struct {
    *BaseService
}

func NewIntelligentCabinet(params ...any) *intelligentCabinetService {
    return &intelligentCabinetService{
        BaseService: newService(params...),
    }
}

// ExchangeUsable 获取换电信息
func (s *intelligentCabinetService) ExchangeUsable(serial string, br model.CabinetBrand) (uid string, info *model.RiderCabinetOperateProcess) {
    var url string

    switch br {
    case model.CabinetBrandKaixin:
        url = ar.Config.Adapter.Kaixin.Api
    default:
        snag.Panic("电柜错误")
    }

    r, err := NewAdapter(url, s.rider).Post("/exchange/usable", &adapter.ExchangeUsableRequest{
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
    res := new(adapter.ResponseStuff[adapter.ExchangeUsableResponse])
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

func (s *intelligentCabinetService) exchangeCacheKey(uid string) string {
    return "INTELLIGENT-CABINET-EXCHANGE-" + uid
}

// Exchange 请求换电
func (s *intelligentCabinetService) Exchange(uid string, ex *ent.Exchange, sub *ent.Subscribe, cab *ent.Cabinet) {
    var (
        stopAt        *time.Time
        duration      float64
        success       bool
        err           error
        afterBattery  string
        beforeBattery string
    )

    defer func() {
        updater := ex.Update()

        if err != nil {
            updater.SetRemark(err.Error())

            // 若换电失败, 标记任务失败
            _ = cache.Set(s.ctx, s.exchangeCacheKey(uid), &model.ExchangeStepResultCache{
                Index: -1,
                Results: []*adapter.ExchangeStepMessage{
                    {Step: adapter.ExchangeStepFirst, Message: err.Error(), Success: false},
                },
            }, 10*time.Minute).Err()
        }

        if stopAt == nil {
            stopAt = silk.Pointer(time.Now())
            duration = stopAt.Sub(ex.StartAt).Seconds()
        }

        // 保存数据库
        _ = updater.
            SetSuccess(success).
            SetFinishAt(*stopAt).
            SetDuration(int(duration)).
            SetAfterBattery(afterBattery).
            Exec(s.ctx)

        if success {
            // 查找电池
            // TODO 是否自动添加电池???
            var bat *ent.Battery
            bat, _ = NewBattery().LoadOrStore(afterBattery, *cab.CityID)
            if bat == nil {
                log.Error("用户订阅更新失败, 未找到电池信息")
            }

            // 更新套餐
            _ = ent.Database.Subscribe.UpdateOneID(ex.SubscribeID).SetBatterySn(bat.Sn).SetBatteryID(bat.ID).Exec(s.ctx)

            // 更新电池
            _ = NewBattery().UpdateRider(afterBattery, ex.RiderID)

            // 清除之前电池用户信息
            _ = NewBattery().UpdateRider(beforeBattery, 0)
        }
    }()

    var r *resty.Response
    r, err = NewAdapter(ar.Config.Adapter.Kaixin.Api, s.rider).Post("/exchange/do", &adapter.ExchangeRequest{
        UUID:    uid,
        Battery: *sub.BatterySn,
        Expires: model.IntelligentScanExpires,
        TimeOut: model.IntelligentExchangeTimeout,
        Minsoc:  cache.Float64(model.SettingExchangeMinBattery),
    })

    if err != nil {
        log.Errorf("换电请求失败: %v", err)
        return
    } else {
        b := r.Body()
        log.Infof("换电请求完成: %s", string(b))

        // 解析换电结果
        res := new(adapter.ResponseStuff[adapter.ExchangeResponse])
        _ = json.Unmarshal(b, res)

        // 若换电成功 直接返回
        if res.Code != http.StatusOK {
            err = fmt.Errorf("%s", res.Message)
            return
        }

        if res.Data.Success {
            beforeBattery = res.Data.BeforeBattery
            afterBattery = res.Data.AfterBattery
            success = true
            steps := res.Data.Results
            n := len(steps)
            for i, result := range steps {
                duration += result.Duration
                if i == n-1 {
                    stopAt = result.StopAt
                }
            }
            return
        }

        if res.Data.Error != "" {
            err = fmt.Errorf("%s", res.Data.Error)
        }
    }
}

func (s *intelligentCabinetService) ExchangeStep(req *adapter.ExchangeStepMessage) {
    if req.Step == 0 {
        return
    }

    // 检查电池是否存在

    key := s.exchangeCacheKey(req.UUID)

    c := &model.ExchangeStepResultCache{}
    _ = cache.Get(s.ctx, key).Scan(c)
    c.Results = append(c.Results, req)
    err := cache.Set(s.ctx, key, c, 10*time.Minute).Err()
    if err != nil {
        log.Error(err)
        return
    }
}

// ExchangeResult 查询换电结果
func (s *intelligentCabinetService) ExchangeResult(uid string) (res *model.RiderExchangeProcessRes) {
    key := s.exchangeCacheKey(uid)
    res = &model.RiderExchangeProcessRes{
        Step:   uint8(ec.ExchangeStepOpenEmpty),
        Status: uint8(ec.TaskStatusProcessing),
    }

    ticker := time.NewTicker(1 * time.Second)
    start := time.Now()
    for {
        select {
        case <-ticker.C:
            c := &model.ExchangeStepResultCache{}
            err := cache.Get(s.ctx, key).Scan(c)
            if err != nil {
                continue
            }

            n := len(c.Results)
            index := c.Index
            if index >= n {
                continue
            }

            if n-1 == index {
                index = n - 1
            } else {
                index += 1
            }

            data := c.Results[index]
            res.Step = uint8(data.Step)
            res.Message = data.Message
            if data.Success {
                res.Status = uint8(ec.TaskStatusSuccess)
            }

            res.Stop = data.Step == adapter.ExchangeStepFourth || !data.Success

            if res.Stop || time.Now().Sub(start).Seconds() > 30 {
                return
            }
        }
    }
}

func (s *intelligentCabinetService) ExchangeCensorX(sub *ent.Subscribe, cab *ent.Cabinet) {
    if cab.Status == model.CabinetStatusMaintenance.Value() {
        snag.Panic("电柜开小差了, 请联系客服")
    }

    if cab.Intelligent {
        if !sub.Intelligent {
            snag.Panic("套餐不匹配")
        }
        if sub.BatterySn == nil || sub.BatteryID == nil {
            snag.Panic("必须是智能电池")
        }

        // 查找电池信息
        bat, _ := NewBattery().QuerySn(*sub.BatterySn)
        if bat == nil {
            snag.Panic("未找到电池信息")
        }

        // 检查电池型号与电柜型号兼容
        if !NewCabinet().ModelInclude(cab, bat.Model) {
            snag.Panic("电池型号不兼容")
        }
    }
}
