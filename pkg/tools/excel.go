// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-07
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
	"fmt"
	"reflect"

	"github.com/xuri/excelize/v2"

	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/utils"
)

type excel struct {
	path  string
	sheet string

	row    int // 当前行
	column int // 当前列

	*excelize.File
}

type ExcelItems [][]any

func (e ExcelItems) Columns() int {
	if len(e) == 0 {
		return 0
	}
	return len(e[0])
}

func NewExcel(path string, args ...any) (e *excel) {
	err := utils.NewFile(path).CreateDirectoryIfNotExist()
	if err != nil {
		snag.Panic(err)
	}

	sheet := "Sheet1"
	if len(args) > 0 {
		sheet = args[0].(string)
	}

	fe := excelize.NewFile()
	_ = fe.SetSheetName("Sheet1", sheet)
	return &excel{
		path:  path,
		sheet: sheet,
		File:  fe,
	}
}

func (e *excel) CellString(row, column int) string {
	name := fmt.Sprintf("%s%d", string(rune(65+column)), row+1)
	return name
}

// AddData 添加数据
// rows, column 从0开始
func (e *excel) AddData(row, column int, data any) *excel {
	cell := e.CellString(row, column)
	err := e.SetCellValue(e.sheet, cell, data)
	if err != nil {
		snag.Panic(err)
	}
	return e
}

// AddValues 批量添加数据
func (e *excel) AddValues(rows ExcelItems) *excel {
	rowFrom := 0
	for _, data := range rows {
		e.column = 0
		rowFrom = e.row

		var mr int // 合并项
		var mergeColumns []int
		for _, v := range data {
			rt := reflect.TypeOf(v)
			if rt.Kind() == reflect.Slice {
				items := v.(ExcelItems)
				colFrom := e.column
				for m, subs := range items {
					for n, sub := range subs {
						e.AddData(e.row+m, colFrom+n, sub)
					}
				}
				// 子项目结束后需要加上列数
				e.column += items.Columns() - 1
				mr = len(items) - 1
			} else {
				e.AddData(e.row, e.column, v)
				mergeColumns = append(mergeColumns, e.column)
			}

			e.column += 1
		}

		e.row += mr

		// 合并单元格
		if mr > 0 && len(mergeColumns) > 0 {
			for _, mc := range mergeColumns {
				err := e.MergeCell(e.sheet, e.CellString(rowFrom, mc), e.CellString(e.row, mc))
				if err != nil {
					snag.Panic(err)
				}
			}
		}

		e.row += 1
	}
	return e
}

// Done 保存文件
func (e *excel) Done() string {
	err := e.SaveAs(e.path)
	if err != nil {
		snag.Panic(err)
	}
	return e.path
}
