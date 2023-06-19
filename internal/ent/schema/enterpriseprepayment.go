package schema

import (
	"ariga.io/atlas/sql/postgres"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

// EnterprisePrepayment holds the schema definition for the EnterprisePrepayment entity.
type EnterprisePrepayment struct {
	ent.Schema
}

// Annotations of the EnterprisePrepayment.
func (EnterprisePrepayment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_prepayment"},
		entsql.WithComments(true),
	}
}

// Fields of the EnterprisePrepayment.
func (EnterprisePrepayment) Fields() []ent.Field {
	return []ent.Field{
		field.Float("amount").Immutable().Comment("预付金额"),
		field.Other("payway", model.PaywayCash).Default(model.PaywayCash).SchemaType(map[string]string{
			dialect.Postgres: postgres.TypeSmallInt,
		}).Comment("支付方式").Annotations(entsql.DefaultExpr("1")),
		field.String("trade_no").Optional().Nillable().Comment("支付平台交易单号"),
	}
}

// Edges of the EnterprisePrepayment.
func (EnterprisePrepayment) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (EnterprisePrepayment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.Modifier{},

		EnterpriseMixin{},
		AgentMixin{Optional: true},
	}
}

func (EnterprisePrepayment) Indexes() []ent.Index {
	return []ent.Index{}
}
