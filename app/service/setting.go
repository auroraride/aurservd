// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "encoding/json"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/setting"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    jsoniter "github.com/json-iterator/go"
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
        orm: ent.Database.Setting,
    }
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
    case model.SettingBatteryFull,
        model.SettingDeposit,
        model.SettingRenewal,
        model.SettingPauseMaxDays,
        model.SettingExchangeInterval,
        model.SettingRescueFee,
        model.SettingReserveDuration:
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

// GetSetting 获取设置
func (s *settingService) GetSetting(key string) (v any) {
    d, ok := model.Settings[key]
    if !ok {
        snag.Panic("未找到设置")
    }

    set, err := s.orm.Query().Where(setting.Key(key)).First(s.ctx)
    if err != nil {
        log.Error(err)
        snag.Panic("未找到设置")
    }

    if set == nil {
        return d.Default
    }

    err = jsoniter.Unmarshal([]byte(set.Content), &d.Default)
    if err != nil {
        log.Error(err)
    }

    return d.Default
}

// SystemMaintain 检查是否维护中
func (s *settingService) SystemMaintain() bool {
    sm, _ := NewSetting().GetSetting(model.SettingMaintain).(bool)
    return sm
}

// SystemMaintainX 检查是否维护中
func (s *settingService) SystemMaintainX() {
    if s.SystemMaintain() {
        snag.Panic("正在唤醒电柜, 请稍后")
    }
}
