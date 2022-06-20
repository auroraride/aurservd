// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/inventory"
    "github.com/auroraride/aurservd/internal/ent/predicate"
    "github.com/auroraride/aurservd/pkg/tools"
)

type inventoryService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.InventoryClient
}

func NewInventory() *inventoryService {
    return &inventoryService{
        ctx: context.Background(),
        orm: ar.Ent.Inventory,
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
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
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
    q := s.orm.QueryNotDeleted().Select(inventory.FieldName, inventory.FieldPurchase, inventory.FieldTransfer, inventory.FieldCount)
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
    bs := NewBattery()
    tp := tools.NewPointer()
    for _, v := range bs.ListVoltages() {
        items = append(items, model.InventoryItem{
            Name:    bs.VoltageName(v),
            Battery: true,
            Voltage: tp.Float64(v),
        })
    }
    for _, item := range s.List(req) {
        if item.Count {
            items = append(items, model.InventoryItem{Name: item.Name})
        }
    }
    return
}

func (s *inventoryService) ListStockInventory(id uint64, req model.InventoryListReq) (res []model.InventoryItemWithNum) {
    inm := NewStock().CurrentMap(id)
    items := s.ListInventory(req)
    res = make([]model.InventoryItemWithNum, len(items))
    for i, item := range items {
        res[i].InventoryItem = item
        if x, ok := inm[item.Name]; ok {
            res[i].Num = x.Num
        }
    }
    return
}
