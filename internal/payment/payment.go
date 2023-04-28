// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

import (
	"net/http"

	"github.com/auroraride/aurservd/app/model"
)

type Payment interface {
	AppPay(prepay *model.PaymentCache) (string, error)
	Native(pc *model.PaymentCache) (string, error)
	Notification(req *http.Request) *model.PaymentCache
	Refund(req *model.PaymentRefund)
}

func Boot() {
	_alipay = &alipayClient{
		Client: newAlipayClient(),
		// refundClient: newAlipayRefundClient(),
	}
	_wechat = newWechatClient()
}
