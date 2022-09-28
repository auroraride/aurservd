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

type CityMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m CityMixin) Fields() []ent.Field {
    f := field.Uint64("city_id").Comment("城市ID")
    if m.Optional {
        f.Optional().Nillable()
    }
    return []ent.Field{f}
}

func (m CityMixin) Edges() []ent.Edge {
    e := edge.To("city", City.Type).Unique().Field("city_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m CityMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("city_id"))
    }
    return
}

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
        edge.To("children", City.Type).From("parent").Field("parent_id").Unique(),

        edge.From("plans", Plan.Type).Ref("cities"),
        edge.From("coupons", Coupon.Type).Ref("cities"),
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
