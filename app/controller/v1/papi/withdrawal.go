package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionWithdrawal struct {
}

var PromotionWithdrawal = new(promotionWithdrawal)

// List
// @ID           PromotionWithdrawalList
// @Router       /promotion/v1/withdrawal [GET]
// @Summary      P2001 会员提现列表
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Param        body  body  promotion.WithdrawalListReq true  "请求参数"
// @Success      200  {object}  []promotion.WithdrawalListRes
func (l *promotionWithdrawal) List(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.WithdrawalListReq](c)
	return ctx.SendResponse(service.NewPromotionWithdrawalService().List(ctx, &promotion.WithdrawalListReq{
		ID:               &ctx.Member.ID,
		PaginationReq:    req.PaginationReq,
		WithdrawalFilter: req.WithdrawalFilter,
	}))
}

// Alter
// @ID           PromotionWithdrawalAlter
// @Router       /promotion/v1/withdrawal/alter [POST]
// @Summary      P2002 申请提现
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Param        body  body  promotion.WithdrawalAlterReq true  "请求参数"
// @Success      200  {object}  model.StatusResponse
func (l *promotionWithdrawal) Alter(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.WithdrawalAlterReq](c)
	service.NewPromotionWithdrawalService().Alter(ctx.Member, req)
	return ctx.SendResponse()
}

// CalculateWithdrawalFee
// @ID           PromotionWithdrawalFee
// @Router       /promotion/v1/withdrawal/fee [POST]
// @Summary      P2003 计算提现手续费
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Param        body  query  promotion.WithdrawalAlterReq  true  "提现金额"
// @Success      200  {object}  promotion.WithdrawalFeeRes
func (l *promotionWithdrawal) CalculateWithdrawalFee(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.WithdrawalAlterReq](c)
	return ctx.SendResponse(service.NewPromotionWithdrawalService().CalculateWithdrawalFee(ctx.Member, req))
}

// Total
// @ID           PromotionWithdrawalTotal
// @Router       /promotion/v1/withdrawal/total [GET]
// @Summary      P2004 查询总提现金额
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Param        body  query  promotion.WithdrawalListReq  true  "提现金额"
// @Success      200  {object}  promotion.TotalRes
func (l *promotionWithdrawal) Total(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.WithdrawalListReq](c)
	return ctx.SendResponse(service.NewPromotionWithdrawalService().TotalWithdrawal(&promotion.WithdrawalListReq{
		ID:               &ctx.Member.ID,
		WithdrawalFilter: req.WithdrawalFilter,
	}))
}
