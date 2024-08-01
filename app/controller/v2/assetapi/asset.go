package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assets struct{}

var Assets = new(assets)

// List
// @ID		AssetList
// @Router	/manager/v2/asset [GET]
// @Summary	资产列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string											true	"管理员校验token"
// @Param	query			query		model.AssetListReq								true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetListRes}	"请求成功"
func (*assets) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetListReq](c)
	return ctx.SendResponse(service.NewAsset().List(ctx.Request().Context(), req))
}

// Detail
// @ID		AssetDetail
// @Router	/manager/v2/asset/{id} [GET]
// @Summary	资产详情
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	id				path		uint64				true	"资产ID"
// @Success	200				{object}	model.AssetListRes	"请求成功"
func (*assets) Detail(c echo.Context) (err error) {
	return nil
}

// Create
// @ID		AssetCreate
// @Router	/manager/v2/asset [POST]
// @Summary	创建资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.AssetCreateReq	true	"创建参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*assets) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetCreateReq](c)
	return ctx.SendResponse(service.NewAsset().Create(ctx.Request().Context(), req, ctx.Modifier))
}

// Update
// @ID		AssetUpdate
// @Router	/manager/v2/asset/{id} [PUT]
// @Summary	修改资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.AssetModifyReq	true	"修改参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*assets) Update(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetModifyReq](c)
	return ctx.SendResponse(service.NewAsset().Modify(ctx.Request().Context(), req, ctx.Modifier))
}

// Delete
// @ID		AssetDelete
// @Router	/manager/v2/asset/{id} [DELETE]
// @Summary	删除资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"资产ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*assets) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewAsset().Delete(ctx.Request().Context(), req.ID))
}

// BatchCreate
// @ID		AssetBatchCreate
// @Router	/manager/v2/asset/batch [POST]
// @Summary	批量创建资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.AssetBatchCreateReq	true	"创建参数"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*assets) BatchCreate(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetBatchCreateReq](c)
	return ctx.SendResponse(service.NewAsset().BatchCreate(ctx, req, ctx.Modifier))
}

// Template
// @ID		AssetTemplate
// @Router	/manager/v2/asset/template [GET]
// @Summary	导出资产模板
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	query	query		model.AssetExportTemplateReq	true	"查询参数"
// @Success	200		{object}	model.ExportRes					"成功"
func (*assets) Template(c echo.Context) (err error) {
	ctx, req := app.ContextBinding[model.AssetExportTemplateReq](c)
	paht, name, err := service.NewAsset().DownloadTemplate(ctx.Request().Context(), req.AssetType)
	if err != nil {
		return err
	}
	return c.Attachment(paht, name+".xlsx")
}

// Export
// @ID		AssetExport
// @Router	/manager/v2/asset/export [GET]
// @Summary	导出资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	query			query		model.AssetListReq	true	"查询参数"
// @Success	200				{object}	model.ExportRes		"成功"
func (*assets) Export(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetListReq](c)
	return ctx.SendResponse(service.NewAsset().Export(ctx.Request().Context(), req, ctx.Modifier))
}
