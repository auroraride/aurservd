// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-20
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/subscribereminder"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
)

type reminderService struct {
    *BaseService
    orm *ent.SubscribeReminderClient
}

func NewReminder(params ...any) *reminderService {
    return &reminderService{
        BaseService: newService(params...),
        orm:         ent.Database.SubscribeReminder,
    }
}

func (s *reminderService) List(req *model.ReminderListReq) *model.PaginationRes {
    q := s.orm.Query().Order(ent.Desc(subscribereminder.FieldCreatedAt))
    if req.RiderID != 0 {
        q.Where(subscribereminder.RiderID(req.RiderID))
    }
    if req.Keyword != "" {
        q.Where(subscribereminder.Or(
            subscribereminder.PhoneContainsFold(req.Keyword),
            subscribereminder.NameContainsFold(req.Keyword),
        ))
    }
    if req.Start != "" {
        q.Where(subscribereminder.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
    }
    if req.End != "" {
        q.Where(subscribereminder.CreatedAtLTE(tools.NewTime().ParseDateStringX(req.End)))
    }
    if req.Days != nil {
        q.Where(subscribereminder.Days(*req.Days))
    }
    if req.Type != "" {
        q.Where(subscribereminder.TypeEQ(subscribereminder.Type(req.Type)))
    }
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.SubscribeReminder) model.ReminderListRes {
        return model.ReminderListRes{
            Phone:      item.Phone,
            Name:       item.Name,
            Success:    item.Success,
            Time:       item.CreatedAt.Format(carbon.DateTimeLayout),
            PlanName:   item.PlanName,
            Fee:        item.Fee,
            FeeFormula: item.FeeFormula,
        }
    })
}
