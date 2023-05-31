// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-06
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/role"
	"github.com/auroraride/aurservd/pkg/snag"
)

type roleService struct {
	ctx      context.Context
	modifier *model.Modifier
	orm      *ent.RoleClient
}

func NewRole() *roleService {
	return &roleService{
		ctx: context.Background(),
		orm: ent.Database.Role,
	}
}

func NewRoleWithModifier(m *model.Modifier) *roleService {
	s := NewRole()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func (s *roleService) Query(id uint64) (*ent.Role, error) {
	return s.orm.Query().Where(role.ID(id)).First(s.ctx)
}

func (s *roleService) QueryX(id uint64) *ent.Role {
	r, _ := s.Query(id)
	if r == nil {
		snag.Panic("角色不存在")
	}
	return r
}

func (s *roleService) Create(req *model.RoleCreateReq) model.Role {
	// 查找是否存在
	if e, _ := s.orm.Query().Where(role.Name(req.Name)).Exist(s.ctx); e {
		snag.Panic("角色已存在")
	}

	q := s.orm.Create().SetName(req.Name)
	if len(req.Permissions) > 0 {
		q.SetPermissions(req.Permissions)
	}

	r := q.SaveX(s.ctx)

	perms := make([]string, 0)
	if len(r.Permissions) > 0 {
		perms = r.Permissions
	}

	return model.Role{
		ID:          r.ID,
		Name:        r.Name,
		Permissions: perms,
	}
}

func (s *roleService) Modify(req *model.RoleModifyReq) model.Role {
	r := s.QueryX(req.ID)

	if r.Buildin || r.Super {
		snag.Panic("内置角色无法编辑")
	}

	q := r.Update()
	if req.Name != "" {
		q.SetName(req.Name)
	}
	if req.Permissions != nil && len(req.Permissions) > 0 {
		q.SetPermissions(req.Permissions)
	}
	r = q.SaveX(s.ctx)
	return model.Role{
		ID:          r.ID,
		Name:        r.Name,
		Permissions: r.Permissions,
	}
}

func (s *roleService) List() []model.Role {
	items, _ := s.orm.Query().All(s.ctx)
	res := make([]model.Role, len(items))
	for i, item := range items {
		res[i] = model.Role{
			ID:          item.ID,
			Name:        item.Name,
			Permissions: item.Permissions,
			Builtin:     item.Buildin,
			Super:       item.Super,
		}
	}
	return res
}

func (s *roleService) Delete(req *model.IDParamReq) {
	// 查找是否有用户
	if e, _ := ent.Database.Manager.QueryNotDeleted().Where(manager.RoleID(req.ID)).Exist(s.ctx); e {
		snag.Panic("角色存在用户, 无法删除")
	}
	err := s.orm.DeleteOneID(req.ID).Exec(s.ctx)
	if err != nil {
		snag.Panic("角色删除失败")
	}
}
