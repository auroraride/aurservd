// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

type commissionService struct {
	ctx      context.Context
	modifier *model.Modifier
	rider    *ent.Rider
}

func NewCommission() *commissionService {
	return &commissionService{
		ctx: context.Background(),
	}
}

func NewCommissionWithModifier(m *model.Modifier) *commissionService {
	s := NewCommission()
	s.ctx = context.WithValue(s.ctx, "modifier", m)
	s.modifier = m
	return s
}
