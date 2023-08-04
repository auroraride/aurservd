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
// @Summary      P1007 会员收益列表
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  promotion.EarningsReq true  "查询请求"
// @Success      200  {object}  []promotion.EarningsRes
func (e *promotionEarnings) List(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.EarningsReq](c)
	return ctx.SendResponse(service.NewPromotionEarningsService().EarningsList(&promotion.EarningsReq{
		ID:             &ctx.Member.ID,
		PaginationReq:  req.PaginationReq,
		EarningsFilter: req.EarningsFilter,
	}))
}
