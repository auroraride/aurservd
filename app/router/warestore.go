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

	// 上班
	auth.POST("/duty/check", wapi.Warestore.CheckDuty) // 检查上班范围
	auth.POST("/duty", wapi.Warestore.OnDuty)          // 上班

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
	auth.GET("/selection/city", wapi.Selection.City)                          // 城市筛选

	// 物资管理
	auth.GET("/assets/count", wapi.Assets.AssetCount)         // 资产统计
	auth.GET("/assets", wapi.Assets.Assets)                   // 资产数据
	auth.GET("/assets/common", wapi.Assets.AssetsCommon)      // 电池/电车资产数据
	auth.PUT("/assets/:id", wapi.Assets.Update)               // 更新资产
	auth.GET("/assets/attributes", wapi.AssetAttributes.List) // 资产属性列表

	// 资产调拨
	auth.POST("/transfer", wapi.AssetTransfer.Transfer)                   // 创建调拨
	auth.GET("/transfer", wapi.AssetTransfer.TransferList)                // 调拨列表
	auth.GET("/transfer/:id", wapi.AssetTransfer.TransferDetail)          // 调拨详情
	auth.POST("/transfer/receive", wapi.AssetTransfer.TransferReceive)    // 调拨批量入库
	auth.GET("/transfer/flow", wapi.AssetTransfer.TransferFlow)           // 资产流转明细
	auth.GET("/transfer/sn/:sn", wapi.AssetTransfer.TransferBySn)         // 根据sn查询调拨信息
	auth.GET("/transfer/details", wapi.AssetTransfer.TransferDetailsList) // 出入库明细
	auth.PUT("/transfer/:id", wapi.AssetTransfer.Modify)                  // 修改调拨记录
	auth.PUT("/transfer/cancel/:id", wapi.AssetTransfer.TransferCancel)   // 取消资产调拨

	// 盘点
	auth.GET("/check/sn/:sn", wapi.AssetCheck.GetAssetBySN)      // 通过SN查询资产
	auth.GET("/check", wapi.AssetCheck.List)                     // 盘点记录
	auth.POST("/check", wapi.AssetCheck.Create)                  // 创建资产盘点
	auth.GET("/check/:id", wapi.AssetCheck.Detail)               // 盘点详情
	auth.GET("/check/asset:id", wapi.AssetCheck.AssetDetailList) // 盘点资产明细

	// 电柜操作
	auth.GET("/cabinet/:serial", wapi.Cabinet.Detail)
	auth.POST("/cabinet/:serial", wapi.Cabinet.Operate)
	auth.POST("/cabinet/:serial/:ordinal", wapi.Cabinet.BinOperate)

}
