// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/adapter/rpc/pb"
    "github.com/auroraride/adapter/rpc/pb/xcpb"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/rpc"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/battery"
    "github.com/auroraride/aurservd/pkg/snag"
    "strconv"
)

type batteryXcService struct {
    *BaseService
}

func NewBatteryXc(params ...any) *batteryXcService {
    return &batteryXcService{
        BaseService: newService(params...),
    }
}

func (s *batteryXcService) Detail(req *model.XcBatteryDetailRequest) (detail *model.XcBatteryDetail) {
    // 请求xcbms rpc
    r, _ := rpc.XcBmsBatch(s.ctx, &pb.BatteryBatchRequest{Sn: []string{req.SN}})
    if r == nil {
        snag.Panic("电池信息查询失败")
    }
    // 查询电池
    bat := ent.Database.Battery.Query().Where(battery.Sn(req.SN)).WithRider().WithCabinet().FirstX(s.ctx)
    var hb *xcpb.Heartbeat
    if len(r.Items[req.SN].Heartbeats) > 0 {
        hb = r.Items[req.SN].Heartbeats[0]
        return
    }

    detail = &model.XcBatteryDetail{}
    if hb != nil {
        detail = &model.XcBatteryDetail{
            UpdatedAt:       hb.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
            DisChargingTime: hb.DisChargingTime,
            XcBmsBattery:    model.NewXcBmsBattery(hb),
            Current:         hb.Current,
            Soh:             uint8(hb.Soh),
            Cycles:          uint16(hb.Cycles),
        }
    }

    detail.Sn = req.SN
    detail.CreatedAt = bat.CreatedAt.Format("2006-01-02 15:04:05")

    if bat.Edges.Cabinet != nil && bat.Ordinal != nil {
        detail.BelongsTo = bat.Edges.Cabinet.Name + "-" + strconv.Itoa(*bat.Ordinal) + "号仓"
    }

    if bat.Edges.Rider != nil {
        detail.BelongsTo = bat.Edges.Rider.Name + "-" + bat.Edges.Rider.Phone
    }

    return
}
