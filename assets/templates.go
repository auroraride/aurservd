// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-01
// Based on aurservd by liasica, magicrolan@qq.com.

package assets

import (
    "embed"
    "github.com/labstack/echo/v4"
    "html/template"
    "io"
    "io/fs"
)

const (
    templatesDir = "views"
)

var (
    //go:embed views/*
    views embed.FS

    Templates *htmlTemplate
)

type htmlTemplate struct {
    Templates map[string]*template.Template
}

func (t *htmlTemplate) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
    return t.Templates[name].ExecuteTemplate(w, name, data)
}

func LoadTemplates() {
    Templates = &htmlTemplate{Templates: make(map[string]*template.Template)}
    tmplFiles, err := fs.ReadDir(views, templatesDir)
    if err != nil {
        return
    }

    for _, tmpl := range tmplFiles {
        if tmpl.IsDir() {
            continue
        }
        var pt *template.Template
        pt, err = template.ParseFS(views, templatesDir+"/"+tmpl.Name())
        if err != nil {
            continue
        }

        Templates.Templates[tmpl.Name()] = pt
    }
}
