// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/auroraride/adapter"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/export"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
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
func (s *exportService) Start(taxonomy string, con any, data map[string]interface{}, remark string, cb func(path string)) model.ExportRes {
	sn := tools.NewUnique().NewSN()
	go func() {
		b, _ := jsoniter.Marshal(con)

		var ex *ent.Export
		var message string
		var err error
		info := make(ar.Map)

		ex, err = s.orm.Create().
			SetCondition(adapter.ConvertBytes2String(b)).
			SetRemark(remark).
			SetTaxonomy(taxonomy).
			SetSn(sn).
			SetManagerID(s.modifier.ID).
			Save(s.ctx)
		if err != nil {
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
			_, _ = ex.Update().SetStatus(uint8(status)).SetMessage(message).SetPath(path).SetDuration(now.Sub(ex.CreatedAt).Milliseconds()).SetInfo(info).SetFinishAt(now).Save(s.ctx)
		}()

		for k, v := range data {
			switch ei := v.(type) {
			case func() string:
				info[k] = ei()
				break
			case *ent.ExportInfo:
				info[k] = ei.GetExportInfoData()
				break
			default:
				info[k] = v
				break
			}
		}

		cb(path)
	}()
	return model.ExportRes{SN: sn}
}

func (s *exportService) List(m *ent.Manager, req *model.ExportListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Order(ent.Desc(export.FieldCreatedAt)).WithManager()
	if req.SN != "" {
		q.Where(export.Sn(req.SN))
	}
	r := m.Edges.Role
	if r == nil || !r.Super {
		q.Where(export.ManagerID(m.ID))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Export) model.ExportListRes {
		res := model.ExportListRes{
			CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
			Operator:  item.Edges.Manager.Name,
			Remark:    item.Remark,
			Taxonomy:  item.Taxonomy,
			SN:        item.Sn,
			Status:    model.ExportStatus(item.Status).String(),
			Message:   item.Message,
			Info:      item.Info,
		}
		if !item.FinishAt.IsZero() {
			res.FinishAt = item.FinishAt.Format(carbon.DateTimeLayout)
		}
		return res
	})
}

func (s *exportService) Download(req *model.ExportDownloadReq) (string, string) {
	e, err := s.orm.QueryNotDeleted().Where(export.Sn(req.SN), export.Status(uint8(model.ExportStatusSuccess))).First(s.ctx)
	if err != nil {
		snag.Panic("未找到文件")
	}
	return e.Path, url.QueryEscape(filepath.Base(e.Path))
}
