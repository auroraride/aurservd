// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-06
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

import (
	"time"

	"github.com/auroraride/adapter/async"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/labstack/echo/v4"
)

type maintain struct{}

var Maintain = new(maintain)

func (*maintain) Update(echo.Context) (err error) {
	// 标记为维护中
	service.NewMaintain().CreateMaintainFile()

	// 查询任务
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		// 是否有进行中的异步业务
		if async.IsDone() {
			ar.Quit <- true
			return
		}
	}

	return
}
