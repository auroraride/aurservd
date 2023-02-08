// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type rider struct{}

var Rider = new(rider)

// List
// @ID           RiderList
// @Router       /manager/v1/rider [GET]
// @Summary      M7001 列举骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.RiderListReq  true  "请求体"
// @Success      200  {object}  model.PaginationRes{items=[]model.RiderItem} "请求成功"
func (*rider) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderListReq](c)
    return ctx.SendResponse(service.NewRiderWithModifier(ctx.Modifier).List(req))
}

// Ban
// @ID           RiderBan
// @Router       /manager/v1/rider/ban [POST]
// @Summary      M7002 封禁/解除封禁身份
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.PersonBanReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Ban(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PersonBanReq](c)
    service.NewPersonWithModifier(ctx.Modifier).Ban(req)
    return ctx.SendResponse()
}

// Block
// @ID           RiderBlock
// @Router       /manager/v1/rider/block [POST]
// @Summary      M7003 封禁/解除封禁骑手账户
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.RiderBlockReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Block(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderBlockReq](c)
    service.NewRiderWithModifier(ctx.Modifier).Block(req)
    return ctx.SendResponse()
}

// Log
// @ID           ManagerRiderLog
// @Router       /manager/v1/rider/log [GET]
// @Summary      M7005 查看骑手操作日志
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.RiderLogReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.LogOperate}  "请求成功"
func (*rider) Log(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderLogReq](c)
    return ctx.SendResponse(service.NewRiderWithModifier(ctx.Modifier).GetLogs(req))
}

// Deposit
// @ID           ManagerSubscribeDeposit
// @Router       /manager/v1/deposit [POST]
// @Summary      M7019 修改押金
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.RiderMgrDepositReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Deposit(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderMgrDepositReq](c)
    service.NewRiderMgrWithModifier(ctx.Modifier).Deposit(req)
    return ctx.SendResponse()
}

// Modify
// @ID           ManagerSubscribeModify
// @Router       /manager/v1/rider/modify [POST]
// @Summary      M7010 修改骑手资料
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.RiderMgrModifyReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderMgrModifyReq](c)
    service.NewRiderMgrWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// Delete
// @ID           ManagerRiderDelete
// @Router       /manager/v1/rider/{id} [DELETE]
// @Summary      M7011 删除骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "骑手ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Delete(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    service.NewRider().Delete(req)
    return ctx.SendResponse()
}

// FollowUpCreate
// @ID           ManagerRiderFollowUpCreate
// @Router       /manager/v1/rider/followup [POST]
// @Summary      M7012 创建骑手跟进
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.RiderFollowUpCreateReq  true  "跟进请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) FollowUpCreate(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderFollowUpCreateReq](c)
    service.NewRiderFollowupWithModifier(ctx.Modifier).Create(req)
    return ctx.SendResponse()
}

// FollowUpList
// @ID           ManagerRiderFollowUpList
// @Router       /manager/v1/rider/followup [GET]
// @Summary      M7013 获取骑手跟进
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.RiderFollowUpListReq  true  "骑手跟进筛选请求"
// @Success      200  {object}  model.PaginationRes{items=[]model.RiderFollowUpListRes}  "请求成功"
func (*rider) FollowUpList(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderFollowUpListReq](c)
    return ctx.SendResponse(service.NewRiderFollowupWithModifier(ctx.Modifier).List(req))
}

// ExchangeLimit
// @ID           ManagerRiderExchangeLimit
// @Router       /manager/v1/rider/exchange-limit [POST]
// @Summary      M7022 设置骑手换电限制
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.RiderExchangeLimitReq  true  "配置项"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) ExchangeLimit(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderExchangeLimitReq](c)
    service.NewRiderWithModifier(ctx.Modifier).ExchangeLimit(req)
    return ctx.SendResponse()
}