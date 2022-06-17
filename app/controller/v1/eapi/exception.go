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

type exception struct{}

var Exception = new(exception)

// Setting
// @ID           EmployeeExceptionSetting
// @Router       /employee/v1/exception/setting [GET]
// @Summary      E3001 物资异常配置
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token   header  string  true  "店员校验token"
// @Success      200  {object}      []model.InventoryItem  "请求成功"
func (*exception) Setting(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(
        service.NewException().Setting(),
    )
}
