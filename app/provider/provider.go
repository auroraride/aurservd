// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    log "github.com/sirupsen/logrus"
    "sync"
    "time"
)

type Provider interface {
    PrepareRequest()
    Cabinets() ([]*ent.Cabinet, error)
    Brand() string
    Logger() *Logger
    UpdateStatus(up *ent.CabinetUpdateOne, item *ent.Cabinet) any
    DoorOperate(code, serial, operation string, door int) bool
    Reboot(code string, serial string) bool
}

func Run(start bool) {
    yd := NewYundong()
    kx := NewKaixin()
    if start {
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
                log.Infof("开始 第%d轮 轮询获取%s电柜状态", times, provider.Brand())
                start := time.Now()

                items, err := provider.Cabinets()
                if err != nil {
                    log.Errorf("%s电柜获取失败: %#v", provider.Brand(), err)
                    return
                }

                length := len(items)

                var wg sync.WaitGroup
                wg.Add(length)

                logs := make([]any, length)
                for i, item := range items {
                    // 未获取到电柜状态设置为离线
                    up := ar.Ent.Cabinet.UpdateOne(item).SetHealth(model.CabinetHealthStatusOffline)
                    go func(i int, item *ent.Cabinet) {
                        l := provider.UpdateStatus(up, item)
                        if l != nil {
                            // 日志归并
                            logs[i] = l
                        }

                        // 存储电柜信息
                        ca := up.SaveX(context.Background())

                        go func() {
                            // 保存历史仓位信息(转换后的)
                            lg := GenerateSlsStatusLogGroup(ca)
                            if lg != nil {
                                err = ali.NewSls().PutLogs(slsCfg.Project, slsCfg.CabinetLog, lg)
                                if err != nil {
                                    log.Errorf("阿里云SLS提交失败: %#v", err)
                                }
                            }
                        }()

                        wg.Done()
                    }(i, item)
                }

                wg.Wait()

                // 写入电柜日志
                provider.Logger().Write(times, logs)
                log.Infof("%s电柜 第%d轮 状态轮询完成, 耗时%.2fs", provider.Brand(), times, time.Now().Sub(start).Seconds())
                time.Sleep(1 * time.Minute)
            }
        }()
    }
}
