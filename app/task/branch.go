// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-17
// Based on aurservd by liasica, magicrolan@qq.com.

package task

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/workwx"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/branchcontract"
    "github.com/golang-module/carbon/v2"
    "github.com/robfig/cron/v3"
    log "github.com/sirupsen/logrus"
)

type branchTask struct {
}

func NewBranchTask() *branchTask {
    return &branchTask{}
}

func (t *branchTask) Start() {
    if !ar.Config.Task.Branch {
        return
    }

    c := cron.New()
    entryID, err := c.AddFunc("0 9 * * *", func() {
        log.Info("开始执行 @daily[branch] 定时任务")
        t.Do()
    })
    if err != nil {
        log.Fatal(err)
        return
    }
    c.Start()
    log.Infof("[BRANCH TASK] started: %d", entryID)
}

func (*branchTask) Do() {
    items, _ := ent.Database.BranchContract.QueryNotDeleted().Where(
        branchcontract.EndTimeGTE(carbon.Now().StartOfDay().AddDays(3).Carbon2Time()),
        branchcontract.EndTimeLTE(carbon.Now().EndOfDay().AddDays(3).Carbon2Time()),
    ).WithBranch(func(bq *ent.BranchQuery) {
        bq.WithCity()
    }).All(context.Background())
    for _, item := range items {
        data := model.BranchExpriesNotice{
            Name: item.Edges.Branch.Name,
            End:  item.EndTime.Format(carbon.DateLayout),
        }
        c := item.Edges.Branch.Edges.City
        if c != nil {
            data.City = c.Name
        }
        go func() {
            workwx.New().SendBranchExpires(data)
        }()
    }
}
