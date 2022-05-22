package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Enterprise holds the schema definition for the Enterprise entity.
type Enterprise struct {
    ent.Schema
}

// Annotations of the Enterprise.
func (Enterprise) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise"},
    }
}

// Fields of the Enterprise.
func (Enterprise) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").Comment("企业名称"),
    }
}

// Edges of the Enterprise.
func (Enterprise) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("riders", Rider.Type),
    }
}

func (Enterprise) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (Enterprise) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
