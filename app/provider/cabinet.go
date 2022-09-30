// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-12
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "context"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/workwx"
    "github.com/auroraride/aurservd/internal/ent"
    log "github.com/sirupsen/logrus"
    "math"
)

type updater struct {
    provider Provider
    ctx      context.Context

    cab  *ent.Cabinet
    old  *ent.Cabinet
    task *ec.Task
}

type CabinetStatus struct {
    Health uint8
    Bins   model.CabinetBins
}

func NewUpdater(cab *ent.Cabinet) *updater {
    var prov Provider
    if cab.Brand == model.CabinetBrandKaixin.Value() {
        prov = NewKaixin()
    } else {
        prov = NewYundong()
    }
    return &updater{
        cab:      cab,
        provider: prov,
        ctx:      context.Background(),
    }
}

func (s *updater) cloneCabinet() *ent.Cabinet {
    item := new(ent.Cabinet)
    *item = *s.cab
    return item
}

func (s *updater) DoUpdate() (err error) {
    // 获取电柜当前执行的任务
    s.task = ec.Obtain(ec.ObtainReq{Serial: s.cab.Serial})
    var bins model.CabinetBins
    var online bool
    online, bins, err = s.provider.FetchStatus(s.cab.Serial)

    // 设置是否离线
    setOfflineTime(s.cab.Serial, online)

    s.old = s.cloneCabinet()

    up := s.cab.Update()

    if err == nil && online {
        var num, full, empty, locked, charging int
        for i, bin := range bins {
            // 电池数量
            if bin.Battery {
                num += 1
            }
            // 仓位备注信息
            if len(s.old.Bin) > i {
                bin.Remark = s.old.Bin[i].Remark
            }
            // 锁仓判定
            if bin.DoorHealth && len(bin.ChargerErrors) == 0 {
                // 仓门正常清除告警设置
                delBinFault(s.cab.Serial, i)
                // 满电充电空仓判定
                if bin.Battery {
                    if bin.Full {
                        // 满电数量
                        full += 1
                    } else {
                        // 充电数量
                        charging += 1
                    }
                } else {
                    // 空仓数量 = 无电池 && 仓门无锁
                    empty += 1
                }
            } else {
                // 锁仓数量
                locked += 1
            }
        }

        up.SetBin(bins).SetBatteryNum(num).SetDoors(len(bins)).SetLockedBinNum(locked).SetEmptyBinNum(empty).SetBatteryFullNum(full).SetBatteryChargingNum(charging)
    }

    health := s.cab.Health
    if online {
        health = model.CabinetHealthStatusOnline
    } else if isOffline(s.cab.Serial) {
        health = model.CabinetHealthStatusOffline
    }

    var item *ent.Cabinet
    item, err = up.SetHealth(health).Save(s.ctx)
    log.Errorf("%s更新写入失败: %v", s.cab.Serial, err)
    if err != nil {
        return
    }

    // 在线变化
    if s.old.Health != health {
        // 电柜在线变动日志
        logging.NewHealthLog(item.Brand, item.Serial, item.UpdatedAt).SetStatus(s.old.Health, health).Send()
        if health == model.CabinetHealthStatusOffline {
            go workwx.New().SendCabinetOffline(item.Name, item.Serial, cabinetCity(item))
        }
    }

    *s.cab = *item

    if item.Health == model.CabinetHealthStatusOnline {
        // 电池变动
        s.batteryMonitor()
    }

    return
}

// batteryMonitor 监控电柜电池变动
//
// 电池异常变动
// 1. 每次电池变动均留存阿里云日志
// 2. 电池变动时判定当前是否有任务和维护状态
// 2.1 如果有任务或维护时变动数量>=2推送
// 2.2 如果没有任务且非维护中时变动数量>=1推送
//
func (s *updater) batteryMonitor() {
    oldNum := s.old.BatteryNum
    oldBins := s.old.Bin
    max := 1.0
    // 判断电柜是否正在执行业务, 若正在执行任务则使用执行任务中的电池数量信息
    if s.task != nil {
        oldNum = s.task.Cabinet.BatteryNum
        oldBins = s.task.Cabinet.Bins
        max = 2.0
    }

    // 监控电池变化
    if oldNum != s.cab.BatteryNum {
        diff := s.cab.BatteryNum - oldNum
        // 当前非业务状态或电池变动数量大于1时
        if math.Abs(float64(diff)) >= max {
            status := model.CabinetStatus(s.cab.Status)
            logging.NewBatteryLog(s.cab.Brand, s.cab.Serial, oldNum, s.cab.BatteryNum, s.cab.UpdatedAt).
                SetTask(s.task).
                SetBin(oldBins, s.cab.Bin).
                SetStatus(status).
                Send()

            // 推送消息
            go func() {
                // 非运营中状态不推送
                if status == model.CabinetStatusNormal {
                    workwx.New().SendBatteryAbnormality(cabinetCity(s.cab), s.cab.Serial, s.cab.Name, oldNum, s.cab.BatteryNum, diff)
                }
            }()
        }
    }
}
