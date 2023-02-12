package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
    jsoniter "github.com/json-iterator/go"
)

// Exchange holds the schema definition for the Exchange entity.
type Exchange struct {
    ent.Schema
}

// Annotations of the Exchange.
func (Exchange) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "exchange"},
        entsql.WithComments(true),
    }
}

// Fields of the Exchange.
func (Exchange) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("rider_id").Comment("骑手ID"),
        field.Uint64("employee_id").Optional().Nillable().Comment("店员ID"),
        field.String("uuid").Unique(),
        field.Uint64("cabinet_id").Optional().Comment("电柜ID"),
        field.Bool("success").Default(true).Comment("是否成功"),
        field.JSON("detail", jsoniter.RawMessage{}).Optional().Comment("电柜换电信息"),
        // field.JSON("info", &model.ExchangeInfo{}).Optional().Comment("电柜换电信息"),
        field.String("model").Comment("电池型号"),
        field.Bool("alternative").Default(false).Comment("是否备用方案"),
        field.Time("start_at").Optional().Comment("换电开始时间"),
        field.Time("finish_at").Optional().Comment("换电结束时间"),
        field.Int("duration").Optional().Comment("换电耗时(s)"),
        field.String("rider_battery").Optional().Nillable().Comment("骑手当前电池编号"),
        field.String("putin_battery").Optional().Nillable().Comment("放入电池编号"),
        field.String("putout_battery").Optional().Nillable().Comment("取出电池编号"),
        field.JSON("cabinet_info", &model.ExchangeCabinetInfo{}).Optional().Comment("电柜信息"),
        field.JSON("empty", &model.BinInfo{}).Optional().Comment("空仓信息"),
        field.JSON("fully", &model.BinInfo{}).Optional().Comment("满仓信息"),
        field.JSON("steps", []*model.ExchangeStepInfo{}).Optional().Comment("步骤信息"),
        field.String("message").Default("").Comment("消息"),
    }
}

// Edges of the Exchange.
func (Exchange) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("cabinet", Cabinet.Type).Unique().Ref("exchanges").Field("cabinet_id"),
        edge.From("rider", Rider.Type).Unique().Required().Ref("exchanges").Field("rider_id"),
        edge.From("employee", Employee.Type).Unique().Ref("exchanges").Field("employee_id"),
    }
}

func (Exchange) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        SubscribeMixin{},
        CityMixin{},

        StoreMixin{Optional: true},
        EnterpriseMixin{Optional: true},
        StationMixin{Optional: true},
    }
}

func (Exchange) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("cabinet_id"),
        index.Fields("rider_id"),
        index.Fields("employee_id"),
        index.Fields("success"),
        index.Fields("model"),
        index.Fields("rider_battery"),
        index.Fields("putin_battery"),
        index.Fields("putout_battery"),
    }
}
