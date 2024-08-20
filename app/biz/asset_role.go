// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assetrole"
	"github.com/auroraride/aurservd/pkg/snag"
)

type assetRoleBiz struct {
	ctx      context.Context
	orm      *ent.AssetRoleClient
	modifier *model.Modifier
}

func NewAssetRole() *assetRoleBiz {
	return &assetRoleBiz{
		ctx: context.Background(),
		orm: ent.Database.AssetRole,
	}
}

func NewAssetRoleWithModifier(m *model.Modifier) *assetRoleBiz {
	b := NewAssetRole()
	b.ctx = context.WithValue(b.ctx, model.CtxModifierKey{}, m)
	b.modifier = m
	return b
}

func (b *assetRoleBiz) Query(id uint64) (*ent.AssetRole, error) {
	return b.orm.Query().Where(assetrole.ID(id)).First(b.ctx)
}

func (b *assetRoleBiz) QueryX(id uint64) *ent.AssetRole {
	r, _ := b.Query(id)
	if r == nil {
		snag.Panic("角色不存在")
	}
	return r
}

func (b *assetRoleBiz) Create(req *definition.AssetRoleCreateReq) definition.AssetRole {
	// 查找是否存在
	if e, _ := b.orm.Query().Where(assetrole.Name(req.Name)).Exist(b.ctx); e {
		snag.Panic("角色已存在")
	}

	q := b.orm.Create().SetName(req.Name)
	if len(req.Permissions) > 0 {
		q.SetPermissions(req.Permissions)
	}

	r := q.SaveX(b.ctx)

	perms := make([]string, 0)
	if len(r.Permissions) > 0 {
		perms = r.Permissions
	}

	return definition.AssetRole{
		ID:          r.ID,
		Name:        r.Name,
		Permissions: perms,
	}
}

func (b *assetRoleBiz) Modify(req *definition.AssetRoleModifyReq) definition.AssetRole {
	r := b.QueryX(req.ID)

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
	r = q.SaveX(b.ctx)
	return definition.AssetRole{
		ID:          r.ID,
		Name:        r.Name,
		Permissions: r.Permissions,
	}
}

func (b *assetRoleBiz) List() []definition.AssetRole {
	items, _ := b.orm.Query().All(b.ctx)
	res := make([]definition.AssetRole, len(items))
	for i, item := range items {
		res[i] = definition.AssetRole{
			ID:          item.ID,
			Name:        item.Name,
			Permissions: item.Permissions,
			Builtin:     item.Buildin,
			Super:       item.Super,
		}
	}
	return res
}

func (b *assetRoleBiz) Delete(req *model.IDParamReq) {
	// 查找是否有用户
	if e, _ := ent.Database.AssetManager.QueryNotDeleted().Where(assetmanager.RoleID(req.ID)).Exist(b.ctx); e {
		snag.Panic("角色存在用户, 无法删除")
	}
	err := b.orm.DeleteOneID(req.ID).Exec(b.ctx)
	if err != nil {
		snag.Panic("角色删除失败")
	}
}

// RoleSelection 筛选角色
func (b *assetRoleBiz) RoleSelection() (items []model.SelectOption) {
	roles, _ := b.orm.Query().All(b.ctx)
	for _, role := range roles {
		items = append(items, model.SelectOption{
			Value: role.ID,
			Label: role.Name,
		})
	}
	return
}
