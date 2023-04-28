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

// EnterpriseBill holds the schema definition for the EnterpriseBill entity.
type EnterpriseBill struct {
	ent.Schema
}

// Annotations of the EnterpriseBill.
func (EnterpriseBill) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_bill"},
		entsql.WithComments(true),
	}
}

// Fields of the EnterpriseBill.
func (EnterpriseBill) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("subscribe_id").Comment("订阅ID"),
		field.Uint64("enterprise_id").Comment("企业ID"),
		field.Uint64("statement_id").Comment("账单ID"),
		field.Time("start").Comment("结算开始日期(包含)"),
		field.Time("end").Comment("结算结束日期(包含)"),
		field.Int("days").Comment("账单天数"),
		field.Float("price").Comment("账单单价"),
		field.Float("cost").Comment("账单金额"),
		field.String("model").Comment("电池型号"),
	}
}

// Edges of the EnterpriseBill.
func (EnterpriseBill) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("enterprise", Enterprise.Type).Ref("bills").Required().Unique().Field("enterprise_id"),
		edge.From("statement", EnterpriseStatement.Type).Ref("bills").Required().Unique().Field("statement_id"),
		edge.From("subscribe", Subscribe.Type).Ref("bills").Required().Unique().Field("subscribe_id"),
	}
}

func (EnterpriseBill) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		RiderMixin{},
		CityMixin{},
		StationMixin{Optional: true},
	}
}

func (EnterpriseBill) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("model"),
		index.Fields("enterprise_id"),
		index.Fields("statement_id"),
		index.Fields("subscribe_id"),
	}
}
