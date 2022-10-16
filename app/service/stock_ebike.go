// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
)

type stockEbikeService struct {
    *BaseService
    orm *ent.StockClient
}

func NewStockEbike(params ...any) *stockEbikeService {
    return &stockEbikeService{
        BaseService: newService(params...),
        orm:         ent.Database.Stock,
    }
}

func (s *stockEbikeService) Loopers(req *model.StockTransferReq) (looppers []model.StockTransferLoopper, failed []string) {
    failed = make([]string, 0)
    // 查询电车信息
    eq := ent.Database.Ebike.Query()
    if req.InboundTarget == model.StockTargetStore {
        eq.Where(
            ebike.Status(model.EbikeStatusInStock),
            ebike.Enable(true),
            ebike.PlateNotNil(),
            ebike.MachineNotNil(),
            ebike.SimNotNil(),
            ebike.StoreIDIsNil(),
        )
    } else {
        // 调拨到平台则状态为 库存中 / 维修中 / 已报废
        eq.Where(
            ebike.StatusIn(model.EbikeStatusInStock, model.EbikeStatusMaintenance, model.EbikeStatusScrapped),
            ebike.StoreIDNotNil(),
        )
    }

    sns := make(map[string]bool)
    for _, e := range req.Ebikes {
        sns[e] = true
    }

    bikes, err := eq.Where(ebike.SnIn(req.Ebikes...)).All(s.ctx)
    if err != nil {
        snag.Panic(err)
    }

    // 车辆信息
    for _, bike := range bikes {
        delete(sns, bike.Sn)

        looppers = append(looppers, model.StockTransferLoopper{
            EbikeSN: silk.String(bike.Sn),
            EbikeID: silk.UInt64(bike.ID),
            Message: bike.Sn,
        })
    }

    // 未找到的车架号
    for k := range sns {
        failed = append(failed, fmt.Sprintf("未找到: %s", k))
    }

    return
}
