// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "github.com/auroraride/aurservd/app/socket"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/cache"
)

type riderSocketService struct {
    ctx context.Context
}

func NewRiderSocket() *riderSocketService {
    return &riderSocketService{
        ctx: context.Background(),
    }
}

func (s *riderSocketService) Prefix() string {
    return "RIDER"
}

func (s *riderSocketService) Connect(hub *socket.WebsocketHub, token string) (uint64, error) {
    id, _ := cache.Get(context.Background(), token).Uint64()
    r, _ := ar.Ent.Rider.QueryNotDeleted().Where(rider.ID(id)).First(s.ctx)
    if r == nil {
        return 0, errors.New("骑手未找到")
    }

    return id, nil
}
