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

type PromotionPrivilegeMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionPrivilegeMixin) Fields() []ent.Field {
	relate := field.Uint64("privilege_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionPrivilegeMixin) Edges() []ent.Edge {
	e := edge.To("privilege", PromotionPrivilege.Type).Unique().Field("privilege_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionPrivilegeMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("privilege_id"))
	}
	return
}

// PromotionPrivilege holds the schema definition for the PromotionPrivilege entity.
type PromotionPrivilege struct {
	ent.Schema
}

// Annotations of the PromotionPrivilege.
func (PromotionPrivilege) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_privilege"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionPrivilege.
func (PromotionPrivilege) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("type").Default(0).Comment("权益类型 0:无权益 1: 佣金提高(%)"),
		field.String("name").Comment("权益名称"),
		field.String("description").Optional().Nillable().Comment("权益描述"),
		field.Uint64("value").Optional().Nillable().Default(0).Comment("权益值"),
	}
}

// Edges of the PromotionPrivilege.
func (PromotionPrivilege) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (PromotionPrivilege) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (PromotionPrivilege) Indexes() []ent.Index {
	return []ent.Index{}
}
