package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "github.com/auroraride/aurservd/internal/ent/internal"
)

// EnterpriseStatement holds the schema definition for the EnterpriseStatement entity.
type EnterpriseStatement struct {
    ent.Schema
}

// Annotations of the EnterpriseStatement.
func (EnterpriseStatement) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "enterprise_statement"},
    }
}

// Fields of the EnterpriseStatement.
func (EnterpriseStatement) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("enterprise_id").Comment("企业ID"),
        field.Float("cost").Default(0).Comment("账单金额"),
        field.Float("amount").Default(0).Comment("总预付金额"),
        field.Float("balance").Default(0).Comment("预付剩余, 负数是欠费"),
        field.Time("settled_at").Optional().Nillable().Comment("清账时间"),
        field.Int("days").Default(0).Comment("账期内使用总天数"),
        field.Int("rider_number").Default(0).Comment("账期内使用总人数"),
        field.Time("bill_time").Optional().Nillable().
            SchemaType(map[string]string{dialect.Postgres: "date"}).
            Comment("对账单计算日期(包含, 例如2022-06-05代表是2022-06-06日计算截止到2022-06-05的账单详情)"),
    }
}

// Edges of the EnterpriseStatement.
func (EnterpriseStatement) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("subscribes", Subscribe.Type),
        edge.From("enterprise", Enterprise.Type).Ref("statements").Unique().Required().Field("enterprise_id"),
    }
}

func (EnterpriseStatement) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},
    }
}

func (EnterpriseStatement) Indexes() []ent.Index {
    return []ent.Index{}
}