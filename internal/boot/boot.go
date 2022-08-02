// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package boot

import (
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/mgo"
    "github.com/auroraride/aurservd/internal/payment"
    "github.com/auroraride/aurservd/pkg/logger"
    "github.com/golang-module/carbon/v2"
    "os"
    "time"

    _ "github.com/auroraride/aurservd/app/permission"
)

func Bootstrap() {
    // 设置全局时区
    tz := "Asia/Shanghai"
    _ = os.Setenv("TZ", tz)
    loc, _ := time.LoadLocation(tz)
    time.Local = loc
    carbon.SetTimezone(tz)

    // 载入配置
    ar.LoadConfig()

    // 配置日志
    l := ar.Config.Logging
    logger.LoadWithConfig(logger.Config{
        Color:    l.Color,
        Level:    l.Level,
        Age:      l.Age,
        Json:     l.Json,
        RootPath: l.RootPath,
    })

    // 加载数据库
    ent.Database = ent.OpenDatabase(ar.Config.Database.Postgres.Dsn, ar.Config.App.SQL)
    mgo.Connect(ar.Config.Database.Mongo.Url, ar.Config.Database.Mongo.DB)

    // 加载其他数据

    // 初始化日志
    logging.Boot()

    // 初始化支付
    payment.Boot()

    // 加载模板
    assets.LoadTemplates()
}
