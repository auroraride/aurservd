// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type setting struct{}

var Setting = new(setting)

// List
//	@ID			ManagerSettingList
//	@Router		/manager/v1/setting [GET]
//	@Summary	M1010 列举设置
//	@Tags		[M]管理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Manager-Token	header		string				true	"管理员校验token"
//	@Success	200				{object}	[]model.SettingReq	"请求成功"
func (*setting) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)

	return ctx.SendResponse(service.NewSetting().List())
}

// Modify
//	@ID			ManagerSettingModify
//	@Router		/manager/v1/setting/{key} [PUT]
//	@Summary	M1011 调整设置
//	@Tags		[M]管理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Manager-Token	header		string					true	"管理员校验token"
//	@Param		key				path		string					true	"设置项"
//	@Param		body			body		model.SettingReq		true	"desc"
//	@Success	200				{object}	model.StatusResponse	"请求成功"
func (*setting) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SettingReq](c)
	service.NewSettingWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// LegalRead
//	@ID			ManagerSettingLegalRead
//	@Router		/manager/v1/setting/legal/{name} [GET]
//	@Summary	M1016 获取法规
//	@Tags		[M]管理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Manager-Token	header		string			true	"管理员校验token"
//	@Param		name			path		string			true	"法规名称, policy: APP隐私政策; agreement: APP服务协议; agent-policy: 代理端小程序隐私政策; agent-agreement: 代理端小程序服务协议; promote-policy: 推广端小程序隐私政策; promote-agreement: 推广端小程序服务协议"
//	@Success	200				{object}	model.LegalRes	"请求成功"
func (*setting) LegalRead(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.LegalName](c)
	return ctx.SendResponse(service.NewLegal().Read(req))
}

// LegalSave
//	@ID			ManagerSettingLegalSave
//	@Router		/manager/v1/setting/legal [POST]
//	@Summary	M1017 保存法规
//	@Tags		[M]管理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Manager-Token	header		string					true	"管理员校验token"
//	@Param		body			body		model.LegalSaveReq		true	"请求参数"
//	@Success	200				{object}	model.StatusResponse	"请求成功"
func (*setting) LegalSave(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.LegalSaveReq](c)
	service.NewLegal().Save(req)
	return ctx.SendResponse()
}
