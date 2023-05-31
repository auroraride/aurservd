// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/22
// Based on aurservd by liasica, magicrolan@qq.com.

package mob

import (
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/pkg/utils"
)

const (
	// pushUrl 推送请求URL
	pushUrl = `http://api.push.mob.com/v3/push/createPush`

	source = "webapi"
)

const (
	broadcast  = iota + 1 // 广播
	alias                 // 别名
	tag                   // 标签
	regid                 // 设备
	geo                   // 地理位置
	group                 // 用户分组
	complexGeo            // 复杂地理
)

const (
	envProduction  = "production"
	envDevelopment = "development"
)

const (
	PlatformAndroid = iota + 1
	PlatformiOS
)

const (
	typeNotify = iota + 1
	typeCustom
)

const (
	iosDevelopment = iota // iOS测试环境
	iosProduction         // iOS生产环境
)

const (
	ChannelSystem  = "system"  // 系统通知
	ChannelOverdue = "overdue" // 逾期通知
)

type mobPush struct {
	appKey        string
	appSecret     string
	iosProduction int
}

func NewPush() *mobPush {
	cfg := ar.Config.Mob.Push
	m := &mobPush{
		appKey:    cfg.AppKey,
		appSecret: cfg.AppSecret,
	}
	if cfg.Env == envProduction {
		m.iosProduction = iosProduction
	}
	return m
}

type Req struct {
	RegId       string
	Platform    int
	Content     string
	Title       string
	MessageData []MessageData
	Channel     string
}

func (m *mobPush) SendMessage(req Req) {
	data := &Message{
		Source: source,
		Appkey: m.appKey,
		PushTarget: &PushTarget{
			Target: regid,
			Rids: []string{
				req.RegId,
			},
		},
		PushNotify: &PushNotify{
			Plats:          []int{req.Platform},
			OfflineSeconds: 7 * 86400,
			Content:        req.Content,
			Title:          req.Title,
			Type:           typeNotify,
			Policy:         1,
			IOSProduction:  m.iosProduction,
		},
	}
	switch req.Platform {
	case PlatformiOS:
		// TODO ios消息结构
		data.PushNotify.IOSNotify = &IOSNotify{
			Badge:     1,
			BadgeType: 2,
		}
	case PlatformAndroid:
		data.PushNotify.AndroidNotify = &AndroidNotify{
			// AndroidChannelId: req.Channel,
			Warn:             "123",
			AndroidBadgeType: 2,
			AndroidBadge:     1,
		}
		data.PushFactoryExtra = &PushFactoryExtra{
			XiaomiExtra: XiaomiExtra{ChannelId: "high_system"},
			OppoExtra:   OppoExtra{ChannelId: "system"},
			VivoExtra:   VivoExtra{Classification: "1"},
		}
	}
	data.PushNotify.ExtrasMapList = append(data.PushNotify.ExtrasMapList, req.MessageData...)
	// 排序并转换json字符串
	b, _ := jsoniter.Marshal(data)
	s := string(b)
	// client.SetHeaders()
	// 生成sign
	sign := utils.Md5String(s + m.appSecret)

	_, _ = resty.New().R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"key":          m.appKey,
			"sign":         sign,
		}).
		SetBody(s).
		Post(pushUrl)
}
