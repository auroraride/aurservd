package schema

import (
    "ariga.io/atlas/sql/postgres"
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type EbikeMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m EbikeMixin) Fields() []ent.Field {
    relate := field.Uint64("ebike_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m EbikeMixin) Edges() []ent.Edge {
    e := edge.To("ebike", Ebike.Type).Unique().Field("ebike_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m EbikeMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("ebike_id"))
    }
    return
}

// Ebike holds the schema definition for the Ebike entity.
type Ebike struct {
    ent.Schema
}

// Annotations of the Ebike.
func (Ebike) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "ebike"},
    }
}

// Fields of the Ebike.
func (Ebike) Fields() []ent.Field {
    return []ent.Field{
        field.Other("status", model.EbikeStatusInStock).Default(model.EbikeStatusInStock).SchemaType(map[string]string{
            dialect.Postgres: postgres.TypeSmallInt,
        }).Comment("状态"),
        field.Bool("enable").Default(true).Comment("是否启用"),
        field.String("sn").Unique().Comment("车架号"),
        field.String("plate").Unique().Optional().Nillable().Comment("车牌号"),
        field.String("machine").Unique().Optional().Nillable().Comment("终端编号"),
        field.String("sim").Unique().Optional().Nillable().Comment("SIM卡号"),
        field.String("color").Default(model.EbikeColorDefault).Comment("颜色"),
        field.String("ex_factory").Comment("生产批次(出厂日期)"),
    }
}

// Edges of the Ebike.
func (Ebike) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (Ebike) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.Modifier{},

        EbikeBrandMixin{},

        RiderMixin{Optional: true},
        StoreMixin{Optional: true},
    }
}

func (Ebike) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("status"),
        index.Fields("ex_factory"),
    }
}
