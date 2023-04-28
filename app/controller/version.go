// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-30
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

import (
	"fmt"
	"strings"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

type version struct{}

var Version = new(version)

func (*version) Get(c echo.Context) (err error) {
	ctx := app.Context(c)
	plaform := ctx.QueryParam("plaform")
	app := ctx.QueryParam("app")
	if app == "" {
		app = "rider"
	}
	key := fmt.Sprintf("%s-%s", strings.ToLower(app), strings.ToLower(plaform))
	set := service.NewSetting().GetSetting(model.SettingAppVersionKey).(map[string]interface{})
	v, ok := set[key]
	if !ok {
		v = model.SettingAppVersionValue{}
	}
	return ctx.SendResponse(v)
}
