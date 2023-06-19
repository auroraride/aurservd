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

// EnterprisePrice holds the schema definition for the EnterprisePrice entity.
type EnterprisePrice struct {
	ent.Schema
}

// Annotations of the EnterprisePrice.
func (EnterprisePrice) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enterprise_price"},
		entsql.WithComments(true),
	}
}

// Fields of the EnterprisePrice.
func (EnterprisePrice) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("enterprise_id"),
		field.Float("price").Comment("单价 元/天"),
		field.String("model").Comment("可用电池型号"),
		field.Bool("intelligent").Default(false).Comment("是否智能电池"),
		field.String("key").Optional().Comment("价格key"),
	}
}

// Edges of the EnterprisePrice.
func (EnterprisePrice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("enterprise", Enterprise.Type).
			Ref("prices").
			Unique().
			Required().
			Field("enterprise_id"),
	}
}

func (EnterprisePrice) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},

		CityMixin{Optional: false},
		EbikeBrandMixin{Optional: true},
	}
}

func (EnterprisePrice) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("model"),
		index.Fields("enterprise_id"),
	}
}
