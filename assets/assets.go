// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package assets

import (
	"embed"
)

var (
	//go:embed city.json
	City []byte

	//go:embed docs/*
	SwaggerSpecs embed.FS

	//go:embed changelog/manager.md
	ChangelogManager []byte

	//go:embed changelog/rider.md
	ChangelogRider []byte

	//go:embed changelog/employee.md
	ChangelogEmployee []byte

	//go:embed octicons.css
	OcticonsCss []byte

	//go:embed sql/stock_overview.sql
	SQLStockOverview string

	//go:embed views/legal.go.html
	LegalTemplate string
)
