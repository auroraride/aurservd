// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-22
// Based on aurservd by liasica, magicrolan@qq.com.

package fix

import (
	"context"
	"fmt"

	"github.com/auroraride/aurservd/internal/baidu"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assistance"
	"github.com/spf13/cobra"
)

func Assistance() *cobra.Command {
	return &cobra.Command{
		Use:   "assistance",
		Short: "修复救援",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			items, _ := ent.Database.Assistance.QueryNotDeleted().WithStore().Where(assistance.HasStore(), assistance.NaviPolylinesIsNil()).All(ctx)
			bm := baidu.NewMap()
			for _, item := range items {
				sto := item.Edges.Store
				seconds, distance, polylines := bm.RidingPlanX(fmt.Sprintf("%f,%f", sto.Lat, sto.Lng), fmt.Sprintf("%f,%f", item.Lat, item.Lng))
				_, _ = item.Update().SetNaviDuration(seconds).
					SetDistance(float64(distance)).
					SetNaviPolylines(polylines).
					Save(ctx)
			}
		},
	}
}
