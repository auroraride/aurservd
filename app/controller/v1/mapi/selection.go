// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type selection struct {
}

var Selection = new(selection)

// Plan
// @ID		ManagerSelectionPlan
// @Router	/manager/v1/selection/plan [GET]
// @Summary	筛选骑士卡
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	query			query		model.PlanSelectionReq			false	"骑士卡筛选项"
// @Success	200				{object}	[]model.CascaderOptionLevel3	"请求成功"
func (*selection) Plan(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.PlanSelectionReq](c)
	return ctx.SendResponse(service.NewSelection().Plan(req))
}

// Rider
// @ID		ManagerSelectionRider
// @Router	/manager/v1/selection/rider [GET]
// @Summary	筛选骑手
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	query			query		model.RiderSelectionReq	true	"骑手筛选项"
// @Success	200				{object}	[]model.SelectOption	"请求成功"
func (*selection) Rider(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.RiderSelectionReq](c)
	return ctx.SendResponse(service.NewSelection().Rider(req))
}

// Store
// @ID		ManagerSelectionStore
// @Router	/manager/v1/selection/store [GET]
// @Summary	筛选门店
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Store(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewSelection().Store())
}

// Employee
// @ID		ManagerSelectionEmployee
// @Router	/manager/v1/selection/employee [GET]
// @Summary	筛选店员
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Employee(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewSelection().Employee())
}

// City
// @ID		ManagerSelectionCity
// @Router	/manager/v1/selection/city [GET]
// @Summary	筛选启用的城市
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) City(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewSelection().City())
}

// Branch
// @ID		ManagerSelectionBranch
// @Router	/manager/v1/selection/branch [GET]
// @Summary	筛选网点
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Branch(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewSelection().Branch())
}

// Enterprise
// @ID		ManagerSelectionEnterprise
// @Router	/manager/v1/selection/enterprise [GET]
// @Summary	筛选企业
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Enterprise(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewSelection().Enterprise())
}

// Cabinet
// @ID		ManagerSelectionCabinet
// @Router	/manager/v1/selection/cabinet [GET]
// @Summary	筛选电柜
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	query			query		model.CabinetSelectionReq		false	"筛选参数"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Cabinet(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CabinetSelectionReq](c)
	return ctx.SendResponse(service.NewSelection().Cabinet(req))
}

// Role
// @ID		ManagerSelectionRole
// @Router	/manager/v1/selection/role [GET]
// @Summary	筛选角色
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Success	200				{object}	[]model.SelectOption	"请求成功"
func (*selection) Role(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewSelection().Role())
}

// WxEmployee
// @ID		ManagerSelectionWxEmployee
// @Router	/manager/v1/selection/wxemployee [GET]
// @Summary	筛选企业微信成员
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Success	200				{object}	[]map[string]interface{}	"请求成功"
func (*selection) WxEmployee(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewSelection().WorkwxEmployee())
}

// PlanModel
// @ID		ManagerSelectionPlanModel
// @Router	/manager/v1/selection/planmodel [GET]
// @Summary	筛选骑行卡电池
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	query			query		model.SelectionPlanModelReq	true	"选项"
// @Success	200				{object}	[]string					"电池型号列表"
func (*selection) PlanModel(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SelectionPlanModelReq](c)
	return ctx.SendResponse(service.NewSelection().PlanModel(req))
}

// CabinetModel
// @ID		ManagerSelectionCabinetModel
// @Router	/manager/v1/selection/cabinetmodel [GET]
// @Summary	筛选电柜电池
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string									true	"管理员校验token"
// @Param	query			query		model.SelectionCabinetModelByCabinetReq	true	"选项"
// @Success	200				{object}	[]string								"电池型号列表"
func (*selection) CabinetModel(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SelectionCabinetModelByCabinetReq](c)
	return ctx.SendResponse(service.NewSelection().CabinetModel(req))
}

// Model
// @ID		ManagerSelectionModel
// @Router	/manager/v1/selection/model [GET]
// @Summary	筛选电池型号
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string		true	"管理员校验token"
// @Success	200				{object}	[]string	"电池型号列表"
func (*selection) Model(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewSelection().Models())
}

// CouponTemplate
// @ID			ManagerSelectionCouponTemplate
// @Router		/manager/v1/selection/coupon/template [GET]
// @Summary		MB014 筛选优惠券模板
// @Description	筛选样式参考 <a target="_blank" href="https://element.eleme.cn/#/zh-CN/component/select#fen-zu">ElementUI-select-分组</a> <a target="_blank" href="https://element.eleme.cn/#/zh-CN/component/select#zi-ding-yi-mo-ban">ElementUI-select-自定义模板</a>
// @Tags		筛选
// @Accept		json
// @Produce		json
// @Param		X-Manager-Token	header		string						true	"管理员校验token"
// @Success		200				{object}	[]model.SelectOptionGroup	"请求成功"
func (*selection) CouponTemplate(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewSelection().CouponTemplate())
}

// EbikeBrand
// @ID		ManagerSelectionEbikeBrand
// @Router	/manager/v1/selection/ebike/brand [GET]
// @Summary	车辆型号列表
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Success	200				{object}	[]model.SelectOption	"请求成功"
func (*selection) EbikeBrand(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewSelection().EbikeBrand())
}

// BatterySerial
// @ID		ManagerSelectionBatterySerial
// @Router	/manager/v1/selection/battery/serial [GET]
// @Summary	按流水号搜索电池
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	query			query		model.BatterySearchReq	true	"请求参数"
// @Success	200				{object}	[]model.Battery			"请求成功"
func (*selection) BatterySerial(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BatterySearchReq](c)
	return ctx.SendResponse(service.NewSelection().BatterySerialSearch(req))
}

// CommissionPlan
// @ID		ManagerSelectionCommissionPlan
// @Router	/manager/v1/selection/commission/plan [GET]
// @Summary	返佣方案筛选骑士卡
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	query	query		model.CommissionPlanSelectionReq	true	"请求参数"
// @Success	200		{object}	[]model.SelectOption				"请求成功"
func (*selection) CommissionPlan(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.CommissionPlanSelectionReq](c)
	return ctx.SendResponse(service.NewSelection().CommissionPlan(req))
}

// QuestionCategory 问题分类
// @ID		ManagerSelectionQuestionCategory
// @Router	/manager/v1/selection/question/category [GET]
// @Summary	问题分类
// @Tags	筛选
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Success	200				{object}	[]model.SelectOption	"请求成功"
func (*selection) QuestionCategory(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewSelection().QuestionCategory())
}
