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
    log "github.com/sirupsen/logrus"
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
    var m *ent.Rider
    m, err = orm.Query().Where(rider.Phone(phone)).Only(ctx)
    if err != nil {
        // 创建骑手
        m, err = orm.Create().SetPhone(phone).Save(ctx)
        if err != nil {
            return
        }
    }

    r.SetDevice(m, device)
    // 判定是否更换了设备
    isNewDevice := r.IsNewDevice(m, device)
    var token string

    if !isNewDevice {
        token = utils.RandTokenString()
        cache := ar.Cache
        key := fmt.Sprintf("RIDER_%d", m.ID)

        // 删除旧的token
        if old := cache.Get(ctx, key).Val(); old != "" {
            cache.Del(ctx, key)
            cache.Del(ctx, old)
        }

        // 设置登录token
        cache.Set(ctx, key, token, 7*24*time.Hour)
        cache.Set(ctx, token, m.ID, 7*24*time.Hour)

        // 更新设备信息
        r.SetDevice(m, device)
    }

    res = &model.RiderSigninRes{
        Id:          m.ID,
        Token:       token,
        IsNewDevice: isNewDevice,
    }
    return
}

// SetDevice 设置用户设备
func (r *riderService) SetDevice(m *ent.Rider, device *app.Device) {
    key := fmt.Sprintf("DEVICE_%d", m.ID)
    err := ar.Cache.Set(context.Background(), key, device, -1).Err()
    log.Println(err)
}

// IsNewDevice 检查是否是新设备
func (r *riderService) IsNewDevice(m *ent.Rider, device *app.Device) bool {
    cache := ar.Cache
    ctx := context.Background()
    key := fmt.Sprintf("DEVICE_%d", m.ID)
    d := new(app.Device)
    err := cache.Get(ctx, key).Scan(d)
    if err != nil {
        log.Errorf("获取用户 %d 设备失败: %v", m.ID, err)
        return true
    }
    return d.Serial != device.Serial
}
