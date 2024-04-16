// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type person struct{}

var Person = new(person)

// CertificationOcrClient
// @ID		PersonCertificationOcrClient
// @Router	/rider/v2/certification/ocr/client [GET]
// @Summary	获取客户端OCR参数
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Success	200	{object}	definition.PersonCertificationOcrClientRes	"请求成功"
func (*person) CertificationOcrClient(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(biz.NewPerson().CertificationOcrClient(ctx.Rider))
}

// CertificationOcrCloud
// @ID		PersonCertificationOcrCloud
// @Router	/rider/v2/certification/ocr/cloud [GET]
// @Summary	获取云端OCR参数
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Param	query	query		definition.PersonCertificationOcrCloudReq	true	"请求参数"
// @Success	200		{object}	definition.PersonCertificationOcrCloudRes	"请求成功"
func (*person) CertificationOcrCloud(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationOcrCloudReq](c)
	return ctx.SendResponse(biz.NewPerson().CertificationOcrCloud(req.Hash))
}

// CertificationFace
// @ID		PersonCertificationFace
// @Router	/rider/v2/certification/face [POST]
// @Summary	提交身份信息并获取人脸核身参数
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Param	body	body		definition.PersonCertificationFaceReq	true	"请求参数"
// @Success	200		{object}	definition.PersonCertificationFaceRes	"请求成功"
func (*person) CertificationFace(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationFaceReq](c)
	return ctx.SendResponse(biz.NewPerson().CertificationFace(ctx.Rider, req))
}

// CertificationFaceResult
// @ID		PersonCertificationFaceResult
// @Router	/rider/v2/certification/face/result [GET]
// @Summary	获取人脸核身结果
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Param	query	query		definition.PersonCertificationFaceResultReq	true	"请求参数"
// @Success	200		{object}	definition.PersonCertificationFaceResultRes	"请求成功"
func (p *person) CertificationFaceResult(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationFaceResultReq](c)
	return ctx.SendResponse(biz.NewPerson().CertificationFaceResult(ctx.Rider, req))
}

// CertificationSupplement
// @ID		PersonCertificationSupplement
// @Router	/rider/v2/certification/supplement [POST]
// @Summary	提交人脸核身补充信息
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Param	body	body		definition.PersonCertificationSupplementReq	true	"请求参数"
// @Success	200		{object}	model.StatusResponse						"请求成功"
func (*person) CertificationSupplement(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationSupplementReq](c)
	return ctx.SendResponse(biz.NewPerson().CertificationSupplement(ctx.Rider, req))
}
