// Copyright (C) liasica. 2022-present.
//
// Created at 2022-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/manager"
    "time"
)

type managerService struct {
    cacheKeyPrefix string
}

func NewManager() *managerService {
    return &managerService{
        cacheKeyPrefix: "MANAGER_",
    }
}

// GetManagerById 根据ID获取管理员
func (*managerService) GetManagerById(id uint64) (u *ent.Manager, err error) {
    return ar.Ent.Manager.
        QueryNotDeleted().
        Where(manager.ID(id)).
        Only(context.Background())
}

// ExtendTokenTime 延长管理员登录有效期「24小时」
func (s *managerService) ExtendTokenTime(id uint64, token string) {
    key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, id)
    cache := ar.Cache
    ctx := context.Background()
    cache.Set(ctx, key, token, 24*time.Hour)
    cache.Set(ctx, token, id, 24*time.Hour)
    _ = ar.Ent.Rider.
        UpdateOneID(id).
        SetLastSigninAt(time.Now()).
        Exec(context.Background())
}
