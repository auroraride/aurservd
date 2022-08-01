// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-01
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import "github.com/auroraride/aurservd/app/controller/tool"

func loadToolRoutes() {
    g := root.Group("tools")
    g.GET("/cabinet/exchange/wR3l9ozNbxvmE8597eHqt0tdeoLiSdwl", tool.Cabinet.Exchange)
}
