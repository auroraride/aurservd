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

type AssetAttributesMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetAttributesMixin) Fields() []ent.Field {
	relate := field.Uint64("attributes_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetAttributesMixin) Edges() []ent.Edge {
	e := edge.To("attributes", AssetAttributes.Type).Unique().Field("attributes_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetAttributesMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("attributes_id"))
	}
	return
}

// AssetAttributes holds the schema definition for the AssetAttributes entity.
type AssetAttributes struct {
	ent.Schema
}

// Annotations of the AssetAttributes.
func (AssetAttributes) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_attributes"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetAttributes.
func (AssetAttributes) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Optional().Comment("名称"),
	}
}

// Edges of the AssetAttributes.
func (AssetAttributes) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("values", AssetAttributeValues.Type),
	}
}

func (AssetAttributes) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetAttributes) Indexes() []ent.Index {
	return []ent.Index{}
}
