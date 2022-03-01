// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package ar

import (
    atlasMigrate "ariga.io/atlas/sql/migrate"
    "context"
    "database/sql"
    "entgo.io/ent/dialect"
    entsql "entgo.io/ent/dialect/sql"
    "entgo.io/ent/dialect/sql/schema"
    "github.com/auroraride/aurservd/internal/ar/db"
    _ "github.com/auroraride/aurservd/internal/ent/runtime"
    log "github.com/sirupsen/logrus"

    "github.com/auroraride/aurservd/internal/ent"

    "github.com/auroraride/aurservd/internal/ent/migrate"
    _ "github.com/jackc/pgx/v4/stdlib"
)

var Ent *ent.Client

type orm struct {
    *ent.Client
}

func OpenDatabase() *orm {
    Ent = openPgx(Config.Database.Postgres.Dsn)
    o := &orm{Ent}
    o.autoMigrate()
    return o
}

func openPgx(dsn string) (c *ent.Client) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatalf("数据库打开失败: %v", err)
    }

    // 从db变量中构造一个ent.Driver对象。
    drv := entsql.OpenDB(dialect.Postgres, db)
    return ent.NewClient(ent.Driver(drv))
}

func (c *orm) autoMigrate() {
    ctx := context.Background()
    o := new(db.Options)
    if err := c.Schema.Create(
        ctx,
        schema.WithAtlas(true),
        migrate.WithDropIndex(true),
        migrate.WithDropColumn(true),
        migrate.WithGlobalUniqueID(true),
        schema.WithApplyHook(func(next schema.Applier) schema.Applier {
            return schema.ApplyFunc(func(ctx context.Context, conn dialect.ExecQuerier, plan *atlasMigrate.Plan) error {
                for _, change := range plan.Changes {
                    if change.Comment == `create "city" table` {
                        o.City = true
                    }
                    if change.Comment == `create "manager" table` {
                        o.Manager = true
                    }
                }
                return next.Apply(ctx, conn, plan)
            })
        }),
    ); err != nil {
        log.Fatalf("数据库迁移失败: %v", err)
    }

    // 初始化数据库
    db.Database(c.Client, o)
}
