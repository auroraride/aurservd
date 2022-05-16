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

func Run(start bool) {
    // yd := NewYundong()
    kx := NewKaixin()
    if start {
        StartCabinetProvider(kx)
    }
}

func StartCabinetProvider(providers ...Provider) {
    for _, p := range providers {
        provider := p
        go func() {
            for {
                provider.PrepareRequest()
                provider.UpdateStatus()
                time.Sleep(1 * time.Minute)
            }
        }()
    }
}
