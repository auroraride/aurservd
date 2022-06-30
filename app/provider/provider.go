// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
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
    UpdateStatus(up *ent.CabinetUpdateOne, item *ent.Cabinet) error
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
                    health := item.Health
                    oldBins := item.Bin
                    oldNum := item.BatteryNum
                    // 未获取到电柜状态设置为离线
                    up := ent.Database.Cabinet.UpdateOne(item).SetHealth(model.CabinetHealthStatusOffline)
                    err = provider.UpdateStatus(up, item)

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
                        // 监控在线变化
                        if health != item.Health {
                            logging.NewHealthLog(item.Brand, item.Serial, item.UpdatedAt).SetStatus(health, item.Health).Send()
                        }
                        // 监控电池变化
                        if oldNum != item.BatteryNum {
                            // 判断电柜是否正在执行业务
                            info, busy := model.CabinetProcessJob(item.Serial)
                            // 当前非业务状态或电池变动数量大于1时
                            if !busy || math.Abs(float64(item.BatteryNum)-float64(oldNum)) > 1 {
                                logging.NewBatteryLog(item.Brand, item.Serial, int(item.BatteryNum)-int(oldNum), item.UpdatedAt).
                                    SetExchangeProcess(info).
                                    SetBin(oldBins, item.Bin).
                                    Send()
                            }
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
