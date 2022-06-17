// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-16
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type exchange struct{}

var Exchange = new(exchange)

// List
// @ID           EmployeeExchangeList
// @Router       /employee/v1/exchange [GET]
// @Summary      E2008 换电记录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        query  query   model.ExchangeListReq  true  "列表请求筛选参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.ExchangeEmployeeListRes}  "请求成功"
func (*exchange) List(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.ExchangeListReq](c)
    return ctx.SendResponse(service.NewExchangeWithEmployee(ctx.Employee).EmployeeList(req))
}
