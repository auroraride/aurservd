// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/service"
)

type misc struct{}

var Misc = new(misc)

// Feedback
// @ID		AgentMiscFeedback
// @Router	/agent/v1/misc/feedback [POST]
// @Summary	AZ001 意见反馈
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string					true	"代理校验token"
// @Param	body			body		definition.FeedbackReq	true	"反馈内容"
// @Success	200				{object}	bool					"请求成功"
func (*misc) Feedback(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.FeedbackReq](c)
	return ctx.SendResponse(service.NewFeedback().Create(req, ctx.Agent))
}

// FeedbackImage
// @ID		AgentMiscFeedbackImage
// @Router	/agent/v1/misc/feedback/image [POST]
// @Summary	AZ002 意见反馈上传图片
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string		true	"代理校验token"
// @Param	images			formData	file		true	"图片文件"
// @Success	200				{object}	[]string	"请求成功"
func (*misc) FeedbackImage(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewFeedback().UploadImage(ctx))
}
