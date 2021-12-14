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
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/baidu"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
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
    return u.Edges.Person != nil && model.PersonAuthStatus(u.Edges.Person.Status).RequireAuth()
}

// IsBlocked 骑手是否被封锁
func (r *riderService) IsBlocked(u *ent.Rider) bool {
    p := u.Edges.Person
    return p != nil && p.Block
}

// IsNewDevice 检查是否是新设备
func (r *riderService) IsNewDevice(u *ent.Rider, device *app.Device) bool {
    return u.LastDevice != device.Serial
}

// Signin 骑手登录
// 当返回字段 isNewDevice 时候需要跳转人脸识别界面
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
        err = errors.New("你已被封禁")
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

    res = &model.RiderSigninRes{
        Id:              u.ID,
        Token:           token,
        IsNewDevice:     r.IsNewDevice(u, device),
        IsContactFilled: u.Contact != nil,
        Contact:         u.Contact,
        IsAuthed:        r.IsAuthed(u),
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

// SetDevice 更新用户设备
func (r *riderService) SetDevice(u *ent.Rider, device *app.Device) error {
    _, err := ar.Ent.Rider.
        UpdateOneID(u.ID).
        SetLastDevice(device.Serial).
        SetDeviceType(device.Type.Raw()).
        Save(context.Background())
    return err
}

// FaceAuthResult 获取并更新人脸实名验证结果
func (r *riderService) FaceAuthResult(u *ent.Rider, token string) (success bool, err error) {
    data, err := baidu.New().AuthenticatorResult(token)
    if err != nil {
        return
    }
    success = data.Success
    if !success {
        return
    }
    // 判断用户是否被 Block

    res := data.Result
    oss := ali.NewOss()
    prefix := fmt.Sprintf("%s-%s/", res.IdcardOcrResult.Name, res.IdcardOcrResult.IdCardNumber)
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
    fm := oss.UploadUrlFile(prefix+"face.jpg", res.FaceImg)
    pm := oss.UploadBase64ImageJpeg(prefix+"portrait.jpg", res.IdcardImages.FrontBase64)
    nm := oss.UploadBase64ImageJpeg(prefix+"national.jpg", res.IdcardImages.BackBase64)

    icNum := vr.IdCardNumber
    id, err := ar.Ent.Person.Create().
        SetStatus(model.PersonAuthSuccess.Raw()).
        SetIcNumber(icNum).
        SetName(vr.Name).
        SetFaceImg(fm).
        SetIcNational(nm).
        SetIcPortrait(pm).
        SetFaceVerifyResult(vr).
        OnConflictColumns(person.FieldIcNumber).
        UpdateNewValues().
        ID(context.Background())
    if err != nil {
        panic(response.NewError(err))
    }

    err = ar.Ent.Rider.
        UpdateOneID(u.ID).
        SetPersonID(id).
        SetLastFace(fm).
        Exec(context.Background())
    if err != nil {
        panic(response.NewError(err))
    }
    return
}

// FaceResult 获取人脸比对结果
func (r *riderService) FaceResult(u *ent.Rider, token, sn string) (success bool, err error) {
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
    fm := ali.NewOss().UploadUrlFile(fmt.Sprintf("%s-%s/face-%s.jpg", p.Name, p.IcNumber, sn), res.Result.Image)
    err = ar.Ent.Rider.
        UpdateOneID(u.ID).
        SetLastFace(fm).
        SetLastDevice(sn).
        Exec(context.Background())
    if err != nil {
        panic(response.NewError(err))
    }
    return
}

// UpdateContact 更新紧急联系人
func (r *riderService) UpdateContact(u *ent.Rider, contact *model.RiderContact) {
    // 判断紧急联系人手机号是否和当前骑手手机号一样
    if u.Phone == contact.Phone {
        panic(response.NewError("紧急联系人手机号不能是当前手机号"))
    }
    err := ar.Ent.Rider.UpdateOneID(u.ID).SetContact(contact).Exec(context.Background())
    if err != nil {
        panic(response.NewError(err))
    }
}
