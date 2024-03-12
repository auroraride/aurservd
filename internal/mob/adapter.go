package mob

import (
	"errors"
	"time"
)

type PushOption struct {
	Title        *string
	Content      *string
	ScheduleTime *time.Time
	MessageData  map[string]string
}

// Push 推送消息
func Push(options *PushOption) (*Response, error) {
	msg := Req{
		Platform:    PlatformAndroid,
		MessageData: []MessageData{},
	}
	if options.Title != nil {
		msg.Title = *options.Title
	}
	if options.Content != nil {
		msg.Content = *options.Content
	}
	if options.ScheduleTime != nil {
		msg.TaskCron = 1
		msg.TaskTime = uint64(options.ScheduleTime.UnixNano() / int64(time.Millisecond))
	}
	if options.MessageData != nil && len(options.MessageData) > 0 {
		for k, v := range options.MessageData {
			msg.MessageData = append(msg.MessageData, MessageData{
				Key:   k,
				Value: v,
			})
		}
	}
	res, err := NewPush().SendMessage(msg)
	if err != nil {
		return nil, err
	}
	if res.Status != 200 {
		return nil, errors.New(res.Rrror)
	}
	return res, nil
}

// Drop 删除推送消息
func Drop(batchId string) (*Response, error) {
	res, err := NewPush().DropMessage(batchId)
	if err != nil {
		return nil, err
	}
	if res.Status != 200 {
		return nil, errors.New(res.Rrror)
	}
	return res, nil
}
