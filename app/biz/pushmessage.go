// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-11, by lisicen

package biz

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/pushmessage"
	"github.com/auroraride/aurservd/internal/mob"
)

var (
	TriggerList []*time.Timer
)

type pushmessageBiz struct {
	orm *ent.PushmessageClient
}

func NewPushmessage() *pushmessageBiz {
	return &pushmessageBiz{
		orm: ent.Database.Pushmessage,
	}
}

// Create 保存推送消息
func (b *pushmessageBiz) Create(req *definition.PushmessageSaveReq) error {
	option := &mob.PushOption{
		Title:   &req.Title,
		Content: &req.Content,
	}
	// 推送状态，默认已推送
	pushStatus := definition.PushStatusPushed
	pushTime := fmtPushTime(req.PushTime)
	// 如果是定时推送设置推送时间
	if req.PushType == definition.PushTypeSchedule {
		option.ScheduleTime = pushTime
		pushStatus = definition.PushStatusUnPushed
	}
	// TODO 如果是主页推送，要同时开启APP内部推送
	res, err := mob.Push(option)
	if err != nil {
		return err
	}
	_, err = b.orm.Create().
		SetTitle(req.Title).
		SetImage(req.Image).
		SetContent(req.Content).
		SetPushType(req.PushType).
		SetPushTime(*pushTime).
		SetIsHome(req.IsHome).
		SetHomeContent(req.HomeContent).
		SetMessageStatus(pushStatus).
		SetMessageType(req.MessageType).
		SetThirdPartyID(res.Res.BatchId).
		Save(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// Modify 修改推送消息
func (b *pushmessageBiz) Modify(req *definition.PushmessageModifyReq) error {
	ctx := context.Background()
	old, err := b.orm.Get(context.Background(), req.ID)
	if err != nil {
		return err
	}
	// 已经推送的消息不允许修改
	// TODO 是否还需要PushStatusPushed
	if old.MessageStatus == definition.PushStatusPushed || old.PushTime.Before(time.Now()) {
		return errors.New("已推送的消息不允许修改")
	}
	// 删除旧的消息
	if _, err = mob.Drop(old.ThirdPartyID); err != nil {
		return err
	}
	// 重新推送新的消息
	option := &mob.PushOption{
		Title:   &req.Title,
		Content: &req.Content,
	}
	// 推送时间
	pushTime := fmtPushTime(req.PushTime)
	if req.PushType == definition.PushTypeSchedule {
		option.ScheduleTime = pushTime
	}
	// TODO
	res, err := mob.Push(option)
	if err != nil {
		return err
	}
	_, err = b.orm.UpdateOneID(req.ID).
		SetTitle(req.Title).
		SetImage(req.Image).
		SetContent(req.Content).
		SetPushType(req.PushType).
		SetPushTime(*pushTime).
		SetIsHome(req.IsHome).
		SetHomeContent(req.HomeContent).
		// SetMessageStatus(req.MessageStatus).
		SetMessageType(req.MessageType).
		SetThirdPartyID(res.Res.BatchId).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除推送消息
func (b *pushmessageBiz) Delete(req *definition.PushmessageDeleteReq) error {
	old, err := b.orm.Get(context.Background(), req.ID)
	if err != nil {
		return err
	}
	// 删除推送消息
	if _, err = mob.Drop(old.ThirdPartyID); err != nil {
		return err
	}
	// 删除数据库消息
	err = b.orm.DeleteOneID(req.ID).Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// Get 获取推送消息
func (b *pushmessageBiz) Get(req *definition.PushmessageGetReq) (*definition.PushmessageDetail, error) {
	data, err := b.orm.Get(context.Background(), req.ID)
	if err != nil {
		return nil, err
	}
	return definition.ToDetail(data), nil
}

// List 分页获取推送消息列表
func (b *pushmessageBiz) List(req *definition.PushmessageListReq) (*model.PaginationRes, error) {
	query := b.orm.Query().Order(ent.Desc(pushmessage.FieldCreatedAt))
	// 关键字
	if req.Keyword != nil {
		query.Where(pushmessage.TitleContains(*req.Keyword))
	}
	// 推送状态
	if req.MessageStatus != nil {
		query.Where(pushmessage.MessageStatusEQ(*req.MessageStatus))
	}
	// 推送时间
	if req.PushBeginTime != nil {
		query.Where(pushmessage.PushTimeGTE(*fmtPushTime(req.PushBeginTime)))
	}
	if req.PushEndTime != nil {
		query.Where(pushmessage.PushTimeLTE(*fmtPushTime(req.PushEndTime)))
	}
	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Pushmessage) *definition.PushmessageDetail {
		return definition.ToDetail(item)
	}), nil
}

// fmtPushTime 格式化推送时间
func fmtPushTime(t *string) *time.Time {
	if t == nil {
		now := time.Now()
		return &now
	}
	tm, _ := time.Parse("2006-01-02 15:04:05", *t)
	return &tm
}
