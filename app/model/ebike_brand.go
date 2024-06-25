// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type EbikeBrand struct {
	ID             uint64                 `json:"id" bson:"id"`
	Name           string                 `json:"name" bson:"name"`                       // 名称
	Cover          string                 `json:"cover,omitempty" bson:"cover,omitempty"` // 封面图
	MainPic        []string               `json:"mainPic"`                                // 主图
	BrandAttribute []*EbikeBrandAttribute `json:"brandAttribute"`                         // 品牌属性
}

type EbikeBrandCreateReq struct {
	Name           string                `json:"name" validate:"required"`    // 名称
	Cover          string                `json:"cover" validate:"required"`   // 封面图
	MainPic        []string              `json:"mainPic" validate:"required"` // 主图
	BrandAttribute []EbikeBrandAttribute `json:"brandAttribute"`              // 品牌属性
}

type EbikeBrandModifyReq struct {
	ID             uint64                `json:"id" param:"id" validate:"required"`
	Name           string                `json:"name"`                        // 名称
	Cover          string                `json:"cover"`                       // 封面图
	MainPic        []string              `json:"mainPic" validate:"required"` // 主图
	BrandAttribute []EbikeBrandAttribute `json:"brandAttribute" `             // 品牌属性
}
