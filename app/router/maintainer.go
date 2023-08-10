package router

import (
	"github.com/auroraride/aurservd/app/controller/v1/oapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadMaintainerRoutes() {
	g := root.Group("maintainer/v1")

	// 无须校验
	guide := g.Group("", middleware.Maintainer())
	guide.POST("/signin", oapi.Maintainer.Signin)

	// 需校验
	auth := g.Group("", middleware.Maintainer(), middleware.MaintainerAuth())
	auth.GET("/cabinets", oapi.Cabinet.List)
}
