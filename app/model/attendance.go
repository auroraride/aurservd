// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type AttendancePrecheck struct {
	Duty    bool     `json:"duty"` // 上下班 `true`上班 `false`下班
	SN      *string  `json:"sn" validate:"required" trans:"门店编号"`
	Lat     *float64 `json:"lat" validate:"required" trans:"纬度"`
	Lng     *float64 `json:"lng" validate:"required" trans:"经度"`
	Address *string  `json:"address" validate:"required" trans:"详细地址"`
}

type AttendanceCreateReq struct {
	Photo     *string        `json:"photo"`                                        // 上班照片
	Inventory map[string]int `json:"inventory" validate:"required" trans:"物资盘点清单"` // 格式为 [名称]:数量
	*AttendancePrecheck
}

type AttendanceInventory struct {
	Name     string `json:"name"`            // 物资名称
	Num      int    `json:"num"`             // 物资数量
	StockNum int    `json:"stockNum"`        // 当前库存该物资数量
	Model    string `json:"model,omitempty"` // 电池型号
}

type AttendanceListReq struct {
	PaginationReq

	Keyword *string `json:"keyword" query:"keyword"`         // 店员姓名或手机号
	Duty    uint8   `json:"duty" enums:"0,1,2" query:"duty"` // 考勤分类 0:全部 1:上班 2:下班
	Start   *string `json:"start" query:"start"`             // 考勤开始日期
	End     *string `json:"end" query:"end"`                 // 考勤结束日期
	CityID  *uint64 `json:"cityId" query:"cityId"`           // 城市ID
	StoreID *uint64 `json:"storeId" query:"storeId"`         // 门店ID
}

type AttendanceListRes struct {
	ID        uint64                `json:"id"`
	City      City                  `json:"city"`      // 城市
	Store     Store                 `json:"store"`     // 门店
	Duty      bool                  `json:"duty"`      // true上班 或 false下班打卡
	Name      string                `json:"name"`      // 店员姓名
	Phone     string                `json:"phone"`     // 店员电话
	Time      string                `json:"time"`      // 盘点时间
	Photo     string                `json:"photo"`     // 照片
	Inventory []AttendanceInventory `json:"inventory"` // 物资盘点情况, 当物资异常的时候需要显示实际库存量
}
