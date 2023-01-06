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
    "github.com/labstack/gommon/log"
    "os"
)

type maintainService struct {
    *BaseService
    maintainFile string
}

func NewMaintain(params ...any) *maintainService {
    return &maintainService{
        BaseService:  newService(params...),
        maintainFile: "runtime/MAINTAIN",
    }
}

func (s *maintainService) SetMaintain(state bool) {
    _ = ent.Database.Setting.Update().Where(setting.Key(model.SettingMaintain)).SetContent(fmt.Sprintf("%v", state)).Exec(s.ctx)
}

func (s *maintainService) UpdateMaintain() {
    _, err := os.Stat(s.maintainFile)
    if err == nil {
        s.SetMaintain(false)
    }
}

func (s *maintainService) CreateMaintainFile() {
    f, err := os.OpenFile(s.maintainFile, os.O_CREATE, 0644)
    if err != nil {
        log.Errorf("%s create failed: %v", s.maintainFile, err)
        return
    }

    defer func(f *os.File) {
        _ = f.Close()
    }(f)
}
