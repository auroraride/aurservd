// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-18
// Based on aurservd by liasica, magicrolan@qq.com.

package fix

import (
	"context"

	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/subscribereminder"
	"github.com/spf13/cobra"
)

func Reminder() *cobra.Command {
	return &cobra.Command{
		Use:   "reminder",
		Short: "修复催费",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			items, _ := ent.Database.SubscribeReminder.Query().Where(subscribereminder.DaysLT(0)).WithSubscribe().All(ctx)

			for _, item := range items {
				f, l := service.NewSubscribe().OverdueFee(item.Edges.Subscribe)
				_, _ = item.Update().SetFee(f).SetFeeFormula(l).Save(ctx)
			}
		},
	}
}
