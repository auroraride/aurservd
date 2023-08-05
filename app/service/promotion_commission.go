package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotionearnings"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/snag"
)

type promotionCommissionService struct {
	*BaseService
	ctx context.Context
}

func NewPromotionCommissionService(params ...any) *promotionCommissionService {
	return &promotionCommissionService{
		ctx:         context.Background(),
		BaseService: newService(params...),
	}
}

// Create 返佣方案创建
func (s *promotionCommissionService) Create(req *promotion.CommissionCreateReq) *ent.PromotionCommission {
	if *req.Type == promotion.CommissionDefault {
		snag.Panic("返佣方案不能为默认方案")
	}
	// 验证方案
	s.CheckCommissionRule(&req.Rule)

	// 自定义方案 一个会员只能有一个自定义方案
	if *req.Type == promotion.CommissionCustom && req.MemberID != nil {
		// 删除会员已有的自定义方案
		ent.Database.PromotionCommission.SoftDelete().Where(promotioncommission.MemberID(*req.MemberID)).ExecX(s.ctx)
	}

	// 创建方案
	return ent.Database.PromotionCommission.Create().
		SetName(req.Name).
		SetType(req.Type.Value()).
		SetRule(&req.Rule).
		SetNillableMemberID(req.MemberID).
		SetNillableDesc(req.Desc).
		SetStartAt(time.Now()).
		SaveX(s.ctx)

}

// Selection 返佣方案选择
func (s *promotionCommissionService) Selection() []promotion.CommissionSelection {
	commission := ent.Database.PromotionCommission.QueryNotDeleted().Where(promotioncommission.TypeNotIn(promotion.CommissionCustom.Value())).AllX(s.ctx)
	res := make([]promotion.CommissionSelection, 0, len(commission))
	for _, itme := range commission {
		res = append(res, promotion.CommissionSelection{
			ID:   itme.ID,
			Name: itme.Name,
		})
	}
	return res
}

// CommissionTaskSelection 返佣任务选择
func (s *promotionCommissionService) CommissionTaskSelection() []promotion.CommissionTaskSelect {
	res := make([]promotion.CommissionTaskSelect, 0, 4)

	for k, v := range promotion.CommissionRuleKeyNames {
		res = append(res, promotion.CommissionTaskSelect{
			Key:  k,
			Name: v,
		})
	}
	return res
}

// PromotionCommissionByID 通过id查询方案
func (s *promotionCommissionService) PromotionCommissionByID(id uint64) (*ent.PromotionCommission, error) {
	return ent.Database.PromotionCommission.QueryNotDeleted().Where(promotioncommission.ID(id)).First(s.ctx)
}

// DefaultPromotionCommission 获取默认方案
func (s *promotionCommissionService) DefaultPromotionCommission() (*ent.PromotionCommission, error) {
	return ent.Database.PromotionCommission.QueryNotDeleted().Where(promotioncommission.Type(promotion.CommissionDefault.Value())).First(s.ctx)
}

// List 返佣方案列表
func (s *promotionCommissionService) List() []promotion.CommissionDetail {
	item := ent.Database.PromotionCommission.
		QueryNotDeleted().
		Where(promotioncommission.TypeNEQ(promotion.CommissionCustom.Value())).
		Order(ent.Asc(promotioncommission.FieldType), ent.Desc(promotioncommission.FieldCreatedAt)).
		AllX(s.ctx)
	res := make([]promotion.CommissionDetail, 0, len(item))
	for _, v := range item {
		res = append(res, s.detail(v))
	}
	return res
}

// Detail 返佣方案详情
func (s *promotionCommissionService) Detail(id uint64) promotion.CommissionDetail {
	return s.detail(ent.Database.PromotionCommission.QueryNotDeleted().Where(promotioncommission.ID(id)).FirstX(s.ctx))
}

func (s *promotionCommissionService) detail(item *ent.PromotionCommission) promotion.CommissionDetail {
	res := promotion.CommissionDetail{
		ID:        item.ID,
		Name:      item.Name,
		Type:      promotion.CommissionType(*item.Type),
		Rule:      *item.Rule,
		Desc:      item.Desc,
		Enable:    item.Enable,
		CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
	}

	if item.StartAt != nil {
		res.StartAt = item.StartAt.Format(carbon.DateTimeLayout)
	}
	if item.EndAt != nil {
		res.EndAt = item.EndAt.Format(carbon.DateTimeLayout)
	}

	amountSum, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", *item.AmountSum), 64)
	res.AmountSum = amountSum
	// 正在使用
	count := ent.Database.PromotionMember.QueryNotDeleted().Where(promotionmember.CommissionID(item.ID)).CountX(s.ctx)
	res.UseCount = uint64(count)
	return res
}

// Update  返佣方案更新
func (s *promotionCommissionService) Update(req *promotion.CommissionCreateReq) {
	commissionInfo, _ := s.PromotionCommissionByID(req.ID)
	if commissionInfo == nil {
		snag.Panic("方案不存在")
	}

	// 校验规则
	s.CheckCommissionRule(&req.Rule)

	startAt := time.Now()

	// 通用方案或者自定义方案 删除重新创建
	ent.Database.PromotionCommission.SoftDeleteOneID(commissionInfo.ID).SetEndAt(startAt).SaveX(s.ctx)

	// 创建新方案
	nc := ent.Database.PromotionCommission.
		Create().
		SetName(req.Name).
		SetRule(&req.Rule).
		SetNillableMemberID(req.MemberID).
		SetNillableDesc(req.Desc).
		SetHistoryID(append(commissionInfo.HistoryID, commissionInfo.ID)).
		SetType(req.Type.Value()).
		SetStartAt(startAt).
		SaveX(s.ctx)

	// 更新会员返佣方案
	ent.Database.PromotionMember.Update().Where(promotionmember.CommissionID(commissionInfo.ID)).SetCommissionID(nc.ID).SaveX(s.ctx)
}

// StatusUpdate 方案状态更新
func (s *promotionCommissionService) StatusUpdate(req *promotion.CommissionEnableReq) {
	commissionInfo, _ := s.PromotionCommissionByID(req.ID)
	if commissionInfo == nil {
		snag.Panic("方案不存在")
	}
	// 全局方案不能禁用
	if *commissionInfo.Type == promotion.CommissionDefault.Value() {
		snag.Panic("全局方案不能修改状态")
	}
	ent.Database.PromotionCommission.UpdateOneID(req.ID).SetEnable(*req.Enable).SaveX(s.ctx)
	// 更新会员返佣方案
	if commissionInfo.Enable && !*req.Enable { // 禁用当前方案
		// 查询默认方案
		defaultCommission, _ := ent.Database.PromotionCommission.QueryNotDeleted().Where(promotioncommission.Type(promotion.CommissionDefault.Value())).First(s.ctx)
		if defaultCommission == nil {
			snag.Panic("默认方案不存在")
		}
		// 禁用当前方案 需要更新使用该方案的会员更换为默认方案

		ent.Database.PromotionMember.Update().Where(promotionmember.CommissionID(commissionInfo.ID)).SetCommissionID(defaultCommission.ID).SaveX(s.ctx)
	}
}

// CheckCommissionRule 返佣规则校验
func (s *promotionCommissionService) CheckCommissionRule(config *promotion.CommissionRule) {
	if len(config.NewUserCommission) == 0 && len(config.RenewalCommission) == 0 {
		snag.Panic("请至少配置一种返佣规则")
	}

	if len(config.NewUserCommission) != 0 {
		// 新签规则校验
		for _, v := range config.NewUserCommission {
			if len(v.Ratio) == 0 {
				snag.Panic(v.Name + " 返佣比例不能为空")
			}
			checkCommissionRatio(v.Ratio[0])
		}
	}

	if len(config.RenewalCommission) != 0 {
		// 续签规则校验
		for _, v := range config.RenewalCommission {
			if len(v.Ratio) == 0 {
				snag.Panic(v.Name + " 续签返佣比例不能为空")
			}
			for _, r := range v.Ratio {
				checkCommissionRatio(r)
			}
		}
	}
}

// CheckCommissionRatio 返佣比例校验
func checkCommissionRatio(ratio float64) {
	if ratio >= 100 {
		snag.Panic("返佣比例不能大于100%")
	}
	if ratio < 0 {
		snag.Panic("返佣比例不能小于0%")
	}
}

// Delete 返佣方案删除
func (s *promotionCommissionService) Delete(id uint64) {
	// 已经启用方案不能删除 // 全局方案不能删除
	info, _ := ent.Database.PromotionCommission.QueryNotDeleted().Where(promotioncommission.ID(id)).
		Where(
			promotioncommission.Or(
				promotioncommission.Enable(true),
				promotioncommission.Type(promotion.CommissionDefault.Value()),
			),
		).First(s.ctx)
	if info != nil {
		snag.Panic("已经启用方案不能删除或全局方案不能删除")
	}
	ent.Database.PromotionCommission.SoftDeleteOneID(id).ExecX(s.ctx)
}

// HistoryList 历史方案
func (s *promotionCommissionService) HistoryList(id uint64) []promotion.CommissionDetail {
	mc, _ := s.PromotionCommissionByID(id)
	if mc == nil {
		snag.Panic("方案不存在")
	}
	item := ent.Database.PromotionCommission.Query().Where(promotioncommission.IDIn(mc.HistoryID...)).Order(ent.Desc(promotioncommission.FieldCreatedAt)).AllX(s.ctx)
	res := make([]promotion.CommissionDetail, 0, len(item))
	for _, v := range item {
		res = append(res, s.detail(v))
	}
	return res
}

// CommissionCalculation 分佣计算
func (s *promotionCommissionService) CommissionCalculation(tx *ent.Tx, req *promotion.CommissionCalculation) (err error) {
	ec := make([]promotion.EarningsCreateReq, 0)

	// 查询会员信息
	member, _ := ent.Database.PromotionMember.QueryNotDeleted().WithReferred().Where(promotionmember.RiderID(req.RiderID)).First(s.ctx)
	if member == nil {
		zap.L().Error("会员不存在", zap.Int64("骑手ID", int64(req.RiderID)))
		return
	}

	referred := member.Edges.Referred

	// 无上级不计算
	if referred.ReferringMemberID == nil {
		zap.L().Error("会员无上级", zap.Int64("会员ID", int64(member.ID)))
		return
	}

	parentMember, _ := NewPromotionMemberService().GetMemberById(*referred.ReferringMemberID)
	if parentMember == nil || (parentMember != nil && parentMember.Edges.Commission == nil) { // 无上级或者上级没有返佣方案
		zap.L().Error("上级会员不存在或者上级会员没有返佣方案", zap.Int64("会员ID", int64(member.ID)))
		return
	}

	// 计算一级收益
	ec = append(ec, s.calculateFirstLevelCommission(req, parentMember))

	// 计算二级收益
	slc := s.calculateSecondLevelCommission(req, parentMember)
	if slc != nil {
		ec = append(ec, *slc)
	}

	// 保存收益
	for _, v := range ec {
		err = s.saveEarningsAndUpdateCommission(tx, v)
		if err != nil {
			zap.L().Error("保存收益失败", zap.Error(err))
		}
	}
	return
}

func (s *promotionCommissionService) calculateFirstLevelCommission(req *promotion.CommissionCalculation, mem *ent.PromotionMember) (res promotion.EarningsCreateReq) {
	res.CommissionID = *mem.CommissionID
	res.MemberID = mem.ID
	res.RiderID = req.RiderID

	conf := mem.Edges.Commission.Rule
	count := 0
	if req.Type == promotion.CommissionTypeNewlySigned { // 新签
		res.Amount = s.calculateCommission(req.CommissionBase, s.getCommissionRatio(conf.NewUserCommission, promotion.FirstLevelNewSubscribeKey, 0))
		res.CommissionRuleKey = promotion.FirstLevelNewSubscribeKey
	} else if req.Type == promotion.CommissionTypeRenewal { // 续签
		if conf.RenewalCommission[promotion.FirstLevelRenewalSubscribeKey].LimitedType == promotion.CommissionLimited { // 有限次数返佣
			// 查询已经返佣的次数
			count, _ = ent.Database.PromotionEarnings.Query().Where(promotionearnings.MemberID(mem.ID), promotionearnings.RiderID(req.RiderID)).Count(s.ctx)
		}

		res.Amount = s.calculateCommission(req.CommissionBase, s.getCommissionRatio(conf.RenewalCommission, promotion.FirstLevelRenewalSubscribeKey, count))
		res.CommissionRuleKey = promotion.FirstLevelRenewalSubscribeKey
	}
	return
}
func (s *promotionCommissionService) calculateSecondLevelCommission(req *promotion.CommissionCalculation, mem *ent.PromotionMember) (res *promotion.EarningsCreateReq) {
	referred := mem.Edges.Referred
	if referred != nil && referred.ReferringMemberID != nil {
		parentMember := referred.Edges.ReferringMember
		if parentMember == nil || parentMember.CommissionID == nil {
			// 无上家或者上家没有返佣方案
			return
		}

		res = &promotion.EarningsCreateReq{}
		res.CommissionID = *parentMember.CommissionID
		res.MemberID = parentMember.ID
		res.RiderID = req.RiderID

		conf := parentMember.Edges.Commission.Rule
		count := 0
		if req.Type == promotion.CommissionTypeNewlySigned { // 新签

			res.Amount = s.calculateCommission(req.CommissionBase, s.getCommissionRatio(parentMember.Edges.Commission.Rule.NewUserCommission, promotion.SecondLevelNewSubscribeKey, 0))
			res.CommissionRuleKey = promotion.SecondLevelNewSubscribeKey

		} else if req.Type == promotion.CommissionTypeRenewal { // 续签

			if conf.RenewalCommission[promotion.SecondLevelRenewalSubscribeKey].LimitedType == promotion.CommissionLimited { // 有限次数返佣
				// 查询已经返佣的次数
				count, _ = ent.Database.PromotionEarnings.Query().Where(promotionearnings.MemberID(parentMember.ID), promotionearnings.RiderID(req.RiderID)).Count(s.ctx)
			}

			res.Amount = s.calculateCommission(req.CommissionBase, s.getCommissionRatio(conf.RenewalCommission, promotion.SecondLevelRenewalSubscribeKey, count))
			res.CommissionRuleKey = promotion.SecondLevelRenewalSubscribeKey
		}
	}
	return
}

func (s *promotionCommissionService) saveEarningsAndUpdateCommission(tx *ent.Tx, req promotion.EarningsCreateReq) (err error) {
	// 保存收益
	err = NewPromotionEarningsService().Create(tx, &req)
	if err != nil {
		zap.L().Error("收益记录创建失败", zap.Error(err), zap.Any("收益记录", req))
		return
	}

	// 更新返佣总收益
	_, err = tx.PromotionCommission.UpdateOneID(req.CommissionID).AddAmountSum(req.Amount).Save(s.ctx)
	if err != nil {
		zap.L().Error("返佣总收益更新失败", zap.Error(err), zap.Any("收益记录", req))
		return
	}

	// 更新会员未结算收益
	_, err = tx.PromotionMember.UpdateOneID(req.MemberID).AddFrozen(req.Amount).Save(s.ctx)
	if err != nil {
		zap.L().Error("会员未结算收益更新失败", zap.Error(err), zap.Any("收益记录", req))
		return
	}

	// 查询任务积分
	lt, _ := NewPromotionLevelTaskService().QueryByKey(req.CommissionRuleKey.Value())
	if lt == nil {
		zap.L().Error("会员任务查询失败 会员成长值记录失败 会员成长值更新失败", zap.Error(err))
		return
	}
	// 记录积分
	err = NewPromotionGrowthService().Create(tx, &promotion.GrowthCreateReq{
		MemberID:    req.MemberID,
		TaksID:      lt.ID,
		GrowthValue: lt.GrowthValue,
		Status:      promotion.GrowthStatusValid.Value(),
	})
	if err != nil {
		zap.L().Error("会员成长值记录失败", zap.Error(err), zap.Any("积分记录", lt))
		return
	}

	// 更新会员成长值并升级
	err = NewPromotionMemberService().UpgradeMemberLevel(tx, req.MemberID, lt.GrowthValue)
	if err != nil {
		zap.L().Error("会员成长值更新失败", zap.Error(err), zap.Any("积分记录", lt))
		return
	}
	return
}

// 计算返佣金额
func (s *promotionCommissionService) calculateCommission(baseAmount, ratio float64) float64 {
	float, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", baseAmount*ratio/100), 64)
	return float
}

// 获取续费返佣比例
func (s *promotionCommissionService) getCommissionRatio(rule map[promotion.CommissionRuleKey]*promotion.CommissionRuleConfig, index promotion.CommissionRuleKey, countIndex int) float64 {
	if rule == nil {
		return 0
	}
	if rule[index].LimitedType == promotion.CommissionLimited { // 有限制次数
		// 判断是否超过最大次数
		if countIndex >= len(rule[index].Ratio) {
			return 0
		}
		return rule[index].Ratio[countIndex]
	}
	return rule[index].Ratio[0]
}

// GetCommissionType 查询骑手是新用户还是续费用户
func (s *promotionCommissionService) GetCommissionType(phone string) (promotion.CommissionCalculationType, error) {

	// 通过实名认证 查询骑手是否是新用户
	ri, _ := ent.Database.Rider.Query().WithPerson().Where(rider.Phone(phone)).First(s.ctx)
	if ri == nil {
		return 0, errors.New("骑手不存在")
	}

	if ri.Edges.Person == nil {
		return 0, errors.New("骑手未实名认证")
	}

	riders, _ := ent.Database.Rider.Query().WithSubscribes().Where(rider.PersonID(ri.Edges.Person.ID)).All(s.ctx)

	for _, v := range riders {
		sub := v.Edges.Subscribes
		if sub != nil && (len(sub) == 1 && sub[0].RenewalDays > 0 || len(sub) > 1) {
			return promotion.CommissionTypeRenewal, nil
		}
	}
	return promotion.CommissionTypeNewlySigned, nil
}

// GetCommissionRule 获取返佣方案
func (s *promotionCommissionService) GetCommissionRule(mem *ent.PromotionMember) promotion.CommissionRuleRes {
	res := promotion.CommissionRuleRes{}

	commission, _ := NewPromotionCommissionService().DefaultPromotionCommission()

	if commission == nil || commission.Rule == nil {
		snag.Panic("默认返佣方案规则不存在")
	}

	dfRule := commission.Rule
	rd := s.appendCommissionRuleDetails(dfRule)
	detailDesc := commission.Desc

	if mem != nil && mem.Edges.Commission != nil {
		detailDesc = mem.Edges.Commission.Desc
		meRule := mem.Edges.Commission.Rule
		if meRule != nil {
			rd = s.appendCommissionRuleDetails(meRule)
		}
	}

	res.Detail = rd
	res.DetailDesc = *detailDesc

	return res
}

// 添加佣金规则详情到结果列表
func (s *promotionCommissionService) appendCommissionRuleDetails(r *promotion.CommissionRule) []promotion.CommissionRuleDetail {
	res := make([]promotion.CommissionRuleDetail, 0, 4)
	cfg := &promotion.CommissionRuleConfig{}

	for k := range promotion.CommissionRuleKeyNames {
		value, ok := r.NewUserCommission[k]
		cfg = value
		if !ok {
			cfg = r.RenewalCommission[k]
		}
		if cfg != nil {
			rd := promotion.CommissionRuleDetail{
				Key:  k,
				Name: cfg.Name,
				Desc: cfg.Desc,
			}

			if cfg.LimitedType == promotion.CommissionUnlimited {
				rd.Ratio = cfg.Ratio[0]
			} else {
				// 取比例最高的值
				rd.Ratio = s.findMaxNumber(cfg.Ratio)
			}

			res = append(res, rd)
		}
	}
	return res
}

func (s *promotionCommissionService) findMaxNumber(numbers []float64) float64 {
	if len(numbers) == 0 {
		// 切片为空时，返回一个默认值或者进行错误处理
		// 这里简单地返回0作为默认值
		return 0
	}

	maxNumber := numbers[0] // 假设第一个数是最大的

	// 遍历切片中的每个元素，更新maxNumber
	for _, num := range numbers {
		if num > maxNumber {
			maxNumber = num
		}
	}

	return maxNumber
}
