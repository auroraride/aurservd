package service

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/internal/ent/promotionbankcard"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionwithdrawal"
	"github.com/auroraride/aurservd/pkg/zip"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type promotionWithdrawalService struct {
	ctx context.Context
	*BaseService
}

func NewPromotionWithdrawalService(params ...any) *promotionWithdrawalService {
	return &promotionWithdrawalService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
}

// List 提现列表
func (s *promotionWithdrawalService) List(req *promotion.WithdrawalListReq) *model.PaginationRes {
	q := ent.Database.PromotionWithdrawal.Query().WithCards().WithMember().Order(ent.Asc(promotionwithdrawal.FieldStatus), ent.Desc(promotionwithdrawal.FieldCreatedAt))

	if req.ID != nil {
		q.Where(promotionwithdrawal.MemberID(*req.ID))
	}

	if req.Account != nil {
		q.Where(promotionwithdrawal.HasCardsWith(promotionbankcard.CardNo(*req.Account)))
	}

	if req.Status != nil {
		q.Where(promotionwithdrawal.Status(*req.Status))
	}

	if req.Keywork != nil {
		q.Where(promotionwithdrawal.HasMemberWith(
			promotionmember.Or(
				promotionmember.NameContains(*req.Keywork),
				promotionmember.PhoneContains(*req.Keywork),
			),
		))
	}

	if req.Status != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			promotionwithdrawal.CreatedAtGTE(start),
			promotionwithdrawal.CreatedAtLTE(end),
		)
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.PromotionWithdrawal) promotion.WithdrawalListRes {
			res := promotion.WithdrawalListRes{
				WithdrawalDetail: promotion.WithdrawalDetail{
					ID:        item.ID,
					Amount:    item.Amount,
					Method:    promotion.WithdrawalMethod(item.Method).String(),
					Status:    item.Status,
					Remark:    item.Remark,
					CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
					ApplyTime: item.ApplyTime.Format(carbon.DateTimeLayout),
				},
			}
			if item.ReviewTime != nil {
				res.ReviewTime = item.ReviewTime.Format(carbon.DateTimeLayout)
			}

			if item.Edges.Cards != nil {
				account := item.Edges.Cards.CardNo

				res.BankCard = &promotion.BankCardRes{
					BankLogoURL: item.Edges.Cards.BankLogoURL,
					Bank:        item.Edges.Cards.Bank,
				}
				if item.Method == promotion.WithdrawalMethodBank.Value() && len(account) > 4 {
					// 截取银行卡号后四位
					res.BankCard.CardNo = account[len(account)-4:]
				}
			}
			if item.Edges.Member != nil {
				res.MemberBaseInfo = &promotion.MemberBaseInfo{
					ID:    item.Edges.Member.ID,
					Name:  item.Edges.Member.Name,
					Phone: item.Edges.Member.Phone,
				}
			}
			return res
		},
	)
}

// Alter 申请提现
func (s *promotionWithdrawalService) Alter(mem *ent.PromotionMember, req *promotion.WithdrawalAlterReq) {

	if ent.Database.PromotionWithdrawal.Query().Where(promotionwithdrawal.MemberID(mem.ID), promotionwithdrawal.Status(promotion.WithdrawalStatusPending.Value())).CountX(s.ctx) > 0 {
		snag.Panic("存在未审批的提现申请")
	}

	if mem.Balance < req.ApplyAmount {
		snag.Panic("余额不足")
	}

	bankCard := mem.Edges.Cards
	if bankCard == nil {
		snag.Panic("请先绑定银行卡")
	}

	if mem.Edges.Person != nil && mem.Edges.Person.Status != promotion.PersonAuthenticated.Value() {
		snag.Panic("请先实名认证")
	}

	wf := s.CalculateWithdrawalFee(mem, req)

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		err = tx.PromotionWithdrawal.Create().
			SetMember(mem).
			SetAccountID(req.AccountID).
			SetApplyAmount(req.ApplyAmount).
			SetAmount(wf.AmountReceived).
			SetFee(wf.Taxable + wf.WithdrawalFee).
			SetMethod(promotion.WithdrawalMethodBank.Value()).
			SetStatus(promotion.WithdrawalStatusPending.Value()).
			SetApplyTime(time.Now()).
			Exec(s.ctx)
		if err != nil {
			snag.Panic("申请提现失败")
		}
		// 扣除余额
		err = tx.PromotionMember.UpdateOneID(mem.ID).AddBalance(-req.ApplyAmount).Exec(s.ctx)
		if err != nil {
			snag.Panic("申请提现失败")
		}
		return
	})
}

// AlterReview 审批提现
func (s *promotionWithdrawalService) AlterReview(req *promotion.WithdrawalApprovalReq) {
	mw := ent.Database.PromotionWithdrawal.Query().WithMember().
		Where(
			promotionwithdrawal.IDIn(req.IDs...),
			promotionwithdrawal.Status(promotion.WithdrawalStatusPending.Value()),
		).AllX(s.ctx)
	if len(mw) == 0 {
		snag.Panic("提现记录不存在")
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		for _, v := range mw {
			if req.Status == promotion.WithdrawalStatusFailed.Value() {
				// 审批不通过 退回余额
				err = tx.PromotionMember.UpdateOneID(v.Edges.Member.ID).AddBalance(v.ApplyAmount).SetRemark(req.Remark).Exec(s.ctx)
				if err != nil {
					snag.Panic("审批不通过 退回余额失败")
				}
			}
			err = tx.PromotionWithdrawal.Update().Where(promotionwithdrawal.ID(v.ID)).SetStatus(req.Status).SetReviewTime(time.Now()).Exec(s.ctx)
			if err != nil {
				snag.Panic("审批失败")
			}
		}
		return
	})

}

// roundToTwoDecimalPlaces 保留两位小数
func (s *promotionWithdrawalService) roundToTwoDecimalPlaces(value float64) float64 {
	cash, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return cash
}

// CalculateWithdrawalFee 计算提现费用
func (s *promotionWithdrawalService) CalculateWithdrawalFee(mem *ent.PromotionMember, req *promotion.WithdrawalAlterReq) promotion.WithdrawalFeeRes {
	bankCard := ent.Database.PromotionBankCard.Query().Where(promotionbankcard.MemberID(mem.ID), promotionbankcard.ID(req.AccountID)).FirstX(s.ctx)
	if bankCard == nil {
		snag.Panic("提现账户不存在")
	}

	var (
		taxableAmount float64 // 应税金额
		fee           float64 // 手续费
		tax           float64 // 税费
		transferFees  float64 // 转账手续费
	)

	// 税费
	if req.ApplyAmount > promotion.TaxExemptAmount {
		taxableAmount = req.ApplyAmount - promotion.TaxExemptAmount
		tax = taxableAmount * promotion.TaxRate
	}

	// 转账手续费
	transferFees = promotion.FransferFee
	if bankCard.City == promotion.FeeExemptionCity && bankCard.Bank == promotion.FeeExemptionBank {
		transferFees = 0
	}

	// 手续费
	fee = (req.ApplyAmount-tax)*promotion.FeeRate + transferFees
	// 实际到账金额
	amount := req.ApplyAmount - fee - tax

	return promotion.WithdrawalFeeRes{
		ApplyAmount:    req.ApplyAmount,
		AmountReceived: s.roundToTwoDecimalPlaces(amount),
		WithdrawalFee:  s.roundToTwoDecimalPlaces(fee),
		Taxable:        s.roundToTwoDecimalPlaces(tax),
	}
}

// Export 导出
func (s *promotionWithdrawalService) Export() (string, string) {
	items := ent.Database.PromotionWithdrawal.Query().WithMember(func(q *ent.PromotionMemberQuery) {
		q.WithPerson()
	}).WithCards().Where(promotionwithdrawal.Status(promotion.WithdrawalStatusPending.Value())).AllX(s.ctx)
	file1 := s.ExportTex(items)
	file2 := s.ExportWithdrawal(items)

	files := []string{file1, file2}
	output := "/tmp/转账代发_报税明细.zip"
	if err := zip.ZipFiles(output, files); err != nil {
		snag.Panic("压缩文件失败")
	}
	return output, url.QueryEscape(filepath.Base("转账代发_报税明细"))
}

// ExportTex 导出个税
func (s *promotionWithdrawalService) ExportTex(items []*ent.PromotionWithdrawal) string {
	var rows tools.ExcelItems
	title := []any{"工号", "姓名", "证照类型", "证照号码", "所得项目", "本期收入", "本期免税收入", "商业健康保险费", "税延养老保险费", "准予扣除的捐赠额", "其他扣除", "减免税额", "备注"}
	rows = append(rows, title)
	for _, item := range items {
		row := []any{
			"",
			"",
			"居民身份证",
			"",
			"劳务报酬",
			item.ApplyAmount,
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		}
		if item.Edges.Member != nil && item.Edges.Member.Edges.Person != nil {
			row[0] = item.Edges.Member.ID
			row[1] = item.Edges.Member.Edges.Person.Name
			row[3] = item.Edges.Member.Edges.Person.IDCardNumber
		}

		rows = append(rows, row)
	}
	tempFile := "/tmp/劳务报酬报税明细" + time.Now().Format(carbon.ShortDateLayout) + ".xlsx"
	tools.NewExcel(tempFile).AddValues(rows).Done()
	return tempFile
}

// ExportWithdrawal 转账代发明细
func (s *promotionWithdrawalService) ExportWithdrawal(items []*ent.PromotionWithdrawal) string {
	var rows tools.ExcelItems
	title := []any{"账号", "户名", "金额", "开户行", "开户地", "汇款备注"}
	rows = append(rows, title)
	for _, item := range items {
		row := []any{
			"",
			"",
			item.Amount,
			"",
			"",
			"推广提现",
		}
		if item.Edges.Cards != nil && item.Edges.Member != nil && item.Edges.Member.Edges.Person != nil {
			row[0] = item.Edges.Cards.CardNo
			row[1] = item.Edges.Member.Edges.Person.Name
			row[3] = item.Edges.Cards.Bank
			row[4] = item.Edges.Cards.City
		}
		rows = append(rows, row)
	}
	tempFile := "/tmp/转账代发明细" + time.Now().Format(carbon.ShortDateLayout) + ".xlsx"
	tools.NewExcel(tempFile).AddValues(rows).Done()
	return tempFile
}
