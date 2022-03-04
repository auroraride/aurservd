// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
)

type cscService struct {
    ctx context.Context
}

func NewCSC() *cscService {
    return &cscService{
        ctx: context.Background(),
    }
}

func ParseNameList() {
    // tel := "02863804608"
    // tmpl := "TTS_235791551"
}
