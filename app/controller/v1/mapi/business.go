// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-24
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type business struct{}

var Business = new(business)

// List
// @ID           ManagerBusinessList
// @Router       /manager/v1/business [GET]
// @Summary      M8004 业务记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.BusinessListReq  true  "列表请求筛选参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.BusinessListRes}  "请求成功"
func (*business) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BusinessListReq](c)
    return ctx.SendResponse(service.NewBusiness().ListManager(req))
}
