// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-06
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/robfig/cron/v3"
    "go.uber.org/zap"
)

type enterpriseTask struct {
}

func NewEnterprise() *enterpriseTask {
    return &enterpriseTask{}
}

func (t *enterpriseTask) Start() {
    if !ar.Config.Task.Enterprise {
        return
    }

    go t.Do()

    c := cron.New()
    _, err := c.AddFunc("@daily", func() {
        zap.L().Info("开始执行 @daily[enterprise] 定时任务")
        go t.Do()
    })
    if err != nil {
        zap.L().Fatal("@daily[enterprise] 定时任务执行失败", zap.Error(err))
        return
    }
    c.Start()
}

// Do 更新企业订单
func (*enterpriseTask) Do() {
    srv := service.NewEnterprise()
    // 获取当前合作中的企业
    items := srv.QueryAllCollaborated()
    for _, item := range items {
        srv.UpdateStatement(item)
    }
}
