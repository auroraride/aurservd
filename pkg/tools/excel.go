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
    return fmt.Sprintf("%s%d", string(rune(65+column)), row+1)
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

// AddValues 批量添加数据
func (e *excel) AddValues(rows [][]any) *excel {
    for row, data := range rows {
        for column, v := range data {
            rt := reflect.TypeOf(v)
            if rt.Kind() == reflect.Slice {
                items := v.([]any)
                for m, subs := range items {
                    for n, sub := range subs.([]any) {
                        if m > 0 {
                            err := e.MergeCell(e.sheet, e.CellString(row, column+n), e.CellString(row+m, column+n))
                            if err != nil {
                                snag.Panic(err)
                            }
                        }
                        e.AddData(row+m, column+n, sub)
                    }
                }
            } else {
                e.AddData(row, column, v)
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
