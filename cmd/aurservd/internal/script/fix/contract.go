// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-28
// Based on aurservd by liasica, magicrolan@qq.com.

package fix

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/esign"
    "github.com/spf13/cobra"
    "time"
)

func Contract() *cobra.Command {
    return &cobra.Command{
        Use:   "contract",
        Short: "修复签约时间",
        Run: func(cmd *cobra.Command, args []string) {
            ctx := context.Background()
            items, _ := ent.Database.Contract.QueryNotDeleted().All(ctx)
            fmt.Printf("即将查询 %d 条结果\n", len(items))
            for _, item := range items {
                fmt.Printf("查询 %d 签约时间\n", item.ID)
                status, sr := esign.New().Result(item.FlowID)
                at, ok := sr.SignValidity.(float64)

                updater := item.Update()

                // 如果签约完成但是签约完成时间为空
                if item.SignedAt == nil && status == model.ContractStatusSuccess {
                    updater.SetSignedAt(sr.EndAt())
                }

                // 如果过期时间为空
                if item.ExpiresAt == nil && ok {
                    t := time.UnixMilli(int64(at))
                    updater.SetExpiresAt(t)
                }

                _ = updater.Exec(ctx)
                time.Sleep(1 * time.Second)
            }
        },
    }
}
