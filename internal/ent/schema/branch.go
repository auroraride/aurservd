package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
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
        field.Float("lng").Comment("经度"),
        field.Float("lat").Comment("纬度"),
        field.String("address").Comment("详细地址"),
        field.Strings("photos").Comment("网点照片"),
        field.Other("geom", &model.Geometry{}).SchemaType(map[string]string{
            dialect.Postgres: "geometry",
        }).Unique().Comment("坐标"),
    }
}

// Edges of the Branch.
func (Branch) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("contracts", BranchContract.Type),
        edge.To("cabinets", Cabinet.Type),
        edge.From("city", City.Type).Ref("branches").Required().Unique().Field("city_id"),
        edge.To("faults", CabinetFault.Type),
        edge.To("stores", Store.Type),
    }
}

func (Branch) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Creator{},
        internal.LastModifier{},
    }
}

func (Branch) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("city_id"),
        index.Fields("lng", "lat"),
        index.Fields("geom").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIST",
            }),
        ),
    }
}
