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

type AssetAttributeValuesMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetAttributeValuesMixin) Fields() []ent.Field {
	relate := field.Uint64("values_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetAttributeValuesMixin) Edges() []ent.Edge {
	e := edge.To("values", AssetAttributeValues.Type).Unique().Field("values_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetAttributeValuesMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("values_id"))
	}
	return
}

// AssetAttributeValues holds the schema definition for the AssetAttributeValues entity.
type AssetAttributeValues struct {
	ent.Schema
}

// Annotations of the AssetAttributeValues.
func (AssetAttributeValues) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_attribute_values"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetAttributeValues.
func (AssetAttributeValues) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("attribute_id").Optional().Comment("属性ID"),
		field.String("value").MaxLen(255).Comment("属性值"),
	}
}

// Edges of the AssetAttributeValues.
func (AssetAttributeValues) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("attribute", AssetAttributes.Type).Ref("values").Unique().Field("attribute_id"),
	}
}

func (AssetAttributeValues) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetAttributeValues) Indexes() []ent.Index {
	return []ent.Index{}
}
