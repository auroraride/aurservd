package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type StoreMixin struct {
    mixin.Schema
}

func (StoreMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("store_id").Optional(),
    }
}

func (StoreMixin) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("store", Store.Type).Unique().Field("store_id"),
    }
}

// Store holds the schema definition for the Store entity.
type Store struct {
    ent.Schema
}

// Annotations of the Store.
func (Store) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "store"},
    }
}

// Fields of the Store.
func (Store) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("branch_id").Comment("网点ID"),
        field.String("sn").Immutable().Comment("门店编号"),
        field.String("name").Comment("门店名称"),
        field.Uint8("status").Default(0).Comment("门店状态 0维护 1营业 2休息 3隐藏"),
    }
}

// Edges of the Store.
func (Store) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("branch", Branch.Type).Ref("stores").Required().Unique().Field("branch_id"),
    }
}

func (Store) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Store) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("status"),
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
