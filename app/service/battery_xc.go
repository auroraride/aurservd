// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/adapter"
    "github.com/auroraride/adapter/defs/xcdef"
    "github.com/auroraride/adapter/rpc/pb"
    "github.com/auroraride/adapter/rpc/pb/xcpb"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/rpc"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/battery"
    "github.com/auroraride/aurservd/pkg/snag"
    "math"
    "strconv"
    "time"
)

type batteryXcService struct {
    *BaseService
}

func NewBatteryXc(params ...any) *batteryXcService {
    return &batteryXcService{
        BaseService: newService(params...),
    }
}

func (s *batteryXcService) Detail(req *model.XcBatterySNRequest) (detail *model.XcBatteryDetail) {
    // 请求xcbms rpc
    r, _ := rpc.XcBmsBatch(s.ctx, &pb.BatteryBatchRequest{Sn: []string{req.SN}})
    if r == nil {
        snag.Panic("电池信息查询失败")
    }

    // 查询电池
    bat := ent.Database.Battery.Query().Where(battery.Sn(req.SN)).WithRider().WithCabinet().FirstX(s.ctx)
    if bat == nil {
        snag.Panic("电池未录入")
    }

    var (
        hb *xcpb.Heartbeat
        rb *xcpb.Battery
    )

    if len(r.Items[req.SN].Heartbeats) > 0 {
        rb = r.Items[req.SN]
        hb = rb.Heartbeats[0]
    }

    detail = &model.XcBatteryDetail{}
    if hb != nil {
        detail = &model.XcBatteryDetail{
            UpdatedAt:            hb.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
            XcBmsBattery:         model.NewXcBmsBattery(hb),
            Current:              hb.Current,
            Soh:                  uint8(hb.Soh),
            Cycles:               uint16(hb.Cycles),
            Geom:                 adapter.NewGeometry(hb.Geom).WGS84toGCJ02(),
            Voltage:              hb.Voltage,
            Power:                math.Round(math.Abs(hb.Current)*hb.Voltage/1000*100) / 100,
            ChargingTime:         hb.ChargingTime,
            DisChargingTime:      hb.DisChargingTime,
            UsingTime:            hb.UsingTime,
            TotalChargingTime:    hb.TotalChargingTime,
            TotalDisChargingTime: hb.TotalDisChargingTime,
            TotalUsingTime:       hb.TotalUsingTime,
            InCabinet:            hb.InCabinet,
            Capacity:             hb.Capacity,
            MonMaxVoltage:        uint16(hb.MonMaxVoltage),
            MonMaxVoltagePos:     uint8(hb.MonMaxVoltagePos),
            MonMinVoltage:        uint16(hb.MonMinVoltage),
            MonMinVoltagePos:     uint8(hb.MonMinVoltagePos),
            MaxTemp:              uint16(hb.MaxTemp),
            MinTemp:              uint16(hb.MinTemp),
            MosStatus:            xcdef.NewMosStatus(hb.MosStatus),
            MonVoltage:           xcdef.NewMonVoltage(hb.MonVoltage),
            Temp:                 xcdef.NewTemperature(hb.Temp),
            MosTemp:              uint16(hb.MosTemp),
            EnvTemp:              uint16(hb.EnvTemp),
            Strength:             uint8(hb.Strength),
            Gps:                  xcdef.GPSStatus(hb.Gps),
            Online:               time.Now().Sub(hb.CreatedAt.AsTime()).Minutes() < 35,
            FaultsOverview:       make([]*pb.BatteryFaultOverview, 0),
        }
    }

    if rb != nil {
        detail.SoftVersion = rb.SoftVersion.Value
        detail.HardVersion = rb.HardVersion.Value
        detail.Soft4gVersion = rb.Soft_4GVersion.Value
        detail.Hard4gVersion = rb.Hard_4GVersion.Value
        detail.Sn4g = rb.Sn_4G.Value
        detail.Iccid = rb.Iccid.Value
    }

    detail.Sn = req.SN
    detail.CreatedAt = bat.CreatedAt.Format("2006-01-02 15:04:05")

    if bat.Edges.Cabinet != nil && bat.Ordinal != nil {
        detail.BelongsTo = bat.Edges.Cabinet.Name + "-" + strconv.Itoa(*bat.Ordinal) + "号仓"
    }

    if bat.Edges.Rider != nil {
        detail.BelongsTo = bat.Edges.Rider.Name + "-" + bat.Edges.Rider.Phone
    }

    fr, _ := rpc.XcBmsFaultOverview(s.ctx, &pb.BatterySnRequest{Sn: req.SN})
    if fr != nil {
        detail.FaultsOverview = fr.Items
    }

    return
}

func (s *batteryXcService) Statistics(req *model.XcBatterySNRequest) (detail *model.XcBatteryStatistics) {
    // 请求xcbms rpc
    r, _ := rpc.XcBmsStatistics(s.ctx, &pb.BatterySnRequest{Sn: req.SN})
    if r == nil {
        snag.Panic("电池数据查询失败")
    }

    return &model.XcBatteryStatistics{
        DateHour:    r.DateHour,
        Voltage:     r.Voltage,
        BatTemp:     r.BatTemp,
        MosTemp:     r.MosTemp,
        EnvTemp:     r.EnvTemp,
        Soc:         r.Soc,
        Strength:    r.Strength,
        Charging:    r.Charging,
        DisCharging: r.DisCharging,
    }
}
