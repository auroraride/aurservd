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

	"github.com/auroraride/aurservd/internal/ent/internal"
)

type AssetManagerMixin struct {
	mixin.Schema
	DisableIndex bool
	Optional     bool
}

func (m AssetManagerMixin) Fields() []ent.Field {
	f := field.Uint64("asset_manager_id").Comment("仓库管理人ID")
	if m.Optional {
		f.Optional().Nillable()
	}
	return []ent.Field{f}
}

func (m AssetManagerMixin) Edges() []ent.Edge {
	e := edge.To("asset_manager", AssetManager.Type).Unique().Field("asset_manager_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m AssetManagerMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("asset_manager_id"))
	}
	return
}

// AssetManager holds the schema definition for the AssetManager entity.
type AssetManager struct {
	ent.Schema
}

// Annotations of the AssetManager.
func (AssetManager) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_manager"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetManager.
func (AssetManager) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(30).Comment("姓名"),
		field.String("phone").MaxLen(30).Comment("账户/手机号"),
		field.String("password").Comment("密码"),
		field.Uint64("role_id").Optional().Nillable().Comment("角色ID"),
		field.Bool("mini_enable").Default(false).Comment("仓管小程序人员是否启用"),
		field.Uint("mini_limit").Default(0).Comment("仓管小程序人员限制范围(m)"),
		field.Time("last_signin_at").Nillable().Optional().Comment("最后登录时间"),
	}
}

// Edges of the AssetManager.
func (AssetManager) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("role", AssetRole.Type).Unique().Ref("asset_managers").Field("role_id"),
		edge.From("warehouses", Warehouse.Type).Ref("asset_managers"),
		edge.To("warehouse", Warehouse.Type).Unique(),
	}
}

func (AssetManager) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (AssetManager) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id"),
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
