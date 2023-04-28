// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package workwx

import (
	"github.com/auroraride/aurservd/internal/ar"
)

type GroupCreateReponse struct {
	baseResponse
	Chatid string `json:"chatid"`
}

// CreateGroup 创建群聊
func (w *Client) CreateGroup(name, owner, chatid string, users []string) error {
	var res GroupCreateReponse
	return w.RequestPost("/appchat/create", ar.Map{
		"name":     name,
		"owner":    owner,
		"userlist": users,
		"chatid":   chatid,
	}, &res)
}
