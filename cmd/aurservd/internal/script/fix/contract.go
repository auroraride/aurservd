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
    "github.com/auroraride/aurservd/pkg/snag"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "time"
)

func Contract() *cobra.Command {
    return &cobra.Command{
        Use:   "contract",
        Short: "修复签约时间",
        Run: func(cmd *cobra.Command, args []string) {
            ctx := context.Background()
            items, _ := ent.Database.Contract.QueryNotDeleted().WithRider().All(ctx)
            log.Printf("即将查询 %d 条结果", len(items))
            s := esign.New()
            for _, item := range items {
                err := snag.WithPanic(func() {

                    log.Printf("查询 %d 签约时间", item.ID)
                    status, sr := s.Result(item.FlowID)

                    updater := item.Update()

                    // 设置骑手
                    r := item.Edges.Rider
                    if r != nil {
                        updater.SetRiderInfo(&model.ContractRider{
                            Phone:        r.Phone,
                            Name:         r.Name,
                            IDCardNumber: r.IDCardNumber,
                        })
                    }

                    // 修复状态
                    if status == model.ContractStatusSuccess {
                        updater.SetSignedAt(sr.EndAt()).SetStatus(model.ContractStatusSuccess.Value())
                        // 文档修复, 查询并下载文档
                        if len(item.Files) == 0 {
                            updater.SetFiles(s.DownloadDocument(fmt.Sprintf("%s-%s/contracts/", r.Name, r.IDCardNumber), item.FlowID))
                        }
                    } else {
                        // 设置过期时间
                        if at, ok := sr.SignValidity.(float64); ok {
                            t := time.UnixMilli(int64(at))
                            updater.SetExpiresAt(t)
                        }
                        if at, ok := sr.ContractValidity.(float64); ok {
                            t := time.UnixMilli(int64(at))
                            updater.SetExpiresAt(t)
                        }
                    }

                    _ = updater.Exec(ctx)
                })

                if err != nil {
                    log.Error(err)
                }

                time.Sleep(1 * time.Second)
            }
        },
    }
}
