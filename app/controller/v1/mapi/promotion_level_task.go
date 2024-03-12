package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/app/service"
)

type promotionLevelTask struct {
}

var PromotionLevelTask = new(promotionLevelTask)

// List
// @ID		ManagerPromotionTaskList
// @Router	/manager/v1/promotion/task [GET]
// @Summary	会员任务列表
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string	true	"管理员校验token"
// @Success	200				{object}	[]promotion.LevelTask
func (m *promotionLevelTask) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewPromotionLevelTaskService().List())
}

// Update
// @ID		ManagerPromotionTaskUpdate
// @Router	/manager/v1/promotion/task/{id} [PUT]
// @Summary	更新会员任务
// @Tags	[PM]推广管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string				true	"管理员校验token"
// @Param	body			body		promotion.LevelTask	true	"更新参数"
// @Success	200				{object}	string
func (m *promotionLevelTask) Update(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[promotion.LevelTask](c)
	service.NewPromotionLevelTaskService(ctx.Modifier).Update(req)
	return ctx.SendResponse()
}
