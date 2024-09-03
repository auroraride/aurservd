// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-12, by aurb

package wapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type warestore struct{}

var Warestore = new(warestore)

// Signin
// @ID		WarestoreSignin
// @Router	/warestore/v2/signin [POST]
// @Summary	登录
// @Tags	Warestore - 仓管接口
// @Accept	json
// @Produce	json
// @Param	body	body		definition.WarestorePeopleSigninReq	true	"登录请求"
// @Success	200		{object}	definition.WarestorePeopleSigninRes	"请求成功"
func (*warestore) Signin(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.WarestorePeopleSigninReq](c)
	return ctx.SendResponse(biz.NewWarestore().Signin(req))
}

// Duty
// @ID		WarestoreDuty
// @Router	/warestore/v2/duty [POST]
// @Summary	上班
// @Tags	Warestore - 仓管接口
// @Accept	json
// @Produce	json
// @Param	body	body		definition.WarestoreDutyReq	true	"登录请求"
// @Success	200		{object}	model.StatusResponse		"请求成功"
func (*warestore) Duty(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.WarestoreDutyReq](c)
	return ctx.SendResponse(biz.NewWarestore().Duty(definition.AssetSignInfo{
		AssetManager: ctx.AssetManager,
		Employee:     ctx.Employee,
	}, req))
}
