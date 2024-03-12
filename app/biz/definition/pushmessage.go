package definition

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

const (
	PushTypeNow      = uint8(iota) + 1 // 立即推送
	PushTypeSchedule                   // 定时推送
)

const (
	PushStatusPushed   = uint8(iota) + 1 // 已推送
	PushStatusUnPushed                   // 未推送
)

const (
	PushMessageTypeAll    = uint8(iota) + 1 // 全部类型
	PushMessageTypeNotice                   // 通知消息
)

// PushmessageDetail 推送消息详情
type PushmessageDetail struct {
	ID            uint64 `json:"id"`             // ID
	Image         string `json:"image"`          // 封面图片
	Title         string `json:"title"`          // 标题
	Content       string `json:"content"`        // 内容
	PushType      uint8  `json:"push_type"`      // 推送类型
	PushTime      string `json:"push_time"`      // 推送时间
	IsHome        bool   `json:"is_home"`        // 是否首页推送
	HomeContent   string `json:"home_content"`   // 首页推送内容
	MessageStatus uint8  `json:"message_status"` // 消息状态, PushStatusPushed:已推送, PushStatusUnPushed:未推送
	MessageType   uint8  `json:"message_type"`   // 消息类型
}

// PushmessageCommon 推送消息通用字段
type PushmessageCommon struct {
	Title       string  `json:"title" form:"title" validate:"required"`         // 标题
	Image       string  `json:"image" form:"image"`                             // 封面图片
	Content     string  `json:"content" form:"content" validate:"required"`     // 内容
	PushType    uint8   `json:"push_type" form:"push_type" validate:"required"` // 推送类型, PushTypeNow:立即推送, PushTypeSchedule:定时推送
	PushTime    *string `json:"push_time" form:"push_time"`                     // 推送时间，如果是定时推算，则必填
	IsHome      bool    `json:"is_home" form:"is_home"`                         // 是否首页推送
	HomeContent string  `json:"home_content" form:"home_content"`               // 首页推送内容
	MessageType uint8   `json:"message_type" form:"type"`                       // 消息类型, PushMessageTypeAll:全部类型, PushMessageTypeNotice:通知消息
}

// PushmessageSaveReq 推送消息保存请求
type PushmessageSaveReq struct {
	PushmessageCommon
}

// PushmessageModifyReq 推送消息修改请求
type PushmessageModifyReq struct {
	model.IDParamReq
	PushmessageCommon
}

// PushmessageDeleteReq 推送消息删除请求
type PushmessageDeleteReq struct {
	model.IDParamReq
}

// PushmessageGetReq 推送消息获取请求
type PushmessageGetReq struct {
	model.IDParamReq
}

// PushmessageListReq 推送消息列表请求
type PushmessageListReq struct {
	model.PaginationReq
	Keyword       *string `json:"keyword"`
	MessageStatus *uint8  `json:"message_status"`
	MessageType   *uint8  `json:"message_type"`
	PushBeginTime *string `json:"push_begin_time"`
	PushEndTime   *string `json:"push_end_time"`
}

// ToDetail 转换为推送消息详情
func ToDetail(pd *ent.Pushmessage) *PushmessageDetail {
	return &PushmessageDetail{
		ID:            pd.ID,
		Image:         pd.Image,
		Title:         pd.Title,
		Content:       pd.Content,
		PushType:      pd.PushType,
		PushTime:      pd.PushTime.Format("2006-01-02 15:04:05"),
		IsHome:        pd.IsHome,
		HomeContent:   pd.HomeContent,
		MessageStatus: pd.MessageStatus,
		MessageType:   pd.MessageType,
	}
}
