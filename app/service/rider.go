// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/baidu"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/rs/xid"
    "time"
)

type riderService struct {
}

func NewRider() *riderService {
    return &riderService{}
}

// GetRiderById 根据ID获取骑手及其实名状态
func (r *riderService) GetRiderById(id uint64) (u *ent.Rider, err error) {
    return ar.Ent.Rider.
        QueryNotDeleted().
        WithPerson().
        Where(rider.ID(id)).
        Only(context.Background())
}

// IsAuthed 是否已认证
func (r *riderService) IsAuthed(u *ent.Rider) bool {
    return u.Edges.Person != nil && !model.PersonAuthStatus(u.Edges.Person.Status).RequireAuth()
}

// IsBlocked 骑手是否被封锁
func (r *riderService) IsBlocked(u *ent.Rider) bool {
    p := u.Edges.Person
    return p != nil && p.Block
}

// IsNewDevice 检查是否是新设备
func (r *riderService) IsNewDevice(u *ent.Rider, device *app.Device) bool {
    return u.LastDevice != device.Serial || u.IsNewDevice
}

// Signin 骑手登录
func (r *riderService) Signin(phone string, device *app.Device) (res *model.RiderSigninRes, err error) {
    ctx := context.Background()
    orm := ar.Ent.Rider
    var u *ent.Rider
    u, err = orm.Query().Where(rider.Phone(phone)).WithPerson().Only(ctx)
    if err != nil {
        // 创建骑手
        u, err = orm.Create().
            SetPhone(phone).
            SetLastDevice(device.Serial).
            SetDeviceType(device.Type.Raw()).
            Save(ctx)
        if err != nil {
            return
        }
    }

    // 判定用户是否被封禁
    if r.IsBlocked(u) {
        err = errors.New(ar.RiderBlockedMessage)
        return
    }

    token := xid.New().String() + utils.RandTokenString()
    cache := ar.Cache
    key := fmt.Sprintf("RIDER_%d", u.ID)

    // 删除旧的token
    if old := cache.Get(ctx, key).Val(); old != "" {
        cache.Del(ctx, key)
        cache.Del(ctx, old)
    }

    // 更新设备
    if u.LastDevice != device.Serial {
        r.SetNewDevice(u, device)
    }

    res = &model.RiderSigninRes{
        Id:              u.ID,
        Token:           token,
        IsNewDevice:     r.IsNewDevice(u, device),
        IsContactFilled: u.Contact != nil,
        IsAuthed:        r.IsAuthed(u),
        Contact:         u.Contact,
        Qrcode:          fmt.Sprintf("https://rider.auroraride.com/%d", u.ID),
    }

    // 设置登录token
    cache.Set(ctx, key, token, 7*24*time.Hour)
    cache.Set(ctx, token, u.ID, 7*24*time.Hour)

    return
}

// Signout 强制登出
func (r *riderService) Signout(u *ent.Rider) {
    cache := ar.Cache
    ctx := context.Background()
    key := fmt.Sprintf("RIDER_%d", u.ID)
    token := cache.Get(ctx, key).Val()
    cache.Del(ctx, key)
    cache.Del(ctx, token)
}

// SetNewDevice 更新用户设备
func (r *riderService) SetNewDevice(u *ent.Rider, device *app.Device) {
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
func (r *riderService) GetFaceAuthUrl(c *app.RiderContext) string {
    uri, token := baidu.New().GetAuthenticatorUrl()
    ar.Cache.Set(context.Background(), token, r.GeneratePrivacy(c), 30*time.Minute)
    return uri
}

// GetFaceUrl 获取人脸校验URL
func (r *riderService) GetFaceUrl(c *app.RiderContext) string {
    p := c.Rider.Edges.Person
    uri, token := baidu.New().GetFaceUrl(p.Name, p.IDCardNumber)
    ar.Cache.Set(context.Background(), token, r.GeneratePrivacy(c), 30*time.Minute)
    return uri
}

// FaceAuthResult 获取并更新人脸实名验证结果
func (r *riderService) FaceAuthResult(c *app.RiderContext) (success bool, err error) {
    if !r.ComparePrivacy(c) {
        return false, errors.New("验证失败")
    }
    token := c.Param("token")
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
    blocked, _ := ar.Ent.Person.Query().Where(person.IDCardNumber(detail.IdCardNumber), person.Block(true)).Exist(context.Background())
    if blocked {
        panic(errors.New(""))
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
        SetStatus(model.PersonAuthSuccess.Raw()).
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
func (r *riderService) FaceResult(c *app.RiderContext) (success bool, err error) {
    if !r.ComparePrivacy(c) {
        return false, errors.New("验证失败")
    }
    token := c.Param("token")
    u := c.Rider
    res, resErr := baidu.New().FaceResult(token)
    err = resErr
    if err != nil {
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
func (r *riderService) UpdateContact(u *ent.Rider, contact *model.RiderContact) {
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
func (r *riderService) GeneratePrivacy(c *app.RiderContext) string {
    return fmt.Sprintf("%s-%d", c.Device.Serial, c.Rider.ID)
}

// ComparePrivacy 比对实名认证或人脸识别限制条件是否满足
func (r *riderService) ComparePrivacy(c *app.RiderContext) bool {
    return ar.Cache.Get(context.Background(), c.Param("token")).Val() == r.GeneratePrivacy(c)
}

// ExtendTokenTime 延长骑手登录有效期
func (r *riderService) ExtendTokenTime(id uint64, token string) {
    key := fmt.Sprintf("RIDER_%d", id)
    cache := ar.Cache
    ctx := context.Background()
    cache.Set(ctx, key, token, 7*24*time.Hour)
    cache.Set(ctx, token, id, 7*24*time.Hour)
    _ = ar.Ent.Rider.
        UpdateOneID(id).
        SetLastSigninAt(time.Now()).
        Exec(context.Background())
}
