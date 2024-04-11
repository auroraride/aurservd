package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
)

type ContractTemplateMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m ContractTemplateMixin) Fields() []ent.Field {
	relate := field.Uint64("template_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m ContractTemplateMixin) Edges() []ent.Edge {
	e := edge.To("template", ContractTemplate.Type).Unique().Field("template_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m ContractTemplateMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("template_id"))
	}
	return
}

// ContractTemplate holds the schema definition for the ContractTemplate entity.
type ContractTemplate struct {
	ent.Schema
}

// Annotations of the ContractTemplate.
func (ContractTemplate) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "contract_template"},
		entsql.WithComments(true),
	}
}

// Fields of the ContractTemplate.
func (ContractTemplate) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("模板名称"),
		field.String("url").Comment("模板文件地址"),
		field.Uint8("aimed").Default(definition.ContractTemplateAimedPersonal.Value()).Comment("用户类型 1:个签 2:团签"),
		field.Uint8("plan_type").Default(model.PlanTypeBattery.Value()).Comment("套餐类型 1:单电 2:车电"),
		field.String("hash").Comment("模板hash"),
		field.Bool("enable").Default(false).Comment("是否启用"),
	}
}

// Edges of the ContractTemplate.
func (ContractTemplate) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (ContractTemplate) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (ContractTemplate) Indexes() []ent.Index {
	return []ent.Index{}
}
