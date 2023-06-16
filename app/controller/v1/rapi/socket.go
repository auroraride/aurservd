// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/app/socket"
)

type socketapi struct{}

var Socket = new(socketapi)

func (*socketapi) Rider(c echo.Context) (err error) {
	srv := service.NewRiderSocket()
	return socket.Wrap(c, srv)
}
