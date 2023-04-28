// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-01
// Based on aurservd by liasica, magicrolan@qq.com.

package tool

import (
	"context"
	"net/http"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/exchange"
	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
)

type cabinet struct{}

var Cabinet = new(cabinet)

type exchangeResult struct {
	Full        int
	Alternative int
}

func (*cabinet) Exchange(c echo.Context) (err error) {
	ctx := app.Context(c)
	date := c.QueryParam("date")
	var d carbon.Carbon
	if date != "" {
		d = carbon.Parse(date)
	}

	if d.IsZero() {
		d = carbon.Yesterday().StartOfDay()
	}

	items, _ := ent.Database.Cabinet.QueryNotDeleted().WithExchanges(func(eq *ent.ExchangeQuery) {
		eq.Where(exchange.CreatedAtGTE(d.Carbon2Time()), exchange.CreatedAtLTE(d.EndOfDay().Carbon2Time()))
	}).WithCity().All(context.Background())

	out := make(map[string]map[string]*exchangeResult)
	cities := make(map[string]*exchangeResult)

	for _, item := range items {
		ci := item.Edges.City
		if ci == nil {
			continue
		}

		if out[ci.Name] == nil {
			out[ci.Name] = make(map[string]*exchangeResult)
			cities[ci.Name] = &exchangeResult{}
		}

		co := out[ci.Name]
		if co[item.Name] == nil {
			co[item.Name] = &exchangeResult{}
		}

		r := co[item.Name]
		cr := cities[ci.Name]

		for _, e := range item.Edges.Exchanges {
			if e.Alternative {
				r.Alternative += 1
				cr.Alternative += 1
			} else {
				r.Full += 1
				cr.Full += 1
			}
		}
	}

	// TODO 导出 https://www.revisitclass.com/css/how-to-export-download-the-html-table-to-excel-using-javascript/
	return ctx.Render(http.StatusOK, "exchange.html", ar.Map{
		"date":   d.Format("Y-m-d"),
		"items":  out,
		"cities": cities,
	})
}
