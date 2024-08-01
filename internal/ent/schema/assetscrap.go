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

type AssetScrapMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m AssetScrapMixin) Fields() []ent.Field {
	relate := field.Uint64("scrap_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m AssetScrapMixin) Edges() []ent.Edge {
	e := edge.To("scrap", AssetScrap.Type).Unique().Field("scrap_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetScrapMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("scrap_id"))
	}
	return
}

// AssetScrap holds the schema definition for the AssetScrap entity.
type AssetScrap struct {
	ent.Schema
}

// Annotations of the AssetScrap.
func (AssetScrap) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_scrap"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetScrap.
func (AssetScrap) Fields() []ent.Field {
	return []ent.Field{
		field.Uint8("scrap_reason_type").Optional().Comment("报废原因 1:丢失 2:损坏 3:其他"),
		field.Time("scrap_at").Optional().Comment("报废时间"),
		field.Uint64("scrap_operate_id").Optional().Nillable().Comment("操作报废人员ID"),
		field.Uint8("scrap_operate_role_type").Optional().Nillable().Comment("报废人员角色类型 1:后台管理员 2:门店管理员 3:运维 4:物资管理员 5:代理管理员"),
		field.String("scrap_batch").Optional().Comment("报废批次"),
		field.Uint64("asset_id").Comment("资产ID"),
	}
}

// Edges of the AssetScrap.
func (AssetScrap) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联后台管理员
		edge.To("manager", Manager.Type).Unique().Field("scrap_operate_id"),
		// 关联门店管理员
		edge.To("employee", Employee.Type).Unique().Field("scrap_operate_id"),
		// 关联运维
		edge.To("maintainer", Maintainer.Type).Unique().Field("scrap_operate_id"),
		// 关联代理管理员
		edge.To("agent", Agent.Type).Unique().Field("scrap_operate_id"),

		edge.From("asset", Asset.Type).Ref("scrap").Unique().Field("asset_id").Required(),
	}
}

func (AssetScrap) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.Modifier{},
	}
}

func (AssetScrap) Indexes() []ent.Index {
	return []ent.Index{}
}
