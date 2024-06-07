package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type rider struct{}

var Rider = new(rider)

// ChangePhone 修改手机号
// @ID		RiderChangePhone
// @Router	/rider/v2/change/phone [POST]
// @Summary	修改手机号
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		definition.RiderChangePhoneReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*rider) ChangePhone(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.RiderChangePhoneReq](c)
	return ctx.SendResponse(biz.NewRiderBiz().ChangePhone(ctx.Rider, req))
}

// Signin
// @ID		Signin
// @Router	/rider/v2/signin [POST]
// @Summary	登录或注册
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	body	body		definition.RiderSignupReq	true	"desc"
// @Success	200		{object}	model.RiderSigninRes		"请求成功"
func (*rider) Signin(c echo.Context) (err error) {
	ctx, req := app.ContextBinding[definition.RiderSignupReq](c)
	return ctx.SendResponse(biz.NewRiderBiz().Signin(ctx.Device, req))
}

// GetOpenid
// @ID		GetOpenid
// @Router	/rider/v2/mini/openid [GET]
// @Summary	获取支付宝小程序openid
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string			true	"骑手校验token"
// @Param	code			query		string			true	"支付宝code"
// @Success	200				{object}	model.OpenidRes	"请求成功"
func (*rider) GetOpenid(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.OpenidReq](c)
	return ctx.SendResponse(biz.NewRiderBiz().GetAlipayOpenid(req))
}

// SetMobPushId
// @ID		RiderSetMobPush
// @Router	/rider/v2/mobpush [POST]
// @Summary	设置推送ID
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		definition.RiderSetMobPushReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*rider) SetMobPushId(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.RiderSetMobPushReq](c)
	return ctx.SendResponse(biz.NewRiderBiz().SetMobPushId(ctx.Rider, req))
}

// Profile
// @ID		Profile
// @Router	/rider/v2/profile [GET]
// @Summary 获取个人信息
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	model.RiderSigninRes	"请求成功"
func (*rider) Profile(c echo.Context) error {
	ctx := c.(*app.RiderContext)
	profile, err := biz.NewRiderBiz().Profile(ctx.Rider, ctx.Device, ctx.Token)
	profile.Token = ctx.Token
	return ctx.SendResponse(profile, err)
}
