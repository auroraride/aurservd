// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-07
// Based on aurservd by liasica, magicrolan@qq.com.

package boot

import (
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
)

func entInit() {
	ent.Database = ent.OpenDatabase(ar.Config.Database.Postgres.Dsn, ar.Config.App.SQL)
	ent.Database.Cabinet.Use(service.NewCabinet().EntHooks()...)
	// 使用stock hook
	// ent.Database.Stock.Use(hook.On(ent.NewStockHook(), ent.OpCreate))
}
