// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-12
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type updater struct {
    provider Provider

    cab  *ent.Cabinet
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
    bins, err = s.provider.FetchStatus(s.cab.Serial)
    if err != nil {
        setOfflineTime(s.cab.Serial, true)
        return
    }

    setOfflineTime(s.cab.Serial, false)

    old := s.cloneCabinet()

    var num, full, empty, locked int
    for i, bin := range bins {
        // 电池数量
        if bin.Battery {
            num += 1
        }
        // 满电数量
        if bin.Full {
            full += 1
        }
        // 锁仓数量
        if !bin.DoorHealth {
            locked += 1
        }
        // 空仓数量 = 无电池 && 仓门无锁
        if bin.DoorHealth && !bin.Battery {
            empty += 1
        }
        // 仓位备注信息
        if len(old.Bin) > i {
            bin.Remark = old.Bin[i].Remark
        }
        // 仓门正常清除告警设置
        if bin.DoorHealth {
            delBinFault(s.cab.Serial, i)
        }
    }

    u := s.cab.
        Update().
        SetBin(bins).
        SetBatteryNum(uint(num)).
        SetBatteryFullNum(uint(full)).
        SetHealth(uint8(res.Data.Isonline)).
        SetDoors(uint(len(bins)))

    // 判断是否处于换电过程中, 如果处于换电过程中则不保存电池数量, 以避免电池变动数量大的情况出现
    if s.task == nil {

    }
}

func (s *updater) save() {

}
