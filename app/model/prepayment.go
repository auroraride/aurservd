// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-04
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type PrepaymentOverview struct {
	Times   int     `json:"times"`   // 充值次数
	Balance float64 `json:"balance"` // 当前余额
	Amount  float64 `json:"amount"`  // 总充值金额
	Cost    float64 `json:"cost"`    // 已使用金额
}

type PrepaymentListReq struct {
	PaginationReq
	Payway       Payway `json:"payway" query:"payway" enums:"0,1,2" validate:"lte=3,gte=0"` // 支付方式 0:全部 1:现金 2:微信
	Start        string `json:"start" query:"start"`                                        // 开始日期
	End          string `json:"end" query:"end"`                                            // 结束日期
	EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId"`                          // 团签ID
}

type PrepaymentListRes struct {
	Amount  float64 `json:"amount"`                                    // 金额
	Name    string  `json:"name"`                                      // 操作人
	Time    string  `json:"time"`                                      // 时间
	Remark  string  `json:"remark"`                                    // 备注信息
	Payway  Payway  `json:"payway" enums:"1,2" validate:"gte=1,lte=2"` // 支付方式 1:现金 2:微信
	TradeNo *string `json:"tradeNo,omitempty"`                         // 支付平台交易单号
}
