// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-27, by aurb

package model

// ActivityImage 活动图片
type ActivityImage struct {
	List  string `json:"list"`  // 列表图片
	Popup string `json:"popup"` // 弹窗图片
	Home  string `json:"home"`  // 首页icon图片
}
