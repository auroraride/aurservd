// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import "github.com/auroraride/aurservd/internal/ent"

// RiderContext 骑手上下文
type RiderContext struct {
    *Context

    Rider *ent.Rider
}
