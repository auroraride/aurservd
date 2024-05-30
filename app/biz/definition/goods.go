// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-29, by Jorjan

package definition

import (
	"github.com/auroraride/aurservd/app/model"
)

type GoodsType uint8

const (
	GoodsTypeEbike GoodsType = iota + 1 // 电车
)

func (t GoodsType) Value() uint8 {
	return uint8(t)
}

type GoodsStatus uint8

const (
	GoodsStatusOffline GoodsStatus = iota // 下架
	GoodsStatusOnline                     // 上架
)

func (t GoodsStatus) Value() uint8 {
	return uint8(t)
}

// GoodsListReq 列表请求
type GoodsListReq struct {
	model.PaginationReq
	Keyword *string      `json:"keyword" query:"keyword"` // 关键字
	Status  *GoodsStatus `json:"status" query:"status"`   // 是否启用

	Start *string `json:"start" query:"start"` // 开始时间
	End   *string `json:"end" query:"end"`     // 结束时间
}

// Goods 商品公共字段
type Goods struct {
	Sn        string        `json:"sn"`        // 商品编号
	Name      string        `json:"name"`      // 商品名称
	Type      GoodsType     `json:"type"`      // 商品类型 1-电车
	Lables    []string      `json:"lables"`    // 商品标签
	Price     float64       `json:"price"`     // 商品价格
	Weight    int           `json:"weight"`    // 商品权重
	HeadPic   string        `json:"headPic"`   // 商品头图
	Photos    []string      `json:"photos"`    // 商品图片
	Intro     []string      `json:"intro"`     // 商品介绍
	Stores    []model.Store `json:"stores"`    // 配置店铺信息
	CreatedAt string        `json:"createdAt"` // 创建时间
	Status    GoodsStatus   `json:"status"`    // 商品状态 0-已下架 1-已上架
	Remark    string        `json:"remark"`    // 备注
}

// GoodsDetail 商品详情
type GoodsDetail struct {
	ID uint64 `json:"id"`
	Goods
}

// GoodsCreateReq 创建
type GoodsCreateReq struct {
	Name     string    `json:"name" validate:"required" trans:"商品名称"`      // 商品名称
	Type     GoodsType `json:"type" validate:"required" trans:"商品类别"`      // 商品类别 1电车
	Lables   []string  `json:"lables"`                                     // 商品标签
	Price    float64   `json:"price" validate:"required" trans:"商品价格"`     // 商品价格
	Weight   int       `json:"weight" validate:"required" trans:"商品权重"`    // 商品权重
	HeadPic  string    `json:"headPic" validate:"required" trans:"商品头图"`   // 商品头图
	Photos   []string  `json:"photos" validate:"max=5" trans:"商品图片"`       // 商品图片
	Intro    []string  `json:"intro" validate:"max=5" trans:"商品介绍"`        // 商品介绍
	StoreIds []uint64  `json:"storeIds" validate:"required" trans:"门店IDS"` // 门店IDS
	Remark   string    `json:"remark"`                                     // 备注
}

// GoodsModifyReq 商品修改请求
type GoodsModifyReq struct {
	model.IDParamReq
	GoodsCreateReq
}

// GoodsUpdateStatusReq 商品更新状态请求
type GoodsUpdateStatusReq struct {
	model.IDParamReq
	Status GoodsStatus `json:"status"` // 商品状态 0下架 1上架
}
