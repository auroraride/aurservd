// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

// CreatePersonAccountReq 创建个人签署账号请求体
type CreatePersonAccountReq struct {
    ThirdPartyUserId string `json:"thirdPartyUserId,omitempty"`
    Name             string `json:"name,omitempty"`
    IdType           string `json:"idType,omitempty"`
    IdNumber         string `json:"idNumber,omitempty"`
    Mobile           string `json:"mobile,omitempty"`
    Email            string `json:"email,omitempty"`
}

// CreatePersonAccountRes 创建个人账号返回体
type CreatePersonAccountRes struct {
    AccountId *string `json:"accountId"`
}

// CreatePersonAccount 创建个人签署账号
func (e *Esign) CreatePersonAccount(req CreatePersonAccountReq) *string {
    res := new(CreatePersonAccountRes)
    e.request(createPersonAccountUrl, "POST", req, res)
    return res.AccountId
}
