// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package assets

import (
    _ "embed"
)

var (
    //go:embed city.json
    City []byte

    //go:embed swagger.redoc.html
    SwaggerRedocUI string

    //go:embed docs/swagger.json
    SwaggerSpec []byte

    //go:embed docs/swagger.yaml
    SwaggerSpecYaml []byte

    //go:embed api.paw
    Paw []byte

    //go:embed changelog/manager.md
    ChangelogManager []byte

    //go:embed changelog/rider.md
    ChangelogRider []byte

    //go:embed octicons.css
    OcticonsCss []byte
)
