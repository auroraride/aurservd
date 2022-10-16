// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-16
// Based on aurservd by liasica, magicrolan@qq.com.

package esign

type Notice struct {
    Action              string `json:"action,omitempty"`
    FlowId              string `json:"flowId,omitempty"`
    AccountId           string `json:"accountId,omitempty"`
    AuthorizedAccountId string `json:"authorizedAccountId,omitempty"`
    SignTime            string `json:"signTime,omitempty"`
    Order               int    `json:"order,omitempty"`
    SignResult          int    `json:"signResult,omitempty"`
    ThirdOrderNo        string `json:"thirdOrderNo,omitempty"`
    ResultDescription   string `json:"resultDescription,omitempty"`
    Timestamp           int64  `json:"timestamp,omitempty"`
    ThirdPartyUserId    string `json:"thirdPartyUserId,omitempty"`
}
