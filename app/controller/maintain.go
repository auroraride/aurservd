// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-06
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/auroraride/adapter/async"
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
)

type maintain struct{}

var Maintain = new(maintain)

func (*maintain) Update(c echo.Context) (err error) {
	log.Println("已请求维护:", c.Request().RemoteAddr)

	// 标记为维护中
	service.NewMaintain().CreateMaintainFile()

	// 查询任务
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		// 是否有进行中的异步业务
		if async.IsDone() {
			ar.Quit <- true
			return c.String(http.StatusOK, ">>> 已设为维护状态 <<<")
		}
	}

	return
}
