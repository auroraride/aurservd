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

type PromotionPersonMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PromotionPersonMixin) Fields() []ent.Field {
	relate := field.Uint64("person_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PromotionPersonMixin) Edges() []ent.Edge {
	e := edge.To("person", PromotionPerson.Type).Unique().Field("person_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PromotionPersonMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("person_id"))
	}
	return
}

// PromotionPerson holds the schema definition for the PromotionPerson entity.
type PromotionPerson struct {
	ent.Schema
}

// Annotations of the PromotionPerson.
func (PromotionPerson) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "promotion_person"},
		entsql.WithComments(true),
	}
}

// Fields of the PromotionPerson.
func (PromotionPerson) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("status").Default(0).Comment("认证状态 0未认证 1已认证 2认证失败"),
		field.String("name").Optional().MaxLen(40).Comment("真实姓名"),
		field.String("id_card_number").Optional().MaxLen(40).Comment("证件号码").Unique(),
		field.String("address").Optional().Comment("地址"),
	}
}

// Edges of the PromotionPerson.
func (PromotionPerson) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("member", PromotionMember.Type),
	}
}

func (PromotionPerson) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
	}
}

func (PromotionPerson) Indexes() []ent.Index {
	return []ent.Index{}
}
