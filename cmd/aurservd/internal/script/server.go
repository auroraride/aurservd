// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package script

import (
    "context"
    "github.com/auroraride/aurservd/app/ec"
    pvd "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/app/router"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/app/task"
    "github.com/auroraride/aurservd/app/task/reminder"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "time"
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

            // 启动电柜服务
            go pvd.Run()

            // 启动 subscribe task
            go task.NewSubscribe().Start()

            // 启动 enterprise task
            go task.NewEnterprise().Start()

            // 启动 cabinet task
            go task.NewCabinetTask().Start()

            // 启动 branch task
            go task.NewBranchTask().Start()

            // 启动 reserve task
            go task.NewReserve().Start()

            // 启动 任务补偿
            compensate()

            // 启动服务器
            router.Run()
        },
    }

    return cmd
}

// compensate 程序启动时自动将之前的所有任务标记为失败
func compensate() {
    now := time.Now()
    msg := "程序异常"
    tasks := ec.GetAllProcessing()
    log.Infof("共获取到%d个进行中的任务日志", len(tasks))
    m := make(map[string]*ec.Task)
    for _, t := range tasks {
        t.Message = msg
        if t.Job == ec.JobExchange {
            t.Exchange.CurrentStep().Time = now
            t.Exchange.CurrentStep().Status = ec.TaskStatusFail
        }
        t.Stop(ec.TaskStatusFail)
        m[t.ID.Hex()] = t
    }

    orm := ent.Database.Exchange
    ctx := context.Background()
    items, _ := orm.QueryNotDeleted().Where(exchange.Success(false), exchange.FinishAtIsNil(), exchange.CabinetIDNotNil(), exchange.StartAtNotNil()).All(ctx)
    log.Infof("共获取到%d个进行中的换电", len(items))
    for _, item := range items {
        u := item.Update().
            SetSuccess(false).
            SetFinishAt(now).
            SetDuration(int(now.Sub(item.CreatedAt).Seconds()))
        if x, ok := m[item.UUID]; ok {
            u.SetInfo(&ec.ExchangeInfo{
                Cabinet:  x.Cabinet,
                Exchange: x.Exchange,
                Message:  x.Message,
            })
        }
        _, _ = u.Save(ctx)
    }
}
