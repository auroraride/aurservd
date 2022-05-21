// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
)

var (
    Rider = new(rider)
)

type rider struct {
}

// Signin
// @ID           RiderSignin
// @Router       /rider/v1/signin [POST]
// @Summary      R10001 登录或注册
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        body  body  model.RiderSignupReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Signin(c echo.Context) (err error) {
    ctx, req := app.ContextBinding[model.RiderSignupReq](c)

    // 校验短信
    debugPhones := ar.Config.App.Debug.Phone
    if !debugPhones[req.Phone] && !service.NewSms().VerifyCode(req.SmsId, req.SmsCode) {
        return errors.New("短信验证码校验失败")
    }

    // 注册+登录
    var data *model.RiderSigninRes
    s := service.NewRider()
    data, err = s.Signin(ctx.Device, req)
    if err != nil {
        return
    }
    return ctx.SendResponse(data)
}

// Contact
// @ID           RiderContact
// @Router       /rider/v1/contact [POST]
// @Summary      R10002 添加紧急联系人
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.RiderContact  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (r *rider) Contact(c echo.Context) error {
    ctx, req := app.RiderContextAndBinding[model.RiderContact](c)
    service.NewRider().UpdateContact(ctx.Rider, req)
    return ctx.SendResponse()
}

// Authenticator 实名认证
func (*rider) Authenticator(c echo.Context) error {
    ctx, req := app.RiderContextAndBinding[model.RiderContact](c)
    r := service.NewRider()
    // 更新紧急联系人
    r.UpdateContact(ctx.Rider, req)
    // 获取人脸识别URL
    return ctx.SendResponse(model.FaceAuthUrlResponse{
        Url: r.GetFaceAuthUrl(ctx),
    })
}

// AuthResult 实名认证结果
// TODO 测试认证失败逻辑
func (r *rider) AuthResult(c echo.Context) error {
    ctx := c.(*app.RiderContext)
    success, err := service.NewRider().FaceAuthResult(ctx)
    if err != nil {
        snag.Panic(err)
    }
    return ctx.SendResponse(model.StatusResponse{Status: success})
}

// FaceResult 获取人脸校验结果
func (r *rider) FaceResult(c echo.Context) error {
    ctx := c.(*app.RiderContext)
    success, err := service.NewRider().FaceResult(ctx)
    if err != nil {
        snag.Panic(err)
    }
    return ctx.SendResponse(model.StatusResponse{Status: success})
}

func (r *rider) Demo(c echo.Context) error {
    ctx := c.(*app.RiderContext)
    return ctx.SendResponse(model.StatusResponse{Status: true})
}

// Profile 用户信息
func (r *rider) Profile(c echo.Context) error {
    ctx := c.(*app.RiderContext)
    return ctx.SendResponse(model.StatusResponse{Status: true})
}
