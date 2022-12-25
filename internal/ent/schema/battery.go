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
    }
}

// Fields of the Battery.
func (Battery) Fields() []ent.Field {
    return []ent.Field{
        field.String("sn").Unique().Comment("电池编号"),
        field.Bool("enable").Default(true).Comment("是否启用"),
        field.String("model").Comment("电池型号"),
    }
}

// Edges of the Battery.
func (Battery) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (Battery) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{},                  // 所在城市
        RiderMixin{Optional: true},   // 所在骑手
        CabinetMixin{Optional: true}, // 所在电柜
    }
}

func (Battery) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("enable"),
        index.Fields("model").StorageKey("index_battery_model"),
        index.Fields("sn").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}
