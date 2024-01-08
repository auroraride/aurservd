// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-08
// Based on aurservd by liasica, magicrolan@qq.com.

package main

import (
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
)

func main() {
	// 获取读取目录
	target := "app/controller/v1/rapi"
	pkgs, _ := parser.ParseDir(token.NewFileSet(), target, nil, parser.ParseComments)
	for _, p := range pkgs {
		dp := doc.New(p, "./", doc.AllDecls)
		fmt.Println(dp)
		for _, fvalue := range p.Files {
			for _, comment := range fvalue.Comments {
				if strings.Contains(comment.Text(), "Copyright (C)") {
					continue
				}
				for _, c := range comment.List {
					if strings.Contains(c.Text, "@Router") {
						fmt.Println(c.Text)
					}
				}
			}
		}
	}
}
