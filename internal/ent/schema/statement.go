package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Statement holds the schema definition for the Statement entity.
type Statement struct {
    ent.Schema
}

// Annotations of the Statement.
func (Statement) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "statement"},
    }
}

// Fields of the Statement.
func (Statement) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("enterprise_id").Comment("企业ID"),
        field.Float("arrearage").Default(0).Comment("欠费金额"),
        field.Float("amount").Default(0).Comment("账单金额"),
        field.Time("settled_at").Optional().Nillable().Comment("清账时间"),
        field.Uint("days").Default(0).Comment("账期内使用总天数"),
        field.Uint("rider_number").Default(0).Comment("账期内使用总人数"),
    }
}

// Edges of the Statement.
func (Statement) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("subscribes", Subscribe.Type),
        edge.From("enterprise", Enterprise.Type).Ref("statements").Unique().Required().Field("enterprise_id"),
    }
}

func (Statement) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Statement) Indexes() []ent.Index {
    return []ent.Index{}
}
