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

func (s *stockEbikeService) Loopers(req *model.StockTransferReq, enterpriseId uint64) (looppers []model.StockTransferLoopper, failed []string) {
	failed = make([]string, 0)
	// 查询电车信息
	eq := ent.Database.Ebike.Query()
	// 调拨至门店或站点调拨 都要求电车未分配
	if req.InboundTarget == model.StockTargetStore || req.InboundTarget == model.StockTargetStation {
		eq = NewEbike().AllocatableBaseFilter()
	} else {
		// 调拨到平台则状态为 库存中 / 维修中 / 已报废
		eq.Where(
			ebike.StatusIn(model.EbikeStatusInStock, model.EbikeStatusMaintenance, model.EbikeStatusScrapped),
			ebike.Or(
				ebike.StoreIDNotNil(),
				ebike.StationIDNotNil(),
			),
		)
	}

	sns := make(map[string]bool)
	for _, e := range req.Ebikes {
		sns[e] = true
	}

	bikes, err := eq.Where(ebike.SnIn(req.Ebikes...)).WithBrand().WithEnterprise().All(s.ctx)
	if err != nil {
		snag.Panic(err)
	}

	// 车辆信息
	for _, bike := range bikes {
		delete(sns, bike.Sn)
		// 当调出为站点 调入为站点 需要判断调出站点和调入站点是否为同一企业 不同企业不能调拨
		if req.InboundTarget == model.StockTargetStation && req.OutboundTarget == model.StockTargetStation &&
			bike.Edges.Enterprise != nil && bike.Edges.Enterprise.ID != enterpriseId {
			failed = append(failed, fmt.Sprintf("电车调拨失败，[%s]不属于当前团签", bike.Sn))
			continue
		}

		looppers = append(looppers, model.StockTransferLoopper{
			EbikeSN:   silk.String(bike.Sn),
			EbikeID:   silk.UInt64(bike.ID),
			Message:   bike.Sn,
			BrandID:   silk.UInt64(bike.BrandID),
			BrandName: silk.String(bike.Edges.Brand.Name),
		})
	}

	// 未找到的车架号
	for k := range sns {
		failed = append(failed, fmt.Sprintf("未找到: %s", k))
	}

	return
}
