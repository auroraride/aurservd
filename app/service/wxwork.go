// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-16
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/workwx"
    "github.com/auroraride/aurservd/internal/ar"
)

type wxworkService struct {
    ctx        context.Context
    agentID    int64
    corpID     string
    corpSecret string
    client     *workwx.Client
}

func NewWxwork() *wxworkService {
    cfg := ar.Config.WxWork
    return &wxworkService{
        ctx:    context.Background(),
        client: workwx.New(cfg.AgentID, cfg.CorpID, cfg.CorpSecret),
    }
}
