// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-07
// Based on aurservd by liasica, magicrolan@qq.com.

package tools

import (
    "fmt"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/labstack/echo/v4"
    "github.com/xuri/excelize/v2"
    "os"
)

type excel struct {
    path  string
    sheet string

    w echo.Context

    *excelize.File
}

func NewExcel(w echo.Context, fp string, args ...any) *excel {
    err := utils.NewFile(fp).CreateDirectoryIfNotExist()
    if err != nil {
        snag.Panic("文件创建失败")
    }

    sheet := "Sheet1"
    if len(args) > 0 {
        sheet = args[0].(string)
    }

    fe := excelize.NewFile()
    fe.SetSheetName("Sheet1", sheet)

    return &excel{
        path:  fp,
        sheet: sheet,
        w:     w,
        File:  fe,
    }
}

func NewExcelExistsExport(w echo.Context, fp string, args ...any) (e *excel) {
    _, err := os.Stat(fp)
    if err == nil {
        err = w.File(fp)
        if err != nil {
            snag.Panic(err)
        }
        return nil
    }
    return NewExcel(w, fp, args...)
}

// AddData 添加数据
// rows, column 从0开始
func (e *excel) AddData(row, column int, data any) *excel {
    err := e.SetCellValue(e.sheet, fmt.Sprintf("%s%d", string(rune(65+column)), row+1), data)
    if err != nil {
        snag.Panic(err)
    }
    return e
}

// AddValues 批量添加数据
func (e *excel) AddValues(rows [][]any) *excel {
    for m, columns := range rows {
        for n, column := range columns {
            e.AddData(m+1, n+1, column)
        }
    }
    return e
}

// Done 保存文件
func (e *excel) Done() *excel {
    err := e.SaveAs(e.path)
    if err != nil {
        snag.Panic(err)
    }
    return e
}

// Export 导出文件
func (e *excel) Export() {
    err := e.w.File(e.path)
    if err != nil {
        snag.Panic(err)
    }
}
