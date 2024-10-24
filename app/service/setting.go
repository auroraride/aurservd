// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/auroraride/adapter"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/setting"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
)

type settingService struct {
	ctx      context.Context
	modifier *model.Modifier
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
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func (s *settingService) ParseKey(key string) string {
	return strings.ToUpper(key)
}

// CacheSettings 缓存设置
func (s *settingService) CacheSettings(sm *ent.Setting) {
	switch sm.Key {
	case model.SettingBatteryFullKey,
		model.SettingDepositKey,
		model.SettingRenewalKey,
		model.SettingPauseMaxDaysKey,
		model.SettingExchangeIntervalKey,
		model.SettingRescueFeeKey,
		model.SettingReserveDurationKey,
		model.SettingExchangeMinBatteryKey:
		f, err := strconv.ParseFloat(strings.ReplaceAll(sm.Content, `"`, ""), 64)
		if err == nil {
			cache.Set(s.ctx, sm.Key, f, 0)
		}
	case model.SettingExchangeLimitKey,
		model.SettingExchangeFrequencyKey:
		cache.Set(s.ctx, sm.Key, adapter.ConvertString2Bytes(sm.Content), -1)
	}
}

// Initialize 初始化
func (s *settingService) Initialize() {
	for k, set := range model.Settings {
		sm, _ := s.orm.Query().Where(setting.Key(s.ParseKey(k))).First(s.ctx)
		if sm == nil {
			// 创建
			var err error
			b, _ := jsoniter.Marshal(set.Default)
			sm, err = s.orm.Create().SetKey(k).
				SetDesc(set.Desc).
				SetContent(string(b)).
				Save(s.ctx)
			if err != nil {
				zap.L().Fatal("设置初始化失败", zap.Error(err))
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
	sm, _ := s.orm.Query().Where(setting.Key(k)).First(s.ctx)
	if sm == nil {
		snag.Panic("未找到设置项")
	}
	var err error
	switch *req.Key {
	case model.SettingExchangeLimitKey:
		if *req.Content == "[]" || *req.Content == "" {
			*req.Content = "{}"
		} else {
			var data map[string]model.RiderExchangeLimit
			err = jsoniter.Unmarshal(adapter.ConvertString2Bytes(*req.Content), &data)
			for key, limit := range data {
				if limit.Duplicate() {
					snag.Panic("设定重复")
				}
				data[key].Sort()
			}
		}
	case model.SettingExchangeFrequencyKey:
		if *req.Content == "[]" || *req.Content == "" {
			*req.Content = "{}"
		} else {
			var data map[string]model.RiderExchangeFrequency
			err = jsoniter.Unmarshal(adapter.ConvertString2Bytes(*req.Content), &data)
			for key, limit := range data {
				if limit.Duplicate() {
					snag.Panic("设定重复")
				}
				data[key].Sort()
			}
		}
	}
	if err != nil {
		snag.Panic(err)
	}

	sm = s.orm.UpdateOne(sm).
		SetContent(*req.Content).
		SaveX(s.ctx)

	s.CacheSettings(sm)
}

func GetSetting[T any](key string) (v T, err error) {
	var set *ent.Setting
	set, err = ent.Database.Setting.Query().Where(setting.Key(key)).First(context.Background())
	if err != nil {
		return
	}

	if set == nil {
		err = errors.New("未找到设置")
		return
	}

	err = jsoniter.Unmarshal([]byte(set.Content), &v)
	return
}

// GetSetting 获取设置
func (s *settingService) GetSetting(key string) (v any) {
	d, ok := model.Settings[key]
	if !ok {
		snag.Panic("未找到设置")
	}

	set, err := s.orm.Query().Where(setting.Key(key)).First(s.ctx)
	if err != nil {
		snag.Panic("未找到设置")
	}

	if set == nil {
		return d.Default
	}

	_ = jsoniter.Unmarshal([]byte(set.Content), &d.Default)

	return d.Default
}

// SystemMaintain 检查是否维护中
func (s *settingService) SystemMaintain() bool {
	sm, _ := s.GetSetting(model.SettingMaintainKey).(bool)
	return sm
}

// SystemMaintainX 检查是否维护中
func (s *settingService) SystemMaintainX() {
	if s.SystemMaintain() {
		snag.Panic("正在唤醒电柜, 请稍后")
	}
}

func (s *settingService) Question() (v []interface{}) {
	v, _ = s.GetSetting(model.SettingQuestionKey).([]interface{})
	if len(v) == 0 {
		v = make([]interface{}, 0)
	}
	return v
}

// DailyRentItems 日租金设定
func (s *settingService) DailyRentItems() (data map[string]float64) {
	items, ok := s.GetSetting(model.SettingDailyRent).(map[string]any)
	if !ok {
		return
	}

	data = make(map[string]float64)
	for k, v := range items {
		data[k] = v.(float64)
	}
	return
}

// DailyRent 获取日租金
func (s *settingService) DailyRent(items map[string]float64, cityID uint64, batteryModel string, brandID *uint64) (key string, v float64) {
	key = strconv.FormatInt(int64(cityID), 10) + "_" + batteryModel
	if brandID != nil && *brandID > 0 {
		key += "_" + strconv.FormatInt(int64(*brandID), 10)
	}
	if items == nil {
		items = s.DailyRentItems()
	}
	v = items[key]
	if v == 0 {
		v = model.DailyRentDefault
	}
	return
}
