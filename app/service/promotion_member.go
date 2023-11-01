package service

import (
	stdsql "database/sql"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-module/carbon/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/person"
	"github.com/auroraride/aurservd/internal/ent/promotioncommissionplan"
	"github.com/auroraride/aurservd/internal/ent/promotionlevel"
	"github.com/auroraride/aurservd/internal/ent/promotionmember"
	"github.com/auroraride/aurservd/internal/ent/promotionperson"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/auroraride/aurservd/pkg/utils"
)

type promotionMemberService struct {
	*BaseService
	tokenCacheKey string
}

func NewPromotionMemberService(params ...any) *promotionMemberService {
	return &promotionMemberService{
		tokenCacheKey: ar.Config.Environment.UpperString() + ":" + "AGENT:TOKEN",
		BaseService:   newService(params...),
	}
}

func (s *promotionMemberService) TokenVerify(token string) (me *ent.PromotionMember) {
	// 获取token对应ID
	id, _ := ar.Redis.HGet(s.ctx, s.tokenCacheKey, token).Uint64()
	if id <= 0 {
		return
	}

	// 反向校验token是否正确
	if ar.Redis.HGet(s.ctx, s.tokenCacheKey, strconv.FormatUint(id, 10)).Val() != token {
		return
	}

	me, _ = NewPromotionMemberService().GetMemberById(id)
	if me == nil {
		return
	}
	return me
}

// Signin 推广会员登录
func (s *promotionMemberService) Signin(req *promotion.MemberSigninReq) *promotion.MemberSigninRes {
	switch req.SigninType {
	case promotion.MemberSigninTypeSms:
		// 校验短信
		NewSms().VerifyCodeX(req.Phone, req.SmsID, req.Code)
	case promotion.MemberSigninTypeWechat:
		// 获取手机号
		req.Phone = NewPromotionMiniProgram().GetPhoneNumber(req.Code)
	default:
		snag.Panic("不支持的登录方式")
	}

	// 获取骑手信息或创建骑手
	ri, _ := ent.Database.Rider.QueryNotDeleted().Where(rider.Phone(req.Phone)).First(s.ctx)
	if ri == nil {
		ri = s.createRider(req.Phone)
	}
	// 查询是否有最新的订阅
	var subscribeID *uint64
	sub := NewSubscribe().Recent(ri.ID)
	if sub != nil {
		subscribeID = &sub.ID
	}

	mem, _ := s.GetMemberByPhone(req.Phone)
	if mem == nil {
		// 创建会员
		mem = s.createMember(&promotion.MemberCreateReq{
			Phone:       req.Phone,
			Name:        req.Name,
			SubscribeID: subscribeID,
			RiderID:     &ri.ID,
		})
	} else {
		s.updateMemberInfo(&promotion.MemberCreateReq{
			Phone:       req.Phone,
			SubscribeID: subscribeID,
			RiderID:     &ri.ID,
		})
	}

	if !mem.Enable {
		snag.Panic("账号已被禁用")
	}

	return s.signin(mem)
}

// Signup 邀请注册
func (s *promotionMemberService) Signup(req *promotion.MemberSigninReq) promotion.MemberInviteRes {
	switch req.SigninType {
	case promotion.MemberSigninTypeSms:
		// 校验短信
		NewSms().VerifyCodeX(req.Phone, req.SmsID, req.Code)
	default:
		snag.Panic("不支持的注册方式")
	}

	res := promotion.MemberInviteRes{}

	ri, _ := ent.Database.Rider.QueryNotDeleted().Where(rider.Phone(req.Phone)).First(s.ctx)
	// 用于判断是新注册还是已经注册过
	isSignup := true
	if ri == nil {
		isSignup = false
		// 获取骑手信息或创建骑手
		ri = s.createRider(req.Phone)
	}
	// 查询是否有最新的订阅
	var subscribeID *uint64
	sub := NewSubscribe().Recent(ri.ID)
	if sub != nil {
		subscribeID = &sub.ID
	}
	//  获取会员信息
	mem, _ := s.GetMemberByPhone(req.Phone)
	if mem == nil {
		// 创建会员
		mem = s.createMember(&promotion.MemberCreateReq{
			Phone:       req.Phone,
			Name:        req.Name,
			SubscribeID: subscribeID,
			RiderID:     &ri.ID,
		})
	} else {
		s.updateMemberInfo(&promotion.MemberCreateReq{
			Phone:       req.Phone,
			SubscribeID: subscribeID,
			RiderID:     &ri.ID,
		})
	}

	referralsProgress := &promotion.Referrals{
		ReferringMemberId: req.ReferringMemberID,
		ReferredMemberId:  &mem.ID,
		Name:              req.Name,
		RiderID:           &ri.ID,
		Status:            promotion.ReferralsStatusSuccess,
		Remark:            "邀请成功",
	}
	// 判断当前账号
	var it promotion.InviteType
	if mem.ID == *req.ReferringMemberID {
		// 判断是否绑定自己
		it = promotion.MemberInviteSelfFail
	}

	if mem.Edges.Referred != nil && mem.Edges.Referred.ReferringMemberID != nil {
		// 判断是否已经被邀请
		it = promotion.MemberInviteFail
	}

	if s.IsActivation([]uint64{ri.ID}) {
		it = promotion.MemberActivationFail
	}
	if it != 0 {
		res.InviteType = it
		referralsProgress.Status = promotion.ReferralsStatusFail
		referralsProgress.Remark = it.String()
		// 邀请失败
		NewPromotionReferralsService().MemberReferralsProgress(referralsProgress)
		return res
	}

	// 判断骑手是否能被绑定
	if ircb := s.IsRiderCanBind(ri); ircb != promotion.MemberAllowBind {
		res.InviteType = ircb
		referralsProgress.Status = promotion.ReferralsStatusFail
		referralsProgress.Remark = ircb.String()
		// 邀请失败
		NewPromotionReferralsService().MemberReferralsProgress(referralsProgress)
		return res
	}

	// 判断骑手是否实名
	if s.IsRiderPerson(req.Phone) {
		// 已经实名的 直接记录推荐关系 并且将推荐关系状态改为已激活
		NewPromotionReferralsService().CreateReferrals(&promotion.Referrals{
			ReferringMemberId: req.ReferringMemberID,
			ReferredMemberId:  &mem.ID,
			RiderID:           &ri.ID,
		})
		referralsProgress.Status = promotion.ReferralsStatusSuccess
		referralsProgress.Remark = "邀请成功"
	} else {
		// 未实名不绑定关系
		referralsProgress.Status = promotion.ReferralsStatusInviting
		referralsProgress.Remark = "待用户实名后判定邀请是否生效"
	}
	NewPromotionReferralsService().MemberReferralsProgress(referralsProgress)

	if isSignup {
		// 没实名提示
		if !s.IsRiderPerson(req.Phone) {
			res.InviteType = promotion.MemberSignSuccessWaitAuth
			return res
		}
	} else {
		res.InviteType = promotion.MemberSignSuccess
		return res
	}
	res.InviteType = promotion.MemberBindSuccess
	return res
}

// IsRiderPerson 判断骑手是否实名
func (s *promotionMemberService) IsRiderPerson(phone string) bool {
	if exists, _ := ent.Database.Rider.QueryNotDeleted().
		Where(rider.Phone(phone)).
		QueryPerson().
		Where(
			person.Status(model.PersonAuthenticated.Value()),
		).Exist(s.ctx); !exists {
		return false
	}
	return true
}

// IsRiderCanBind 判断被邀请的骑手是否能被绑定
func (s *promotionMemberService) IsRiderCanBind(riderInfo *ent.Rider) (it promotion.InviteType) {
	it = promotion.MemberAllowBind

	q := ent.Database.Rider.Query()
	if riderInfo.PersonID != nil {
		// 查询实名用户 和 电话
		q.Where(
			rider.Or(
				rider.Phone(riderInfo.Phone),
				rider.PersonID(*riderInfo.PersonID),
			),
		)
	} else {
		q.Where(rider.Phone(riderInfo.Phone))
	}

	// 判断这个人名下所有账号
	riAll, _ := q.IDs(s.ctx)
	if len(riAll) > 0 {
		// 查询所有 注册会员
		pm := ent.Database.PromotionMember.Query().Where(promotionmember.RiderIDIn(riAll...)).WithReferred().AllX(s.ctx)
		for _, v := range pm {
			if v.Edges.Referred != nil && v.Edges.Referred.ReferringMemberID != nil {
				// 判断是否已经被邀请
				it = promotion.MemberInviteOtherFail
				return
			}
		}

		// 判断是否已经激活
		if s.IsActivation(riAll) {
			it = promotion.MemberActivationOtherFail
		}
	}

	return it
}

// IsActivation 判断是否已经激活
func (s *promotionMemberService) IsActivation(riderId []uint64) bool {
	// 判断是否已经被激活
	sub := ent.Database.Subscribe.Query().Where(subscribe.RiderIDIn(riderId...), subscribe.StatusNEQ(model.SubscribeStatusCanceled)).Order(ent.Desc(subscribe.FieldCreatedAt)).AllX(s.ctx)

	for _, v := range sub {

		if v.EndAt != nil { // 判断最新的退订时间是否超出设置天数
			past := int(carbon.CreateFromStdTime(*v.EndAt).AddDay().DiffInDays(carbon.Now()))
			// 判定退订时间是否超出设置天数
			if !model.NewRecentSubscribePastDays(past).Commission() {
				return true
			}

		} else if v.Status != model.SubscribeStatusInactive { // 判断是否已经激活
			return true
		}

	}
	return false
}

// 获取骑手信息或创建骑手
func (s *promotionMemberService) createRider(phone string) *ent.Rider {
	// 创建骑手
	rinfo := ent.Database.Rider.Create().SetPhone(phone).SaveX(s.ctx)
	return rinfo
}

// 更新会员信息
func (s *promotionMemberService) updateMemberInfo(req *promotion.MemberCreateReq) {

	ent.Database.PromotionMember.Update().
		Where(promotionmember.Phone(req.Phone)).
		SetNillableRiderID(req.RiderID).
		SetNillableSubscribeID(req.SubscribeID).
		ExecX(s.ctx)
}

// GetMemberByPhone  获取会员信息
func (s *promotionMemberService) GetMemberByPhone(phone string) (*ent.PromotionMember, error) {
	return ent.Database.PromotionMember.QueryNotDeleted().
		Where(promotionmember.Phone(phone)).
		WithReferred().
		WithLevel().
		WithPerson().
		WithRider().
		First(s.ctx)
}

// 创建会员
func (s *promotionMemberService) createMember(req *promotion.MemberCreateReq) *ent.PromotionMember {
	// 创建会员
	mem := NewPromotionMemberService().Create(&promotion.MemberCreateReq{
		Phone:   req.Phone,
		Name:    req.Name,
		RiderID: req.RiderID,
	})
	return mem
}

func (s *promotionMemberService) signin(mem *ent.PromotionMember) *promotion.MemberSigninRes {
	idstr := strconv.FormatUint(mem.ID, 10)
	// 查询并删除旧token key
	exists := ar.Redis.HGet(s.ctx, s.tokenCacheKey, idstr).Val()
	if exists != "" {
		ar.Redis.HDel(s.ctx, s.tokenCacheKey, exists)
	}

	// 生成token
	token := utils.NewEcdsaToken()

	// 存储登录token和ID进行对应
	ar.Redis.HSet(s.ctx, s.tokenCacheKey, token, mem.ID)
	ar.Redis.HSet(s.ctx, s.tokenCacheKey, idstr, token)

	return &promotion.MemberSigninRes{
		Profile: s.MemberProfile(mem.ID),
		Token:   token,
	}
}

// MemberProfile  会员信息
func (s *promotionMemberService) MemberProfile(id uint64) *promotion.MemberProfile {
	// 会员信息
	mem, _ := NewPromotionMemberService().GetMemberById(id)
	if mem == nil {
		snag.Panic("会员不存在")
	}

	res := &promotion.MemberProfile{
		MemberBaseInfo: promotion.MemberBaseInfo{
			ID:    mem.ID,
			Phone: s.MaskSensitiveInfo(mem.Phone, 3, 4),
			Name:  mem.Name,
		},
		AvatarURL: mem.AvatarURL,
	}
	if mem.Edges.Level != nil {
		res.Level = mem.Edges.Level.Level
	}
	if mem.Edges.Person != nil {
		res.AuthStatusName = promotion.PersonAuthStatus(mem.Edges.Person.Status).String()
		res.AuthStatus = promotion.PersonAuthStatus(mem.Edges.Person.Status)
		res.IDCardNumber = s.MaskSensitiveInfo(mem.Edges.Person.IDCardNumber, 2, 2)
	}
	return res
}

// MaskSensitiveInfo 脱敏敏感信息，保留前prefixLen位和后suffixLen位，中间用"*"脱敏
func (s *promotionMemberService) MaskSensitiveInfo(sensitiveInfo string, prefixLen, suffixLen int) string {
	sensitiveInfo = strings.TrimSpace(sensitiveInfo)
	infoLen := utf8.RuneCountInString(sensitiveInfo)

	if infoLen <= prefixLen+suffixLen {
		return sensitiveInfo
	}

	firstPart := sensitiveInfo[:prefixLen]
	lastPart := sensitiveInfo[infoLen-suffixLen:]
	middlePart := strings.Repeat("*", infoLen-prefixLen-suffixLen)
	return firstPart + middlePart + lastPart
}

// MaskName 脱敏姓名，保留姓的第一个字和名的最后一个字，中间用"*"脱敏
func (s *promotionMemberService) MaskName(name string) string {
	name = strings.TrimSpace(name)
	runes := []rune(name)
	nameLen := len(runes)

	if nameLen <= 1 {
		return name
	}

	firstPart := string(runes[0]) // 留下姓的第一个字

	if nameLen == 2 {
		secondPart := string(runes[1]) // 脱敏名的第一个字
		return firstPart + "*" + secondPart
	}

	lastPart := string(runes[nameLen-1])         // 留下最后一个名字的第一个字
	middlePart := strings.Repeat("*", nameLen-2) // 脱敏中间名字的所有字
	return firstPart + middlePart + lastPart
}

// List 会员列表
func (s *promotionMemberService) List(req *promotion.MemberReq) *model.PaginationRes {
	q := ent.Database.PromotionMember.QueryNotDeleted().WithPerson().WithLevel().Order(ent.Desc(promotionmember.FieldCreatedAt))
	if req.Keyword != nil {
		q.Where(
			promotionmember.Or(
				promotionmember.PhoneContains(*req.Keyword),
				promotionmember.NameContains(*req.Keyword),
			),
		)
	}
	if req.Enable != nil {
		q.Where(promotionmember.Enable(*req.Enable))
	}

	if req.LevelID != nil {
		q.Where(promotionmember.HasLevelWith(promotionlevel.IDEQ(*req.LevelID)))
	}

	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(promotionmember.CreatedAtGTE(start), promotionmember.CreatedAtLTE(end))
	}

	if req.AuthStatus != nil {
		if *req.AuthStatus == promotion.PersonUnauthenticated {
			q.Where(
				promotionmember.Or(
					promotionmember.HasPersonWith(promotionperson.StatusEQ(req.AuthStatus.Value())),
					promotionmember.PersonIDIsNil(),
				),
			)
		} else {
			q.Where(promotionmember.HasPersonWith(promotionperson.StatusEQ(req.AuthStatus.Value())))
		}
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.PromotionMember) (res promotion.MemberRes) {
			return s.detail(item)
		},
	)
}

// Detail 会员详情
func (s *promotionMemberService) Detail(req *promotion.MemberReq) promotion.MemberRes {
	info, _ := ent.Database.PromotionMember.QueryNotDeleted().WithReferred(
		func(query *ent.PromotionReferralsQuery) {
			query.WithReferringMember()
		}).WithLevel().WithPerson().Where(promotionmember.IDEQ(req.ID)).First(s.ctx)
	if info == nil {
		snag.Panic("会员不存在")
	}
	return s.detail(info)
}

func (s *promotionMemberService) detail(item *ent.PromotionMember) (res promotion.MemberRes) {
	res = promotion.MemberRes{
		ID: item.ID,
		MemberInfo: promotion.MemberBaseInfo{
			ID:    item.ID,
			Phone: item.Phone,
			Name:  item.Name,
		},
		Enable: item.Enable,

		CreatedAt:          item.CreatedAt.Format(carbon.DateTimeLayout),
		CurrentGrowthValue: item.CurrentGrowthValue,
	}
	level := item.Edges.Level
	if level != nil {
		res.Level = level.Level
		res.PrivilegeCommission = level.CommissionRatio
	}

	res.TotalBalance = tools.NewDecimal().Sum(item.Balance, item.Frozen)

	if item.Edges.Referred != nil && item.Edges.Referred.Edges.ReferringMember != nil {
		res.ParentInfo = &promotion.MemberBaseInfo{
			ID:    item.Edges.Referred.Edges.ReferringMember.ID,
			Name:  item.Edges.Referred.Edges.ReferringMember.Name,
			Phone: item.Edges.Referred.Edges.ReferringMember.Phone,
		}
	}

	res.AuthStatus = promotion.PersonUnauthenticated
	res.AuthStatusName = promotion.PersonUnauthenticated.String()
	if item.Edges.Person != nil {
		res.AuthStatusName = promotion.PersonAuthStatus(item.Edges.Person.Status).String()
		res.AuthStatus = promotion.PersonAuthStatus(item.Edges.Person.Status)
	}

	if item.Edges.Cards != nil {
		for _, card := range item.Edges.Cards {
			res.BankCard = append(res.BankCard, &promotion.BankCardRes{
				CardNo:      NewPromotionBankCardService().EncryptCardNo(card.CardNo),
				IsDefault:   card.IsDefault,
				Bank:        card.Bank,
				BankLogoURL: card.BankLogoURL,
			})
		}
	}

	team := s.TeamStatistics(&promotion.MemberTeamReq{
		ID: item.ID,
	})
	// 团队人数
	res.FirstLevel = team.FirstLevelCount
	res.SecondLevel = team.SecondLevelCount

	return res
}

// Create 创建会员
func (s *promotionMemberService) Create(req *promotion.MemberCreateReq) *ent.PromotionMember {
	mem := ent.Database.PromotionMember.Create().
		SetName(req.Name).
		SetPhone(req.Phone).
		SetNillableRiderID(req.RiderID).
		SaveX(s.ctx)

	return mem
}

// GetMemberById 通过id获取会员信息
func (s *promotionMemberService) GetMemberById(id uint64) (*ent.PromotionMember, error) {
	return ent.Database.PromotionMember.QueryNotDeleted().
		WithReferred(
			func(query *ent.PromotionReferralsQuery) {
				query.WithReferringMember()
			},
		).
		WithLevel().
		WithCommissions(func(query *ent.PromotionMemberCommissionQuery) {
			query.WithCommission()
		}).
		WithPerson().
		WithCards().
		Where(promotionmember.IDEQ(id)).
		First(s.ctx)
}

// Update 编辑会员
func (s *promotionMemberService) Update(req *promotion.MemberUpdateReq) {
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		err = tx.PromotionMember.
			UpdateOneID(req.ID).
			SetNillableEnable(req.Enable).
			SetNillableName(req.Name).
			Exec(s.ctx)
		if err != nil {
			snag.Panic("更新失败")
		}

		if req.Enable != nil && !*req.Enable {
			// 禁用会员 删除token禁止登录
			ar.Redis.HDel(s.ctx, s.tokenCacheKey, strconv.FormatUint(req.ID, 10))
		}
		return
	})
}

// TeamList 会员团队列表
func (s *promotionMemberService) TeamList(ctx echo.Context, req *promotion.MemberTeamReq) model.PaginationRes {
	sqls := `
			WITH RECURSIVE member_hierarchy AS (
                SELECT referred_member_id, 1 AS level, re.rider_id, re.created_at, r.person_id
			    FROM promotion_referrals re
			    LEFT JOIN rider r ON r.id = re.rider_id
				` + fmt.Sprintf(" WHERE referring_member_id = %d", req.ID) + `
                UNION ALL
                SELECT mr.referred_member_id, mh.level + 1 AS level, mr.rider_id, mr.created_at, r.person_id
                FROM member_hierarchy mh
                INNER JOIN promotion_referrals as mr ON mh.referred_member_id = mr.referring_member_id
				LEFT JOIN rider r ON r.id = mh.rider_id
				WHERE mh.level < 2
			)
			SELECT referred_member_id, m.phone, COALESCE(r.name, m.name) AS name, mh.level, s.status as subscribeStatus, s.start_at as subscribeStartAt, renew_count AS renewalCount
				FROM member_hierarchy mh
				JOIN promotion_member m ON m.id = mh.referred_member_id
				LEFT JOIN rider r ON m.rider_id = r.id
				LEFT JOIN subscribe s ON r.id = s.rider_id
			    AND s.created_at = (
			        SELECT MAX(created_at)
			        FROM subscribe
			        WHERE rider_id = r.id
			    )
			WHERE mh.referred_member_id IN (
			    SELECT referred_member_id
			    FROM member_hierarchy
			)`

	// 条件筛选
	s.MemberTeamFilter(req, &sqls)

	// 排序 asc
	sqls += " ORDER BY mh.created_at DESC "

	if req.Current == 0 && req.PageSize == 0 {
		// 默认分页
		req.Current = 1
		req.PageSize = 20
	}
	sqls += fmt.Sprintf(" LIMIT %d OFFSET %d ", req.PageSize, (req.Current-1)*req.PageSize)

	// 参数
	rows, err := ent.Database.QueryContext(s.ctx, sqls)
	if err != nil {
		zap.L().Error("查询失败", zap.Error(err))
		snag.Panic("查询失败")
	}
	defer func(rows *stdsql.Rows) {
		err = rows.Close()
		if err != nil {
			snag.Panic("rows close error")
		}
	}(rows)

	data := make([]*promotion.MemberTeamRes, 0)

	for rows.Next() {
		item := &promotion.MemberTeamRows{}
		err = rows.Scan(&item.ID, &item.Phone, &item.Name, &item.Level, &item.SubscribeStatus, &item.SubscribeStartAt, &item.RenewalCount)
		if err != nil {
			zap.L().Error("查询失败", zap.Error(err))
			snag.Panic("查询失败")
		}

		row := &promotion.MemberTeamRes{
			ID:           item.ID,
			Level:        item.Level.String(),
			RenewalCount: item.RenewalCount,
		}

		row.Name = item.Name.String
		row.Phone = item.Phone
		// 根据路由判断是否需要Mask Name
		if ctx.Path() != "/manager/v1/promotion/member/team/:id" {
			row.Name = s.MaskName(item.Name.String)
			row.Phone = s.MaskSensitiveInfo(item.Phone, 3, 4)
		}

		if item.SubscribeStatus.Valid {
			status := uint8(item.SubscribeStatus.Int64)
			var name string
			if status == model.SubscribeStatusInactive || status == model.SubscribeStatusCanceled {
				status = promotion.SubscribeStatusInactive.Value()
				name = promotion.SubscribeStatusInactive.String()
			} else if status == model.SubscribeStatusUsing || status == model.SubscribeStatusPaused || status == model.SubscribeStatusOverdue {
				status = promotion.SubscribeStatusUsing.Value()
				name = promotion.SubscribeStatusUsing.String()
			} else {
				status = promotion.SubscribeStatusUnSubscribed.Value()
				name = promotion.SubscribeStatusUnSubscribed.String()
			}
			row.SubscribeStatus = status
			row.SubscribeStatusName = name
		}
		if item.SubscribeStartAt.Valid {
			row.SubscribeStartAt = item.SubscribeStartAt.Time.Format("2006-01-02")
		}

		data = append(data, row)
	}

	// 统计总数
	statistics := s.TeamStatistics(req)

	pageReq := req.PaginationReq
	pages := pageReq.GetPages(int(statistics.Total))
	res := model.PaginationRes{
		Pagination: model.Pagination{
			Current: pageReq.GetCurrent(),
			Pages:   pages,
			Total:   int(statistics.Total),
		},
		Items: data,
	}

	return res
}

// MemberTeamFilter 筛选
func (s *promotionMemberService) MemberTeamFilter(req *promotion.MemberTeamReq, sql *string) {
	// 查询条件
	if req.Keyword != nil {
		*sql += fmt.Sprintf(" AND (m.phone LIKE '%%%s%%' OR m.name LIKE '%%%s%%' OR r.name LIKE '%%%s%%') ", *req.Keyword, *req.Keyword, *req.Keyword)
	}
	if req.SubscribeStatus != nil {
		if *req.SubscribeStatus == promotion.SubscribeStatusInactive.Value() {
			*sql += fmt.Sprintf(" AND (s.status = %d OR s.status IS NULL OR s.status = %d)", model.SubscribeStatusInactive, model.SubscribeStatusCanceled)
		} else if *req.SubscribeStatus == promotion.SubscribeStatusUsing.Value() {
			*sql += fmt.Sprintf(" AND (s.status = %d OR s.status = %d OR s.status = %d)", model.SubscribeStatusUsing, model.SubscribeStatusPaused, model.SubscribeStatusOverdue)
		} else {
			*sql += fmt.Sprintf(" AND s.status = %d ", *req.SubscribeStatus)
		}
	}
	if req.Level != nil {
		*sql += fmt.Sprintf(" AND mh.level = %d ", *req.Level)
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start).Format(carbon.DateTimeLayout)
		end := tools.NewTime().ParseNextDateStringX(*req.End).Format(carbon.DateTimeLayout)
		*sql += fmt.Sprintf(" AND s.start_at >= '%s' AND s.start_at < '%s' ", start, end)
	}
}

// TeamStatistics 会员团队统计
func (s *promotionMemberService) TeamStatistics(req *promotion.MemberTeamReq) *promotion.MemberTeamStatisticsRes {
	sqls := `
			WITH RECURSIVE member_hierarchy AS (
                SELECT referred_member_id, 1 AS level,rider_id,created_at
                FROM promotion_referrals
				` + fmt.Sprintf(" WHERE referring_member_id = %d", req.ID) + `
                UNION ALL
                SELECT mr.referred_member_id, mh.level + 1 AS level,mr.rider_id,mr.created_at
                FROM member_hierarchy mh
                INNER JOIN promotion_referrals as mr ON mh.referred_member_id = mr.referring_member_id
				WHERE mh.level < 2
			)
			SELECT
			  COUNT(*) AS total,
			  COALESCE(SUM(CASE WHEN level = 1 THEN 1 ELSE 0 END),0) AS firstLevelCount,
			  COALESCE(SUM(CASE WHEN level = 2 THEN 1 ELSE 0 END),0) AS secondLevelCount
			FROM member_hierarchy mh
			JOIN promotion_member m ON m.id = mh.referred_member_id
			LEFT JOIN rider r ON m.rider_id = r.id
			LEFT JOIN subscribe s ON r.id = s.rider_id
			    AND s.created_at = (
			        SELECT MAX(created_at)
			        FROM subscribe
			        WHERE rider_id = r.id
			    )
			WHERE mh.referred_member_id IN (
			    SELECT referred_member_id
			    FROM member_hierarchy
			)
`
	s.MemberTeamFilter(req, &sqls)

	rows, err := ent.Database.QueryContext(s.ctx, sqls)
	if err != nil {
		snag.Panic(err)
	}
	defer func(rows *stdsql.Rows) {
		err = rows.Close()
		if err != nil {
			snag.Panic(err)
		}
	}(rows)

	res := &promotion.MemberTeamStatisticsRes{}
	for rows.Next() {
		err = rows.Scan(&res.Total, &res.FirstLevelCount, &res.SecondLevelCount)
		if err != nil {
			return nil
		}
	}
	return res
}

// SetCommission 会员设置返佣方案
func (s *promotionMemberService) SetCommission(req *promotion.MemberCommissionReq) {
	mem, _ := ent.Database.PromotionMember.QueryNotDeleted().WithCommissions(
		func(query *ent.PromotionMemberCommissionQuery) {
			query.WithCommission()
		}).Where(promotionmember.IDEQ(req.ID)).First(s.ctx)
	if mem == nil {
		snag.Panic("会员不存在")
	}
	if req.Rule == nil {
		snag.Panic("规则不能为空")
	}

	s.CommissionPlanIsConfigured(mem, req)

	// 自定义返佣方案
	commissionTypeValue := promotion.CommissionType(2)
	mc := NewPromotionCommissionService().Create(&promotion.CommissionCreateReq{
		Name:     "自定义返佣",
		Rule:     *req.Rule,
		Type:     &commissionTypeValue,
		MemberID: &req.ID,
		Desc:     req.Desc,
	})

	bulk := make([]*ent.PromotionCommissionPlanCreate, len(req.PlanID))
	for i, v := range req.PlanID {
		bulk[i] = ent.Database.PromotionCommissionPlan.Create().SetPlanID(v).SetCommissionID(mc.ID).SetMemberID(req.ID)
	}
	ent.Database.PromotionCommissionPlan.CreateBulk(bulk...).ExecX(s.ctx)

	ent.Database.PromotionMemberCommission.Create().SetCommissionID(mc.ID).SetMemberID(req.ID).ExecX(s.ctx)

}

// CommissionPlanIsConfigured 判断返佣方案骑士卡是否已被配置
func (s *promotionMemberService) CommissionPlanIsConfigured(mem *ent.PromotionMember, req *promotion.MemberCommissionReq) {

	// 查询返佣方案
	pc, _ := ent.Database.PromotionCommissionPlan.QueryNotDeleted().Where(
		promotioncommissionplan.MemberID(mem.ID),
		promotioncommissionplan.PlanIDIn(req.PlanID...),
	).Select(promotioncommissionplan.FieldPlanID).WithPlan().All(s.ctx)
	if len(pc) >= len(req.PlanID) {
		planinfo := make(map[uint64]string)
		for _, v := range pc {
			planinfo[v.ID] = fmt.Sprintf("骑士卡 %s - %s - %.2f", v.Edges.Plan.Name, v.Edges.Plan.Model, v.Edges.Plan.Price)
		}
		snag.Panic(fmt.Sprintf("骑士卡 %v 已被配置", planinfo))
	}

}

// UploadAvatar 更新会员头像
func (s *promotionMemberService) UploadAvatar(ctx *app.PromotionContext) promotion.UploadAvatar {
	f, err := ctx.FormFile("avatar")
	if err != nil {
		snag.Panic("上传失败")
	}

	src, err := f.Open()
	if err != nil {
		snag.Panic("上传图片失败")
	}

	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	// 确保只接受指定的图片格式
	ext := filepath.Ext(f.Filename)
	if !NewFeedback().IsValidImageExtension(ext) {
		snag.Panic("只支持jpg、jpeg、png格式的图片")
	}

	// 生成相对路径
	randomNum := rand.Intn(1000) // 生成一个随机数，用于防止同一秒钟上传多个文件时的冲突
	r := filepath.Join("promotion", "avatar", fmt.Sprintf("%s%d%s", time.Now().Format(carbon.ShortDateTimeLayout), randomNum, ext))

	err = ali.NewOss().Bucket.PutObject(r, src)
	if err != nil {
		zap.L().Error("上传图片失败", zap.Error(err))
		snag.Panic("上传图片失败")
	}

	// 更新会员头像
	ent.Database.PromotionMember.UpdateOneID(ctx.Member.ID).SetAvatarURL(r).SaveX(s.ctx)
	return promotion.UploadAvatar{Avatar: r}
}

// UpgradeMemberLevel  会员升级
func (s *promotionMemberService) UpgradeMemberLevel(tx *ent.Tx, memberID, addGrowthValue uint64) error {
	var curLevel uint64

	mem, _ := s.GetMemberById(memberID)
	if mem == nil {
		return errors.New("会员不存在")
	}

	curLevel = 0
	if mem.Edges.Level != nil {
		curLevel = mem.Edges.Level.Level
	}
	// 获取下一级会员等级 大于当前等级 并且按等级升序排列
	setLevel, _ := ent.Database.PromotionLevel.QueryNotDeleted().Where(promotionlevel.LevelGT(curLevel)).Order(ent.Asc(promotionlevel.FieldLevel)).First(s.ctx)
	if setLevel == nil {
		return errors.New("无下一级会员等级")
	}

	q := tx.PromotionMember.UpdateOne(mem)
	// 当前成长值
	currentGrowthValue := mem.CurrentGrowthValue + addGrowthValue
	if setLevel != nil { // 有下一级会员等级 则判断是否需要升级 无下一级会员等级 则不需要升级
		// 判断当前成长值是否大于等于下一级会员等级所需成长值
		if currentGrowthValue >= setLevel.GrowthValue {
			// 升级后的成长值
			currentGrowthValue = currentGrowthValue - setLevel.GrowthValue
			// 升级会员等级
			q.SetLevelID(setLevel.ID)
		}
	}
	return q.AddTotalGrowthValue(int64(addGrowthValue)).SetCurrentGrowthValue(currentGrowthValue).Exec(s.ctx)
}
