// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package biz

import (
	"context"
	"errors"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
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
	// 判定版本是否存在
	q, _ := s.orm.QueryNotDeleted().Where(version.Platform(req.AppPlatform), version.Version(req.Version)).First(s.ctx)
	if q != nil {
		return errors.New("版本已存在")
	}
	_, err = s.orm.Create().
		SetPlatform(req.AppPlatform).
		SetContent(req.Content).
		SetVersion(req.Version).
		SetForce(req.Force).
		SetNillableEnable(req.Enable).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return
}

// Modify 修改
func (s *versionBiz) Modify(req *definition.VersionModifyReq) (err error) {
	// 判定版本是否存在 当前版本不判定
	q, _ := s.orm.QueryNotDeleted().Where(version.Platform(req.AppPlatform), version.Version(req.Version), version.IDNEQ(req.ID)).First(s.ctx)
	if q != nil {
		return errors.New("版本已存在")
	}
	_, err = s.orm.UpdateOneID(req.ID).
		SetPlatform(req.AppPlatform).
		SetContent(req.Content).
		SetVersion(req.Version).
		SetNillableEnable(req.Enable).
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
				Enable:      item.Enable,
			}
		},
	)
}

// LatestVersion 获取最新版本
func (s *versionBiz) LatestVersion(req *definition.LatestVersionReq) *definition.Version {
	q, _ := s.orm.QueryNotDeleted().
		Where(
			version.PlatformEQ(req.AppPlatform),
			version.Enable(true),
		).
		Order(ent.Desc(version.FieldCreatedAt)).First(s.ctx)
	if q == nil {
		return nil
	}

	var link string
	if q.Platform == model.AppPlatformAndroid {
		link = ar.Config.AppDownloadUrl
	}
	return &definition.Version{
		AppPlatform:  q.Platform,
		Version:      q.Version,
		Content:      q.Content,
		Force:        q.Force,
		CreatedAt:    q.CreatedAt.Format(carbon.DateTimeLayout),
		DownloadLink: link,
		Enable:       q.Enable,
	}
}
