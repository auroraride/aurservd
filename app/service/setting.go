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
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    log "github.com/sirupsen/logrus"
    "strconv"
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

// CacheSettings 缓存设置
func (s *settingService) CacheSettings(sm *ent.Setting) {
    switch sm.Key {
    case model.SettingBatteryFull, model.SettingDeposit:
        f, err := strconv.ParseFloat(strings.ReplaceAll(sm.Content, `"`, ""), 10)
        if err == nil {
            cache.Set(s.ctx, sm.Key, f, 0)
        }
        break
    }
}

// Initialize 初始化
func (s *settingService) Initialize() {
    for k, set := range model.Settings {
        sm, _ := s.orm.Query().Where(setting.Key(s.ParseKey(k))).Only(s.ctx)
        if sm == nil {
            // 创建
            var err error
            b, _ := json.Marshal(set.Default)
            sm, err = s.orm.Create().SetKey(k).
                SetDesc(set.Desc).
                SetContent(string(b)).
                Save(s.ctx)
            if err != nil {
                log.Fatal(err)
            }
        }
        s.CacheSettings(sm)
    }
}

// List 列举设置
func (s *settingService) List() (items []model.SettingRes) {
    s.orm.Query().Select(setting.FieldKey, setting.FieldContent, setting.FieldDesc).ScanX(s.ctx, &items)
    return
}

// Modify 修改设置
func (s *settingService) Modify(req *model.SettingReq) {
    k := s.ParseKey(*req.Key)
    sm := s.orm.Query().Where(setting.Key(k)).OnlyX(s.ctx)
    if sm == nil {
        snag.Panic("未找到设置项")
    }
    sm = s.orm.UpdateOne(sm).
        SetContent(*req.Content).
        SaveX(s.ctx)

    s.CacheSettings(sm)
}
