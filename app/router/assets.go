package router

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/service/asset"
)

func loadAssersRoutes() {
	// 电池
	g := root.Group("assets/v1")

	// 电池列表
	g.GET("/download/template", func(c echo.Context) error {
		paht, err := asset.NewBattery().DownloadTemplate()
		if err != nil {
			return err
		}
		return c.Attachment(paht, "导入电池模版.xlsx")
	})
}
