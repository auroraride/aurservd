// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/h2non/filetype"
    "github.com/h2non/filetype/matchers"
    "github.com/h2non/filetype/types"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "github.com/xuri/excelize/v2"
    "mime/multipart"
    "strings"
)

type BaseService struct {
    ctx      context.Context
    rider    *model.Rider
    modifier *model.Modifier
    employee *model.Employee
}

func newService(params ...any) (bs *BaseService) {
    bs = &BaseService{}
    ctx := context.Background()
    for _, param := range params {
        switch p := param.(type) {
        case *ent.Rider:
            bs.rider = &model.Rider{
                ID:    p.ID,
                Phone: p.Phone,
                Name:  p.Name,
            }
            ctx = context.WithValue(ctx, "rider", bs.rider)
        case *ent.Manager:
            bs.modifier = &model.Modifier{
                ID:    p.ID,
                Phone: p.Phone,
                Name:  p.Name,
            }
            ctx = context.WithValue(ctx, "modifier", bs.modifier)
        case *model.Modifier:
            bs.modifier = p
            ctx = context.WithValue(ctx, "modifier", bs.modifier)
        case *ent.Employee:
            bs.employee = &model.Employee{
                ID:    p.ID,
                Name:  p.Name,
                Phone: p.Phone,
            }
            ctx = context.WithValue(ctx, "employee", bs.employee)
        }
    }

    bs.ctx = ctx

    return
}

func (s *BaseService) GetXlsxDataX(c echo.Context) (rows [][]string) {
    source, err := c.FormFile("file")
    if err != nil {
        snag.Panic("未获取到上传的文件")
        return
    }

    var f multipart.File
    f, err = source.Open()
    if err != nil {
        snag.Panic(err)
    }
    defer func(f multipart.File) {
        _ = f.Close()
    }(f)

    var kind types.Type
    kind, err = filetype.MatchReader(f)
    if err != nil {
        log.Errorf("文件格式获取失败：%v", err)
        snag.Panic(err)
    }
    if kind != matchers.TypeXlsx {
        snag.Panic(fmt.Sprintf("文件格式错误，必须为标准xlsx格式，当前为：%s", kind.Extension))
    }
    _, _ = f.Seek(0, 0)

    var r *excelize.File
    r, err = excelize.OpenReader(f)
    if err != nil {
        snag.Panic(err)
    }
    defer func(r *excelize.File) {
        _ = r.Close()
    }(r)

    sheet := r.GetSheetName(0)
    rows, err = r.GetRows(sheet)

    if err != nil {
        snag.Panic(err)
    }

    for i, columns := range rows {
        for j, column := range columns {
            rows[i][j] = strings.TrimSpace(column)
        }
    }

    return
}
