package schema

import (
	"ariga.io/atlas/sql/postgres"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/app/model"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type VersionMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m VersionMixin) Fields() []ent.Field {
	relate := field.Uint64("version_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m VersionMixin) Edges() []ent.Edge {
	e := edge.To("version", Version.Type).Unique().Field("version_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m VersionMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("version_id"))
	}
	return
}

// Version holds the schema definition for the Version entity.
type Version struct {
	ent.Schema
}

// Annotations of the Version.
func (Version) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "version"},
		entsql.WithComments(true),
	}
}

// Fields of the Version.
func (Version) Fields() []ent.Field {
	return []ent.Field{
		field.Other("platform", model.AppPlatform("")).Default(model.AppPlatform("")).SchemaType(map[string]string{
			dialect.Postgres: postgres.TypeCharVar,
		}).Comment("平台"),
		field.String("version").Unique().Comment("版本号"),
		field.String("content").Comment("更新内容"),
		field.Bool("force").Comment("是否强制更新"),
	}
}

// Edges of the Version.
func (Version) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Version) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Version) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("platform", "version").Unique(),
	}
}
