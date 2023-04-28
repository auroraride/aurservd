// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-26
// Based on aurservd by liasica, magicrolan@qq.com.

package kit

import (
	"context"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/labstack/echo/v4"
)

type cabinet struct{}

var Cabinet = new(cabinet)

func (*cabinet) Name(c echo.Context) (err error) {
	ctx := app.Context(c)
	serial := c.Param("serial")
	name, _ := ar.Redis.HGet(context.Background(), ar.CabinetNameCacheKey, serial).Result()

	return ctx.SendResponse(map[string]string{"name": name})
}
