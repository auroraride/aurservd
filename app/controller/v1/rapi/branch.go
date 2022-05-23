// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type branch struct{}

var Branch = new(branch)

// List
// @ID           RiderBranchList
// @Router       /rider/v1/branch [GET]
// @Summary      R20001 列举网点
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query  model.BranchWithDistanceReq  true  "desc"
// @Success      200  {object}  []model.BranchWithDistanceRes  "请求成功"
func (*branch) List(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BranchWithDistanceReq](c)
    return ctx.SendResponse(service.NewBranch().ListByDistance(req))
}
