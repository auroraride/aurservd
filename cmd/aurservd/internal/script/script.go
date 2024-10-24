// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
	"github.com/spf13/cobra"

	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/boot"
	"github.com/auroraride/aurservd/pkg/cache"
)

func Execute() {
	var (
		configFile string
	)

	rootCmd := &cobra.Command{
		Use:               "aurservd",
		Short:             "极光出行管理端控制台",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			ar.SetConfigFile(configFile)

			boot.Bootstrap()

			// 初始化缓存
			cache.CreateClient(ar.Redis)

			// 初始化系统设置
			service.NewSetting().Initialize()

			// 初始化营销设置
			service.NewPromotionSettingService().Initialize()

			// 初始化 买前必读 优惠券 积分 阅读相关配置
			biz.NewInstructions().InitInstructions()

			// 初始化资产属性
			service.NewAssetAttributes().Initialize()

			// todo 临时测试用 之后删除
			// fix.AssetDo()

		},
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "配置文件")

	rootCmd.AddCommand(
		cityCmd,
		managerCmd,
		importCmd,
		serverCommand(),
		customerCommand(),
		fixCommand(),
		personCommand(),
		keyCommand(),
	)

	_ = rootCmd.Execute()
}
