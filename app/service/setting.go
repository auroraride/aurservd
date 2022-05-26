// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "encoding/json"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/setting"
    "github.com/auroraride/aurservd/pkg/snag"
    log "github.com/sirupsen/logrus"
    "strings"
)

type settingService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.SettingClient
}

func NewSetting() *settingService {
    return &settingService{
        ctx: context.Background(),
        orm: ar.Ent.Setting,
    }
}

func NewSettingWithRider(rider *ent.Rider) *orderService {
    s := NewOrder()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewSettingWithModifier(m *model.Modifier) *settingService {
    s := NewSetting()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *settingService) ParseKey(key string) string {
    return strings.ToUpper(key)
}

// Initialize 初始化
func (s *settingService) Initialize() {
    for k, set := range model.Settings {
        exist, err := s.orm.Query().Where(setting.Key(s.ParseKey(k))).Exist(s.ctx)
        if err == nil && !exist {
            // 创建
            b, _ := json.Marshal(set.Default)
            _, err = s.orm.Create().SetKey(k).
                SetDesc(set.Desc).
                SetContent(string(b)).
                Save(s.ctx)
            if err != nil {
                log.Fatal(err)
            }
        }
    }
}

// List 列举设置
func (s *settingService) List() (items []model.SettingReq) {
    s.orm.Query().Select(setting.FieldKey, setting.FieldContent).ScanX(s.ctx, &items)
    return
}

// Modify 修改设置
func (s *settingService) Modify(req *model.SettingReq) {
    k := s.ParseKey(*req.Key)
    _, ok := model.Settings[k]
    if !ok {
        snag.Panic("未找到设置项")
    }
    s.orm.Update().
        Where(setting.Key(k)).
        SetContent(*req.Content).
        SaveX(s.ctx)
}
