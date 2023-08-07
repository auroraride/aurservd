package mapi

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
// @ID           ManagerPromotionEarningsList
// @Router       /manager/v1/promotion/earnings/{id} [GET]
// @Summary      PM2001 会员收益列表
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  promotion.EarningsReq true  "查询请求"
// @Success      200  {object}  []promotion.EarningsRes
func (m *promotionEarnings) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.EarningsReq](c)
	return ctx.SendResponse(service.NewPromotionEarningsService().List(req))
}

// Cancel
// @ID           ManagerPromotionEarningsCancel
// @Router       /manager/v1/promotion/earnings/cancel [POST]
// @Summary      PM2002 取消会员收益
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "会员收益ID"
// @Success      200  {object}  model.StatusResponse
func (m *promotionEarnings) Cancel(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.EarningsCancelReq](c)
	service.NewPromotionEarningsService(ctx.Modifier).Cancel(req)
	return ctx.SendResponse()
}
