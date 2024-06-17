package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

// StoreGoods holds the schema definition for the StoreGoods entity.
type StoreGoods struct {
	ent.Schema
}

// Annotations of the StoreGoods.
func (StoreGoods) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "store_goods"},
		entsql.WithComments(true),
	}
}

// Fields of the StoreGoods.
func (StoreGoods) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("goods_id").Optional(),
		field.Uint64("store_id").Optional(),
	}
}

// Edges of the StoreGoods.
func (StoreGoods) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("goods", Goods.Type).Ref("stores").Unique().Field("goods_id"),
		edge.From("store", Store.Type).Ref("goods").Unique().Field("store_id"),
	}
}

func (StoreGoods) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
	}
}
