// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

// 优惠券码最佳使用算法, 贪心算法参考
// https://blog.csdn.net/qq_44112474/article/details/123616038
// https://segmentfault.com/q/1010000015942183/a-1020000015946180

package model

var (
    CouponShuffle = []int{9, 0, 1, 4, 3, 5, 6, 8, 7, 2}
)

type CouponGenerateReq struct {
    TemplateID uint64   `json:"templateId" validate:"required" trans:"优惠券模板"`
    Phones     []string `json:"phones" validate:"required_without=Number" trans:"骑手电话"`
    Number     int      `json:"number" validate:"required_without=Phones" trans:"数量"`
    Amount     float64  `json:"amount" validate:"required" trans:"金额"`
    Remark     string   `json:"remark" validate:"required" trans:"备注"`
}
