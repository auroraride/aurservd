package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Commission holds the schema definition for the Commission entity.
type Commission struct {
    ent.Schema
}

// Annotations of the Commission.
func (Commission) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "commission"},
    }
}

// Fields of the Commission.
func (Commission) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("order_id").Comment("订单ID"),
        field.Float("amount").Immutable().Comment("提成金额"),
        field.Uint8("status").Default(0).Comment("提成状态 0未发放 1已发放"),
        field.Uint64("employee_id").Optional().Comment("员工ID"),
    }
}

// Edges of the Commission.
func (Commission) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("order", Order.Type).Ref("commission").Unique().Required().Field("order_id"),
    }
}

func (Commission) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Commission) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("status"),
    }
}