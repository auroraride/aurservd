// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
)

var (
    Rider = new(rider)

    // 登录debug
    debugPhones = map[string]bool{
        "18501358308": true,
        "19381630638": true,
    }
)

type rider struct {
}

// Signin 骑手登录
func (*rider) Signin(c echo.Context) (err error) {
    req := new(model.RiderSignupReq)
    ctx := c.(*app.GlobalContext)
    ctx.BindValidate(req)

    // 校验短信
    if !debugPhones[req.Phone] && !service.NewSms().VerifyCode(req.SmsId, req.SmsCode) {
        return errors.New("短信验证码校验失败")
    }

    // 注册+登录
    var data *model.RiderSigninRes
    s := service.NewRider()
    data, err = s.Signin(req.Phone, ctx.Device)
    if err != nil {
        return
    }
    return response.New(c).SetData(data).Send()
}

// Contact 添加紧急联系人
func (r *rider) Contact(c echo.Context) error {
    ctx := c.(*app.RiderContext)
    contact := new(model.RiderContact)
    ctx.BindValidate(contact)
    service.NewRider().UpdateContact(ctx.Rider, contact)
    return response.New(c).Success().Send()
}

// Authenticator 实名认证
func (*rider) Authenticator(c echo.Context) error {
    ctx := c.(*app.RiderContext)
    contact := new(model.RiderContact)
    ctx.BindValidate(contact)
    r := service.NewRider()
    // 更新紧急联系人
    r.UpdateContact(ctx.Rider, contact)
    // 获取人脸识别URL
    return response.New(c).Success().SetData(ar.Map{"url": r.GetFaceAuthUrl(ctx)}).Send()
}

// AuthResult 实名认证结果
// TODO 测试认证失败逻辑
func (r *rider) AuthResult(c echo.Context) error {
    success, err := service.NewRider().FaceAuthResult(c.(*app.RiderContext))
    if err != nil {
        return err
    }
    return response.New(c).Success().SetData(ar.Map{"status": success}).Send()
}

// FaceResult 获取人脸验证结果
func (r *rider) FaceResult(c echo.Context) error {
    success, err := service.NewRider().FaceResult(c.(*app.RiderContext))
    if err != nil {
        return err
    }
    return response.New(c).Success().SetData(ar.Map{"status": success}).Send()
}