// Copyright (C) liasica. 2022-present.
//
// Created at 2022-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/manager"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/rs/xid"
    log "github.com/sirupsen/logrus"
    "time"
)

type managerService struct {
    cacheKeyPrefix string

    orm      *ent.ManagerClient
    ctx      context.Context
    modifier *model.Modifier
}

func NewManager() *managerService {
    return &managerService{
        cacheKeyPrefix: "MANAGER_",
        orm:            ar.Ent.Manager,
        ctx:            context.Background(),
    }
}

func NewManagerWithModifier(m *model.Modifier) *managerService {
    s := NewManager()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// Create 新增管理员
func (s *managerService) Create(req *model.ManagerCreateReq) error {
    if exists, _ := s.orm.QueryNotDeleted().Where(manager.Phone(req.Phone)).Exist(s.ctx); exists {
        return errors.New("用户已存在")
    }
    password, _ := utils.PasswordGenerate(req.Password)
    return s.orm.Create().SetName(req.Name).SetPhone(req.Phone).SetPassword(password).Exec(s.ctx)
}

// Signin 管理员登录
func (s *managerService) Signin(req *model.ManagerSigninReq) (res *model.ManagerSigninRes, err error) {
    var u *ent.Manager
    u, err = s.orm.QueryNotDeleted().Where(manager.Phone(req.Phone)).Only(s.ctx)
    if err != nil {
        log.Errorf("[M] 管理员查询失败: %v", err)
        return nil, errors.New(ar.UserNotFound)
    }

    // 比对密码
    if !utils.PasswordCompare(req.Password, u.Password) {
        return nil, errors.New(ar.UserAuthenticationFailed)
    }

    token := xid.New().String() + utils.RandTokenString()
    key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, u.ID)

    // 删除旧的token
    if old := cache.Get(s.ctx, key).Val(); old != "" {
        cache.Del(s.ctx, key)
        cache.Del(s.ctx, old)
    }

    // 设置登录token，更新最后登录时间
    s.ExtendTokenTime(u.ID, token)

    return &model.ManagerSigninRes{
        ID:    u.ID,
        Token: token,
        Name:  u.Name,
        Phone: u.Phone,
    }, err
}

// GetManagerById 根据ID获取管理员
func (s *managerService) GetManagerById(id uint64) (u *ent.Manager, err error) {
    return s.orm.
        QueryNotDeleted().
        Where(manager.ID(id)).
        Only(context.Background())
}

// ExtendTokenTime 延长管理员登录有效期「24小时」
func (s *managerService) ExtendTokenTime(id uint64, token string) {
    key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, id)
    ctx := context.Background()
    cache.Set(ctx, key, token, 24*time.Hour)
    cache.Set(ctx, token, id, 24*time.Hour)
    _ = ar.Ent.Rider.
        UpdateOneID(id).
        SetLastSigninAt(time.Now()).
        Exec(context.Background())
}
