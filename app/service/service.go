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
    ctx context.Context

    rider    *model.Rider
    entRider *ent.Rider

    modifier *model.Modifier

    employee    *model.Employee
    entEmployee *ent.Employee

    entStore *ent.Store
}

func newService(params ...any) (bs *BaseService) {
    bs = &BaseService{}
    ctx := context.Background()
    for _, param := range params {
        if param == nil {
            continue
        }
        switch p := param.(type) {
        case *ent.Rider:
            bs.entRider = p
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
            bs.entEmployee = p
            bs.entStore, _ = p.QueryStore().First(ctx)
            bs.employee = &model.Employee{
                ID:    p.ID,
                Name:  p.Name,
                Phone: p.Phone,
            }
            ctx = context.WithValue(ctx, "employee", bs.employee)
        case *ent.Store:
            bs.entStore = p
        }
    }

    bs.ctx = ctx

    return
}

// GetXlsxRows 从xlsx文件中读取数据
// start 从第几行开始为数据
// columnsNumber 每行数据数量
// pkIndex 主键下标(以此排重)
func (s *BaseService) GetXlsxRows(c echo.Context, start, columnsNumber int, pkIndex int) (rows [][]string, pks, failed []string) {
    failed = make([]string, 0)
    source, err := c.FormFile("file")
    if err != nil {
        log.Errorf("GetXlsxDataX error: %s", err)
        snag.Panic("未获取到上传的文件: " + err.Error())
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

    // 主键 => 行数(i+1)
    m := make(map[string]int)
    for i, columns := range rows {
        // 排错
        if len(columns) < columnsNumber {
            failed = append(failed, fmt.Sprintf("格式错误:%s", strings.Join(columns, ",")))
            continue
        }
        for j, column := range columns {
            t := strings.TrimSpace(column)
            // 去重
            if j == pkIndex {
                if target, ok := m[t]; ok {
                    failed = append(failed, fmt.Sprintf("第%d行和第%d行重复", i+1, target))
                    continue
                }
                m[t] = i + 1
                pks = append(pks, t)
            }
            rows[i][j] = t
        }
    }

    if len(rows) < start {
        snag.Panic("至少有一条有效信息")
    }

    return
}
