// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/h2non/filetype"
    "github.com/h2non/filetype/matchers"
    log "github.com/sirupsen/logrus"
    "github.com/xuri/excelize/v2"
    "mime/multipart"
)

type cscService struct {
    ctx context.Context
}

func NewCSC() *cscService {
    return &cscService{
        ctx: context.Background(),
    }
}

// ParseNameListShiguangju 时光驹催缴工具
func (*cscService) ParseNameListShiguangju(source *multipart.FileHeader) []*model.ShiguangjuIVRItem {
    tel := "02863804608"
    tmpl := "TTS_235791551"

    f, err := source.Open()
    if err != nil {
        snag.Panic(err)
    }
    defer func(f multipart.File) {
        _ = f.Close()
    }(f)

    kind, err := filetype.MatchReader(f)
    if err != nil {
        log.Errorf("文件格式获取失败：%v", err)
        snag.Panic(err)
    }
    if kind != matchers.TypeXlsx {
        snag.Panic(fmt.Sprintf("文件格式错误，必须为标准xlsx格式，当前为：%s", kind.Extension))
    }

    _, _ = f.Seek(0, 0)
    r, err := excelize.OpenReader(f)
    if err != nil {
        snag.Panic(err)
    }
    defer func(r *excelize.File) {
        _ = r.Close()
    }(r)

    sheet := r.GetSheetName(0)
    rows, _ := r.GetRows(sheet)
    if len(rows) < 2 {
        snag.Panic("至少有一条外呼请求，请严格遵循模板制作文档")
    }
    items := make([]*model.ShiguangjuIVRItem, len(rows)-1)
    for i, row := range rows {
        if i > 0 {
            if len(row) != 3 {
                snag.Panic("文档格式错误，请严格遵循模板制作文档")
            }
            item := &model.ShiguangjuIVRItem{
                Name:    row[0],
                Phone:   row[1],
                Product: row[2],
            }
            item.Status = ali.NewVms().SendVoiceMessageByTts(tools.NewPointerInterface(item.Phone), tools.NewPointerInterface(fmt.Sprintf(`{"name":"%s","product": "%s"}`, item.Name, item.Product)), &tel, &tmpl)
            items[i-1] = item
        }
    }
    return items
}
