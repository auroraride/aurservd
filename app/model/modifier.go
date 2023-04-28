// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-23
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "context"

// Modifier 修改人
type Modifier struct {
	ID    uint64 `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type ModifierLogger interface {
	SetLastModifier(m *Modifier)
	SetCreator(m *Modifier)
	Creator() (m *Modifier, exists bool)
	LastModifier() (m *Modifier, exists bool)
}

func GetModifierFromContext(ctx context.Context) *Modifier {
	v, ok := ctx.Value("modifier").(*Modifier)
	if ok {
		return v
	}
	return nil
}
