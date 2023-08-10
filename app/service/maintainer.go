// Copyright (C) liasica. 2023-present.
//
// Created at 2023-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

type maintainerService struct {
	*BaseService

	tokenCacheKey string
	orm           *ent.MaintainerClient
}

func NewMaintainer(params ...any) *maintainerService {
	return &maintainerService{
		BaseService:   newService(params...),
		tokenCacheKey: ar.Config.Environment.UpperString() + ":" + "MAINTAINER:TOKEN",
		orm:           ent.Database.Maintainer,
	}
}

func (s *maintainerService) QueryID(id uint64, enable ...bool) (*ent.Maintainer, error) {
	q := s.orm.Query().Where(maintainer.ID(id))
	if len(enable) > 0 {
		q.Where(maintainer.Enable(enable[0]))
	}
	return q.First(s.ctx)
}

func (s *maintainerService) QueryIDX(id uint64, enable ...bool) *ent.Maintainer {
	x, _ := s.QueryID(id, enable...)
	if x == nil {
		snag.Panic("未找到运维信息")
	}
	return x
}

func (s *maintainerService) QueryPhone(phone string, enable ...bool) (*ent.Maintainer, error) {
	q := s.orm.Query().Where(maintainer.Phone(phone))
	if len(enable) > 0 {
		q.Where(maintainer.Enable(enable[0]))
	}
	return q.First(s.ctx)
}

func (s *maintainerService) QueryPhoneX(phone string, enable ...bool) *ent.Maintainer {
	x, _ := s.QueryPhone(phone, enable...)
	if x == nil {
		snag.Panic("未找到运维信息")
	}
	return x
}

func (s *maintainerService) Detail(item *ent.Maintainer) (data *model.Maintainer) {
	data = &model.Maintainer{
		ID:     item.ID,
		Name:   item.Name,
		Enable: item.Enable,
		Phone:  item.Phone,
		Cities: make([]model.City, len(item.Edges.Cities)),
	}

	for i, c := range item.Edges.Cities {
		data.Cities[i] = model.City{
			ID:   c.ID,
			Name: c.Name,
		}
	}
	return
}

func (s *maintainerService) List(req *model.MaintainerListReq) *model.PaginationRes {
	q := s.orm.Query().WithCities()

	if req.CityID != 0 {
		q.WithCities(func(query *ent.CityQuery) {
			query.Where(city.ID(req.CityID))
		})
	}

	if req.Keyword != "" {
		q.Where(maintainer.Or(
			maintainer.PhoneContainsFold(req.Keyword),
			maintainer.NameContainsFold(req.Keyword),
		))
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Maintainer) *model.Maintainer {
		return s.Detail(item)
	})
}

func (s *maintainerService) hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func (s *maintainerService) checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *maintainerService) Create(req *model.MaintainerCreateReq) {
	q := s.orm.Create().
		SetPhone(req.Phone).
		SetPassword(s.hashPassword(req.Password)).
		SetName(req.Name).
		AddCityIDs(req.CityIDs...)
	if req.Enable != nil {
		q.SetEnable(*req.Enable)
	} else {
		q.SetEnable(true)
	}
	snag.PanicIfError(q.Exec(s.ctx))
}

func (s *maintainerService) Modify(req *model.MaintainerModifyReq) {
	m := s.QueryIDX(req.ID)
	updater := m.Update()

	// 修改密码
	if req.Password != nil {
		// 校对密码，不允许重复
		if s.checkPassword(*req.Password, m.Password) {
			snag.Panic("密码不能重复")
		}
		updater.SetPassword(s.hashPassword(*req.Password))
	}

	// 修改电话
	if req.Phone != nil {
		updater.SetPhone(*req.Phone)
	}

	// 修改姓名
	if req.Name != nil {
		updater.SetName(*req.Name)
	}

	// 修改城市列表
	if len(req.CityIDs) > 0 {
		updater.ClearCities().AddCityIDs(req.CityIDs...)
	}

	snag.PanicIfError(updater.Exec(s.ctx))
}

// Signin 登录
func (s *maintainerService) Signin(req *model.MaintainerSigninReq) (res *model.MaintainerSigninRes) {
	m, _ := s.orm.Query().Where(maintainer.Phone(req.Phone), maintainer.Enable(true)).WithCities().First(s.ctx)
	if m == nil {
		snag.Panic("未找到有效运维")
	}
	if !s.checkPassword(req.Password, m.Password) {
		snag.Panic("电话或密码错误")
	}
	res.Maintainer = s.Detail(m)
	// 生成token
	res.Token = utils.NewEcdsaToken()

	// 存储登录token和ID进行对应
	ar.Redis.HSet(s.ctx, s.tokenCacheKey, res.Token, m.ID)
	ar.Redis.HSet(s.ctx, s.tokenCacheKey, strconv.FormatUint(m.ID, 10), res.Token)
	return
}

// TokenVerify Token校验
func (s *maintainerService) TokenVerify(token string) (m *ent.Maintainer) {
	// 获取token对应ID
	id, _ := ar.Redis.HGet(s.ctx, s.tokenCacheKey, token).Uint64()
	if id <= 0 {
		return
	}

	// 反向校验token是否正确
	if ar.Redis.HGet(s.ctx, s.tokenCacheKey, strconv.FormatUint(id, 10)).Val() != token {
		return
	}

	// 获取运维人员
	m, _ = s.orm.Query().Where(maintainer.ID(id)).WithCities().First(s.ctx)
	if m == nil {
		return
	}
	return
}
