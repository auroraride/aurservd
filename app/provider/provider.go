// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    log "github.com/sirupsen/logrus"
    "time"
)

type LogCallback func(data any)

type Provider interface {
    PrepareRequest()
    Cabinets() ([]*ent.Cabinet, error)
    Brand() string
    Logger() *Logger
    UpdateStatus(up *ent.CabinetUpdateOne, item *ent.Cabinet)
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
                provider.Logger().Write(fmt.Sprintf("开始 第%d轮 轮询获取%s电柜状态", times, provider.Brand()))
                start := time.Now()

                items, err := provider.Cabinets()
                if err != nil {
                    log.Errorf("%s电柜获取失败: %#v", provider.Brand(), err)
                    return
                }

                for _, item := range items {
                    // 未获取到电柜状态设置为离线
                    up := ar.Ent.Cabinet.UpdateOne(item).SetHealth(model.CabinetHealthStatusOffline)
                    provider.UpdateStatus(up, item)

                    // 存储电柜信息
                    ca, _ := up.Save(context.Background())

                    // 提交日志
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
                }

                // 写入电柜日志
                provider.Logger().Write(fmt.Sprintf("%s电柜 第%d轮 状态轮询完成, 耗时%.2fs\n\n", provider.Brand(), times, time.Now().Sub(start).Seconds()))

                time.Sleep(1 * time.Minute)
            }
        }()
    }
}
