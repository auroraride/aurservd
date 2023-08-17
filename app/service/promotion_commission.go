package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotioncommissionplan"
	"github.com/auroraride/aurservd/internal/ent/promotionearnings"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionmembercommission"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type promotionCommissionService struct {
	*BaseService
}

func NewPromotionCommissionService(params ...any) *promotionCommissionService {
	return &promotionCommissionService{
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

	s.CheckCommissionMutex(req)

	// 创建方案
	res := ent.Database.PromotionCommission.Create().
		SetName(req.Name).
		SetType(req.Type.Value()).
		SetRule(&req.Rule).
		SetNillableMemberID(req.MemberID).
		SetNillableDesc(req.Desc).
		SetStartAt(time.Now()).
		SaveX(s.ctx)

	bulk := make([]*ent.PromotionCommissionPlanCreate, len(req.PlanID))
	for i, v := range req.PlanID {
		bulk[i] = ent.Database.PromotionCommissionPlan.Create().SetPlanID(v).SetCommissionID(res.ID)
	}
	ent.Database.PromotionCommissionPlan.CreateBulk(bulk...).ExecX(s.ctx)

	return res
}

// CommissionPlanList 返回自定义方案中列表
func (s *promotionCommissionService) CommissionPlanList(req *model.IDParamReq) (res []*promotion.CommissionPlanListRes) {
	res = make([]*promotion.CommissionPlanListRes, 0)
	all, _ := ent.Database.PromotionCommission.
		QueryNotDeleted().
		WithPlans(
			func(query *ent.PromotionCommissionPlanQuery) {
				query.WithPlan()
			},
		).Where(
		promotioncommission.Or(
			promotioncommission.MemberID(req.ID),
			promotioncommission.TypeNEQ(promotion.CommissionCustom.Value()),
		),
	).All(s.ctx)
	if all == nil {
		return res
	}

	for _, v := range all {
		rows := promotion.CommissionPlanListRes{
			CommissionID:   v.ID,
			CommissionName: v.Name,
			Type:           promotion.CommissionType(*v.Type),
			CreatedAt:      v.CreatedAt.Format(carbon.DateTimeLayout),
			Rule:           v.Rule,
		}
		plans := v.Edges.Plans
		if len(plans) != 0 {
			for _, p := range plans {
				if p.Edges.Plan != nil {
					rows.Plan = append(rows.Plan, &promotion.CommissionPlan{
						ID:     p.Edges.Plan.ID,
						Name:   p.Edges.Plan.Name,
						Amount: p.Edges.Plan.Price,
					})
				}
			}
		}
		res = append(res, &rows)
	}
	return res
}

// CheckCommissionMutex 验证方案互斥
func (s *promotionCommissionService) CheckCommissionMutex(req *promotion.CommissionCreateReq) {
	// 验证方案互斥
	commissionPlan, _ := ent.Database.PromotionCommissionPlan.QueryNotDeleted().
		Where(
			promotioncommissionplan.PlanIDIn(req.PlanID...),
		).Where(
		promotioncommissionplan.HasPromotionCommissionWith(
			promotioncommission.TypeNEQ(promotion.CommissionCustom.Value()),
		),
	).First(s.ctx)

	if commissionPlan != nil {
		snag.Panic("方案已经存在")
	}
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
		WithPlans(
			func(q *ent.PromotionCommissionPlanQuery) {
				q.WithPlan()
			}).
		AllX(s.ctx)
	res := make([]promotion.CommissionDetail, 0, len(item))
	for _, v := range item {
		res = append(res, s.detail(v))
	}
	return res
}

// Detail 返佣方案详情
func (s *promotionCommissionService) Detail(id uint64) promotion.CommissionDetail {
	return s.detail(ent.Database.PromotionCommission.QueryNotDeleted().WithPlans(
		func(q *ent.PromotionCommissionPlanQuery) {
			q.WithPlan()
		}).Where(promotioncommission.ID(id)).FirstX(s.ctx))
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
		AmountSum: item.AmountSum,
	}

	if item.StartAt != nil {
		res.StartAt = item.StartAt.Format(carbon.DateTimeLayout)
	}
	if item.EndAt != nil {
		res.EndAt = item.EndAt.Format(carbon.DateTimeLayout)
	}
	if item.Edges.Plans != nil {
		res.Plan = make([]*promotion.CommissionPlan, 0)
		for _, v := range item.Edges.Plans {
			res.Plan = append(res.Plan, &promotion.CommissionPlan{
				ID:     v.Edges.Plan.ID,
				Name:   v.Edges.Plan.Name,
				Amount: v.Edges.Plan.Price,
			})
		}
	}

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
	commissionPlan, _ := ent.Database.PromotionCommissionPlan.QueryNotDeleted().Where(promotioncommissionplan.HasPromotionCommissionWith(promotioncommission.TypeNEQ(promotion.CommissionCustom.Value()))).Where(
		promotioncommissionplan.PlanIDIn(req.PlanID...),
		promotioncommissionplan.CommissionIDNEQ(req.ID),
	).First(s.ctx)
	if commissionPlan != nil && *commissionInfo.Type != promotion.CommissionCustom.Value() {
		snag.Panic("方案已经存在")
	}

	startAt := time.Now()

	// 通用方案 删除重新创建
	if *commissionInfo.Type != promotion.CommissionCustom.Value() {
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
		req.ID = nc.ID
	}

	ent.Database.PromotionCommissionPlan.SoftDelete().Where(promotioncommissionplan.CommissionID(commissionInfo.ID)).ExecX(s.ctx)
	bulk := make([]*ent.PromotionCommissionPlanCreate, len(req.PlanID))
	for i, v := range req.PlanID {
		bulk[i] = ent.Database.PromotionCommissionPlan.Create().SetPlanID(v).SetCommissionID(req.ID)
	}
	ent.Database.PromotionCommissionPlan.CreateBulk(bulk...).ExecX(s.ctx)

	// 更新会员返佣方案
	if *commissionInfo.Type == promotion.CommissionCustom.Value() {
		ent.Database.PromotionMemberCommission.Update().Where(promotionmembercommission.CommissionID(commissionInfo.ID)).SetCommissionID(req.ID).SaveX(s.ctx)
	}
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
}

// CheckCommissionRule 返佣规则校验
func (s *promotionCommissionService) CheckCommissionRule(config *promotion.CommissionRule) {
	if len(config.NewUserCommission) == 0 && len(config.RenewalCommission) == 0 {
		snag.Panic("请至少配置一种返佣规则")
	}

	if len(config.NewUserCommission) != 0 {
		// 新签规则校验
		for _, v := range config.NewUserCommission {
			if len(v.Value) == 0 {
				snag.Panic(v.Name + " 返佣设置不能为空")
			}
			checkCommissionRatio(v.Value[0])
		}
	}

	if len(config.RenewalCommission) != 0 {
		// 续签规则校验
		for _, v := range config.RenewalCommission {
			if len(v.Value) == 0 {
				snag.Panic(v.Name + " 续签返佣比例不能为空")
			}
			for _, r := range v.Value {
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
	if info != nil && *info.Type != promotion.CommissionCustom.Value() {
		snag.Panic("已经启用方案不能删除或全局方案不能删除")
	}
	ent.Database.PromotionCommission.SoftDeleteOneID(id).ExecX(s.ctx)
	ent.Database.PromotionCommissionPlan.SoftDelete().Where(promotioncommissionplan.CommissionID(id)).ExecX(s.ctx)
	ent.Database.PromotionMemberCommission.SoftDelete().Where(promotionmembercommission.CommissionID(id)).ExecX(s.ctx)
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
		zap.L().Error("分佣计算 会员不存在", zap.Int64("骑手ID", int64(req.RiderID)))
		return
	}

	referred := member.Edges.Referred

	// 无上级不计算
	if referred.ReferringMemberID == nil {
		zap.L().Error("分佣计算 会员无上级", zap.Int64("会员ID", int64(member.ID)))
		return
	}

	parentMember, _ := NewPromotionMemberService().GetMemberById(*referred.ReferringMemberID)
	if parentMember == nil { // 无上级
		zap.L().Error("分佣计算 上级会员不存", zap.Int64("会员ID", int64(member.ID)))
		return
	}

	// 计算一级收益
	ec = append(ec, s.calculateFirstLevelCommission(parentMember, req))

	// 计算二级收益
	slc := s.calculateSecondLevelCommission(parentMember, req)
	if slc != nil {
		ec = append(ec, *slc)
	}

	// 新增续费次数统计
	q := ent.Database.PromotionMember.Update().Where(promotionmember.RiderID(req.RiderID))
	if req.Type == promotion.CommissionTypeRenewal {
		q.AddRenewCount(1)
	} else {
		q.AddNewSignCount(1)
	}
	err = q.Exec(s.ctx)
	if err != nil {
		zap.L().Error("更新会员续费次数失败", zap.Error(err), zap.Int64("会员ID", int64(member.ID)), zap.Int("类型", int(req.Type)))
	}

	// 保存收益
	for _, v := range ec {
		err = s.saveEarningsAndUpdateCommission(tx, v)
		if err != nil {
			zap.L().Error("保存收益失败", zap.Error(err))
		}
		zap.L().Info("返佣成功", zap.Any(fmt.Sprintf("收益记录 member_id: %d", v.MemberID), v), zap.Any(fmt.Sprintf("订单信息 order_id: %d", req.OrderID), req))
	}
	return
}

// 判断使用哪种返佣方案 优先查询自定义返佣方案 再查询通用返佣方案
func (s *promotionCommissionService) getCommissionRule(mem *ent.PromotionMember, planID uint64) *ent.PromotionCommission {

	commission, _ := ent.Database.PromotionCommission.QueryNotDeleted().
		Where(
			promotioncommission.Enable(true),
			promotioncommission.TypeEQ(promotion.CommissionCustom.Value()),
		).Where(
		promotioncommission.HasPlansWith(
			promotioncommissionplan.PlanID(planID),
			promotioncommissionplan.MemberID(mem.ID),
		)).First(s.ctx)
	if commission != nil {
		return commission
	}
	commission, _ = ent.Database.PromotionCommission.QueryNotDeleted().
		Where(
			promotioncommission.Enable(true),
			promotioncommission.TypeNEQ(promotion.CommissionCustom.Value()),
		).Where(promotioncommission.HasPlansWith(promotioncommissionplan.PlanID(planID))).First(s.ctx)
	if commission != nil {
		return commission

	}
	return nil
}

func (s *promotionCommissionService) calculateFirstLevelCommission(mem *ent.PromotionMember, req *promotion.CommissionCalculation) (res promotion.EarningsCreateReq) {

	// 判断使用哪种返佣方案 优先查询个人返佣方案 再查询通用返佣方案
	commissions := s.getCommissionRule(mem, req.PlanID)
	if commissions == nil {
		zap.L().Error("会员返佣方案不存在", zap.Any("会员", mem), zap.Any("订单", req))
		return
	}

	res.CommissionID = commissions.ID
	res.MemberID = mem.ID
	res.RiderID = req.RiderID
	res.OrderID = req.OrderID

	conf := commissions.Rule
	count := 0
	if req.Type == promotion.CommissionTypeNewlySigned { // 新签

		value, optionType := s.getCommissionRatio(conf.NewUserCommission, promotion.FirstLevelNewSubscribeKey, 0)
		res.Amount = s.calculateCommission(req.ActualAmount, value, optionType)

		res.CommissionRuleKey = promotion.FirstLevelNewSubscribeKey

	} else if req.Type == promotion.CommissionTypeRenewal { // 续签

		if conf.RenewalCommission[promotion.FirstLevelRenewalSubscribeKey].LimitedType == promotion.CommissionLimited { // 有限次数返佣
			// 查询已经返佣的次数
			count, _ = NewPromotionEarningsService().CountCommission(mem.ID, req.RiderID)
		}

		value, optionType := s.getCommissionRatio(conf.RenewalCommission, promotion.FirstLevelRenewalSubscribeKey, count)
		res.Amount = s.calculateCommission(req.ActualAmount, value, optionType)

		res.CommissionRuleKey = promotion.FirstLevelRenewalSubscribeKey
	}
	return
}
func (s *promotionCommissionService) calculateSecondLevelCommission(mem *ent.PromotionMember, req *promotion.CommissionCalculation) (res *promotion.EarningsCreateReq) {
	referred := mem.Edges.Referred
	if referred != nil && referred.ReferringMemberID != nil {
		parentMember := referred.Edges.ReferringMember
		if parentMember == nil {
			zap.L().Error("上级会员不存在", zap.Any("会员", mem), zap.Any("订单", req))
			return
		}

		// 判断使用哪种返佣方案 优先查询个人返佣方案 再查询通用返佣方案
		commissions := s.getCommissionRule(mem, req.PlanID)
		if commissions == nil {
			return
		}

		res = &promotion.EarningsCreateReq{}
		res.CommissionID = commissions.ID
		res.MemberID = parentMember.ID
		res.RiderID = req.RiderID
		res.OrderID = req.OrderID

		conf := commissions.Rule
		count := 0
		if req.Type == promotion.CommissionTypeNewlySigned { // 新签

			value, optionType := s.getCommissionRatio(commissions.Rule.NewUserCommission, promotion.SecondLevelNewSubscribeKey, 0)
			res.Amount = s.calculateCommission(req.ActualAmount, value, optionType)

			res.CommissionRuleKey = promotion.SecondLevelNewSubscribeKey

		} else if req.Type == promotion.CommissionTypeRenewal { // 续签

			if conf.RenewalCommission[promotion.SecondLevelRenewalSubscribeKey].LimitedType == promotion.CommissionLimited { // 有限次数返佣
				// 查询已经返佣的次数
				count, _ = NewPromotionEarningsService().CountCommission(mem.ID, req.RiderID)
			}

			value, optionType := s.getCommissionRatio(conf.RenewalCommission, promotion.SecondLevelRenewalSubscribeKey, count)
			res.Amount = s.calculateCommission(req.ActualAmount, value, optionType)

			res.CommissionRuleKey = promotion.SecondLevelRenewalSubscribeKey
		}
	}
	return
}

func (s *promotionCommissionService) saveEarningsAndUpdateCommission(tx *ent.Tx, req promotion.EarningsCreateReq) (err error) {

	if req.Amount != 0 {
		// 查询是是否重复返佣
		es, _ := ent.Database.PromotionEarnings.Query().Where(promotionearnings.OrderID(req.OrderID), promotionearnings.CommissionRuleKey(req.CommissionRuleKey.Value())).First(s.ctx)
		if es != nil {
			zap.L().Error("返佣失败,重复返佣", zap.Any("收益记录", req))
			return
		}

		// 保存收益
		err = NewPromotionEarningsService().Create(tx, &req)
		if err != nil {
			zap.L().Error("收益记录创建失败", zap.Error(err), zap.Any("收益记录", req))
			return
		}

		// 更新返佣总收益
		com, _ := ent.Database.PromotionCommission.Query().Where(promotioncommission.ID(req.CommissionID)).First(s.ctx)
		if com == nil {
			zap.L().Error("返佣方案不存在", zap.Error(err), zap.Any("收益记录", req))
			return
		}
		_, err = tx.PromotionCommission.UpdateOneID(req.CommissionID).SetAmountSum(tools.NewDecimal().Sum(com.AmountSum, req.Amount)).Save(s.ctx)
		if err != nil {
			zap.L().Error("返佣总收益更新失败", zap.Error(err), zap.Any("收益记录", req))
			return
		}

		// 更新会员未结算收益
		mem, _ := ent.Database.PromotionMember.QueryNotDeleted().Where(promotionmember.ID(req.MemberID)).First(s.ctx)
		if mem == nil {
			zap.L().Error("会员不存在", zap.Error(err), zap.Any("收益记录", req))
			return
		}
		_, err = tx.PromotionMember.UpdateOneID(req.MemberID).SetFrozen(tools.NewDecimal().Sum(mem.Frozen, req.Amount)).Save(s.ctx)
		if err != nil {
			zap.L().Error("会员未结算收益更新失败", zap.Error(err), zap.Any("收益记录", req), zap.Any("收益记录", req))
			return
		}
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
		RiderID:     req.RiderID,
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
func (s *promotionCommissionService) calculateCommission(actualAmount, value float64, optionType promotion.CommissionOptionType) float64 {
	if optionType == promotion.Percentage {
		dl := tools.NewDecimal()
		value = dl.Mul(actualAmount, value/100)
	}
	zap.L().Info("计算返佣金额", zap.Any("计算参数", map[string]interface{}{
		"actualAmount": actualAmount,
		"value":        value,
		"optionType":   optionType,
	}))
	return value
}

// 获取续费返佣比例
func (s *promotionCommissionService) getCommissionRatio(rule map[promotion.CommissionRuleKey]*promotion.CommissionRuleConfig, index promotion.CommissionRuleKey, countIndex int) (float64, promotion.CommissionOptionType) {
	if rule == nil {
		return 0, 0
	}
	if rule[index].LimitedType == promotion.CommissionLimited { // 有限制次数
		// 判断是否超过最大次数 或者countIndex小于0
		if countIndex >= len(rule[index].Value) || countIndex < 0 {
			return 0, 0
		}
		return rule[index].Value[countIndex], rule[index].OptionType
	}
	return rule[index].Value[0], rule[index].OptionType
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

	if len(riders) > 1 {
		return promotion.CommissionTypeRenewal, nil
	}

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

	// // 查询最高返佣金额
	// var commissionBase float64
	// pl, _ := ent.Database.Plan.QueryNotDeleted().Where(plan.CommissionBaseNEQ(0)).Order(ent.Desc(plan.FieldCommissionBase)).First(s.ctx)
	// if pl == nil {
	// 	commissionBase = 0
	// } else {
	// 	commissionBase = pl.CommissionBase
	// }

	dfRule := commission.Rule
	rd := s.appendCommissionRuleDetails(dfRule)
	detailDesc := commission.Desc

	// if mem != nil && mem.Edges.Commission != nil {
	// 	detailDesc = mem.Edges.Commission.Desc
	// 	meRule := mem.Edges.Commission.Rule
	// 	if meRule != nil {
	// 		rd = s.appendCommissionRuleDetails(meRule)
	// 	}
	// }

	res.Detail = rd
	res.DetailDesc = detailDesc

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
				rd.Ratio = cfg.Value[0]
			} else {
				// 取比例最高的值
				rd.Ratio = s.findMaxNumber(cfg.Value)
			}

			// 计算返佣金额
			// rd.Amount = uint64(math.Round(tools.NewDecimal().Mul(amount, rd.Ratio/100)))

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

// RiderActivateCommission 骑手激活分佣
func (s *promotionCommissionService) RiderActivateCommission(sub *ent.Subscribe) {
	// 查询订单是否支付,如果未支付不返佣
	do, _ := ent.Database.Order.Query().Where(order.SubscribeID(sub.ID), order.Status(model.OrderStatusPaid)).First(s.ctx)
	if do == nil {
		zap.L().Error(fmt.Sprintf("激活成功获取返佣失败,订单不存在或未支付 subid:%d", sub.ID), zap.Any("order", do))
		return
	}

	// TODO edges查询写成公共方法

	// 查询骑手
	r := sub.Edges.Rider
	if r == nil {
		r, _ = sub.QueryRider().First(s.ctx)
	}

	if r == nil {
		zap.L().Error(fmt.Sprintf("激活成功获取返佣失败,骑手不存在 subid:%d", sub.ID))
		return
	}

	// 查询订阅
	p := sub.Edges.Plan
	if p == nil {
		p, _ = sub.QueryPlan().First(s.ctx)
	}

	// 判断返佣类型 新签有可能是续签
	commissionType, err := NewPromotionCommissionService().GetCommissionType(r.Phone)
	if err != nil || p == nil {
		zap.L().Error("激活成功获取返佣失败", zap.Error(err))
		return
	}

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		err = NewPromotionCommissionService().CommissionCalculation(tx, &promotion.CommissionCalculation{
			RiderID:        r.ID,
			CommissionBase: p.CommissionBase,
			Type:           commissionType,
			OrderID:        do.ID,
			ActualAmount:   do.Total,
			Price:          p.Price,
		})
		return
	})
}
