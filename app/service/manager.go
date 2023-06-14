// Copyright (C) liasica. 2022-present.
//
// Created at 2022-02-28
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/permission"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
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
		orm:            ent.Database.Manager,
		ctx:            context.Background(),
	}
}

func NewManagerWithModifier(m *model.Modifier) *managerService {
	s := NewManager()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

// Create 新增管理员
func (s *managerService) Create(req *model.ManagerCreateReq) error {
	if exists, _ := s.orm.QueryNotDeleted().Where(manager.Phone(req.Phone)).Exist(s.ctx); exists {
		return errors.New("用户已存在")
	}
	password, _ := utils.PasswordGenerate(req.Password)
	return s.orm.Create().SetName(req.Name).SetPhone(req.Phone).SetPassword(password).SetRoleID(req.RoleID).Exec(s.ctx)
}

func (s *managerService) Modify(req *model.ManagerModifyReq) {
	m := s.Query(req.ID)
	u := m.Update()
	if req.Phone != "" {
		u.SetPhone(req.Phone)
	}
	if req.Name != "" {
		u.SetName(req.Name)
	}
	if req.RoleID != 0 {
		u.SetRoleID(req.RoleID)
	}
	if req.Password != "" {
		password, _ := utils.PasswordGenerate(req.Password)
		u.SetPassword(password)
	}

	err := u.Exec(s.ctx)
	if err != nil {
		snag.Panic("管理员编辑失败")
		return
	}
}

// Signin 管理员登录
func (s *managerService) Signin(req *model.ManagerSigninReq) (res *model.ManagerSigninRes, err error) {
	var u *ent.Manager
	u, err = s.orm.QueryNotDeleted().Where(manager.Phone(req.Phone)).WithRole().First(s.ctx)
	if err != nil {
		return nil, errors.New(ar.UserNotFound)
	}

	// 比对密码
	if !utils.PasswordCompare(req.Password, u.Password) {
		return nil, errors.New(ar.UserAuthenticationFailed)
	}

	token := xid.New().String() + "/" + utils.NewEcdsaToken()
	key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, u.ID)

	// 删除旧的token
	if old := cache.Get(s.ctx, key).Val(); old != "" {
		cache.Del(s.ctx, key)
		cache.Del(s.ctx, old)
	}

	// 设置登录token，更新最后登录时间
	s.ExtendTokenTime(u.ID, token)

	perms, super := s.GetPermissions(u)

	return &model.ManagerSigninRes{
		ID:          u.ID,
		Token:       token,
		Name:        u.Name,
		Phone:       u.Phone,
		Permissions: perms,
		Super:       super,
	}, err
}

func (s *managerService) GetPermissions(u *ent.Manager) (perms []string, super bool) {
	r := u.Edges.Role
	if r != nil {
		if r.Super {
			return permission.Keys, r.Super
		} else {
			return r.Permissions, r.Super
		}
	}
	return make([]string, 0), false
}

// GetManagerById 根据ID获取管理员
func (s *managerService) GetManagerById(id uint64) (u *ent.Manager, err error) {
	return s.orm.
		QueryNotDeleted().
		WithRole().
		Where(manager.ID(id)).
		First(context.Background())
}

// ExtendTokenTime 延长管理员登录有效期「24小时」
func (s *managerService) ExtendTokenTime(id uint64, token string) {
	key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, id)
	ctx := context.Background()
	cache.Set(ctx, key, token, 24*time.Hour)
	cache.Set(ctx, token, id, 24*time.Hour)
	_, _ = s.orm.
		UpdateOneID(id).
		SetLastSigninAt(time.Now()).
		Save(ctx)
}

// List 列举管理员
func (s *managerService) List(req *model.ManagerListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Order(ent.Desc(manager.FieldCreatedAt)).WithRole()
	if req.Keyword != nil {
		q.Where(
			manager.Or(
				manager.PhoneContainsFold(*req.Keyword),
				manager.NameContainsFold(*req.Keyword),
			),
		)
	}
	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Manager) model.ManagerListRes {
			res := model.ManagerListRes{
				ID:    item.ID,
				Name:  item.Name,
				Phone: item.Phone,
				Role: model.Role{
					ID:   1,
					Name: "无角色",
				},
			}
			r := item.Edges.Role
			if r != nil {
				res.Role = model.Role{
					ID:   r.ID,
					Name: r.Name,
				}
			}
			return res
		},
	)
}

func (s *managerService) Query(id uint64) *ent.Manager {
	mgr, _ := s.GetManagerById(id)
	if mgr == nil {
		snag.Panic("未找到有效的管理员")
	}
	return mgr
}

// Delete 删除管理员
func (s *managerService) Delete(req *model.IDParamReq) {
	s.orm.SoftDeleteOne(s.Query(req.ID)).SaveX(s.ctx)
}
