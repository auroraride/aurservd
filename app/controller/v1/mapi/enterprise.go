// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-05
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type enterprise struct{}

var Enterprise = new(enterprise)

// Create
// @ID           ManagerEnterpriseCreate
// @Router       /manager/v1/enterprise [POST]
// @Summary      M9001 创建企业
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EnterpriseDetail  true  "desc"
// @Success      200  {object}  int  "请求成功"
func (*enterprise) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseDetail](c)
    return ctx.SendResponse(service.NewEnterpriseWithModifier(ctx.Modifier).Create(req))
}

// Modify
// @ID           ManagerEnterpriseModify
// @Router       /manager/v1/enterprise/{id} [PUT]
// @Summary      M9002 修改企业
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.EnterpriseDetailWithID  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*enterprise) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseDetailWithID](c)
    service.NewEnterpriseWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// List
// @ID           ManagerEnterpriseList
// @Router       /manager/v1/enterprise [GET]
// @Summary      M9003 列举企业
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.EnterpriseListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.EnterpriseRes} "请求成功"
func (*enterprise) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseListReq](c)
    return ctx.SendResponse(
        service.NewEnterpriseWithModifier(ctx.Modifier).List(req),
    )
}

// Detail
// @ID           ManagerEnterpriseDetail
// @Router       /manager/v1/enterprise/{id} [GET]
// @Summary      M9004 企业详情
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "企业ID"
// @Success      200  {object}  model.EnterpriseRes  "请求成功"
func (*enterprise) Detail(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    return ctx.SendResponse(service.NewEnterpriseWithModifier(ctx.Modifier).GetDetail(req))
}

// Prepayment
// @ID           ManagerEnterprisePrepayment
// @Router       /manager/v1/enterprise/{id}/prepayment [POST]
// @Summary      M9005 企业预付费
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EnterprisePrepaymentReq  true  "desc"
// @Success      200  {object}  float64  "当前余额"
func (*enterprise) Prepayment(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterprisePrepaymentReq](c)
    return ctx.SendResponse(
        service.NewEnterpriseWithModifier(ctx.Modifier).Prepayment(req),
    )
}

// CreateStation
// @ID           ManagerEnterpriseCreateStation
// @Router       /manager/v1/enterprise/station [POST]
// @Summary      M9006 创建站点
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EnterpriseStationCreateReq  true  "desc"
// @Success      200  {object}  int64  "请求成功"
func (*enterprise) CreateStation(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseStationCreateReq](c)
    return ctx.SendResponse(service.NewEnterpriseStationWithModifier(ctx.Modifier).Create(req))
}

// ModifyStation
// @ID           ManagerEnterpriseModifyStation
// @Router       /manager/v1/enterprise/station/{id} [PUT]
// @Summary      M9007 编辑站点
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EnterpriseStationModifyReq  true  "desc"
// @Param        id  path  uint64  true  "站点ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*enterprise) ModifyStation(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseStationModifyReq](c)
    service.NewEnterpriseStationWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// ListStation
// @ID           ManagerEnterpriseListStation
// @Router       /manager/v1/enterprise/station [GET]
// @Summary      M9008 列举站点
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        enterpriseId  query  uint64  true  "企业ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*enterprise) ListStation(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseStationListReq](c)
    return ctx.SendResponse(
        service.NewEnterpriseStationWithModifier(ctx.Modifier).List(req),
    )
}

// CreateRider
// @ID           ManagerEnterpriseCreateRider
// @Router       /manager/v1/enterprise/rider [POST]
// @Summary      M9009 添加骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.EnterpriseRiderCreateReq  true  "desc"
// @Success      200  {object}  model.EnterpriseRider  "请求成功"
func (*enterprise) CreateRider(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseRiderCreateReq](c)
    return ctx.SendResponse(service.NewEnterpriseRiderWithModifier(ctx.Modifier).Create(req))
}

// ListRider
// @ID           ManagerEnterpriseListRider
// @Router       /manager/v1/enterprise/rider [GET]
// @Summary      M9010 列举骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.EnterpriseRiderListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.EnterpriseRider}  "请求成功"
func (*enterprise) ListRider(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseRiderListReq](c)
    return ctx.SendResponse(service.NewEnterpriseRiderWithModifier(ctx.Modifier).List(req))
}

// ModifyPrice
// @ID           ManagerEnterpriseModifyPrice
// @Router       /manager/v1/enterprise/price [POST]
// @Summary      M9016 编辑价格
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.EnterprisePriceModifyReq  true  "价格详情"
// @Success      200  {object}  model.EnterprisePriceWithCity  "请求成功"
func (*enterprise) ModifyPrice(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterprisePriceModifyReq](c)
    return ctx.SendResponse(service.NewEnterpriseWithModifier(ctx.Modifier).ModifyPrice(req))
}

// DeletePrice
// @ID           ManagerEnterpriseDeletePrice
// @Router       /manager/v1/enterprise/price/{id} [POST]
// @Summary      M9017 删除价格
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "价格ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*enterprise) DeletePrice(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    service.NewEnterpriseWithModifier(ctx.Modifier).DeletePrice(req)
    return ctx.SendResponse()
}

// ModifyContract
// @ID           ManagerEnterpriseModifyContract
// @Router       /manager/v1/enterprise/contract [POST]
// @Summary      M9018 编辑合同
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.EnterpriseContractModifyReq  true  "合同字段"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*enterprise) ModifyContract(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseContractModifyReq](c)
    service.NewEnterpriseWithModifier(ctx.Modifier).ModifyContract(req)
    return ctx.SendResponse()
}

// DeleteContract
// @ID           ManagerEnterpriseDeleteContract
// @Router       /manager/v1/enterprise/contract/{id} [DELETE]
// @Summary      M9019 删除合同
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "合同ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*enterprise) DeleteContract(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    service.NewEnterpriseWithModifier(ctx.Modifier).DeleteContract(req)
    return ctx.SendResponse()
}
