// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/app/model"
)

type HookModifier struct {
	mixin.Schema
}

type HookModifierMutator interface {
	SetModifier(value *model.Modifier)
}

// Fields of the PointLog.
func (HookModifier) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("modifier", &model.Modifier{}).Optional().Comment("管理"),
	}
}

func (HookModifier) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				ml, ok := m.(HookModifierMutator)
				if ok {
					value, _ := ctx.Value(model.CtxModifierKey{}).(*model.Modifier)
					if value != nil {
						switch op := m.Op(); {
						case op.Is(ent.OpCreate):
							ml.SetModifier(value)
						case op.Is(ent.OpUpdateOne | ent.OpUpdate):
							// TODO: 更新?
						}
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
