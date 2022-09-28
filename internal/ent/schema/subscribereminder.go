package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// SubscribeReminder holds the schema definition for the SubscribeReminder entity.
type SubscribeReminder struct {
    ent.Schema
}

// Annotations of the SubscribeReminder.
func (SubscribeReminder) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "subscribe_reminder"},
    }
}

// Fields of the SubscribeReminder.
func (SubscribeReminder) Fields() []ent.Field {
    return []ent.Field{
        field.Enum("type").Values("sms", "vms").Comment("催费类型"),
        field.String("phone").Comment("电话"),
        field.String("name").Comment("姓名"),
        field.Bool("success").Comment("是否成功"),
        field.Int("days").Comment("剩余天数"),
        field.String("plan_name").Comment("套餐名称"),
        field.String("date").Comment("发送日期"),
        field.Float("fee").Default(0).Comment("逾期费用"),
        field.String("fee_formula").Optional().Comment("逾期费用计算公式"),
    }
}

// Edges of the SubscribeReminder.
func (SubscribeReminder) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (SubscribeReminder) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{DisableIndex: true},
        SubscribeMixin{},
        PlanMixin{},
    }
}

func (SubscribeReminder) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("type"),
        index.Fields("date"),
    }
}
