// Copyright (C) liasica. 2023-present.
//
// Created at 2023-07-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
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

	data := strings.Replace(assets.LegalTemplate, "{{- .Title -}}", req.LegalName.Name.Title(), 1)
	data = strings.Replace(data, "{{- .Content -}}", req.Content, 1)

	_, err := f.Write([]byte(data))
	if err != nil {
		snag.Panic(err)
	}
}
