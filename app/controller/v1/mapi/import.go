// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-18
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/labstack/echo/v4"
)

type importApi struct{}

var Import = new(importApi)

// Rider
// @ID		ManagerImportRider
// @Router	/manager/v1/import/rider [POST]
// @Summary	ME001 单个导入骑手
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.ImportRiderCreateReq	true	"骑手信息"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*importApi) Rider(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ImportRiderCreateReq](c)
	err = service.NewImportRiderWithModifier(ctx.Modifier).Create(req)
	if err != nil {
		snag.Panic(err)
	}
	return ctx.SendResponse()
}
