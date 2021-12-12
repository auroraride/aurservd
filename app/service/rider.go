// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/utils"
    "time"
)

type riderService struct {
}

func NewRider() *riderService {
    return &riderService{}
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

    // 判定是否更换了设备
    var token string

    token = utils.RandTokenString()
    cache := ar.Cache
    key := fmt.Sprintf("RIDER_%d", u.ID)

    // 删除旧的token
    if old := cache.Get(ctx, key).Val(); old != "" {
        cache.Del(ctx, key)
        cache.Del(ctx, old)
    }

    res = &model.RiderSigninRes{
        Id:              u.ID,
        IsNewDevice:     false,
        Token:           token,
        TokenPermission: r.GetTokenPermission(u, device),
    }

    // 设置登录token
    cache.Set(ctx, key, token, 7*24*time.Hour)
    cache.Set(ctx, token, u.ID, 7*24*time.Hour)

    return
}

// GetRiderById 根据ID获取骑手及其实名状态
func (r *riderService) GetRiderById(id uint64) (u *ent.Rider, err error) {
    return ar.Ent.Rider.Query().
        WithPerson().
        Where(rider.ID(id)).
        Only(context.Background())
}

// GetTokenPermission 获取token权限
func (r *riderService) GetTokenPermission(u *ent.Rider, device *app.Device) model.RiderTokenPermission {
    perm := model.RiderTokenPermissionCommon
    // 判断是否实名
    p := u.Edges.Person
    if p == nil || model.PersonAuthStatus(p.Status).RequireAuth() {
        perm = model.RiderTokenPermissionAuth
    } else if r.IsNewDevice(u, device) {
        // 判断是否新设备
        perm = model.RiderTokenPermissionNewDevice
    }
    return perm
}

// GetTokenPermissionResponse 获取token权限响应
func (r *riderService) GetTokenPermissionResponse(perm model.RiderTokenPermission) (status int, message string) {
    switch perm {
    case model.RiderTokenPermissionAuth:
        status = response.StatusNotAcceptable
        message = "需要实名认证"
    case model.RiderTokenPermissionNewDevice:
        status = response.StatusLocked
        message = "需要人脸验证"
    }
    return
}

// SetDevice 更新用户设备
func (r *riderService) SetDevice(u *ent.Rider, device *app.Device) error {
    orm := ar.Ent.Rider
    _, err := orm.UpdateOneID(u.ID).
        SetLastDevice(device.Serial).
        SetDeviceType(device.Type.Raw()).
        Save(context.Background())
    return err
}

// IsNewDevice 检查是否是新设备
func (r *riderService) IsNewDevice(u *ent.Rider, device *app.Device) bool {
    return u.LastDevice == device.Serial
}
