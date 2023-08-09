package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type member struct {
}

var Member = new(member)

// List
// @ID           ManagerPromotionMemberList
// @Router       /manager/v1/promotion/member [GET]
// @Summary      PM1001 会员列表
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body promotion.MemberReq  true  "会员列表请求"
// @Success      200  {object}  []promotion.MemberRes  "请求成功"
func (m member) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.MemberReq](c)
	return ctx.SendResponse(service.NewPromotionMemberService().List(req))
}

// Detail
// @ID           ManagerPromotionMemberDetail
// @Router       /manager/v1/promotion/member/{id} [GET]
// @Summary      PM1002 会员详情
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "会员ID"
// @Success      200  {object}  promotion.MemberRes
func (m member) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.MemberReq](c)
	return ctx.SendResponse(service.NewPromotionMemberService().Detail(req))
}

// Update
// @ID           ManagerPromotionMemberUpdate
// @Router       /manager/v1/promotion/member/{id} [PUT]
// @Summary      PM1003 更新会员信息
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body promotion.MemberUpdateReq  true  "会员更新请求"
// @Success      200  {object}  model.StatusResponse
func (m member) Update(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.MemberUpdateReq](c)
	service.NewPromotionMemberService(ctx.Modifier).Update(req)
	return ctx.SendResponse()
}

// TeamList
// @ID           ManagerPromotionMemberTeamList
// @Router       /manager/v1/promotion/member/team/{id} [GET]
// @Summary      PM1004 会员团队列表
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body promotion.MemberTeamReq true  "会员团队列表请求"
// @Success      200  {object}  model.PaginationRes{items=[]promotion.MemberTeamRes}
func (m *member) TeamList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.MemberTeamReq](c)
	return ctx.SendResponse(service.NewPromotionMemberService().TeamList(ctx, req))
}

// SetCommission
// @ID           ManagerPromotionMemberSetCommission
// @Router       /manager/v1/promotion/member/setcommission [POST]
// @Summary      PM1005 设置会员佣金方案
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body promotion.MemberCommissionReq true  "设置会员佣金请求"
// @Success      200  {object}  model.StatusResponse
func (m *member) SetCommission(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.MemberCommissionReq](c)
	service.NewPromotionMemberService(ctx.Modifier).SetCommission(req)
	return ctx.SendResponse()
}
