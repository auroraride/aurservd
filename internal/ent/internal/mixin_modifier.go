// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/auroraride/aurservd/app/model"
)

// Modifier 修改或创建人
type Modifier struct {
	mixin.Schema
	IndexCreator bool // 是否索引创建人
}

func (Modifier) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("creator", &model.Modifier{}).Immutable().Optional().Comment("创建人"),
		field.JSON("last_modifier", &model.Modifier{}).Optional().Comment("最后修改人"),
		field.String("remark").Optional().Comment("管理员改动原因/备注"),
	}
}

func (m Modifier) Indexes() []ent.Index {
	var indexes []ent.Index
	if m.IndexCreator {
		indexes = append(indexes, index.Fields("creator").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
		))
	}
	return indexes
}

func (Modifier) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				ml, ok := m.(model.ModifierLogger)
				if ok {
					mod := model.GetModifierFromContext(ctx)
					if mod != nil {
						switch op := m.Op(); {
						case op.Is(ent.OpCreate):
							ml.SetCreator(mod)
							ml.SetLastModifier(mod)
						case op.Is(ent.OpUpdateOne | ent.OpUpdate):
							ml.SetLastModifier(mod)
						}
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
