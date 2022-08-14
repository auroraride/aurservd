// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    PersonUnauthenticated      PersonAuthStatus = iota // 未认证
    PersonAuthPending                                  // 认证中
    PersonAuthenticated                                // 已认证
    PersonAuthenticationFailed                         // 认证失败
)

type PersonAuthStatus uint8

func (s PersonAuthStatus) String() string {
    switch s {
    case PersonUnauthenticated:
        return "未认证"
    case PersonAuthPending:
        return "认证中"
    case PersonAuthenticated:
        return "已认证"
    }
    return "认证失败"
}

func (s PersonAuthStatus) Value() uint8 {
    return uint8(s)
}

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

// RequireAuth 是否需要认证
func (s PersonAuthStatus) RequireAuth() bool {
    switch s {
    case PersonAuthPending, PersonUnauthenticated, PersonAuthenticationFailed:
        return true
    }
    return false
}

// PersonBanReq 封禁或解封骑手身份
type PersonBanReq struct {
    ID  uint64 `json:"id" ` // 骑手ID
    Ban bool   `json:"ban"` // `true`封禁 `false`解封
}

type Person struct {
    // 证件号码
    IDCardNumber string `json:"id_card_number,omitempty"`
    // 证件人像面
    IDCardPortrait string `json:"id_card_portrait,omitempty"`
    // 证件国徽面
    IDCardNational string `json:"id_card_national,omitempty"`
    // 实名认证人脸照片
    AuthFace string `json:"auth_face,omitempty"`
}
