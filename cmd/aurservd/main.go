// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on ard by liasica, magicrolan@qq.com.

package main

import (
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/cmd/aurservd/internal"
    "github.com/auroraride/aurservd/cmd/aurservd/internal/script"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/boot"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
)

func main() {
    snag.WithPanicStack(func() {
        boot.Bootstrap()

        // 初始化缓存
        createCache()

        // 初始化系统设置
        service.NewSetting().Initialize()

        // Demo
        internal.Demo()

        // 启动脚本
        script.Execute()
    })
}

func createCache() {
    cfg := ar.Config.Database.Redis
    cache.CreateClient(cfg.Addr, cfg.Password, cfg.DB)
}
