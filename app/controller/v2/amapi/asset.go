package amapi

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
// @Param	X-Asset-Manager-Token	header		string											true	"管理员校验token"
// @Param	query					query		model.AssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetListRes}	"请求成功"
func (*assets) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetListReq](c)
	return ctx.SendResponse(service.NewAsset().List(ctx.Request().Context(), req))
}

// Detail
// @ID		AssetDetail
// @Router	/manager/v2/asset/{id} [GET]
// @Summary	资产详情
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string				true	"管理员校验token"
// @Param	id						path		uint64				true	"资产ID"
// @Success	200						{object}	model.AssetListRes	"请求成功"
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
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	body					body		model.AssetCreateReq	true	"创建参数"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assets) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetCreateReq](c)
	return ctx.SendResponse(service.NewAsset(ctx.Operator).Create(ctx.Request().Context(), req, ctx.Modifier))
}

// Update
// @ID		AssetUpdate
// @Router	/manager/v2/asset/{id} [PUT]
// @Summary	修改资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		uint64					true	"资产ID"
// @Param	body					body		model.AssetModifyReq	true	"修改参数"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assets) Update(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetModifyReq](c)
	return ctx.SendResponse(service.NewAsset().Modify(ctx.Request().Context(), req, ctx.Modifier))
}

// BatchCreate
// @ID		AssetBatchCreate
// @Router	/manager/v2/asset/batch [POST]
// @Summary	批量创建资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	assetType				formData	uint8					true	"资产类型 1:电车 2:智能电池 3:非智能电池 4:电柜配件 5:电车配件 6:其它"
// @Param	file					formData	file					true	"文件"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assets) BatchCreate(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetBatchCreateReq](c)
	return ctx.SendResponse(service.NewAsset(ctx.Operator).BatchCreate(ctx, req, ctx.Modifier))
}

// Template
// @ID		AssetTemplate
// @Router	/manager/v2/asset/template [GET]
// @Summary	导出资产模板
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	query					query		model.AssetExportTemplateReq	true	"查询参数"
// @Success	200						{object}	model.ExportRes					"成功"
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
// @Router	/manager/v2/asset/export [POST]
// @Summary	导出资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string				true	"管理员校验token"
// @Param	body					body		model.AssetListReq	true	"查询参数"
// @Success	200						{object}	model.ExportRes		"成功"
func (*assets) Export(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetListReq](c)
	return ctx.SendResponse(service.NewAsset().Export(ctx.Request().Context(), req, ctx.Modifier))
}

// Count
// @ID		AssetCount
// @Router	/manager/v2/asset/count [GET]
// @Summary	查询有效的资产数量
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string				true	"管理员校验token"
// @Param	query					query		model.AssetFilter	true	"查询参数"
// @Success	200						{object}	model.AssetNumRes	"请求成功"
func (*assets) Count(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetFilter](c)
	return ctx.SendResponse(service.NewAsset().Count(ctx.Request().Context(), req))
}
