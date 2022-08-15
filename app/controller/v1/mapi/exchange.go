// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type exchange struct{}

var Exchange = new(exchange)

// List
// @ID           ManagerExchangeList
// @Router       /manager/v1/exchange [GET]
// @Summary      MG003 换电记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.ExchangeManagerListReq  false  "筛选选项"
// @Success      200  {object}  model.PaginationRes{items=[]model.ExchangeManagerListRes}  "请求成功"
func (*exchange) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.ExchangeManagerListReq](c)
    return ctx.SendResponse(service.NewExchangeWithModifier(ctx.Modifier).List(req))
}
