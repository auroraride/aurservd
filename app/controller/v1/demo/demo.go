// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package demo

import (
    "context"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/baidu"
    "github.com/labstack/echo/v4"
)

func BaiduFace(c echo.Context) error {
    return response.New(c).Success().SetData(baidu.New().Faceprint()).Send()
}

func BaiduFaceResult(c echo.Context) error {
    token := c.Param("token")
    u, _ := ar.Ent.Rider.Query().First(context.Background())
    service.NewRider().FaceAuthResult(u, token)
    return response.New(c).Success().Send()
}

func BaiduFaceSuccess(c echo.Context) error {
    return nil
}

func BaiduFacefail(c echo.Context) error {
    return nil
}

func AliyunOss(c echo.Context) error {
    ali.NewOss()
    return nil
}
