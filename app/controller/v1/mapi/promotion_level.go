package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionLevel struct {
}

var PromotionLevel = new(promotionLevel)

// List
// @ID           ManagerPromotionLevelList
// @Router       /manager/v1/promotion/level [GET]
// @Summary      PM3006 会员等级配置列表
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []promotion.Level
func (l *promotionLevel) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionLevelService().MemberLevel())
}

// Update
// @ID           ManagerPromotionLevelUpdate
// @Router       /manager/v1/promotion/level [PUT]
// @Summary      PM3007 更新会员等级配置
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body promotion.Level true "会员等级配置"
// @Success      200  {object}  model.StatusResponse
func (l *promotionLevel) Update(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.Level](c)
	service.NewPromotionLevelService(ctx.Modifier).UpdateMemberLevel(req)
	return ctx.SendResponse()
}

// Create
// @ID           ManagerPromotionLevelCreate
// @Router       /manager/v1/promotion/level [POST]
// @Summary      PM3008 创建会员等级配置
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body promotion.Level true "会员等级配置"
// @Success      200  {object}  model.StatusResponse
func (l *promotionLevel) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.Level](c)
	service.NewPromotionLevelService(ctx.Modifier).CreateMemberLevel(req)
	return ctx.SendResponse()
}

// Delete
// @ID           ManagerPromotionLevelDelete
// @Router       /manager/v1/promotion/level [DELETE]
// @Summary      PM3009 删除会员等级配置
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  query int true "会员等级配置id"
// @Success      200  {object}  model.StatusResponse
func (l *promotionLevel) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewPromotionLevelService(ctx.Modifier).DeleteMemberLevel(req.ID)
	return ctx.SendResponse()
}

// Selection
// @ID           ManagerPromotionLevelSelection
// @Router       /manager/v1/promotion/level/selection [GET]
// @Summary      PM3010 会员等级配置下拉选择列表
// @Tags         [PM]推广管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []promotion.Level
func (l *promotionLevel) Selection(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionLevelService().LevelSelection())
}