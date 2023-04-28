// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-28
// Based on aurservd by liasica, magicrolan@qq.com.

package fix

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/esign"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/spf13/cobra"
)

func contractRider() *cobra.Command {
	return &cobra.Command{
		Use:   "common",
		Short: "修复签约骑手和过期时间",
		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()
			items, _ := ent.Database.Contract.QueryNotDeleted().WithRider().All(ctx)
			for _, item := range items {
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
				if item.Status != model.ContractStatusSuccess.Value() && item.CreatedAt.Before(time.Now().Add(-31*time.Minute)) {
					updater.SetExpiresAt(item.CreatedAt.Add(model.ContractExpiration * time.Minute))
				}
				_ = updater.Exec(ctx)
			}
		},
	}
}

func contractStatus() *cobra.Command {
	var (
		sleep time.Duration
	)
	cmd := &cobra.Command{
		Use:   "status",
		Short: "修复签约状态",
		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()
			items, _ := ent.Database.Contract.QueryNotDeleted().WithRider().Where(
				contract.SignedAtIsNil(),
				contract.ExpiresAtIsNil(),
			).All(ctx)
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
						updater.SetExpiresAt(item.CreatedAt.Add(model.ContractExpiration * time.Minute))
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
					log.Println(err)
				}

				time.Sleep(sleep)
			}
		},
	}
	cmd.Flags().DurationVarP(&sleep, "sleep", "s", 30*time.Second, "休眠时间")
	return cmd
}

func Contract() *cobra.Command {
	c := &cobra.Command{
		Use:   "contract",
		Short: "修复签约时间",
	}
	c.AddCommand(
		contractRider(),
		contractStatus(),
	)
	return c
}
