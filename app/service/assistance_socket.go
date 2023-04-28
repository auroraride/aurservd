// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/socket"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/silk"
	"go.uber.org/zap"
)

type assistanceSocketService struct {
	ctx context.Context
}

func NewAssistanceSocket() *assistanceSocketService {
	return &assistanceSocketService{
		ctx: context.Background(),
	}
}

func (s *assistanceSocketService) Detail(ass *ent.Assistance) (message *model.AssistanceSocketMessage, err error) {
	st := ass.Edges.Store
	if st == nil && ass.StoreID != nil {
		st, _ = ent.Database.Store.QueryNotDeleted().Where(store.ID(*ass.StoreID)).First(s.ctx)
		if st == nil {
			err = errors.New("未找到门店")
			return
		}
	}

	e := ass.Edges.Employee
	if e == nil && ass.EmployeeID != nil {
		e, _ = ent.Database.Employee.QueryNotDeleted().Where(employee.ID(*ass.EmployeeID)).First(s.ctx)
		if e == nil {
			err = errors.New("未找到店员")
			return
		}
	}

	message = &model.AssistanceSocketMessage{
		ID:        ass.ID,
		Status:    ass.Status,
		Breakdown: ass.Breakdown,
		Rider: model.LngLat{
			Lng: ass.Lng,
			Lat: ass.Lat,
		},
	}

	if st != nil {
		message.Store = &model.StoreLngLat{
			Store: model.Store{
				ID:   st.ID,
				Name: st.Name,
			},
			Lng: st.Lng,
			Lat: st.Lat,
		}
	}

	if e != nil {
		message.Employee = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
	}

	if ass.Status == model.AssistanceStatusPending {
		message.Seconds = int(time.Now().Sub(ass.CreatedAt).Seconds())
	} else {
		message.Seconds = ass.Wait
	}

	return
}

// SendRider 发送消息给骑手端
func (s *assistanceSocketService) SendRider(riderID uint64, ass *ent.Assistance) {
	message, err := s.Detail(ass)
	if err != nil {
		zap.L().Error("发送消息给骑手失败", zap.Error(err))
		return
	}

	socket.SendMessage(NewRiderSocket(), riderID, &model.RiderSocketMessage{Assistance: message})
}

// SenderEmployee 发送消息给门店端
func (s *assistanceSocketService) SenderEmployee(employeeID uint64, ass *ent.Assistance) {
	socket.SendMessage(NewEmployeeSocket(), employeeID, &model.EmployeeSocketMessage{
		Speech:       "您有一条救援任务",
		AssistanceID: silk.Pointer(ass.ID),
	})
}
