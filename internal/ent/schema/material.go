package schema

import (
	"fmt"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type MaterialMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
	Prefix       string
}

func (m MaterialMixin) prefield() (string, string) {
	if m.Prefix == "" {
		return "material_id", "material"
	}
	return fmt.Sprintf("%s_material_id", m.Prefix), fmt.Sprintf("%sMaterial", m.Prefix)
}

func (m MaterialMixin) Fields() []ent.Field {
	pf, _ := m.prefield()
	f := field.Uint64(pf).Comment("物资ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m MaterialMixin) Edges() []ent.Edge {
	pf, pn := m.prefield()
	e := edge.To(pn, Material.Type).Unique().Field(pf)
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m MaterialMixin) Indexes() (arr []ent.Index) {
	pf, _ := m.prefield()
	if !m.DisableIndex {
		arr = append(arr, index.Fields(pf))
	}
	return
}

// Material holds the schema definition for the Material entity.
type Material struct {
	ent.Schema
}

// Annotations of the Material.
func (Material) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "material"},
		entsql.WithComments(true),
	}
}

// Fields of the Material.
func (Material) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("物资名称"),
		field.Uint8("type").Comment("物资类型 1电柜配件 2电车配件 3.其他"),
		field.String("statement").Optional().Comment("说明"),
		field.Bool("allot").Comment("是否可调拨"),
	}
}

// Edges of the Material.
func (Material) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Material) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Material) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}
