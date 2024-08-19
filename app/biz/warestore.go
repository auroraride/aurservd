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
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

type warestoreBiz struct {
	ctx                    context.Context
	warehouseTokenCacheKey string
	storeTokenCacheKey     string
}

func NewWarestore() *warestoreBiz {
	return &warestoreBiz{
		ctx:                    context.Background(),
		warehouseTokenCacheKey: ar.Config.Environment.UpperString() + ":" + "WAREHOUSE:TOKEN",
		storeTokenCacheKey:     ar.Config.Environment.UpperString() + ":" + "STORE:TOKEN",
	}
}

// signin 仓库、门店登录
func (b *warestoreBiz) signin(am *ent.AssetManager, ep *ent.Employee, platType definition.PlatType) *definition.WarestorePeopleSigninRes {

	var tokenKey string

	switch platType {
	case definition.PlatTypeWarehouse:
		tokenKey = b.warehouseTokenCacheKey
	case definition.PlatTypeStore:
		tokenKey = b.storeTokenCacheKey
	default:
		return nil
	}

	idstr := strconv.FormatUint(am.ID, 10)
	// 查询并删除旧token key
	exists := ar.Redis.HGet(b.ctx, tokenKey, idstr).Val()
	if exists != "" {
		ar.Redis.HDel(b.ctx, tokenKey, exists)
	}

	// 生成token
	token := utils.NewEcdsaToken()

	// 存储登录token和ID进行对应
	ar.Redis.HSet(b.ctx, tokenKey, token, am.ID)
	ar.Redis.HSet(b.ctx, tokenKey, idstr, token)

	return &definition.WarestorePeopleSigninRes{
		Profile: b.Profile(am, ep, platType),
		Token:   token,
	}
}

// Signin 登录
func (b *warestoreBiz) Signin(req *definition.WarestorePeopleSigninReq) *definition.WarestorePeopleSigninRes {
	switch req.SigninType {
	case model.SigninTypeSms:
		// 校验短信
		service.NewSms().VerifyCodeX(req.Phone, req.SmsId, req.Code)
	case model.SigninTypeAuth:
		// 获取手机号
		req.Phone = service.NewminiProgram().GetPhoneNumber(req.Code)
	}

	am := new(ent.AssetManager)
	ep := new(ent.Employee)
	var err error
	switch req.PlatType {
	case definition.PlatTypeWarehouse:
		am, err = ent.Database.AssetManager.QueryNotDeleted().
			Where(
				assetmanager.Phone(req.Phone),
				assetmanager.MiniEnable(true),
			).First(b.ctx)
		if am == nil || err != nil {
			snag.Panic("账号不存在")
		}
	case definition.PlatTypeStore:
		ep, err = ent.Database.Employee.QueryNotDeleted().Where(employee.Phone(req.Phone)).First(b.ctx)
		if am == nil || err != nil {
			snag.Panic("账号不存在")
		}
	}

	return b.signin(am, ep, req.PlatType)
}

// Profile 仓管资料
func (b *warestoreBiz) Profile(am *ent.AssetManager, ep *ent.Employee, platType definition.PlatType) definition.WarestorePeopleProfile {
	switch {
	case platType == definition.PlatTypeWarehouse && am != nil:
		return definition.WarestorePeopleProfile{
			ID:    am.ID,
			Phone: am.Phone,
			Name:  am.Name,
		}

	case platType == definition.PlatTypeStore && ep != nil:
		return definition.WarestorePeopleProfile{
			ID:    ep.ID,
			Phone: ep.Phone,
			Name:  ep.Name,
		}
	default:
		return definition.WarestorePeopleProfile{}
	}
}

// TokenVerify Token校验
func (b *warestoreBiz) TokenVerify(token string) (am *ent.AssetManager, ep *ent.Employee) {
	// 获取token对应ID
	amId, _ := ar.Redis.HGet(b.ctx, b.warehouseTokenCacheKey, token).Uint64()
	if amId > 0 {
		// 反向校验token是否正确
		if ar.Redis.HGet(b.ctx, b.warehouseTokenCacheKey, strconv.FormatUint(amId, 10)).Val() != token {
			return
		}
		// 获取库管人员
		am, _ = ent.Database.AssetManager.QueryNotDeleted().Where(assetmanager.ID(amId), assetmanager.MiniEnable(true)).First(b.ctx)
	}

	epId, _ := ar.Redis.HGet(b.ctx, b.storeTokenCacheKey, token).Uint64()
	if amId > 0 {
		// 反向校验token是否正确
		if ar.Redis.HGet(b.ctx, b.storeTokenCacheKey, strconv.FormatUint(amId, 10)).Val() != token {
			return
		}
		// 获取门店人员
		ep, _ = ent.Database.Employee.QueryNotDeleted().Where(employee.ID(epId)).First(b.ctx)

	}
	return
}

// TransferList 调拨记录列表
func (b *warestoreBiz) TransferList(am *ent.AssetManager, ep *ent.Employee, req *definition.TransferListReq) (res *model.PaginationRes, err error) {
	newReq := model.AssetTransferListReq{
		PaginationReq:       req.PaginationReq,
		AssetTransferFilter: req.AssetTransferFilter,
	}

	if am != nil {
		newReq.AssetManagerID = am.ID
	}

	if ep != nil {
		newReq.EmployeeID = ep.ID
	}

	return service.NewAssetTransfer().TransferList(context.Background(), &newReq)
}

// TransferReceive 接收资产
func (b *warestoreBiz) TransferReceive(am *ent.AssetManager, ep *ent.Employee, req *model.AssetTransferReceiveBatchReq) (err error) {
	var md model.Modifier

	if am != nil {
		md = model.Modifier{
			ID:    am.ID,
			Name:  am.Name,
			Phone: am.Phone,
		}
	}
	if ep != nil {
		md = model.Modifier{
			ID:    ep.ID,
			Name:  ep.Name,
			Phone: ep.Phone,
		}
	}

	return service.NewAssetTransfer().TransferReceive(b.ctx, req, &md)
}
