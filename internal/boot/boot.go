// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package boot

import (
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/pkg/logger"
)

func init() {
    // 载入配置
    ar.LoadConfig()

    // 配置日志
    l := ar.Config.Logging
    logger.LoadWithConfig(logger.Config{
        Color: l.Color,
        Level: l.Level,
        Age:   l.Age,
        Json:  l.Json,
    })
}
