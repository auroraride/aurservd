// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/robfig/cron/v3"
    log "github.com/sirupsen/logrus"
)

type subscribeTask struct {
}

func NewSubscribe() *subscribeTask {
    return &subscribeTask{}
}

func (t *subscribeTask) Start() {
    go t.Do()
    entryID, err := cron.New().AddFunc("0 0 * * *", func() {
        t.Do()
    })
    if err != nil {
        log.Fatal(err)
        return
    }
    log.Infof("[SUBSCRIBE TASK] started: %d", entryID)
}

// Do 检查逾期状态
func (*subscribeTask) Do() {
    srv := service.NewSubscribe()
    // 获取当前所有生效的订阅
    items := srv.QueryAllEffective()
    for _, item := range items {
        go srv.UpdateStatus(item)
    }
}
