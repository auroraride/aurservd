// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "errors"
    "fmt"
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/defs/cabdef"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/go-resty/resty/v2"
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
    playload := &cabdef.ExchangeUsableRequest{
        Serial: serial,
        Minsoc: cache.Float64(model.SettingExchangeMinBattery),
        Lock:   10,
        Model:  bm,
    }

    v, err := adapter.Post[cabdef.CabinetBinUsableResponse](s.GetCabinetAdapterUrlX(br, "/exchange/usable"), s.GetAdapterUserX(), playload)

    if err != nil {
        snag.Panic(err)
    }

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
        empty    *model.BinInfo
        bs       = NewBattery()
    )

    defer func() {
        updater := ex.Update()

        if err != nil {
            updater.SetRemark(err.Error())

            // 若换电失败, 标记任务失败
            _ = cache.Set(s.ctx, s.exchangeCacheKey(uid), &model.ExchangeStepResultCache{
                Index: 0,
                Results: []*cabdef.ExchangeStepMessage{
                    {Step: 1, Message: err.Error(), Success: false},
                },
            }, 10*time.Minute).Err()
        }

        if stopAt == nil {
            stopAt = silk.Pointer(time.Now())
            duration = stopAt.Sub(ex.StartAt).Seconds()
        }

        if empty != nil {
            ex.Info.Exchange.Empty = empty
        }

        // 保存数据库
        _ = updater.
            SetSuccess(success).
            SetFinishAt(*stopAt).
            SetDuration(int(duration)).
            SetPutoutBattery(putout).
            SetPutinBattery(putin).
            Exec(s.ctx)

    }()

    playload := &cabdef.ExchangeRequest{
        UUID:    id,
        Battery: *sub.BatterySn,
        Expires: model.IntelligentBusinessScanExpires,
        Timeout: model.IntelligentBusinessStepTimeout,
        Minsoc:  cache.Float64(model.SettingExchangeMinBattery),
    }

    var v cabdef.ExchangeResponse
    v, err = adapter.Post[cabdef.ExchangeResponse](s.GetCabinetAdapterUrlX(model.CabinetBrand(cab.Brand), "/exchange/do"), s.GetAdapterUserX(), playload, func(r *resty.Response) {
        log.Infof("换电请求完成: %s", string(r.Body()))
    })

    if err != nil {
        log.Errorf("换电请求失败: %v", err)
        return
    }

    // 查询结果
    for _, result := range v.Results {
        after := result.After
        before := result.Before

        duration += result.Duration
        stopAt = result.StopAt

        if result.Success {
            // 记录用户放入的电池
            if ec.ExchangeStepPutInto.EqualInt(result.Step) && after != nil {
                putin = after.BatterySN
                empty = &model.BinInfo{
                    Index:       after.Ordinal - 1,
                    Electricity: model.BatteryElectricity(after.Current),
                    Voltage:     after.Voltage,
                }

                // 更新电池
                bat, _ := bs.LoadOrCreate(putin)
                if bat != nil {
                    _ = NewBattery().UpdateSubscribe(bat, nil)
                    _ = NewSubscribeWithRider(s.entRider).UpdateBattery(sub, nil)
                }

                go bs.RiderBusiness(true, putin, s.rider, cab, after.Ordinal)
            }

            // 记录用户取走的电池
            if result.Step == ec.ExchangeStepOpenFull.Int() && before != nil {
                putout = before.BatterySN

                go bs.RiderBusiness(false, putout, s.rider, cab, before.Ordinal)

                bat, _ := bs.LoadOrCreate(putout)
                if bat != nil {
                    // 更新订阅
                    _ = NewSubscribe().UpdateBattery(sub, bat)
                    // 更新电池
                    _ = NewBattery().UpdateSubscribe(bat, sub)
                }
            }
        }
    }

    // 若换电成功 直接返回
    if v.Success {
        success = true
        return
    }

    if v.Error != "" {
        err = fmt.Errorf("%s", v.Error)
    }
}

// ExchangeStepSync 换电步骤同步
func (s *intelligentCabinetService) ExchangeStepSync(req *cabdef.ExchangeStepMessage) {
    if req.Step == 0 {
        return
    }

    // TODO 检查电池是否存在???
    key := s.exchangeCacheKey(req.UUID)

    c := &model.ExchangeStepResultCache{}
    _ = cache.Get(s.ctx, key).Scan(c)
    c.Results = append(c.Results, req)

    // 排序
    slices.SortFunc(c.Results, func(a, b *cabdef.ExchangeStepMessage) bool {
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
func (s *intelligentCabinetService) BusinessUsable(br model.CabinetBrand, bus adapter.Business, serial, bm string) (uid string, index int, err error) {
    playload := &cabdef.BusinuessUsableRequest{
        Minsoc:   cache.Float64(model.SettingExchangeMinBattery),
        Business: bus,
        Serial:   serial,
        Model:    bm,
    }

    var v cabdef.CabinetBinUsableResponse
    v, err = adapter.Post[cabdef.CabinetBinUsableResponse](s.GetCabinetAdapterUrlX(br, "/business/usable"), s.GetAdapterUserX(), playload)
    if err != nil {
        return
    }

    uid = v.UUID
    index = v.BusinessBin.Ordinal - 1
    return
}

// DoBusiness 请求办理业务
func (s *intelligentCabinetService) DoBusiness(br model.CabinetBrand, uidstr string, bus adapter.Business, sub *ent.Subscribe, riderBat *ent.Battery, cab *ent.Cabinet) (info *model.BinInfo, batinfo *model.Battery, err error) {
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

    playload := &cabdef.BusinessRequest{
        UUID:     uid,
        Business: bus,
        Serial:   cab.Serial,
        Timeout:  model.IntelligentBusinessStepTimeout,
        Battery:  batterySN,
        Model:    sub.Model,
    }

    var v cabdef.BusinessResponse
    v, err = adapter.Post[cabdef.BusinessResponse](s.GetCabinetAdapterUrlX(br, "/business/do"), s.GetAdapterUserX(), playload)

    if err != nil {
        return
    }

    // TODO 失败后电池信息是否更新
    if v.Error != "" {
        err = errors.New(v.Error)
        return
    }

    var sn string
    var putin bool
    results := v.Results

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

    // 获取电池
    var bat *ent.Battery
    bat, err = NewBattery().LoadOrCreate(sn)
    if err != nil {
        log.Errorf("业务记录失败: %v", err)
        return
    }

    batinfo = &model.Battery{
        ID:    bat.ID,
        SN:    sn,
        Model: bat.Model,
    }

    // 更新订阅和电池
    ss := NewSubscribeWithRider(s.entRider)
    bs := NewBattery(s.rider)
    if putin {
        // 放入电池
        _ = ss.UpdateBattery(sub, nil)
        _ = bs.UpdateSubscribe(bat, nil)
    } else {
        // 取走电池
        _ = ss.UpdateBattery(sub, bat)
        _ = bs.UpdateSubscribe(bat, sub)
    }

    return
}

func (s *intelligentCabinetService) Operate(br model.CabinetBrand, op cabdef.Operate, serial string, ordinal int, remark string) bool {
    playload := &cabdef.OperateBinRequest{
        Operate: op,
        Ordinal: silk.Int(ordinal),
        Serial:  serial,
        Remark:  remark,
    }

    _, err := adapter.Post[[]*cabdef.BinOperateResult](s.GetCabinetAdapterUrlX(br, "/business/do"), s.GetAdapterUserX(), playload)

    return err == nil
}