package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetScrap struct{}

var AssetScrap = new(assetScrap)

// Scrap
// @ID		AssetScrap
// @Router	/manager/v2/asset/scrap [POST]
// @Summary	报废资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	body					body		model.AssetScrapReq		true	"报废参数"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assetScrap) Scrap(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetScrapReq](c)
	return ctx.SendResponse(service.NewAssetScrap().Scrap(ctx.Request().Context(), req, ctx.Modifier))
}

// ScrapList
// @ID		AssetScrapList
// @Router	/manager/v2/asset/scrap [GET]
// @Summary	资产报废列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string													true	"管理员校验token"
// @Param	query					query		model.AssetScrapListReq									true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetScrapListRes}	"请求成功"
func (*assetScrap) ScrapList(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetScrapListReq](c)
	return ctx.SendResponse(service.NewAssetScrap().ScrapList(ctx.Request().Context(), req))
}

// ScrapBatchRestore
// @ID		AssetScrapBatchRestore
// @Router	/manager/v2/asset/scrap/batch/restore [POST]
// @Summary	资产报废还原
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		model.AssetScrapBatchRestoreReq	true	"还原参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*assetScrap) ScrapBatchRestore(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetScrapBatchRestoreReq](c)
	return ctx.SendResponse(service.NewAssetScrap().ScrapBatchRestore(ctx.Request().Context(), req, ctx.Modifier))
}

// ScrapReasonSelect
// @ID		AssetScrapReasonSelect
// @Router	/manager/v2/asset/scrap/reason [GET]
// @Summary	报废理由Select
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Success	200						{object}	[]model.SelectOption	"请求成功"
func (*assetScrap) ScrapReasonSelect(c echo.Context) (err error) {
	ctx := app.GetAssetManagerContext(c)
	return ctx.SendResponse(service.NewAssetScrap().ScrapReasonSelect(ctx.Request().Context()))
}
