package amapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetExport struct{}

var AssetExport = new(assetExport)

// List
// @ID		AssetManagerExportList
// @Router	/manager/v2/asset/export [GET]
// @Summary	导出列表
// @Tags	AssetExport - 导出
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string													true	"管理员校验token"
// @Param	query					query		model.AssetExportListReq								false	"分页信息"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetExportListRes}	"请求成功"
func (*assetExport) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetExportListReq](c)
	return ctx.SendResponse(service.NewAssetExportWithModifier(ctx.Modifier).List(ctx.AssetManager, req))
}

// Download
// @ID		AssetManagerExportDownload
// @Router	/manager/v2/asset/export/download/{sn} [GET]
// @Summary	下载文件
// @Tags	AssetExport - 导出
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	sn						path		string					true	"编号"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assetExport) Download(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetExportDownloadReq](c)
	return ctx.Attachment(service.NewAssetExportWithModifier(ctx.Modifier).Download(req))
}
