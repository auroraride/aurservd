// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-11, by lisicen

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type pushmessage struct{}

var Pushmessage = new(pushmessage)

// Create
// @ID		PushmessageCreate
// @Router	/manager/v1/push/message [POST]
// @Summary	保存推送消息
// @Tags	Pushmessage - 推送消息
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.PushmessageSaveReq	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*pushmessage) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.PushmessageSaveReq](c)
	return ctx.SendResponse(biz.NewPushmessage().Create(req))
}

// Modify
// @ID		PushmessageModify
// @Router	/manager/v1/push/message [PUT]
// @Summary	更新推送消息
// @Tags	Pushmessage - tags
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.PushmessageModifyReq	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*pushmessage) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.PushmessageModifyReq](c)
	return ctx.SendResponse(biz.NewPushmessage().Modify(req))
}

// Delete
// @ID		PushmessageDelete
// @Router	/manager/v1/push/message/{id} [DELETE]
// @Summary	删除推送消息
// @Tags	Pushmessage - tags
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"消息ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*pushmessage) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.PushmessageDeleteReq](c)
	return ctx.SendResponse(biz.NewPushmessage().Delete(req))
}

// Get
// @ID		PushmessageGet
// @Router	/manager/v1/push/message/{id} [GET]
// @Summary	获取推送消息
// @Tags	Pushmessage - tags
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	id				path		string							true	"消息ID"
// @Success	200				{object}	definition.PushmessageDetail	"请求成功"
func (*pushmessage) Get(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.PushmessageGetReq](c)
	return ctx.SendResponse(biz.NewPushmessage().Get(req))
}
