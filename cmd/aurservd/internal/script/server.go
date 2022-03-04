// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    "github.com/auroraride/aurservd/app/router"
    "github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
    Use:   "server",
    Short: "启动API服务",
    Run: func(cmd *cobra.Command, args []string) {
        // 启动服务器
        router.Run()
    },
}
