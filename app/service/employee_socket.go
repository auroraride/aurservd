// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/pkg/cache"
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

func (s *employeeSocketService) Connect(token string) (uint64, error) {
    id, _ := cache.Get(context.Background(), token).Uint64()
    emr, _ := ar.Ent.Employee.QueryNotDeleted().Where(employee.ID(id)).WithStore().First(s.ctx)
    if emr == nil {
        return 0, errors.New("店员未找到")
    }

    eet := emr.Edges.Store
    if eet == nil {
        return 0, errors.New("店员未上班")
    }

    return id, nil
}
