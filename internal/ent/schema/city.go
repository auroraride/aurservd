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

// City holds the schema definition for the City entity.
type City struct {
    ent.Schema
}

// Annotations of the City.
func (City) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "city"},
    }
}

// Fields of the City.
func (City) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("id"),
        field.Bool("open").Optional().Nillable().Comment("启用"),
        field.String("name").MaxLen(100).Comment("城市"),
        field.String("code").MaxLen(10).Comment("城市编号"),
        field.Uint64("parent_id").Optional().Nillable().Comment("父级"),
        field.Float("lng").Optional().Comment("经度"),
        field.Float("lat").Optional().Comment("纬度"),
    }
}

// Edges of the City.
func (City) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("plans", Plan.Type).Ref("cities"),
        edge.To("children", City.Type).From("parent").Field("parent_id").Unique(),
        edge.To("branches", Branch.Type),
        edge.To("faults", CabinetFault.Type),
        edge.To("orders", Order.Type),
        // edge.To("commissions", Commission.Type),
    }
}

func (City) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (City) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("open"),
        index.Fields("parent_id"),
    }
}
