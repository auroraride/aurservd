package router

import (
	"github.com/auroraride/aurservd/app/controller/v2/assetapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadAssersRoutes() {
	g := root.Group("manager/v2")

	asset := g.Group("/asset")

	asset.Use(middleware.ManagerMiddleware())

	// 基础档案
	asset.POST("", assetapi.Assets.Create)       // 创建资产
	asset.GET("", assetapi.Assets.List)          // 资产列表
	asset.GET("/:id", assetapi.Assets.Detail)    // 资产详情
	asset.PUT("/:id", assetapi.Assets.Update)    // 更新资产
	asset.DELETE("/:id", assetapi.Assets.Delete) // 删除资产
	asset.GET("/count", assetapi.Assets.Count)   // 资产数量

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
	asset.GET("/template", assetapi.Assets.Template)  // 导出模版

	// 仓库
	asset.POST("/warehouse", assetapi.Warehouse.Create)        // 创建仓库
	asset.GET("/warehouse", assetapi.Warehouse.List)           // 仓库列表
	asset.GET("/warehouse/:id", assetapi.Warehouse.Detail)     // 仓库详情
	asset.PUT("/warehouse/:id", assetapi.Warehouse.Modify)     // 更新仓库
	asset.DELETE("/warehouse/:id", assetapi.Warehouse.Delete)  // 删除仓库
	asset.GET("/warehouse_assets ", assetapi.Warehouse.Assets) // 仓库物资

	// 门店物资
	asset.GET("/store_assets", assetapi.Store.StoreAsset) // 门店物资列表

	// 运维物资
	asset.GET("/maintainer_assets", assetapi.Maintainer.MaintainerAsset) // 运维物资列表

	// 电柜物资
	asset.GET("/cabinet_assets", assetapi.Cabinet.CabinetAsset) // 电柜物资列表

	// 团签物资
	asset.GET("/cabinet_assets", assetapi.Cabinet.CabinetAsset) // 团签物资列表

	// 其他物资
	asset.POST("/material", assetapi.Material.Create)       // 创建仓库
	asset.GET("/material", assetapi.Material.List)          // 仓库列表
	asset.PUT("/material/:id", assetapi.Material.Modify)    // 更新仓库
	asset.DELETE("/material/:id", assetapi.Material.Delete) // 删除仓库

	// 资产调拨
	asset.POST("/transfer", assetapi.AssetTransfer.Transfer)                 // 资产调拨
	asset.GET("/transfer", assetapi.AssetTransfer.TransferList)              // 资产调拨列表
	asset.GET("/transfer/:id", assetapi.AssetTransfer.TransferDetail)        // 资产调拨详情
	asset.PUT("/transfer/cancel/:id", assetapi.AssetTransfer.TransferCancel) // 取消资产调拨
	asset.POST("/transfer/receive", assetapi.AssetTransfer.TransferReceive)  // 接收资产
	asset.GET("/transfer/sn/:sn", assetapi.AssetTransfer.GetTransferBySN)    // 根据sn查询调拨单
}
