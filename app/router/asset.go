package router

import (
	"github.com/auroraride/aurservd/app/controller/v2/mapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadAssersRoutes() {
	g := root.Group("manager/v2")

	asset := g.Group("/asset")

	asset.GET("/template", mapi.Assets.Template) // 导出模版

	asset.Use(middleware.ManagerMiddleware())

	// 基础档案
	asset.POST("", mapi.Assets.Create)       // 创建资产
	asset.GET("", mapi.Assets.List)          // 资产列表
	asset.GET("/:id", mapi.Assets.Detail)    // 资产详情
	asset.PUT("/:id", mapi.Assets.Update)    // 更新资产
	asset.DELETE("/:id", mapi.Assets.Delete) // 删除资产

	// 资产属性
	asset.GET("/attributes", mapi.AssetAttributes.List) // 资产属性列表

	// 报废
	asset.POST("/scrap", mapi.AssetScrap.Scrap)                           // 报废资产
	asset.POST("/scrap/batch/restore", mapi.AssetScrap.ScrapBatchRestore) // 批量还原报废
	asset.GET("/scrap", mapi.AssetScrap.ScrapList)                        // 报废列表
	asset.GET("/scrap/reason", mapi.AssetScrap.ScrapReasonSelect)         // 报废理由列表
	// 导入导出
	asset.POST("/batch", mapi.Assets.BatchCreate) // 导入资产
	asset.POST("/export", mapi.Assets.Export)     // 导出资产

}
