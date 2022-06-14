// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
    "strconv"
    "strings"
)

type business struct{}

var Business = new(business)

// Rider
// @ID           EmployeeBusinessRider
// @Router       /employee/v1/business/rider [GET]
// @Summary      E2003 骑手业务详情
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token   header  string  true  "店员校验token"
// @Param        qrcode query       string  true  "骑手二维码, 最好把`https://rider.auroraride.com/`删除"
// @Success      200    {object}    model.SubscribeBusiness  "业务详情返回"
func (*business) Rider(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.QRQueryReq](c)
    qr := strings.ReplaceAll(req.Qrcode, "https://rider.auroraride.com/", "")
    id, _ := strconv.ParseUint(qr, 10, 64)
    return ctx.SendResponse(service.NewBusinessWithEmployee(ctx.Employee).Detail(id))
}
