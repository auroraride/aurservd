// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/internal/ar"
)

func loadDocRoutes() {
	g := e.Group("/docs")

	items := map[string]string{
		"/agent/v1":     "代理接口 ver1",
		"/employee/v1":  "店员接口 ver1",
		"/manager/v1":   "管理接口 ver1",
		"/operator/v1":  "运维接口 ver1",
		"/promotion/v1": "推广接口 ver1",
		"/rider/v1":     "骑手接口 ver1",
		"/rider/v2":     "骑手接口 ver2",
	}

	g.GET("", func(c echo.Context) error {
		return c.Render(http.StatusOK, "docs.html", ar.Map{"items": items})
	})

	for path, name := range items {
		k := path
		v := name
		spec, _ := assets.SwaggerSpecs.ReadFile("docs" + k + "/swagger.json")

		g.GET(k, func(c echo.Context) error {
			// return c.HTML(200, assets.SwaggerRedocUI)
			return c.Render(http.StatusOK, "swagger.redoc.html", ar.Map{
				"path": template.URL("/docs" + k + "/oai3.json"),
				"name": v,
			})
		})

		g.GET(k+"/swagger.json", func(c echo.Context) error {
			return c.JSONBlob(200, spec)
		})

		g.GET(k+"/oai3.json", func(c echo.Context) (err error) {
			var doc2 openapi2.T
			if err = jsoniter.Unmarshal(spec, &doc2); err != nil {
				return
			}
			doc, _ := openapi2conv.ToV3(&doc2)
			b, _ := json.MarshalIndent(doc, "", "  ")
			return c.Blob(200, "application/json", b)
		})
	}

	g.GET("/swagger/*", echoSwagger.WrapHandler)

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
