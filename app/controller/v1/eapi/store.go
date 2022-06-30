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

type store struct{}

var Store = new(store)

// Status
// @ID           EmployeeStoreStatus
// @Router       /employee/v1/store/status [POST]
// @Summary      E1006 切换门店状态
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        body  body  model.StoreSwtichStatusReq  true  "状态请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*store) Status(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.StoreSwtichStatusReq](c)
    service.NewStoreWithEmployee(ctx.Employee).SwitchStatus(req)
    return ctx.SendResponse()
}
