// Copyright (C) liasica. 2023-present.
//
// Created at 2023-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package oapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// List
// @ID           MaintainerCabinetList
// @Router       /maintainer/v1/cabinet [GET]
// @Summary      O2001 获取电柜列表
// @Tags         [O]运维接口
// @Accept       json
// @Produce      json
// @Param        X-Maintainer-Token  header  string  true  "运维校验token"
// @Success      200  {object}  model.Pagination{items=[]model.CabinetListByDistanceRes}  "请求成功"
func (*cabinet) List(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[model.PaginationReq](c)
	return ctx.SendResponse(service.NewMaintainerCabinet().List(ctx.CityIDs(), req))
}
