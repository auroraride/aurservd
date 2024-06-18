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
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type simTask struct {
}

func NewSimTask() *simTask {
	return &simTask{}
}

func (t *simTask) Start() {
	if !ar.Config.Task.Sim {
		return
	}

	c := cron.New()
	_, err := c.AddFunc("0 9 * * *", func() {
		zap.L().Info("开始执行 @daily[sim] 定时任务")
		t.Do()
	})
	if err != nil {
		zap.L().Fatal("@daily[sim] 定时任务执行失败", zap.Error(err))
		return
	}
	c.Start()
}

// Do 检查SIM卡过期
func (*simTask) Do() {
	items, _ := ent.Database.Cabinet.QueryNotDeleted().Where(
		cabinet.SimDateGTE(carbon.Now().StartOfDay().AddDays(3).StdTime()),
		cabinet.SimDateLTE(carbon.Now().EndOfDay().AddDays(3).StdTime()),
	).WithCity().All(context.Background())
	for _, item := range items {
		data := model.CabinetSimNotice{
			Serial: item.Serial,
			Name:   item.Name,
			Sim:    item.SimSn,
			End:    item.SimDate.Format(carbon.DateLayout),
		}
		c := item.Edges.City
		if c != nil {
			data.City = c.Name
		}
		go func() {
			workwx.New().SendSimExpires(data)
		}()
	}
}
