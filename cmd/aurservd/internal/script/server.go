// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    pvd "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/app/router"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/app/task"
    "github.com/spf13/cobra"
)

func serverCommand() *cobra.Command {

    cmd := &cobra.Command{
        Use:   "server",
        Short: "启动API服务",
        Run: func(cmd *cobra.Command, args []string) {
            // 初始化数据
            service.DatabaseInitial()

            // 启动电柜服务
            go pvd.Run()

            // 启动subscribe task
            go task.NewSubscribe().Start()

            // 启动enterprise task
            go task.NewEnterprise().Start()

            // 启动服务器
            router.Run()
        },
    }

    return cmd
}
