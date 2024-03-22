// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type setting struct{}

var Setting = new(setting)

// LatestVersion
// @ID		SettingLatestVersion
// @Router	/rider/v2/setting/version [GET]
// @Summary	获取最新版本
// @Tags	Setting - 设置
// @Accept	json
// @Produce	json
// @Param	query	query		definition.LatestVersionReq	true	"请求参数"
// @Success	200		{object}	definition.VersionRes		"请求成功"
func (*setting) LatestVersion(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.LatestVersionReq](c)
	return ctx.SendResponse(biz.NewVersion().LatestVersion(req))
}
