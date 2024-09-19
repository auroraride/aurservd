package router

import (
	"github.com/auroraride/aurservd/app/controller/v2/amapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadAssetsRoutes() {
	g := root.Group("manager/v2")

	asset := g.Group("/asset")

	asset.POST("/user/signin", amapi.AssetManager.Signin) // 登录

	asset.Use(middleware.AssetManagerMiddleware())

	// 管理员
	asset.POST("/user", amapi.AssetManager.Create)
	asset.GET("/user", amapi.AssetManager.List)
	asset.DELETE("/user/:id", amapi.AssetManager.Delete)
	asset.PUT("/user/:id", amapi.AssetManager.Modify)
	asset.GET("/user/profile", amapi.AssetManager.Profile)

	// 店员
	asset.POST("/employee", amapi.Employee.Create)       // 创建店员
	asset.GET("/employee", amapi.Employee.List)          // 店员列表
	asset.PUT("/employee/:id", amapi.Employee.Modify)    // 更新店员
	asset.DELETE("/employee/:id", amapi.Employee.Delete) // 删除店员

	// 角色权限
	asset.GET("/permission", amapi.AssetPermission.List)
	asset.GET("/permission/role", amapi.AssetPermission.ListRole)
	asset.POST("/permission/role", amapi.AssetPermission.CreateRole)
	asset.PUT("/permission/role/:id", amapi.AssetPermission.ModifyRole)
	asset.DELETE("/permission/role/:id", amapi.AssetPermission.DeleteRole)

	// 筛选数据
	asset.GET("/selection/warehouse_city", amapi.Selection.WarehouseByCity) // 城市仓库筛选
	asset.GET("/selection/city", amapi.Selection.City)                      // 城市筛选
	asset.GET("/selection/ebike/brand", amapi.Selection.EbikeBrand)         // 电车品牌筛选
	asset.GET("/selection/store", amapi.Selection.Store)                    // 城市门店筛选
	asset.GET("/selection/enterprise", amapi.Selection.Enterprise)          // 企业筛选
	asset.GET("/selection/role", amapi.Selection.AssetRole)                 // 角色筛选
	asset.GET("/selection/model", amapi.Selection.Model)                    // 电池型号筛选
	asset.GET("/selection/maintainer", amapi.Selection.Maintainer)          // 运维人员筛选
	asset.GET("/selection/station", amapi.Selection.Station)                // 站点筛选
	asset.GET("/selection/material", amapi.Selection.Material)              // 物资类型筛选
	asset.GET("/selection/cabinet", amapi.Selection.Cabinet)                // 筛选电柜

	// 基础档案
	asset.POST("", amapi.Assets.Create)     // 创建资产
	asset.GET("", amapi.Assets.List)        // 资产列表
	asset.GET("/:id", amapi.Assets.Detail)  // 资产详情
	asset.PUT("/:id", amapi.Assets.Update)  // 更新资产
	asset.GET("/count", amapi.Assets.Count) // 资产数量

	// 资产属性
	asset.GET("/attributes", amapi.AssetAttributes.List) // 资产属性列表

	// 报废
	asset.POST("/scrap", amapi.AssetScrap.Scrap)                           // 报废资产
	asset.POST("/scrap/batch/restore", amapi.AssetScrap.ScrapBatchRestore) // 批量还原报废
	asset.GET("/scrap", amapi.AssetScrap.ScrapList)                        // 报废列表
	asset.GET("/scrap/reason", amapi.AssetScrap.ScrapReasonSelect)         // 报废理由列表
	// 导入导出
	asset.POST("/batch", amapi.Assets.BatchCreate) // 导入资产
	asset.POST("/export", amapi.Assets.Export)     // 导出资产
	asset.GET("/template", amapi.Assets.Template)  // 导出模版

	// 仓库
	asset.POST("/warehouse", amapi.Warehouse.Create)                // 创建仓库
	asset.GET("/warehouse", amapi.Warehouse.List)                   // 仓库列表
	asset.GET("/warehouse/:id", amapi.Warehouse.Detail)             // 仓库详情
	asset.PUT("/warehouse/:id", amapi.Warehouse.Modify)             // 更新仓库
	asset.DELETE("/warehouse/:id", amapi.Warehouse.Delete)          // 删除仓库
	asset.GET("/warehouse_assets", amapi.Warehouse.Assets)          // 仓库物资
	asset.GET("/warehouse_assets/:id", amapi.Warehouse.AssetDetail) // 仓库物资详情

	// 城市
	asset.GET("/city", amapi.City.List)       // 城市列表
	asset.PUT("/city/:id", amapi.City.Modify) // 启用或关闭城市

	// 门店集合
	asset.GET("/store_group", amapi.StoreGroup.List)          // 门店集合列表
	asset.POST("/store_group", amapi.StoreGroup.Create)       // 创建门店集合
	asset.DELETE("/store_group/:id", amapi.StoreGroup.Delete) // 删除门店集合

	// 门店物资
	asset.GET("/store_assets", amapi.Store.Asset)           // 门店物资列表
	asset.GET("/store_assets/:id", amapi.Store.AssetDetail) // 门店物资详情

	// 运维物资
	asset.GET("/maintainer_assets", amapi.Maintainer.Asset)           // 运维物资列表
	asset.GET("/maintainer_assets/:id", amapi.Maintainer.AssetDetail) // 运维物资详情

	// 电柜物资
	asset.GET("/cabinet_assets", amapi.Cabinet.Asset)           // 电柜物资列表
	asset.GET("/cabinet_assets/:id", amapi.Cabinet.AssetDetail) // 电柜物资详情

	// 团签物资
	asset.GET("/enterprise_assets", amapi.Enterprise.Asset)           // 团签物资列表
	asset.GET("/enterprise_assets/:id", amapi.Enterprise.AssetDetail) // 团签物资详情

	// 电池型号
	asset.POST("/batterymodel", amapi.BatteryModel.Create)       // 创建电池型号
	asset.GET("/batterymodel", amapi.BatteryModel.List)          // 电池型号列表
	asset.DELETE("/batterymodel/:id", amapi.BatteryModel.Delete) // 删除电池型号

	// 电车型号
	asset.GET("/ebike/brand", amapi.EbikeBrand.List)          // 电池型号列表
	asset.POST("/ebike/brand", amapi.EbikeBrand.Create)       // 创建电池型号
	asset.PUT("/ebike/brand/:id", amapi.EbikeBrand.Modify)    // 更新电池型号
	asset.DELETE("/ebike/brand/:id", amapi.EbikeBrand.Delete) // 删除电池型号

	// 其他物资
	asset.POST("/material", amapi.Material.Create)       // 创建其他物资
	asset.GET("/material", amapi.Material.List)          // 其他物资列表
	asset.PUT("/material/:id", amapi.Material.Modify)    // 更新其他物资
	asset.DELETE("/material/:id", amapi.Material.Delete) // 删除其他物资

	// 资产调拨
	asset.POST("/transfer", amapi.AssetTransfer.Transfer)                   // 资产调拨
	asset.GET("/transfer", amapi.AssetTransfer.TransferList)                // 资产调拨列表
	asset.GET("/transfer/:id", amapi.AssetTransfer.TransferDetail)          // 资产调拨详情
	asset.PUT("/transfer/cancel/:id", amapi.AssetTransfer.TransferCancel)   // 取消资产调拨
	asset.POST("/transfer/receive", amapi.AssetTransfer.TransferReceive)    // 接收资产
	asset.GET("/transfer/flow", amapi.AssetTransfer.TransferFlow)           // 调拨流转记录
	asset.GET("/transfer/details", amapi.AssetTransfer.TransferDetailsList) // 调拨详情列表(出入库明细)
	asset.PUT("/transfer/:id", amapi.AssetTransfer.Modify)                  // 修改调拨记录

	// 资产维修
	asset.GET("/maintenance", amapi.AssetMaintenance.List)       // 资产维修列表
	asset.PUT("/maintenance/:id", amapi.AssetMaintenance.Modify) // 修改维修记录

	// 资产盘点
	asset.GET("/check", amapi.AssetCheck.List)                                 // 盘点列表
	asset.GET("/check/:id", amapi.AssetCheck.Detail)                           // 盘点详情
	asset.GET("/check/asset/:id", amapi.AssetCheck.AssetDetailList)            // 盘点资产列表
	asset.POST("/check", amapi.AssetCheck.Create)                              // 盘点资产
	asset.GET("/check/sn/:sn", amapi.AssetCheck.GetAssetBySN)                  // 通过SN查询资产
	asset.GET("/check/abnormal/:id", amapi.AssetCheck.Abnormal)                // 盘点异常
	asset.PUT("/check/abnormal/operate/:id", amapi.AssetCheck.AbnormalOperate) // 盘点异常操作

}
