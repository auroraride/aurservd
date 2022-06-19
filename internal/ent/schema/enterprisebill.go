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

// EnterpriseBill holds the schema definition for the EnterpriseBill entity.
type EnterpriseBill struct {
    ent.Schema
}

// Annotations of the EnterpriseBill.
func (EnterpriseBill) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "EnterpriseBill"},
    }
}

// Fields of the EnterpriseBill.
func (EnterpriseBill) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("enterprise_id").Comment("企业ID"),
        field.Uint64("statement_id").Comment("账单ID"),
        field.Time("start").SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("结算开始日期(包含)"),
        field.Time("end").SchemaType(map[string]string{dialect.Postgres: "date"}).Comment("结算结束日期(包含)"),
        field.Int("days").Comment("账单日期"),
        field.Float("price").Comment("账单单价"),
        field.Float("cost").Comment("账单金额"),
        field.Float("voltage").Comment("电压型号"),
    }
}

// Edges of the EnterpriseBill.
func (EnterpriseBill) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("enterprise", Enterprise.Type).Ref("bills").Required().Unique().Field("enterprise_id"),
        edge.From("statement", EnterpriseStatement.Type).Ref("bills").Required().Unique().Field("statement_id"),
    }
}

func (EnterpriseBill) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        RiderMixin{},
        SubscribeMixin{},
        CityMixin{},
    }
}

func (EnterpriseBill) Indexes() []ent.Index {
    return []ent.Index{}
}
