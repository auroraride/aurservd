// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-17, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
)

type selection struct{}

var Selection = new(selection)

// WarehouseByCity
// @ID		SelectionWarehouseByCity
// @Router	/manager/v2/asset/selection/warehouse_city [GET]
// @Summary	城市仓库列表
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-AssetManager-Token	header		string							true	"管理员校验token"
// @Success	200						{object}	[]definition.WarehouseByCityRes	"请求成功"
func (*selection) WarehouseByCity(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(biz.NewWarehouse().ListByCity())
}
