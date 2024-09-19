// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-21, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

type EbikeBrandListReq struct {
	model.PaginationReq
	Keyword *string `json:"keyword" query:"keyword"` // 关键字
}
