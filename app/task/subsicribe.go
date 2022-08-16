// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/robfig/cron/v3"
    log "github.com/sirupsen/logrus"
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
    entryID, err := c.AddFunc("@daily", func() {
        log.Info("开始执行 @daily[subscribe] 定时任务")
        go t.Do()
    })
    if err != nil {
        log.Fatal(err)
        return
    }
    c.Start()
    log.Infof("[SUBSCRIBE TASK] started: %d", entryID)
}

// Do 检查逾期状态
func (*subscribeTask) Do() {
    srv := service.NewSubscribe()
    // 获取当前所有生效的订阅
    items := srv.QueryAllRidersEffective()
    for _, item := range items {
        _ = srv.UpdateStatus(item, true)
    }
}
