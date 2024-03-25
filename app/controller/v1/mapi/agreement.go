package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type agreement struct{}

var Agreement = new(agreement)

// List
// @ID		ManagerAgreementList
// @Router	/manager/v1/agreement [GET]
// @Summary	协议列表
// @Tags	协议
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Success	200				{object}	[]definition.Agreement	"请求成功"
func (*agreement) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(biz.NewAgreement().List())
}

// Detail
// @ID		ManagerAgreementDetail
// @Router	/manager/v1/agreement/{id} [GET]
// @Summary	协议详情
// @Tags	协议
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"协议ID"
// @Success	200				{object}	definition.Agreement	"请求成功"
func (*agreement) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewAgreement().Detail(req.ID))
}

// Create
// @ID		ManagerAgreementCreate
// @Router	/manager/v1/agreement [POST]
// @Summary	创建协议
// @Tags	协议
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.AgreementCreateReq	true	"创建协议请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*agreement) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.AgreementCreateReq](c)
	return ctx.SendResponse(biz.NewAgreement().Create(req))
}

// Modify
// @ID		ManagerAgreementModify
// @Router	/manager/v1/agreement/{id} [PUT]
// @Summary	修改协议
// @Tags	协议
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.AgreementModifyReq	true	"修改协议请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*agreement) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.AgreementModifyReq](c)
	return ctx.SendResponse(biz.NewAgreement().Modify(req))
}

// Delete
// @ID		ManagerAgreementDelete
// @Router	/manager/v1/agreement/{id} [DELETE]
// @Summary	删除协议
// @Tags	协议
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"协议ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*agreement) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewAgreement().Delete(req.ID))
}

// Selection
// @ID		ManagerAgreementSelection
// @Router	/manager/v1/agreement/selection [GET]
// @Summary	协议下拉列表
// @Tags	协议
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// query		    query		definition.AgreementSearchReq	true	"协议查询请求参数"
// @Success	200				{object}	[]definition.AgreementSelectionRes	"请求成功"
func (*agreement) Selection(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.AgreementSearchReq](c)
	return ctx.SendResponse(biz.NewAgreement().AgreementSelection(req))
}
