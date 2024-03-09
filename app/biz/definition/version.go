// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package definition

import "github.com/auroraride/aurservd/app/model"

type VersionReq struct {
	AppPlatform model.AppPlatform `json:"appPlatform" query:"appPlatform" validate:"required,oneof=android ios" trans:"平台" enums:"android,ios"`
}

type VersionRes struct {
	Version
}

type Version struct {
	Version string `json:"version"` // 版本号
	Content string `json:"content"` // 更新内容
	Force   bool   `json:"force"`   // 是否强制更新
}
