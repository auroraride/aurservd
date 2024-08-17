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

type AssetCheckDetailsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetCheckDetailsMixin) Fields() []ent.Field {
	relate := field.Uint64("details_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetCheckDetailsMixin) Edges() []ent.Edge {
	e := edge.To("details", AssetCheckDetails.Type).Unique().Field("details_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetCheckDetailsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("details_id"))
	}
	return
}

// AssetCheckDetails holds the schema definition for the AssetCheckDetails entity.
type AssetCheckDetails struct {
	ent.Schema
}

// Annotations of the AssetCheckDetails.
func (AssetCheckDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_check_details"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetCheckDetails.
func (AssetCheckDetails) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("asset_id").Optional().Comment("资产ID"),
		field.Uint64("check_id").Optional().Comment("盘点ID"),
	}
}

// Edges of the AssetCheckDetails.
func (AssetCheckDetails) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset", Asset.Type).Ref("check_details").Unique().Field("asset_id"),
		edge.From("check", AssetCheck.Type).Ref("check_details").Unique().Field("check_id"),
	}
}

func (AssetCheckDetails) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetCheckDetails) Indexes() []ent.Index {
	return []ent.Index{}
}
