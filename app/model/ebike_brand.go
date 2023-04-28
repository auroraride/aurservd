// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type EbikeBrand struct {
	ID    uint64 `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`                       // 名称
	Cover string `json:"cover,omitempty" bson:"cover,omitempty"` // 封面图
}

type EbikeBrandCreateReq struct {
	Name  string `json:"name" validate:"required" validate:"名称"`   // 名称
	Cover string `json:"cover" validate:"required" validate:"封面图"` // 封面图
}

type EbikeBrandModifyReq struct {
	ID    uint64 `json:"id" param:"id" validate:"required"`
	Name  string `json:"name"`  // 名称
	Cover string `json:"cover"` // 封面图
}
