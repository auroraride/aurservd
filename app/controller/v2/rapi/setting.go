// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type setting struct{}

var Setting = new(setting)

// LatestVersion
// TODO: 未完成
// @ID		SettingLatestVersion
// @Router	/rider/v2/setting/version [GET]
// @Summary	获取最新版本
// @Tags	Setting - 设置
// @Accept	json
// @Produce	json
// @Param	query	query		definition.VersionReq	true	"请求参数"
// @Success	200		{object}	definition.VersionRes	"请求成功"
func (*setting) LatestVersion(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(&definition.VersionRes{
		Version: definition.Version{
			Version: "2.0.0",
			Content: "更新内容:\n1. 修复了一些bug\n2. 优化了一些功能\n3. 增加了一些新功能",
			Force:   false,
		},
	})
}