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
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/assetexport"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type assetExportService struct {
	ctx      context.Context
	modifier *model.Modifier
	orm      *ent.AssetExportClient
}

func NewAssetExportWithModifier(m *model.Modifier) *assetExportService {
	s := &assetExportService{
		ctx: context.Background(),
		orm: ent.Database.AssetExport,
	}
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

// Start 开始执行任务
func (s *assetExportService) Start(taxonomy string, con any, data map[string]interface{}, remark string, cb func(path string)) model.AssetExportRes {
	sn := tools.NewUnique().NewSN()
	go func() {
		b, _ := jsoniter.Marshal(con)

		var ex *ent.AssetExport
		var message string
		var err error
		info := make(ar.Map)

		ex, err = s.orm.Create().
			SetCondition(adapter.ConvertBytes2String(b)).
			SetRemark(remark).
			SetTaxonomy(taxonomy).
			SetSn(sn).
			SetAssetManagerID(s.modifier.ID).
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
			case *ent.ExportInfo:
				info[k] = ei.GetExportInfoData()
			default:
				info[k] = v
			}
		}

		cb(path)
	}()
	return model.AssetExportRes{SN: sn}
}

func (s *assetExportService) List(m *ent.AssetManager, req *model.AssetExportListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Order(ent.Desc(assetexport.FieldCreatedAt)).WithAssetManager()
	if req.SN != "" {
		q.Where(assetexport.Sn(req.SN))
	}
	r := m.Edges.Role
	if r == nil || !r.Super {
		q.Where(assetexport.AssetManagerID(m.ID))
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.AssetExport) model.AssetExportListRes {
		opName := ""
		if item.Edges.AssetManager != nil {
			opName = item.Edges.AssetManager.Name
		}
		res := model.AssetExportListRes{
			CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
			Operator:  opName,
			Remark:    item.Remark,
			Taxonomy:  item.Taxonomy,
			SN:        item.Sn,
			Status:    model.AssetExportStatus(item.Status).String(),
			Message:   item.Message,
			Info:      item.Info,
		}
		if !item.FinishAt.IsZero() {
			res.FinishAt = item.FinishAt.Format(carbon.DateTimeLayout)
		}
		return res
	})
}

func (s *assetExportService) Download(req *model.AssetExportDownloadReq) (string, string) {
	e, err := s.orm.QueryNotDeleted().Where(assetexport.Sn(req.SN), assetexport.Status(uint8(model.AssetExportStatusSuccess))).First(s.ctx)
	if err != nil {
		snag.Panic("未找到文件")
	}
	return e.Path, url.QueryEscape(filepath.Base(e.Path))
}
