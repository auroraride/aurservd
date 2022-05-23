// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-23
// Based on aurservd by liasica, magicrolan@qq.com.

package amap

import "github.com/auroraride/aurservd/internal/ar"

type amap struct {
    Key string
}

type LngLat struct {
    Lng float64 `json:"lng"`
    Lat float64 `json:"lat"`
}

func New() *amap {
    cfg := ar.Config.Amap
    return &amap{Key: cfg.Key}
}
