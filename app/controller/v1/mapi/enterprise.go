// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-05
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type enterprise struct{}

var Enterprise = new(enterprise)

// Create
// @ID		ManagerEnterpriseCreate
// @Router	/manager/v1/enterprise [POST]
// @Summary	创建企业
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.EnterpriseDetail	true	"desc"
// @Success	200				{object}	int						"请求成功"
func (*enterprise) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseDetail](c)
	return ctx.SendResponse(service.NewEnterpriseWithModifier(ctx.Modifier).Create(req))
}

// Modify
// @ID		ManagerEnterpriseModify
// @Router	/manager/v1/enterprise/{id} [PUT]
// @Summary	修改企业
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.EnterpriseDetailWithID	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*enterprise) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseDetailWithID](c)
	service.NewEnterpriseWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// List
// @ID		ManagerEnterpriseList
// @Router	/manager/v1/enterprise [GET]
// @Summary	列举企业
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.EnterpriseListReq								true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.EnterpriseRes}	"请求成功"
func (*enterprise) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseListReq](c)
	return ctx.SendResponse(
		service.NewEnterpriseWithModifier(ctx.Modifier).List(req),
	)
}

// Detail
// @ID		ManagerEnterpriseDetail
// @Router	/manager/v1/enterprise/{id} [GET]
// @Summary	企业详情
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	id				path		uint64				true	"企业ID"
// @Success	200				{object}	model.EnterpriseRes	"请求成功"
func (*enterprise) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewEnterpriseWithModifier(ctx.Modifier).GetDetail(req))
}

// Prepayment
// @ID		ManagerEnterprisePrepayment
// @Router	/manager/v1/enterprise/{id}/prepayment [POST]
// @Summary	企业预付费
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.EnterprisePrepaymentReq	true	"desc"
// @Success	200				{object}	float64							"当前余额"
func (*enterprise) Prepayment(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterprisePrepaymentReq](c)
	return ctx.SendResponse(
		service.NewEnterpriseWithModifier(ctx.Modifier).Prepayment(req),
	)
}

// CreateStation
// @ID		ManagerEnterpriseCreateStation
// @Router	/manager/v1/enterprise/station [POST]
// @Summary	创建站点
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	body			body		model.EnterpriseStationCreateReq	true	"desc"
// @Success	200				{object}	int64								"请求成功"
func (*enterprise) CreateStation(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseStationCreateReq](c)
	return ctx.SendResponse(service.NewEnterpriseStationWithModifier(ctx.Modifier).Create(req))
}

// ModifyStation
// @ID		ManagerEnterpriseModifyStation
// @Router	/manager/v1/enterprise/station/{id} [PUT]
// @Summary	编辑站点
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	body			body		model.EnterpriseStationModifyReq	true	"desc"
// @Param	id				path		uint64								true	"站点ID"
// @Success	200				{object}	model.StatusResponse				"请求成功"
func (*enterprise) ModifyStation(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseStationModifyReq](c)
	service.NewEnterpriseStationWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// ListStation
// @ID		ManagerEnterpriseListStation
// @Router	/manager/v1/enterprise/station [GET]
// @Summary	列举站点
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	enterpriseId	query		uint64								true	"企业ID"
// @Success	200				{object}	[]model.EnterpriseStationListRes	"请求成功"
func (*enterprise) ListStation(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseStationListReq](c)
	return ctx.SendResponse(
		service.NewEnterpriseStationWithModifier(ctx.Modifier).List(req),
	)
}

// CreateRider
// @ID		ManagerEnterpriseCreateRider
// @Router	/manager/v1/enterprise/rider [POST]
// @Summary	添加骑手
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.EnterpriseRiderCreateReq	true	"desc"
// @Success	200				{object}	model.EnterpriseRider			"请求成功"
func (*enterprise) CreateRider(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseRiderCreateReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithModifier(ctx.Modifier).Create(req))
}

// ListRider
// @ID		ManagerEnterpriseListRider
// @Router	/manager/v1/enterprise/rider [GET]
// @Summary	列举骑手
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.EnterpriseRiderListReq						true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.EnterpriseRider}	"请求成功"
func (*enterprise) ListRider(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseRiderListReq](c)
	return ctx.SendResponse(service.NewEnterpriseRiderWithModifier(ctx.Modifier).List(req))
}

// Price
// @ID		ManagerEnterprisePrice
// @Router	/manager/v1/enterprise/price [POST]
// @Summary	团签单价设定
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.EnterprisePriceReq		true	"价格详情"
// @Success	200				{object}	model.EnterprisePriceWithCity	"请求成功"
func (*enterprise) Price(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterprisePriceReq](c)
	return ctx.SendResponse(service.NewEnterpriseWithModifier(ctx.Modifier).Price(req))
}

// DeletePrice
// @ID		ManagerEnterpriseDeletePrice
// @Router	/manager/v1/enterprise/price/{id} [POST]
// @Summary	删除价格
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"价格ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*enterprise) DeletePrice(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewEnterpriseWithModifier(ctx.Modifier).DeletePrice(req)
	return ctx.SendResponse()
}

// ModifyContract
// @ID		ManagerEnterpriseModifyContract
// @Router	/manager/v1/enterprise/contract [POST]
// @Summary	编辑合同
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	body			body		model.EnterpriseContractModifyReq	true	"合同字段"
// @Success	200				{object}	model.StatusResponse				"请求成功"
func (*enterprise) ModifyContract(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseContractModifyReq](c)
	service.NewEnterpriseWithModifier(ctx.Modifier).ModifyContract(req)
	return ctx.SendResponse()
}

// DeleteContract
// @ID		ManagerEnterpriseDeleteContract
// @Router	/manager/v1/enterprise/contract/{id} [DELETE]
// @Summary	删除合同
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"合同ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*enterprise) DeleteContract(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewEnterpriseWithModifier(ctx.Modifier).DeleteContract(req)
	return ctx.SendResponse()
}

// AgentList
// @ID		ManagerEnterpriseAgentList
// @Router	/manager/v1/enterprise/agent [GET]
// @Summary	代理账号列表
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	enterpriseId	query		uint64				true	"团签ID"
// @Success	200				{object}	[]model.AgentMeta	"请求成功"
func (*enterprise) AgentList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AgentListReq](c)
	return ctx.SendResponse(service.NewAgent().List(req))
}

// AgentCreate
// @ID		ManagerEnterpriseAgentCreate
// @Router	/manager/v1/enterprise/agent [POST]
// @Summary	创建代理账号
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.AgentCreateReq	true	"账号属性"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*enterprise) AgentCreate(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AgentCreateReq](c)
	service.NewAgent(ctx.Modifier).Create(req)
	return ctx.SendResponse()
}

// AgentModify
// @ID		ManagerEnterpriseAgentModify
// @Router	/manager/v1/enterprise/agent/{id} [PUT]
// @Summary	修改代理账号
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"账号ID"
// @Param	body			body		model.AgentModifyReq	true	"账号属性"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*enterprise) AgentModify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AgentModifyReq](c)
	service.NewAgent(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// AgentDelete
// @ID		ManagerEnterpriseAgentDelete
// @Router	/manager/v1/enterprise/agent{id} [DELETE]
// @Summary	删除代理账号
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"账号ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*enterprise) AgentDelete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewAgent(ctx.Modifier).Delete(req)
	return ctx.SendResponse()
}

// BindCabinet
// @ID		ManagerEnterpriseBindCabinet
// @Router	/manager/v1/enterprise/bind/cabinet [POST]
// @Summary	团签绑定电柜
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	body	body		model.EnterpriseBindCabinetReq	true	"绑定电柜请求"
// @Success	200		{object}	model.StatusResponse			"请求成功"
func (*enterprise) BindCabinet(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EnterpriseBindCabinetReq](c)
	service.NewCabinetWithModifier(ctx.Modifier).BindCabinet(req)
	return ctx.SendResponse(model.StatusResponse{Status: true})
}

// UnbindCabinet
// @ID		ManagerEnterpriseUnbindCabinet
// @Router	/manager/v1/enterprise/unbind/cabinet/{id} [GET]
// @Summary	团签解绑电柜
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	body	body		model.IDParamReq		true	"解绑电柜请求"
// @Success	200		{object}	model.StatusResponse	"请求成功"
func (*enterprise) UnbindCabinet(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewCabinetWithModifier(ctx.Modifier).UnbindCabinet(req)
	return ctx.SendResponse()
}

// PrepaymentList
// @ID		ManagerEnterprisePrepaymentList
// @Router	/manager/v1/enterprise/prepayment [GET]
// @Summary	充值记录
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	query	query		model.PrepaymentListReq									true	"请求参数"
// @Success	200		{object}	model.PaginationRes{items=[]model.PrepaymentListRes}	"请求成功"
func (*enterprise) PrepaymentList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.PrepaymentListReq](c)
	return ctx.SendResponse(service.NewPrepayment(ctx.Modifier).List(req.EnterpriseID, req))
}

// SubscribeAlterList
// @ID		ManagerEnterpriseSubscribeAlterList
// @Router	/manager/v1/enterprise/subscribe/alter/{enterpriseId} [GET]
// @Summary	加时申请列表
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	query	query		model.SubscribeAlterListReq	false	"请求参数"
// @Success	200		{object}	model.PaginationRes{items=[]model.SubscribeAlterApplyListRes}
func (*enterprise) SubscribeAlterList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SubscribeAlterListReq](c)
	return ctx.SendResponse(service.NewSubscribeAlter(ctx.Modifier).List(req))
}

// SubscribeApply
// @ID		ManagerEnterpriseSubscribeApply
// @Router	/manager/v1/enterprise/subscribe/alter [POST]
// @Summary	审批订阅申请
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	body	body		model.SubscribeAlterReviewReq	true	"审批请求"
// @Success	200		{object}	model.StatusResponse
func (*enterprise) SubscribeApply(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SubscribeAlterReviewReq](c)
	service.NewAgentSubscribe(ctx.Modifier).AlterReview(req)
	return ctx.SendResponse()
}

// FeedbackList 反馈列表
// @ID		ManagerEnterpriseFeedbackList
// @Router	/manager/v1/enterprise/feedback [GET]
// @Summary	反馈列表
// @Tags	企业
// @Accept	json
// @Produce	json
// @Param	body	body		model.FeedbackListReq	true	"反馈列表请求"
// @Success	200		{object}	model.PaginationRes{items=[]model.FeedbackDetail}
func (*enterprise) FeedbackList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.FeedbackListReq](c)
	return ctx.SendResponse(service.NewFeedback().FeedbackList(req))
}
