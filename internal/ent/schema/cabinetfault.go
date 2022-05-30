package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// CabinetFault holds the schema definition for the CabinetFault entity.
type CabinetFault struct {
    ent.Schema
}

// Annotations of the CabinetFault.
func (CabinetFault) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "cabinet_fault"},
    }
}

// Fields of the CabinetFault.
func (CabinetFault) Fields() []ent.Field {
    return []ent.Field{
        field.Uint8("status").Default(0).Comment("故障状态 0未处理 1已处理"),
        // field.JSON("city", model.City{}).Comment("城市"),
        field.Uint64("branch_id").Comment("网点ID"),
        field.Uint64("cabinet_id").Comment("电柜ID"),
        field.Uint64("rider_id").Comment("骑手ID"),
        // field.String("cabinet_name").Comment("电柜名称"),
        // field.String("brand").Comment("电柜品牌"),
        // field.String("serial").Comment("电柜编号"),
        // field.JSON("models", []model.BatteryModel{}).Comment("电池型号"),
        field.String("fault").Optional().Comment("故障内容"),
        field.JSON("attachments", []string{}).Optional().Comment("附件"),
        field.String("description").Optional().Comment("故障留言"),
    }
}

// Edges of the CabinetFault.
func (CabinetFault) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("branch", Branch.Type).Ref("faults").Required().Unique().Field("branch_id"),
        edge.From("cabinet", Cabinet.Type).Ref("faults").Required().Unique().Field("cabinet_id"),
        edge.From("rider", Rider.Type).Ref("faults").Required().Unique().Field("rider_id"),
    }
}

func (CabinetFault) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{},
    }
}

func (CabinetFault) Indexes() []ent.Index {
    return []ent.Index{
        // index.Fields("city").Annotations(
        //     entsql.IndexTypes(map[string]string{
        //         dialect.Postgres: "GIN",
        //     }),
        // ),
        // index.Fields("cabinet_name"),
        index.Fields("status"),
        // index.Fields("models").Annotations(
        //     entsql.IndexTypes(map[string]string{
        //         dialect.Postgres: "GIN",
        //     }),
        // ),
    }
}
