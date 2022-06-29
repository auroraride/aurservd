// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-17
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type rider struct{}

var Rider = new(rider)

// Detail
// @ID           EmployeeRiderDetail
// @Router       /employee/v1/rider [GET]
// @Summary      E4001 获取骑手信息
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        phone  query   string  true  "骑手手机号"
// @Success      200  {object}  model.RiderEmployeeSearchRes  "请求成功"
func (*rider) Detail(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.RiderPhoneSearchReq](c)
    return ctx.SendResponse(
        service.NewRiderMgrWithEmployee(ctx.Employee).QueryPhone(req.Phone),
    )
}

// Exchange
// @ID           EmployeeRiderExchange
// @Router       /employee/v1/rider/exchange [GET]
// @Summary      E4002 骑手换电记录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        query  query   model.RiderExchangeReq  true  "换电记录筛选项"
// @Success      200  {object}  model.PaginationRes{items=[]model.ExchangeRiderListRes}  "请求成功"
func (*rider) Exchange(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.RiderExchangeReq](c)
    return ctx.SendResponse(
        service.NewExchange().RiderList(req.RiderID, req.PaginationReq),
    )
}
