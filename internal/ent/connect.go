// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-30
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "context"
    "database/sql"
    "entgo.io/ent/dialect"
    entsql "entgo.io/ent/dialect/sql"
    "entgo.io/ent/dialect/sql/schema"
    "fmt"
    "github.com/auroraride/aurservd/internal/ent/migrate"
    log "github.com/sirupsen/logrus"

    _ "github.com/auroraride/aurservd/internal/ent/runtime"
    _ "github.com/jackc/pgx/v4/stdlib"
)

var Database *Client

func OpenDatabase(dsn string, debug bool) *Client {
    pgx, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatalf("数据库打开失败: %v", err)
    }

    // 从db变量中构造一个ent.Driver对象。
    drv := entsql.OpenDB(dialect.Postgres, pgx)
    c := NewClient(Driver(drv))
    if debug {
        c = c.Debug()
    }

    autoMigrate(c)

    return c
}

func autoMigrate(c *Client) {
    ctx := context.Background()
    if err := c.Schema.Create(
        ctx,
        schema.WithAtlas(true),
        migrate.WithDropIndex(true),
        migrate.WithDropColumn(true),
        migrate.WithGlobalUniqueID(true),
        // migrate.WithForeignKeys(false),
    ); err != nil {
        log.Fatalf("数据库迁移失败: %v", err)
    }
}

func WithTx(ctx context.Context, fn func(tx *Tx) error) error {
    tx, err := Database.Tx(ctx)
    if err != nil {
        return err
    }
    defer func() {
        if v := recover(); v != nil {
            tx.Rollback()
            panic(v)
        }
    }()
    if err = fn(tx); err != nil {
        if rerr := tx.Rollback(); rerr != nil {
            err = fmt.Errorf("rolling back transaction: %w", rerr)
        }
        return err
    }
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("committing transaction: %w", err)
    }
    return nil
}
