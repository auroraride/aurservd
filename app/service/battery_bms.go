// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"math"
	"strconv"
	"time"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/defs/batdef"
	"github.com/auroraride/adapter/defs/xcdef"
	"github.com/auroraride/adapter/rpc/pb"
	"github.com/auroraride/adapter/rpc/pb/timestamppb"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/baidu"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type batteryBmsService struct {
	*BaseService
	orm *ent.AssetClient
}

func NewBatteryBms(params ...any) *batteryBmsService {
	return &batteryBmsService{
		BaseService: newService(params...),
		orm:         ent.Database.Asset,
	}
}

func (s *batteryBmsService) Sync(data []*batdef.BatteryFlow) {
	go func() {
		bb, _ := jsoniter.Marshal(data)
		zap.L().Info("电池同步消息", zap.ByteString("sync-battery", bb))
	}()

	// 获取对应的battery rpc
	for _, bf := range data {
		// 获取电柜信息
		cab := NewCabinet().QueryOneSerial(bf.Serial)
		if cab == nil {
			zap.L().Error("未找到电柜信息: " + bf.Serial)
			continue
		}
		if bf.In != nil {
			// 放入电池
			_, _ = s.SyncPutin(bf.In.SN, cab, bf.Ordinal)
		}
		if bf.Out != nil {
			// 取出电池
			s.SyncPutout(cab, bf.Ordinal)
		}
	}
}

// SyncPutout 同步消息 - 从电柜中取出
func (s *batteryBmsService) SyncPutout(cab *ent.Cabinet, ordinal int) {
	// 查询电池
	bat, _ := s.orm.Query().Where(asset.Ordinal(ordinal), asset.LocationsID(cab.ID), asset.LocationsType(model.AssetLocationsTypeCabinet.Value())).First(s.ctx)
	if bat == nil {
		return
	}
	// 查询调拨
	// ent.Database.AssetTransferDetails.QueryNotDeleted().Where(assettransferdetails.AssetID(bat.ID))

	// _ = s.orm.Update().Where(battery.CabinetID(cab.ID), battery.Ordinal(ordinal)).ClearCabinetID().ClearOrdinal().Exec(s.ctx)
}

// SyncPutin 同步消息 - 放入电柜中
func (s *batteryBmsService) SyncPutin(sn string, cab *ent.Cabinet, ordinal int) (bat *ent.Asset, err error) {
	// 查询电池
	bat, err = NewBattery().QuerySn(sn)
	if err != nil {
		zap.L().Error("电池信息创建失败", zap.Error(err))
		return
	}

	// 查询电池有无调拨单 如果无调拨单 则更新电池信息
	// at, _ := NewAssetTransfer().QueryTransferByAssetID(context.Background(), bat.ID)
	// if at == nil {
	// 	zap.L().Error("电池调拨查询失败", zap.Error(err))
	// 	return
	// }

	// if time.Since(bat.UpdatedAt).Seconds() < 20 {
	// 	rid := ""
	// 	if bat.RiderID != nil {
	// 		rid = strconv.FormatUint(*bat.RiderID, 10)
	// 	}
	// 	zap.L().Error("电池解绑过快, sn=" + bat.Sn + ", updated_at=" + bat.UpdatedAt.Format("2006-01-02 15:04:05.000") + ", rider_id=" + rid + ", serial=" + cab.Serial + ", ordinal=" + strconv.Itoa(ordinal))
	// 	return
	// }

	err = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
		// 	updater := bat.Update().
		// 		SetCabinetID(cab.ID).
		// 		SetOrdinal(ordinal).
		// 		ClearRiderID().
		// 		ClearSubscribeID()
		//
		// 	if cab.StationID != nil {
		// 		// 当前电柜属于代理站点时, 设置新的站点信息
		// 		updater.SetNillableStationID(cab.StationID).SetNillableEnterpriseID(cab.EnterpriseID)
		// 	} else {
		// 		// 当前电柜属于平台时, 清除原有站点信息
		// 		updater.ClearStationID().ClearEnterpriseID()
		// 	}
		//
		// 	// 更新电池电柜信息
		// 	_, err = updater.Save(s.ctx)
		// 	if err != nil {
		// 		zap.L().Error("放入电柜更新电池失败", zap.Error(err))
		// 		return
		// 	}
		//
		// 更新电池流转
		NewBatteryFlow().Create(tx, bat, model.BatteryFlowCreateReq{
			CabinetID: silk.Pointer(cab.ID),
			Ordinal:   silk.Pointer(ordinal),
			Serial:    silk.Pointer(cab.Serial),
		})
		return
	})
	return
}

func (s *batteryBmsService) Detail(req *model.BatterySNRequest) (detail *model.BatteryBmsDetail) {
	ab, err := adapter.ParseBatterySN(req.SN)
	if err != nil {
		snag.Panic(err)
	}
	// 请求bms rpc
	r := rpc.BmsBatch(ab.Brand, &pb.BatteryBatchRequest{Sn: []string{req.SN}})
	if r == nil || r.Items[req.SN] == nil {
		snag.Panic("电池信息查询失败")
	}

	// 查询电池
	bat := ent.Database.Asset.Query().Where(asset.Sn(req.SN)).WithRider().WithCabinet().FirstX(s.ctx)
	if bat == nil {
		snag.Panic("电池未录入")
	}

	var (
		hb *pb.BatteryHeartbeat
		rb *pb.BatteryItem
	)

	if len(r.Items[req.SN].Heartbeats) > 0 {
		rb = r.Items[req.SN]
		hb = rb.Heartbeats[0]
	}

	detail = &model.BatteryBmsDetail{}
	if hb != nil {
		detail = &model.BatteryBmsDetail{
			UpdatedAt:            hb.CreatedAt.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05"),
			BmsBattery:           model.NewBmsBattery(hb),
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
			Online:               time.Since(hb.CreatedAt.AsTime().In(ar.TimeLocation)).Minutes() < 35,
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

	fr := rpc.BmsFaultOverview(ab.Brand, &pb.BatterySnRequest{Sn: req.SN})
	if fr != nil {
		detail.FaultsOverview = fr.Items
	}

	return
}

func (s *batteryBmsService) Statistics(req *model.BatterySNRequest) (detail *model.BatteryStatistics) {
	ab, err := adapter.ParseBatterySN(req.SN)
	if err != nil {
		snag.Panic(err)
	}
	// 请求xcbms rpc
	r := rpc.BmsStatistics(ab.Brand, &pb.BatterySnRequest{Sn: req.SN})
	if r == nil {
		snag.Panic("电池数据查询失败")
	}

	return &model.BatteryStatistics{
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

func (s *batteryBmsService) Position(req *model.BatteryPositionReq) (res *model.BatteryPositionRes) {
	ab, err := adapter.ParseBatterySN(req.SN)
	if err != nil {
		snag.Panic(err)
	}
	var start, end *timestamppb.Timestamp
	if req.Start != "" {
		reqStart := carbon.ParseByLayout(req.Start, carbon.DateTimeLayout).StdTime()
		start = timestamppb.New(reqStart)
	}
	if req.End != "" {
		reqEnd := carbon.ParseByLayout(req.End, carbon.DateTimeLayout).StdTime()
		end = timestamppb.New(reqEnd)
	}
	r := rpc.BmsPosition(ab.Brand, &pb.BatteryPositionRequest{
		Sn:    req.SN,
		Start: start,
		End:   end,
	})
	if r == nil {
		return &model.BatteryPositionRes{
			Positions:  make([]*model.BatteryPosition, 0),
			Stationary: make([]*model.BatteryStationary, 0),
		}
	}
	res = &model.BatteryPositionRes{
		Start:      r.Start.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05"),
		End:        r.End.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05"),
		Positions:  make([]*model.BatteryPosition, len(r.Positions)),
		Stationary: make([]*model.BatteryStationary, len(r.Stationary)),
	}
	for i, p := range r.Positions {
		res.Positions[i] = &model.BatteryPosition{
			InCabinet:  p.InCabinet,
			Stationary: p.Stationary,
			Soc:        p.Soc,
			Lng:        p.Lng,
			Lat:        p.Lat,
			Voltage:    p.Voltage,
			Gsm:        p.Gsm,
		}
		if p.At != nil {
			res.Positions[i].At = p.At.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05")
		}
	}
	for i, sa := range r.Stationary {
		res.Stationary[i] = &model.BatteryStationary{
			InCabinet: sa.InCabinet,
			Duration:  sa.Duration,
			StartSoc:  sa.StartSoc,
			EndSoc:    sa.EndSoc,
			Lng:       sa.Lng,
			Lat:       sa.Lat,
		}
		if sa.StartAt != nil {
			res.Stationary[i].StartAt = sa.StartAt.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05")
		}
		if sa.EndAt != nil {
			res.Stationary[i].EndAt = sa.EndAt.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05")
		}
	}
	return
}

func (s *batteryBmsService) FaultList(req *model.BatteryFaultReq) *model.PaginationRes {
	ab, err := adapter.ParseBatterySN(*req.SN)
	if err != nil {
		snag.Panic(err)
	}

	pq := &pb.BatteryFaultListRequest{
		Sn:    req.SN,
		Fault: req.Fault,
		Pagination: &pb.PaginationRequest{
			Current:  int64(req.Current),
			PageSize: int64(req.PageSize),
		},
	}
	if req.Start != "" {
		pq.BeginAt = timestamppb.New(tools.NewTime().ParseDateStringX(req.Start))
	}
	if req.End != "" {
		pq.EndAt = timestamppb.New(tools.NewTime().ParseNextDateStringX(req.End))
	}

	r := rpc.BmsFaultList(ab.Brand, pq)

	page := model.Pagination{
		Current: req.Current,
	}

	items := make([]*model.BatteryFaultRes, 0)
	if r != nil {
		page.Pages = int(r.Pagination.Pages)
		page.Total = int(r.Pagination.Total)

		items = make([]*model.BatteryFaultRes, len(r.Items))
		for i, item := range r.Items {
			items[i] = &model.BatteryFaultRes{
				Sn:      item.Sn,
				Fault:   item.Fault,
				BeginAt: item.BeginAt.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05"),
			}
			if item.EndAt != nil {
				items[i].EndAt = item.EndAt.AsTime().In(ar.TimeLocation).Format("2006-01-02 15:04:05")
			}
		}

	}
	return &model.PaginationRes{
		Pagination: page,
		Items:      items,
	}
}

// TrackRectify 轨迹纠偏
func (s *batteryBmsService) TrackRectify(req *model.BatteryTrackReq) (res *model.BatteryTrackRes, err error) {
	r := baidu.Track{}
	for _, v := range req.Points {
		r = append(r, baidu.Point{
			LocTime:        v.LocTime,
			Longitude:      v.Longitude,
			Latitude:       v.Latitude,
			CoordTypeInput: baidu.MapCoordType,
		})
	}

	rectify, err := baidu.NewMap().TrackRectify(r)
	if err != nil {
		return nil, err
	}
	return rectify, nil
}
