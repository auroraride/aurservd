// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-23, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

type EmployeeListReq struct {
	model.PaginationReq
	Status  *bool   `json:"status" query:"status"`   // 启用状态
	Keyword *string `json:"keyword" query:"keyword"` // 搜索关键词, 门店、手机号或姓名
	CityID  *uint64 `json:"cityId" query:"cityId"`   // 城市ID
}

type EmployeeListRes struct {
	ID            uint64       `json:"id"`
	Enable        bool         `json:"enable"`                  // 是否启用
	Name          string       `json:"name"`                    // 姓名
	Phone         string       `json:"phone"`                   // 电话
	City          model.City   `json:"city"`                    // 城市
	Group         *StoreGroup  `json:"group,omitempty"`         // 门店集合
	Stores        []*StoreInfo `json:"stores,omitempty"`        // 门店数据
	EmployeeStore *StoreInfo   `json:"employeeStore,omitempty"` // 当前上班门店
	Limit         uint         `json:"limit"`                   // 限制范围(m)
}
type StoreGroup struct {
	ID   uint64 `json:"id"`   // 集合ID
	Name string `json:"name"` // 集合名称
}

type StoreInfo struct {
	ID   uint64 `json:"id"`   // 门店ID
	Name string `json:"name"` // 门店名称
}

type EmployeeCreateReq struct {
	CityID   uint64   `json:"cityId" validate:"required" trans:"城市ID"`
	Name     string   `json:"name" validate:"required" trans:"姓名"`
	Phone    string   `json:"phone" validate:"required,phone" trans:"手机号"`
	Password string   `json:"password" validate:"required" trans:"密码"`
	GroupID  uint64   `json:"groupID"`  // 门店集合ID
	StoreIDs []uint64 `json:"storeIDs"` // 门店IDS
	Limit    uint     `json:"limit" trans:"限制范围"`
}

type EmployeeModifyReq struct {
	ID       *uint64  `json:"id" validate:"required" param:"id" trans:"店员ID"`
	CityID   *uint64  `json:"cityId"`
	Name     *string  `json:"name"`
	Phone    *string  `json:"phone"`
	Password *string  `json:"password"` // 密码
	GroupID  *uint64  `json:"groupID"`  // 门店集合ID
	StoreIDs []uint64 `json:"storeIDs"` //  门店IDS
	Limit    *uint    `json:"limit"`
	Enable   *bool    `json:"enable"` // 是否启用
}
