package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// Branch holds the schema definition for the Branch entity.
type Branch struct {
    ent.Schema
}

// Annotations of the Branch.
func (Branch) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "branch"},
    }
}

// Fields of the Branch.
func (Branch) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("city_id").Comment("城市ID"),
        field.String("name").Comment("网点名称"),
        field.Uint64("lng").Comment("经度"),
        field.Uint64("lat").Comment("纬度"),
        field.String("address").Comment("详细地址"),
    }
}

// Edges of the Branch.
func (Branch) Edges() []ent.Edge {
    return nil
}

func (Branch) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Creator{},
        internal.LastModify{},
    }
}

func (Branch) Indexes() []ent.Index {
    return nil
}
