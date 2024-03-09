// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package biz

import "github.com/auroraride/aurservd/internal/ent"

type versionBiz struct {
	orm *ent.VersionClient
}

func NewVersion() *versionBiz {
	return &versionBiz{
		orm: ent.Database.Version,
	}
}
