// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-02, by liasica

package tools

import (
	"os"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	jsoniter "github.com/json-iterator/go"
)

type SorterTest struct {
	Price     float64 `json:"price"`
	BrandName string  `json:"brandName"`
	Rto       bool    `json:"rto"`
	Daily     bool    `json:"daily"`
}

func TestSorter(t *testing.T) {
	var items []SorterTest

	_ = jsoniter.Unmarshal([]byte(`[{"price":300,"brandName":"时光1代","rto":false,"daily":false},{"price":360,"brandName":"时光1代","rto":true,"daily":false},{"price":20,"brandName":"时光1代","rto":false,"daily":true},{"price":400,"brandName":"时光1代","rto":false,"daily":false},{"price":310,"brandName":"时光2代","rto":false,"daily":false},{"price":370,"brandName":"时光2代","rto":true,"daily":false},{"price":21,"brandName":"时光2代","rto":false,"daily":true},{"price":410,"brandName":"时光2代","rto":false,"daily":false},{"price":320,"brandName":"时光3代","rto":false,"daily":false},{"price":380,"brandName":"时光3代","rto":true,"daily":false},{"price":22,"brandName":"时光3代","rto":false,"daily":true},{"price":420,"brandName":"时光3代","rto":false,"daily":false}]`), &items)
	NewSorter().
		AddInt(func(i interface{}) int {
			v := i.(SorterTest)
			if v.Rto {
				return 3
			} else if v.Daily {
				return 2
			} else {
				return 1
			}
		}).
		AddStr(func(i interface{}) string { return i.(SorterTest).BrandName }).
		// AddInt(func(i interface{}) int {
		// 	runes := []rune(i.(SorterTest).BrandName)
		// 	fmt.Printf("%s -> %v\n", i.(SorterTest).BrandName, runes)
		// 	var sum int
		// 	for _, r := range runes {
		// 		sum += int(r)
		// 	}
		// 	return sum
		// }).
		AddFloat(func(i interface{}) float64 { return i.(SorterTest).Price }).
		SortStable(items)

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{"#", "类型", "车型", "价格"})

	var rows []table.Row

	for i, item := range items {
		typ := "月租"
		if item.Rto {
			typ = "以租代购"
		} else if item.Daily {
			typ = "日租"
		}

		rows = append(rows, table.Row{i + 1, typ, item.BrandName, item.Price})
	}
	tw.AppendRows(rows)
	tw.AppendSeparator()
	tw.Render()
}
