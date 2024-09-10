// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-22, by aurb

package definition

type StoreGroupCreateRep struct {
	Name string `json:"name" validate:"required"` // 门店集合名称
}

type StoreGroupListRes struct {
	ID   uint64 `json:"id"`   // 门店集合id
	Name string `json:"name"` // 门店集合名称
}
