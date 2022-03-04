// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
)

type csc struct{}

var Csc = new(csc)

// IvrShiguangju
// @ID           CscIvrShiguangju
// @Router       /manager/v1/csc/irv [POST]
// @Summary      MT1. 时光驹催费工具
// @Tags         [M]管理接口
// @Accept       mpfd
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        file  formData  file  true  "外呼文件"
// @Success      200  {object}  model.ShiguangjuIVRRes  "请求成功"
func (*csc) IvrShiguangju(c echo.Context) (err error) {
    file, err := c.FormFile("file")
    if err != nil {
        snag.Panic(err)
    }

    return app.NewResponse(c).
        SetData(service.NewCSC().ParseNameListShiguangju(file)).
        Send()
}
