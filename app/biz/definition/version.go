// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-09, by liasica

package definition

import "github.com/auroraride/aurservd/app/model"

type VersionReq struct {
	Platform []string `json:"platform" validate:"required,dive,required,oneof=android ios" trans:"平台"` // 平台
	Version  string   `json:"version" validate:"required"`                                             // 版本号
	Content  string   `json:"content" validate:"required"`                                             // 更新内容
	Force    bool     `json:"force" `                                                                  // 是否强制更新
}

type VersionRes struct {
	Version
}

type Version struct {
	ID        uint64   `json:"id,omitempty"` // ID
	Platform  []string `json:"platform"`     // 平台
	Version   string   `json:"version"`      // 版本号
	Content   string   `json:"content"`      // 更新内容
	Force     bool     `json:"force"`        // 是否强制更新
	CreatedAt string   `json:"createdAt"`    // 创建时间
}

type VersionModifyReq struct {
	model.IDParamReq
	VersionReq
}

type VersionListReq struct {
	model.PaginationReq
}
