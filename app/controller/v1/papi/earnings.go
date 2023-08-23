package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionEarnings struct {
}

var PromotionEarnings = new(promotionEarnings)

// List
// @ID           PromotionEarningsList
// @Router       /promotion/v1/earnings [GET]
// @Summary      P3001 会员收益列表
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "管理员校验token"
// @Param        body  body  promotion.EarningsReq true  "查询请求"
// @Success      200  {object}  []promotion.EarningsRes
func (m *promotionEarnings) List(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.EarningsReq](c)
	return ctx.SendResponse(service.NewPromotionEarningsService().List(&promotion.EarningsReq{
		ID:             &ctx.Member.ID,
		PaginationReq:  req.PaginationReq,
		EarningsFilter: req.EarningsFilter,
	}))
}

// Total 用户总收益
// @ID           PromotionEarningsTotal
// @Router       /promotion/v1/earnings/total [GET]
// @Summary      P3002 获取总收益
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "管理员校验token"
// @Param        body  body  promotion.EarningsReq true  "查询请求"
// @Success      200  {object}  []promotion.EarningsRes
func (m *promotionEarnings) Total(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.EarningsReq](c)
	return ctx.SendResponse(service.NewPromotionEarningsService().TotalEarnings(&promotion.EarningsReq{
		ID:             &ctx.Member.ID,
		EarningsFilter: req.EarningsFilter,
	}))
}
