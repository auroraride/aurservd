package amapi

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
// @Tags	AssetMaintenance - 维修记录
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string														true	"管理员校验token"
// @Param	query					query		model.AssetMaintenanceListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetMaintenanceListRes}	"请求成功"
func (*assetMaintenance) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetMaintenanceListReq](c)
	return ctx.SendResponse(service.NewAssetMaintenance().List(ctx.Request().Context(), req))
}

// Modify
// @ID		AssetMaintenanceModify
// @Router	/manager/v2/asset/maintenance/{id} [PUT]
// @Summary	修改维修记录
// @Tags	AssetMaintenance - 维修记录
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		int								true	"维修记录ID"
// @Param	body					body		model.AssetMaintenanceModifyReq	true	"修改参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*assetMaintenance) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetMaintenanceModifyReq](c)
	return ctx.SendResponse(service.NewAssetMaintenance().Modify(ctx.Request().Context(), req, ctx.Modifier))
}

// Export
// @ID		AssetMaintenanceExport
// @Router	/manager/v2/asset/maintenance/export [POST]
// @Summary	维修记录列表导出
// @Tags	AssetMaintenance - 维修记录
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		model.AssetMaintenanceListReq	true	"查询参数"
// @Success	200						{object}	model.ExportRes					"成功"
func (*assetMaintenance) Export(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetMaintenanceListReq](c)
	return ctx.SendResponse(service.NewAssetMaintenance().Export(ctx.Request().Context(), req, ctx.Modifier))
}
