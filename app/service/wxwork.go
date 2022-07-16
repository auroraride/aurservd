// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/workwx"
)

type wxworkService struct {
    ctx        context.Context
    agentID    int64
    corpID     string
    corpSecret string
    client     *workwx.Client
}

func NewWxwork() *wxworkService {
    return &wxworkService{
        ctx:    context.Background(),
        client: workwx.New(),
    }
}
