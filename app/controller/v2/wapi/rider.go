// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-10-24, by aurb

package wapi

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/pkg/silk"
)

type rider struct{}

var Rider = new(rider)

// List
// @ID		RiderList
// @Router	/warestore/v2/rider [GET]
// @Summary	骑手列表
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string											true	"管理员校验token"
// @Param	query				query		model.RiderListReq								true	"查询参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.RiderItem}	"请求成功"
func (*rider) List(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.RiderListReq](c)
	req.EmployeeID = silk.UInt64(0)
	if ctx.Employee != nil {
		req.EmployeeID = silk.UInt64(ctx.Employee.ID)
	}
	fmt.Println(service.NewRider().GetQrcode(98784280650))
	return ctx.SendResponse(service.NewRider().List(req))
}

// Info
// @ID		RiderInfo
// @Router	/warestore/v2/rider/info [GET]
// @Summary	根据二维码获取骑手信息
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string			true	"仓管校验token"
// @Param	qrcode				query		string			true	"二维码"
// @Success	200					{object}	model.RiderItem	"请求成功"
func (*rider) Info(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.QRQueryReq](c)
	id := service.NewRider().ParseQrcode(req.Qrcode)
	return ctx.SendResponse(service.NewRider().GetRiderInfoByID(id))
}
