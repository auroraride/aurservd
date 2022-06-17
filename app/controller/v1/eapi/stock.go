// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-17
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type stock struct{}

var Stock = new(stock)

// Overview
// @ID           EmployeeStockOverview
// @Router       /employee/v1/stock/overview [GET]
// @Summary      E2009 物资概览
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Success      200  {object}  model.StockEmployeeOverview  "请求成功"
func (*stock) Overview(c echo.Context) (err error) {
    ctx := app.ContextX[app.EmployeeContext](c)
    return ctx.SendResponse(service.NewStockWithEmployee(ctx.Employee).EmployeeOverview())
}
