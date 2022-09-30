// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/inventory"
    "github.com/auroraride/aurservd/internal/ent/predicate"
)

type inventoryService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    orm          *ent.InventoryClient
    employeeInfo *model.Employee
}

func NewInventory() *inventoryService {
    return &inventoryService{
        ctx: context.Background(),
        orm: ent.Database.Inventory,
    }
}

func NewInventoryWithRider(r *ent.Rider) *inventoryService {
    s := NewInventory()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewInventoryWithModifier(m *model.Modifier) *inventoryService {
    s := NewInventory()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewInventoryWithEmployee(e *ent.Employee) *inventoryService {
    s := NewInventory()
    if e != nil {
        s.employee = e
        s.employeeInfo = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
        s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    }
    return s
}

func (s *inventoryService) CreateOrModify(req *model.Inventory) {
    s.orm.Create().
        SetName(req.Name).
        SetCount(req.Count).
        SetPurchase(req.Purchase).
        SetTransfer(req.Transfer).
        OnConflictColumns(inventory.FieldName).
        UpdateNewValues().
        ExecX(s.ctx)
}

func (s *inventoryService) List(params ...model.InventoryListReq) (items []model.Inventory) {
    q := s.orm.QueryNotDeleted().
        Order(ent.Asc(inventory.FieldCreatedAt)).
        Select(inventory.FieldName, inventory.FieldPurchase, inventory.FieldTransfer, inventory.FieldCount)
    if len(params) > 0 {
        req := params[0]
        var or []predicate.Inventory
        if req.Count {
            or = append(or, inventory.Count(true))
        }
        if req.Transfer {
            or = append(or, inventory.Transfer(true))
        }
        if req.Purchase {
            or = append(or, inventory.Purchase(true))
        }
        q.Where(inventory.Or(or...))
    }
    _ = q.Scan(s.ctx, &items)
    if len(items) < 1 {
        items = make([]model.Inventory, 0)
    }
    return
}

func (s *inventoryService) Delete(req *model.InventoryDelete) {
    s.orm.SoftDelete().Where(inventory.Name(*req.Name)).SaveX(s.ctx)
}

// ListInventory 获取物资列表
func (s *inventoryService) ListInventory(req model.InventoryListReq) (items []model.InventoryItem) {
    // 电池型号列表
    bs := NewBatteryModel()
    for _, v := range bs.Models() {
        items = append(items, model.InventoryItem{
            Name:    v,
            Battery: true,
            Model:   v,
        })
    }
    for _, item := range s.List(req) {
        items = append(items, model.InventoryItem{Name: item.Name})
    }
    return
}

func (s *inventoryService) ListStockInventory(id uint64, req model.InventoryListReq) (res []model.InventoryNum) {
    inm := NewStock().StoreCurrentMap(id)
    items := s.ListInventory(req)
    res = make([]model.InventoryNum, len(items))
    for i, item := range items {
        res[i].InventoryItem = item
        if x, ok := inm[item.Name]; ok {
            res[i].Num = x.Num
        }
    }
    return
}
