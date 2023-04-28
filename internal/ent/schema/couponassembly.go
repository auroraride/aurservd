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
)

type CouponAssemblyMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m CouponAssemblyMixin) Fields() []ent.Field {
	relate := field.Uint64("assembly_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m CouponAssemblyMixin) Edges() []ent.Edge {
	e := edge.To("assembly", CouponAssembly.Type).Unique().Field("assembly_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m CouponAssemblyMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("assembly_id"))
	}
	return
}

// CouponAssembly holds the schema definition for the CouponAssembly entity.
type CouponAssembly struct {
	ent.Schema
}

// Annotations of the CouponAssembly.
func (CouponAssembly) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "coupon_assembly"},
		entsql.WithComments(true),
	}
}

// Fields of the CouponAssembly.
func (CouponAssembly) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("名称"),
		field.Int("number").Comment("数量"),
		field.Float("amount").Comment("金额"),
		field.Uint8("target").Comment("发送对象"),
		field.JSON("meta", &model.CouponTemplateMeta{}).Comment("详情"),
	}
}

// Edges of the CouponAssembly.
func (CouponAssembly) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CouponAssembly) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.Modifier{},
		CouponTemplateMixin{},
	}
}

func (CouponAssembly) Indexes() []ent.Index {
	return []ent.Index{}
}
