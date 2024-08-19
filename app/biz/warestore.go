// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package biz

import (
	"context"
	"strconv"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

type warestoreBiz struct {
	ctx           context.Context
	tokenCacheKey string
}

func NewWarestore() *warestoreBiz {
	return &warestoreBiz{
		ctx:           context.Background(),
		tokenCacheKey: ar.Config.Environment.UpperString() + ":" + "WARESTORE:TOKEN",
	}
}

// 仓库登录, 返回仓库端资料
func (b *warestoreBiz) signin(am *ent.AssetManager) *definition.WarehousePeopleSigninRes {
	idstr := strconv.FormatUint(am.ID, 10)
	// 查询并删除旧token key
	exists := ar.Redis.HGet(b.ctx, b.tokenCacheKey, idstr).Val()
	if exists != "" {
		ar.Redis.HDel(b.ctx, b.tokenCacheKey, exists)
	}

	// 生成token
	token := utils.NewEcdsaToken()

	// 存储登录token和ID进行对应
	ar.Redis.HSet(b.ctx, b.tokenCacheKey, token, am.ID)
	ar.Redis.HSet(b.ctx, b.tokenCacheKey, idstr, token)

	return &definition.WarehousePeopleSigninRes{
		Profile: b.Profile(am),
		Token:   token,
	}
}

// Signin 登录
func (b *warestoreBiz) Signin(req *definition.WarehousePeopleSigninReq) *definition.WarehousePeopleSigninRes {
	switch req.SigninType {
	case model.SigninTypeSms:
		// 校验短信
		service.NewSms().VerifyCodeX(req.Phone, req.SmsId, req.Code)
	case model.SigninTypeAuth:
		// 获取手机号
		req.Phone = service.NewminiProgram().GetPhoneNumber(req.Code)
	}
	am, err := ent.Database.AssetManager.QueryNotDeleted().Where(assetmanager.Phone(req.Phone)).First(b.ctx)
	if err != nil {
		snag.Panic("账号不存在")
	}
	return b.signin(am)
}

// Profile 仓管资料
func (b *warestoreBiz) Profile(am *ent.AssetManager) definition.WarehousePeopleProfile {
	return definition.WarehousePeopleProfile{
		ID:    am.ID,
		Phone: am.Phone,
		Name:  am.Name,
	}
}

// TokenVerify Token校验
func (b *warestoreBiz) TokenVerify(token string) (m *ent.AssetManager) {
	// 获取token对应ID
	id, _ := ar.Redis.HGet(b.ctx, b.tokenCacheKey, token).Uint64()
	if id <= 0 {
		return
	}

	// 反向校验token是否正确
	if ar.Redis.HGet(b.ctx, b.tokenCacheKey, strconv.FormatUint(id, 10)).Val() != token {
		return
	}

	// 获取库管人员
	m, _ = ent.Database.AssetManager.QueryNotDeleted().Where(assetmanager.ID(id)).First(b.ctx)
	if m == nil {
		return
	}
	return
}
