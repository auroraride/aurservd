// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"strconv"
	"strings"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type businessEmployeeService struct {
	*BaseService
}

func NewBusinessEmployee(params ...any) *businessEmployeeService {
	return &businessEmployeeService{
		BaseService: newService(params...),
	}
}

// Inactive 获取骑手待激活订阅详情
func (s *businessEmployeeService) Inactive(qr string) (*model.SubscribeActiveInfo, *ent.Subscribe) {
	if strings.HasPrefix(qr, "SUBSCRIBE:") {
		qr = strings.ReplaceAll(qr, "SUBSCRIBE:", "")
	}
	id, _ := strconv.ParseUint(strings.TrimSpace(qr), 10, 64)
	return NewBusinessRiderWithEmployee(s.entEmployee).Inactive(id)
}

// Active 激活订阅
func (s *businessEmployeeService) Active(em *ent.Employee, req *model.EmployeeAllocateCreateReq) (res model.AllocateCreateRes) {
	if s.entStore == nil {
		snag.Panic("当前未上班")
	}

	// 解析二维码
	subscribeID := req.SubscribeID
	if req.SubscribeID == nil {
		str := strings.ReplaceAll(*req.Qrcode, "SUBSCRIBE:", "")
		id, _ := strconv.ParseUint(str, 10, 64)
		subscribeID = silk.UInt64(id)
	}

	if subscribeID == nil {
		snag.Panic("请选择骑手订阅")
	}

	// 获取当前门店
	st := NewEmployee().QueryStoreX(em)
	storeID := &st.ID
	employeeID := &em.ID

	return NewAllocate().Create(&model.AllocateCreateParams{
		SubscribeID: subscribeID,
		StoreID:     storeID,
		EmployeeID:  employeeID,
		BatteryID:   req.BatteryID,
		EbikeParam: model.AllocateCreateEbikeParam{
			ID: req.EbikeID,
		},
	})
}

func (s *businessEmployeeService) UnSubscribe(req *model.UnsubscribeEmployeeReq) {
	NewBusinessRiderWithEmployee(s.entEmployee).UnSubscribe(&model.BusinessSubscribeReq{ID: req.SubscribeID}, func(sub *ent.Subscribe) {
		if sub.BrandID != nil {
			if req.Qrcode == "" {
				snag.Panic("必须提交车辆信息")
			}
			// 校验车辆
			bike, _ := ent.Database.Ebike.Query().Where(
				ebike.RiderID(sub.RiderID),
				ebike.Or(
					ebike.Sn(req.Qrcode),
					ebike.Plate(req.Qrcode),
				),
			).First(s.ctx)
			if bike == nil {
				snag.Panic("未找到对应车辆")
			}
		}
	})
}
