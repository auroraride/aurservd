package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/auroraride/aurservd/internal/ent/internal"
)

// AssetExport holds the schema definition for the AssetExport entity.
type AssetExport struct {
	ent.Schema
}

// Annotations of the AssetExport.
func (AssetExport) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "asset_export"},
		entsql.WithComments(true),
	}
}

// Fields of the AssetExport.
func (AssetExport) Fields() []ent.Field {
	return []ent.Field{
		field.String("taxonomy").Comment("分类"),
		field.String("sn").Comment("编号"),
		field.Uint8("status").Default(0).Comment("状态"),
		field.String("path").Optional().Comment("文件路径"),
		field.String("message").Optional().Comment("失败原因"),
		field.Time("finish_at").Optional().Comment("生成时间"),
		field.Int64("duration").Optional().Comment("耗时"),
		field.Text("condition").Comment("筛选条件"),
		field.JSON("info", map[string]interface{}{}).Optional().Comment("详细信息"),
		field.String("remark").Comment("备注信息"),
	}
}

// Edges of the AssetExport.
func (AssetExport) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (AssetExport) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		AssetManagerMixin{},
	}
}

func (AssetExport) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("sn"),
		index.Fields("status"),
	}
}
