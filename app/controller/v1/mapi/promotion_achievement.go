package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionAchievement struct {
}

var PromotionAchievement = new(promotionAchievement)

// List
// @ID		ManagerPromotionAchievementList
// @Router	/manager/v1/promotion/achievement [GET]
// @Summary	会员成就列表
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Success	200				{object}	[]promotion.Achievement
func (a *promotionAchievement) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionAchievementService().List())
}

// Create
// @ID		ManagerPromotionAchievementCreate
// @Router	/manager/v1/promotion/achievement [POST]
// @Summary	创建会员成就
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		promotion.Achievement	true	"创建请求"
// @Success	200				{object}	model.StatusResponse
func (a *promotionAchievement) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.Achievement](c)
	service.NewPromotionAchievementService(ctx.Modifier).Create(req)
	return ctx.SendResponse()
}

// Update
// @ID		ManagerPromotionAchievementUpdate
// @Router	/manager/v1/promotion/achievement [PUT]
// @Summary	更新会员成就
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		promotion.Achievement	true	"更新请求"
// @Success	200				{object}	model.StatusResponse
func (a *promotionAchievement) Update(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.Achievement](c)
	service.NewPromotionAchievementService(ctx.Modifier).Update(req)
	return ctx.SendResponse()
}

// Delete
// @ID		ManagerPromotionAchievementDelete
// @Router	/manager/v1/promotion/achievement [DELETE]
// @Summary	删除会员成就
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Param	id				query		string	true	"成就id"
// @Success	200				{object}	model.StatusResponse
func (a *promotionAchievement) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewPromotionAchievementService(ctx.Modifier).Delete(req.ID)
	return ctx.SendResponse()
}

// UploadIcon
// @ID		ManagerPromotionAchievementUploadIcon
// @Router	/manager/v1/promotion/achievement/icon [POST]
// @Summary	上传会员成就图标
// @Tags	[PM]推广管理接口
// @Accept	multipart/form-data
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Param	icon			formData	file	true	"成就图标"
// @Success	200				{object}	promotion.UploadIcon
func (a *promotionAchievement) UploadIcon(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionAchievementService(ctx.Modifier).UploadIcon(ctx))
}
