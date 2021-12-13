// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    PersonAuthNot PersonAuthStatus = iota
    PersonAuthPending
    PersonAuthSuccess
    PersonAuthFailed
)

type FaceVerifyResult struct {
    Birthday       string  `json:"birthday"`
    IssueAuthority string  `json:"issueAuthority"`
    Address        string  `json:"address"`
    Gender         string  `json:"gender"`
    Nation         string  `json:"nation"`
    ExpireTime     string  `json:"expireTime"`
    Name           string  `json:"name"`
    IssueTime      string  `json:"issueTime"`
    IdCardNumber   string  `json:"idCardNumber"`
    Score          float64 `json:"score"`
    LivenessScore  float64 `json:"livenessScore"`
    Spoofing       float64 `json:"spoofing"`
}

type PersonAuthStatus uint8

func (s PersonAuthStatus) String() string {
    switch s {
    case PersonAuthNot:
        return "未认证"
    case PersonAuthPending:
        return "认证中"
    case PersonAuthSuccess:
        return "已认证"
    }
    return "认证失败"
}

func (s PersonAuthStatus) Raw() uint8 {
    return uint8(s)
}

// RequireAuth 是否需要认证
func (s PersonAuthStatus) RequireAuth() bool {
    switch s {
    case PersonAuthNot, PersonAuthFailed:
        return true
    }
    return false
}
