package amapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetCheck struct{}

var AssetCheck = new(assetCheck)

// Create 盘点资产
// @ID		AssetCheckCreate
// @Router	/manager/v2/asset/check [POST]
// @Summary	盘点资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	body					body		model.AssetCheckCreateReq	true	"盘点参数"
// @Success	200						{object}	model.StatusResponse		"请求成功"
func (*assetCheck) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetCheckCreateReq](c)
	req.OperatorType = model.OperatorTypeAssetManager
	req.OperatorID = ctx.Modifier.ID
	return ctx.SendResponse(service.NewAssetCheck().CreateAssetCheck(ctx.Request().Context(), req, ctx.Modifier))
}

// GetAssetBySN 通过SN查询资产
// @ID		AssetCheckGetAssetBySN
// @Router	/manager/v2/asset/check/sn/{sn} [GET]
// @Summary	通过SN查询资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	sn						path		string							true	"资产SN"
// @Success	200						{object}	model.AssetCheckByAssetSnRes	"请求成功"
func (*assetCheck) GetAssetBySN(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetCheckByAssetSnReq](c)
	req.OperatorType = model.OperatorTypeAssetManager
	req.OperatorID = ctx.Modifier.ID
	return ctx.SendResponse(service.NewAssetCheck().GetAssetBySN(ctx.Request().Context(), req))
}

// List
// @ID		AssetCheckList
// @Router	/manager/v2/asset/check [GET]
// @Summary	盘点列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string													true	"管理员校验token"
// @Param	query					query		model.AssetCheckListReq									true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetCheckListRes}	"请求成功"
func (*assetCheck) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetCheckListReq](c)
	return ctx.SendResponse(service.NewAssetCheck().List(ctx.Request().Context(), req))
}

// Abnormal
// @ID		AssetCheckAbnormal
// @Router	/manager/v2/asset/check/abnormal/{id} [GET]
// @Summary	盘点异常
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		uint64							true	"盘点ID"
// @Param	query					query		model.AssetCheckListAbnormalReq	true	"查询参数"
// @Success	200						{object}	[]model.AssetCheckAbnormal		"请求成功"
func (*assetCheck) Abnormal(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetCheckListAbnormalReq](c)
	return ctx.SendResponse(service.NewAssetCheck().ListAbnormal(ctx.Request().Context(), req))
}

// AssetDetailList
// @ID		AssetCheckDetailList
// @Router	/manager/v2/asset/check/asset/{id} [GET]
// @Summary	盘点资产明细
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string												true	"管理员校验token"
// @Param	id						path		uint64												true	"盘点ID"
// @Param	query					query		model.AssetCheckDetailListReq						true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetCheckDetail}	"请求成功"
func (*assetCheck) AssetDetailList(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetCheckDetailListReq](c)
	return ctx.SendResponse(service.NewAssetCheck().AssetDetailList(ctx.Request().Context(), req))
}

// AbnormalOperate
// @ID		AssetCheckAbnormalOperate
// @Router	/manager/v2/asset/check/abnormal/operate/{id} [PUT]
// @Summary	盘点操作
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string								true	"管理员校验token"
// @Param	id						path		uint64								true	"盘点详情ID"
// @Param	body					body		model.AssetCheckAbnormalOperateReq	true	"操作参数"
// @Success	200						{object}	model.StatusResponse				"请求成功"
func (*assetCheck) AbnormalOperate(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetCheckAbnormalOperateReq](c)
	return ctx.SendResponse(service.NewAssetCheck().AssetCheckAbnormalOperate(ctx.Request().Context(), req, ctx.Modifier))
}

// Detail
// @ID		AssetCheckDetail
// @Router	/manager/v2/asset/check/{id} [GET]
// @Summary	盘点详情
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		uint64					true	"盘点ID"
// @Success	200						{object}	model.AssetCheckListRes	"请求成功"
func (*assetCheck) Detail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewAssetCheck().Detail(ctx.Request().Context(), req.ID))
}
