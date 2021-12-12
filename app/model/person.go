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

// RequireAuth 是否需要认证
func (s PersonAuthStatus) RequireAuth() bool {
    switch s {
    case PersonAuthNot, PersonAuthFailed:
        return true
    }
    return false
}