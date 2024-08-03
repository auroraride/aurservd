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

type AssetScrapDetailsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetScrapDetailsMixin) Fields() []ent.Field {
	relate := field.Uint64("details_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetScrapDetailsMixin) Edges() []ent.Edge {
	e := edge.To("details", AssetScrapDetails.Type).Unique().Field("details_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetScrapDetailsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("details_id"))
	}
	return
}

// AssetScrapDetails holds the schema definition for the AssetScrapDetails entity.
type AssetScrapDetails struct {
	ent.Schema
}

// Annotations of the AssetScrapDetails.
func (AssetScrapDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_scrap_details"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetScrapDetails.
func (AssetScrapDetails) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("asset_id").Comment("资产ID"),
		field.Uint64("scrap_id").Optional().Comment("报废ID"),
	}
}

// Edges of the AssetScrapDetails.
func (AssetScrapDetails) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("asset", Asset.Type).Ref("scrap_details").Unique().Field("asset_id").Required(),
		edge.From("scrap", AssetScrap.Type).Ref("scrap_details").Unique().Field("scrap_id"),
	}
}

func (AssetScrapDetails) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		MaterialMixin{Optional: true},
	}
}

func (AssetScrapDetails) Indexes() []ent.Index {
	return []ent.Index{}
}
