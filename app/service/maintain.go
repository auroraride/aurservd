// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-06
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/setting"
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
    ent.Database.Setting.Update().Where(setting.Key(model.SettingMaintain)).SetContent(fmt.Sprintf("%v", state)).Exec(s.ctx)
}
