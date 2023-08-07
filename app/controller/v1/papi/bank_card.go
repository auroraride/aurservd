package papi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionBankCard struct {
}

var PromotionBankCard = new(promotionBankCard)

// Create
// @ID           PromotionBankCardCreate
// @Router       /promotion/v1/bank/card [POST]
// @Summary      P4001 创建银行卡
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Param        body  body  promotion.BankCardReq true  "请求参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (m *promotionBankCard) Create(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[promotion.BankCardReq](c)
	service.NewPromotionBankCardService().Create(ctx.Member, req)
	return ctx.SendResponse()
}

// Update
// @ID           PromotionBankCardUpdate
// @Router       /promotion/v1/bank/card/{id} [PUT]
// @Summary      P4002 修改银行卡默认状态
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Param        body  body  promotion.BankCardReq true  "请求参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (m *promotionBankCard) Update(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[model.IDParamReq](c)
	service.NewPromotionBankCardService().Update(ctx.Member, req)
	return ctx.SendResponse()
}

// List
// @ID           PromotionBankCardList
// @Router       /promotion/v1/bank/card [GET]
// @Summary      P4003 获取银行卡列表
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Success      200  {object}  []promotion.BankCardRes  "请求成功"
func (m *promotionBankCard) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.PromotionContext](c)
	return ctx.SendResponse(service.NewPromotionBankCardService().List(ctx.Member))
}

// Delete
// @ID           PromotionBankCardDelete
// @Router       /promotion/v1/bank/card/{id} [DELETE]
// @Summary      P4004 删除银行卡
// @Tags         [P]推广接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Success      200  {object}  []promotion.BankCardRes  "请求成功"
func (m *promotionBankCard) Delete(c echo.Context) (err error) {
	ctx, req := app.PromotionContextAndBinding[model.IDParamReq](c)
	service.NewPromotionBankCardService().Delete(ctx.Member, req)
	return ctx.SendResponse()
}
