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

type CabinetEcMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m CabinetEcMixin) Fields() []ent.Field {
	relate := field.Uint64("ec_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m CabinetEcMixin) Edges() []ent.Edge {
	e := edge.To("ec", CabinetEc.Type).Unique().Field("ec_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m CabinetEcMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("ec_id"))
	}
	return
}

// CabinetEc holds the schema definition for the CabinetEc entity.
type CabinetEc struct {
	ent.Schema
}

// Annotations of the CabinetEc.
func (CabinetEc) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "cabinet_ec"},
		entsql.WithComments(true),
	}
}

// Fields of the CabinetEc.
func (CabinetEc) Fields() []ent.Field {
	return []ent.Field{
		field.String("serial").Comment("电柜原始编号"),
		field.Time("date").Comment("日期"),
		field.Float("start").Comment("开始电量"),
		field.Float("end").Optional().Comment("结束电量"),
		field.Float("total").Comment("耗电量"),
	}
}

// Edges of the CabinetEc.
func (CabinetEc) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (CabinetEc) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
	}
}

func (CabinetEc) Indexes() []ent.Index {
	return []ent.Index{}
}
