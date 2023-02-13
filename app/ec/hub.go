// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-13
// Based on aurservd by liasica, magicrolan@qq.com.

package ec

import (
    "context"
    "github.com/auroraride/aurservd/internal/ar"
    "time"
)

var hub = &Hub{
    tasks:   make(map[string]*time.Timer),
    setter:  make(chan string),
    updater: make(chan string),
}

type Hub struct {
    tasks   map[string]*time.Timer
    setter  chan string
    updater chan string
}

func Start() {
    for {
        select {
        case id := <-hub.setter:
            hub.tasks[id] = time.AfterFunc(DeactivateTime*time.Second, func() {
                ar.Redis.HDel(context.Background(), ar.TaskCacheKey, id)
            })
        case id := <-hub.updater:
            if t, ok := hub.tasks[id]; ok {
                t.Reset(10 * time.Minute)
            }
        }
    }
}
