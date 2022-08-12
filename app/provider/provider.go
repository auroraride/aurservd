// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/workwx"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
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
    FetchStatus(serial string) (bool, model.CabinetBins, error)
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

func cabinetCity(cab *ent.Cabinet) string {
    c := cab.Edges.City
    if c == nil && cab.CityID != nil {
        c, _ = ent.Database.City.Query().Where(city.ID(*cab.CityID)).First(context.Background())
    }
    return c.Name
}

// monitor 监控电柜变动
// bins health num 旧数据
func monitor(oldBins model.CabinetBins, oldHealth uint8, oldNum uint, item *ent.Cabinet) {
    // 监控在线变化
    if oldHealth != item.Health {
        // 电柜在线变动日志
        logging.NewHealthLog(item.Brand, item.Serial, item.UpdatedAt).SetStatus(oldHealth, item.Health).Send()
        go func() {
            workwx.New().SendCabinetOffline(item.Name, item.Serial, cabinetCity(item))
        }()
    }

    // 监控电池变化
    if oldNum != item.BatteryNum {
        // 判断电柜是否正在执行业务
        task := ec.Obtain(ec.ObtainReq{Serial: item.Serial})
        if task != nil {
            oldNum = task.Cabinet.BatteryNum
        }
        diff := int(item.BatteryNum) - int(oldNum)
        // 当前非业务状态或电池变动数量大于1时
        if task == nil || math.Abs(float64(diff)) > 1 {
            status := model.CabinetStatus(item.Status)
            logging.NewBatteryLog(item.Brand, item.Serial, int(oldNum), int(item.BatteryNum), item.UpdatedAt).
                SetExchangeProcess(task).
                SetBin(oldBins, item.Bin).
                SetStatus(status).
                Send()

            // 推送消息
            go func() {
                // 非运营中状态不推送
                if status == model.CabinetStatusNormal {
                    workwx.New().SendBatteryAbnormality(cabinetCity(item), item.Serial, item.Name, oldNum, item.BatteryNum, diff)
                }
            }()
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
                    if ec.Busy(item.Serial) {
                        continue
                    }

                    // 更新电柜信息
                    err = provider.FetchStatus(item)

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
