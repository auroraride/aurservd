// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package biz

import (
	"context"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/version"
)

type versionBiz struct {
	orm      *ent.VersionClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewVersion() *versionBiz {
	return &versionBiz{
		orm: ent.Database.Version,
		ctx: context.Background(),
	}
}

func NewVersionWithModifierBiz(m *model.Modifier) *versionBiz {
	s := NewVersion()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// Create 创建版本
func (s *versionBiz) Create(req *definition.VersionReq) (err error) {
	_, err = s.orm.Create().
		SetPlatform(req.AppPlatform).
		SetContent(req.Content).
		SetVersion(req.Version).
		SetForce(req.Force).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return
}

// Modify 修改
func (s *versionBiz) Modify(req *definition.VersionModifyReq) (err error) {
	_, err = s.orm.UpdateOneID(req.ID).
		SetPlatform(req.AppPlatform).
		SetContent(req.Content).
		SetVersion(req.Version).
		SetForce(req.Force).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return
}

// Delete 删除
func (s *versionBiz) Delete(id uint64) (err error) {
	err = s.orm.DeleteOneID(id).Exec(s.ctx)
	if err != nil {
		return err
	}
	return
}

// List 列表
func (s *versionBiz) List(req *definition.VersionListReq) (res *model.PaginationRes) {
	q := s.orm.QueryNotDeleted().Order(ent.Desc(version.FieldCreatedAt))
	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Version) (res *definition.Version) {
			return &definition.Version{
				ID:          item.ID,
				AppPlatform: item.Platform,
				Version:     item.Version,
				Content:     item.Content,
				Force:       item.Force,
				CreatedAt:   item.CreatedAt.Format(carbon.DateTimeLayout),
			}
		},
	)
}

// LatestVersion 获取最新版本
func (s *versionBiz) LatestVersion(req *definition.LatestVersionReq) *definition.Version {
	q, _ := s.orm.QueryNotDeleted().
		Where(version.PlatformEQ(req.AppPlatform)).
		Order(ent.Desc(version.FieldCreatedAt)).First(s.ctx)
	if q == nil {
		return nil
	}
	return &definition.Version{
		AppPlatform: q.Platform,
		Version:     q.Version,
		Content:     q.Content,
		Force:       q.Force,
		CreatedAt:   q.CreatedAt.Format(carbon.DateTimeLayout),
	}
}
