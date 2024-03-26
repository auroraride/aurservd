// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-20, by liasica

package biz

import "github.com/auroraride/aurservd/internal/ent"

type exchangeBiz struct {
	orm *ent.ExchangeClient
}

func NewExchange() *exchangeBiz {
	return &exchangeBiz{
		orm: ent.Database.Exchange,
	}
}

func (b *exchangeBiz) Start() {

}
