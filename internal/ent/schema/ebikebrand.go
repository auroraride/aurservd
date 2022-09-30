package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type EbikeBrandMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m EbikeBrandMixin) Fields() []ent.Field {
    relate := field.Uint64("brand_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m EbikeBrandMixin) Edges() []ent.Edge {
    e := edge.To("brand", EbikeBrand.Type).Unique().Field("brand_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m EbikeBrandMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("brand_id"))
    }
    return
}

// EbikeBrand holds the schema definition for the EbikeBrand entity.
type EbikeBrand struct {
    ent.Schema
}

// Annotations of the EbikeBrand.
func (EbikeBrand) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "ebike_brand"},
    }
}

// Fields of the EbikeBrand.
func (EbikeBrand) Fields() []ent.Field {
    return []ent.Field{}
}

// Edges of the EbikeBrand.
func (EbikeBrand) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (EbikeBrand) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (EbikeBrand) Indexes() []ent.Index {
    return []ent.Index{}
}
