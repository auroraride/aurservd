// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
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
    var m *ent.Rider
    m, err = orm.Query().Where(rider.Phone(phone)).Only(ctx)
    if err != nil {
        // 创建骑手
        m, err = orm.Create().SetPhone(phone).SetDeviceSn(device.Sn).SetDeviceType(device.Type.Raw()).Save(ctx)
        if err != nil {
            return
        }
    }
    token := utils.RandTokenString()

    // 设置登录token
    ar.Cache.Set(ctx, "RIDER_"+phone, token, 7*24*time.Hour)
    ar.Cache.Set(ctx, token, m.ID, 7*24*time.Hour)

    // 判定是否更换了设备
    res = &model.RiderSigninRes{
        Id:          m.ID,
        Token:       token,
        IsNewDevice: m.DeviceSn == device.Sn,
    }
    return
}
