package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/internal"
	"github.com/auroraride/aurservd/pkg/snag"
)

// Stock holds the schema definition for the Stock entity.
type Stock struct {
	ent.Schema
}

// Annotations of the Stock.
func (Stock) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "stock"},
		entsql.WithComments(true),
	}
}

// Fields of the Stock.
func (Stock) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("parent_id").Optional().Nillable().Comment("父级"),
		field.String("sn").Comment("调拨编号"),
		field.Uint8("type").Default(0).Comment("类型 0:调拨 1:骑手激活 2:骑手寄存 3:骑手结束寄存 4:骑手退租"),
		field.Uint64("store_id").Optional().Nillable().Comment("入库至 或 出库自 门店ID"),
		field.Uint64("cabinet_id").Optional().Nillable().Comment("入库至 或 出库自 电柜ID"),
		field.Uint64("rider_id").Optional().Nillable().Comment("对应骑手ID"),
		field.Uint64("employee_id").Optional().Nillable().Comment("操作店员ID"),
		field.Uint64("enterprise_id").Optional().Nillable().Comment("团签ID"),
		field.Uint64("station_id").Optional().Nillable().Comment("站点ID"),
		field.String("name").Comment("物资名称"),
		field.String("model").Optional().Nillable().Comment("电池型号"),
		field.Int("num").Immutable().Comment("物资数量: 正值调入 / 负值调出"),
		field.Enum("material").Values("battery", "ebike", "others").Comment("物资种类"),
	}
}

// Edges of the Stock.
func (Stock) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).Unique().Ref("stocks").Field("store_id"),
		edge.From("cabinet", Cabinet.Type).Unique().Ref("stocks").Field("cabinet_id"),
		edge.From("rider", Rider.Type).Unique().Ref("stocks").Field("rider_id"),
		edge.From("employee", Employee.Type).Unique().Ref("stocks").Field("employee_id"),

		edge.To("spouse", Stock.Type).Unique(),

		edge.To("children", Stock.Type).From("parent").Field("parent_id").Unique(),
		edge.From("enterprise", Enterprise.Type).Unique().Ref("stocks").Field("enterprise_id"),
		edge.From("station", EnterpriseStation.Type).Unique().Ref("stocks").Field("station_id"),
	}
}

func (Stock) Mixin() []ent.Mixin {
	return []ent.Mixin{
		internal.TimeMixin{},
		internal.DeleteMixin{},
		internal.Modifier{
			IndexCreator: true,
		},
		CityMixin{Optional: true},
		SubscribeMixin{Optional: true},

		// 电车
		EbikeMixin{Optional: true},
		EbikeBrandMixin{Optional: true},

		// 电池
		BatteryMixin{Optional: true},

		// 代理商
		AgentMixin{Optional: true},
	}
}

func (Stock) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("store_id"),
		index.Fields("cabinet_id"),
		index.Fields("rider_id"),
		index.Fields("employee_id"),
		index.Fields("model"),
		index.Fields("sn"),
		index.Fields("parent_id"),
		index.Fields("name").Annotations(
			entsql.IndexTypes(map[string]string{
				dialect.Postgres: "GIN",
			}),
			entsql.OpClass("gin_trgm_ops"),
		),
	}
}

func (Stock) Hooks() []ent.Hook {
	type intr interface {
		Num() (r int, exists bool)
		StoreID() (r uint64, exists bool)
		GetType() (r uint8, exists bool)
	}

	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if m.Op().Is(ent.OpCreate) {
					if st, ok := m.(intr); ok {
						n, _ := st.Num()
						if t, ok := st.GetType(); ok {
							// 骑手业务是否满足调拨电池数量要求
							x := model.StockNumberOfRiderBusiness(t)
							if x != 0 && n != x {
								snag.Panic("电池数量错误")
							}
						}
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
