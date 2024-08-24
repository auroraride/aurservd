// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-12, by aurb

package router

import (
	"github.com/auroraride/aurservd/app/controller/v2/wapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadWarestoreRoutes() {
	g := root.Group("warestore/v2")

	// 无须校验
	guide := g.Group("", middleware.Warestore())
	guide.POST("/signin", wapi.Warestore.Signin)

	// 需校验
	auth := g.Group("", middleware.Warestore(), middleware.WarestoreAuth())

	// 筛选项
	auth.GET("/selection/warehouse", wapi.Selection.Warehouse)                // 城市仓库
	auth.GET("/selection/manager_warehouse", wapi.Selection.ManagerWarehouse) // 仓管仓库
	auth.GET("/selection/store", wapi.Selection.Store)                        // 城市门店
	auth.GET("/selection/employee_store", wapi.Selection.EmployeeStore)       // 店员门店
	auth.GET("/selection/enterprise", wapi.Selection.Enterprise)              // 城市团签企业
	auth.GET("/selection/maintainer", wapi.Selection.Maintainer)              // 运维
	auth.GET("/selection/station", wapi.Selection.Station)                    // 企业站点
	auth.GET("/selection/model", wapi.Selection.Model)                        // 电池型号筛选
	auth.GET("/selection/ebike/brand", wapi.Selection.EbikeBrand)             // 电车品牌筛选

	// 资产统计
	auth.GET("/asset/count", wapi.Assets.AssetCount) // 资产统计

	// 物资管理
	auth.GET("/assets", wapi.Assets.Assets)              // 资产数据
	auth.GET("/assets/common", wapi.Assets.AssetsCommon) // 电池/电车资产数据

	// 资产调拨
	auth.POST("/transfer", wapi.Warestore.Transfer)                       // 创建调拨
	auth.GET("/transfer", wapi.Warestore.TransferList)                    // 调拨列表
	auth.GET("/transfer/:id", wapi.Warestore.TransferDetail)              // 调拨详情
	auth.POST("/transfer/receive", wapi.Warestore.TransferReceive)        // 调拨批量入库
	auth.GET("/transfer/flow", wapi.Warestore.TransferFlow)               // 资产流转明细
	auth.GET("/transfer/sn/:sn", wapi.Warestore.TransferBySn)             // 根据sn查询调拨信息
	auth.GET("/transfer/details", wapi.AssetTransfer.TransferDetailsList) // 出入库明细

	// 盘点
	auth.GET("/check/sn/:sn", wapi.AssetCheck.GetAssetBySN)      // 通过SN查询资产
	auth.GET("/check", wapi.AssetCheck.List)                     // 盘点记录
	auth.POST("/check", wapi.AssetCheck.Create)                  // 创建资产盘点
	auth.GET("/check/:id", wapi.AssetCheck.Detail)               // 盘点详情
	auth.GET("/check/asset:id", wapi.AssetCheck.AssetDetailList) // 盘点资产明细

}
