// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-16
// Based on aurservd by liasica, magicrolan@qq.com.

package reminder

import (
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/subscribereminder"
)

type Task struct {
    Name              string
    Phone             string
    Type              subscribereminder.Type
    Days              int
    SubscribeID       uint64
    PlanName          string
    SubscribeReminder *ent.SubscribeReminder
    Success           bool
}
