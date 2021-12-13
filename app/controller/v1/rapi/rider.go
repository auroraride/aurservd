// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "context"
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/baidu"
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
func (*rider) Signin(ctx echo.Context) (err error) {
    req := new(model.RiderSignupReq)
    c := ctx.(*app.GlobalContext)
    c.BindValidate(req)

    // 校验短信
    if !debugPhones[req.Phone] && !service.NewSms().VerifyCode(req.SmsId, req.SmsCode) {
        return errors.New("短信验证码校验失败")
    }

    // 注册+登录
    var data *model.RiderSigninRes
    s := service.NewRider()
    data, err = s.Signin(req.Phone, c.Device)
    if err != nil {
        return
    }
    return response.New(ctx).SetData(data).Send()
}

// Contact 添加紧急联系人
func (r *rider) Contact(ctx echo.Context) error {
    c := ctx.(*app.RiderContext)
    contact := new(model.RiderContact)
    c.BindValidate(contact)
    // 更新紧急联系人
    err := ar.Ent.Rider.UpdateOneID(c.Rider.ID).SetContact(contact).Exec(context.Background())
    if err != nil {
        return err
    }
    return nil
}

// Authenticator 实名认证
func (*rider) Authenticator(ctx echo.Context) error {
    // 获取人脸识别URL
    return response.New(ctx).Success().SetData(map[string]string{"url": baidu.New().Faceprint()}).Send()
}
