package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionCommission struct {
}

var PromotionCommission = new(promotionCommission)

// List
// @ID		ManagerPromotionCommissionList
// @Router	/manager/v1/promotion/commission [GET]
// @Summary	推广返佣方案列表
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Success	200				{object}	[]promotion.CommissionDetail	"请求成功"
func (p *promotionCommission) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionCommissionService().List())
}

// Detail
// @ID		ManagerPromotionCommissionDetail
// @Router	/manager/v1/promotion/commission/{id} [GET]
// @Summary	推广返佣方案详情
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Param	id				path		int		true	"推广返佣方案ID"
// @Success	200				{object}	promotion.CommissionDetail
func (p *promotionCommission) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewPromotionCommissionService().Detail(req.ID))
}

// Create
// @ID		ManagerPromotionCommissionCreate
// @Router	/manager/v1/promotion/commission [POST]
// @Summary	创建推广返佣方案
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		promotion.CommissionCreateReq	true	"创建请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (p *promotionCommission) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.CommissionCreateReq](c)
	service.NewPromotionCommissionService(ctx.Modifier).Create(req)
	return ctx.SendResponse()
}

// Update
// @ID		ManagerPromotionCommissionUpdate
// @Router	/manager/v1/promotion/commission [PUT]
// @Summary	更新推广返佣方案
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		promotion.CommissionCreateReq	true	"修改请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (p *promotionCommission) Update(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.CommissionCreateReq](c)
	service.NewPromotionCommissionService(ctx.Modifier).Update(req)
	return ctx.SendResponse()
}

// Delete
// @ID		ManagerPromotionCommissionDelete
// @Router	/manager/v1/promotion/commission/{id} [DELETE]
// @Summary	删除推广返佣方案
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		int						true	"推广返佣方案ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (p *promotionCommission) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewPromotionCommissionService(ctx.Modifier).Delete(req.ID)
	return ctx.SendResponse()
}

// Enable
// @ID		ManagerPromotionCommissionUpdateEnable
// @Router	/manager/v1/promotion/commission/enable [POST]
// @Summary	更新推广返佣方案状态
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		promotion.CommissionEnableReq	true	"修改请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (p *promotionCommission) Enable(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.CommissionEnableReq](c)
	service.NewPromotionCommissionService(ctx.Modifier).StatusUpdate(req)
	return ctx.SendResponse()
}

// HistoryList
// @ID		ManagerPromotionCommissionHistoryList
// @Router	/manager/v1/promotion/commission/history/{id} [GET]
// @Summary	推广返佣方案历史列表
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	id				path		int								true	"推广返佣方案ID"
// @Success	200				{object}	[]promotion.CommissionDetail	"请求成功"
func (p *promotionCommission) HistoryList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewPromotionCommissionService().HistoryList(req.ID))
}

// Selection
// @ID		ManagerPromotionCommissionSelection
// @Router	/manager/v1/promotion/commission/selection [GET]
// @Summary	推广返佣方案选择
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Success	200				{object}	[]promotion.CommissionSelection
func (p *promotionCommission) Selection(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionCommissionService().Selection())
}

// TaskSelection
// @ID		ManagerPromotionCommissionTaskSelection
// @Router	/manager/v1/promotion/commission/task/selection [GET]
// @Summary	返佣方案任务选择
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Success	200				{object}	[]promotion.CommissionTaskSelect
func (p *promotionCommission) TaskSelection(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionCommissionService().CommissionTaskSelection())
}

// CommissionPlanList
// @ID		ManagerPromotionCommissionPlanList
// @Router	/manager/v1/promotion/commission/plan/list/{:id} [GET]
// @Summary	返佣方案骑士卡列表
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Param	id				path		int		true	"会员id"
// @Success	200				{object}	[]promotion.CommissionPlanListRes
func (p *promotionCommission) CommissionPlanList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewPromotionCommissionService().CommissionPlanList(req))
}
