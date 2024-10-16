package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type PurchaseFollowMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
	Prefix       string
}

func (m PurchaseFollowMixin) prefield() (string, string) {
	if m.Prefix == "" {
		return "follow_id", "follow"
	}
	return fmt.Sprintf("%s_follow_id", m.Prefix), fmt.Sprintf("%sPurchaseFollow", m.Prefix)
}

func (m PurchaseFollowMixin) Fields() []ent.Field {
	pf, _ := m.prefield()
	f := field.Uint64(pf).Comment("订单跟进ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m PurchaseFollowMixin) Edges() []ent.Edge {
	pf, pn := m.prefield()
	e := edge.To(pn, PurchaseFollow.Type).Unique().Field(pf)
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PurchaseFollowMixin) Indexes() (arr []ent.Index) {
	pf, _ := m.prefield()
	if !m.DisableIndex {
		arr = append(arr, index.Fields(pf))
	}
	return
}

// PurchaseFollow holds the schema definition for the PurchaseFollow entity.
type PurchaseFollow struct {
	ent.Schema
}

// Annotations of the PurchaseFollow.
func (PurchaseFollow) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "purchase_follow"},
		entsql.WithComments(true),
	}
}

// Fields of the PurchaseFollow.
func (PurchaseFollow) Fields() []ent.Field {
	return []ent.Field{
		field.String("content").Comment("跟进内容"),
		field.Strings("pics").Comment("跟进图片"),
		field.Uint64("order_id").Optional().Comment("订单id"),
	}
}

// Edges of the PurchaseFollow.
func (PurchaseFollow) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", PurchaseOrder.Type).Ref("follows").Unique().Field("order_id"),
	}
}

func (PurchaseFollow) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (PurchaseFollow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("content").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
