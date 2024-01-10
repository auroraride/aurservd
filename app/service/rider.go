// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"golang.org/x/exp/slices"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ali"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/baidu"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/business"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/contract"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/internal/ent/person"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/riderfollowup"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
	"github.com/auroraride/aurservd/pkg/utils"
)

type riderService struct {
	cacheKeyPrefix string

	ctx      context.Context
	orm      *ent.RiderClient
	modifier *model.Modifier
	rider    *ent.Rider
}

func NewRider() *riderService {
	return &riderService{
		cacheKeyPrefix: "RIDER_",
		ctx:            context.Background(),
		orm:            ent.Database.Rider,
	}
}

func NewRiderWithModifier(m *model.Modifier) *riderService {
	s := NewRider()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func NewRiderWithRider(u *ent.Rider) *riderService {
	s := NewRider()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, u)
	s.rider = u
	return s
}

// GetRiderById 根据ID获取骑手及其实名状态
func (s *riderService) GetRiderById(id uint64) (u *ent.Rider, err error) {
	return ent.Database.Rider.
		QueryNotDeleted().
		WithPerson().
		WithEnterprise().
		WithEnterprise().
		WithStation().
		Where(rider.ID(id)).
		First(s.ctx)
}

// IsAuthed 是否已认证
func (s *riderService) IsAuthed(u *ent.Rider) bool {
	return u.Edges.Person != nil && !model.PersonAuthStatus(u.Edges.Person.Status).RequireAuth()
}

// IsNewDevice 检查是否是新设备
func (s *riderService) IsNewDevice(u *ent.Rider, device *model.Device) bool {
	return u.LastDevice != device.Serial || u.IsNewDevice
}

// IsBanned 骑手是否被拉黑
func (s *riderService) IsBanned(u *ent.Rider) bool {
	p := u.Edges.Person
	return p != nil && p.Banned
}

// IsBlocked 骑手是否被封禁
func (s *riderService) IsBlocked(u *ent.Rider) bool {
	return u.Blocked
}

// Signin 骑手登录
func (s *riderService) Signin(device *model.Device, req *model.RiderSignupReq) (res *model.RiderSigninRes) {
	NewSms().VerifyCodeX(req.Phone, req.SmsId, req.SmsCode)

	ctx := context.Background()
	orm := ent.Database.Rider
	var u *ent.Rider
	var err error

	u, err = orm.QueryNotDeleted().Where(rider.Phone(req.Phone)).WithPerson().WithEnterprise().First(ctx)
	if err != nil {
		// 创建骑手
		u, err = orm.Create().
			SetPhone(req.Phone).
			SetLastDevice(device.Serial).
			SetDeviceType(device.Type.Value()).
			Save(ctx)
		if err != nil {
			snag.Panic(err)
		}
	}

	// 判定用户是否被封禁
	if s.IsBanned(u) {
		snag.Panic(snag.StatusForbidden, ar.BannedMessage)
	}

	token := xid.New().String() + utils.RandTokenString()
	key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, u.ID)

	// 删除旧的token
	if old := cache.Get(ctx, key).Val(); old != "" {
		cache.Del(ctx, key)
		cache.Del(ctx, old)
	}

	// 更新设备
	if u.LastDevice != device.Serial {
		s.SetNewDevice(u, device)
	}

	res = s.Profile(u, device, token)

	// 设置登录token
	s.ExtendTokenTime(u.ID, token)

	return
}

// Signout 强制登出
func (s *riderService) Signout(u *ent.Rider) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, u.ID)
	token := cache.Get(ctx, key).Val()
	cache.Del(ctx, key)
	cache.Del(ctx, token)
}

// SetNewDevice 更新用户设备
func (s *riderService) SetNewDevice(u *ent.Rider, device *model.Device) {
	// isNew := true
	// if ar.Config.App.Debug.Phone[u.Phone] {
	// 	isNew = false
	// }
	// // TODO 暂时跳过人脸校验
	// isNew = false
	_, err := ent.Database.Rider.
		UpdateOneID(u.ID).
		SetLastDevice(device.Serial).
		SetDeviceType(device.Type.Value()).
		SetIsNewDevice(false). // 暂时跳过人脸校验 --by: 曹博文 2022-10-24 13:01
		Save(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	u.IsNewDevice = true
}

// GetFaceAuthUrl 获取实名验证URL
func (s *riderService) GetFaceAuthUrl(c *app.RiderContext) string {
	uri, token := baidu.NewFace().GetAuthenticatorUrl()
	cache.Set(context.Background(), token, s.GeneratePrivacy(c), 30*time.Minute)
	return uri
}

// GetFaceUrl 获取人脸校验URL
func (s *riderService) GetFaceUrl(c *app.RiderContext) string {
	p := c.Rider.Edges.Person
	uri, token := baidu.NewFace().GetFaceUrl(p.Name, p.IDCardNumber)
	cache.Set(context.Background(), token, s.GeneratePrivacy(c), 30*time.Minute)
	return uri
}

// FaceAuthResult 获取并更新人脸实名验证结果
func (s *riderService) FaceAuthResult(c *app.RiderContext, token string) (success bool) {
	if !s.ComparePrivacy(c) {
		snag.Panic("验证失败")
	}
	u := c.Rider
	data, err := baidu.NewFace().AuthenticatorResult(token)
	if err != nil {
		return
	}

	status := model.PersonAuthenticated.Value()
	success = data.Success

	if !success {
		status = model.PersonAuthenticationFailed.Value()
	}

	res := data.Result
	detail := res.IdcardOcrResult

	var remark string
	// 判定是否年满18周岁
	birthday := carbon.Parse(detail.Birthday).ToStdTime().AddDate(18, 0, 0)
	// 未年满18岁认证标记为失败
	if birthday.After(time.Now()) {
		remark = "未年满18岁"
		status = model.PersonAuthenticationFailed.Value()
	}

	vr := &model.FaceVerifyResult{
		Birthday:       detail.Birthday,
		IssueAuthority: detail.IssueAuthority,
		Address:        detail.Address,
		Gender:         detail.Gender,
		Nation:         detail.Nation,
		ExpireTime:     detail.ExpireTime,
		Name:           detail.Name,
		IssueTime:      detail.IssueTime,
		IdCardNumber:   detail.IdCardNumber,
		Score:          res.VerifyResult.Score,
		LivenessScore:  res.VerifyResult.LivenessScore,
		Spoofing:       res.VerifyResult.Spoofing,
	}

	// 上传图片到七牛云
	var fm, pm, nm string
	oss := ali.NewOss()
	prefix := fmt.Sprintf("%s-%s/%s-", res.IdcardOcrResult.Name, res.IdcardOcrResult.IdCardNumber, time.Now().Format(carbon.ShortDateTimeLayout))
	if res.FaceImg != "" {
		fm = oss.UploadUrlFile(prefix+"face.jpg", res.FaceImg)
	}
	if res.IdcardImages.FrontBase64 != "" {
		pm = oss.UploadBase64ImageJpeg(prefix+"portrait.jpg", res.IdcardImages.FrontBase64)
	}
	if res.IdcardImages.BackBase64 != "" {
		nm = oss.UploadBase64ImageJpeg(prefix+"national.jpg", res.IdcardImages.BackBase64)
	}

	icNum := vr.IdCardNumber
	var id uint64
	id, err = ent.Database.Person.
		Create().
		SetStatus(status).
		SetIDCardNumber(icNum).
		SetName(vr.Name).
		SetAuthFace(fm).
		SetIDCardNational(nm).
		SetIDCardPortrait(pm).
		SetAuthResult(vr).
		SetAuthAt(time.Now()).
		OnConflictColumns(person.FieldIDCardNumber).
		UpdateNewValues().
		SetBaiduLogID(data.LogId).
		SetBaiduVerifyToken(token).
		SetRemark(remark).
		ID(context.Background())
	if err != nil {
		snag.Panic(err)
	}

	if success || u.PersonID == nil {
		// 判断ID是否等于实名认证的ID, 如果不是, 则删除
		if u.PersonID != nil && *u.PersonID != id {
			_ = ent.Database.Person.DeleteOneID(*u.PersonID).Exec(s.ctx)
		}
		// 更新骑手信息
		ri, err := ent.Database.Rider.
			UpdateOneID(u.ID).
			SetPersonID(id).
			SetLastFace(fm).
			SetIsNewDevice(false).
			SetName(vr.Name).
			SetIDCardNumber(icNum).
			Save(context.Background())
		if err != nil {
			snag.Panic(err)
		}
		// 如果骑手注册过推广小程序 判断是否需要绑定推荐关系
		NewPromotionReferralsService().RiderBindReferrals(ri)
	}

	return success
}

// FaceResult 获取人脸比对结果
func (s *riderService) FaceResult(c *app.RiderContext, token string) (success bool) {
	if !s.ComparePrivacy(c) {
		snag.Panic("验证失败")
	}
	u := c.Rider
	res, err := baidu.NewFace().FaceResult(token)
	if err != nil {
		snag.Panic(err)
		return
	}
	success = res.Success
	if !success {
		return
	}
	// 上传人脸图
	p := u.Edges.Person
	fm := ali.NewOss().UploadUrlFile(fmt.Sprintf("%s-%s/face-%s.jpg", p.Name, p.IDCardNumber, u.LastDevice), res.Result.Image)
	err = ent.Database.Rider.
		UpdateOneID(u.ID).
		SetLastFace(fm).
		SetIsNewDevice(false).
		Exec(context.Background())
	if err != nil {
		snag.Panic(err)
	}
	return
}

// UpdateContact 更新紧急联系人
func (s *riderService) UpdateContact(u *ent.Rider, contact *model.RiderContact) {
	// 判断紧急联系人手机号是否和当前骑手手机号一样
	if u.Phone == contact.Phone {
		snag.Panic("紧急联系人手机号不能是当前手机号")
	}
	err := ent.Database.Rider.UpdateOneID(u.ID).SetContact(contact).Exec(context.Background())
	if err != nil {
		snag.Panic(err)
	}
}

// GeneratePrivacy 获取实名认证或人脸识别限制条件
func (s *riderService) GeneratePrivacy(c *app.RiderContext) string {
	return fmt.Sprintf("%s-%d", c.Device.Serial, c.Rider.ID)
}

// ComparePrivacy 比对实名认证或人脸识别限制条件是否满足
func (s *riderService) ComparePrivacy(c *app.RiderContext) bool {
	return cache.Get(context.Background(), c.Param("token")).Val() == s.GeneratePrivacy(c)
}

// ExtendTokenTime 延长骑手登录有效期
func (s *riderService) ExtendTokenTime(id uint64, token string) {
	key := fmt.Sprintf("%s%d", s.cacheKeyPrefix, id)
	ctx := context.Background()
	cache.Set(ctx, key, token, 7*24*time.Hour)
	cache.Set(ctx, token, id, 7*24*time.Hour)
	_ = ent.Database.Rider.
		UpdateOneID(id).
		SetLastSigninAt(time.Now()).
		Exec(context.Background())
}

// GetRiderSampleInfo 获取骑手简单信息
func (*riderService) GetRiderSampleInfo(r *ent.Rider) model.RiderSampleInfo {
	return model.RiderSampleInfo{
		ID:    r.ID,
		Name:  r.Name,
		Phone: r.Phone,
	}
}

func (s *riderService) Status(u *ent.Rider) uint8 {
	status := model.RiderStatusNormal
	if u.Blocked {
		status = model.RiderStatusBlocked
	}
	p := u.Edges.Person
	if p != nil {
		if p.Banned {
			status = model.RiderStatusBanned
		}
	}
	return status
}

func (s *riderService) listFilter(req model.RiderListFilter) (q *ent.RiderQuery, info ar.Map) {
	var subqs []predicate.Subscribe

	info = make(ar.Map)
	q = ent.Database.Rider.
		QueryNotDeleted().
		WithPerson().
		WithOrders(func(oq *ent.OrderQuery) {
			oq.Where(
				order.Type(model.OrderTypeDeposit),
				order.Status(model.OrderStatusPaid),
			)
		}).
		WithSubscribes(func(sq *ent.SubscribeQuery) {
			sq.WithCity().Order(ent.Desc(subscribe.FieldCreatedAt)).WithEbike().WithBrand().WithPlan()
		}).
		WithContracts(func(cq *ent.ContractQuery) {
			cq.Where(contract.DeletedAtIsNil(), contract.Status(model.ContractStatusSuccess.Value())).Order(ent.Desc(contract.FieldCreatedAt))
		}).
		WithEnterprise().
		WithStation().
		WithBattery().
		Order(ent.Desc(rider.FieldCreatedAt))
	if req.Keyword != nil {
		info["关键词"] = *req.Keyword
		// 判定是否id字段
		q.Where(
			rider.Or(
				rider.NameContainsFold(*req.Keyword),
				rider.IDCardNumberContainsFold(*req.Keyword),
				rider.PhoneContainsFold(*req.Keyword),
			),
		)
	}
	if req.Start != nil {
		info["开始日期"] = *req.Start
		q.Where(rider.CreatedAtGTE(tools.NewTime().ParseDateStringX(*req.Start)))
	}
	if req.End != nil {
		info["结束日期"] = *req.End
		q.Where(rider.CreatedAtLT(tools.NewTime().ParseNextDateStringX(*req.End)))
	}
	if req.Modified != nil {
		m := *req.Modified
		if m {
			info["修改状态"] = "已被修改"
			q.Where(rider.LastModifierNotNil())
		} else {
			info["修改状态"] = "未被修改"
			q.Where(rider.LastModifierIsNil())
		}
	}
	if req.Status != nil {
		rs := *req.Status
		switch rs {
		case model.RiderStatusNormal:
			info["用户状态"] = "未认证"
			q.Where(
				rider.Blocked(false),
				rider.Or(
					rider.PersonIDIsNil(),
					rider.HasPersonWith(person.Banned(false)),
				),
			)
		case model.RiderStatusBlocked:
			info["用户状态"] = "已禁用"
			q.Where(rider.Blocked(true))
		case model.RiderStatusBanned:
			info["用户状态"] = "已封禁"
			q.Where(rider.HasPersonWith(person.Banned(true)))
		}
	}
	if req.AuthStatus != nil {
		ra := *req.AuthStatus
		info["认证状态"] = req.AuthStatus.String()
		switch ra {
		case model.PersonUnauthenticated:
			q.Where(
				rider.Or(
					rider.PersonIDIsNil(),
					rider.HasPersonWith(person.Status(model.PersonUnauthenticated.Value())),
				),
			)
		default:
			q.Where(rider.HasPersonWith(person.Status(ra.Value())))
		}
	}
	if req.SubscribeStatus != nil {
		key := "业务状态"
		rss := *req.SubscribeStatus
		switch rss {
		case 11:
			// 即将到期
			subqs = append(
				subqs,
				subscribe.Or(
					subscribe.And(
						subscribe.EnterpriseIDIsNil(),
						subscribe.StatusIn(model.SubscribeNotUnSubscribed()...),
						subscribe.RemainingLTE(3),
					),
					subscribe.And(
						subscribe.EnterpriseIDNotNil(),
						subscribe.Status(model.SubscribeStatusUsing),
						subscribe.AgentEndAtLTE(carbon.CreateFromStdTime(tools.NewTime().WillEnd(time.Now(), 3, true)).EndOfDay().ToStdTime()),
					),
				),
			)
		case 99:
			q.Where(rider.Not(rider.HasSubscribes()))
		case model.SubscribeStatusUnSubscribed:
			subqs = append(
				subqs,
				subscribe.Status(model.SubscribeStatusUnSubscribed),
			)
		default:
			subqs = append(subqs, subscribe.Status(rss))
		}
		info[key] = map[uint8]string{
			0:  "未激活",
			1:  "计费中",
			2:  "寄存中",
			3:  "已逾期",
			4:  "已退订",
			5:  "已取消",
			11: "即将到期",
			99: "未使用",
		}[rss]
	}
	if req.PlanID != nil {
		info["骑士卡"] = ent.NewExportInfo(*req.PlanID, plan.Table)
		// q.Where(rider.HasSubscribesWith(subscribe.PlanID(*req.PlanID)))
		subqs = append(subqs, subscribe.PlanID(*req.PlanID))
	}
	if req.Enterprise != nil && *req.Enterprise != 0 {
		key := "是否团签"
		var value string
		if *req.Enterprise == 1 {
			value = "是"
			if req.EnterpriseID == nil {
				q.Where(rider.EnterpriseIDNotNil())
			} else {
				info["团签企业"] = ent.NewExportInfo(*req.EnterpriseID, enterprise.Table)
				q.Where(rider.EnterpriseID(*req.EnterpriseID))
			}
		} else {
			value = "否"
			q.Where(rider.EnterpriseIDIsNil())
		}
		info[key] = value
	}

	if req.CityID != nil {
		info["城市"] = ent.NewExportInfo(*req.CityID, city.Table)
		subqs = append(subqs, subscribe.CityID(*req.CityID))
	}

	if req.Remaining != nil {
		arr := strings.Split(*req.Remaining, ",")

		subqs = append(subqs, subscribe.StatusNotIn(model.SubscribeStatusUnSubscribed, model.SubscribeStatusCanceled))
		r1, _ := strconv.Atoi(strings.TrimSpace(arr[0]))
		subqs = append(subqs, subscribe.RemainingGTE(r1))
		if len(arr) > 1 {
			r2, _ := strconv.Atoi(strings.TrimSpace(arr[1]))
			if r2 > 0 {
				if r1 > r2 {
					snag.Panic("区间错误")
				}
				subqs = append(subqs, subscribe.RemainingLTE(r2))
			}
			info["骑士卡剩余天数"] = fmt.Sprintf("%d - %d", r1, r2)
		} else {
			info["骑士卡剩余天数"] = fmt.Sprintf("> %d", r1)
		}
	}

	if req.Suspend != nil {
		if *req.Suspend {
			subqs = append(subqs, subscribe.SuspendAtNotNil(), subscribe.EndAtIsNil())
			info["暂停扣费"] = "是"
		} else {
			subqs = append(subqs, subscribe.SuspendAtIsNil(), subscribe.EndAtIsNil())
			info["暂停扣费"] = "否"
		}
	}

	if req.Model != nil {
		info["电池型号"] = req.Model
		subqs = append(subqs, subscribe.Model(*req.Model))
		// 如果筛选了电池型号但是未筛选订阅状态时, 默认订阅状态为: 非退订 && 非取消
		if req.SubscribeStatus == nil {
			subqs = append(subqs, subscribe.StatusNotIn(model.SubscribeStatusCanceled, model.SubscribeStatusUnSubscribed))
		}
	}

	if req.EbikeBrandID != nil {
		info["电车型号"] = ent.NewExportInfo(*req.EbikeBrandID, ebikebrand.Table)
		subqs = append(subqs, subscribe.BrandID(*req.EbikeBrandID))
		// 如果筛选了电车型号但是未筛选订阅状态时, 默认订阅状态为: 非退订 && 非取消
		if req.SubscribeStatus == nil {
			subqs = append(subqs, subscribe.StatusNotIn(model.SubscribeStatusCanceled, model.SubscribeStatusUnSubscribed))
		}
	}

	if req.BatteryID != nil {
		info["电池编码"] = ent.NewExportInfo(*req.BatteryID, battery.Table)
		q.Where(rider.HasBatteryWith(battery.ID(*req.BatteryID)))
	}

	if len(subqs) > 0 {
		q.Where(rider.HasSubscribesWith(subqs...))
	}
	return
}

// List 骑手列表
func (s *riderService) List(req *model.RiderListReq) *model.PaginationRes {
	q, _ := s.listFilter(req.RiderListFilter)

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Rider) model.RiderItem {
			return s.detailRiderItem(item)
		},
	)
}

func (s *riderService) detailRiderItem(item *ent.Rider) model.RiderItem {
	p := item.Edges.Person
	ri := model.RiderItem{
		ID:                item.ID,
		Phone:             item.Phone,
		Status:            model.RiderStatusNormal,
		AuthStatus:        model.PersonUnauthenticated,
		Contact:           item.Contact,
		Points:            item.Points,
		Balance:           0,
		ExchangeLimit:     item.ExchangeLimit,
		ExchangeFrequency: item.ExchangeFrequency,
	}
	e := item.Edges.Enterprise
	if e != nil {
		ri.Enterprise = &model.Enterprise{
			ID:    e.ID,
			Name:  e.Name,
			Agent: e.Agent,
		}
	}

	st := item.Edges.Station
	if st != nil {
		ri.Station = &model.EnterpriseStation{
			ID:   st.ID,
			Name: st.Name,
		}
	}

	if item.Blocked {
		ri.Status = model.RiderStatusBlocked
	}
	if p != nil {
		ri.Name = p.Name
		ri.AuthStatus = model.PersonAuthStatus(p.Status)
		if p.Banned {
			ri.Status = model.RiderStatusBanned
		}
		if p.AuthResult != nil {
			ri.Address = p.AuthResult.Address
		}
		ri.Person = &model.Person{
			IDCardNumber:   p.IDCardNumber,
			IDCardPortrait: p.IDCardPortrait,
			IDCardNational: p.IDCardNational,
			AuthFace:       p.AuthFace,
		}
	}

	// 获取合同
	contracts := item.Edges.Contracts
	if len(contracts) > 0 && contracts[0].Files != nil {
		ri.Contract = contracts[0].Files[0]
	}

	if item.Edges.Orders != nil && len(item.Edges.Orders) > 0 {
		ri.Deposit = item.Edges.Orders[0].Amount
	}

	if item.Edges.Subscribes != nil && len(item.Edges.Subscribes) > 0 {
		sub := item.Edges.Subscribes[0]

		remaining := sub.Remaining

		// 代理骑手未退组剩余天数
		if sub.AgentEndAt != nil && sub.EndAt == nil {
			remaining = tools.NewTime().LastDaysToNow(*sub.AgentEndAt)
		}

		ri.Subscribe = &model.RiderItemSubscribe{
			ID:          sub.ID,
			Status:      sub.Status,
			Remaining:   remaining,
			Model:       sub.Model,
			Suspend:     sub.SuspendAt != nil,
			Formula:     sub.Formula,
			Type:        model.SubscribeTypeBattery,
			Ebike:       NewEbike().Detail(sub.Edges.Ebike, sub.Edges.Brand),
			Intelligent: sub.Intelligent,
		}
		if sub.BrandID != nil {
			ri.Subscribe.Type = model.SubscribeTypeEbike
		}
		if sub.AgentEndAt != nil {
			ri.Subscribe.AgentEndAt = sub.AgentEndAt.Format(carbon.DateLayout)
		}
		ri.City = &model.City{
			ID: sub.CityID,
		}
		if sub.Edges.City != nil {
			ri.City.Name = sub.Edges.City.Name
		}

		pl := sub.Edges.Plan
		pn := "单电"

		if sub.BrandID != nil {
			pn = "车电"
		}

		if pl != nil {
			pn = fmt.Sprintf("%s [%s]", pn, pl.GetExportInfo())
		}

		ri.PlanName = pn
	}
	if item.DeletedAt != nil {
		ri.DeletedAt = item.DeletedAt.Format(carbon.DateTimeLayout)
		ri.Remark = item.Remark
	}

	// 获取电池
	bat := item.Edges.Battery
	if bat != nil {
		ri.Battery = &model.Battery{
			ID:    bat.ID,
			SN:    bat.Sn,
			Model: bat.Model,
		}
	}
	return ri
}

func (s *riderService) ListExport(req *model.RiderListExport) model.ExportRes {
	q, info := s.listFilter(req.RiderListFilter)
	return NewExportWithModifier(s.modifier).Start("骑手列表", req, info, req.Remark, func(path string) {
		items, _ := q.All(s.ctx)

		var rows tools.ExcelItems
		title := []any{
			"城市",     // 0
			"门店",     // 1
			"激活办理人",  // 2
			"退租办理人",  // 3
			"骑手",     // 4
			"电话",     // 5
			"证件",     // 6
			"户籍",     // 7
			"团签",     // 8
			"押金",     // 9
			"订阅",     // 10
			"暂停",     // 11
			"订单开始时间", // 12
			"订单到期时间", // 13
			"电池",     // 14
			"剩余",     // 15
			"逾期费用",   // 16
			"状态",     // 17
			"认证",     // 18
			"紧急联系",   // 19
			"注册时间",   // 20
			"电车型号",   // 21
			"车牌号",    // 22
			"车架号",    // 23
			"跟进详情",   // 24
		}
		rows = append(rows, title)
		for _, item := range items {
			detail := s.detailRiderItem(item)
			row := []any{
				"",
				"",
				"",
				"",
				detail.Name,
				detail.Phone,
				"",
				detail.Address,
				"",
				detail.Deposit,
				"",
				"否",
				"",
				"",
				"",
				"",
				"",
				[]string{"正常", "正常", "禁用", "黑名单"}[detail.Status],
				detail.AuthStatus.String(),
				"",
				item.CreatedAt.Format(carbon.DateTimeLayout),
				"",
				"",
				"",
				"",
			}
			if detail.City != nil {
				row[0] = detail.City.Name
			}
			if detail.Person != nil {
				row[6] = detail.Person.IDCardNumber
			}
			// 团签
			var group string
			if detail.Enterprise != nil {
				group = detail.Enterprise.Name
				if detail.Station != nil {
					group += "-" + detail.Station.Name
				}
			}
			row[8] = group
			if detail.Subscribe != nil {
				row[10] = model.SubscribeStatusText(detail.Subscribe.Status)
				if detail.Subscribe.Suspend {
					row[11] = "是"
				}
				row[14] = detail.Subscribe.Model
				row[15] = detail.Subscribe.Remaining
				bike := detail.Subscribe.Ebike
				if bike != nil {
					row[21] = bike.Brand.Name
					if bike.SN != "" {
						row[23] = bike.SN
					}
					if bike.Plate != nil {
						row[22] = *bike.Plate
					}
				}
			}
			if item.Contact != nil {
				row[19] = item.Contact.String()
			}
			// 激活办理人
			var activeOperator string
			// 退租办理人
			var unsubscribeOperator string

			bizList, _ := ent.Database.Business.QueryNotDeleted().
				WithEmployee().WithStore().WithAgent().
				Where(business.RiderID(item.ID), business.TypeIn(business.TypeActive, business.TypeUnsubscribe)).
				Order(ent.Desc(business.FieldCreatedAt)).
				Limit(2).
				All(s.ctx)

			// 取2条数据，判定最后一条数据若为退租业务，则需要轮询判定和展示激活办理人和退租办理人；
			// 若最后一条为激活业务，只需要展示激活办理人，不做其他判定，不展示退租办理人
			for i, biz := range bizList {
				// 获取操作人员
				var operator string
				switch {
				case biz.EmployeeID == nil && biz.CabinetID == nil && biz.AgentID == nil && biz.Creator != nil:
					// 操作人是平台
					operator = biz.Creator.Name + "-" + biz.Creator.Phone
				case biz.CabinetID != nil:
					// 操作人是骑手
					operator = item.Name + "-" + item.Phone
				case biz.EmployeeID != nil:
					// 操作人是店员
					if biz.Edges.Employee != nil {
						operator = biz.Edges.Employee.Name + "-" + biz.Edges.Employee.Phone
					}
				case biz.AgentID != nil:
					// 操作人是代理
					if biz.Edges.Agent != nil {
						operator = biz.Edges.Agent.Name + "-" + biz.Edges.Agent.Phone
					}
				}

				// （最新记录）记录门店
				if biz.Edges.Store != nil && i == 0 {
					row[1] = biz.Edges.Store.Name
				}

				// 记录激活操作人
				if biz.Type == business.TypeActive {
					activeOperator = operator
				}

				// 最新记录为退租业务时记录退租操作人
				if biz.Type == business.TypeUnsubscribe && i == 0 {
					unsubscribeOperator = operator
				}
			}

			// 办理人
			row[2] = activeOperator
			row[3] = unsubscribeOperator

			// 订单开始时间、订单结束时间
			subs := item.Edges.Subscribes
			if len(subs) > 0 {
				sub := subs[0]
				if sub.StartAt != nil {
					row[12] = sub.StartAt.Format(carbon.DateLayout)
				}
				if detail.Subscribe.AgentEndAt != "" {
					// 代理商处到期时间
					row[13] = detail.Subscribe.AgentEndAt
				} else {
					// 到期时间 当前时间+订阅剩余天数
					row[13] = carbon.Now().AddDays(sub.Remaining).ToDateString()
				}
				// 逾期费用
				if sub.Remaining < 0 {
					fee, _ := NewSubscribe().OverdueFee(sub)
					row[16] = fee
				}
			}
			// 跟进详情
			riderFollowUps, _ := ent.Database.RiderFollowUp.QueryNotDeleted().Where(riderfollowup.RiderID(item.ID)).All(s.ctx)
			if len(riderFollowUps) > 0 {
				var temp = make([]string, len(riderFollowUps))
				for k, v := range riderFollowUps {
					temp[k] = v.Remark
				}
				var riderFollowUpsDetail = strings.Join(temp, "-")
				row[24] = riderFollowUpsDetail
			}
			rows = append(rows, row)
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}

func (s *riderService) Query(id uint64) *ent.Rider {
	item, err := ent.Database.Rider.QueryNotDeleted().Where(rider.ID(id)).WithPerson().First(s.ctx)
	if err != nil || item == nil {
		snag.Panic("未找到骑手")
	}
	return item
}

// QueryForBusinessID 查找骑手并判定是否满足业务办理条件
func (s *riderService) QueryForBusinessID(riderID uint64) (u *ent.Rider, err error) {
	u = s.Query(riderID)
	err = s.Permission(u)
	return
}

// CheckForBusiness 骑手是否可办理业务
func (s *riderService) CheckForBusiness(u *ent.Rider) {
	err := s.Permission(u)
	if err != nil {
		snag.Panic(err)
	}
}

func (s *riderService) Permission(u *ent.Rider) (err error) {
	if u.Edges.Person == nil {
		u.Edges.Person, _ = u.QueryPerson().First(s.ctx)
	}
	if u.IsNewDevice {
		err = errors.New("骑手未人脸识别")
	}
	if !s.IsAuthed(u) {
		err = errors.New("骑手未实名")
	}
	if NewAssistance().Unpaid(u.ID) != nil {
		err = errors.New("救援订单未支付")
	}
	if s.IsBlocked(u) {
		err = errors.New("骑手被封禁")
	}
	if s.IsBanned(u) {
		err = errors.New("骑手被拉黑")
	}
	return
}

// Block 封锁/解封骑手账户
func (s *riderService) Block(req *model.RiderBlockReq) {
	item := s.Query(req.ID)
	if req.Block == item.Blocked {
		snag.Panic("骑手已是封禁状态")
	}
	_, err := s.orm.UpdateOne(item).SetBlocked(req.Block).Save(s.ctx)
	if err != nil {
		snag.Panic(err)
	}
	nb := "未封禁"
	bd := "已封禁"
	ol := logging.NewOperateLog().SetRef(item).SetModifier(s.modifier)
	if req.Block {
		// 封禁
		ol.SetOperate(model.OperateRiderBLock).SetDiff(nb, bd)
	} else {
		ol.SetOperate(model.OperateRiderUnBLock).SetDiff(bd, nb)
	}
	ol.Send()
}

// DepositOrder 获取骑手押金订单
func (s *riderService) DepositOrder(riderID uint64) *ent.Order {
	o, _ := ent.Database.Order.QueryNotDeleted().Where(
		order.RiderID(riderID),
		order.Status(model.OrderStatusPaid),
		order.Type(model.OrderTypeDeposit),
		order.DeletedAtIsNil(),
	).First(s.ctx)
	return o
}

// DepositPaid 已缴押金
func (s *riderService) DepositPaid(riderID uint64) model.RiderDepositRes {
	o := s.DepositOrder(riderID)
	res := model.RiderDepositRes{
		Deposit: 0,
	}
	if o != nil {
		res.Deposit = o.Amount
	}
	return res
}

// Deposit 获取用户应交押金
func (s *riderService) Deposit(riderID uint64) float64 {
	o := s.DepositOrder(riderID)
	if o != nil {
		return 0
	}
	f, _ := cache.Get(s.ctx, model.SettingDepositKey).Float64()
	return f
}

func (s *riderService) GetQrcode(id uint64) string {
	b, _ := tools.NewAESCrypto().Encrypt([]byte(fmt.Sprintf("%d", id)))
	return b
}

func (s *riderService) ParseQrcode(qrcode string) uint64 {
	b, _ := base64.StdEncoding.DecodeString(qrcode)
	str, _ := tools.NewAESCrypto().Decrypt(b)
	id, _ := strconv.ParseUint(str, 10, 64)
	return id
}

// Profile 获取用户资料
func (s *riderService) Profile(u *ent.Rider, device *model.Device, token string) *model.RiderSigninRes {
	subd, sub := NewSubscribe().RecentDetail(u.ID)
	profile := &model.RiderSigninRes{
		ID:              u.ID,
		Phone:           u.Phone,
		IsNewDevice:     s.IsNewDevice(u, device),
		IsContactFilled: u.Contact != nil,
		IsAuthed:        s.IsAuthed(u),
		Contact:         u.Contact,
		Qrcode:          s.GetQrcode(u.ID),
		Token:           token,
	}

	// 站点
	stat := u.Edges.Station
	if stat != nil {
		profile.Station = &model.EnterpriseStation{
			ID:   stat.ID,
			Name: stat.Name,
		}
	}

	// 订阅
	if sub != nil && slices.Contains(model.SubscribeNotUnSubscribed(), sub.Status) {
		profile.Subscribe = subd
		profile.CabinetBusiness = u.EnterpriseID == nil && sub.BrandID == nil
	}
	if u.Edges.Person != nil {
		profile.Name = u.Edges.Person.Name
	}
	en := u.Edges.Enterprise
	if en != nil {
		profile.Enterprise = &model.Enterprise{
			ID:    en.ID,
			Name:  en.Name,
			Agent: en.Agent,
		}
		profile.UseStore = !en.Agent || en.UseStore
		if en.Agent {
			profile.EnterpriseContact = &model.EnterpriseContact{
				Name:  en.ContactName,
				Phone: en.ContactPhone,
			}
		}
		// 判断是否能退出团签
		if subd != nil && (subd.Status == model.SubscribeStatusInactive || subd.Status == model.SubscribeStatusUnSubscribed) {
			profile.ExitEnterprise = true
		}

	} else {
		profile.OrderNotActived = silk.Bool(subd != nil && subd.Status == model.SubscribeStatusInactive)
		profile.Deposit = s.Deposit(u.ID)
		profile.UseStore = true
	}
	return profile
}

// GetLogs 获取用户操作日志
func (s *riderService) GetLogs(req *model.RiderLogReq) *model.PaginationRes {
	cfg := ar.Config.Aliyun.Sls

	u := s.Query(req.ID)
	query := fmt.Sprintf(`refTable:'rider' AND refId:%d`, u.ID)
	if req.Type != model.RiderLogTypeAll {
		ts, ok := model.RiderLogTypes[req.Type]
		if !ok {
			snag.Panic("类型错误")
		}
		and := make([]string, len(ts))
		for i, t := range ts {
			and[i] = fmt.Sprintf(`operate:%s`, t)
		}
		query += fmt.Sprintf(" AND %s", strings.Join(and, " OR "))
	}

	// 分页获取
	total := logging.GetCount(cfg.OperateLog, query, u.CreatedAt)
	pageReq := req.PaginationReq
	pages := pageReq.GetPages(total)

	// 查询结果
	result := logging.NewOperateLog().GetLogs(u.CreatedAt, query, int64(pageReq.GetOffset()), int64(pageReq.GetLimit()))
	b, _ := jsoniter.Marshal(result)
	items := make([]model.LogOperate, 0)
	_ = jsoniter.Unmarshal(b, &items)
	return &model.PaginationRes{
		Pagination: model.Pagination{
			Current: pageReq.GetCurrent(),
			Pages:   pages,
			Total:   total,
		},
		Items: items,
	}
}

// Delete 删除账户
func (s *riderService) Delete(req *model.IDParamReq) {
	u := s.Query(req.ID)
	q := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.RiderID(req.ID))

	if u.EnterpriseID != nil {
		q.Where(
			subscribe.StatusIn(model.SubscribeStatusUsing, model.SubscribeStatusOverdue),
		)
	} else {
		q.Where(
			subscribe.StatusNotIn(model.SubscribeStatusUnSubscribed, model.SubscribeStatusCanceled),
		)
	}

	sub, _ := q.First(s.ctx)
	if sub != nil {
		snag.Panic("骑手当前有使用中的订阅")
	}

	_, err := s.orm.SoftDeleteOneID(req.ID).Save(s.ctx)

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 软删除用户
		err = tx.Rider.SoftDeleteOne(u).Exec(s.ctx)
		if err != nil {
			return
		}

		// 软删除订阅
		if sub != nil && sub.EnterpriseID != nil && sub.Status == model.SubscribeStatusInactive {
			err = tx.Subscribe.DeleteOne(sub).Exec(s.ctx)
		}
		return
	})

	s.Signout(u)

	if err != nil {
		snag.Panic(err)
	}
}

func (s *riderService) NameFromID(id uint64) string {
	r, _ := ent.Database.Rider.QueryNotDeleted().WithPerson().Where(rider.ID(id)).First(s.ctx)
	if r == nil {
		return "-"
	}
	str := r.Phone
	p := r.Edges.Person
	if p != nil {
		str += " - " + p.Name
	}
	return str
}

func (s *riderService) QueryPhones(phones []string) (riders []*ent.Rider, ids []uint64, notfound []string) {
	notfound = make([]string, 0)
	m := make(map[string]struct{})
	for _, phone := range phones {
		m[phone] = struct{}{}
	}
	// 查询骑手
	riders, _ = ent.Database.Rider.QueryNotDeleted().Where(rider.PhoneIn(phones...)).All(s.ctx)
	if len(riders) == 0 {
		snag.Panic("全部骑手查询失败")
	}
	// 找到的手机号
	for _, r := range riders {
		delete(m, r.Phone)
		ids = append(ids, r.ID)
	}
	// 未找到的手机号
	for k := range m {
		notfound = append(notfound, k)
	}
	return
}

func (s *riderService) QueryPhone(phone string) (*ent.Rider, error) {
	return s.orm.Query().Where(rider.Phone(phone), rider.DeletedAtIsNil()).Order(ent.Desc(rider.FieldCreatedAt)).First(s.ctx)
}

func (s *riderService) QueryPhoneX(phone string) (rd *ent.Rider) {
	rd, _ = s.QueryPhone(phone)
	if rd == nil {
		snag.Panic("未找到有效骑手")
	}
	return
}

// ExchangeLimit 设置骑手换电限制
func (s *riderService) ExchangeLimit(req *model.RiderExchangeLimitReq) {
	r := s.Query(req.ID)
	updater := r.Update()
	if len(req.ExchangeLimit) == 0 {
		updater.ClearExchangeLimit()
	} else {
		if req.ExchangeLimit.Duplicate() {
			snag.Panic("设定重复")
		}
		req.ExchangeLimit.Sort()
		updater.SetExchangeLimit(req.ExchangeLimit)
	}
	_ = updater.Exec(s.ctx)

	// 记录日志
	go logging.NewOperateLog().
		SetRef(r).
		SetModifier(s.modifier).
		SetOperate(model.OperateExchangeLimit).
		SetDiff(r.ExchangeLimit.String(), req.ExchangeLimit.String()).
		Send()
}

// ExchangeFrequency 设置骑手换电频次
func (s *riderService) ExchangeFrequency(req *model.RiderExchangeFrequencyReq) {
	r := s.Query(req.ID)
	updater := r.Update()
	if len(req.ExchangeFrequency) == 0 {
		updater.ClearExchangeLimit()
	} else {
		if req.ExchangeFrequency.Duplicate() {
			snag.Panic("设定重复")
		}
		req.ExchangeFrequency.Sort()
		updater.SetExchangeFrequency(req.ExchangeFrequency)
	}
	_ = updater.Exec(s.ctx)

	// 记录日志
	go logging.NewOperateLog().
		SetRef(r).
		SetModifier(s.modifier).
		SetOperate(model.OperateExchangeFrequency).
		SetDiff(r.ExchangeFrequency.String(), req.ExchangeFrequency.String()).
		Send()
}

// GetRiderNameById 根据ID获取骑手信息
func (s *riderService) GetRiderNameById(id uint64) *model.RiderSampleInfo {
	ri, _ := ent.Database.Rider.
		QueryNotDeleted().
		Where(rider.ID(id)).
		First(s.ctx)
	if ri == nil {
		snag.Panic("未找到骑手")
	}
	return &model.RiderSampleInfo{
		ID:    ri.ID,
		Name:  ri.Name,
		Phone: ri.Phone,
	}
}
