// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/socket"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/assistance"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/tools"
)

type employeeSocketService struct {
    ctx context.Context
}

func NewEmployeeSocket() *employeeSocketService {
    return &employeeSocketService{
        ctx: context.Background(),
    }
}

func (s *employeeSocketService) Prefix() string {
    return "EMPLOYEE"
}

func (s *employeeSocketService) Connect(hub *socket.WebsocketHub, token string) (uint64, error) {
    id, _ := cache.Get(context.Background(), token).Uint64()
    emr, _ := ent.Database.Employee.QueryNotDeleted().Where(employee.ID(id)).WithStore().First(s.ctx)
    if emr == nil {
        return 0, errors.New("店员未找到")
    }

    // 查询最近的救援订单
    ass, _ := ent.Database.Assistance.QueryNotDeleted().Where(assistance.EmployeeID(id), assistance.Status(model.AssistanceStatusAllocated)).First(s.ctx)
    if ass != nil {
        hub.SendMessage(&model.EmployeeSocketMessage{
            Speech:       "您有一条救援任务正在进行中",
            AssistanceID: tools.NewPointerInterface(ass.ID),
        })
    }

    return id, nil
}
