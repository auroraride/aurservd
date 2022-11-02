// Copyright (C) liasica. 2022-present.
//
// Created at 2022-11-02
// Based on aurservd by liasica, magicrolan@qq.com.

package fix

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/spf13/cobra"
)

func Sql() *cobra.Command {
    var (
        sql string
    )

    c := &cobra.Command{
        Use:   "sql",
        Short: "SQL修复",
        Run: func(_ *cobra.Command, _ []string) {
            _, err := ent.Database.ExecContext(context.Background(), sql)
            if err != nil {
                fmt.Printf("执行失败: %v", err)
            }
        },
    }

    c.Flags().StringVar(&sql, "raw", "", "要执行的SQL语句")

    _ = c.MarkFlagRequired("sql")

    return c
}
