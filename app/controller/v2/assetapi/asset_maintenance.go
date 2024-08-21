package assetapi

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
// @Router	/manager/v2/asset/maintenance [GET]
// @Summary	维修记录列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string														true	"管理员校验token"
// @Param	query					query		model.AssetMaintenanceListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetMaintenanceListRes}	"请求成功"
func (*assetMaintenance) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetMaintenanceListReq](c)
	return ctx.SendResponse(service.NewAssetMaintenance().List(ctx.Request().Context(), req))
}

// Create
// @ID		AssetMaintenanceCreate
// @Router	/manager/v2/asset/maintenance [POST]
// @Summary	创建维修记录
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		model.AssetMaintenanceCreateReq	true	"创建参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*assetMaintenance) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetMaintenanceCreateReq](c)
	return ctx.SendResponse(service.NewAssetMaintenance().Create(ctx.Request().Context(), req, ctx.Modifier))
}
