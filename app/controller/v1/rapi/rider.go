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
    c := ctx.(*app.Context)
    err = c.BindValidate(req)
    if err != nil {
        return
    }

    // 校验短信
    if !debugPhones[req.Phone] && !service.NewSms().VerifyCode(req.SmsId, req.SmsCode) {
        return errors.New("短信验证码校验失败")
    }

    // 注册+登录
    var data *model.RiderSigninRes
    data, err = service.NewRider().Signin(req.Phone, c.Device)
    if err != nil {
        return
    }
    res := response.New(ctx).SetData(data)
    switch data.TokenPermission {
    case model.RiderTokenPermissionAuth:
        res.Error(response.StatusNotAcceptable).SetMessage("需要实名认证")
    case model.RiderTokenPermissionNewDevice:
        res.Error(response.StatusForbidden).SetMessage("需要验证本人")
    }
    return res.Send()
}
