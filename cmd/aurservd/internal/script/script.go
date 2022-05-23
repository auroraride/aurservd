// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import "github.com/spf13/cobra"

func Execute() {
    rootCmd := &cobra.Command{
        Use:   "aurservd",
        Short: "极光出行管理端控制台",
    }
    rootCmd.AddCommand(
        cityCmd,
        managerCmd,
        serverCommand(),
    )
    _ = rootCmd.Execute()
}
