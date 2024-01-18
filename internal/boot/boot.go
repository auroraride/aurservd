// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package boot

import (
	"os"
	"time"

	"github.com/auroraride/adapter/log"
	"github.com/golang-module/carbon/v2"
	"github.com/redis/go-redis/v9"

	"github.com/auroraride/aurservd/app/logging"
	_ "github.com/auroraride/aurservd/app/permission"
	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/payment"
	"github.com/auroraride/aurservd/internal/tencent"
)

func Bootstrap() {
	// 设置全局时区
	tz := "Asia/Shanghai"
	_ = os.Setenv("TZ", tz)
	loc, _ := time.LoadLocation(tz)
	ar.TimeLocation = loc
	time.Local = loc
	carbon.SetTimezone(tz)

	// 载入配置
	ar.LoadConfig()

	// 创建redis客户端
	ar.Redis = redis.NewClient(&redis.Options{
		Addr:     ar.Config.Redis.Address,
		Password: ar.Config.Redis.Password,
		DB:       ar.Config.Redis.DB,
	})

	// 初始化日志
	logcfg := &log.Config{
		FormatJson: !ar.Config.LoggerDebug,
		Stdout:     ar.Config.LoggerDebug,
		LoggerName: ar.Config.LoggerName,
		NoCaller:   true,
	}
	if !ar.Config.LoggerDebug {
		logcfg.Writers = append(logcfg.Writers, log.NewRedisWriter(ar.Redis))
	}
	log.New(logcfg)

	// 加载数据库
	entInit()

	// 初始化日志
	logging.Boot()

	// 初始化支付
	payment.Boot()

	// 加载模板
	assets.LoadTemplates()

	// 初始化腾讯人身核验
	tencent.BootWbFace(ar.Redis)
}
