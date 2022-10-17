// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-22
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/socket"
)

type speechService struct {
    ctx context.Context
}

func NewSpeech() *speechService {
    return &speechService{
        ctx: context.Background(),
    }
}

// SendSpeech 发送播报内容
func (s *speechService) SendSpeech(employeeID uint64, message string) {
    res := &model.EmployeeSocketMessage{
        Speech: message,
    }
    socket.SendMessage(NewEmployeeSocket(), employeeID, res)
}
