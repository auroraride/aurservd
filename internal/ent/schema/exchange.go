package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
    jsoniter "github.com/json-iterator/go"
)

type ExchangeMixin struct {
    mixin.Schema
    DisableIndex bool
    Optional     bool
}

func (m ExchangeMixin) Fields() []ent.Field {
    f := field.Uint64("exchange_id")
    if m.Optional {
        f.Optional().Nillable()
    }
    return []ent.Field{f}
}

func (m ExchangeMixin) Edges() []ent.Edge {
    e := edge.To("exchange", Exchange.Type).Unique().Field("exchange_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m ExchangeMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("exchange_id"))
    }
    return
}

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
        field.JSON("info", &model.ExchangeInfo{}).Optional().Comment("电柜换电信息"),
        field.String("model").Comment("电池型号"),
        field.Bool("alternative").Default(false).Comment("是否备用方案"),
        field.Time("start_at").Optional().Comment("换电开始时间"),
        field.Time("finish_at").Optional().Comment("换电结束时间"),
        field.Int("duration").Optional().Comment("换电耗时(s)"),
        field.String("rider_battery").Optional().Nillable().Comment("骑手当前电池编号"),
        field.String("putin_battery").Optional().Nillable().Comment("放入电池编号"),
        field.String("putout_battery").Optional().Nillable().Comment("取出电池编号"),
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
