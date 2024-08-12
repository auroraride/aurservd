package service

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

type assetTransferDetailService struct {
	orm *ent.AssetTransferDetailsClient
}

func NewAssetTransferDetail() *assetTransferDetailService {
	return &assetTransferDetailService{
		orm: ent.Database.AssetTransferDetails,
	}
}

// TransferDetailCount 调拨记录详情的资产统计
func (s *assetTransferDetailService) TransferDetailCount(items map[string]*model.TransferAssetDetail, key string, isIn bool) {
	if _, ok := items[key]; !ok {
		items[key] = &model.TransferAssetDetail{
			Name:     key,
			Outbound: 0,
			Inbound:  0,
		}
	}

	// 资产进入一次就计入一次出库数量
	items[key].Outbound += 1
	// 资产是否已入库
	if isIn {
		items[key].Inbound += 1
	}
}
