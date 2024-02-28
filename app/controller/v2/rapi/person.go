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

// CertificationOcr
// @ID		CertificationOcr
// @Router	/rider/v2/certification/ocr [GET]
// @Summary	获取人脸核身OCR参数
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Success	200	{object}	definition.PersonCertificationOcrRes	"请求成功"
func (*person) CertificationOcr(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(biz.NewPerson().CertificationOcr(ctx.Rider))
}

// CertificationFace
// @ID		CertificationFace
// @Router	/rider/v2/certification/face [POST]
// @Summary	提交身份信息并获取人脸核身参数
// @Tags	Person - 实人
// @Accept	json
// @Produce	json
// @Param	query	query		definition.PersonCertificationFaceReq	true	"请求参数"
// @Success	200		{object}	definition.PersonCertificationFaceRes	"请求成功"
func (*person) CertificationFace(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PersonCertificationFaceReq](c)
	return ctx.SendResponse(biz.NewPerson().CertificationFace(ctx.Rider, req))
}

// CertificationFaceResult
// @ID		CertificationFaceResult
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
