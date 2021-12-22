// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/22
// Based on aurservd by liasica, magicrolan@qq.com.

package mob

type PushTarget struct {
    Target int      `json:"target,omitempty"`
    Rids   []string `json:"rids,omitempty"`
}

type MessageData struct {
    Key   string `json:"key,omitempty"`
    Value string `json:"value,omitempty"`
}

type AndroidNotify struct {
    Warn             string `json:"warn,omitempty"`
    Sound            string `json:"sound,omitempty"`
    AndroidChannelId string `json:"androidChannelId,omitempty"`
    AndroidBadgeType int    `json:"androidBadgeType,omitempty"`
    AndroidBadge     int    `json:"androidBadge,omitempty"`
}

type IOSNotify struct {
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
}

type PushNotify struct {
    Plats          []int          `json:"plats,omitempty"`
    IOSProduction  int            `json:"iosProduction"`
    OfflineSeconds int            `json:"offlineSeconds,omitempty"`
    Content        string         `json:"content,omitempty"`
    Title          string         `json:"title,omitempty"`
    Type           int            `json:"type,omitempty"`
    AndroidNotify  *AndroidNotify `json:"androidNotify,omitempty"`
    IOSNotify      *IOSNotify     `json:"iosNotify,omitempty"`
    TaskCron       int            `json:"taskCron,omitempty"`
    TaskTime       uint64         `json:"taskTime,omitempty"`
    Policy         int            `json:"policy,omitempty"` // 推送策略： * 1:先走tcp，再走厂商 * 2:先走厂商，再走tcp * 3:只走厂商 * 4:只走tcp (厂商透传policy只支持策略3或4)
    ExtrasMapList  []MessageData `json:"extrasMapList,omitempty"`
}

type PushCallback struct {
    Url    string            `json:"url,omitempty"`
    Params map[string]string `json:"params,omitempty"`
}

type PushForward struct {
    Url            string        `json:"url,omitempty"`
    Scheme         string        `json:"scheme,omitempty"`
    NextType       int           `json:"nextType,omitempty"`
    SchemeDataList []MessageData `json:"schemeDataList,omitempty"`
}

type XiaomiExtra struct {
    ChannelId string `json:"channelId,omitempty"`
}

type OppoExtra struct {
    ChannelId string `json:"channelId,omitempty"`
}

type VivoExtra struct {
    Classification string `json:"classification,omitempty"`
}

type PushFactoryExtra struct {
    XiaomiExtra XiaomiExtra `json:"xiaomiExtra,omitempty"`
    OppoExtra   OppoExtra   `json:"oppoExtra,omitempty"`
    VivoExtra   VivoExtra   `json:"vivoExtra,omitempty"`
}

// Message 推送消息结构体
// @doc https://mob.com/wiki/detailed?wiki=234&id=136
type Message struct {
    Workno           string            `json:"workno,omitempty"`
    Source           string            `json:"source,omitempty"`
    Appkey           string            `json:"appkey,omitempty"`
    PushTarget       *PushTarget       `json:"pushTarget,omitempty"`
    PushNotify       *PushNotify       `json:"pushNotify,omitempty"`
    PushCallback     *PushCallback     `json:"pushCallback,omitempty"`
    PushForward      *PushForward      `json:"pushForward,omitempty"`
    PushFactoryExtra *PushFactoryExtra `json:"pushFactoryExtra,omitempty"`
}
