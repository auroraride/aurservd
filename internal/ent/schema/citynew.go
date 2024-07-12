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

type CityNewMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m CityNewMixin) Fields() []ent.Field {
	relate := field.Uint64("new_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m CityNewMixin) Edges() []ent.Edge {
	e := edge.To("new", CityNew.Type).Unique().Field("new_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m CityNewMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("new_id"))
	}
	return
}

// CityNew holds the schema definition for the CityNew entity.
type CityNew struct {
	ent.Schema
}

// Annotations of the CityNew.
func (CityNew) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "city_new"},
		entsql.WithComments(true),
	}
}

// Fields of the CityNew.
func (CityNew) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Bool("open").Optional().Nillable().Comment("启用"),
		field.String("name").MaxLen(100).Comment("城市"),
		field.String("code").MaxLen(10).Comment("城市编号"),
		field.Uint64("parent_id").Optional().Nillable().Comment("父级"),
		field.Float("lng").Optional().Comment("经度"),
		field.Float("lat").Optional().Comment("纬度"),
	}
}

// Edges of the CityNew.
func (CityNew) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CityNew) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (CityNew) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("open"),
		index.Fields("parent_id"),
	}
}
