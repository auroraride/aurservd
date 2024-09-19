// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-17, by aurb

package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"

	assetpermission "github.com/auroraride/aurservd/app/assetpermission"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assetrole"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/utils"
)

type assetManagerBiz struct {
	cacheKeyPrefix string
	orm            *ent.AssetManagerClient
	ctx            context.Context
	modifier       *model.Modifier
}

func NewAssetManager() *assetManagerBiz {
	return &assetManagerBiz{
		cacheKeyPrefix: "ASSET_MANAGER_",
		orm:            ent.Database.AssetManager,
		ctx:            context.Background(),
	}
}

func NewAssetManagerWithModifier(m *model.Modifier) *assetManagerBiz {
	b := NewAssetManager()
	b.ctx = context.WithValue(b.ctx, model.CtxModifierKey{}, m)
	b.modifier = m
	return b
}

// Create 新增管理员
func (b *assetManagerBiz) Create(req *definition.AssetManagerCreateReq) error {
	if exists, _ := b.orm.QueryNotDeleted().Where(assetmanager.Phone(req.Phone)).Exist(b.ctx); exists {
		return errors.New("用户已存在")
	}
	password, _ := utils.PasswordGenerate(req.Password)
	return b.orm.Create().SetName(req.Name).SetPhone(req.Phone).SetPassword(password).SetRoleID(req.RoleID).Exec(b.ctx)
}

func (b *assetManagerBiz) Modify(req *definition.AssetManagerModifyReq) error {
	u := b.orm.UpdateOneID(req.ID)
	if req.Phone != nil {
		u.SetPhone(*req.Phone)
	}
	if req.Name != nil {
		u.SetName(*req.Name)
	}
	if req.RoleID != nil {
		u.SetRoleID(*req.RoleID)
	}
	if req.Password != nil {
		password, _ := utils.PasswordGenerate(*req.Password)
		u.SetPassword(password)
	}
	if req.MiniEnable != nil {
		u.SetMiniEnable(*req.MiniEnable)
	}
	if req.MiniLimit != nil {
		u.SetMiniLimit(*req.MiniLimit)
	}
	if len(req.WarehouseIDs) != 0 {
		u.ClearBelongWarehouses()
		u.AddBelongWarehouseIDs(req.WarehouseIDs...)
	}

	_, err := u.Save(b.ctx)
	if err != nil {
		return errors.New("管理员编辑失败")
	}
	return nil
}

// Signin 管理员登录
func (b *assetManagerBiz) Signin(req *definition.AssetManagerSigninReq) (res *definition.AssetManagerSigninRes, err error) {
	var u *ent.AssetManager
	u, err = b.orm.QueryNotDeleted().Where(assetmanager.Phone(req.Phone)).WithRole().First(b.ctx)
	if err != nil {
		return nil, errors.New(ar.UserNotFound)
	}

	// 比对密码
	if !utils.PasswordCompare(req.Password, u.Password) {
		return nil, errors.New(ar.UserAuthenticationFailed)
	}

	token := xid.New().String() + "/" + utils.NewEcdsaToken()
	key := fmt.Sprintf("%s%d", b.cacheKeyPrefix, u.ID)

	// 删除旧的token
	if old := cache.Get(b.ctx, key).Val(); old != "" {
		cache.Del(b.ctx, key)
		cache.Del(b.ctx, old)
	}

	// 设置登录token，更新最后登录时间
	b.ExtendTokenTime(u.ID, token)

	perms, super := b.GetAssetPermissions(u)

	return &definition.AssetManagerSigninRes{
		ID:          u.ID,
		Token:       token,
		Name:        u.Name,
		Phone:       u.Phone,
		Permissions: perms,
		Super:       super,
	}, err
}

func (b *assetManagerBiz) GetAssetPermissions(u *ent.AssetManager) (perms []string, super bool) {
	r := u.Edges.Role
	if r != nil {
		if r.Super {
			return assetpermission.Keys, r.Super
		} else {
			return r.Permissions, r.Super
		}
	}
	return make([]string, 0), false
}

// GetAssetManagerById 根据ID获取管理员
func (b *assetManagerBiz) GetAssetManagerById(id uint64) (u *ent.AssetManager, err error) {
	return b.orm.
		QueryNotDeleted().
		WithRole().
		Where(assetmanager.ID(id)).
		First(context.Background())
}

// ExtendTokenTime 延长管理员登录有效期「24小时」
func (b *assetManagerBiz) ExtendTokenTime(id uint64, token string) {
	key := fmt.Sprintf("%s%d", b.cacheKeyPrefix, id)
	ctx := context.Background()
	cache.Set(ctx, key, token, 24*time.Hour)
	cache.Set(ctx, token, id, 24*time.Hour)
	_, _ = b.orm.
		UpdateOneID(id).
		SetLastSigninAt(time.Now()).
		Save(ctx)
}

// List 列举管理员
func (b *assetManagerBiz) List(req *definition.AssetManagerListReq) *model.PaginationRes {
	q := b.orm.QueryNotDeleted().Order(ent.Desc(assetmanager.FieldCreatedAt)).WithRole().WithBelongWarehouses(func(query *ent.WarehouseQuery) {
		query.WithCity()
	})
	if req.Keyword != nil {
		q.Where(
			assetmanager.Or(
				assetmanager.PhoneContainsFold(*req.Keyword),
				assetmanager.NameContainsFold(*req.Keyword),
			),
		)
	}

	if req.Warestore != nil && *req.Warestore {
		q.Where(
			assetmanager.HasRoleWith(assetrole.Name("仓库管理员"), assetrole.Buildin(true)),
		)
	}

	if req.WarehouseID != nil {
		q.Where(
			assetmanager.HasBelongWarehousesWith(
				warehouse.ID(*req.WarehouseID),
			),
		)
	}

	if req.Enable != nil {
		q.Where(
			assetmanager.MiniEnable(*req.Enable),
		)
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.AssetManager) definition.AssetManagerListRes {
			res := definition.AssetManagerListRes{
				ID:    item.ID,
				Name:  item.Name,
				Phone: item.Phone,
				Role: definition.AssetRole{
					ID:   1,
					Name: "无角色",
				},
				MiniEnable: item.MiniEnable,
				MiniLimit:  item.MiniLimit,
			}
			r := item.Edges.Role
			if r != nil {
				res.Role = definition.AssetRole{
					ID:   r.ID,
					Name: r.Name,
				}
			}
			whs := item.Edges.BelongWarehouses
			wRes := make([]*definition.AssetManagerWarehouse, 0)
			for _, wh := range whs {
				var cityName string
				if wh.Edges.City != nil {
					cityName = wh.Edges.City.Name
				}
				wRes = append(wRes, &definition.AssetManagerWarehouse{
					ID:       wh.ID,
					Name:     wh.Name,
					CityName: cityName,
				})

			}
			res.Warehouses = wRes
			return res
		},
	)
}

func (b *assetManagerBiz) Query(id uint64) (*ent.AssetManager, error) {
	mgr, _ := b.GetAssetManagerById(id)
	if mgr == nil {
		return nil, errors.New("未找到有效的管理员")
	}
	return mgr, nil
}

// Delete 删除管理员
func (b *assetManagerBiz) Delete(req *model.IDParamReq) error {
	am, err := b.Query(req.ID)
	if err != nil {
		return err
	}
	b.orm.SoftDeleteOne(am).SaveX(b.ctx)
	return nil
}

// Profile 管理员信息
func (b *assetManagerBiz) Profile(m *ent.AssetManager) definition.AssetManagerSigninRes {
	perms, super := b.GetAssetPermissions(m)
	return definition.AssetManagerSigninRes{
		ID:          m.ID,
		Name:        m.Name,
		Phone:       m.Phone,
		Permissions: perms,
		Super:       super,
	}
}
