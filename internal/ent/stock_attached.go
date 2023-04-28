// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-16
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

func (sc *StockCreate) Clone() (creator *StockCreate) {
	mutation := new(StockMutation)
	*mutation = *sc.mutation
	return &StockCreate{
		config:   sc.config,
		mutation: mutation,
		hooks:    sc.hooks,
		conflict: sc.conflict,
	}
}
