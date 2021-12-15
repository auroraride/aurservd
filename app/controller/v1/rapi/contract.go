// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import "github.com/labstack/echo/v4"

type contract struct {
}

var Contract = new(contract)

func (*contract) Sign(c echo.Context) error {
    return nil
}
