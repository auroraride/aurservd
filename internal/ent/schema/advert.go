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

type AdvertMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AdvertMixin) Fields() []ent.Field {
	relate := field.Uint64("advert_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AdvertMixin) Edges() []ent.Edge {
	e := edge.To("advert", Advert.Type).Unique().Field("advert_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AdvertMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("advert_id"))
	}
	return
}

// Advert holds the schema definition for the Advert entity.
type Advert struct {
	ent.Schema
}

// Annotations of the Advert.
func (Advert) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "advert"},
		entsql.WithComments(true),
	}
}

// Fields of the Advert.
func (Advert) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("名称"),
		field.String("image").Comment("图片"),
		field.String("link").Comment("连接"),
		field.Int("sort").Default(0).Comment("排序"),
	}
}

// Edges of the Advert.
func (Advert) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Advert) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Advert) Indexes() []ent.Index {
	return []ent.Index{}
}
