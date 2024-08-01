package router

import (
	"github.com/auroraride/aurservd/app/controller/v2/assetapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadAssersRoutes() {
	g := root.Group("manager/v2")

	asset := g.Group("/asset")

	asset.GET("/template", assetapi.Assets.Template) // 导出模版

	asset.Use(middleware.ManagerMiddleware())

	// 基础档案
	asset.POST("", assetapi.Assets.Create)       // 创建资产
	asset.GET("", assetapi.Assets.List)          // 资产列表
	asset.GET("/:id", assetapi.Assets.Detail)    // 资产详情
	asset.PUT("/:id", assetapi.Assets.Update)    // 更新资产
	asset.DELETE("/:id", assetapi.Assets.Delete) // 删除资产

	// 资产属性
	asset.GET("/attributes", assetapi.AssetAttributes.List) // 资产属性列表

	// 报废
	asset.POST("/scrap", assetapi.AssetScrap.Scrap)                           // 报废资产
	asset.POST("/scrap/batch/restore", assetapi.AssetScrap.ScrapBatchRestore) // 批量还原报废
	asset.GET("/scrap", assetapi.AssetScrap.ScrapList)                        // 报废列表
	asset.GET("/scrap/reason", assetapi.AssetScrap.ScrapReasonSelect)         // 报废理由列表
	// 导入导出
	asset.POST("/batch", assetapi.Assets.BatchCreate) // 导入资产
	asset.POST("/export", assetapi.Assets.Export)     // 导出资产

	// 资产调拨
	asset.POST("/transfer", assetapi.AssetTransfer.Transfer) // 资产调拨

}
