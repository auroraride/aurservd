package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type MaintainerMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m MaintainerMixin) Fields() []ent.Field {
	relate := field.Uint64("maintainer_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m MaintainerMixin) Edges() []ent.Edge {
	e := edge.To("maintainer", Maintainer.Type).Unique().Field("maintainer_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m MaintainerMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("maintainer_id"))
	}
	return
}

// Maintainer holds the schema definition for the Maintainer entity.
type Maintainer struct {
	ent.Schema
}

// Annotations of the Maintainer.
func (Maintainer) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "maintainer"},
		entsql.WithComments(true),
	}
}

// Fields of the Maintainer.
func (Maintainer) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enable").Comment("是否启用"),
		field.String("name").Comment("姓名"),
		field.String("phone").Unique().Comment("电话"),
		field.String("password").Comment("密码"),
	}
}

// Edges of the Maintainer.
func (Maintainer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("cities", City.Type).Ref("maintainers"),
	}
}

func (Maintainer) Mixin() []ent.Mixin {
	return []ent.Mixin{}
}

func (Maintainer) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("phone").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
