package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type PurchaseCommodityMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PurchaseCommodityMixin) Fields() []ent.Field {
	relate := field.Uint64("commodity_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PurchaseCommodityMixin) Edges() []ent.Edge {
	e := edge.To("commodity", PurchaseCommodity.Type).Unique().Field("commodity_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PurchaseCommodityMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("commodity_id"))
	}
	return
}

// PurchaseCommodity holds the schema definition for the PurchaseCommodity entity.
type PurchaseCommodity struct {
	ent.Schema
}

// Annotations of the PurchaseCommodity.
func (PurchaseCommodity) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "purchase_commodity"},
		entsql.WithComments(true),
	}
}

// Fields of the PurchaseCommodity.
func (PurchaseCommodity) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the PurchaseCommodity.
func (PurchaseCommodity) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PurchaseCommodity) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (PurchaseCommodity) Indexes() []ent.Index {
	return []ent.Index{}
}
