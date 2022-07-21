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
    "net/url"
    "os"
    "path/filepath"
    "reflect"
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

func sendAttachement(w echo.Context, path string) {
    // w.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    // w.Response().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename*=utf-8''%s`, url.QueryEscape(filepath.Base(path))))
    // err := w.File(path)
    err := w.Attachment(path, url.QueryEscape(filepath.Base(path)))
    if err != nil {
        snag.Panic(err)
    }
}

func NewExcelExistsExport(w echo.Context, fp string, args ...any) (e *excel) {
    _, err := os.Stat(fp)
    if err == nil {
        sendAttachement(w, fp)
        return nil
    }
    return NewExcel(w, fp, args...)
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
func (e *excel) Done() *excel {
    err := e.SaveAs(e.path)
    if err != nil {
        snag.Panic(err)
    }
    return e
}

// Export 导出文件
func (e *excel) Export() {
    sendAttachement(e.w, e.path)
}
