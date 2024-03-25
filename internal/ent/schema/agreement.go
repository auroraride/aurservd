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

type AgreementMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AgreementMixin) Fields() []ent.Field {
	relate := field.Uint64("agreement_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AgreementMixin) Edges() []ent.Edge {
	e := edge.To("agreement", Agreement.Type).Unique().Field("agreement_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AgreementMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("agreement_id"))
	}
	return
}

// Agreement holds the schema definition for the Agreement entity.
type Agreement struct {
	ent.Schema
}

// Annotations of the Agreement.
func (Agreement) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "agreement"},
		entsql.WithComments(true),
	}
}

// Fields of the Agreement.
func (Agreement) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("协议名称"),
		field.Text("content").Comment("协议内容"),
		field.Uint8("user_type").Default(1).Comment("用户类型 1:个签 2:团签"),
		field.Uint8("force_read_time").Default(0).Comment("强制阅读时间"),
		field.Bool("is_default").Default(false).Comment("是否为默认协议"),
		field.String("hash").Optional().Comment("hash"),
		field.String("url").Optional().Comment("url"),
	}
}

// Edges of the Agreement.
func (Agreement) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Agreement) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Agreement) Indexes() []ent.Index {
	return []ent.Index{}
}
