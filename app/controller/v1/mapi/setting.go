// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type setting struct{}

var Setting = new(setting)

// List
// @ID		ManagerSettingList
// @Router	/manager/v1/setting [GET]
// @Summary	列举设置
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Success	200				{object}	[]model.SettingReq	"请求成功"
func (*setting) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewSetting().List())
}

// Modify
// @ID		ManagerSettingModify
// @Router	/manager/v1/setting/{key} [PUT]
// @Summary	调整设置
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	key				path		string					true	"设置项"
// @Param	body			body		model.SettingReq		true	"desc"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*setting) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SettingReq](c)
	service.NewSettingWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// LegalRead
// @ID		ManagerSettingLegalRead
// @Router	/manager/v1/setting/legal/{name} [GET]
// @Summary	获取法规
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string			true	"管理员校验token"
// @Param	name			path		string			true	"法规名称, policy: APP隐私政策; agreement: APP服务协议; agent-policy: 代理端小程序隐私政策; agent-agreement: 代理端小程序服务协议; promote-policy: 推广端小程序隐私政策; promote-agreement: 推广端小程序服务协议"
// @Success	200				{object}	model.LegalRes	"请求成功"
func (*setting) LegalRead(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.LegalName](c)
	return ctx.SendResponse(service.NewLegal().Read(req))
}

// LegalSave
// @ID		ManagerSettingLegalSave
// @Router	/manager/v1/setting/legal [POST]
// @Summary	保存法规
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.LegalSaveReq		true	"请求参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*setting) LegalSave(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.LegalSaveReq](c)
	service.NewLegal().Save(req)
	return ctx.SendResponse()
}

// ActivityList
// @ID		SettingActivityList
// @Router	/manager/v1/setting/activity [GET]
// @Summary	获取活动列表
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	query			query		definition.ActivityListReq	true	"请求参数"
// @Success	200				{object}	[]definition.ActivityDetail	"请求成功"
func (*setting) ActivityList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.ActivityListReq](c)
	return ctx.SendResponse(biz.NewActivity().List(req))
}

// ActivityDetail
// @ID		SettingActivityDetail
// @Router	/manager/v1/setting/activity/{id} [GET]
// @Summary	获取指定活动
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	id				path		string						true	"活动ID"
// @Success	200				{object}	definition.ActivityDetail	"请求成功"
func (*setting) ActivityDetail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewActivity().Detail(req.ID))
}

// ActivityModify
// @ID		SettingActivityModify
// @Router	/manager/v1/setting/activity [PUT]
// @Summary	修改活动
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.ActivityModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*setting) ActivityModify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.ActivityModifyReq](c)
	return ctx.SendResponse(biz.NewActivityWithModifierBiz(ctx.Modifier).Modify(req))
}

// ActivityCreate
// @ID		SettingActivityCreate
// @Router	/manager/v1/setting/activity [POST]
// @Summary	创建活动
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.ActivityCreateReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*setting) ActivityCreate(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.ActivityCreateReq](c)
	return ctx.SendResponse(biz.NewActivityWithModifierBiz(ctx.Modifier).Create(req))
}

// VersionCreate
// @ID		SettingVersionCreate
// @Router	/manager/v1/setting/version [POST]
// @Summary	创建版本
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		definition.VersionReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*setting) VersionCreate(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.VersionReq](c)
	return ctx.SendResponse(biz.NewVersionWithModifierBiz(ctx.Modifier).Create(req))
}

// VersionModify
// @ID		SettingVersionModify
// @Router	/manager/v1/setting/version/{id} [PUT]
// @Summary	修改版本
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		definition.VersionModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*setting) VersionModify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.VersionModifyReq](c)
	return ctx.SendResponse(biz.NewVersionWithModifierBiz(ctx.Modifier).Modify(req))
}

// VersionList
// @ID		SettingVersionList
// @Router	/manager/v1/setting/version [GET]
// @Summary	获取版本列表
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	query			query		definition.VersionListReq	true	"请求参数"
// @Success	200				{object}	[]definition.Version		"请求成功"
func (*setting) VersionList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.VersionListReq](c)
	return ctx.SendResponse(biz.NewVersion().List(req))
}

// VersionDelete
// @ID		SettingVersionDelete
// @Router	/manager/v1/setting/version/{id} [DELETE]
// @Summary	删除版本
// @Tags	设置
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"版本ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*setting) VersionDelete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewVersionWithModifierBiz(ctx.Modifier).Delete(req.ID))
}
