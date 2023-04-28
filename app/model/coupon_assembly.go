// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type CouponAssemblyListReq struct {
	PaginationReq
	TemplateID uint64 `json:"templateId" query:"templateId"` // 模板
}

type CouponAssembly struct {
	ID      uint64             `json:"id"`
	Name    string             `json:"name"`               // 名称
	Target  uint8              `json:"target" enums:"1,2"` // 对象, 1:骑手 2:库存
	Amount  float64            `json:"amount"`             // 金额
	Number  int                `json:"number"`             // 数量
	Creator *Modifier          `json:"creator"`            // 操作人
	Time    string             `json:"time"`               // 时间
	Remark  string             `json:"remark"`             // 备注
	Meta    CouponTemplateMeta `json:"meta"`               // 详细信息
}
