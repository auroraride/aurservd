// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-22
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/gorilla/websocket"
    "sync"
)

var (
    speechers sync.Map
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
    speecher, ok := speechers.Load(storeID)
    if !ok {
        return
    }

    res := &model.EmployeeSocketMessage{
        Success: true,
        Message: "OK",
        Speech:  message,
    }

    switch speecher.(type) {
    case *websocket.Conn:
        speecher.(*websocket.Conn).WriteMessage(websocket.TextMessage, res.Bytes())
        break
    }
}

func (s *speechService) SpeecherDisConnect(conn *websocket.Conn) {
    speechers.Range(func(key, value any) bool {
        if _, ok := value.(*websocket.Conn); ok {
            speechers.Delete(key)
            return false
        }
        return true
    })
}

func (s *speechService) SpeecherDisConnectByStoreID(storeID uint64) {
    speecher, ok := speechers.Load(storeID)
    if !ok {
        return
    }
    switch speecher.(type) {
    case *websocket.Conn:
        speecher.(*websocket.Conn).Close()
        break
    }
}

func (s *speechService) SpeecherConnect(conn *websocket.Conn, token string) *model.EmployeeSocketMessage {
    id, _ := cache.Get(context.Background(), token).Uint64()
    emr, _ := ar.Ent.Employee.QueryNotDeleted().Where(employee.ID(id)).WithStore().First(s.ctx)
    if emr == nil {
        return &model.EmployeeSocketMessage{
            Success: false,
            Message: "店员未找到",
        }
    }

    eet := emr.Edges.Store
    if eet == nil {
        return &model.EmployeeSocketMessage{
            Success: false,
            Message: "店员未上班",
        }
    }

    // 存储连接信息
    s.SpeecherDisConnectByStoreID(eet.ID)
    speechers.Store(eet.ID, conn)

    return &model.EmployeeSocketMessage{
        Success: true,
        Message: "OK",
    }
}
