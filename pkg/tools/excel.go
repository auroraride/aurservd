// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-07
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "fmt"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/xuri/excelize/v2"
    "reflect"
)

type excel struct {
    path  string
    sheet string

    *excelize.File
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
    fe.SetSheetName("Sheet1", sheet)
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
    err := e.SetCellValue(e.sheet, e.CellString(row, column), data)
    if err != nil {
        snag.Panic(err)
    }
    return e
}

// AddValuesFromStruct 从结构添加 TODO
func (e *excel) AddValuesFromStruct() *excel {
    return e
}

// AddValues 批量添加数据
func (e *excel) AddValues(rows [][]any) *excel {
    row := -1
    rowFrom := 0
    rowEnd := 0
    for _, data := range rows {
        row += 1
        rowFrom = row

        var mr int // 合并项
        var mergeColumns []int
        for column, v := range data {
            rt := reflect.TypeOf(v)
            if rt.Kind() == reflect.Slice {
                items := v.([]any)
                for m, subs := range items {
                    for n, sub := range subs.([]any) {
                        e.AddData(row+m, column+n, sub)
                    }
                }
                mr = len(items) - 1
                row += mr
            } else {
                e.AddData(row, column, v)
                mergeColumns = append(mergeColumns, column)
            }
        }

        rowEnd = row

        // 合并单元格
        if mr > 0 && len(mergeColumns) > 0 {
            for _, mc := range mergeColumns {
                err := e.MergeCell(e.sheet, e.CellString(rowFrom, mc), e.CellString(rowEnd, mc))
                if err != nil {
                    snag.Panic(err)
                }
            }
        }
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
