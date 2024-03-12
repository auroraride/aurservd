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

type PushmessageMixin struct {
	mixin.Schema
	Optional     bool
	DisableIndex bool
}

func (m PushmessageMixin) Fields() []ent.Field {
	relate := field.Uint64("pushmessage_id")
	if m.Optional {
		relate.Optional().Nillable()
	}
	return []ent.Field{
		relate,
	}
}

func (m PushmessageMixin) Edges() []ent.Edge {
	e := edge.To("pushmessage", Pushmessage.Type).Unique().Field("pushmessage_id")
	if !m.Optional {
		e.Required()
	}
	return []ent.Edge{e}
}

func (m PushmessageMixin) Indexes() (arr []ent.Index) {
	if !m.DisableIndex {
		arr = append(arr, index.Fields("pushmessage_id"))
	}
	return
}

// Pushmessage holds the schema definition for the Pushmessage entity.
type Pushmessage struct {
	ent.Schema
}

// Annotations of the Pushmessage.
func (Pushmessage) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "pushmessage"},
		entsql.WithComments(true),
	}
}

// Fields of the Pushmessage.
func (Pushmessage) Fields() []ent.Field {
	// 标题 封面图片 内容 推送类型 推送时间 是否首页推送 首页推送内容 消息状态 消息类型
	return []ent.Field{
		field.String("title").Comment("标题"),
		field.String("image").Comment("封面图片"),
		field.Text("content").Comment("内容"),
		field.Uint8("push_type").Comment("推送类型"),
		field.Time("push_time").Optional().Nillable().Comment("推送时间"),
		field.Bool("is_home").Comment("是否首页推送"),
		field.String("home_content").Comment("首页推送内容"),
		field.Uint8("message_status").Comment("消息状态"),
		field.Uint8("message_type").Comment("消息类型"),
		field.String("third_party_id").Comment("第三方推送平台消息ID"),
	}
}

// Edges of the Pushmessage.
func (Pushmessage) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Pushmessage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{},
	}
}

func (Pushmessage) Indexes() []ent.Index {
	return []ent.Index{}
}
