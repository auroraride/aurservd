// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
	v1 "github.com/auroraride/aurservd/app/router/v1"
	v2 "github.com/auroraride/aurservd/app/router/v2"
)

func loadRiderRoutes() {
	v1.LoadRiderV1Routes(root)
	v2.LoadRiderV2Routes(root)
}
