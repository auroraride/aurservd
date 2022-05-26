// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package payment

func Boot() {
    _alipay = &alipayClient{
        Client:       newAlipayClient(),
        notifyClient: newNotifyClient(),
    }
    _wechat = newWechatClient()
}
