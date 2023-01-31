// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-31
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "log"
    "net/http"
    _ "net/http/pprof"
)

func RunPprof() {
    go func() {
        if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
            log.Fatal(err)
        }
    }()
}
