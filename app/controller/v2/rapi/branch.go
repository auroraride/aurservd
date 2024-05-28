// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-28, by Jorjan

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type branch struct{}

var Branch = new(branch)

// List
// @ID		BranchList
// @Router	/rider/v2/branch [GET]
// @Summary	列举网点
// @Tags	Branch - 网点
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	query			query		definition.BranchWithDistanceReq		true	"根据距离获取网点请求参数"
// @Success	200				{object}	[]definition.BranchWithDistanceRes	"请求成功"
func (*branch) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.BranchWithDistanceReq](c)
	return ctx.SendResponse(biz.NewBranchWithRider(ctx.Rider).ListByDistanceRider(req))
}
