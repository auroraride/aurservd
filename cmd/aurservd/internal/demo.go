// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
	"context"
	"log"

	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
)

func Demo() {
	// TODO 修复需要删除
	log.Println("修复站点城市 (默认为企业城市ID)")
	ctx := context.Background()
	items, _ := ent.Database.EnterpriseStation.QueryNotDeleted().Where(enterprisestation.CityIDIsNil()).WithEnterprise().All(ctx)
	for _, item := range items {
		_ = item.Update().SetCityID(item.Edges.Enterprise.CityID).Exec(ctx)
	}
}
