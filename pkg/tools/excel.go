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
            e.AddData(m, n, column)
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
