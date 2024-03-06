package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type member struct {
}

var Member = new(member)

// Signin
// @ID		PromotionMemberSignin
// @Router	/promotion/v1/member/signin [GET]
// @Summary	P1001 登录
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	body	body		promotion.MemberSigninReq	true	"登录请求"
// @Success	200		{object}	promotion.MemberSigninRes	"请求成功"
func (m *member) Signin(c echo.Context) (err error) {
	ctx, req := app.ContextBinding[promotion.MemberSigninReq](c)
	return ctx.SendResponse(service.NewPromotionMemberService().Signin(req))
}

// Signup
// @ID		PromotionMemberSignup
// @Router	/promotion/v1/member/signup [POST]
// @Summary	P1002 邀请注册
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string						true	"会员校验token"
// @Param	body				body		promotion.MemberSigninReq	true	"查询请求"
// @Success	200					{object}	promotion.MemberInviteRes	"请求成功"
func (m *member) Signup(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.MemberSigninReq](c)
	return ctx.SendResponse(service.NewPromotionMemberService().Signup(req))
}

// Profile
// @ID		PromotionMemberProfile
// @Router	/promotion/v1/member/profile [GET]
// @Summary	P1003 会员信息
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string	true	"会员校验token"
// @Success	200					{object}	promotion.MemberProfile
func (m *member) Profile(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionMemberService().MemberProfile(ctx.Member.ID))
}

// ShareQrcode
// @ID		PromotionMemberShareQrcode
// @Router	/promotion/v1/member/share/qrcode [GET]
// @Summary	P1004 获取推广二维码
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string	true	"会员校验token"
// @Success	200					{object}	string	"请求成功"
func (m *member) ShareQrcode(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionMiniProgram().PromotionQrcode(ctx.Member.ID))
}

// UpdateAvatar
// @ID		PromotionMemberUpdateAvatar
// @Router	/promotion/v1/member/avatar [POST]
// @Summary	P1005 更新头像
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string					true	"会员校验token"
// @Param	avatar				formData	file					true	"头像"
// @Success	200					{object}	promotion.UploadAvatar	"请求成功"
func (m *member) UpdateAvatar(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionMemberService().UploadAvatar(ctx))
}

// Team
// @ID		PromotionMemberTeam
// @Router	/promotion/v1/member/team [GET]
// @Summary	P1006 我的团队列表
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string					true	"会员校验token"
// @Param	body				body		promotion.MemberTeamReq	true	"查询请求"
// @Success	200					{object}	model.PaginationRes{items=[]promotion.MemberTeamRes}
func (m *member) Team(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.MemberTeamReq](c)
	return ctx.SendResponse(service.NewPromotionMemberService().TeamList(ctx, &promotion.MemberTeamReq{
		ID:               ctx.Member.ID,
		PaginationReq:    req.PaginationReq,
		MemberTeamFilter: req.MemberTeamFilter,
	}))
}
