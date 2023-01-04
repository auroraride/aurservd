// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/14
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

import (
    "fmt"
    jsoniter "github.com/json-iterator/go"
)

// PersonAccountReq 创建个人签署账号请求体
type PersonAccountReq struct {
    ThirdPartyUserId string `json:"thirdPartyUserId,omitempty"`
    Name             string `json:"name,omitempty"`
    IdType           string `json:"idType,omitempty"`
    IdNumber         string `json:"idNumber,omitempty"`
    Mobile           string `json:"mobile,omitempty"`
    Email            string `json:"email,omitempty"`
}

// CreatePersonAccountRes 创建个人账号返回体
type CreatePersonAccountRes struct {
    AccountId        string `json:"accountId,omitempty"`
    ThirdPartyUserId string `json:"thirdPartyUserId,omitempty"`
    Name             string `json:"name,omitempty"`
    IdType           string `json:"idType,omitempty"`
    IdNumber         string `json:"idNumber,omitempty"`
    Mobile           string `json:"mobile,omitempty"`
    Email            string `json:"email,omitempty"`
}

// CreatePersonAccount 创建个人签署账号
func (e *Esign) CreatePersonAccount(req PersonAccountReq) string {
    res := new(CreatePersonAccountRes)
    e.request(createPersonAccountUrl, "POST", req, res)
    return res.AccountId
}

func (e *Esign) ModifyAccount(id string, req PersonAccountReq) string {
    res := new(CreatePersonAccountRes)
    e.request(fmt.Sprintf(modifyAccountUrl, id), "PUT", req, res)
    b, _ := jsoniter.Marshal(res)
    fmt.Printf("修改骑手信息: %s\n", b)
    return res.AccountId
}
