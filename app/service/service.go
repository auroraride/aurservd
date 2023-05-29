// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/auroraride/adapter"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

type BaseService struct {
	ctx context.Context

	rider    *model.Rider
	entRider *ent.Rider

	modifier *model.Modifier

	employee    *model.Employee
	entEmployee *ent.Employee

	entStore *ent.Store

	agent      *ent.Agent
	enterprise *ent.Enterprise
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
		case *model.Rider:
			bs.rider = p
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
		case *ent.Agent:
			bs.agent = p
		case *ent.Enterprise:
			bs.enterprise = p
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

	var rawRows [][]string
	rawRows, err = r.GetRows(sheet)

	if err != nil {
		snag.Panic(err)
	}

	// 主键 => 行数(i+1)
	m := make(map[string]int)
	for i, columns := range rawRows {
		if i < start-1 {
			continue
		}
		// 排错
		if len(columns) < columnsNumber {
			failed = append(failed, fmt.Sprintf("格式错误:%s", strings.Join(columns, ",")))
			continue
		}

		column := make([]string, columnsNumber)
		for j, rc := range columns {
			t := strings.TrimSpace(rc)
			// 去重
			if j == pkIndex {
				if target, ok := m[t]; ok {
					failed = append(failed, fmt.Sprintf("第%d行和第%d行重复", i+1, target))
					continue
				}
				m[t] = i + 1
				pks = append(pks, t)
			}
			column[j] = t
		}

		rows = append(rows, column)
	}

	if len(rows) < start {
		snag.Panic("至少有一条有效信息")
	}

	return
}

func (s *BaseService) GetCabinetAdapterUrl(cab *ent.Cabinet, apiurl string) (url string, err error) {
	switch cab.Brand {
	case adapter.CabinetBrandKaixin:
		url = ar.Config.Sync.Kxcab.Api
		if !cab.Intelligent {
			url = ar.Config.Sync.Kxnicab.Api
		}
	case adapter.CabinetBrandYundong:
		url = ar.Config.Sync.Ydcab.Api
	case adapter.CabinetBrandTuobang:
		url = ar.Config.Sync.Tbcab.Api
	case adapter.CabinetBrandXiliulouServer:
		url = ar.Config.Sync.Xllscab.Api
	default:
		return "", adapter.ErrorCabinetBrand
	}

	return url + apiurl, nil
}

func (s *BaseService) GetCabinetAdapterUrlX(cab *ent.Cabinet, apiurl string) string {
	url, err := s.GetCabinetAdapterUrl(cab, apiurl)
	if err != nil {
		snag.Panic(err)
	}
	return url
}
func (s *BaseService) GetAdapterUser() (user *adapter.User, err error) {
	switch true {
	default:
		return nil, adapter.ErrorUserRequired
	case s.rider != nil:
		return &adapter.User{
			Type: adapter.UserTypeRider,
			ID:   s.rider.Phone,
		}, nil
	case s.employee != nil:
		return &adapter.User{
			Type: adapter.UserTypeEmployee,
			ID:   s.employee.Phone,
		}, nil
	case s.modifier != nil:
		return &adapter.User{
			Type: adapter.UserTypeManager,
			ID:   s.modifier.Phone,
		}, nil
	}
}

func (s *BaseService) GetAdapterUserX() *adapter.User {
	user, err := s.GetAdapterUser()
	if err != nil {
		snag.Panic(err)
	}
	return user
}
