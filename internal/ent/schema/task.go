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
    "github.com/rs/xid"
)

type TaskMixin struct {
    mixin.Schema
    Optional     bool
    DisableIndex bool
}

func (m TaskMixin) Fields() []ent.Field {
    relate := field.Uint64("task_id")
    if m.Optional {
        relate.Optional().Nillable()
    }
    return []ent.Field{
        relate,
    }
}

func (m TaskMixin) Edges() []ent.Edge {
    e := edge.To("task", Task.Type).Unique().Field("task_id")
    if !m.Optional {
        e.Required()
    }
    return []ent.Edge{e}
}

func (m TaskMixin) Indexes() (arr []ent.Index) {
    if !m.DisableIndex {
        arr = append(arr, index.Fields("task_id"))
    }
    return
}

// Task holds the schema definition for the Task entity.
type Task struct {
    ent.Schema
}

// Annotations of the Task.
func (Task) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "task"},
        entsql.WithComments(true),
    }
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
    return []ent.Field{
        field.Other("uuid", xid.ID{}).Default(xid.New).SchemaType(map[string]string{
            dialect.Postgres: postgres.TypeBytea,
        }).Immutable(),
        field.Uint64("exchange_id").Optional().Nillable(),
        field.Uint64("cabinet_id").Optional().Nillable(),
        field.String("serial").Comment("电柜编码"),
        field.Enum("job").GoType(model.TaskJobExchange).Comment("任务类别"),
        field.Other("status", model.TaskStatusNotStart).Default(model.TaskStatusNotStart).SchemaType(map[string]string{
            dialect.Postgres: postgres.TypeSmallInt,
        }).Comment("任务状态"),
        field.Time("start_at").Optional().Nillable().Comment("开始时间"),
        field.Time("stop_at").Optional().Nillable().Comment("结束时间"),
        field.String("message").Optional().Comment("失败消息"),
        field.JSON("exchange", &model.ExchangeTaskInfo{}).Optional().Comment("换电信息"),
        field.JSON("business_bin_info", &model.BinInfo{}).Optional().Comment("仓位信息"),
        field.JSON("cabinet", &model.ExchangeTaskCabinet{}).Comment("电柜信息"),
    }
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
    return []ent.Edge{}
}

func (Task) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},

        // CabinetMixin{},
        RiderMixin{},
        // ExchangeMixin{Optional: true},
    }
}

func (Task) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("serial"),
        index.Fields("status"),
        index.Fields("start_at"),
    }
}
