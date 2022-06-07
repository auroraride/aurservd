package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// EnterpriseStation holds the schema definition for the EnterpriseStation entity.
type EnterpriseStation struct {
    ent.Schema
}

// Annotations of the EnterpriseStation.
func (EnterpriseStation) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise_station"},
    }
}

// Fields of the EnterpriseStation.
func (EnterpriseStation) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("enterprise_id").Comment("企业ID"),
        field.String("name").Comment("站点名称"),
    }
}

// Edges of the EnterpriseStation.
func (EnterpriseStation) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("enterprise", Enterprise.Type).Ref("stations").Unique().Required().Field("enterprise_id"),
    }
}

func (EnterpriseStation) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (EnterpriseStation) Indexes() []ent.Index {
    return []ent.Index{}
}
