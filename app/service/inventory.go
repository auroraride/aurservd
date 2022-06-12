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
)

type inventoryService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
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

func NewInventoryWithEmployee(e *model.Employee) *inventoryService {
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

func (s *inventoryService) List() (items []model.Inventory) {
    _ = s.orm.QueryNotDeleted().Select(inventory.FieldName, inventory.FieldPurchase, inventory.FieldTransfer, inventory.FieldCount).Scan(s.ctx, &items)
    if len(items) < 1 {
        items = make([]model.Inventory, 0)
    }
    return
}

func (s *inventoryService) Delete(req *model.InventoryDelete) {
    s.orm.SoftDelete().Where(inventory.Name(*req.Name)).SaveX(s.ctx)
}
