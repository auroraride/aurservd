// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-23, by aurb

package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/utils"
)

type employeeBiz struct {
	orm      *ent.EmployeeClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewEmployee() *employeeBiz {
	return &employeeBiz{
		orm: ent.Database.Employee,
		ctx: context.Background(),
	}
}

func NewEmployeeWithModifier(m *model.Modifier) *employeeBiz {
	b := NewEmployee()
	if m != nil {
		b.ctx = context.WithValue(b.ctx, model.CtxModifierKey{}, m)
		b.modifier = m
	}
	return b
}

// List 店员列表
func (b *employeeBiz) List(req *definition.EmployeeListReq) *model.PaginationRes {
	q := b.orm.QueryNotDeleted().
		WithCity().
		WithStore().
		WithStores().
		WithGroup()
	b.filter(q, req)

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Employee) definition.EmployeeListRes {
		return b.detail(item)
	})
}

func (b *employeeBiz) filter(q *ent.EmployeeQuery, req *definition.EmployeeListReq) {
	if req.Keyword != nil {
		q.Where(
			employee.Or(
				employee.NameContainsFold(*req.Keyword),
				employee.PhoneContainsFold(*req.Keyword),
				employee.HasStoresWith(store.NameContainsFold(*req.Keyword)),
			),
		)
	}

	if req.CityID != nil {
		q.Where(employee.CityID(*req.CityID))
	}

	if req.Status != nil {
		q.Where(employee.Enable(*req.Status))
	}
}

func (b *employeeBiz) detail(item *ent.Employee) definition.EmployeeListRes {
	res := definition.EmployeeListRes{
		ID:     item.ID,
		Name:   item.Name,
		Phone:  item.Phone,
		Enable: item.Enable,
		Limit:  item.Limit,
	}

	ec := item.Edges.City
	if ec != nil {
		res.City = model.City{
			ID:   ec.ID,
			Name: ec.Name,
		}
	}

	sto := item.Edges.Store
	if sto != nil {
		res.EmployeeStore = &definition.StoreInfo{
			ID:   sto.ID,
			Name: sto.Name,
		}
	}

	sts := item.Edges.Stores
	ests := make([]*definition.StoreInfo, 0)
	for _, st := range sts {
		ests = append(ests, &definition.StoreInfo{
			ID:   st.ID,
			Name: st.Name,
		})
	}
	res.Stores = ests

	gp := item.Edges.Group
	if gp != nil {
		res.Group = &definition.StoreGroup{
			ID:   gp.ID,
			Name: gp.Name,
		}
		stis := make([]*definition.StoreInfo, 0)
		all, _ := ent.Database.Store.QueryNotDeleted().Where(store.GroupID(gp.ID)).All(b.ctx)
		for _, s := range all {
			stis = append(stis, &definition.StoreInfo{
				ID:   s.ID,
				Name: s.Name,
			})
		}
		res.Stores = stis
	}

	return res
}

// Create 添加店员
func (b *employeeBiz) Create(req *definition.EmployeeCreateReq) (err error) {
	// 判断重复
	em := b.QueryByPhone(req.Phone)
	if em != nil {
		return errors.New("店员已存在")
	}
	password, _ := utils.PasswordGenerate(req.Password)
	emc := b.orm.Create().
		SetCityID(req.CityID).
		SetName(req.Name).
		SetPhone(req.Phone).
		SetPassword(password).
		SetLimit(req.Limit)

	if req.GroupID != 0 {
		emc.SetGroupID(req.GroupID)
	}
	if len(req.StoreIDs) != 0 {
		emc.AddStoreIDs(req.StoreIDs...)
	}

	em, _ = emc.Save(b.ctx)
	if em == nil {
		return errors.New("店员添加失败")
	}
	return nil
}

// QueryByPhone 根据phone查找店员
func (b *employeeBiz) QueryByPhone(phone string) *ent.Employee {
	item, _ := b.orm.QueryNotDeleted().Where(employee.Phone(phone)).First(b.ctx)
	return item
}

// Modify 修改店员
func (b *employeeBiz) Modify(req *definition.EmployeeModifyReq) error {
	emu := b.orm.UpdateOneID(req.ID)

	if req.CityID != nil {
		emu.SetCityID(*req.CityID)
	}

	if req.Name != nil {
		emu.SetName(*req.Name)
	}

	if req.Phone != nil {
		em := b.QueryByPhone(*req.Phone)
		if em != nil && em.ID != req.ID {
			return errors.New("手机号已存在")
		}
		emu.SetPhone(*req.Phone)
	}

	if req.Limit != nil {
		emu.SetLimit(*req.Limit)
	}

	if req.Password != nil {
		password, _ := utils.PasswordGenerate(*req.Password)
		emu.SetPassword(password)
	}
	if req.GroupID != nil {
		emu.SetGroupID(*req.GroupID)
	}
	if len(req.StoreIDs) != 0 {
		emu.AddStoreIDs(req.StoreIDs...)
	}
	if req.Enable != nil {
		emu.SetEnable(*req.Enable)
	}

	nem, _ := emu.Save(b.ctx)
	if nem == nil {
		return errors.New("店员更新失败")
	}
	return nil
}

// Delete 删除店员
func (b *employeeBiz) Delete(id uint64) error {
	_, err := b.orm.SoftDeleteOneID(id).Save(b.ctx)
	if err != nil {
		return errors.New("店员删除失败")
	}
	return nil
}
