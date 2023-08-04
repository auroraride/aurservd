// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
	"context"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/auroraride/aurservd/app/ec"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/router"
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/app/sync"
	"github.com/auroraride/aurservd/app/task"
	"github.com/auroraride/aurservd/app/task/reminder"
	"github.com/auroraride/aurservd/cmd/aurservd/internal"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/exchange"
)

func serverCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "server",
		Short: "启动API服务",
		Run: func(cmd *cobra.Command, args []string) {
			// 初始化数据
			service.DatabaseInitial()

			// 启动催费服务
			go reminder.Run()

			// 启动 subscribe task
			go task.NewSubscribe().Start()

			// 启动 enterprise task
			go task.NewEnterprise().Start()

			// 启动 sim task
			go task.NewSimTask().Start()

			// 启动 branch task
			go task.NewBranchTask().Start()

			// 启动 reserve task
			go task.NewReserve().Start()

			// 启动 earnings task
			go task.NewPromotionEarnings().Start()

			// 启动电柜任务
			go ec.Start()

			// 启动 任务补偿
			go compensate()

			// 启动sync
			go sync.Run()

			// 启动rpc服务端
			go rpc.Run()

			// 启动服务器
			go router.Run()

			// 判定是否有维护任务
			go service.NewMaintain().UpdateMaintain()

			// 缓存所有电柜名称
			go service.NewCabinet().CacheAll()

			// Demo
			go internal.Demo()

			<-ar.Quit
		},
	}

	return cmd
}

// compensate 程序启动时自动将之前的所有任务标记为失败
func compensate() {
	now := time.Now()
	msg := "程序异常"
	m := ec.DeleteRange(func(x *ec.Task) bool {
		return x.Status != model.TaskStatusProcessing
	})

	orm := ent.Database.Exchange
	ctx := context.Background()
	items, _ := orm.QueryNotDeleted().Where(exchange.Success(false), exchange.FinishAtIsNil(), exchange.CabinetIDNotNil(), exchange.StartAtNotNil()).All(ctx)
	log.Printf("共获取到%d个进行中的换电", len(items))
	for _, item := range items {
		u := item.Update().
			SetSuccess(false).
			SetFinishAt(now).
			SetDuration(int(now.Sub(item.CreatedAt).Seconds())).
			SetMessage(msg)
		if len(item.Steps) > 0 {
			last := len(item.Steps) - 1
			item.Steps[last].Time = now
			item.Steps[last].Status = model.TaskStatusFail
			u.SetSteps(item.Steps)
		}
		if t, ok := m[item.UUID]; ok {
			u.SetEmpty(t.Exchange.Empty).SetFully(t.Exchange.Fully)
			t.Delete()
		}
		_, _ = u.Save(ctx)
	}
	ec.Clear()
}
