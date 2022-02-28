// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on ard by liasica, magicrolan@qq.com.

package main

import (
    _ "github.com/auroraride/aurservd/internal/boot"
)

import (
    "github.com/auroraride/aurservd/app/router"
    "github.com/auroraride/aurservd/cmd/aurservd/internal"
)

func main() {
    internal.Demo()
    // 启动服务器
    router.Run()
}
