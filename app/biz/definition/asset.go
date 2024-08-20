// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-20, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

type ReqType uint8

const (
	ReqTypeBattery ReqType = iota + 1 // 电池
	ReqTypeEbike                      // 电车
)

// AssetListReq 资产列表请求
type AssetListReq struct {
	model.PaginationReq
	WareHousID *uint64            `json:"wareHousID" query:"wareHousID"`           // 仓库ID
	StoreID    *uint64            `json:"storeID" query:"storeID"`                 // 门店ID
	Status     *model.AssetStatus `json:"status" query:"status" enums:"1,2,3,4,5"` // 资产状态 0:待入库 1:库存中 2:配送中 3:使用中 4:故障 5:报废
	ModelID    *uint64            `json:"modelID" query:"modelID"`                 // 电池型号ID
	BrandID    *uint64            `json:"brandId" query:"brandId"`                 // 电车品牌ID
	SN         *string            `json:"sn" query:"sn"`                           // 编号(电池编号或车架号或车牌号)
	Type       ReqType            `json:"type" query:"type" enums:"1,2"`           // 请求资产类型 1:电池 2:电车
}
