// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-06
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/robfig/cron/v3"
    log "github.com/sirupsen/logrus"
)

type enterpriseTask struct {
}

func NewEnterprise() *enterpriseTask {
    return &enterpriseTask{}
}

func (t *enterpriseTask) Start() {
    go t.Do()

    c := cron.New()
    entryID, err := c.AddFunc("1 0 * * *", func() {
        log.Info("开始执行 [enterprise] 定时任务")
        t.Do()
    })
    if err != nil {
        log.Fatal(err)
        return
    }
    c.Start()
    log.Infof("[ENTERPRISE TASK] started: %d", entryID)
}

// Do 检查逾期状态
func (*enterpriseTask) Do() {
    srv := service.NewEnterprise()
    // 获取当前所有生效的订阅
    items := srv.QueryAllCollaborated()
    for _, item := range items {
        go srv.UpdateStatement(item)
    }
}