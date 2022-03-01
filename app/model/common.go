// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// Modifier 修改人
type Modifier struct {
    ID    uint64 `json:"id,omitempty"`
    Name  string `json:"name,omitempty"`
    Phone string `json:"phone,omitempty"`
}

type AliyunOssStsRes struct {
    AccessKeySecret string `json:"accessKeySecret,omitempty"`
    Expiration      string `json:"expiration,omitempty"`
    AccessKeyId     string `json:"accessKeyId,omitempty"`
    StsToken        string `json:"stsToken,omitempty"`
    Bucket          string `json:"bucket,omitempty"`
    Region          string `json:"region,omitempty"`
}
