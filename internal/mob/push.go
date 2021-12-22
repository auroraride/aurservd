// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/22
// Based on aurservd by liasica, magicrolan@qq.com.

package mob

import "github.com/auroraride/aurservd/internal/ar"

type mobPush struct {
    AppKey    string
    AppSecret string
}

// Message 推送消息结构体
// @doc https://mob.com/wiki/detailed?wiki=234&id=136
type Message struct {
    Workno     string `json:"workno,omitempty"`
    Source     string `json:"source,omitempty"`
    Appkey     string `json:"appkey,omitempty"`
    PushTarget struct {
        Target int      `json:"target,omitempty"`
        Rids   []string `json:"rids,omitempty"`
    } `json:"pushTarget"`
    PushNotify struct {
        Plats          []int  `json:"plats,omitempty"`
        IosProduction  int    `json:"iosProduction,omitempty"`
        OfflineSeconds int    `json:"offlineSeconds,omitempty"`
        Content        string `json:"content,omitempty"`
        Title          string `json:"title,omitempty"`
        Type           int    `json:"type,omitempty"`
        AndroidNotify  struct {
            Warn              string `json:"warn,omitempty"`
            Sound             string `json:"sound,omitempty"`
            AndroidChannelId  string `json:"androidChannelId,omitempty"`
            AandroidBadgeType int    `json:"aandroidBadgeType,omitempty"`
            AndroidBadge      int    `json:"androidBadge,omitempty"`
        } `json:"androidNotify"`
        IosNotify struct {
            Badge            int    `json:"badge,omitempty"`
            BadgeType        int    `json:"badgeType,omitempty"`
            Category         string `json:"category,omitempty"`
            Sound            string `json:"sound,omitempty"`
            Subtitle         string `json:"subtitle,omitempty"`
            SlientPush       int    `json:"slientPush,omitempty"`
            ContentAvailable int    `json:"contentAvailable,omitempty"`
            MutableContent   int    `json:"mutableContent,omitempty"`
            AttachmentType   int    `json:"attachmentType,omitempty"`
            Attachment       string `json:"attachment,omitempty"`
        } `json:"iosNotify"`
        TaskCron      int               `json:"taskCron,omitempty"`
        TaskTime      uint64            `json:"taskTime,omitempty"`
        Policy        int               `json:"policy,omitempty"` // 推送策略： * 1:先走tcp，再走厂商 * 2:先走厂商，再走tcp * 3:只走厂商 * 4:只走tcp (厂商透传policy只支持策略3或4)
        ExtrasMapList []MessageKeyValue `json:"extrasMapList,omitempty"`
    } `json:"pushNotify"`
    PushCallback struct {
        Url    string            `json:"url,omitempty"`
        Params map[string]string `json:"params,omitempty"`
    } `json:"pushCallback"`
    PushForward struct {
        Url            string            `json:"url,omitempty"`
        Scheme         string            `json:"scheme,omitempty"`
        NextType       int               `json:"nextType,omitempty"`
        SchemeDataList []MessageKeyValue `json:"schemeDataList,omitempty"`
    } `json:"pushForward"`
    PushFactoryExtra struct {
        XiaomiExtra struct {
            ChannelId string `json:"channelId,omitempty"`
        } `json:"xiaomiExtra"`
        OppoExtra struct {
            ChannelId string `json:"channelId,omitempty"`
        } `json:"oppoExtra"`
        VivoExtra struct {
            Classification string `json:"classification,omitempty"`
        } `json:"vivoExtra"`
    } `json:"pushFactoryExtra"`
}

type MessageKeyValue struct {
    Key   string
    Value string
}

func NewPush() *mobPush {
    cfg := ar.Config.Mob.Push
    return &mobPush{
        AppKey:    cfg.AppKey,
        AppSecret: cfg.AppSecret,
    }
}

func (m *mobPush) SendMessage() {
}
