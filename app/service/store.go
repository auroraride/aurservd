// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/jinzhu/copier"
    "github.com/lithammer/shortuuid/v4"
    "strings"
    "time"
)

type storeService struct {
    ctx      context.Context
    orm      *ent.StoreClient
    employee *ent.Employee
    modifier *model.Modifier
}

func NewStore() *storeService {
    return &storeService{
        ctx: context.Background(),
        orm: ar.Ent.Store,
    }
}

func NewStoreWithModifier(m *model.Modifier) *storeService {
    s := NewStore()
    s.modifier = m
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    return s
}

func NewStoreWithEmployee(e *ent.Employee) *storeService {
    s := NewStore()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *storeService) Query(id uint64) *ent.Store {
    item, err := s.orm.QueryNotDeleted().WithEmployee().Where(store.ID(id)).Only(s.ctx)
    if err != nil {
        snag.Panic("未找到有效门店")
    }
    return item
}

func (s *storeService) QuerySn(sn string) *ent.Store {
    if strings.HasPrefix(sn, "STORE:") {
        sn = strings.ReplaceAll(sn, "STORE:", "")
    }
    item, err := s.orm.QueryNotDeleted().WithEmployee().Where(store.Sn(sn)).Only(s.ctx)
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
    q := s.orm.UpdateOne(item)
    if req.Status != nil {
        q.SetStatus(*req.Status)
    }
    if req.Name != nil {
        q.SetName(*req.Name)
    }
    if req.BranchID != nil {
        b := NewBranch().Query(*req.BranchID)
        q.SetBranchID(*req.BranchID).SetCityID(b.CityID)
    }
    q.SaveX(s.ctx)
    return s.Detail(item.ID)
}

// Detail 获取门店详情
// TODO 店员
func (s *storeService) Detail(id uint64) model.StoreItem {
    item, err := s.orm.QueryNotDeleted().
        Where(store.ID(id)).
        WithEmployee().
        WithCity().
        Only(s.ctx)
    if err != nil {
        snag.Panic("未找到有效门店")
    }
    city := item.Edges.City
    res := model.StoreItem{
        ID:     item.ID,
        Name:   item.Name,
        Status: item.Status,
        QRCode: fmt.Sprintf("STORE:%s", item.Sn),
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
