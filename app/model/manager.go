// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// ManagerSigninReq 管理员登录请求
type ManagerSigninReq struct {
    Phone    string `json:"phone" validate:"required,phone"`
    Password string `json:"password" validate:"required"`
}

// ManagerSigninRes 管理员登录返回
type ManagerSigninRes struct {
    Id    uint64 `json:"id"`
    Token string `json:"token"`
}

// ManagerAddReq 管理员新增
type ManagerAddReq struct {
    ManagerSigninReq
    Name string `json:"name" validate:"required" trans:"姓名"`
}
