// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-15
// Based on aurservd by liasica, magicrolan@qq.com.

package provider

import (
    "time"
)

type Provider interface {
    PrepareRequest()
    UpdateStatus()
}

func Run() {
    Yundong = NewYundong()
}

func StartCabinetProvider(providers ...Provider) {
    for _, p := range providers {
        provider := p
        ticker := time.NewTicker(time.Minute * 1)
        go func() {
            for {
                select {
                case <-ticker.C:
                    provider.PrepareRequest()
                    provider.UpdateStatus()
                }
            }
        }()
    }
}
