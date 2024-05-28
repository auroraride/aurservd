// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/lithammer/shortuuid/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/snag"
)

type storeService struct {
	ctx          context.Context
	orm          *ent.StoreClient
	employee     *ent.Employee
	modifier     *model.Modifier
	employeeInfo *model.Employee
}

func NewStore() *storeService {
	return &storeService{
		ctx: context.Background(),
		orm: ent.Database.Store,
	}
}

func NewStoreWithModifier(m *model.Modifier) *storeService {
	s := NewStore()
	s.modifier = m
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	return s
}

func NewStoreWithEmployee(e *ent.Employee) *storeService {
	s := NewStore()
	if e != nil {
		s.employee = e
		s.employeeInfo = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
		s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
	}
	return s
}

func (s *storeService) Query(id uint64) *ent.Store {
	item, _ := s.orm.QueryNotDeleted().WithEmployee().Where(store.ID(id)).First(s.ctx)
	if item == nil {
		snag.Panic("未找到有效门店")
	}
	return item
}

func (s *storeService) QuerySn(sn string) *ent.Store {
	if strings.HasPrefix(sn, "STORE:") {
		sn = strings.ReplaceAll(sn, "STORE:", "")
	}
	item, err := s.orm.QueryNotDeleted().WithEmployee().Where(store.Sn(sn)).First(s.ctx)
	if err != nil {
		snag.Panic("未找到有效门店")
	}
	return item
}

// Create 创建门店
func (s *storeService) Create(req *model.StoreCreateReq) model.StoreItem {
	b := NewBranch().Query(*req.BranchID)

	item := s.orm.Create().
		SetName(*req.Name).
		SetStatus(req.Status).
		SetBranch(b).
		SetCityID(b.CityID).
		SetSn(shortuuid.New()).
		SetLng(b.Lng).
		SetLat(b.Lat).
		SetAddress(b.Address).
		SetEbikeRepair(req.EbikeRepair).
		SetEbikeObtain(req.EbikeObtain).
		SetEbikeSale(req.EbikeSale).
		SetBusinessHours(req.BusinessHours).
		SetRest(req.Rest).
		SetPhotos(req.Photos).
		SaveX(s.ctx)

	if len(req.Materials) > 0 {
		for _, m := range req.Materials {
			tf := &model.StockTransferReq{
				OutboundID: 0,
				InboundID:  item.ID,
				Num:        m.Num,
			}
			if m.Model != "" {
				tf.Model = m.Model
			} else {
				tf.Name = m.Name
			}
			NewStockWithModifier(s.modifier).Transfer(tf)
		}
	}

	return s.Detail(item.ID)
}

// Modify 修改门店
func (s *storeService) Modify(req *model.StoreModifyReq) model.StoreItem {
	item := s.Query(req.ID)
	q := s.orm.UpdateOne(item).
		SetNillableEbikeObtain(req.EbikeObtain).
		SetNillableEbikeRepair(req.EbikeRepair).
		SetNillableEbikeSale(req.EbikeSale).
		SetNillableRest(req.Rest)
	if req.Status != nil {
		q.SetStatus(*req.Status)
	}
	if req.Name != nil {
		q.SetName(*req.Name)
	}
	if req.BranchID != nil {
		b := NewBranch().Query(*req.BranchID)
		q.SetLng(b.Lng).
			SetLat(b.Lat).
			SetAddress(b.Address).
			SetBranchID(*req.BranchID).
			SetCityID(b.CityID)
	}
	if req.BusinessHours != nil {
		q.SetBusinessHours(*req.BusinessHours)
	}
	if req.Photos != nil {
		q.SetPhotos(*req.Photos)
	}
	q.SaveX(s.ctx)
	return s.Detail(item.ID)
}

// Detail 获取门店详情
func (s *storeService) Detail(id uint64) model.StoreItem {
	item, _ := s.orm.QueryNotDeleted().
		Where(store.ID(id)).
		WithEmployee().
		WithCity().
		First(s.ctx)
	if item == nil {
		snag.Panic("未找到有效门店")
	}
	city := item.Edges.City
	res := model.StoreItem{
		ID:            item.ID,
		Name:          item.Name,
		Status:        item.Status,
		QRCode:        fmt.Sprintf("STORE:%s", item.Sn),
		EbikeRepair:   item.EbikeRepair,
		EbikeObtain:   item.EbikeObtain,
		EbikeSale:     item.EbikeSale,
		BusinessHours: item.BusinessHours,
		City: model.City{
			ID:   city.ID,
			Name: city.Name,
		},
	}
	if item.Edges.Employee != nil {
		ee := item.Edges.Employee
		res.Employee = &model.Employee{
			ID:    ee.ID,
			Name:  ee.Name,
			Phone: ee.Phone,
		}
	}
	return res
}

// Delete 删除门店
func (s *storeService) Delete(req *model.IDParamReq) {
	item := s.Query(req.ID)
	s.orm.UpdateOne(item).SetDeletedAt(time.Now()).ClearEmployeeID().SaveX(s.ctx)
}

// List 列举门店
func (s *storeService) List(req *model.StoreListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().WithCity()
	if req.CityID != nil {
		q.Where(store.CityID(*req.CityID))
	}
	if req.Name != nil {
		q.Where(store.NameContainsFold(*req.Name))
	}
	if req.Status != nil {
		q.Where(store.Status(*req.Status))
	}

	if req.BusinessType != nil {
		switch *req.BusinessType {
		case model.StoreBusinessTypeObtain: // 租车
			q.Where(store.EbikeObtain(true))
		case model.StoreBusinessTypeRepair:
			// 维修
			q.Where(store.EbikeRepair(true))
		case model.StoreBusinessTypeSale:
			// 买车
			q.Where(store.EbikeSale(true))
		case model.StoreBusinessTypeRest:
			// 驿站
			q.Where(store.Rest(true))
		}
	}

	return model.ParsePaginationResponse[model.StoreItem, ent.Store](q, req.PaginationReq, func(item *ent.Store) (res model.StoreItem) {
		_ = copier.Copy(&res, item)
		city := item.Edges.City
		res.City = model.City{
			ID:   city.ID,
			Name: city.Name,
		}
		res.QRCode = fmt.Sprintf("STORE:%s", item.Sn)
		return
	})
}

func (s *storeService) SwitchStatus(req *model.StoreSwtichStatusReq) {
	st := s.employee.Edges.Store

	if st == nil {
		snag.Panic("当前未上班")
	}

	if req.Status != model.StoreStatusOpen && req.Status != model.StoreStatusClose {
		snag.Panic("状态错误")
	}
	_, err := st.Update().SetStatus(req.Status).Save(s.ctx)
	if err != nil {
		snag.Panic(err)
	}
}
