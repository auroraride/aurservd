// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-18
// Based on aurservd by liasica, magicrolan@qq.com.

package fix

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/business"
)

func Commission() *cobra.Command {
	return &cobra.Command{
		Use:   "commission",
		Short: "修复提成",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			items, _ := ent.Database.Commission.QueryNotDeleted().WithOrder(func(query *ent.OrderQuery) {
				query.WithSubscribe()
			}).All(ctx)

			for _, item := range items {
				o := item.Edges.Order
				if o == nil {
					fmt.Printf("%d 未找到订单\n", item.ID)
					continue
				}
				sub := o.Edges.Subscribe
				if sub == nil {
					fmt.Printf("%d [O:%d] 未找到订单\n", item.ID, o.ID)
					continue
				}
				b, _ := ent.Database.Business.QueryNotDeleted().Where(business.SubscribeID(sub.ID), business.TypeEQ(model.BusinessTypeActive)).First(ctx)
				if b == nil {
					fmt.Printf("%d [O:%d, S:%d] 未找到业务\n", item.ID, o.ID, sub.ID)
					continue
				}

				_, _ = item.Update().SetSubscribeID(o.SubscribeID).SetNillablePlanID(o.PlanID).SetSubscribeID(sub.ID).SetBusiness(b).SetRiderID(sub.RiderID).Save(ctx)
			}
		},
	}
}
