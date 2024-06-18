package service

import (
	"net/url"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"

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
	*BaseService
}

func NewPromotionWithdrawalService(params ...any) *promotionWithdrawalService {
	return &promotionWithdrawalService{
		BaseService: newService(params...),
	}
}

// List 提现列表
func (s *promotionWithdrawalService) List(ctx echo.Context, req *promotion.WithdrawalListReq) *model.PaginationRes {
	q := ent.Database.PromotionWithdrawal.Query().WithCards().WithMember().Order(ent.Desc(promotionwithdrawal.FieldCreatedAt))
	s.FilterWithdrawal(q, req)

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.PromotionWithdrawal) promotion.WithdrawalListRes {
			res := promotion.WithdrawalListRes{
				WithdrawalDetail: promotion.WithdrawalDetail{
					ID:          item.ID,
					ApplyAmount: item.ApplyAmount,
					Amount:      item.Amount,
					Fee:         item.Fee,
					Tax:         item.Tex,
					Method:      promotion.WithdrawalMethod(item.Method).String(),
					Status:      item.Status,
					Remark:      item.Remark,
					CreatedAt:   item.CreatedAt.Format(carbon.DateTimeLayout),
					ApplyTime:   item.ApplyTime.Format(carbon.DateTimeLayout),
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

				res.BankCard.CardNo = account
				if item.Method == promotion.WithdrawalMethodBank.Value() && len(account) > 4 && ctx.Path() != "/manager/v1/promotion/withdrawal" {
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

func (s *promotionWithdrawalService) FilterWithdrawal(q *ent.PromotionWithdrawalQuery, req *promotion.WithdrawalListReq) {

	if req.ID != nil {
		q.Where(promotionwithdrawal.MemberID(*req.ID))
	}

	if req.Account != nil {
		q.Where(promotionwithdrawal.HasCardsWith(promotionbankcard.CardNoContains(*req.Account)))
	}

	if req.Status != nil {
		q.Where(promotionwithdrawal.Status(*req.Status))
	}

	if req.Keyword != nil {
		q.Where(promotionwithdrawal.HasMemberWith(
			promotionmember.Or(
				promotionmember.NameContains(*req.Keyword),
				promotionmember.PhoneContains(*req.Keyword),
			),
		))
	}

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			promotionwithdrawal.CreatedAtGTE(start),
			promotionwithdrawal.CreatedAtLTE(end),
		)
	}
}

func (s *promotionWithdrawalService) TotalWithdrawal(req *promotion.WithdrawalListReq) promotion.TotalRes {
	var v []promotion.TotalRes
	q := ent.Database.PromotionWithdrawal.Query().Where(promotionwithdrawal.StatusNEQ(promotion.WithdrawalStatusFailed.Value()))
	s.FilterWithdrawal(q, req)
	err := q.Aggregate(ent.Sum(promotionwithdrawal.FieldApplyAmount)).Scan(s.ctx, &v)
	if err != nil {
		snag.Panic("查询总提现失败")
	}
	return v[0]
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
			SetFee(wf.WithdrawalFee).
			SetTex(wf.Taxable).
			SetMethod(promotion.WithdrawalMethodBank.Value()).
			SetStatus(promotion.WithdrawalStatusPending.Value()).
			SetApplyTime(time.Now()).
			Exec(s.ctx)
		if err != nil {
			snag.Panic("申请提现失败")
		}

		// 扣除余额
		err = tx.PromotionMember.UpdateOneID(mem.ID).SetBalance(tools.NewDecimal().Sub(mem.Balance, req.ApplyAmount)).Exec(s.ctx)
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
				if v.Edges.Member == nil {
					snag.Panic("会员不存在")
				}
				// 审批不通过 退回余额
				err = tx.PromotionMember.UpdateOneID(v.Edges.Member.ID).SetBalance(tools.NewDecimal().Sum(v.Edges.Member.Balance, v.ApplyAmount)).Exec(s.ctx)
				if err != nil {
					snag.Panic("审批不通过 退回余额失败")
				}
			}
			err = tx.PromotionWithdrawal.Update().Where(promotionwithdrawal.ID(v.ID)).SetStatus(req.Status).SetReviewTime(time.Now()).SetRemark(req.Remark).Exec(s.ctx)
			if err != nil {
				snag.Panic("审批失败")
			}
		}
		return
	})

}

// CalculateWithdrawalFee 计算提现费用
func (s *promotionWithdrawalService) CalculateWithdrawalFee(mem *ent.PromotionMember, req *promotion.WithdrawalAlterReq) promotion.WithdrawalFeeRes {
	bankCard, _ := ent.Database.PromotionBankCard.Query().Where(promotionbankcard.MemberID(mem.ID), promotionbankcard.ID(req.AccountID)).First(s.ctx)
	if bankCard == nil {
		snag.Panic("提现账户不存在")
	}

	var (
		taxableAmount float64 // 应税金额
		fee           float64 // 手续费
		tax           float64 // 税费
		transferFees  float64 // 转账手续费
	)

	dl := tools.NewDecimal()

	// 税费
	if req.ApplyAmount > promotion.TaxExemptAmount {
		taxableAmount = dl.Sub(req.ApplyAmount, promotion.TaxExemptAmount)
		tax = dl.Mul(taxableAmount, promotion.TaxRate)
	}

	// 转账手续费
	transferFees = promotion.FransferFee
	if bankCard.City == promotion.FeeExemptionCity && bankCard.Bank == promotion.FeeExemptionBank {
		transferFees = 0
	}

	// 手续费 = (申请金额 - 税费) * 手续费率 + 转账手续费
	fee = dl.Sum(dl.Mul(dl.Sub(req.ApplyAmount, tax), promotion.FeeRate), transferFees)
	// 实际到账金额 = 申请金额 - 手续费 - 税费
	amount := dl.Sub(dl.Sub(req.ApplyAmount, fee), tax)

	return promotion.WithdrawalFeeRes{
		ApplyAmount:    req.ApplyAmount,
		AmountReceived: amount,
		WithdrawalFee:  fee,
		Taxable:        tax,
	}
}

// Export 导出
func (s *promotionWithdrawalService) Export(req *promotion.WithdrawalExportReq) (string, string) {
	q := ent.Database.PromotionWithdrawal.Query().WithMember(func(q *ent.PromotionMemberQuery) {
		q.WithPerson()
	}).WithCards()

	if req.Status != nil {
		q.Where(promotionwithdrawal.Status(*req.Status))
	}

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			promotionwithdrawal.CreatedAtGTE(start),
			promotionwithdrawal.CreatedAtLTE(end),
		)
	} else {
		// 默认查询最近一个月的数据
		q.Where(
			promotionwithdrawal.CreatedAtGTE(carbon.Now().SubMonth().StdTime()),
			promotionwithdrawal.CreatedAtLTE(carbon.Now().StdTime()),
		)

	}

	if req.Account != nil {
		q.Where(promotionwithdrawal.HasCardsWith(promotionbankcard.CardNoContains(*req.Account)))
	}
	if req.Status != nil {
		q.Where(promotionwithdrawal.Status(*req.Status))
	}
	if req.Keyword != nil {
		q.Where(promotionwithdrawal.HasMemberWith(
			promotionmember.Or(
				promotionmember.NameContainsFold(*req.Keyword),
				promotionmember.PhoneContainsFold(*req.Keyword),
			),
		))
	}

	items, _ := q.All(s.ctx)
	if len(items) == 0 {
		snag.Panic("没有可导出的数据")
	}
	file1 := s.ExportList(items)

	// 待审核的提现
	items, _ = ent.Database.PromotionWithdrawal.Query().WithMember(func(q *ent.PromotionMemberQuery) {
		q.WithPerson()
	}).WithCards().Where(promotionwithdrawal.Status(promotion.WithdrawalStatusPending.Value())).All(s.ctx)

	file2 := s.ExportTex(items)
	file3 := s.ExportWithdrawal(items)

	files := []string{file1, file2, file3}
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

func (s *promotionWithdrawalService) ExportList(items []*ent.PromotionWithdrawal) string {
	var rows tools.ExcelItems
	title := []any{"姓名电话", "银行信息", "申请金额", "手续费", "税费", "提现金额", "提现状态", "申请时间", "审核时间", "备注"}
	rows = append(rows, title)
	for _, item := range items {
		row := []any{
			"",
			"",
			item.ApplyAmount,
			item.Fee,
			item.Tex,
			item.Amount,
			promotion.WithdrawalStatus(item.Status).String(),
			item.ApplyTime.Format(carbon.DateTimeLayout),
			"",
			item.Remark,
		}

		var name, phone, account, bank string
		if item.Edges.Member != nil {
			name = item.Edges.Member.Name
			phone = item.Edges.Member.Phone
		}

		if item.Edges.Cards != nil {
			account = item.Edges.Cards.CardNo
			bank = item.Edges.Cards.Bank
		}

		row[0] = name + "-" + phone
		row[1] = bank + "-" + account
		if item.ReviewTime != nil {
			row[8] = item.ReviewTime.Format(carbon.DateTimeLayout)
		}
		rows = append(rows, row)
	}
	tempFile := "/tmp/提现列表" + time.Now().Format(carbon.ShortDateLayout) + ".xlsx"
	tools.NewExcel(tempFile).AddValues(rows).Done()
	return tempFile

}
