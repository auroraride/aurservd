package biz

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/fault"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/setting"
	"github.com/auroraride/aurservd/pkg/tools"
)

type faultBiz struct {
	orm *ent.FaultClient
	ctx context.Context
}

func NewFaultBiz() *faultBiz {
	return &faultBiz{
		orm: ent.Database.Fault,
		ctx: context.Background(),
	}
}

// List 获取故障列表
func (s *faultBiz) List(req *definition.FaultListReq) (res *model.PaginationRes) {
	q := s.orm.QueryNotDeleted().Order(ent.Desc(fault.FieldCreatedAt)).WithRider().WithCity().WithCabinet().WithBattery().WithEbike()

	if req.Type != nil {
		q.Where(fault.TypeEQ(*req.Type))
	}

	if req.Status != nil {
		q.Where(fault.StatusEQ(*req.Status))
	}

	if req.Keyword != nil {
		q.Where(
			fault.HasRiderWith(
				rider.Or(
					rider.NameContains(*req.Keyword),
					rider.PhoneContains(*req.Keyword),
				),
			),
		)
	}

	if req.City != nil {
		q.Where(fault.HasCityWith(city.NameEQ(*req.City)))
	}

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(fault.CreatedAtGTE(start), fault.CreatedAtLTE(end))
	}

	if req.DeviceNo != nil {
		q.Where(
			fault.Or(
				fault.HasCabinetWith(cabinet.SerialContains(*req.DeviceNo)),
				fault.HasBatteryWith(battery.SnContains(*req.DeviceNo)),
				fault.HasEbikeWith(ebike.SnContains(*req.DeviceNo)),
			))
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Fault) (res *definition.Fault) {
		res = &definition.Fault{
			ID:          item.ID,
			Type:        definition.FaultType(item.Type),
			FaultCause:  item.Fault,
			Status:      item.Status,
			Description: item.Description,
			Attachments: item.Attachments,
			CreatedAt:   item.CreatedAt.Format(carbon.DateTimeLayout),
			Remark:      item.Remark,
		}
		if item.Edges.City != nil {
			res.City = item.Edges.City.Name
		}

		if item.Edges.Cabinet != nil {
			res.DeviceNo = item.Edges.Cabinet.Serial
		}
		if item.Edges.Battery != nil {
			res.DeviceNo = item.Edges.Battery.Sn
		}
		if item.Edges.Ebike != nil {
			res.DeviceNo = item.Edges.Ebike.Sn
		}
		if item.Edges.Rider != nil {
			res.Rider = model.Rider{
				ID:    item.Edges.Rider.ID,
				Phone: item.Edges.Rider.Phone,
				Name:  item.Edges.Rider.Name,
			}
		}
		return res
	})
}

// Create 创建故障
func (s *faultBiz) Create(r *ent.Rider, req *definition.FaultCreateReq) error {
	q := s.orm.Create().
		SetType(req.Type.Value()).
		SetFault(req.FaultCause).
		SetRider(r).
		SetDescription(req.Description).
		SetAttachments(req.Attachments).
		SetCityID(req.CityID)

	// 查找设备
	switch req.Type {
	case definition.FaultTypeBattery:
		d, err := ent.Database.Battery.QueryNotDeleted().Where(battery.SnEQ(req.DeviceNo)).First(s.ctx)
		if err != nil {
			return errors.New("电池不存在")
		}
		q.SetBatteryID(d.ID)
	case definition.FaultTypeCabinet:
		d, err := ent.Database.Cabinet.QueryNotDeleted().Where(cabinet.SerialEQ(req.DeviceNo)).First(s.ctx)
		if err != nil {
			return errors.New("电柜不存在")
		}
		q.SetCabinetID(d.ID)
	case definition.FaultTypeEbike:
		d, err := ent.Database.Ebike.Query().Where(ebike.SnEQ(req.DeviceNo)).First(s.ctx)
		if err != nil {
			return errors.New("车辆不存在")
		}
		q.SetEbikeID(d.ID)
	default:
		return errors.New("故障类型错误")
	}

	_, err := q.Save(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// ModifyStatus 修改故障状态
func (s *faultBiz) ModifyStatus(req *definition.FaultModifyStatusReq) (err error) {
	_, err = s.orm.UpdateOneID(req.ID).
		SetStatus(req.Status).
		SetRemark(req.Remark).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// FaultCause 故障原因
func (s *faultBiz) FaultCause() (items []definition.FaultCauseRes, err error) {
	res, _ := ent.Database.Setting.Query().Where(setting.KeyIn(
		"EBIKE_FAULT",
		"BATTERY_FAULT",
		"OTHER_FAULT",
		"CABINET_FAULT",
	)).All(s.ctx)
	if res == nil {
		return items, nil
	}

	for _, v := range res {
		var content []string
		err = json.Unmarshal([]byte(v.Content), &content)
		if err != nil {
			return nil, err
		}
		items = append(items, definition.FaultCauseRes{
			Key:   v.Key,
			Value: content,
		})
	}

	return items, nil
}
