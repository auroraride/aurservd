package mapi

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
// @ID           ManagerPromotionWithdrawalList
// @Router       /manager/v1/promotion/withdrawal [GET]
// @Summary      PM3001 会员提现列表
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  promotion.WithdrawalListReq true  "请求参数"
// @Success      200  {object}  []promotion.WithdrawalListRes
func (p *promotionWithdrawal) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.WithdrawalListReq](c)
	return ctx.SendResponse(service.NewPromotionWithdrawalService().List(ctx, req))
}

// AlterReview
// @ID           ManagerPromotionWithdrawalAlterReview
// @Router       /manager/v1/promotion/withdrawal/alter/review [POST]
// @Summary      PM3002  审批提现
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  promotion.WithdrawalApprovalReq true  "请求参数"
// @Success      200  {object}  model.StatusResponse
func (p *promotionWithdrawal) AlterReview(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.WithdrawalApprovalReq](c)
	service.NewPromotionWithdrawalService().AlterReview(req)
	return ctx.SendResponse()
}

// Export
// @ID           ManagerPromotionWithdrawalExport
// @Router       /manager/v1/promotion/withdrawal/export [GET]
// @Summary       PM3003  导出提现列表
// @Tags          [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Promotion-Token  header  string  true  "会员校验token"
// @Param        body  body  promotion.WithdrawalExportReq true  "请求参数"
// @Success      200  {object}  string
func (p *promotionWithdrawal) Export(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.WithdrawalExportReq](c)
	return ctx.Attachment(service.NewPromotionWithdrawalService().Export(req))
}
