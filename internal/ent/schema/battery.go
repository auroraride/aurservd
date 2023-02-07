package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

type BatteryMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m BatteryMixin) Fields() []ent.Field {
    relate := field.Uint64("battery_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m BatteryMixin) Edges() []ent.Edge {
    e := edge.To("battery", Battery.Type).Unique().Field("battery_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m BatteryMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("battery_id"))
    }
    return
}

// Battery holds the schema definition for the Battery entity.
type Battery struct {
    ent.Schema
}

// Annotations of the Battery.
func (Battery) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "battery"},
        entsql.WithComments(true),
    }
}

// Fields of the Battery.
func (Battery) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Optional().Nillable().Comment("骑手ID"),
        field.Uint64("cabinet_id").Optional().Nillable().Comment("电柜ID"),
        field.Uint64("subscribe_id").Optional().Nillable().Comment("订阅ID"),
        field.String("sn").Unique().Comment("电池编号"),
        field.Bool("enable").Default(true).Comment("是否启用"),
        field.String("model").Comment("电池型号"),
        field.Int("ordinal").Optional().Nillable().Comment("所在智能柜仓位序号"),
    }
}

// Edges of the Battery.
func (Battery) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("rider", Rider.Type).Ref("battery").Unique().Field("rider_id").Comment("所属骑手"),
        edge.From("cabinet", Cabinet.Type).Ref("batteries").Unique().Field("cabinet_id").Comment("所属电柜"),
        edge.From("subscribe", Subscribe.Type).Ref("battery").Unique().Field("subscribe_id").Comment("所属订阅"),

        edge.To("flows", BatteryFlow.Type).Annotations(
            entsql.WithComments(true),
        ).Comment("流转记录"),
    }
}

func (Battery) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{Optional: true}, // 所在城市
    }
}

func (Battery) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("cabinet_id", "ordinal"),
        index.Fields("enable"),
        index.Fields("model").StorageKey("index_battery_model"),
        index.Fields("sn").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
            entsql.OpClass("gin_trgm_ops"),
        ),
    }
}
