// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/app/task/reminder"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/robfig/cron/v3"
    "go.uber.org/zap"
)

type subscribeTask struct {
}

func NewSubscribe() *subscribeTask {
    return &subscribeTask{}
}

func (t *subscribeTask) Start() {
    if !ar.Config.Task.Subscribe {
        return
    }

    go t.Do()

    c := cron.New()
    _, err := c.AddFunc("@daily", func() {
        zap.L().Info("开始执行 @daily[subscribe] 定时任务")
        go t.Do()
    })
    if err != nil {
        zap.L().Fatal("@daily[subscribe] 定时任务执行失败", zap.Error(err))
        return
    }
    c.Start()
}

// Do 检查逾期状态
func (*subscribeTask) Do() {
    // 重置催费任务
    reminder.Reset()
    srv := service.NewSubscribe()
    // 获取当前所有生效的订阅
    items := srv.QueryAllRidersEffective()
    for _, item := range items {
        _ = srv.UpdateStatus(item, true)
    }
}
