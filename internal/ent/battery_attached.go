// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-08
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "context"
    "github.com/auroraride/aurservd/internal/ent/battery"
    "github.com/jackc/pgconn"
)

// Allocate 将电池分配给骑手
func (buo *BatteryUpdateOne) Allocate(sub *Subscribe) (err error) {
    err = buo.SetSubscribeID(sub.ID).SetRiderID(sub.RiderID).Exec(context.Background())
    if err != nil && IsConstraintError(err) {
        switch v := err.(*ConstraintError).wrap.(type) {
        case *pgconn.PgError:
            // battery_subscribe_id_key
            if v.ConstraintName == "battery_subscribe_id_key" || v.ConstraintName == "battery_rider_id_key" {
                // TODO 删除原有信息
                ctx := context.Background()
                err = NewBatteryClient(buo.config).Update().Where(battery.SubscribeID(sub.ID)).ClearRiderID().ClearSubscribeID().Exec(ctx)
                if err != nil {
                    return
                }
                return buo.SetSubscribeID(sub.ID).SetRiderID(sub.RiderID).Exec(context.Background())
            }
        }
    }

    return
}

// Unallocate 清除骑手信息
func (buo *BatteryUpdateOne) Unallocate() error {
    return buo.ClearSubscribeID().ClearRiderID().Exec(context.Background())
}
