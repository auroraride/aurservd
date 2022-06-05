// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/baidu"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    jsoniter "github.com/json-iterator/go"
    "github.com/rs/xid"
    log "github.com/sirupsen/logrus"
    "strconv"
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
        orm:            ar.Ent.Rider,
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
    return ar.Ent.Rider.
        QueryNotDeleted().
        WithPerson().
        WithEnterprise().
        Where(rider.ID(id)).
        Only(context.Background())
}

// IsAuthed 是否已认证
func (s *riderService) IsAuthed(u *ent.Rider) bool {
    return u.Edges.Person != nil && !model.PersonAuthStatus(u.Edges.Person.Status).RequireAuth()
}

// IsBanned 骑手是否被封禁
func (s *riderService) IsBanned(u *ent.Rider) bool {
    p := u.Edges.Person
    return p != nil && p.Banned
}

// IsNewDevice 检查是否是新设备
func (s *riderService) IsNewDevice(u *ent.Rider, device *model.Device) bool {
    return u.LastDevice != device.Serial || u.IsNewDevice
}

// Signin 骑手登录
func (s *riderService) Signin(device *model.Device, req *model.RiderSignupReq) (res *model.RiderSigninRes, err error) {
    ctx := context.Background()
    orm := ar.Ent.Rider
    var u *ent.Rider

    u, err = orm.QueryNotDeleted().Where(rider.Phone(req.Phone)).WithPerson().Only(ctx)
    if err != nil {
        // 创建骑手
        u, err = orm.Create().
            SetPhone(req.Phone).
            SetLastDevice(device.Serial).
            SetDeviceType(device.Type.Raw()).
            Save(ctx)
        if err != nil {
            return
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
    _, err := ar.Ent.Rider.
        UpdateOneID(u.ID).
        SetLastDevice(device.Serial).
        SetDeviceType(device.Type.Raw()).
        SetIsNewDevice(true).
        Save(context.Background())
    if err != nil {
        snag.Panic(err)
    }
    u.IsNewDevice = true
}

// GetFaceAuthUrl 获取实名验证URL
func (s *riderService) GetFaceAuthUrl(c *app.RiderContext) string {
    uri, token := baidu.New().GetAuthenticatorUrl()
    cache.Set(context.Background(), token, s.GeneratePrivacy(c), 30*time.Minute)
    return uri
}

// GetFaceUrl 获取人脸校验URL
func (s *riderService) GetFaceUrl(c *app.RiderContext) string {
    p := c.Rider.Edges.Person
    uri, token := baidu.New().GetFaceUrl(p.Name, p.IDCardNumber)
    cache.Set(context.Background(), token, s.GeneratePrivacy(c), 30*time.Minute)
    return uri
}

// FaceAuthResult 获取并更新人脸实名验证结果
func (s *riderService) FaceAuthResult(c *app.RiderContext, token string) (success bool) {
    if !s.ComparePrivacy(c) {
        snag.Panic("验证失败")
    }
    u := c.Rider
    data, err := baidu.New().AuthenticatorResult(token)
    if err != nil {
        return
    }
    success = data.Success
    if !success {
        return
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
    // 判断用户是否被封禁
    banned, _ := ar.Ent.Person.QueryNotDeleted().Where(person.IDCardNumber(detail.IdCardNumber), person.Banned(true)).Exist(context.Background())
    if banned {
        snag.Panic(snag.StatusForbidden, ar.BannedMessage)
    }

    // 上传图片到七牛云
    oss := ali.NewOss()
    prefix := fmt.Sprintf("%s-%s/", res.IdcardOcrResult.Name, res.IdcardOcrResult.IdCardNumber)
    fm := oss.UploadUrlFile(prefix+"face.jpg", res.FaceImg)
    pm := oss.UploadBase64ImageJpeg(prefix+"portrait.jpg", res.IdcardImages.FrontBase64)
    nm := oss.UploadBase64ImageJpeg(prefix+"national.jpg", res.IdcardImages.BackBase64)

    icNum := vr.IdCardNumber
    id, err := ar.Ent.Person.
        Create().
        SetStatus(model.PersonAuthenticated.Raw()).
        SetIDCardNumber(icNum).
        SetName(vr.Name).
        SetAuthFace(fm).
        SetIDCardNational(nm).
        SetIDCardPortrait(pm).
        SetAuthResult(vr).
        SetAuthAt(time.Now()).
        OnConflictColumns(person.FieldIDCardNumber).
        UpdateNewValues().
        ID(context.Background())
    if err != nil {
        snag.Panic(err)
    }

    err = ar.Ent.Rider.
        UpdateOneID(u.ID).
        SetPersonID(id).
        SetLastFace(fm).
        SetIsNewDevice(false).
        Exec(context.Background())
    if err != nil {
        snag.Panic(err)
    }
    return
}

// FaceResult 获取人脸比对结果
func (s *riderService) FaceResult(c *app.RiderContext, token string) (success bool) {
    if !s.ComparePrivacy(c) {
        snag.Panic("验证失败")
    }
    u := c.Rider
    res, err := baidu.New().FaceResult(token)
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
    err = ar.Ent.Rider.
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
    err := ar.Ent.Rider.UpdateOneID(u.ID).SetContact(contact).Exec(context.Background())
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
    _ = ar.Ent.Rider.
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

// List 骑手列表
func (s *riderService) List(req *model.RiderListReq) *model.PaginationRes {
    q := ar.Ent.Rider.
        Query().
        WithPerson().
        WithOrders(func(oq *ent.OrderQuery) {
            oq.Where(
                order.Type(model.OrderTypeDeposit),
                order.Status(model.OrderStatusPaid),
            )
        }).
        WithSubscribes(func(sq *ent.SubscribeQuery) {
            sq.Order(ent.Desc(subscribe.FieldCreatedAt)).Limit(1)
        })
    if req.Keyword != nil {
        // 判定是否id字段
        if id, err := strconv.ParseUint(*req.Keyword, 10, 64); err == nil && id > 0 {
            q.Where(rider.ID(id))
        } else {
            q.Where(
                rider.Or(
                    rider.HasPersonWith(person.NameContainsFold(*req.Keyword)),
                    rider.PushIDContainsFold(*req.Keyword),
                ),
            )
        }
    }
    if req.Start != nil {
        start := carbon.ParseByLayout(*req.Start, carbon.DateLayout)
        if start.Error != nil {
            snag.Panic("日期格式错误")
        }
        q.Where(rider.CreatedAtGTE(start.Carbon2Time()))
    }
    if req.End != nil {
        end := carbon.ParseByLayout(*req.End, carbon.DateLayout)
        if end.Error != nil {
            snag.Panic("日期格式错误")
        }
        end.AddDay()
        q.Where(rider.CreatedAtLT(end.Carbon2Time()))
    }
    if req.Modified != nil {
        m := *req.Modified
        if m {
            q.Where(rider.LastModifierNotNil())
        } else {
            q.Where(rider.LastModifierIsNil())
        }
    }
    if req.Status != nil {
        rs := *req.Status
        switch rs {
        case model.RiderStatusNormal:
            q.Where(
                rider.Blocked(false),
                rider.Or(
                    rider.PersonIDIsNil(),
                    rider.HasPersonWith(person.Banned(false)),
                ),
            )
            break
        case model.RiderStatusBlocked:
            q.Where(rider.Blocked(true))
            break
        case model.RiderStatusBanned:
            q.Where(rider.HasPersonWith(person.Banned(true)))
            break
        }
    }
    if req.AuthStatus != nil {
        ra := *req.AuthStatus
        switch ra {
        case model.PersonUnauthenticated:
            q.Where(
                rider.Or(
                    rider.PersonIDIsNil(),
                    rider.HasPersonWith(person.Status(model.PersonUnauthenticated.Raw())),
                ),
            )
            break
        default:
            q.Where(rider.HasPersonWith(person.Status(ra.Raw())))
            break
        }
    }
    if req.SubscribeStatus != nil {
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
        default:
            q.Where(rider.HasSubscribesWith(subscribe.Status(rss)))
            break
        }
    }

    return model.ParsePaginationResponse[model.RiderItem, ent.Rider](
        q,
        req.PaginationReq,
        func(item *ent.Rider) model.RiderItem {
            p := item.Edges.Person
            ri := model.RiderItem{
                ID:         item.ID,
                Enterprise: nil,
                Phone:      item.Phone,
                Status:     model.RiderStatusNormal,
                AuthStatus: model.PersonUnauthenticated,
            }
            if item.Blocked {
                ri.Status = model.RiderStatusBlocked
            }
            if p != nil {
                ri.Name = p.Name
                ri.IDCardNumber = p.IDCardNumber
                ri.AuthStatus = model.PersonAuthStatus(p.Status)
                if p.Banned {
                    ri.Status = model.RiderStatusBanned
                }
                if p.AuthResult != nil {
                    ri.Address = p.AuthResult.Address
                }
            }

            if item.Edges.Orders != nil && len(item.Edges.Orders) > 0 {
                ri.Deposit = item.Edges.Orders[0].Amount
            }
            if item.Edges.Subscribes != nil {
                sub := item.Edges.Subscribes[0]
                ri.Subscribe = &model.RiderItemSubscribe{
                    Status:    sub.Status,
                    Remaining: sub.Remaining,
                    Voltage:   sub.Voltage,
                }
                if sub.Status == model.SubscribeStatusUsing && sub.Remaining <= 3 {
                    ri.Subscribe.Status = 11
                }
            }
            if item.DeletedAt != nil {
                ri.DeletedAt = item.DeletedAt.Format(carbon.DateTimeLayout)
                ri.Remark = item.Remark
            }
            return ri
        },
    )
}

func (s *riderService) ModifyUserPlanDays() {

}

func (s *riderService) Query(id uint64) *ent.Rider {
    item, err := ar.Ent.Rider.QueryNotDeleted().Where(rider.ID(id)).Only(s.ctx)
    if err != nil || item == nil {
        snag.Panic("未找到骑手")
    }
    return item
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
        ol.SetOperate(logging.OperateRiderBLock).SetDiff(nb, bd)
    } else {
        ol.SetOperate(logging.OperateRiderUnBLock).SetDiff(bd, nb)
    }
    ol.Send()
}

// DepositOrder 获取骑手押金订单
func (s *riderService) DepositOrder(riderID uint64) *ent.Order {
    o, err := ar.Ent.Order.QueryNotDeleted().Where(
        order.RiderID(riderID),
        order.Status(model.OrderStatusPaid),
        order.Type(model.OrderTypeDeposit),
    ).First(s.ctx)
    if err != nil {
        log.Error(err)
    }
    return o
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

// Profile 获取用户资料
func (s *riderService) Profile(u *ent.Rider, device *model.Device, token string) *model.RiderSigninRes {
    sub := NewSubscribe().Recent(u.ID)
    profile := &model.RiderSigninRes{
        ID:              u.ID,
        IsNewDevice:     s.IsNewDevice(u, device),
        IsContactFilled: u.Contact != nil,
        IsAuthed:        s.IsAuthed(u),
        Contact:         u.Contact,
        Qrcode:          fmt.Sprintf("https://rider.auroraride.com/%d", u.ID),
        Deposit:         s.Deposit(u.ID),
        Subscribe:       sub,
        OrderNotActived: sub != nil && sub.Status == model.SubscribeStatusInactive,
        Token:           token,
    }
    return profile
}

// GetLogs 获取用户操作日志
func (s *riderService) GetLogs(req *model.RiderLogReq) (items []model.LogOperate) {
    u := s.Query(req.ID)
    b, _ := jsoniter.Marshal(logging.NewOperateLog().GetLogs(u.CreatedAt, fmt.Sprintf(`refTable = 'rider' AND refId = %d`, u.ID), req.Offset))
    _ = jsoniter.Unmarshal(b, &items)
    return
}
