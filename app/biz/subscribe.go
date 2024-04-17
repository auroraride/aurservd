package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
)

type subscribeBiz struct {
	orm *ent.SubscribeClient
	ctx context.Context
}

func NewSubscribe() *subscribeBiz {
	return &subscribeBiz{
		orm: ent.Database.Subscribe,
		ctx: context.Background(),
	}
}

// StoreModify 车电套餐修改激活门店
func (s *subscribeBiz) StoreModify(r *ent.Rider, req *definition.SubscribeStoreModifyReq) (err error) {
	// 判定套餐是否是车电套餐
	sub, _ := s.orm.QueryNotDeleted().
		Where(subscribe.ID(
			req.SubscribeID),
			subscribe.RiderID(r.ID),
		).WithPlan().First(s.ctx)
	if sub == nil || sub.Edges.Plan == nil {
		return errors.New("未找到订阅,或套餐不存在")
	}

	// 订阅不处于未激活状态不能更改门店
	if sub.Status != model.SubscribeStatusInactive {
		return errors.New("订阅状态不正确")
	}

	if sub.ID == req.StoreID {
		return errors.New("门店未变更")
	}

	if sub.Edges.Plan.Type != model.PlanTypeEbikeWithBattery.Value() {
		return errors.New("非车电套餐")
	}
	_, err = s.orm.UpdateOneID(req.SubscribeID).
		SetStoreID(req.StoreID).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return
}

// SubscribeStatus 查询订阅是否激活
func (s *subscribeBiz) SubscribeStatus(r *ent.Rider, req *model.EnterpriseRiderSubscribeStatusReq) (res bool) {
	sub, _ := s.orm.QueryNotDeleted().
		Where(subscribe.ID(req.ID),
			subscribe.RiderID(r.ID),
		).First(s.ctx)
	if sub == nil {
		return false
	}
	if sub.StartAt != nil {
		return true
	}
	return false
}

// SubscribeStatus 查询订阅是否激活
func (s *subscribeBiz) SubscribeStatus(r *ent.Rider, req *model.EnterpriseRiderSubscribeStatusReq) (res bool) {
	sub, _ := s.orm.QueryNotDeleted().
		Where(subscribe.ID(req.ID),
			subscribe.RiderID(r.ID),
		).First(s.ctx)
	if sub == nil {
		return false
	}
	if sub.StartAt != nil {
		return true
	}
	return false
}

// SubscribeStatus 查询订阅是否激活
func (s *subscribeBiz) SubscribeStatus(r *ent.Rider, req *model.EnterpriseRiderSubscribeStatusReq) (res bool) {
	sub, _ := s.orm.QueryNotDeleted().
		Where(subscribe.ID(req.ID),
			subscribe.RiderID(r.ID),
		).First(s.ctx)
	if sub == nil {
		return false
	}
	if sub.StartAt != nil {
		return true
	}
	return false
}
