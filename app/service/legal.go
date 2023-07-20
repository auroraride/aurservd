// Copyright (C) liasica. 2023-present.
//
// Created at 2023-07-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"html/template"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/exp/slices"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/assets"
	"github.com/auroraride/aurservd/pkg/snag"
)

type legalService struct {
	*BaseService
}

func NewLegal(params ...any) *legalService {
	return &legalService{
		BaseService: newService(params...),
	}
}

func (s *legalService) open(name model.Legal) *os.File {
	if !slices.Contains(model.Legals, name) {
		snag.Panic("请求参数错误")
	}

	// 创建或读取文件
	f, err := os.OpenFile(name.Filepath(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		snag.Panic(err)
	}
	return f
}

func (s *legalService) Read(req *model.LegalName) model.LegalRes {
	f := s.open(req.Name)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		snag.Panic("读取失败")
	}

	var content string
	content, err = doc.Find(".container").Html()
	if err != nil {
		snag.Panic("读取失败")
	}

	return model.LegalRes{
		Title:   req.Name.Title(),
		Content: content,
		Url:     req.Name.Url(),
	}
}

func (s *legalService) Save(req *model.LegalSaveReq) {
	f := s.open(req.LegalName.Name)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	t, err := template.New("layout").ParseFS(assets.LegalTemplateFS, "views/legal.go.html")
	if err != nil {
		snag.Panic(err)
	}

	err = t.ExecuteTemplate(f, "legal.go.html", map[string]any{
		"name":    req.Name.Title(),
		"content": template.HTML(strings.Join(strings.Split(req.Content, "\n"), "")),
	})
	if err != nil {
		snag.Panic(err)
	}
}
