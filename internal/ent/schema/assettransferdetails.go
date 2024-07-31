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

type AssetTransferDetailsMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetTransferDetailsMixin) Fields() []ent.Field {
	relate := field.Uint64("details_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetTransferDetailsMixin) Edges() []ent.Edge {
	e := edge.To("details", AssetTransferDetails.Type).Unique().Field("details_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetTransferDetailsMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("details_id"))
	}
	return
}

// AssetTransferDetails holds the schema definition for the AssetTransferDetails entity.
type AssetTransferDetails struct {
	ent.Schema
}

// Annotations of the AssetTransferDetails.
func (AssetTransferDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_transfer_details"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetTransferDetails.
func (AssetTransferDetails) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("transfer_id").Optional().Comment("调拨ID"),
	}
}

// Edges of the AssetTransferDetails.
func (AssetTransferDetails) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transfer", AssetTransfer.Type).Ref("details").Unique().Field("transfer_id"),
	}
}

func (AssetTransferDetails) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
		AssetMixin{Optional: true},
	}
}

func (AssetTransferDetails) Indexes() []ent.Index {
	return []ent.Index{}
}
