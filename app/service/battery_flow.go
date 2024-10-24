// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/rpc/pb"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/internal/ent"
)

type batteryFlowService struct {
	*BaseService

	orm *ent.BatteryFlowClient
}

func NewBatteryFlow(params ...any) *batteryFlowService {
	return &batteryFlowService{
		BaseService: newService(params...),

		orm: ent.Database.BatteryFlow,
	}
}

func (s *batteryFlowService) Create(tx *ent.Tx, bat *ent.Asset, req model.BatteryFlowCreateReq) {
	updater := tx.BatteryFlow.Create().
		SetSn(bat.Sn).
		SetBatteryID(bat.ID).
		SetNillableRiderID(req.RiderID).
		SetNillableSubscribeID(req.SubscribeID).
		SetNillableCabinetID(req.CabinetID).
		SetNillableOrdinal(req.Ordinal).
		SetNillableSerial(req.Serial)

	sr := rpc.BmsSample(adapter.BatteryBrand(bat.BrandName), &pb.BatterySnRequest{Sn: bat.Sn})
	if sr != nil {
		updater.SetSoc(float64(sr.Soc)).SetGeom(adapter.NewGeometry(sr.Geom).WGS84toGCJ02())
	}
	err := updater.Exec(s.ctx)
	if err != nil {
		zap.L().Error("电池流转创建失败", zap.Error(err))
	}
}

// func (s *batteryFlowService) QueryFromBin(cabinetID uint64, ordinal int, sn string) {
//
// }
//
// func (s *batteryFlowService) Sync(cab *ent.Cabinet, old map[int]*model.CabinetBin, b *cabdef.Bin) {
//     ob, ok := old[b.Ordinal-1]
//     // 如果旧仓位信息未找到, 则需要初始化一个仓位信息
//     if !ok {
//         ob = &model.CabinetBin{BatterySN: ""}
//     }
//
//     // 如果新旧仓位电池编号相等, 则直接跳过
//     if ob.BatterySN == b.BatterySn {
//         return
//     }
//
//     // 放入 (新有旧无)
//     if ob.BatterySN == "" {
//         bat, err = NewBattery().LoadOrCreate(sn, &model.BatteryInCabinet{
//             CabinetID: cabinetID,
//             Ordinal:   ordinal,
//         })
//         s.orm.Create().
//             SetSn(b.BatterySn)
//     }
//
//     // 取出 (旧有新无)
//
//     // 替换 (旧有新有)
// }
//
// func (s *batteryFlowService) doFlow() {
//
// }
