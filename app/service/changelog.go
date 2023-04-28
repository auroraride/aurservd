// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"

	"github.com/auroraride/aurservd/assets/docs"
)

type changelogService struct {
	ctx context.Context
}

func NewChangelog() *changelogService {
	return &changelogService{
		ctx: context.Background(),
	}
}

func (s *changelogService) Generate() {
	tmpl := docs.SwaggerInfo.SwaggerTemplate
	fmt.Println(tmpl)
}
