// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/10
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "context"
    "database/sql"
    "entgo.io/ent/dialect"
    entsql "entgo.io/ent/dialect/sql"
    "github.com/auroraride/aurservd/internal/ent/migrate"
    _ "github.com/jackc/pgx/v4/stdlib"
    log "github.com/sirupsen/logrus"
)

func OpenPgx(dsn string) (c *Client) {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatalf("数据库打开失败: %v", err)
    }

    // 从db变量中构造一个ent.Driver对象。
    drv := entsql.OpenDB(dialect.Postgres, db)
    return NewClient(Driver(drv))
}

func (c *Client) AutoMigrate() {
    ctx := context.Background()
    if err := c.Debug().Schema.Create(
        ctx,
        migrate.WithDropIndex(true),
        migrate.WithDropColumn(true),
        migrate.WithGlobalUniqueID(true),
    ); err != nil {
        log.Fatalf("数据库迁移失败: %v", err)
    }
}

func (c *Client) DB() *sql.DB {
    return c.driver.(*entsql.Driver).DB()
}
