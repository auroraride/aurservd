package schema

import (
    "context"
    "entgo.io/ent"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/entc/integration/ent/hook"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/edge"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"
    "entgo.io/ent/schema/mixin"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/internal"
    "time"
)

type CabinetMixin struct {
    mixin.Schema
}

func (CabinetMixin) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("cabinet_id").Optional(),
    }
}

func (CabinetMixin) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("cabinet", Cabinet.Type).Unique().Field("cabinet_id"),
    }
}

// Cabinet holds the schema definition for the Cabinet entity.
type Cabinet struct {
    ent.Schema
}

// Annotations of the Cabinet.
func (Cabinet) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entsql.Annotation{Table: "cabinet"},
    }
}

// Fields of the Cabinet.
func (Cabinet) Fields() []ent.Field {
    return []ent.Field{
        field.Uint64("branch_id").Optional().Nillable().Comment("网点"),
        field.String("sn").Unique().Comment("编号"),
        field.String("brand").Comment("品牌"),
        field.String("serial").Comment("原始编号"),
        field.String("name").Comment("名称"),
        field.Uint("doors").Comment("柜门数量"),
        field.Uint8("status").Comment("投放状态"),
        field.Uint8("health").Default(0).Comment("健康状态 0未知 1正常 2离线 3故障"),
        field.JSON("bin", []model.CabinetBin{}).Optional().Comment("仓位信息"),
        field.Uint("battery_num").Default(0).Comment("电池总数"),
        field.Uint("battery_full_num").Default(0).Comment("满电电池数"),
        field.Float("lng").Optional().Comment("经度"),
        field.Float("lat").Optional().Comment("纬度"),
        field.String("address").Optional().Comment("详细地址"),
    }
}

// Edges of the Cabinet.
func (Cabinet) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("branch", Branch.Type).
            Ref("cabinets").
            Unique().
            Field("branch_id"),
        edge.To("bms", BatteryModel.Type),
        edge.To("faults", CabinetFault.Type),
        edge.To("exchanges", Exchange.Type),
    }
}

func (Cabinet) Mixin() []ent.Mixin {
    return []ent.Mixin{
        internal.TimeMixin{},
        internal.DeleteMixin{},
        internal.Modifier{},

        CityMixin{Optional: true},
    }
}

func (Cabinet) Indexes() []ent.Index {
    return []ent.Index{
        index.Fields("branch_id"),
        index.Fields("brand"),
        index.Fields("serial").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
        index.Fields("name").Annotations(
            entsql.IndexTypes(map[string]string{
                dialect.Postgres: "GIN",
            }),
        ),
    }
}

func (Cabinet) Hooks() []ent.Hook {
    type intr interface {
        Health() (r uint8, exists bool)
        OldHealth(ctx context.Context) (v uint8, err error)
        OldBrand(ctx context.Context) (v string, err error)
        OldSerial(ctx context.Context) (v string, err error)
        UpdatedAt() (r time.Time, exists bool)
    }

    return []ent.Hook{
        hook.On(
            func(next ent.Mutator) ent.Mutator {
                return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
                    if mt, ok := m.(intr); ok {
                        // 监听状态变化
                        if from, err := mt.OldHealth(ctx); err == nil {
                            if to, exists := mt.Health(); exists && from != to {
                                fmt.Println("状态发生了变化:", from, to)
                                b, _ := mt.OldBrand(ctx)
                                s, _ := mt.OldSerial(ctx)
                                u, _ := mt.UpdatedAt()
                                logging.NewHealthLog(b, s, u).SetStatus(from, to).Send()
                            }
                        }
                    }
                    return next.Mutate(ctx, m)
                })
            },
            ent.OpUpdate|ent.OpUpdateOne,
        ),
    }
}
