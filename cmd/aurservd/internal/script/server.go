// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    pvd "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/app/router"
    "github.com/spf13/cobra"
)

func serverCommand() *cobra.Command {
    var (
        provider bool
    )

    cmd := &cobra.Command{
        Use:   "server",
        Short: "启动API服务",
        Run: func(cmd *cobra.Command, args []string) {
            pvd.Run(provider)
            // 启动服务器
            router.Run()
        },
    }

    cmd.Flags().BoolVarP(&provider, "provider", "p", false, "启动电柜状态轮询")

    return cmd
}
