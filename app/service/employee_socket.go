// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/gorilla/websocket"
    "sync"
)

var (
    employeeClients sync.Map
)

type employeeSocketService struct {
    ctx          context.Context
    employee     *ent.Employee
    employeeInfo *model.Employee
}

func NewEmployeeSocket() *employeeSocketService {
    return &employeeSocketService{
        ctx: context.Background(),
    }
}

func NewEmployeeSocketWithEmployee(e *ent.Employee) *employeeSocketService {
    s := NewEmployeeSocket()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    s.employeeInfo = &model.Employee{
        ID:    s.employee.ID,
        Name:  s.employee.Name,
        Phone: s.employee.Phone,
    }
    return s
}

func (s *employeeSocketService) DisConnect(conn *websocket.Conn) {
    employeeClients.Range(func(key, value any) bool {
        if _, ok := value.(*websocket.Conn); ok {
            employeeClients.Delete(key)
            return false
        }
        return true
    })
}

func (s *employeeSocketService) DisConnectByStoreID(storeID uint64) {
    client, ok := employeeClients.Load(storeID)
    if !ok {
        return
    }
    switch client.(type) {
    case *websocket.Conn:
        _ = client.(*websocket.Conn).Close()
        break
    }
}

func (s *employeeSocketService) Connect(conn *websocket.Conn, token string) *model.EmployeeSocketMessage {
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
    s.DisConnectByStoreID(eet.ID)
    employeeClients.Store(eet.ID, conn)

    return &model.EmployeeSocketMessage{
        Success: true,
        Message: "OK",
    }
}

// Send 发送消息
func (s *employeeSocketService) Send(storeID uint64, res *model.EmployeeSocketMessage) {
    client, ok := employeeClients.Load(storeID)
    if !ok {
        return
    }

    switch client.(type) {
    case *websocket.Conn:
        _ = client.(*websocket.Conn).WriteMessage(websocket.TextMessage, res.Bytes())
        break
    }
}
