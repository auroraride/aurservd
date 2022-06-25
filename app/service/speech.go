// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-22
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
)

type speechService struct {
    ctx context.Context
}

func NewSpeech() *speechService {
    return &speechService{
        ctx: context.Background(),
    }
}

func (s *speechService) SendSpeech(storeID uint64, message string) {
    res := &model.EmployeeSocketMessage{
        Success: true,
        Message: "OK",
        Speech:  message,
    }

    NewEmployeeSocket().Send(storeID, res)
}