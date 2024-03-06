package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
)

type promotionCommission struct {
}

var PromotionCommission = new(promotionCommission)

// CommissionRule
// @ID		PromotionCommissionRule
// @Router	/promotion/v1/commission/rule [GET]
// @Summary	P1007 获取会员分佣规则
// @Tags	[P]推广接口
// @Accept	json
// @Produce	json
// @Param	X-Promotion-Token	header		string							true	"会员校验token"
// @Success	200					{object}	[]promotion.CommissionRuleRes	"请求成功"
func (m *promotionCommission) CommissionRule(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionCommissionService().GetCommissionRule(ctx.Member))
}
