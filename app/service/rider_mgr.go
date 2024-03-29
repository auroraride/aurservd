// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type riderMgrService struct {
	ctx          context.Context
	modifier     *model.Modifier
	rider        *ent.Rider
	employee     *ent.Employee
	employeeInfo *model.Employee
}

func NewRiderMgr() *riderMgrService {
	return &riderMgrService{
		ctx: context.Background(),
	}
}

func NewRiderMgrWithRider(r *ent.Rider) *riderMgrService {
	s := NewRiderMgr()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, r)
	s.rider = r
	return s
}

func NewRiderMgrWithModifier(m *model.Modifier) *riderMgrService {
	s := NewRiderMgr()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func NewRiderMgrWithEmployee(e *ent.Employee) *riderMgrService {
	s := NewRiderMgr()
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

// Deposit 手动调整押金
func (s *riderMgrService) Deposit(req *model.RiderMgrDepositReq) {
	r := NewRider().Query(req.ID)
	o, _ := ent.Database.Order.QueryNotDeleted().
		Where(
			order.RiderID(req.ID),
			order.Status(model.OrderStatusPaid),
			order.Type(model.OrderTypeDeposit),
			order.DeletedAtIsNil(),
		).
		First(s.ctx)
	var before float64
	// 判断押金是否骑手自行缴纳
	if o != nil && o.Creator == nil {
		snag.Panic("用户有实际缴纳的押金订单, 无法继续修改")
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		if o != nil {
			before = o.Amount
			_, err = tx.Order.SoftDeleteOne(o).Save(s.ctx)
			snag.PanicIfError(err)
		}

		if req.Amount > 0 {
			_, err = tx.Order.Create().
				SetRiderID(req.ID).
				SetType(model.OrderTypeDeposit).
				SetStatus(model.OrderStatusPaid).
				SetRemark("管理员修改").
				SetAmount(req.Amount).
				SetTotal(req.Amount).
				SetPayway(model.OrderPaywayManual).
				SetOutTradeNo(tools.NewUnique().NewSN28()).
				SetTradeNo(tools.NewUnique().NewSN28()).
				Save(s.ctx)
			snag.PanicIfError(err)
		}
		return
	})

	// 记录日志
	go logging.NewOperateLog().
		SetRef(r).
		SetModifier(s.modifier).
		SetOperate(model.OperateDeposit).
		SetDiff(fmt.Sprintf("%.2f元", before), fmt.Sprintf("%.2f元", req.Amount)).
		Send()
}

// Modify 修改骑手资料
func (s *riderMgrService) Modify(req *model.RiderMgrModifyReq) {
	if req.Contact == nil && req.Phone == nil && req.AuthStatus == nil {
		snag.Panic("参数错误")
	}

	r := NewRiderWithModifier(s.modifier).Query(req.ID)
	var before, after []string

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {

		ru := tx.Rider.UpdateOne(r)

		if req.Phone != nil {
			if tx.Rider.QueryNotDeleted().Where(rider.Phone(*req.Phone)).ExistX(s.ctx) && r.Phone != *req.Phone {
				snag.Panic("电话已存在")
			}
			ru.SetPhone(*req.Phone)
			before = append(before, fmt.Sprintf("电话: %s", r.Phone))
			after = append(after, fmt.Sprintf("电话: %s", *req.Phone))
		}

		if req.Contact != nil {
			ru.SetContact(req.Contact)
			if r.Contact == nil {
				before = append(before, "联系人: 无")
			} else {
				before = append(before, fmt.Sprintf("联系人: %s, %s, %s", r.Contact.Relation, r.Contact.Phone, r.Contact.Name))
			}
			after = append(after, fmt.Sprintf("联系人: %s, %s, %s", req.Contact.Relation, req.Contact.Phone, req.Contact.Name))
		}

		if req.AuthStatus != nil {
			// 查询实名信息是否已存在
			p := r.Edges.Person
			if *req.AuthStatus != model.PersonUnauthenticated && req.Name == nil && req.IdCardNumber == nil && req.IdCardPortrait == nil && req.IdCardNational == nil {
				snag.Panic("修改实名状态时, 身份证信息不能为空")
			}
			var portrait, national *string
			if req.IdCardPortrait != nil {
				portrait = req.IdCardPortrait
				if !strings.HasPrefix(*portrait, "http") {
					portrait = silk.String("https://cdn.auroraride.com/" + *portrait)
				}
			}
			if req.IdCardNational != nil {
				national = req.IdCardNational
				if !strings.HasPrefix(*national, "http") {
					national = silk.String("https://cdn.auroraride.com/" + *national)
				}
			}
			// 更新或创建实名信息
			if p != nil {
				// 更新实名信息
				p.Update().
					SetNillableIDCardNumber(req.IdCardNumber).
					SetNillableIDCardPortrait(portrait).
					SetNillableIDCardNational(national).
					SetStatus(req.AuthStatus.Value()).
					SetNillableName(req.Name).
					SaveX(s.ctx)
			} else {
				// 创建实名信息
				p = ent.Database.Person.Create().
					SetNillableIDCardNumber(req.IdCardNumber).
					SetNillableIDCardPortrait(portrait).
					SetNillableIDCardNational(national).
					SetStatus(req.AuthStatus.Value()).
					SetName(*req.Name).
					SaveX(s.ctx)
				// 更新骑手实名信息
				ru.SetName(*req.Name).SetPersonID(p.ID)
			}
			before = append(before, fmt.Sprintf("认证状态: %s 身份证号: %s 正面: %s 国徽面:%s 姓名:%s ", model.PersonAuthStatus(p.Status).String(), p.IDCardNumber, p.IDCardPortrait, p.IDCardNational, p.Name))
			if *req.AuthStatus != model.PersonUnauthenticated {
				after = append(after, fmt.Sprintf("身份证号: %s 正面: %s 国徽面:%s 姓名: %s", *req.IdCardNumber, *req.IdCardPortrait, *req.IdCardNational, *req.Name))
			}
			after = append(after, fmt.Sprintf("认证状态: %s ", req.AuthStatus.String()))
		}
		_, err = ru.Save(s.ctx)
		snag.PanicIfError(err)
		return
	})

	// 记录日志
	go logging.NewOperateLog().
		SetRef(r).
		SetModifier(s.modifier).
		SetOperate(model.OperateProfile).
		SetDiff(strings.Join(before, "\n"), strings.Join(after, "\n")).
		Send()
}

// EmployeeQueryPhone 店员根据手机号查询骑手
func (s *riderMgrService) EmployeeQueryPhone(phone string) model.RiderEmployeeSearchRes {
	r, _ := ent.Database.Rider.QueryNotDeleted().WithPerson().Where(rider.Phone(phone)).WithEnterprise().First(s.ctx)
	if r == nil {
		snag.Panic("未找到骑手")
	}

	subd, _ := NewSubscribe().RecentDetail(r.ID)

	res := model.RiderEmployeeSearchRes{
		ID:              r.ID,
		Phone:           r.Phone,
		Overview:        NewExchange().Overview(r.ID),
		Status:          NewRider().Status(r),
		SubscribeStatus: subd.Status,
	}

	p := r.Edges.Person
	if p != nil {
		res.Name = p.Name
		res.AuthStatus = model.PersonAuthStatus(p.Status)
	}

	e := r.Edges.Enterprise
	if e != nil {
		res.Enterprise = &model.Enterprise{
			ID:    e.ID,
			Name:  e.Name,
			Agent: e.Agent,
		}
	}

	res.Plan = subd.Plan

	return res
}
