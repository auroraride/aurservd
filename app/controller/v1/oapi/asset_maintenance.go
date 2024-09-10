// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-28, by aurb

package oapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetMaintenance struct{}

var AssetMaintenance = new(assetMaintenance)

// List
// @ID		AssetMaintenanceList
// @Router	/maintainer/v1/asset/maintenance [GET]
// @Summary	维保记录
// @Tags	AssetMaintenance - 维保
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string														true	"管理员校验token"
// @Param	query				query		model.AssetMaintenanceListReq								true	"查询参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetMaintenanceListRes}	"请求成功"
func (*assetMaintenance) List(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[model.AssetMaintenanceListReq](c)
	return ctx.SendResponse(service.NewAssetMaintenance().List(ctx.Request().Context(), req))
}
