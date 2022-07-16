// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/workwx"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/pkg/cache"
    log "github.com/sirupsen/logrus"
    "math"
    "time"
)

type LogCallback func(data any)

type Provider interface {
    PrepareRequest()
    Cabinets() ([]*ent.Cabinet, error)
    Brand() string
    Logger() *Logger
    UpdateStatus(item *ent.Cabinet, params ...any) error
    DoorOperate(code, serial, operation string, door int) bool
    Reboot(code string, serial string) bool
}

func Run() {
    yd := NewYundong()
    kx := NewKaixin()
    if ar.Config.Cabinet.Provider {
        StartCabinetProvider(yd, kx)
    }
}

// getOfflineTime 获取离线时间
func getOfflineTime(serial string) time.Time {
    t, _ := cache.Get(context.Background(), fmt.Sprintf("OFFLINE-%s", serial)).Time()
    return t
}

// setOfflineTime 设置离线时间
func setOfflineTime(serial string, offline bool) {
    key := fmt.Sprintf("OFFLINE-%s", serial)
    if offline {
        t := getOfflineTime(serial)
        if t.IsZero() {
            cache.Set(context.Background(), key, time.Now(), -1)
        }
    } else {
        // 在线则删除
        cache.Del(context.Background(), key)
    }
}

// isOffline 判定电柜是否离线, 3分钟以上算作离线
func isOffline(serial string) bool {
    t := getOfflineTime(serial)
    return !t.IsZero() && time.Now().Sub(t).Minutes() > 3
}

func cabinetCity(cab *ent.Cabinet) string {
    c := cab.Edges.City
    if c == nil && cab.CityID != nil {
        c, _ = ent.Database.City.Query().Where(city.ID(*cab.CityID)).First(context.Background())
    }
    return c.Name
}

// AutoBinFault 自动处理换电仓门操作失败
// 每次都推送, 超过两次锁仓
func AutoBinFault(operator model.CabinetDoorOperator, cab *ent.Cabinet, index int, status bool, lock func()) {
    binKey := fmt.Sprintf("AUTO-BINFAULT-%s-%d", cab.Serial, index)
    ctx := context.Background()
    if status {
        cache.Del(ctx, binKey)
    } else {
        // 查询并保存失败次数
        times, _ := cache.Get(ctx, binKey).Int()
        times += 1
        cache.Set(ctx, binKey, times, -1)
        // 推送消息
        go func() {
            _ = workwx.New().ExchangeBinFault(cabinetCity(cab), cab.Name, cab.Serial, fmt.Sprintf("%d号仓", index+1), operator.Name, operator.Phone, times)
        }()
        // 锁仓
        if times > 2 {
            lock()
        }
    }
}

// delBinFault 删除仓门自动故障
func delBinFault(serial string, index int) {
    cache.Del(context.Background(), fmt.Sprintf("AUTO-BINFAULT-%s-%d", serial, index))
}

// monitor 监控电柜变动
// bins health num 旧数据
func monitor(oldBins []model.CabinetBin, oldHealth uint8, oldNum uint, item *ent.Cabinet) {
    // 监控在线变化
    if oldHealth != item.Health {
        // 电柜在线变动日志
        logging.NewHealthLog(item.Brand, item.Serial, item.UpdatedAt).SetStatus(oldHealth, item.Health).Send()
        go func() {
            _ = workwx.New().SendCabinetOffline(item.Name, item.Serial, cabinetCity(item))
        }()
    }

    // 监控电池变化
    if oldNum != item.BatteryNum {
        // 判断电柜是否正在执行业务
        info, busy := model.CabinetProcessJob(item.Serial)
        if busy {
            oldNum = info.BatteryNum
        }
        diff := int(item.BatteryNum) - int(oldNum)
        // 当前非业务状态或电池变动数量大于1时
        if !busy || math.Abs(float64(diff)) > 1 {
            logging.NewBatteryLog(item.Brand, item.Serial, int(oldNum), int(item.BatteryNum), item.UpdatedAt).
                SetExchangeProcess(info).
                SetBin(oldBins, item.Bin).
                Send()
        }
    }
}

// StartCabinetProvider 开始执行任务
func StartCabinetProvider(providers ...Provider) {
    slsCfg := ar.Config.Aliyun.Sls
    for _, p := range providers {
        provider := p
        times := 0
        go func() {
            for {
                times += 1
                provider.PrepareRequest()
                provider.Logger().Write(fmt.Sprintf("开始第%d次%s电柜状态轮询\n", times, provider.Brand()))
                start := time.Now()

                items, err := provider.Cabinets()
                if err != nil {
                    log.Errorf("%s电柜获取失败: %#v", provider.Brand(), err)
                    return
                }

                for _, item := range items {
                    // 换电过程不查询状态
                    if model.CabinetBusying(item.Serial) {
                        continue
                    }

                    // 更新电柜信息
                    err = provider.UpdateStatus(item)

                    if err == nil {
                        // 提交日志
                        if item.Health == model.CabinetHealthStatusOnline {
                            go func() {
                                // 保存历史仓位信息(转换后的)
                                lg := GenerateSlsStatusLogGroup(item)
                                if lg != nil {
                                    err = ali.NewSls().PutLogs(slsCfg.Project, slsCfg.CabinetLog, lg)
                                    if err != nil {
                                        log.Errorf("阿里云SLS提交失败: %#v", err)
                                    }
                                }
                            }()
                        }
                    }

                    time.Sleep(time.Duration((60000-int(time.Now().Sub(start).Milliseconds()))/len(items)) * time.Millisecond)
                }

                // 写入电柜日志
                provider.Logger().Write(fmt.Sprintf("完成第%d次%s电柜状态轮询, 耗时%.2fs\n\n", times, provider.Brand(), time.Now().Sub(start).Seconds()))

                time.Sleep(1 * time.Minute)
            }
        }()
    }
}
