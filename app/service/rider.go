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

    // 设置登录token
    cache.Set(ctx, key, token, 7*24*time.Hour)
    cache.Set(ctx, token, u.ID, 7*24*time.Hour)

    res = &model.RiderSigninRes{
        Id:          u.ID,
        IsNewDevice: false,
        Token:       token,
    }

    // 判断是否实名
    p := u.Edges.Person
    if p == nil || model.PersonAuthStatus(p.Status).RequireAuth() {
        res.TokenPermission = model.TokenPermissionAuth
        return
    }

    // 判断是否新设备
    if r.IsNewDevice(u, device) {
        res.TokenPermission = model.TokenPermissionNewDevice
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
