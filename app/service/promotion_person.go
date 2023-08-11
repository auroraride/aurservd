package service

import (
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionperson"
	"github.com/auroraride/aurservd/pkg/snag"
)

type promotionPersonService struct {
	*BaseService
}

func NewPromotionPersonService(params ...any) *promotionPersonService {
	return &promotionPersonService{
		BaseService: newService(params...),
	}
}

// RealNameAuth 实名认证
func (s *promotionPersonService) RealNameAuth(mem *ent.PromotionMember, req *promotion.RealNameAuthReq) promotion.RealNameAuthRes {
	per := mem.Edges.Person
	if per != nil && per.Status == promotion.PersonAuthenticated.Value() {
		snag.Panic("已实名认证")
	}

	// 验证身份证号
	v := ali.NewIdVerify().Verify(req.Name, req.IdCard)
	if v == nil || v.Code != "0" || v.Result.Res != "1" {
		return promotion.RealNameAuthRes{Success: false}
	}

	// 查询是否有重复的身份证号
	per, _ = ent.Database.PromotionPerson.Query().Where(promotionperson.IDCardNumber(v.Result.IDCard), promotionperson.Name(v.Result.Name)).First(s.ctx)

	if per == nil {
		// 创建新的认证信息
		per = ent.Database.PromotionPerson.Create().
			SetName(v.Result.Name).
			SetIDCardNumber(v.Result.IDCard).
			SetAddress(v.Result.Address).
			SetStatus(promotion.PersonAuthenticated.Value()).
			SaveX(s.ctx)
	}

	// 修改会员信息
	ent.Database.PromotionMember.UpdateOne(mem).
		SetPersonID(per.ID).
		SetName(v.Result.Name).
		ExecX(s.ctx)

	return promotion.RealNameAuthRes{Success: true}
}
