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
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
    "time"
)

type BatteryFaultMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m BatteryFaultMixin) Fields() []ent.Field {
    relate := field.Uint64("fault_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m BatteryFaultMixin) Edges() []ent.Edge {
    e := edge.To("fault", BatteryFault.Type).Unique().Field("fault_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m BatteryFaultMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("fault_id"))
    }
    return
}

// BatteryFault holds the schema definition for the BatteryFault entity.
type BatteryFault struct {
    ent.Schema
}

// Annotations of the BatteryFault.
func (BatteryFault) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "battery_fault"},
        entsql.WithComments(true),
    }
}

// Fields of the BatteryFault.
func (BatteryFault) Fields() []ent.Field {
    return []ent.Field{
        field.String("sn").Comment("电池编号"),
        field.Uint64("battery_id").Comment("电池ID"),
        field.Enum("fault").GoType(model.BatteryFault("")).Comment("故障"),
        field.Time("begin_at").Immutable().Default(time.Now).Comment("开始时间"),
        field.Time("end_at").Optional().Comment("结束时间"),
    }
}

// Edges of the BatteryFault.
func (BatteryFault) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("battery", Battery.Type).Field("battery_id").
            Ref("faults").Unique().Required().
            Annotations(
                entsql.WithComments(true),
            ).
            Comment("所属电池"),
    }
}

func (BatteryFault) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
    }
}

func (BatteryFault) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("begin_at"),
        index.Fields("end_at"),
        index.Fields("battery_id"),
        index.Fields("sn").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
            entsql.OpClass("gin_trgm_ops"),
        ),
        index.Fields("fault"),
    }
}
