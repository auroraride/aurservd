// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-21
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
	"github.com/auroraride/aurservd/app/model"
)

func (pl *Plan) BasicInfo() *model.Plan {
	return &model.Plan{
		ID:             pl.ID,
		Name:           pl.Name,
		Days:           pl.Days,
		Intelligent:    pl.Intelligent,
		CommissionBase: pl.CommissionBase,
		Original:       pl.Original,
	}
}

func (pc *PlanCreate) Clone() (creator *PlanCreate) {
	mutation := new(PlanMutation)
	*mutation = *pc.mutation
	return &PlanCreate{
		config:   pc.config,
		mutation: mutation,
		hooks:    pc.hooks,
		conflict: pc.conflict,
	}
}
