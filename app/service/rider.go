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
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/baidu"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/ent/enterprise"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/predicate"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    jsoniter "github.com/json-iterator/go"
    "github.com/rs/xid"
    "strconv"
    "strings"
    "time"
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
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewRiderWithRider(u *ent.Rider) *riderService {
    s := NewRider()
    s.ctx = context.WithValue(s.ctx, "rider", u)
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
        Where(rider.ID(id)).
        Only(context.Background())
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
    ctx := context.Background()
    orm := ent.Database.Rider
    var u *ent.Rider
    var err error

    u, err = orm.QueryNotDeleted().Where(rider.Phone(req.Phone)).WithPerson().WithEnterprise().Only(ctx)
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
    isNew := true
    if ar.Config.App.Debug.Phone[u.Phone] {
        isNew = false
    }
    _, err := ent.Database.Rider.
        UpdateOneID(u.ID).
        SetLastDevice(device.Serial).
        SetDeviceType(device.Type.Value()).
        SetIsNewDevice(isNew).
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
        ID(context.Background())
    if err != nil {
        snag.Panic(err)
    }

    if success || u.PersonID == nil {
        // 判断ID是否等于实名认证的ID, 如果不是, 则删除
        if u.PersonID != nil && *u.PersonID != id {
            _ = ent.Database.Person.DeleteOneID(*u.PersonID).Exec(s.ctx)
        }
        err = ent.Database.Rider.
            UpdateOneID(u.ID).
            SetPersonID(id).
            SetLastFace(fm).
            SetIsNewDevice(false).
            Exec(context.Background())
        if err != nil {
            snag.Panic(err)
        }
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
func (*riderService) GetRiderSampleInfo(rider *ent.Rider) model.RiderSampleInfo {
    return model.RiderSampleInfo{
        ID:    rider.ID,
        Name:  rider.Edges.Person.Name,
        Phone: rider.Phone,
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
            sq.WithCity().Order(ent.Desc(subscribe.FieldCreatedAt))
        }).
        WithContracts(func(cq *ent.ContractQuery) {
            cq.Where(contract.DeletedAtIsNil(), contract.Status(model.ContractStatusSuccess.Value()))
        }).
        WithEnterprise().
        Order(ent.Desc(rider.FieldCreatedAt))
    if req.Keyword != nil {
        info["关键词"] = *req.Keyword
        // 判定是否id字段
        q.Where(
            rider.Or(
                rider.HasPersonWith(
                    person.Or(
                        person.NameContainsFold(*req.Keyword),
                        person.IDCardNumberContainsFold(*req.Keyword),
                    ),
                ),
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
            break
        case model.RiderStatusBlocked:
            info["用户状态"] = "已禁用"
            q.Where(rider.Blocked(true))
            break
        case model.RiderStatusBanned:
            info["用户状态"] = "已封禁"
            q.Where(rider.HasPersonWith(person.Banned(true)))
            break
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
            break
        default:
            q.Where(rider.HasPersonWith(person.Status(ra.Value())))
            break
        }
    }
    if req.SubscribeStatus != nil {
        key := "业务状态"
        rss := *req.SubscribeStatus
        switch rss {
        case 11:
            // 即将到期
            q.Where(rider.HasSubscribesWith(
                subscribe.Status(model.SubscribeStatusUsing),
                subscribe.RemainingLTE(3),
            ))
            break
        case 99:
            q.Where(rider.Not(rider.HasSubscribes()))
            break
        case model.SubscribeStatusUnSubscribed:
            q.Where(
                rider.And(
                    rider.HasSubscribesWith(
                        subscribe.EndAtNotNil(),
                        subscribe.Status(model.SubscribeStatusUnSubscribed),
                    ),
                    rider.Not(rider.HasSubscribesWith(subscribe.StatusIn(model.SubscribeNotUnSubscribed()...))),
                ),
            )
        default:
            q.Where(rider.HasSubscribesWith(subscribe.Status(rss)))
            break
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
        }[*req.SubscribeStatus]
    }
    if req.PlanID != nil {
        info["骑士卡"] = ent.NewExportInfo(*req.PlanID, plan.Table)
        q.Where(rider.HasSubscribesWith(subscribe.PlanID(*req.PlanID)))
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
        q.Where(rider.HasSubscribesWith(subscribe.CityID(*req.CityID)))
    }

    if req.Remaining != nil {
        arr := strings.Split(*req.Remaining, ",")

        var subqs []predicate.Subscribe
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

        q.Where(rider.HasSubscribesWith(subqs...))
    }

    if req.Suspend != nil {
        if *req.Suspend {
            q.Where(rider.HasSubscribesWith(subscribe.SuspendAtNotNil()))
            info["暂停扣费"] = "是"
        } else {
            q.Where(rider.HasSubscribesWith(subscribe.SuspendAtIsNil()))
            info["暂停扣费"] = "否"
        }
    }
    return
}

// List 骑手列表
func (s *riderService) List(req *model.RiderListReq) *model.PaginationRes {
    q, _ := s.listFilter(req.RiderListFilter)

    return model.ParsePaginationResponse[model.RiderItem, ent.Rider](
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
        ID:         item.ID,
        Phone:      item.Phone,
        Status:     model.RiderStatusNormal,
        AuthStatus: model.PersonUnauthenticated,
        Contact:    item.Contact,
    }
    e := item.Edges.Enterprise
    if e != nil {
        ri.Enterprise = &model.Enterprise{
            ID:    e.ID,
            Name:  e.Name,
            Agent: e.Agent,
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
    if contracts != nil && len(contracts) > 0 {
        ri.Contract = contracts[0].Files[0]
    }

    if item.Edges.Orders != nil && len(item.Edges.Orders) > 0 {
        ri.Deposit = item.Edges.Orders[0].Amount
    }
    if item.Edges.Subscribes != nil && len(item.Edges.Subscribes) > 0 {
        sub := item.Edges.Subscribes[0]
        ri.Subscribe = &model.RiderItemSubscribe{
            ID:        sub.ID,
            Status:    sub.Status,
            Remaining: sub.Remaining,
            Model:     sub.Model,
            Suspend:   sub.SuspendAt != nil,
            Formula:   sub.Formula,
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
    }
    if item.DeletedAt != nil {
        ri.DeletedAt = item.DeletedAt.Format(carbon.DateTimeLayout)
        ri.Remark = item.Remark
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
            "骑手",     // 1
            "电话",     // 2
            "证件",     // 3
            "户籍",     // 4
            "企业",     // 5
            "押金",     // 6
            "订阅",     // 7
            "暂停",     // 8
            "电池",     // 9
            "剩余",     // 10
            "状态",     // 11
            "认证",     // 12
            "紧急联系", // 13
            "注册时间", // 14
        }
        rows = append(rows, title)
        for _, item := range items {
            detail := s.detailRiderItem(item)
            row := []any{
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
                []string{"正常", "正常", "禁用", "黑名单"}[detail.Status],
                detail.AuthStatus.String(),
                "",
                item.CreatedAt.Format(carbon.DateTimeLayout),
            }
            if detail.City != nil {
                row[0] = detail.City.Name
            }
            if detail.Person != nil {
                row[3] = detail.Person.IDCardNumber
            }
            if detail.Enterprise != nil {
                row[5] = detail.Enterprise.Name
            }
            if detail.Subscribe != nil {
                row[7] = model.SubscribeStatusText(detail.Subscribe.Status)
                if detail.Subscribe.Suspend {
                    row[8] = "是"
                }
                row[9] = detail.Subscribe.Model
                row[10] = detail.Subscribe.Remaining
            }
            if item.Contact != nil {
                row[13] = item.Contact.String()
            }
            rows = append(rows, row)
        }
        tools.NewExcel(path).AddValues(rows).Done()
    })
}

func (s *riderService) Query(id uint64) *ent.Rider {
    item, err := ent.Database.Rider.QueryNotDeleted().Where(rider.ID(id)).WithPerson().Only(s.ctx)
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
    f, _ := cache.Get(s.ctx, model.SettingDeposit).Float64()
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
    subd, _ := NewSubscribe().RecentDetail(u.ID)
    profile := &model.RiderSigninRes{
        ID:              u.ID,
        Phone:           u.Phone,
        IsNewDevice:     s.IsNewDevice(u, device),
        IsContactFilled: u.Contact != nil,
        IsAuthed:        s.IsAuthed(u),
        Contact:         u.Contact,
        Qrcode:          s.GetQrcode(u.ID),
        Token:           token,
        Subscribe:       subd,
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
    } else {
        profile.Subscribe = subd
        profile.OrderNotActived = tools.NewPointer().Bool(subd != nil && subd.Status == model.SubscribeStatusInactive)
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
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
        subscribe.RiderID(req.ID),
        subscribe.StatusNotIn(model.SubscribeStatusUnSubscribed, model.SubscribeStatusCanceled),
    ).First(s.ctx)
    if sub != nil {
        snag.Panic("骑手当前有订阅")
    }
    _, err := s.orm.SoftDeleteOneID(req.ID).Save(s.ctx)
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
