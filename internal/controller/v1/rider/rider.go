// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package rider

import (
    "github.com/labstack/echo/v4"
)

var ApiRider = new(apiRider)

type apiRider struct {
}

type signupReq struct {
    Phone string `json:"phone"`
    Sms   string `json:"sms"`
}

func (*apiRider) Signup(c echo.Context) {
}
