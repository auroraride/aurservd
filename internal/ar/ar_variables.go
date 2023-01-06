// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-10
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import "sync"

var (
    Quit chan bool

    AsynchronousTask *sync.Map
)

func init() {
    Quit = make(chan bool)
    AsynchronousTask = &sync.Map{}
}
