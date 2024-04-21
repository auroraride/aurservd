// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-20, by liasica

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type business struct{}

var Business = new(business)

// Active
// @ID		BusinessActive
// @Router	/rider/v2/business/active [POST]
// @Summary	激活骑士卡
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Active(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
	defer func() {
		if r := recover(); r != nil {
			c.Error(r.(error))
		}
	}()
	return ctx.SendResponse(service.NewRiderBusiness(ctx.Rider).Active(req, model.RouteVersionV2))
}
