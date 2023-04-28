// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/app/socket"
	"github.com/labstack/echo/v4"
)

type socketapi struct{}

var Socket = new(socketapi)

func (*socketapi) Employee(c echo.Context) error {
	srv := service.NewEmployeeSocket()
	return socket.Wrap(c, srv)
}
