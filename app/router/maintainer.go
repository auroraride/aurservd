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
	auth.GET("/cabinet/:serial", oapi.Cabinet.Detail)
	auth.POST("/cabinet/:serial", oapi.Cabinet.Operate)
	auth.POST("/cabinet/:serial/:ordinal", oapi.Cabinet.BinOperate)
	auth.POST("/cabinet/pause/:serial", oapi.Cabinet.Pause)

	// 筛选项
	auth.GET("/selection/warehouse", oapi.Selection.Warehouse)      // 城市仓库
	auth.GET("/selection/store", oapi.Selection.Store)              // 城市门店
	auth.GET("/selection/city_station", oapi.Selection.CityStation) // 城市站点
	auth.GET("/selection/maintainer", oapi.Selection.Maintainer)    // 运维
	auth.GET("/selection/model", oapi.Selection.Model)              // 电池型号筛选
	auth.GET("/selection/material", oapi.Selection.Material)        // 物资类型筛选

	// 资产调拨
	auth.POST("/transfer", oapi.AssetTransfer.Transfer)                   // 创建调拨
	auth.GET("/transfer", oapi.AssetTransfer.TransferList)                // 调拨列表
	auth.GET("/transfer/:id", oapi.AssetTransfer.TransferDetail)          // 调拨详情
	auth.POST("/transfer/receive", oapi.AssetTransfer.TransferReceive)    // 调拨批量入库
	auth.GET("/transfer/sn/:sn", oapi.AssetTransfer.TransferBySn)         // 根据sn查询调拨信息
	auth.GET("/transfer/flow", oapi.AssetTransfer.TransferFlow)           // 资产流转明细
	auth.GET("/transfer/details", oapi.AssetTransfer.TransferDetailsList) // 出入库明细
	auth.PUT("/transfer/:id", oapi.AssetTransfer.Modify)                  // 修改调拨记录
	auth.PUT("/transfer/cancel/:id", oapi.AssetTransfer.TransferCancel)   // 取消资产调拨

	// 物资管理
	auth.GET("/assets/common", oapi.Asset.AssetsCommon) // 电池/电车资产数据
	auth.GET("/assets/count", oapi.Asset.AssetCount)    // 资产统计
	auth.GET("/assets", oapi.Asset.Assets)              // 资产数据
	auth.GET("/assets/:sn", oapi.Asset.AssetBySn)       // 通过SN查询资产信息

	auth.GET("/check/sn/:sn", oapi.AssetCheck.GetAssetBySN) // 通过SN查询资产

	// 资产维修
	auth.GET("/maintenance", oapi.AssetMaintenance.List) // 资产维修列表

}
