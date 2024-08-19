// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	jsoniter "github.com/json-iterator/go"

	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assetrole"
	"github.com/auroraride/aurservd/internal/ent/manager"
	"github.com/auroraride/aurservd/internal/ent/role"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

func DatabaseInitial() {
	cityInitial()
	managerInitial(roleInitial())
	assetManagerInitial(assetRoleInitial())
}

func managerInitial(r *ent.Role) {
	client := ent.Database.Manager
	p := "18888888888"

	if e, _ := client.QueryNotDeleted().
		Where(manager.Phone(p)).
		Exist(context.Background()); e {
		return
	}

	password, _ := utils.PasswordGenerate("AuroraAdmin@2022#!")
	client.Create().
		SetName("超级管理员").
		SetPhone(p).
		SetPassword(password).
		SetRoleID(r.ID).
		ExecX(context.Background())
}

func roleInitial() *ent.Role {
	client := ent.Database.Role
	ctx := context.Background()
	name := "超级管理员"
	if e, _ := client.Query().Where(role.Buildin(true)).First(ctx); e != nil {
		return e
	}
	return client.Create().SetName(name).SetSuper(true).SetBuildin(true).SaveX(ctx)
}

func cityInitial() {
	type R struct {
		Adcode   uint64  `json:"adcode"`
		Name     string  `json:"name"`
		Code     string  `json:"code"`
		Lng      float64 `json:"lng,omitempty"`
		Lat      float64 `json:"lat,omitempty"`
		Children []R     `json:"children,omitempty"`
	}

	ctx := context.Background()

	if c, _ := ent.Database.City.Query().Count(context.Background()); c > 0 {
		return
	}

	// 导入城市
	var items []R
	err := jsoniter.Unmarshal(assets.City, &items)
	if err == nil {
		ent.WithTxPanic(ctx, func(tx *ent.Tx) error {
			for _, item := range items {
				parent, err := tx.City.Create().
					SetID(item.Adcode).
					SetName(item.Name).
					SetCode(item.Code).
					Save(ctx)
				snag.PanicIfError(err)
				for _, child := range item.Children {
					_, err = tx.City.Create().
						SetID(child.Adcode).
						SetName(child.Name).
						SetCode(child.Code).
						SetOpen(false).
						SetParent(parent).
						SetLat(child.Lat).
						SetLng(child.Lng).
						Save(ctx)
					snag.PanicIfError(err)
				}
			}
			return nil
		})
	}
}

func assetRoleInitial() *ent.AssetRole {
	client := ent.Database.AssetRole
	ctx := context.Background()
	name := "超级管理员"
	if e, _ := client.Query().Where(assetrole.Buildin(true)).First(ctx); e != nil {
		return e
	}
	return client.Create().SetName(name).SetSuper(true).SetBuildin(true).SaveX(ctx)
}

func assetManagerInitial(r *ent.AssetRole) {
	client := ent.Database.AssetManager
	p := "18888888888"

	if e, _ := client.QueryNotDeleted().
		Where(assetmanager.Phone(p)).
		Exist(context.Background()); e {
		return
	}

	password, _ := utils.PasswordGenerate("AuroraAdmin@2022#!")
	client.Create().
		SetName("超级管理员").
		SetPhone(p).
		SetPassword(password).
		SetRoleID(r.ID).
		ExecX(context.Background())
}
