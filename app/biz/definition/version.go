// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package definition

import "github.com/auroraride/aurservd/app/model"

type VersionReq struct {
	AppPlatform model.AppPlatform `json:"appPlatform" query:"appPlatform" validate:"required,oneof=android ios" trans:"平台" enums:"android,ios"`
	Version     string            `json:"version" validate:"required"` // 版本号
	Content     string            `json:"content" validate:"required"` // 更新内容
	Force       bool              `json:"force"`                       // 是否强制更新
	Enable      bool              `json:"enable"  validate:"required"` // 是否启用
}

type VersionRes struct {
	Version
}

type Version struct {
	ID           uint64            `json:"id,omitempty"`           // ID
	AppPlatform  model.AppPlatform `json:"appPlatform,omitempty"`  // 平台
	Version      string            `json:"version"`                // 版本号
	Content      string            `json:"content"`                // 更新内容
	Force        bool              `json:"force"`                  // 是否强制更新
	CreatedAt    string            `json:"createdAt"`              // 创建时间
	DownloadLink string            `json:"downloadLink,omitempty"` // 下载链接(Android)
	Enable       bool              `json:"enable"`                 // 是否启用
}

type VersionModifyReq struct {
	model.IDParamReq
	VersionReq
}

type VersionListReq struct {
	model.PaginationReq
}

type LatestVersionReq struct {
	AppPlatform model.AppPlatform `json:"appPlatform" query:"appPlatform" validate:"required,oneof=android ios" enums:"android,ios"` // 平台
}
