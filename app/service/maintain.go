// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-06
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"

	"github.com/auroraride/adapter/maintain"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/setting"
	"github.com/labstack/gommon/log"
)

type maintainService struct {
	*BaseService
}

func NewMaintain(params ...any) *maintainService {
	return &maintainService{
		BaseService: newService(params...),
	}
}

func (s *maintainService) SetMaintain(state bool) {
	_ = ent.Database.Setting.Update().Where(setting.Key(model.SettingMaintainKey)).SetContent(fmt.Sprintf("%v", state)).Exec(s.ctx)
}

func (s *maintainService) UpdateMaintain() {
	if maintain.Exists() {
		s.SetMaintain(false)
		_ = maintain.Remove()
	}
}

func (s *maintainService) CreateMaintainFile() {
	err := maintain.Create()
	if err != nil {
		log.Errorf("%s create failed: %v", maintain.File(), err)
		return
	}
	NewMaintain().SetMaintain(true)
}
