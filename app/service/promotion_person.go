package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionperson"
	"github.com/auroraride/aurservd/pkg/snag"
)

type promotionPersonService struct {
	ctx context.Context
	*BaseService
}

func NewPromotionPersonService(params ...any) *promotionPersonService {
	return &promotionPersonService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
}

// RealNameAuth 实名认证
func (s *promotionPersonService) RealNameAuth(mem *ent.PromotionMember, req *promotion.RealNameAuthReq) promotion.RealNameAuthRes {
	per := mem.Edges.Person
	if per != nil && per.Status == promotion.PersonAuthenticated.Value() {
		snag.Panic("已实名认证")
	}

	success := promotion.RealNameAuthRes{Success: true}
	v := ali.NewIdVerify().Verify(req.Name, req.IdCard)
	if v == nil {
		success.Success = false
		return success
	}

	status := promotion.PersonAuthenticated.Value()
	if v.Code != "0" || v.Result.Res != "1" {
		status = promotion.PersonAuthenticationFailed.Value()
		success.Success = false
	}
	// 会员原有的认证信息
	if per != nil {
		// 修改认证信息
		ent.Database.PromotionPerson.UpdateOneID(per.ID).
			SetName(v.Result.Name).
			SetIDCardNumber(v.Result.IDCard).
			SetAddress(v.Result.Address).
			SetStatus(status).
			ExecX(s.ctx)

		// 修改会员信息
		ent.Database.PromotionMember.UpdateOne(mem).
			SetName(v.Result.Name).
			ExecX(s.ctx)
	}

	// 查询是否有重复的身份证号
	per, _ = ent.Database.PromotionPerson.Query().Where(promotionperson.IDCardNumber(v.Result.IDCard), promotionperson.Name(v.Result.Name)).First(s.ctx)

	if per == nil {
		// 创建新的认证信息
		per = ent.Database.PromotionPerson.Create().
			SetName(v.Result.Name).
			SetIDCardNumber(v.Result.IDCard).
			SetAddress(v.Result.Address).
			SetStatus(status).
			SaveX(s.ctx)
	}

	// 修改会员信息
	ent.Database.PromotionMember.UpdateOne(mem).
		SetPersonID(per.ID).
		SetName(v.Result.Name).
		ExecX(s.ctx)
	return success
}
