// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-30
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import "github.com/auroraride/aurservd/internal/ent"

type AgentContext struct {
    *BaseContext

    Enterprise *ent.Enterprise
}
