// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-19
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

// @title		极光出行API - 骑手端api
// @version		2.0
// @BasePath	/

var Rider = new(rider)

// Feedback
// @ID		RiderFeedback
// @Router	/rider/v2/feedback [POST]
// @Summary	AZ001 意见反馈
// @Tags	[R]骑手接口
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string				true	"骑手校验token"
// @Param	body			body		model.FeedbackReq	true	"反馈内容"
// @Success	200				{object}	bool				"请求成功"
func (*rider) Feedback(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.FeedbackReq](c)
	return ctx.SendResponse(service.NewFeedback().RiderCreate(req, ctx.Rider))
}
