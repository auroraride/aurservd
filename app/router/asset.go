package router

import (
	"github.com/auroraride/aurservd/app/controller/v2/assetapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadAssetsRoutes() {
	g := root.Group("manager/v2")

	asset := g.Group("/asset")

	asset.POST("/user/signin", assetapi.AssetManager.Signin) // 登录

	asset.Use(middleware.AssetManagerMiddleware())

	// 管理员
	asset.POST("/user", assetapi.AssetManager.Create)
	asset.GET("/user", assetapi.AssetManager.List)
	asset.DELETE("/user/:id", assetapi.AssetManager.Delete)
	asset.PUT("/user/:id", assetapi.AssetManager.Modify)
	asset.GET("/user/profile", assetapi.AssetManager.Profile)

	// 角色权限
	asset.GET("/permission", assetapi.AssetPermission.List)
	asset.GET("/permission/role", assetapi.AssetPermission.ListRole)
	asset.POST("/permission/role", assetapi.AssetPermission.CreateRole)
	asset.PUT("/permission/role/:id", assetapi.AssetPermission.ModifyRole)
	asset.DELETE("/permission/role/:id", assetapi.AssetPermission.DeleteRole)

	// 筛选数据
	asset.GET("/selection/warehouse_city", assetapi.Selection.WarehouseByCity) // 仓库城市筛选
	asset.GET("/selection/city", assetapi.Selection.City)                      // 城市筛选
	asset.GET("/selection/ebike/brand", assetapi.Selection.EbikeBrand)         // 电车品牌筛选
	asset.GET("/selection/store", assetapi.Selection.Store)                    // 门店筛选
	asset.GET("/selection/enterprise", assetapi.Selection.Enterprise)          // 企业筛选
	asset.GET("/selection/role", assetapi.Selection.AssetRole)                 // 角色筛选

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
	asset.POST("/warehouse", assetapi.Warehouse.Create)       // 创建仓库
	asset.GET("/warehouse", assetapi.Warehouse.List)          // 仓库列表
	asset.GET("/warehouse/:id", assetapi.Warehouse.Detail)    // 仓库详情
	asset.PUT("/warehouse/:id", assetapi.Warehouse.Modify)    // 更新仓库
	asset.DELETE("/warehouse/:id", assetapi.Warehouse.Delete) // 删除仓库
	asset.GET("/warehouse_assets", assetapi.Warehouse.Assets) // 仓库物资

	// 门店物资
	asset.GET("/store_assets", assetapi.Store.StoreAsset) // 门店物资列表

	// 运维物资
	asset.GET("/maintainer_assets", assetapi.Maintainer.MaintainerAsset) // 运维物资列表

	// 电柜物资
	asset.GET("/cabinet_assets", assetapi.Cabinet.CabinetAsset) // 电柜物资列表

	// 团签物资
	asset.GET("/enterprise_assets", assetapi.Enterprise.EnterpriseAsset) // 团签物资列表

	// 其他物资
	asset.POST("/material", assetapi.Material.Create)       // 创建其他物资
	asset.GET("/material", assetapi.Material.List)          // 其他物资列表
	asset.PUT("/material/:id", assetapi.Material.Modify)    // 更新其他物资
	asset.DELETE("/material/:id", assetapi.Material.Delete) // 删除其他物资

	// 资产调拨
	asset.POST("/transfer", assetapi.AssetTransfer.Transfer)                   // 资产调拨
	asset.GET("/transfer", assetapi.AssetTransfer.TransferList)                // 资产调拨列表
	asset.GET("/transfer/:id", assetapi.AssetTransfer.TransferDetail)          // 资产调拨详情
	asset.PUT("/transfer/cancel/:id", assetapi.AssetTransfer.TransferCancel)   // 取消资产调拨
	asset.POST("/transfer/receive", assetapi.AssetTransfer.TransferReceive)    // 接收资产
	asset.GET("/transfer/sn/:sn", assetapi.AssetTransfer.GetTransferBySN)      // 根据sn查询调拨单
	asset.GET("/transfer/flow/:sn", assetapi.AssetTransfer.TransferFlow)       // 调拨流转记录
	asset.GET("/transfer/details", assetapi.AssetTransfer.TransferDetailsList) // 调拨详情列表(出入库明细)

	// 资产维修
	asset.GET("/maintenance", assetapi.AssetMaintenance.List)    // 资产维修列表
	asset.POST("/maintenance", assetapi.AssetMaintenance.Create) // 创建维修记录

	// 资产盘点
	asset.GET("/check", assetapi.AssetCheck.List)                                 // 盘点列表
	asset.GET("/check/:id", assetapi.AssetCheck.Detail)                           // 盘点详情
	asset.GET("/check/asset/:id", assetapi.AssetCheck.AssetDetailList)            // 盘点资产列表
	asset.POST("/check", assetapi.AssetCheck.Create)                              // 盘点资产
	asset.GET("/check/sn/:sn", assetapi.AssetCheck.GetAssetBySN)                  // 通过SN查询资产
	asset.GET("/check/abnormal/:id", assetapi.AssetCheck.Abnormal)                // 盘点异常
	asset.PUT("/check/abnormal/operate/:id", assetapi.AssetCheck.AbnormalOperate) // 盘点异常操作

}
