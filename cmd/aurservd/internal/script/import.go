// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-27
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
    Use:   "import",
    Short: "导入助手",
}

func importRiderCmd() *cobra.Command {
    var (
        path string
    )

    cmd := &cobra.Command{
        Use:   "rider",
        Short: "从csv中导入骑手",
        Run: func(cmd *cobra.Command, args []string) {
            service.NewImportRider().ParseCSV(path)
        },
    }

    cmd.Flags().StringVarP(&path, "path", "p", "", "文件路径")
    _ = cmd.MarkFlagRequired("path")
    return cmd
}

func init() {
    importCmd.AddCommand(importRiderCmd())
}
