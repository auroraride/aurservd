package service

import (
	"context"
	"strings"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotionbankcard"
	"github.com/auroraride/aurservd/pkg/snag"
)

type promotionBankCardService struct {
	ctx context.Context
}

func NewPromotionBankCardService() *promotionBankCardService {
	return &promotionBankCardService{
		ctx: context.Background(),
	}
}

// List 获取银行卡列表
func (s *promotionBankCardService) List(mem *ent.PromotionMember) (res []*promotion.BankCardRes) {
	res = make([]*promotion.BankCardRes, 0)
	list, _ := ent.Database.PromotionBankCard.Query().Where(promotionbankcard.MemberID(mem.ID)).Order(ent.Desc(promotionbankcard.FieldCreatedAt)).All(s.ctx)
	if len(list) == 0 {
		snag.Panic("获取银行卡列表失败")
	}
	for _, item := range list {
		data := &promotion.BankCardRes{
			ID:          item.ID,
			IsDefault:   item.IsDefault,
			BankLogoURL: item.BankLogoURL,
			Bank:        item.Bank,
		}
		if item.CardNo != "" && len(item.CardNo) > 8 {
			data.CardNo = item.CardNo[len(item.CardNo)-4:]
		}
		res = append(res, data)
	}
	return
}

// Create 创建银行卡
func (s *promotionBankCardService) Create(mem *ent.PromotionMember, req *promotion.BankCardReq) {
	count := ent.Database.PromotionBankCard.Query().Where(promotionbankcard.MemberID(mem.ID)).CountX(s.ctx)
	if count > promotion.BankCardLimit {
		snag.Panic("银行卡数量已达上限")
	}

	// 判断是否已经存在
	bankCard, _ := ent.Database.PromotionBankCard.Query().Where(promotionbankcard.CardNo(req.CardNo), promotionbankcard.MemberID(mem.ID)).First(s.ctx)
	if bankCard != nil {
		snag.Panic("银行卡已经存在")
	}

	// 校验银行卡
	addr := ali.NewBankAddr().GetBankAddr(req.CardNo)
	if addr == nil || (addr != nil && addr.Data == nil) {
		snag.Panic("银行卡校验失败")
	}
	if addr.Data.Type != "借记卡" {
		snag.Panic("请使用借记卡")
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 将其他银行卡设置为非默认
		tx.PromotionBankCard.Update().Where(promotionbankcard.MemberID(mem.ID)).SetIsDefault(false).ExecX(s.ctx)

		tx.PromotionBankCard.Create().
			SetCardNo(req.CardNo).
			SetIsDefault(true).
			SetMember(mem).
			SetBankLogoURL(addr.Data.Logo).
			SetBank(addr.Data.Bank).
			SetProvince(addr.Data.Province).
			SetCity(addr.Data.City).
			SaveX(s.ctx)
		return
	})

}

// Update 修改银行卡默认状态
func (s *promotionBankCardService) Update(mem *ent.PromotionMember, req *model.IDParamReq) {
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 将其他银行卡设置为非默认
		tx.PromotionBankCard.Update().Where(promotionbankcard.MemberID(mem.ID)).SetIsDefault(false).ExecX(s.ctx)
		// 将当前银行卡设置为默认
		tx.PromotionBankCard.Update().Where(promotionbankcard.ID(req.ID)).SetIsDefault(true).ExecX(s.ctx)
		return
	})
}

// EncryptCardNo 加密银行卡号
func (s *promotionBankCardService) EncryptCardNo(cardNo string) string {
	if len(cardNo) < 8 {
		return cardNo
	}
	return cardNo[0:4] + strings.Repeat("*", len(cardNo)-8) + cardNo[len(cardNo)-4:]
}

// Delete 删除银行卡
func (s *promotionBankCardService) Delete(mem *ent.PromotionMember, req *model.IDParamReq) {
	ent.Database.PromotionBankCard.SoftDeleteOneID(req.ID).Where(promotionbankcard.MemberID(mem.ID)).ExecX(s.ctx)
}
