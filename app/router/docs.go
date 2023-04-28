// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/assets/docs"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title                极光出行API
// @version              1.0
// @BasePath             /
// @description.markdown
// @doc https://github.com/swaggo/swag/issues/386 https://github.com/swaggo/swag/issues/548 https://github.com/go-openapi/runtime/blob/master/middleware/redoc.go
func loadDocRoutes() {
	g := e.Group("/docs")

	docs.SwaggerInfo.Host = ar.Config.App.Host

	g.GET("", func(c echo.Context) error {
		return c.HTML(200, assets.SwaggerRedocUI)
	})

	g.GET("/swagger.json", func(c echo.Context) error {
		return c.JSONBlob(200, assets.SwaggerSpec)
	})

	g.GET("/swagger/*", echoSwagger.WrapHandler)

	g.GET("/oai3.json", func(c echo.Context) (err error) {
		var doc2 openapi2.T
		if err = jsoniter.Unmarshal(assets.SwaggerSpec, &doc2); err != nil {
			return
		}
		doc, err := openapi2conv.ToV3(&doc2)
		b, _ := jsoniter.Marshal(doc)
		return c.Blob(200, "application/json", b)
	})

	g.GET("/api.paw", func(c echo.Context) error {
		return c.Blob(200, "application/octet-stream", assets.Paw)
	})

	g.GET("/docs/octicons.css", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "text/css; charset=utf-8", assets.OcticonsCss)
	})

	g.GET("/assets/*", echo.WrapHandler(http.StripPrefix("/docs/assets/", http.FileServer(gfmstyle.Assets))))

	g.GET("/changelog/manager", func(c echo.Context) error {
		renderChangelog(c, "后台接口更新日志", assets.ChangelogManager)
		return nil
	})

	g.GET("/changelog/rider", func(c echo.Context) error {
		renderChangelog(c, "骑手接口更新日志", assets.ChangelogRider)
		return nil
	})

	g.GET("/changelog/employee", func(c echo.Context) error {
		renderChangelog(c, "店员接口更新日志", assets.ChangelogEmployee)
		return nil
	})
}

func renderChangelog(c echo.Context, title string, b []byte) {
	w := c.Response()
	b = github_flavored_markdown.Markdown(b)
	b = bytes.ReplaceAll(b, []byte("http://localhost:5533"), []byte(""))
	b = bytes.ReplaceAll(b, []byte("<a "), []byte(`<a target="_blank" `))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`<html><head><title>%s</title><meta charset="utf-8"><link href="/docs/assets/gfm.css" media="all" rel="stylesheet" type="text/css" /><link href="/docs/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`, title)))
	_, _ = w.Write(b)
	_, _ = w.Write([]byte(`</article></body></html>`))
}
