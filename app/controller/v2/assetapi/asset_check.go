package assetapi

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
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.AssetCheckCreateReq	true	"盘点参数"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*assetCheck) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetCheckCreateReq](c)
	req.OpratorType = model.AssetOperateRoleTypeManager
	req.OpratorID = ctx.Modifier.ID
	return ctx.SendResponse(service.NewAssetCheck().CreateAssetCheck(ctx.Request().Context(), req, ctx.Modifier))
}

// GetAssetBySN 通过SN查询资产
// @ID		AssetCheckGetAssetBySN
// @Router	/manager/v2/asset/check/{sn} [GET]
// @Summary	通过SN查询资产
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	sn				path		string							true	"资产SN"
// @Success	200				{object}	model.AssetCheckByAssetSnRes	"请求成功"
func (*assetCheck) GetAssetBySN(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetCheckByAssetSnReq](c)
	req.OpratorType = model.AssetOperateRoleTypeManager
	req.OpratorID = ctx.Modifier.ID
	return ctx.SendResponse(service.NewAssetCheck().GetAssetBySN(ctx.Request().Context(), req))
}

// List
// @ID		AssetCheckList
// @Router	/manager/v2/asset/check [GET]
// @Summary	盘点列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string													true	"管理员校验token"
// @Param	query			query		model.AssetCheckListReq									true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetCheckListRes}	"请求成功"
func (*assetCheck) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetCheckListReq](c)
	return ctx.SendResponse(service.NewAssetCheck().List(ctx.Request().Context(), req))
}

// Abnormal
// @ID		AssetCheckAbnormal
// @Router	/manager/v2/asset/check/abnormal [GET]
// @Summary	盘点异常
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	query			query		model.AssetCheckListAbnormalReq			true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetCheckAbnormal}	"请求成功"
func (*assetCheck) Abnormal(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetCheckListAbnormalReq](c)
	return ctx.SendResponse(service.NewAssetCheck().ListAbnormal(ctx.Request().Context(), req))
}
