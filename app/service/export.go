// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/tools"
    jsoniter "github.com/json-iterator/go"
    log "github.com/sirupsen/logrus"
    "path/filepath"
    "time"
)

type exportService struct {
    ctx      context.Context
    modifier *model.Modifier
    orm      *ent.ExportClient
}

func NewExportWithModifier(m *model.Modifier) *exportService {
    s := &exportService{
        ctx: context.Background(),
        orm: ent.Database.Export,
    }
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// Start 开始执行任务
func (s *exportService) Start(taxonomy string, con any, info map[string]interface{}, remark string, cb func(path string)) model.ExportRes {
    sn := tools.NewUnique().NewSN()
    go func() {
        for k, v := range info {
            switch v.(type) {
            case func() string:
                info[k] = v.(func() string)()
                break
            }
        }

        b, _ := jsoniter.Marshal(con)

        var ex *ent.Export
        var message string
        var err error

        ex, err = s.orm.Create().
            SetCondition(b).
            SetRemark(remark).
            SetTaxonomy(taxonomy).
            SetInfo(info).
            SetSn(sn).
            Save(s.ctx)
        if err != nil {
            log.Error(err)
            return
        }

        path := filepath.Join("runtime/export", fmt.Sprintf("%s-%s.xlsx", ex.Sn, taxonomy))

        defer func() {
            status := model.ExportStatusSuccess
            if path == "" {
                status = model.ExportStatusFail
            }

            if r := recover(); r != nil {
                message = fmt.Sprintf("%v", r)
                status = model.ExportStatusFail
            }

            now := time.Now()
            _, _ = ex.Update().SetStatus(status).SetMessage(message).SetPath(path).SetDuration(now.Sub(ex.CreatedAt).Milliseconds()).SetFinishAt(now).Save(s.ctx)
        }()

        cb(path)
    }()
    return model.ExportRes{SN: sn}
}
