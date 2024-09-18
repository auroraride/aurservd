package amapi

type export struct{}

var Export = new(export)

// // List
// // @ID		AssetExportList
// // @Router	/manager/v1/export [GET]
// // @Summary	导出列表
// // @Tags	导出
// // @Accept	json
// // @Produce	json
// // @Param	X-Manager-Token	header		string												true	"管理员校验token"
// // @Param	query			query		model.ExportListReq									false	"分页信息"
// // @Success	200				{object}	model.PaginationRes{items=[]model.ExportListRes}	"请求成功"
// func (*export) List(c echo.Context) (err error) {
// 	ctx, req := app.AssetManagerContextAndBinding[model.ExportListReq](c)
// 	return ctx.SendResponse(service.NewExportWithModifier(ctx.Modifier).List(ctx.AssetManager, req))
// }
//
// // Download
// // @ID		ManagerExportDownload
// // @Router	/manager/v1/export/download/{sn} [GET]
// // @Summary	下载文件
// // @Tags	导出
// // @Accept	json
// // @Produce	json
// // @Param	X-Manager-Token	header		string					true	"管理员校验token"
// // @Param	sn				path		string					true	"编号"
// // @Success	200				{object}	model.StatusResponse	"请求成功"
// func (*export) Download(c echo.Context) (err error) {
// 	ctx, req := app.ManagerContextAndBinding[model.ExportDownloadReq](c)
// 	return ctx.Attachment(service.NewExportWithModifier(ctx.Modifier).Download(req))
// }
