// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "bytes"
    "errors"
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/docs"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
    "html/template"
    "io/ioutil"
    "net/http"
    "strings"
)

// ErrSpecNotFound error for when spec file not found
var ErrSpecNotFound = errors.New("spec not found")

// Redoc configuration
type Redoc struct {
    DocsPath    string
    SpecUrl     string
    SpecFile    string
    Title       string
    Description string
}

// Body returns the final html with the js in the body
func (r Redoc) body() ([]byte, error) {
    buf := bytes.NewBuffer(nil)
    tpl, err := template.New("redoc").Parse(assets.SwaggerRedocUI)
    if err != nil {
        return nil, err
    }

    if err = tpl.Execute(buf, map[string]string{
        "title":       r.Title,
        "url":         r.SpecUrl,
        "description": r.Description,
    }); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

// handler sets some defaults and returns a HandlerFunc
func (r Redoc) handler() http.HandlerFunc {
    data, err := r.body()
    if err != nil {
        panic(err)
    }

    specFile := r.SpecFile
    if specFile == "" {
        panic(ErrSpecNotFound)
    }

    spec, err := ioutil.ReadFile(specFile)
    if err != nil {
        panic(err)
    }

    docsPath := r.DocsPath
    return func(w http.ResponseWriter, req *http.Request) {

        if strings.HasSuffix(req.URL.Path, r.SpecUrl) {
            w.WriteHeader(200)
            w.Header().Set("content-type", "application/json")
            _, _ = w.Write(spec)
            return
        }

        if docsPath == "" || docsPath == req.URL.Path {
            w.WriteHeader(200)
            w.Header().Set("content-type", "text/html")
            _, _ = w.Write(data)
        }
    }
}

// @title 极光出行API
// @version 1.0
// @description.markdown 极光出行所有API接口文档
// @basePath /
// @doc https://github.com/swaggo/swag/issues/386 https://github.com/swaggo/swag/issues/548 https://github.com/go-openapi/runtime/blob/master/middleware/redoc.go
func newRedoc() echo.MiddlewareFunc {
    docs.SwaggerInfo.Host = ar.Config.App.Host

    apiDoc := Redoc{
        Title:       "极光出行API",
        Description: "极光出行所有API接口文档",
        SpecFile:    "./docs/swagger.json",
        SpecUrl:     "/docs/swagger.json",
        DocsPath:    "/docs",
    }

    handle := apiDoc.handler()

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(ctx echo.Context) error {
            handle(ctx.Response(), ctx.Request())
            return nil
        }
    }
}
