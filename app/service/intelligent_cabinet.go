// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "errors"
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
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
    "golang.org/x/exp/slices"
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
func (s *intelligentCabinetService) ExchangeUsable(bm, serial string, br model.CabinetBrand) (uid string, info *model.RiderCabinetOperateProcess) {
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
        Model:  bm,
    })
    if err != nil {
        log.Errorf("换电信息请求失败: %v", err)
        snag.Panic("请求失败")
    }

    b := r.Body()
    log.Printf("换电请求成功: %s", string(b))

    // 解析换电信息
    res := new(adapter.ResponseStuff[adapter.CabinetBinUsableResponse])
    err = json.Unmarshal(b, res)
    if err != nil {
        log.Errorf("换电信息结果解析失败: %v", err)
        snag.Panic("请求失败")
    }

    err = res.VerifyResponse()
    if err != nil {
        snag.Panic(err.Error())
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
    id, err := uuid.Parse(uid)
    if err != nil {
        snag.Panic("请求参数错误")
    }

    var (
        stopAt   *time.Time
        duration float64
        success  bool
        putout   string
        putin    string
    )

    defer func() {
        updater := ex.Update()

        if err != nil {
            updater.SetRemark(err.Error())

            // 若换电失败, 标记任务失败
            _ = cache.Set(s.ctx, s.exchangeCacheKey(uid), &model.ExchangeStepResultCache{
                Index: 0,
                Results: []*adapter.ExchangeStepMessage{
                    {Step: 1, Message: err.Error(), Success: false},
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
            SetPutoutBattery(putout).
            SetPutinBattery(putin).
            Exec(s.ctx)

        bs := NewBattery()

        // 更新用户取走的电池
        if putout != "" {
            bat := bs.RiderPutout(putout, sub)
            // 更新订阅
            _ = ent.Database.Subscribe.UpdateOneID(sub.ID).SetBatterySn(bat.Sn).SetBatteryID(bat.ID).Exec(s.ctx)
        }

        // 更新放入的电池
        if putin != "" {
            bs.RiderPutin(putin, sub, cab)
        }
    }()

    var r *resty.Response
    r, err = NewAdapter(ar.Config.Adapter.Kaixin.Api, s.rider).Post("/exchange/do", &adapter.ExchangeRequest{
        UUID:    id,
        Battery: *sub.BatterySn,
        Expires: model.IntelligentBusinessScanExpires,
        Timeout: model.IntelligentBusinessStepTimeout,
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

        err = res.VerifyResponse()
        if err != nil {
            return
        }

        // 若换电成功 直接返回
        if res.Data.Success {
            putin = res.Data.PutinBattery
            putout = res.Data.PutoutBattery
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

// ExchangeStepSync 换电步骤同步
func (s *intelligentCabinetService) ExchangeStepSync(req *adapter.ExchangeStepMessage) {
    if req.Step == 0 {
        return
    }

    // TODO 检查电池是否存在???
    key := s.exchangeCacheKey(req.UUID)

    c := &model.ExchangeStepResultCache{}
    _ = cache.Get(s.ctx, key).Scan(c)
    c.Results = append(c.Results, req)

    // 排序
    slices.SortFunc(c.Results, func(a, b *adapter.ExchangeStepMessage) bool {
        return a.Step <= b.Step
    })

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

    start := time.Now()

    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for ; true; <-ticker.C {
        if time.Now().Sub(start).Seconds() > 30 {
            return
        }

        c := &model.ExchangeStepResultCache{}
        err := cache.Get(s.ctx, key).Scan(c)
        if err != nil {
            continue
        }

        n := len(c.Results)
        index := c.Index

        // 当前的数据
        if index == n {
            s.exchangeStepResultFromCache(index-1, c, res)
            if res.Step < uint8(ec.ExchangeStepPutOut) {
                res.Step += 1
            }
            res.Status = uint8(ec.TaskStatusProcessing)
        }

        if index >= n {
            continue
        }

        s.exchangeStepResultFromCache(index, c, res)

        if !res.Stop {
            ttl, _ := cache.TTL(s.ctx, key).Result()
            c.Index += 1
            cache.Set(s.ctx, key, c, ttl)
        }
        return
    }

    return
}

// stepResult 获取换电步骤结果
func (s *intelligentCabinetService) exchangeStepResultFromCache(index int, c *model.ExchangeStepResultCache, res *model.RiderExchangeProcessRes) {
    data := c.Results[index]
    res.Step = uint8(data.Step)
    res.Message = data.Message
    if data.Success {
        res.Status = uint8(ec.TaskStatusSuccess)
    }

    res.Stop = data.Step == 4 || !data.Success
}

// BusinessCensorX 校验用户是否可以使用智能柜办理业务
func (s *intelligentCabinetService) BusinessCensorX(bus adapter.Business, sub *ent.Subscribe, cab *ent.Cabinet) (bat *ent.Battery) {
    if !cab.Intelligent {
        return
    }

    // 判定电柜状态
    if cab.Status == model.CabinetStatusMaintenance.Value() {
        snag.Panic("电柜开小差了, 请联系客服")
    }

    // 判定是否智能电柜套餐
    if !sub.Intelligent {
        snag.Panic("套餐不匹配")
    }

    // 业务如果需要电池, 查找电池信息
    if bus.BatteryNeed() {
        // 判定是否智能电池
        if sub.BatterySn == nil || sub.BatteryID == nil {
            snag.Panic("必须是智能电池")
        }

        bat, _ = NewBattery().QuerySn(*sub.BatterySn)
        if bat == nil {
            snag.Panic("未找到电池信息")
        }

        // 检查电池型号与电柜型号兼容
        if !NewCabinet().ModelInclude(cab, bat.Model) {
            snag.Panic("电池型号不兼容")
        }
    }

    return
}

// BusinessUsable 获取可用的业务仓位信息
func (s *intelligentCabinetService) BusinessUsable(bus adapter.Business, serial, bm string) (uid string, index int, err error) {
    var r *resty.Response
    r, err = NewAdapter(ar.Config.Adapter.Kaixin.Api, s.rider).Post("/business/usable", &adapter.BusinuessUsableRequest{
        Minsoc:   cache.Float64(model.SettingExchangeMinBattery),
        Business: bus,
        Serial:   serial,
        Model:    bm,
    })
    if err != nil {
        return
    }

    res := new(adapter.ResponseStuff[adapter.CabinetBinUsableResponse])
    err = json.Unmarshal(r.Body(), &res)
    if err != nil {
        return
    }

    err = res.VerifyResponse()
    if err != nil {
        return
    }

    uid = res.Data.UUID
    index = res.Data.BusinessBin.Ordinal - 1
    return
}

// DoBusiness 请求办理业务
func (s *intelligentCabinetService) DoBusiness(uidstr string, bus adapter.Business, sub *ent.Subscribe, riderBat *ent.Battery, cab *ent.Cabinet) (info *model.BinInfo, batinfo *model.Battery, err error) {
    defer func() {
        // 缓存任务返回
        data := &model.BusinessCabinetStatusRes{
            Success: err == nil,
            Stop:    true,
        }
        if err != nil {
            data.Message = err.Error()
        }
        cache.Set(s.ctx, uidstr, data, 10*time.Minute)
    }()

    var uid uuid.UUID
    uid, err = uuid.Parse(uidstr)
    if err != nil {
        return
    }

    var batterySN string
    if riderBat != nil {
        batterySN = riderBat.Sn
    }

    var r *resty.Response
    r, err = NewAdapter(ar.Config.Adapter.Kaixin.Api, s.rider).Post("/business/do", &adapter.BusinessRequest{
        UUID:     uid,
        Business: bus,
        Serial:   cab.Serial,
        Timeout:  model.IntelligentBusinessStepTimeout,
        Battery:  batterySN,
        Model:    sub.Model,
    })

    res := new(adapter.ResponseStuff[adapter.BusinessResponse])
    err = json.Unmarshal(r.Body(), &res)
    if err != nil {
        return
    }

    err = res.VerifyResponse()
    if err != nil {
        return
    }

    // TODO 失败后电池信息是否更新
    if res.Data.Error != "" {
        err = errors.New(res.Data.Error)
        return
    }

    var sn string
    var putin bool
    results := res.Data.Results

    switch bus {
    case adapter.BusinessActive, adapter.BusinessContinue:
        sn = results[0].Before.BatterySN
    case adapter.BusinessPause, adapter.BusinessUnsubscribe:
        sn = results[1].After.BatterySN
        putin = true
    }

    b := results[1].After
    info = &model.BinInfo{
        Index:       b.Ordinal - 1,
        Electricity: model.BatteryElectricity(b.Soc),
        Voltage:     b.Voltage,
    }

    bs := NewBattery(s.rider)
    var bat *ent.Battery
    if putin {
        bat = bs.RiderPutin(sn, sub, cab)
        // 更新订阅
        _ = ent.Database.Subscribe.UpdateOneID(sub.ID).ClearBatterySn().ClearBatteryID().Exec(s.ctx)
    } else {
        bat = bs.RiderPutout(sn, sub)
        // 更新订阅
        _ = ent.Database.Subscribe.UpdateOneID(sub.ID).SetBatterySn(bat.Sn).SetBatteryID(bat.ID).Exec(s.ctx)
    }

    batinfo = &model.Battery{
        ID:    bat.ID,
        SN:    sn,
        Model: bat.Model,
    }

    return
}
