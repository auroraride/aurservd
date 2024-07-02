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

type EbikeBrandAttributeMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m EbikeBrandAttributeMixin) Fields() []ent.Field {
	relate := field.Uint64("attribute_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m EbikeBrandAttributeMixin) Edges() []ent.Edge {
	e := edge.To("attribute", EbikeBrandAttribute.Type).Unique().Field("attribute_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m EbikeBrandAttributeMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("attribute_id"))
	}
	return
}

// EbikeBrandAttribute holds the schema definition for the EbikeBrandAttribute entity.
type EbikeBrandAttribute struct {
	ent.Schema
}

// Annotations of the EbikeBrandAttribute.
func (EbikeBrandAttribute) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "ebike_brand_attribute"},
		entsql.WithComments(true),
	}
}

// Fields of the EbikeBrandAttribute.
func (EbikeBrandAttribute) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("属性名称"),
		field.String("value").Comment("属性值"),
		field.Uint64("brand_id").Optional().Comment("品牌ID"),
	}
}

// Edges of the EbikeBrandAttribute.
func (EbikeBrandAttribute) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("brand", EbikeBrand.Type).Ref("brand_attribute").Unique().Field("brand_id"),
	}
}

func (EbikeBrandAttribute) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (EbikeBrandAttribute) Indexes() []ent.Index {
	return []ent.Index{}
}
